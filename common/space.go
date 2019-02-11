package common

import "gitlocal/gome"

type SpaceComponent struct {
	Position gome.Vector
	Size     gome.Vector
}

func (*SpaceComponent) Name() string { return "Space" }
