package gome

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	Args    WindowArguments
	Scenes  []*Scene
	Current uint
	window  *sdl.Window
}

type WindowArguments struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
	Title  string
}

func (win *Window) SetScene(scene uint) {
	win.Current = scene

	// initialize all systems
	win.Scenes[win.Current].Init()
}

func (win *Window) AddScene(scene *Scene) {
	win.Scenes = append(win.Scenes, scene)
}

func (win *Window) AddScenes(scenes ...*Scene) {
	win.Scenes = append(win.Scenes, scenes...)
}

func (win *Window) GetScene(scene uint) *Scene {
	return win.Scenes[scene]
}

// Spawn spawns the window and makes it visible
func (win *Window) Spawn() {
	// init sdl (we use sdl for the window since it has more options than glfw)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// create a new window with sdl
	win.window, err = sdl.CreateWindow(win.Args.Title, win.Args.X, win.Args.Y,
		win.Args.Height, win.Args.Height, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	win.window.GLCreateContext()
	defer win.window.Destroy()

	win.runCurrentScene()
}

// handleEvents handles all the SDL events
func (win *Window) handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			return
		case *sdl.KeyboardEvent:
			kEvent := event.(*sdl.KeyboardEvent)
			MailBox.Send(KeyboardMessage{
				Key:       kEvent.Keysym,
				State:     kEvent.State,
				Timestamp: kEvent.GetTimestamp(),
			})
		}
	}
}

// runScene initializes and then updates the scene for as long as it's running
func (win *Window) runCurrentScene() {
	// initialize all systems
	win.Scenes[win.Current].Init()

	// marks when the last frame occured
	last := time.Now()

	for {
		win.handleEvents()

		// calculate time since last frame
		delta := time.Since(last)
		last.Add(delta)

		// update scene (systems)
		win.Scenes[win.Current].Update(delta)

		// update the window
		win.window.GLSwap()
	}
}
