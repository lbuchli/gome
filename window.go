package gome

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	Args        WindowArguments
	scenes      []*Scene
	current     uint
	window      *sdl.Window
	stopCurrent chan bool
}

type WindowArguments struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
	Title  string
	Debug  bool
}

func (win *Window) SetScene(scene uint) {
	// request stopping current scene
	win.stopCurrent <- true

	// wait until it has stopped
	<-win.stopCurrent

	// set new scene
	win.current = scene

	// and start it
	win.runCurrentScene()
}

func (win *Window) AddScene(scene *Scene) {
	scene.Init(win.Args)
	win.scenes = append(win.scenes, scene)
}

func (win *Window) AddScenes(scenes ...*Scene) {
	win.scenes = append(win.scenes, scenes...)
}

func (win *Window) GetScene(scene uint) *Scene {
	return win.scenes[scene]
}

func (win *Window) Init() {
	// init sdl (we use sdl for the window since it has more options than glfw)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	// create a new window with sdl
	win.window, err = sdl.CreateWindow(win.Args.Title, win.Args.X, win.Args.Y,
		win.Args.Height, win.Args.Height, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	win.window.GLCreateContext()

	win.scenes = make([]*Scene, 0)
}

// Spawn spawns the window and makes it visible
func (win *Window) Spawn() {
	defer sdl.Quit()
	defer win.window.Destroy()

	win.stopCurrent = make(chan bool, 1)
	win.runCurrentScene()
}

// handleEvents handles all the SDL events
func (win *Window) handleEvents(quit chan bool) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			quit <- true
			return
		case *sdl.KeyboardEvent:
			kEvent := event.(*sdl.KeyboardEvent)
			MailBox.Send(KeyboardMessage{
				Key:       kEvent.Keysym,
				State:     kEvent.State,
				Timestamp: kEvent.GetTimestamp(),
			})
		case *sdl.MouseButtonEvent:
			mEvent := event.(*sdl.MouseButtonEvent)
			MailBox.Send(MouseButtonMessage{
				Button:    mEvent.Button,
				State:     mEvent.State,
				X:         float32(mEvent.X)/float32(win.Args.Height)*2 - 1,
				Y:         -float32(mEvent.Y)/float32(win.Args.Width)*2 - 1,
				Timestamp: mEvent.GetTimestamp(),
			})
		case *sdl.MouseMotionEvent:
			mEvent := event.(*sdl.MouseMotionEvent)
			MailBox.Send(MouseMotionMessage{
				X:         float32(mEvent.X)/float32(win.Args.Height)*2 - 1,
				Y:         -float32(mEvent.Y)/float32(win.Args.Width)*2 - 1,
				XRel:      float32(mEvent.XRel)/float32(win.Args.Height)*2 - 1,
				YRel:      -float32(mEvent.YRel)/float32(win.Args.Width)*2 - 1,
				Timestamp: mEvent.GetTimestamp(),
			})
		case *sdl.MouseWheelEvent:
			mEvent := event.(*sdl.MouseWheelEvent)
			MailBox.Send(MouseScrollMessage{
				X:         float32(-mEvent.X),
				Y:         float32(mEvent.Y),
				Timestamp: mEvent.GetTimestamp(),
			})
		}
	}

	quit <- false
	return
}

// runScene initializes and then updates the scene for as long as it's running
func (win *Window) runCurrentScene() {
	// initialize all systems
	win.scenes[win.current].Init(win.Args)

	// marks when the last frame occured
	last := time.Now()

	running := true

	eventQuit := make(chan bool, 1)

	for running {
		// handle events; quit if requested
		go win.handleEvents(eventQuit)

		// calculate time since last frame
		delta := time.Since(last)
		last.Add(delta)

		// update scene (systems)
		win.scenes[win.current].Update(delta)

		// update the window
		win.window.GLSwap()

		// check if we should break the loop without blocking
		select {
		case stop := <-win.stopCurrent:
			running = stop
		default:
			// block until handleEvents returned something
			running = !<-eventQuit
		}
	}

	win.stopCurrent <- false
}
