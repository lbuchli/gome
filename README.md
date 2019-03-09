# Gome
## A Go 3D Game Engine

Gome is a simple 3D game engine following the ECS (Entity-Component-System) approach
to game engines. Gome was mostly designed for learning purposes and is not in a usable state yet.
It can currently read basic .obj files, render, move, rotate, resize and texture them in form of entities.
You can create your own Entities, Components and Systems to expand upon the ones already given by the
_gome/common_ package.

## TODO
Things not yet implemented include:
 - Animations
 - Light Sources
 - Particles
 - HUD Shader
 - Text Rendering
 - Flexible .OBJ file reading
 - Common physics system
 - Asset Management

## Installation
If for some reason you would want to install this hacky project, you can do it like this:

`go get -u github.com/phoenixdevelops/gome`

## Hello World
"Simple" HelloWorld: (Spinning textured cube, save your own texture under texture.jpg)

**game.go**
```go
package main

import (
	"github.com/phoenixdevelops/gome"
	"github.com/phoenixdevelops/gome/common"
	"time"
)

/*
	Entity
*/

type CubeEntity struct {
	gome.BaseEntity
	Path     string
	Position gome.FloatVector3
}

// New fills the entity with initial data.
func (ce *CubeEntity) New() error {
	spaceComponent := &common.SpaceComponent{}
	spaceComponent.SetPosition(ce.Position)

	ce.BaseEntity.Components = map[string]gome.Component{
		"Render": &common.RenderComponent{
			OBJPath: ce.Path,
		},
		"Space": spaceComponent,
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

// Update gets called every frame.
func (rs *RotationSystem) Update(delta time.Duration) {
	rotation := float32(delta.Seconds())

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

	// make a new camera
	camera := &common.CameraEntity{}
	camera.New()
	// set the position of the camera
	camera.Components["Space"].(*common.SpaceComponent).SetPosition(gome.FloatVector3{4, 3, 3})

	// add the entities to the scene
	scene.AddEntities(
		entity,
		camera,
	)

	// add some required systems to the scene
	scene.AddSystems(
		&common.RenderSystem{}, // renders the cube entity
		&common.CameraSystem{}, // perspective
		&RotationSystem{},      // our system that rotates every visible entity
	)

	// start the window
	win.Spawn()
}
```

**cube.obj**
```
usemtl texture.jpg
v 1.000000 -1.000000 -1.000000
v 1.000000 -1.000000 1.000000
v -1.000000 -1.000000 1.000000
v -1.000000 -1.000000 -1.000000
v 1.000000 1.000000 -1.000000
v 0.999999 1.000000 1.000001
v -1.000000 1.000000 1.000000
v -1.000000 1.000000 -1.000000
vt 0.748573 0.750412
vt 0.749279 0.501284
vt 0.999110 0.501077
vt 0.999455 0.750380
vt 0.250471 0.500702
vt 0.249682 0.749677
vt 0.001085 0.750380
vt 0.001517 0.499994
vt 0.499422 0.500239
vt 0.500149 0.750166
vt 0.748355 0.998230
vt 0.500193 0.998728
vt 0.498993 0.250415
vt 0.748953 0.250920
vn 0.000000 0.000000 -1.000000
vn -1.000000 -0.000000 -0.000000
vn -0.000000 -0.000000 1.000000
vn -0.000001 0.000000 1.000000
vn 1.000000 -0.000000 0.000000
vn 1.000000 0.000000 0.000001
vn 0.000000 1.000000 -0.000000
vn -0.000000 -1.000000 0.000000
f 5/1/1 1/2/1 4/3/1
f 5/1/1 4/3/1 8/4/1
f 3/5/2 7/6/2 8/7/2
f 3/5/2 8/7/2 4/8/2
f 2/9/3 6/10/3 3/5/3
f 6/10/4 7/6/4 3/5/4
f 1/2/5 5/1/5 2/9/5
f 5/1/6 6/10/6 2/9/6
f 5/1/7 8/11/7 6/10/7
f 8/11/7 7/12/7 6/10/7
f 1/2/8 2/9/8 3/13/8
f 1/2/8 3/13/8 4/14/8
```
