package gome

type Entity interface {
	GetID() uint
	GetComponents() map[string]Component
	New(interface{}) error
}

type BaseEntity struct {
	ID         uint
	Components map[string]Component
}

func (be *BaseEntity) GetID() uint { return be.ID }

func (be *BaseEntity) GetComponents() map[string]Component { return be.Components }
