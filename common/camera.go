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
	projectionMatrix mgl32.Mat4
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
	spaceComponent.SetPosition(gome.FloatVector3{X: 0, Y: 0, Z: 0})
	spaceComponent.SetSize(gome.FloatVector3{X: 1, Y: 1, Z: 1})

	ce.BaseEntity.Components = map[string]gome.Component{
		"Space": spaceComponent,
	}

	ce.Lens(
		mgl32.DegToRad(100),
		16/9,
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
		projectionMatrix: mgl32.Perspective(fov, ratio, ncp, fcp),
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
func (cs *CameraSystem) projectionViewMatrix() mgl32.Mat4 {
	if cs.SingleSystem.Active {
		spaceComponent := cs.SingleSystem.Components[1].(*SpaceComponent)
		position := spaceComponent.GetPosition()
		rotation := spaceComponent.GetRotation()

		projectionMatrix := cs.SingleSystem.Components[0].(*CameraComponent).projectionMatrix
		viewMatrix := mgl32.LookAt(
			position.X, position.Y, position.Z,
			rotation.X, rotation.Y, rotation.Z,
			0, 1, 0,
		)
		return projectionMatrix.Mul4(viewMatrix)
	}

	return mgl32.Ident4()
}

func (*CameraSystem) Name() string { return "Camera" }

func (*CameraSystem) RequiredComponents() []string { return []string{"Camera", "Space"} }

func (*CameraSystem) Update(delta time.Duration) {}
