package common

import (
	"gitlocal/gome"

	"github.com/go-gl/mathgl/mgl32"
)

// A SpaceComponent gives an Entity a size and a position.
type SpaceComponent struct {
	modelMatrix mgl32.Mat4
	position    gome.FloatVector3
	rotation    gome.FloatVector3
}

func (*SpaceComponent) Name() string { return "Space" }

// LookAt positions and rotates the entity in space.
//
// pos:      The entities' position.
// target:   The target the entity looks at.
// rotation: The rotation of the entity. ((0, 1, 0) most of the time)
func (sc *SpaceComponent) LookAt(pos, target, rotation gome.FloatVector3) {
	sc.modelMatrix = mgl32.LookAt(
		pos.X, pos.Y, pos.Z,
		target.X, target.Y, target.Z,
		rotation.X, rotation.Y, rotation.Z,
	)

	sc.position = pos
	sc.rotation = target
}

func (sc *SpaceComponent) Rotate(rotation gome.FloatVector3) {
	// TODO
}

func (sc *SpaceComponent) Move(position gome.FloatVector3) {
	// TODO
}

func (sc *SpaceComponent) GetPosition() gome.FloatVector3 {
	return sc.position
}

func (sc *SpaceComponent) GetRotation() gome.FloatVector3 {
	return sc.rotation
}
