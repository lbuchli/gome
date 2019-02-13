package main

import (
	"gitlocal/gome"
	"gitlocal/gome/common"
	"time"

	"github.com/veandco/go-sdl2/sdl"
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
			Position: gome.FloatVector{X: 0, Y: 0},
			Size:     gome.FloatVector{X: 1, Y: 1},
		},
	}

	return nil
}

/*
	ControlSystem
*/

type ControlComponent struct{}

func (*ControlComponent) Name() string { return "Control" }

type ControlSystem struct {
	gome.SingleSystem
}

func (*ControlSystem) RequiredComponents() []string { return []string{"Control", "Space"} }

func (cs *ControlSystem) Init(scene *gome.Scene) {
	cs.SingleSystem.Init(scene)

	gome.MailBox.Listen("Keyboard", func(msg gome.Message) {
		key := msg.(gome.KeyboardMessage).Key
		spaceComponent := cs.SingleSystem.Components[1].(*common.SpaceComponent)

		if msg.(gome.KeyboardMessage).State == sdl.PRESSED {
			switch key.Sym {
			case sdl.K_w:
				spaceComponent.Position.Y += .01
			case sdl.K_s:
				spaceComponent.Position.Y -= .01
			case sdl.K_a:
				spaceComponent.Position.X -= .01
			case sdl.K_d:
				spaceComponent.Position.X += .01
			}
		}
	})
}

func (cs *ControlSystem) Update(delta time.Duration) {}

func TestSpawn() {
	pEntity := &PolygonEntity{}
	pEntity.New([]float32{
		-0.25, -0.25, 0.0,
		0.25, -0.25, 0.0,
		0.0, 0.25, 0.0,
	})

	// make the entity controllable
	pEntity.BaseEntity.Components["Control"] = &ControlComponent{}

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
					&ControlSystem{},
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
