package main

import (
	"gitlocal/gome"
	"gitlocal/gome/common"
	"path/filepath"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

type PolygonEntity struct {
	gome.BaseEntity
	Path string
}

func (pe *PolygonEntity) New() error {
	spaceComponent := &common.SpaceComponent{}
	spaceComponent.SetPosition(gome.FloatVector3{X: 0, Y: 0, Z: 0})
	spaceComponent.SetSize(gome.FloatVector3{X: 1, Y: 1, Z: 1})

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
	currentPos := gome.FloatVector3{X: 0, Y: 0, Z: 0}

	gome.MailBox.Listen("MouseScroll", func(msg gome.Message) {
		mmsg := msg.(gome.MouseScrollMessage)
		spaceComponent := cs.SingleSystem.Components[1].(*common.SpaceComponent)

		Y := float32(mmsg.X) / 100
		X := float32(mmsg.Y) / 100
		spaceComponent.AddRotation(gome.FloatVector3{X: 1, Y: 0, Z: 0}, X)
		spaceComponent.AddRotation(gome.FloatVector3{X: 0, Y: 1, Z: 0}, Y)
	})

	gome.MailBox.Listen("Keyboard", func(msg gome.Message) {
		kmsg := msg.(gome.KeyboardMessage)

		if kmsg.State == sdl.PRESSED {
			switch kmsg.Key.Sym {
			case sdl.K_w:
				currentPos.X -= .1
			case sdl.K_s:
				currentPos.X += .1
			case sdl.K_d:
				currentPos.Z -= .1
			case sdl.K_a:
				currentPos.Z += .1
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
			Width:  1080,
			Height: 1080,
			Title:  "TEST",
			Debug:  true,
		},
	}

	win.Init()

	scene1 := &gome.Scene{}
	win.AddScene(scene1)

	path, _ := filepath.Abs("../testfiles/test1.obj")
	pEntity := &PolygonEntity{
		Path: path,
	}

	pEntity.New()

	cameraEntity := &common.CameraEntity{}
	cameraEntity.New()

	cameraEntity.Lens(
		mgl32.DegToRad(120),
		1,
		0.1,
		100,
	)

	cameraEntity.BaseEntity.Components["Space"].(*common.SpaceComponent).SetPosition(gome.FloatVector3{X: 4, Y: 3, Z: 3})

	pEntity.BaseEntity.Components["Control"] = &ControlComponent{}

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
