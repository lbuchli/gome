package common

import (
	"fmt"
	"gitlocal/gome"
	"gitlocal/gome/common/graphics"
	"os"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

/*
	RenderComponent
*/

// A RenderComponent is a component used to render a texture of its
// entity
type RenderComponent struct {
	array        graphics.VertexArray
	OBJPath      string
	ModelUpdated bool
}

func (rc *RenderComponent) Name() string { return "Render" }

/*
	RenderSystem
*/

// A RenderSystem renders the texture of its entities
type RenderSystem struct {
	gome.MultiSystem

	graphics.Shader
	cameraSystem *CameraSystem
}

func (*RenderSystem) RequiredComponents() []string { return []string{"Render", "Space"} }

func (rs *RenderSystem) Init(scene *gome.Scene) {
	// initialize the base system
	rs.MultiSystem.Init(scene)

	// initialize OpenGL
	gl.Init()

	// Configure global opengl settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0, 0, 0, 0) // set the clear color

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

	// init shader
	rs.Shader.Init("/home/lukas/go/src/gitlocal/gome/common/graphics/default.shader")

	// get the camera system, and if there isn't one, add a new instance to the scene.
	if scene.HasSystem("Camera") {
		rs.cameraSystem = scene.GetSystem("Camera").(*CameraSystem)
	} else {
		rs.cameraSystem = &CameraSystem{}
		scene.AddSystem(rs.cameraSystem)

		// there's probably also no camera entity, so add it as well
		cameraEntity := &CameraEntity{}
		cameraEntity.New()
		scene.AddEntity(cameraEntity)
	}
}

func (rs *RenderSystem) Add(id uint, components []gome.Component) {
	rs.MultiSystem.Add(id, components)

	renderComponent := components[0].(*RenderComponent)

	f, err := os.Open(renderComponent.OBJPath)
	if err != nil {
		gome.Throw(err, "Could not open OBJ file")
	}

	reader := &graphics.OBJFileReader{}
	va, err := reader.Data(f)
	if err != nil {
		gome.Throw(err, "Could not read data from object file")
	}

	renderComponent.array = va
}

func (rs *RenderSystem) Update(delta time.Duration) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // apply clear color

	gl.UseProgram(rs.Shader.Program)

	// Projection View Matrix
	PVM := rs.cameraSystem.projectionViewMatrix()

	for _, components := range rs.MultiSystem.Entities {
		renderComponent := components[0].(*RenderComponent)
		spaceComponent := components[1].(*SpaceComponent)
		VAO := &renderComponent.array

		MVP := PVM.Mul4(spaceComponent.modelMatrix())
		rs.Shader.SetUniformFMat4("u_MVP", MVP)

		test := MVP.Mul4x1(mgl32.Vec4{1, 1, 1, 1})
		_ = fmt.Sprint(test)

		VAO.Draw()
	}
}

func (*RenderSystem) Name() string { return "Render" }
