package main

import (
	"gitlocal/gome"
	"gitlocal/gome/common"
	"time"
)

/*
	Entity
*/

type CubeEntity struct {
	gome.BaseEntity
	Path string
}

// New fills the entity with initial data.
func (ce *CubeEntity) New() error {
	ce.BaseEntity.Components = map[string]gome.Component{
		"Render": &common.RenderComponent{
			OBJPath: ce.Path,
		},
		"Space": &common.SpaceComponent{},
	}

	return nil
}

/*
	System
*/

// RotationSystem rotates all its entities over time.
type RotationSystem struct {
	// A MultiSystem is a system that can contain 0 or more entities.
	gome.MultiSystem
}

// The RotationSystem requires a the SpaceComponent, because it sets the rotation there.
// The RenderComponent is required because else we would also rotate the camera.
func (rs *RotationSystem) RequiredComponents() []string { return []string{"Space", "Render"} }

func (rs *RotationSystem) Name() string { return "Rotation" }

func (rs *RotationSystem) Update(delta time.Duration) {
	rotation := float32(delta.Seconds() * 2)

	// iterate through the systems entities
	for _, components := range rs.MultiSystem.Entities {
		// the component order is always the same we gave in RequiredComponents
		spaceComponent := components[0].(*common.SpaceComponent)

		// rotate the entity by rotating its SpaceComponent
		spaceComponent.AddRotation(gome.FloatVector3{X: 1, Y: 0, Z: 0}, rotation)
	}
}

/*
	Main
*/

func main() {
	// create a new window
	win := gome.Window{
		Args: gome.WindowArguments{
			X:      0,
			Y:      0,
			Width:  1080,
			Height: 1080,
			Title:  "Hello World",
			Debug:  false,
		},
	}

	// ... and initialize it
	win.Init()

	// make a new Scene
	scene := &gome.Scene{}
	// add it to the window
	win.AddScene(scene)

	// make a new entity
	entity := &CubeEntity{
		Path: "cube.obj",
	}
	entity.New()

	// add the entity to the scene
	scene.AddEntity(entity)

	// add some required systems to the scene
	scene.AddSystems(
		&common.RenderSystem{}, // renders the cube entity
		&RotationSystem{},      // our system that rotates every visible entity
	)

	// start the window
	win.Spawn()
}
