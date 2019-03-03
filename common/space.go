package common

import (
	"gitlocal/gome"

	"github.com/go-gl/mathgl/mgl32"
)

// A SpaceComponent gives an Entity a size and a position.
type SpaceComponent struct {
	translationMatrix mgl32.Mat4
	rotationMatrix    mgl32.Mat4
	scaleMatrix       mgl32.Mat4
}

func (*SpaceComponent) Name() string { return "Space" }

func (sc *SpaceComponent) modelMatrix() mgl32.Mat4 {
	// make empty matrices to identity matrices so they don't change
	// the result of the multiplication
	empty := mgl32.Mat4{}
	if sc.translationMatrix == empty {
		sc.translationMatrix = mgl32.Ident4()
	}
	if sc.rotationMatrix == empty {
		sc.rotationMatrix = mgl32.Ident4()
	}
	if sc.scaleMatrix == empty {
		sc.rotationMatrix = mgl32.Ident4()
	}

	return sc.translationMatrix.
		Mul4(sc.rotationMatrix).
		Mul4(sc.scaleMatrix)
}

func (sc *SpaceComponent) SetPosition(pos gome.FloatVector3) {
	sc.translationMatrix = mgl32.Translate3D(pos.X, pos.Y, pos.Z)
}

func (sc *SpaceComponent) SetRotation(axis gome.FloatVector3, angle float32) {
	sc.rotationMatrix = sc.rotationMatrix.Mul4(mgl32.QuatRotate(angle,
		mgl32.Vec3{
			axis.X,
			axis.Y,
			axis.Z,
		}).Mat4())
}

func (sc *SpaceComponent) SetSize(size gome.FloatVector3) {
	sc.scaleMatrix = mgl32.Scale3D(size.X, size.Y, size.Z)
}

func (sc *SpaceComponent) GetPosition() gome.FloatVector3 {
	lastRow := sc.translationMatrix.Row(3)
	return gome.FloatVector3{
		X: lastRow.X(),
		Y: lastRow.Y(),
		Z: lastRow.Z(),
	}
}

func (sc *SpaceComponent) GetSize() gome.FloatVector3 {
	diag := sc.rotationMatrix.Diag()
	return gome.FloatVector3{
		X: diag.X(),
		Y: diag.Y(),
		Z: diag.Z(),
	}
}
