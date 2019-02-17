package common

import (
	"errors"
	"fmt"
	"gitlocal/gome"
	"strings"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Color struct {
	A, R, G, B float32
}

/*
	RenderComponent
*/

// A RenderComponent is a component used to render a texture of its
// entity
type RenderComponent struct {
	Stride       int32
	Vertices     []float32
	Indices      []uint32
	ModelUpdated bool
	Color
}

func (rc *RenderComponent) Name() string { return "Render" }

/*
	RenderSystem
*/

const vertexShaderSrc string = `
#version 130

uniform vec3 u_Size;
uniform vec3 u_Position;

in vec3 vPos;

out vec4 vColor;

void main() {
    gl_Position = vec4(vPos * u_Size + u_Position, 1.0);
}
` + "\x00"

const fragmentShaderSrc string = `
#version 130

uniform vec4 u_Color;

out vec4 fColor;

void main() {
    fColor = u_Color;
}
` + "\x00"

// A glObject contains all relevant pointers for a texture
type glObject struct {

	// Vertex Buffer Object
	vbo uint32

	// Index Buffer Object
	ibo uint32

	// Vertex Array Object
	vao uint32
}

// A RenderSystem renders the texture of its entities
type RenderSystem struct {
	gome.MultiSystem

	shaderProgram uint32

	// glObjects contains pointer to OpenGL Objects.
	// VBO, VAO
	glObjects map[uint]glObject
}

func (*RenderSystem) RequiredComponents() []string { return []string{"Render", "Space"} }

func (rs *RenderSystem) Init(scene *gome.Scene) {
	// initialize the base system
	rs.MultiSystem.Init(scene)

	rs.glObjects = make(map[uint]glObject)

	// initialize OpenGL
	gl.Init()

	// if debug is enabled, show debug output
	if scene.WindowArgs.Debug {
		// opengl version
		fmt.Println("OpenGL Version:", gl.GoStr(gl.GetString(gl.VERSION)))

		// error and debug outptut
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.DebugMessageCallback(func(
			source uint32,
			gltype uint32,
			id uint32,
			severity uint32,
			length int32,
			message string,
			userParam unsafe.Pointer) {

			// warn if it's an error
			errWarning := ""
			if gltype == gl.DEBUG_TYPE_ERROR {
				errWarning = "** ERROR **"
			}

			fmt.Printf("GL CALLBACK: %s type = 0x%x, severity = 0x%x, message = %s\n",
				errWarning, gltype, severity, message)
		}, gl.Ptr(nil))
	}

	// compile the shaders
	if scene.WindowArgs.Debug {
		fmt.Println("Compiling vertex shader...")
	}
	vertexShader, err := compileShader(vertexShaderSrc, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	if scene.WindowArgs.Debug {
		fmt.Println("Compiling fragment shader...")
	}
	fragmentShader, err := compileShader(fragmentShaderSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	// attach shaders
	rs.shaderProgram = gl.CreateProgram()
	gl.AttachShader(rs.shaderProgram, vertexShader)
	gl.AttachShader(rs.shaderProgram, fragmentShader)
	gl.LinkProgram(rs.shaderProgram) // links the shaders

	// check if linking failed
	var status int32
	gl.GetProgramiv(rs.shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(rs.shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		// create empty string that can hold the log content
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(rs.shaderProgram, logLength, nil, gl.Str(log)) // returns the shader compile log

		// display error
		panic("Failed to link shader program:\n" + log)
	}

	// once the program is linked, we don't need the single shaders anymore
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
}

func (rs *RenderSystem) Add(id uint, components []gome.Component) {
	rs.MultiSystem.Add(id, components)

	renderComponent := components[0].(*RenderComponent)

	var VAO uint32
	gl.GenVertexArrays(1, &VAO) // generates the vertex array (or multiple)
	gl.BindVertexArray(VAO)     // binds the vertex array

	// Vertex Buffer Object
	var VBO uint32
	gl.GenBuffers(1, &VBO)              // generates the buffer (or multiple)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO) // tells OpenGL what kind of buffer this is

	// BufferData assigns data to the buffer.
	// there can only be one ARRAY_BUFFER, so OpenGL knows which buffer we mean if we
	// tell it what type of buffer it is.
	//			  type			   size (in bytes)   pointer to data	usage
	gl.BufferData(gl.ARRAY_BUFFER, len(renderComponent.Vertices)*4, gl.Ptr(renderComponent.Vertices), gl.DYNAMIC_DRAW)

	// define an array of generic vertex attribute data
	// index, size, type, normalized, stride of vertex (in bytes), pointer (offset)
	// point positions
	gl.VertexAttribPointer(0, int32(len(renderComponent.Vertices))/renderComponent.Stride,
		gl.FLOAT, false, renderComponent.Stride*4, nil)
	gl.EnableVertexAttribArray(0) // enable the defined vertex attribute pointer

	// Index Buffer Object
	var IBO uint32
	gl.GenBuffers(1, &IBO)                      // generates the buffer (or multiple)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, IBO) // tells OpenGL what kind of buffer this is

	// BufferData assigns data to the buffer.
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(renderComponent.Indices)*4, gl.Ptr(renderComponent.Indices), gl.DYNAMIC_DRAW)

	// unbind
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	// save the pointers for later
	rs.glObjects[id] = glObject{vbo: VBO, ibo: IBO, vao: VAO}
}

func (rs *RenderSystem) Update(delta time.Duration) {
	gl.ClearColor(0, 0, 0, 0)     // set the clear color
	gl.Clear(gl.COLOR_BUFFER_BIT) // apply clear color

	gl.UseProgram(rs.shaderProgram)

	// get uniform locations
	colorLoc := gl.GetUniformLocation(rs.shaderProgram, gl.Str("u_Color\x00"))
	sizeLoc := gl.GetUniformLocation(rs.shaderProgram, gl.Str("u_Size\x00"))
	positionLoc := gl.GetUniformLocation(rs.shaderProgram, gl.Str("u_Position\x00"))

	for id, components := range rs.MultiSystem.Entities {
		glObj := rs.glObjects[id]
		renderComponent := components[0].(*RenderComponent)
		spaceComponent := components[1].(*SpaceComponent)

		// if the model was updated, adjust the buffers accordingly
		if renderComponent.ModelUpdated {

			// VBO //
			// bind the buffer
			gl.BindBuffer(gl.ARRAY_BUFFER, glObj.vbo)

			// write new data into buffer
			gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(renderComponent.Vertices)*4, gl.Ptr(renderComponent.Vertices))

			// IBO //
			// bind the buffer
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, glObj.ibo)

			// write new data into buffer
			gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, len(renderComponent.Indices)*4, gl.Ptr(renderComponent.Indices))

			// reset ModelUpdated flag
			renderComponent.ModelUpdated = false
		}

		// VAO //
		// bind the vertex array
		gl.BindVertexArray(glObj.vao)

		// set color, size, position by uniform
		// color
		color := renderComponent.Color
		gl.Uniform4f(colorLoc, color.R, color.G, color.B, color.A)
		// size
		size := spaceComponent.Size
		gl.Uniform3f(sizeLoc, size.X, size.Y, size.Z)
		// position
		position := spaceComponent.Position
		gl.Uniform3f(positionLoc, position.X, position.Y, position.Z)

		// draw
		indicesCount := int32(len(renderComponent.Indices))
		gl.DrawElements(gl.TRIANGLES, indicesCount, gl.UNSIGNED_INT, nil)

		// unbind
		gl.BindVertexArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	}
}

// compileShader compiles a shader from source and returns the shader ID and, if occurred,
// an error. The shaderType can be any OpenGL shader type, e.g. gl.VERTEX_SHADER
func compileShader(source string, shaderType uint32) (shader uint32, err error) {
	// create a shader from source (returns shader ID)
	shader = gl.CreateShader(shaderType)
	csource, free := gl.Strs(source) // returns a C String and a function to free the memory
	//				shader, count, source string, length (unused)
	gl.ShaderSource(shader, 1, csource, nil)
	free()                   // frees the memory used by csource
	gl.CompileShader(shader) // compile the shader

	// check if compiling failed
	var status int32
	//			   shader  info type		  pointer
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status) // returns shader info
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		// create empty string that can hold the log content
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log)) // returns the shader compile log

		// set error message
		err = errors.New("Failed to compile OpenGL shader:\n" + log)
	}

	return
}
