package gome

import (
	"io"
	"log"
	"runtime/debug"

	"github.com/veandco/go-sdl2/sdl"
)

// Throw outputs an error to console and stops execution
// it should only be used in top-level functions, as returning the
// error is required for unit testing
func Throw(err error, msg string) {
	if err != nil {
		debug.PrintStack()
		log.Fatalf("=> %v: %s", err, msg)
	} else {
		log.Fatalln(msg)
	}
}

/*
	MailBox
*/

// A Message is a piece of information sendable through the MailBox.
type Message interface {
	Name() string
}

type mailBox struct {
	listeners map[string][]func(Message)
}

// The MailBox is used to communicate between systems. Through the MailBox, one
// can send Messages and listen for them.
var MailBox *mailBox

// Send sends a Message through the MailBox to functions listening for
// that type of Message.
func (mb *mailBox) Send(msg Message) {
	for _, fun := range mb.listeners[msg.Name()] {
		fun(msg)
	}
}

// Listen adds the function to the group listening for a Message of a specific type.
func (mb *mailBox) Listen(msgName string, fun func(Message)) {
	mb.listeners[msgName] = append(mb.listeners[msgName], fun)
}

// open (re-)opens the mailbox, forgetting all current listeners.
func (mb *mailBox) open() {
	MailBox = &mailBox{make(map[string][]func(Message))}
}

/*
	Default Messages
*/

// A KeyboardMessage is sent when keyboard events occur.
type KeyboardMessage struct {
	Key       sdl.Keysym
	State     uint8
	Timestamp uint32
}

func (KeyboardMessage) Name() string { return "Keyboard" }

// A MouseButtonMessage is sent when a mouse button gets pressed or released.
type MouseButtonMessage struct {
	Button    uint8
	State     uint8
	X, Y      float32
	Timestamp uint32
}

func (MouseButtonMessage) Name() string { return "MouseButton" }

// A MouseMotionMessage is sent when the mouse gets moved.
type MouseMotionMessage struct {
	X, Y, XRel, YRel float32
	Timestamp        uint32
}

func (MouseMotionMessage) Name() string { return "MouseMotion" }

// A MouseScrollMessage is sent when the mouse wheel gets moved.
type MouseScrollMessage struct {
	X, Y      float32
	Timestamp uint32
}

func (MouseScrollMessage) Name() string { return "MouseScroll" }

// A ChangeSceneMessage will change the current scene
type ChangeSceneMessage struct {
	NewScene int
	Relative bool
}

func (ChangeSceneMessage) Name() string { return "ChangeScene" }

/*
	File Reader
*/

// A FileReader reads a file of a specific type. A FileReader normally
// has a Data(file io.Reader) method returning the whole file data.
type FileReader interface {
	// Check checks if a file is of the right type.
	Check(io.Reader) bool

	// Extension returns the default file extention for this file type.
	Extension() string
}
