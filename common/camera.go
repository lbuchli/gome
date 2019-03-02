package common

import (
	"gitlocal/gome"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

/*
	CameraComponent
*/

type CameraComponent struct {
	viewMatrix mgl32.Mat4
}

func (*CameraComponent) Name() string { return "Camera" }

/*
	CameraEntity
*/

// A CameraEntity represents an invisible camera in the world.
type CameraEntity struct {
	gome.BaseEntity
}

func (ce *CameraEntity) New() error {
	spaceComponent := &SpaceComponent{}
	spaceComponent.SetPosition(gome.FloatVector3{X: 0.5, Y: 0, Z: 0.25})
	spaceComponent.SetSize(gome.FloatVector3{X: 1, Y: 1, Z: 1})

	ce.BaseEntity.Components = map[string]gome.Component{
		"Space": spaceComponent,
	}

	ce.Lens(
		mgl32.DegToRad(120),
		1,
		0.1,
		100,
	)

	return nil
}

// Lens sets the perspective settings for the CameraSystem.
//
// fov:   The vertical Field of View, in radians: the amount of "zoom". Think "camera lens".
// ratio: Aspect Ratio. Depends on the size of your window.
// ncp:   Near clipping plane. Keep as big as possible, or you'll get precision issues.
// fcp:   Far clipping plane. Keep as little as possible.
func (ce *CameraEntity) Lens(fov, ratio, ncp, fcp float32) {
	ce.BaseEntity.Components["Camera"] = &CameraComponent{
		viewMatrix: mgl32.Perspective(fov, ratio, ncp, fcp),
	}
}

/*
	CameraSystem
*/

// A CameraSystem defines how the world is viewed.
type CameraSystem struct {
	gome.SingleSystem
}

// viewProjectionMatrix returns the current View Projection Matrix
func (cs *CameraSystem) viewProjectionMatrix() mgl32.Mat4 {
	if cs.Active {
		viewMatrix := cs.SingleSystem.Components[0].(*CameraComponent).viewMatrix
		projectionMatrix := cs.SingleSystem.Components[1].(*SpaceComponent).modelMatrix()
		return viewMatrix.Mul4(projectionMatrix)
	}

	return mgl32.Ident4()
}

func (*CameraSystem) Name() string { return "Camera" }

func (*CameraSystem) RequiredComponents() []string { return []string{"Camera, Space"} }

func (*CameraSystem) Update(delta time.Duration) {}
