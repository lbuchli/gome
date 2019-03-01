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

// Lens sets the perspective settings for the CameraSystem.
//
// fov:   The vertical Field of View, in radians: the amount of "zoom". Think "camera lens".
// ratio: Aspect Ratio. Depends on the size of your window.
// ncp:   Near clipping plane. Keep as big as possible, or you'll get precision issues.
// fcp:   Far clipping plane. Keep as little as possible.
func (cc *CameraComponent) Lens(fov, ratio, ncp, fcp float32) {
	cc.viewMatrix = mgl32.Perspective(fov, ratio, ncp, fcp)
}

func (*CameraComponent) Name() string { return "Camera" }

/*
	CameraEntity
*/

// A CameraEntity represents an invisible camera in the world.
type CameraEntity struct {
	gome.BaseEntity
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
	viewMatrix := cs.SingleSystem.GetComponent(cs.SingleSystem.ID, "Camera").(*CameraComponent).viewMatrix
	projectionMatrix := cs.SingleSystem.GetComponent(cs.SingleSystem.ID, "Space").(*SpaceComponent).modelMatrix
	return viewMatrix.Mul4(projectionMatrix)
}

func (*CameraSystem) Name() string { return "Camera" }

func (*CameraSystem) RequiredComponents() []string { return []string{"Camera, Space"} }

func (*CameraSystem) Update(delta time.Duration) {}
