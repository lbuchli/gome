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
	spaceComponent.SetPosition(gome.FloatVector3{X: .5, Y: .5, Z: .5})
	spaceComponent.SetSize(gome.FloatVector3{X: .5, Y: .5, Z: .5})

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
	currentcontrolled := uint(3)
	currentRot := float32(0)

	gome.MailBox.Listen("MouseScroll", func(msg gome.Message) {
		mmsg := msg.(gome.MouseScrollMessage)
		spaceComponent := cs.SingleSystem.Components[1].(*common.SpaceComponent)

		currentRot += mmsg.X
		spaceComponent.SetRotation(gome.FloatVector3{X: 1, Y: 0, Z: 0}, currentRot)
	})

	gome.MailBox.Listen("MouseButton", func(msg gome.Message) {
		//mmsg := msg.(gome.MouseButtonMessage)
		//spaceComponent := cs.SingleSystem.Components[1].(*common.SpaceComponent)

		//if mmsg.State == sdl.PRESSED {
		//spaceComponent.Position.X = mmsg.X - 0.25
		//spaceComponent.Position.Y = mmsg.Y + 1.5
		//}
	})

	gome.MailBox.Listen("Keyboard", func(msg gome.Message) {
		kmsg := msg.(gome.KeyboardMessage)

		if kmsg.State == sdl.PRESSED {
			switch kmsg.Key.Sym {
			case sdl.K_SPACE:
				scene.RemoveComponent(currentcontrolled, "Control")
				currentcontrolled = currentcontrolled%3 + 1
				scene.AddComponent(currentcontrolled, &ControlComponent{})
			case sdl.K_TAB:
				// switch scenes
				gome.MailBox.Send(gome.ChangeSceneMessage{
					NewScene: 1,
					Relative: true,
				})
			}
		}
	})
}

func (cs *ControlSystem) Update(delta time.Duration) {}

func TestSpawn() {
	win := &gome.Window{
		Args: gome.WindowArguments{
			X:      0,
			Y:      0,
			Width:  1024,
			Height: 1024,
			Title:  "TEST",
			Debug:  true,
		},
	}

	win.Init()

	scene1 := &gome.Scene{}
	win.AddScene(scene1)

	pEntity := &PolygonEntity{
		Path: "/home/lukas/go/src/gitlocal/gome/testfiles/test1.obj",
	}

	pEntity.New()

	cameraEntity := &common.CameraEntity{}
	cameraEntity.New()

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
