package gome

import "time"

// A scene contains all Entities and Systems used for one game scene.
type Scene struct {
	Entities []Entity
	Systems  []System
}

// Update gets called every frame.
func (s *Scene) Update(delta time.Duration) {
	for _, system := range s.Systems {
		system.Update(delta)
	}
}

// Init initializes the Scene, initializing the systems and adding
// entities to them.
func (s *Scene) Init() {
	// initialize the systems
	for _, system := range s.Systems {
		system.Init(s)
	}

	// add the enitites to the systems
	for _, system := range s.Systems {
		required := system.RequiredComponents()

		for _, entity := range s.Entities {
			components := entity.GetComponents()
			supply := []Component{}

			for _, requirement := range required {
				if val, ok := components[requirement]; ok {
					supply = append(supply, val)
				}
			}

			// if all the dependencies are satisfied, add the entity to the system
			if len(required) == len(supply) {
				system.Add(entity.GetID(), supply)
			}
		}
	}
}
