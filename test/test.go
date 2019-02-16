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
	vdata := data.([]float32)
	pe.BaseEntity.Components = map[string]gome.Component{
		"Render": &common.RenderComponent{
			Stride:   3,
			Vertices: vdata[:len(vdata)-4],
			Indices: []uint32{
				0, 1, 2,
			},
			Color: common.Color{
				R: vdata[len(vdata)-4],
				G: vdata[len(vdata)-3],
				B: vdata[len(vdata)-2],
				A: vdata[len(vdata)-1],
			},
		},
		"Space": &common.SpaceComponent{
			Position: gome.FloatVector{X: 0, Y: 0, Z: 0},
			Size:     gome.FloatVector{X: 1, Y: 1, Z: 1},
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
	currentcontrolled := uint(3)

	gome.MailBox.Listen("MouseScroll", func(msg gome.Message) {
		mmsg := msg.(gome.MouseScrollMessage)
		spaceComponent := cs.SingleSystem.Components[1].(*common.SpaceComponent)

		spaceComponent.Position.X += mmsg.X / 128
		spaceComponent.Position.Y += mmsg.Y / 128
	})

	gome.MailBox.Listen("MouseButton", func(msg gome.Message) {
		mmsg := msg.(gome.MouseButtonMessage)
		spaceComponent := cs.SingleSystem.Components[1].(*common.SpaceComponent)

		if mmsg.State == sdl.PRESSED {
			spaceComponent.Position.X = mmsg.X - 0.25
			spaceComponent.Position.Y = mmsg.Y + 1.5
		}
	})

	gome.MailBox.Listen("Keyboard", func(msg gome.Message) {
		kmsg := msg.(gome.KeyboardMessage)

		if kmsg.State == sdl.PRESSED && kmsg.Key.Sym == sdl.K_SPACE {
			scene.RemoveComponent(currentcontrolled, "Control")
			currentcontrolled = currentcontrolled%3 + 1
			scene.AddComponent(currentcontrolled, &ControlComponent{})
		}
	})
}

func (cs *ControlSystem) Update(delta time.Duration) {}

func TestSpawn() {
	window := &gome.Window{
		Args: gome.WindowArguments{
			X:      0,
			Y:      0,
			Width:  1024,
			Height: 1024,
			Title:  "TEST",
			Debug:  true,
		},
	}

	window.Init()

	scene1 := &gome.Scene{}
	window.AddScene(scene1)

	pEntity := &PolygonEntity{}
	pEntity.New([]float32{
		-0.25, -0.25, 0.0,
		0.25, -0.25, 0.0,
		0.0, 0.25, 0.0,
		1.0, 0.0, 0.0, 1.0, // color
	})

	pEntity2 := &PolygonEntity{}
	pEntity2.New([]float32{
		0.25, -0.25, 0.0,
		0.75, -0.25, 0.0,
		0.5, 0.25, 0.0,
		0.0, 1.0, 0.0, 1.0, // color
	})

	pEntity3 := &PolygonEntity{}
	pEntity3.New([]float32{
		0.0, 0.25, 0.0,
		0.5, 0.25, 0.0,
		0.25, 0.75, 0.0,
		0.0, 0.0, 1.0, 1.0, // color
	})

	// make the entity controllable
	pEntity3.BaseEntity.Components["Control"] = &ControlComponent{}

	scene1.AddEntities(
		pEntity,
		pEntity2,
		pEntity3,
	)

	scene1.AddSystems(
		&common.RenderSystem{},
		&ControlSystem{},
	)

	window.Spawn()
}

func main() {
	TestSpawn()
}
