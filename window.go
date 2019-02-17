package gome

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	Args        WindowArguments
	scenes      []*Scene
	current     int
	window      *sdl.Window
	stopCurrent bool
}

type WindowArguments struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
	Title  string
	Debug  bool
}

func (win *Window) AddScene(scene *Scene) {
	context, err := win.window.GLCreateContext()
	if err != nil {
		Throw(err, "Could not create OpenGL context")
	}
	scene.glcontext = context

	win.scenes = append(win.scenes, scene)
}

func (win *Window) AddScenes(scenes ...*Scene) {
	for _, scene := range scenes {
		context, err := win.window.GLCreateContext()
		if err != nil {
			Throw(err, "Could not create OpenGL contexts")
		}
		scene.glcontext = context
	}

	win.scenes = append(win.scenes, scenes...)
}

func (win *Window) GetScene(scene int) *Scene {
	return win.scenes[scene]
}

func (win *Window) Current() int {
	return win.current
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

	MailBox.open()

	win.scenes = make([]*Scene, 0)
}

// Spawn spawns the window and makes it visible
func (win *Window) Spawn() {
	defer sdl.Quit()
	defer win.window.Destroy()

	// run current scene
	// when the current scene switches, the method will terminate
	// and it will loop to the next current scene.
	for {
		win.runCurrentScene()
	}
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
	if win.Args.Debug {
		fmt.Println("Displaying Scene:", win.current)
	}

	// listen for change scene message
	MailBox.Listen("ChangeScene", func(msg Message) {
		cmsg := msg.(ChangeSceneMessage)
		win.stopCurrent = true

		if cmsg.Relative {
			win.current += cmsg.NewScene

			// make scene increase wrap around
			win.current = win.current % len(win.scenes)
		} else {
			win.current = cmsg.NewScene
		}

		// reopen the mailbox
		MailBox.open()
	})

	// set OpenGL context
	err := win.window.GLMakeCurrent(win.scenes[win.current].glcontext)
	if err != nil {
		Throw(err, "Could not set OpenGL context")
	}

	// initialize all systems
	win.scenes[win.current].Init(win.Args)

	// marks when the last frame occured
	last := time.Now()

	running := true

	eventQuit := make(chan bool, 1)

	for running && !win.stopCurrent {
		// handle events; quit if requested
		go win.handleEvents(eventQuit)

		// calculate time since last frame
		delta := time.Since(last)
		last.Add(delta)

		// update scene (systems)
		win.scenes[win.current].Update(delta)

		// update the window
		win.window.GLSwap()

		// wait for event handling to finish
		<-eventQuit
	}

	// reset
	win.stopCurrent = false
}
