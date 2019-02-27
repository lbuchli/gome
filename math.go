package gome

const precision = float32(1*10 ^ -5)

/*
	Vectors
*/

// A Vector is an array of int32s with a fixed size
type Vector interface {
	// ToArray returns an array of the data contained in the vector.
	ToArray() []int32
}

// A FloatVector is an array of float32s with a fixed size
type FloatVector interface {
	// ToArray returns an array of the data contained in the vector.
	ToArray() []float32
}

// A Vector2 holds a X and a Y value and can be
// used for positions or directions.
type Vector2 struct {
	X, Y int32
}

func (v Vector2) ToArray() []int32 { return []int32{v.X, v.Y} }

// A FloatVector2 holds a X and a Y float value and can be
// used for positions or directions.
type FloatVector2 struct {
	X, Y float32
}

func (v FloatVector2) ToArray() []float32 { return []float32{v.X, v.Y} }

// A Vector3 holds a X, Y and a Z value and can be
// used for positions or directions.
type Vector3 struct {
	X, Y, Z int32
}

func (v Vector3) ToArray() []int32 { return []int32{v.X, v.Y, v.Z} }

// A FloatVector3 holds a X, Y and a Z float value and can be
// used for positions or directions.
type FloatVector3 struct {
	X, Y, Z float32
}

func (v FloatVector3) ToArray() []float32 { return []float32{v.X, v.Y, v.Z} }

// A Vector4 holds a W, X, Y and a Z value and can be
// used for positions, directions or colors.
type Vector4 struct {
	W, X, Y, Z int32
}

func (v Vector4) ToArray() []int32 { return []int32{v.W, v.X, v.Y, v.Z} }

// A FloatVector4 holds a W, X, Y and a Z float value and can be
// used for positions, directions or colors.
type FloatVector4 struct {
	W, X, Y, Z float32
}

func (v FloatVector4) ToArray() []float32 { return []float32{v.W, v.X, v.Y, v.Z} }

// isNear checks if two floats are almost equal.
func isNear(a float32, b float32) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}

	return diff < precision
}

// IsSimilarTo returns if a float vector is similar to another float vector.
func (fv FloatVector2) IsSimilarTo(other FloatVector2) bool {
	return isNear(fv.X, other.X) &&
		isNear(fv.Y, other.Y)
}

// IsSimilarTo returns if a float vector is similar to another float vector.
func (fv FloatVector3) IsSimilarTo(other FloatVector3) bool {
	return isNear(fv.X, other.X) &&
		isNear(fv.Y, other.Y) &&
		isNear(fv.Z, other.Z)
}

// IsSimilarTo returns if a float vector is similar to another float vector.
func (fv FloatVector4) IsSimilarTo(other FloatVector4) bool {
	return isNear(fv.W, other.W) &&
		isNear(fv.X, other.X) &&
		isNear(fv.Y, other.Y) &&
		isNear(fv.Z, other.Z)
}
