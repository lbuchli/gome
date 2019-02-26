package graphics

import (
	"bufio"
	"errors"
	"gitlocal/gome"
	"io"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Shader struct {
	filePath    string
	program     uint32
	shaders     []uint32
	uniformLocs map[string]int32
}

// init initializes the shader and compiles the source
func (s *Shader) Init(shaderPath string) (err error) {
	f, err := os.Open(shaderPath)
	if err != nil {
		return
	}

	reader := bufio.NewReader(f)

	for {
		shader := ""
		var shaderType uint32

		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return err
			}

			// start a new shader string with a new type if the line starts with "#shader"
			if strings.HasPrefix(line, "#shader") {
				typeStr := strings.Split(line, " ")[1]
				switch strings.ToLower(typeStr) {
				case "vertex\n":
					shaderType = gl.VERTEX_SHADER
				case "fragment\n":
					shaderType = gl.FRAGMENT_SHADER
				default:
					err = errors.New("Shader type " + typeStr + " not known.")
					return err
				}

				break
			} else {
				shader += line
			}
		}

		// if the shader is not empty, compile it
		if len(shader) > 0 {
			err = s.compile(shader+"\x00", shaderType)
			if err != nil {
				return
			}
		}
	}
}

// compile compiles a shader from source and returns the shader ID and, if occurred,
// an error. The shaderType can be any OpenGL shader type, e.g. gl.VERTEX_SHADER
func (s *Shader) compile(source string, shaderType uint32) (err error) {
	// create a shader from source (returns shader ID)
	shader := gl.CreateShader(shaderType)
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

	s.shaders = append(s.shaders, shader)

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
	location = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
	s.uniformLocs[name] = location
	return location
}

// Sets a uniform value.
func (s *Shader) SetUniformFVec2(name string, value gome.FloatVector2) {
	loc := s.getUniformLocation(name)
	gl.Uniform2f(loc, value.X, value.Y)
}

// Sets a uniform value.
func (s *Shader) SetUniformFVec3(name string, value gome.FloatVector3) {
	loc := s.getUniformLocation(name)
	gl.Uniform3f(loc, value.X, value.Y, value.Z)
}

// Sets a uniform value.
func (s *Shader) SetUniformFVec4(name string, value gome.FloatVector4) {
	loc := s.getUniformLocation(name)
	gl.Uniform4f(loc, value.W, value.X, value.Y, value.Z)
}

// Sets a uniform value.
func (s *Shader) SetUniformVec2(name string, value gome.Vector2) {
	loc := s.getUniformLocation(name)
	gl.Uniform2i(loc, value.X, value.Y)
}

// Sets a uniform value.
func (s *Shader) SetUniformVec3(name string, value gome.Vector3) {
	loc := s.getUniformLocation(name)
	gl.Uniform3i(loc, value.X, value.Y, value.Z)
}

// Sets a uniform value.
func (s *Shader) SetUniformVec4(name string, value gome.Vector4) {
	loc := s.getUniformLocation(name)
	gl.Uniform4i(loc, value.W, value.X, value.Y, value.Z)
}
