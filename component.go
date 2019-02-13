package gome

// A Component is part of an Entity, storing the Entities' data
// to be used by a System.
// A Component has a name, e.g. the struct 'RenderComponent'
// has the name 'Render'. By convention, the Compontent struct names have the
// suffix 'Component', and the name of the Component is the struct name without that suffix.
type Component interface {

	// Name returns the name of the Component, e.g. 'Render'
	Name() string
}
