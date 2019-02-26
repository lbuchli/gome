package common

import "gitlocal/gome"

// A SpaceComponent gives an Entity a size and a position.
type SpaceComponent struct {
	// the (left bottom) position (normalized)
	Position gome.FloatVector3

	// the size factor
	Size gome.FloatVector3
}

func (*SpaceComponent) Name() string { return "Space" }
