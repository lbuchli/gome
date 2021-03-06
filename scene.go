package gome

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// A scene contains all Entities and Systems used for one game scene.
type Scene struct {
	entities      []Entity
	systems       []System
	idCount       uint
	isInitialized bool
	glcontext     sdl.GLContext
	WindowArgs    WindowArguments
}

// Update gets called every frame.
func (s *Scene) Update(delta time.Duration) {
	for _, system := range s.systems {
		system.Update(delta)
	}
}

// Init initializes the Scene, initializing the systems and adding
// entities to them.
func (s *Scene) Init(args WindowArguments) {
	// only initialize if it's not already initialized
	if !s.isInitialized {
		s.WindowArgs = args
		s.isInitialized = true

		for _, system := range s.systems {
			s.addSystemAfterInit(system)
		}

		for _, entity := range s.entities {
			s.addEntityAfterInit(entity)
		}
	}

	for _, system := range s.systems {
		system.Focus(s)
	}
}

// AddSystem adds a System to the Scene
func (s *Scene) AddSystem(system System) {
	s.systems = append(s.systems, system)

	if s.isInitialized {
		s.addSystemAfterInit(system)
	}
}

// HasSystem checks if there is a system of a specific type in this scene.
func (s *Scene) HasSystem(name string) bool {
	for _, system := range s.systems {
		if system.Name() == name {
			return true
		}
	}

	return false
}

// GetSystem returns a specific system in the scene. Returns nil if the system
// is not found.
func (s *Scene) GetSystem(name string) System {
	for _, system := range s.systems {
		if system.Name() == name {
			return system
		}
	}

	return nil
}

func (s *Scene) addSystemAfterInit(system System) {
	system.Init(s)

	required := system.RequiredComponents()

	for _, entity := range s.entities {
		components := entity.GetComponents()
		supply := []Component{}

		for _, requirement := range required {
			if val, ok := components[requirement]; ok {
				supply = append(supply, val)
			}
		}

		// if all the dependencies are satisfied, add the entity to the system
		if len(required) == len(supply) {
			system.Add(uint(entity.GetID()), supply)
		}
	}
}

// AddSystems adds multiple Systems to the Scene
func (s *Scene) AddSystems(systems ...System) {
	for _, system := range systems {
		s.AddSystem(system)
	}
}

// An entityID uniquely identifies an Entity
type entityID uint

// newEntityID returns  a new entity id.
// IDs start at 1 so you can check if an ID is initialized by checking if it's 0.
func (s *Scene) newEntityID() uint {
	s.idCount++
	return s.idCount
}

// AddEntity adds an Entity to the Scene
func (s *Scene) AddEntity(entity Entity) {
	entity.setID(s.newEntityID())
	s.entities = append(s.entities, entity)

	if s.isInitialized {
		s.addEntityAfterInit(entity)
	}
}

func (s *Scene) addEntityAfterInit(entity Entity) {
	components := entity.GetComponents()
	for _, system := range s.systems {
		// don't add to systems the entity is already in
		if system.Has(entity.GetID()) {
			continue
		}

		required := system.RequiredComponents()
		supply := []Component{}

		for _, requirement := range required {
			if val, ok := components[requirement]; ok {
				supply = append(supply, val)
			}
		}

		// if all the dependencies are satisfied, add the entity to the system
		if len(required) == len(supply) {
			system.Add(uint(entity.GetID()), supply)
		}
	}
}

// AddEntities adds multiple Entities to the Scene.
func (s *Scene) AddEntities(entities ...Entity) {
	for _, entity := range entities {
		s.AddEntity(entity)
	}
}

// RemoveEntity removes the Entity from all current systems and the entity list
func (s *Scene) RemoveEntity(id uint) {
	for i, entity := range s.entities {
		if entity.GetID() == id {
			// delete entity (for efficiency without preserving entity order)
			s.entities[i] = s.entities[len(s.entities)-1]
			s.entities[len(s.entities)-1] = nil
			s.entities = s.entities[:len(s.entities)-1]

			// there should be only one entity with that id
			break
		}
	}

	// remove entity from systems
	for _, system := range s.systems {
		system.Remove(id)
	}
}

// AddComponent adds a component to an existing entity.
func (s *Scene) AddComponent(id uint, component Component) {
	// add the component
	var components map[string]Component
	for i, entity := range s.entities {
		if entity.GetID() == id {
			entity.addComponent(component)
			components = s.entities[i].GetComponents()
			break
		}
	}

	// update the systems
	for _, system := range s.systems {
		// if the entity is already in the system, skip
		if system.Has(id) {
			continue
		}

		required := system.RequiredComponents()
		supply := []Component{}

		for _, requirement := range required {
			if val, ok := components[requirement]; ok {
				supply = append(supply, val)
			}
		}

		// if all the dependencies are satisfied, add the entity to the system
		if len(required) == len(supply) {
			system.Add(id, supply)
		}
	}
}

// RemoveComponent removes a Component from an existing Entity.
func (s *Scene) RemoveComponent(entityID uint, componentName string) {
	// update the systems
	var components map[string]Component
	for i, entity := range s.entities {
		if entity.GetID() == entityID {
			components = s.entities[i].GetComponents()
			break
		}
	}

	for _, system := range s.systems {
		required := system.RequiredComponents()
		possiblyAffected := false
		for _, reqCName := range required {
			if reqCName == componentName {
				possiblyAffected = true
				break
			}
		}

		// if the system seems affected, check with the requirements
		affected := possiblyAffected
		if possiblyAffected {
			for _, requirement := range required {
				if _, ok := components[requirement]; !ok {
					affected = false
				}
			}
		}

		if affected {
			// remove entity from system
			system.Remove(entityID)
		}
	}

	// remove the component
	for _, entity := range s.entities {
		if entity.GetID() == entityID {
			entity.removeComponent(componentName)
			break
		}
	}
}
