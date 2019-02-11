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
	// RequiredComponents demands.
	Add(id uint, components []Component)

	// Remove removes a single entity added to the system and its components.
	Remove(id uint)

	// Init initializes the system.
	Init(scene *Scene)

	// Update gets called every frame and given the time since the last frame
	// in miliseconds.
	Update(delta time.Duration)
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

func (ms *MultiSystem) Init(scene *Scene) {
	ms.Entities = make(map[uint][]Component)
}

// A SingleSystem is a base system that can only hold one entity.
// Use case would be a scrolling background for example.
type SingleSystem struct {
	id         uint
	components []Component
	active     bool
}

func (ss *SingleSystem) GetComponent(id uint, name string) Component {
	if ss.id == id {
		var i int
		req := ss.RequiredComponents()
		for i = 0; i < len(req); i++ {
			if req[i] == name {
				break
			}
		}

		return ss.components[i]

	} else {
		return nil
	}
}

// RequiredComponents should be overwritten.
func (*SingleSystem) RequiredComponents() []string { return []string{} }

func (ss *SingleSystem) Add(id uint, components []Component) {
	ss.id = id
	ss.components = components
	ss.active = true
}

func (ss *SingleSystem) Remove(id uint) {
	ss.active = false
}

func (ss *SingleSystem) Init(scene *Scene) {
	ss.active = false
}
