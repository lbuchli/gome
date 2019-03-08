package graphics

import (
	"bufio"
	"errors"
	"fmt"
	"gitlocal/gome"
	"io"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	Program uint32

	// uniform locations
	uniformLocs map[string]int32
	// unform block indices
	uniformBIndices map[string]uint32

	// uniform buffer objects
	uniformBOs map[string]uint32
}

// init initializes the shader and compiles the source
func (s *Shader) Init(shaderPath string) (err error) {
	s.uniformLocs = make(map[string]int32)
	s.uniformBIndices = make(map[string]uint32)
	s.uniformBOs = make(map[string]uint32)
	shaders := []uint32{}

	f, err := os.Open(shaderPath)
	if err != nil {
		return
	}

	reader := bufio.NewReader(f)

	shaderTypeLine, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// eof gets set to true if we reached the end of the file
	eof := false

	for {
		shader := ""

		// decide on the shader type
		var shaderType uint32
		typeStr := strings.Split(shaderTypeLine, " ")[1]
		switch strings.ToLower(typeStr) {
		case "vertex\n":
			shaderType = gl.VERTEX_SHADER
		case "fragment\n":
			shaderType = gl.FRAGMENT_SHADER
		default:
			err = errors.New("Shader type " + typeStr + " not known.")
			return err
		}

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				eof = true
				break
			}
			if err != nil {
				return err
			}

			// start a new shader string with a new type if the line starts with "#shader"
			if strings.HasPrefix(line, "#shader") {
				// tell the next iteration the information on shader type we read
				shaderTypeLine = line
				break
			} else {
				shader += line
			}
		}

		// if the shader is not empty, compile it
		if len(shader) > 0 {
			shaderptr, err := s.compile(shader+"\x00", shaderType)
			if err != nil {
				return err
			}

			shaders = append(shaders, shaderptr)
		} else {
			break
		}

		if eof {
			break
		}
	}

	// link shaders
	s.Program = gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(s.Program, shader)
	}
	gl.LinkProgram(s.Program)

	// delete the singke shaders. we won't need them anymore
	for _, shader := range shaders {
		gl.DeleteShader(shader)
	}

	return
}

// compile compiles a shader from source and returns the shader ID and, if occurred,
// an error. The shaderType can be any OpenGL shader type, e.g. gl.VERTEX_SHADER
func (s *Shader) compile(source string, shaderType uint32) (shader uint32, err error) {
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

// getUniformLocation gets the location of a uniform in the shader.
// May return -1 if the uniform is not found.
func (s *Shader) getUniformLocation(name string) (location int32) {
	// if we already saved the location, return it
	if location, ok := s.uniformLocs[name]; ok {
		return location
	}

	// if it's not in our location cache, get it from opengl and save it in the cache
	location = gl.GetUniformLocation(s.Program, gl.Str(name+"\x00"))
	if location == -1 {
		fmt.Println("ERROR: Could not find uniform:", name)
		return
	}

	s.uniformLocs[name] = location
	return
}

// getUniformLocation gets the location of a uniform in the shader.
// May return -1 if the uniform is not found.
func (s *Shader) getUniformBlockLocation(name string) (index uint32) {
	// if we already saved the location, return it
	if index, ok := s.uniformBIndices[name]; ok {
		return index
	}

	// if it's not in our index cache, get it from opengl and save it in the cache
	index = gl.GetUniformBlockIndex(s.Program, gl.Str(name+"\x00"))
	s.uniformBIndices[name] = index
	return
}

// Sets a uniform value.
func (s *Shader) SetUniformFVec2(name string, value gome.FloatVector2) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform2f(loc, value.X, value.Y)
	}
}

// Sets a uniform value.
func (s *Shader) SetUniformFVec3(name string, value gome.FloatVector3) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform3f(loc, value.X, value.Y, value.Z)
	}
}

// Sets a uniform value.
func (s *Shader) SetUniformFVec4(name string, value gome.FloatVector4) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform4f(loc, value.W, value.X, value.Y, value.Z)
	}
}

// Sets a uniform value.
func (s *Shader) SetUniformVec2(name string, value gome.Vector2) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform2i(loc, value.X, value.Y)
	}
}

// Sets a uniform value.
func (s *Shader) SetUniformVec3(name string, value gome.Vector3) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform3i(loc, value.X, value.Y, value.Z)
	}
}

// Sets a uniform value.
func (s *Shader) SetUniformVec4(name string, value gome.Vector4) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform4i(loc, value.W, value.X, value.Y, value.Z)
	}
}

// Sets a uniform value.
func (s *Shader) SetUniformFMat4(name string, value mgl32.Mat4) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.UniformMatrix4fv(loc, 1, false, &value[0])
	}
}

// SetUniformBlock sets the data of an active named uniform block.
func (s *Shader) SetUniformBlock(name string, value interface{}, size int) {
	// check if the uniform buffer object already exists
	ubo, ok := s.uniformBOs[name]
	if !ok {
		// if not, generate one
		gl.GenBuffers(1, &ubo)
		s.uniformBOs[name] = ubo
	}

	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)

	// set the new data
	gl.BufferData(gl.UNIFORM_BUFFER, size, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, s.getUniformBlockLocation(name), ubo)
}
