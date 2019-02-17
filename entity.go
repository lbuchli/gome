package gome

// An Entity groups components into one object. For example a player or
// a sheep are Entities.
type Entity interface {

	// GetID returns the unique identifier of the entity.
	GetID() uint

	// setID sets the entity ID.
	setID(uint)

	// GetComponents returns a map of Components the Entity contains.
	GetComponents() map[string]Component

	// addComponent adds a component to the entity
	addComponent(Component)

	// removeComponent removes a component from the entity
	removeComponent(string)

	// New fills the Entity with default and additional data. What data is accepted by the
	// entity depends on the Entity and should be documented by the method implementation.
	New() error
}

// A BaseEntity is a helper struct one can implement into their Entity. The BaseEntity is not
// an Entity on its own. (It lacks the 'New' function)
type BaseEntity struct {

	// The unique identifier of the Entity instance.
	id uint

	// The list of Components. Components should not be added after initialization (New function)
	// because they will be ignored by the active systems. TODO consider making this private
	Components map[string]Component
}

// GetID returns the unique identifier of the entity.
func (be *BaseEntity) GetID() uint { return be.id }

// SetID sets the entity ID. This method should not be overriden or used by the game designer.
func (be *BaseEntity) setID(eid uint) { be.id = eid }

// GetComponents returns a map of Components the Entity contains.
func (be *BaseEntity) GetComponents() map[string]Component { return be.Components }

func (be *BaseEntity) addComponent(component Component) {
	be.Components[component.Name()] = component
}

func (be *BaseEntity) removeComponent(name string) {
	delete(be.Components, name)
}
