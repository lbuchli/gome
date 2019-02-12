package common

import "gitlocal/gome"

type SpaceComponent struct {
	Position gome.FloatVector
	Size     gome.FloatVector
}

func (*SpaceComponent) Name() string { return "Space" }
