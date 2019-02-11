package main

import (
	"gitlocal/gome"
	"gitlocal/gome/common"
)

type PolygonEntity struct {
	gome.BaseEntity
}

func (pe *PolygonEntity) New(data interface{}) error {
	pe.BaseEntity.Components = map[string]gome.Component{
		"Render": &common.RenderComponent{
			Stride:   3,
			Vertices: data.([]float32),
		},
		"Space": &common.SpaceComponent{
			Position: gome.Vector{X: 0, Y: 0},
			Size:     gome.Vector{X: 200, Y: 200},
		},
	}

	return nil
}

func TestSpawn() {
	pEntity := &PolygonEntity{}
	pEntity.New([]float32{
		-0.25, -0.25, 0.0,
		0.25, -0.25, 0.0,
		0.0, 0.25, 0.0,
	})

	pEntity2 := &PolygonEntity{}
	pEntity2.New([]float32{
		0.25, -0.25, 0.0,
		0.75, -0.25, 0.0,
		0.5, 0.25, 0.0,
	})
	pEntity2.BaseEntity.ID = 1

	pEntity3 := &PolygonEntity{}
	pEntity3.New([]float32{
		0.0, 0.25, 0.0,
		0.5, 0.25, 0.0,
		0.25, 0.75, 0.0,
	})
	pEntity3.BaseEntity.ID = 2

	window := &gome.Window{
		Args: gome.WindowArguments{
			X:      0,
			Y:      0,
			Width:  1024,
			Height: 1024,
			Title:  "TEST",
		},
		Scenes: []*gome.Scene{
			&gome.Scene{
				Entities: []gome.Entity{
					pEntity,
					pEntity2,
					pEntity3,
				},
				Systems: []gome.System{
					&common.RenderSystem{},
				},
			},
		},
		Current: 0,
	}

	window.Spawn()
}

func main() {
	TestSpawn()
}
