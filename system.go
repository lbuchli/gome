package gome

import (
	"time"
)

// A System performs actions on Entities in the current scene. It gets initialized
// and can do changes every frame.
type System interface {
	// GetComponent gets a component of an entity added to the system.
	// If the entity or the component doesn't exist in the system, nil
	// will be returned.
	GetComponent(id uint, name string) Component

	// RequiredComponents gets a list of the components required from
	// an entity to be added to the system.
	RequiredComponents() []string

	// Add adds an entity to the system  or overwrites an existing one using
	// its ID and the components required by the system in the same order as
	// RequiredComponents demands. It should only be called when adding an Enttiy
	// to the scene.
	Add(id uint, components []Component)

	// remove removes a single entity added to the system and its components.
	// remove should only be called when removing an Entity from the scene.
	Remove(id uint)

	// Has checks if the entity is in the system.
	Has(id uint) bool

	// Init initializes the system.
	Init(scene *Scene)

	// Focus gets called when the scene gets shown.
	Focus(scene *Scene)

	// Update gets called every frame and given the time since the last frame
	// in miliseconds.
	Update(delta time.Duration)

	// Name returns the name of the system. (E.g. "Render")
	Name() string
}

// A MultiSystem is a base system that can hold multiple entites.
type MultiSystem struct {
	Entities map[uint][]Component
}

// GetComponent returns a specific component of a specific entity.
// This function should not be used by systems, because it's more efficient
// to directly access the component
func (ms *MultiSystem) GetComponent(id uint, name string) Component {
	var i int
	req := ms.RequiredComponents()
	for i = 0; i < len(req); i++ {
		if req[i] == name {
			break
		}
	}

	return ms.Entities[id][i]
}

// RequiredComponents should be overwritten.
func (*MultiSystem) RequiredComponents() []string { return []string{} }

func (ms *MultiSystem) Add(id uint, components []Component) {
	ms.Entities[id] = components
}

func (ms *MultiSystem) Remove(id uint) {
	delete(ms.Entities, id)
}

func (ms *MultiSystem) Has(id uint) bool {
	_, exists := ms.Entities[id]
	return exists
}

func (ms *MultiSystem) Init(scene *Scene) {
	ms.Entities = make(map[uint][]Component)
}

func (ms *MultiSystem) Focus(scene *Scene) {}

// A SingleSystem is a base system that can only hold one entity.
// Use case would be a scrolling background for example.
type SingleSystem struct {
	ID         uint
	Components []Component
	Active     bool
}

func (ss *SingleSystem) GetComponent(id uint, name string) Component {
	if ss.ID == id {
		var i int
		req := ss.RequiredComponents()
		for i = 0; i < len(req); i++ {
			if req[i] == name {
				break
			}
		}

		return ss.Components[i]

	} else {
		return nil
	}
}

// RequiredComponents should be overwritten.
func (*SingleSystem) RequiredComponents() []string { return []string{} }

func (ss *SingleSystem) Add(id uint, components []Component) {
	ss.ID = id
	ss.Components = components
	ss.Active = true
}

func (ss *SingleSystem) Remove(id uint) {
	ss.Active = false
}

func (ss *SingleSystem) Has(id uint) bool {
	return ss.ID == id
}

func (ss *SingleSystem) Init(scene *Scene) {
	ss.Active = false
}

func (ss *SingleSystem) Focus(scene *Scene) {}
