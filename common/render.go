package common

import (
	"errors"
	"fmt"
	"gitlocal/gome"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
)

/*
	RenderComponent
*/

// A RenderComponent is a component used to render a texture of its
// entity
type RenderComponent struct {
	Stride   int32
	Vertices []float32
}

func (rc *RenderComponent) Name() string { return "Render" }

// VertexArrayToObj generates a VAO from vertex data and its stride size (dimensions)
func VertexArrayToObj(array []float32, stride int32) (VAO uint32) {

	return
}

/*
	RenderSystem
*/

const fragmentShaderSrc string = `
#version 130

in vec4 vColor;

out vec4 frag_color;

void main() {
    frag_color = vColor;
}
` + "\x00"

const vertexShaderSrc string = `
#version 130

in vec3 vPos;

out vec4 vColor;

void main() {
    gl_Position = vec4(vPos, 1.0);
	vColor = vec4(1.0, vPos);
}
` + "\x00"

// A RenderSystem renders the texture of its entities
type RenderSystem struct {
	gome.MultiSystem

	shaderProgram uint32

	// glObjects contains pointer to OpenGL Objects.
	// VBO, VAO
	glObjects map[uint][]uint32
}

func (*RenderSystem) RequiredComponents() []string { return []string{"Render", "Space"} }

func (rs *RenderSystem) Init(scene *gome.Scene) {
	// initialize the base system
	rs.MultiSystem.Init(scene)

	rs.glObjects = make(map[uint][]uint32)

	// initialize OpenGL
	gl.Init()
	fmt.Println("OpenGL Version:", gl.GoStr(gl.GetString(gl.VERSION)))

	// compile the shaders
	fmt.Println("Compiling vertex shader...")
	vertexShader, err := compileShader(vertexShaderSrc, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fmt.Println("Compiling fragment shader...")
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

	// Vertex Buffer Object
	var VBO uint32
	gl.GenBuffers(1, &VBO)              // generates the buffer (or multiple)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO) // tells OpenGL what kind of buffer this is

	converted := rs.convertVertexArray(id)

	// BufferData assigns data to the buffer.
	// there can only be one ARRAY_BUFFER, so OpenGL knows which buffer we mean if we
	// tell it what type of buffer it is.
	//			  type			   size (in bytes)   pointer to data	usage
	gl.BufferData(gl.ARRAY_BUFFER, len(converted)*4, gl.Ptr(converted), gl.DYNAMIC_DRAW)

	var VAO uint32
	gl.GenVertexArrays(1, &VAO) // generates the vertex array (or multiple)
	gl.BindVertexArray(VAO)     // binds the vertex array

	// define an array of generic vertex attribute data
	// index, size, type, normalized, stride of vertex (in bytes), pointer (offset)
	// point positions
	gl.VertexAttribPointer(0, int32(len(converted))/renderComponent.Stride,
		gl.FLOAT, false, renderComponent.Stride*4, nil)

	gl.EnableVertexAttribArray(0) // enable the defined vertex attribute pointer
	gl.EnableVertexAttribArray(1) // enable the defined vertex attribute pointer
	gl.BindVertexArray(0)         // unbind the array (there is no vertex array at 0)

	// save the pointers for later
	rs.glObjects[id] = []uint32{VBO, VAO}
}

func (rs *RenderSystem) Update(delta time.Duration) {
	gl.ClearColor(0, 0, 0, 0)     // set the clear color
	gl.Clear(gl.COLOR_BUFFER_BIT) // apply clear color

	gl.UseProgram(rs.shaderProgram)

	for id, components := range rs.MultiSystem.Entities {
		renderComponent := components[0].(*RenderComponent)

		converted := rs.convertVertexArray(id)
		glObj := rs.glObjects[id]

		// bind the buffer
		gl.BindBuffer(gl.ARRAY_BUFFER, glObj[0])

		// write new data into buffer
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(converted)*4, gl.Ptr(converted))

		// bind the vertex array
		gl.BindVertexArray(glObj[1])

		vertexCount := int32(len(renderComponent.Vertices)) / renderComponent.Stride
		gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	}
}

// convertVertexArray applies the SpaceComponent of an enttity to the vertex array of the entities texture
func (rs *RenderSystem) convertVertexArray(entityID uint) (converted []float32) {
	renderComponent := rs.MultiSystem.Entities[entityID][0].(*RenderComponent)
	spaceComponent := rs.MultiSystem.Entities[entityID][1].(*SpaceComponent)

	// copy original data
	// a "deep" copy is needed so we copy the data and not the pointer
	converted = make([]float32, len(renderComponent.Vertices))
	copy(converted, renderComponent.Vertices)

	// for every vertex
	for i := 0; i < len(converted); i += int(renderComponent.Stride) {
		// convert size
		// multiply x and y by size factor
		converted[i+0] *= spaceComponent.Size.X
		converted[i+1] *= spaceComponent.Size.Y

		// add position
		converted[i+0] += spaceComponent.Position.X
		converted[i+1] += spaceComponent.Position.Y
	}

	return
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
