package main

import (
	"gitlocal/gome"
	"gitlocal/gome/common"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type PolygonEntity struct {
	gome.BaseEntity
	Path string
}

func (pe *PolygonEntity) New() error {
	spaceComponent := &common.SpaceComponent{}
	spaceComponent.SetPosition(gome.FloatVector3{X: 0, Y: 0, Z: -.5})
	spaceComponent.SetSize(gome.FloatVector3{X: 2, Y: 2, Z: 2})

	pe.BaseEntity.Components = map[string]gome.Component{
		"Render": &common.RenderComponent{
			OBJPath: pe.Path,
		},
		"Space": spaceComponent,
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

func (*ControlSystem) Name() string { return "Control" }

func (cs *ControlSystem) Focus(scene *gome.Scene) {
	cs.SingleSystem.Focus(scene)
	currentRot := gome.FloatVector2{X: 0, Y: 0}
	currentPos := gome.FloatVector3{X: 0, Y: 0, Z: 0}

	gome.MailBox.Listen("MouseScroll", func(msg gome.Message) {
		mmsg := msg.(gome.MouseScrollMessage)
		spaceComponent := cs.SingleSystem.Components[1].(*common.SpaceComponent)

		currentRot.X += float32(mmsg.X) / 1000000
		currentRot.Y += float32(mmsg.Y) / 1000000
		spaceComponent.SetRotation(gome.FloatVector3{X: 1, Y: 0, Z: 0}, currentRot.X)
		spaceComponent.SetRotation(gome.FloatVector3{X: 0, Y: 1, Z: 0}, currentRot.Y)
	})

	gome.MailBox.Listen("Keyboard", func(msg gome.Message) {
		kmsg := msg.(gome.KeyboardMessage)

		if kmsg.State == sdl.PRESSED {
			switch kmsg.Key.Sym {
			case sdl.K_w:
				currentPos.Z += .0001
			case sdl.K_s:
				currentPos.Z -= .0001
			case sdl.K_d:
				currentPos.X += .0001
			case sdl.K_a:
				currentPos.X -= .0001
			}

			spaceComponent := cs.SingleSystem.Components[1].(*common.SpaceComponent)
			spaceComponent.SetPosition(currentPos)
		}
	})
}

func (cs *ControlSystem) Update(delta time.Duration) {}

func TestSpawn() {
	win := &gome.Window{
		Args: gome.WindowArguments{
			X:      0,
			Y:      0,
			Width:  1920,
			Height: 1080,
			Title:  "TEST",
			Debug:  true,
		},
	}

	win.Init()

	scene1 := &gome.Scene{}
	win.AddScene(scene1)

	pEntity := &PolygonEntity{
		Path: "/home/lukas/go/src/gitlocal/gome/testfiles/test2.obj",
	}

	pEntity.New()

	cameraEntity := &common.CameraEntity{}
	cameraEntity.New()

	cameraEntity.Lens(
		.5,
		1,
		0.1,
		100,
	)

	cameraEntity.BaseEntity.Components["Control"] = &ControlComponent{}

	scene1.AddEntities(
		pEntity,
		cameraEntity,
	)
	scene1.AddSystems(
		&common.RenderSystem{},
		&ControlSystem{},
		&common.CameraSystem{},
	)

	win.Spawn()
}

func main() {
	TestSpawn()
}
