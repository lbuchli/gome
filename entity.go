package gome

// An Entity groups components into one object. For example a player or
// a sheep are Entities.
type Entity interface {

	// GetID returns the unique identifier of the entity.
	GetID() uint

	// GetComponents returns a map of Components the Entity contains.
	GetComponents() map[string]Component

	// New fills the Entity with default and additional data. What data is accepted by the
	// entity depends on the Entity and should be documented by the method implementation.
	New(data interface{}) error
}

// A BaseEntity is a helper struct one can implement into their Entity. The BaseEntity is not
// an Entity on its own. (It lacks the 'New' function)
type BaseEntity struct {

	// The unique identifier of the Entity instance. TODO assure identifiers are unique
	ID uint

	// The list of Components. Components should not be added after initialization (New function)
	// because they will be ignored by the active systems. TODO consider making this private
	Components map[string]Component
}

// GetID returns the unique identifier of the entity.
func (be *BaseEntity) GetID() uint { return be.ID }

// GetComponents returns a map of Components the Entity contains.
func (be *BaseEntity) GetComponents() map[string]Component { return be.Components }
