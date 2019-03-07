package common

import (
	"gitlocal/gome"

	"github.com/go-gl/mathgl/mgl32"
)

// A SpaceComponent gives an Entity a size and a position.
type SpaceComponent struct {
	translationMatrix mgl32.Mat4
	rotationQuat      mgl32.Quat
	scaleMatrix       mgl32.Mat4
}

func (*SpaceComponent) Name() string { return "Space" }

// modelMatrix calculates the model matrix of the entity based on it's position,
// rotation and scale.
func (sc *SpaceComponent) modelMatrix() mgl32.Mat4 {
	// make empty matrices to identity matrices so they don't change
	// the result of the multiplication
	empty := mgl32.Mat4{}
	emptyQuat := mgl32.Quat{}
	if sc.translationMatrix == empty {
		sc.SetPosition(gome.FloatVector3{X: 0, Y: 0, Z: 0})
	}
	if sc.rotationQuat == emptyQuat {
		sc.rotationQuat = mgl32.QuatIdent()
	}
	if sc.scaleMatrix == empty {
		sc.SetSize(gome.FloatVector3{X: 1, Y: 1, Z: 1})
	}

	return sc.translationMatrix.
		Mul4(sc.rotationQuat.Mat4()).
		Mul4(sc.scaleMatrix)
}

// SetPosition sets the entities position in 3-dimensional space.
func (sc *SpaceComponent) SetPosition(pos gome.FloatVector3) {
	sc.translationMatrix = mgl32.Translate3D(pos.X, pos.Y, pos.Z)
}

// AddRotation adds to the rotation of the entity.
func (sc *SpaceComponent) AddRotation(axis gome.FloatVector3, angle float32) {
	mglAxis := mgl32.Vec3{axis.X, axis.Y, axis.Z}
	sc.rotationQuat = sc.rotationQuat.Mul(mgl32.QuatRotate(angle, mglAxis))
}

// SetSize sets the 3-dimensional scale of the entity.
func (sc *SpaceComponent) SetSize(size gome.FloatVector3) {
	sc.scaleMatrix = mgl32.Scale3D(size.X, size.Y, size.Z)
}

// GetPosition returns the current position of the entity.
func (sc *SpaceComponent) GetPosition() gome.FloatVector3 {
	lastCol := sc.translationMatrix.Col(3)
	return gome.FloatVector3{
		X: lastCol.X(),
		Y: lastCol.Y(),
		Z: lastCol.Z(),
	}
}

// GetRotation returns the current  rotation of the entity.
func (sc *SpaceComponent) GetRotation() gome.FloatVector3 {
	rotation := sc.rotationQuat.V
	return gome.FloatVector3{X: rotation[0], Y: rotation[1], Z: rotation[2]}
}

// GetSize returns the size of the entity.
func (sc *SpaceComponent) GetSize() gome.FloatVector3 {
	diag := sc.scaleMatrix.Diag()
	return gome.FloatVector3{
		X: diag.X(),
		Y: diag.Y(),
		Z: diag.Z(),
	}
}
