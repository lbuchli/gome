package common

import "gitlocal/gome"

/*
	LightComponent
*/

type LightType uint32

const (
	POINT_LIGHT = iota
	DIRECTIONAL_LIGHT
)

type LightComponent struct {
	Color       gome.FloatVector3
	Attenuation float32
	Type        LightType
}

func (lc *LightComponent) Name() string { return "Light" }

/*
	LightSystem
*/

type LightSource struct {
	Type        LightType
	Position    gome.FloatVector3
	Direction   gome.FloatVector3
	Color       gome.FloatVector3
	Attenuation float32
}

type LightSystem struct {
	gome.MultiSystem
}

// getLightSources returns all the registered light sources:
// their poition, strength and color
func (ls *LightSystem) getLightSources() []LightSource {
	sources := make([]LightSource, len(ls.MultiSystem.Entities))
	index := 0 // cant use the range index, because the for iterates over a map
	for _, components := range ls.MultiSystem.Entities {
		lightComponent := components[0].(*LightComponent)
		spaceComponent := components[1].(*SpaceComponent)

		sources[index] = LightSource{
			Type:        lightComponent.Type,
			Position:    spaceComponent.GetPosition(),
			Direction:   spaceComponent.GetRotation(),
			Attenuation: lightComponent.Attenuation,
			Color:       lightComponent.Color,
		}

		index++
	}

	return sources
}

func (ls *LightSystem) Name() string { return "Light" }

func (ls *LightSystem) RequiredComponents() []string { return []string{"Light", "Space"} }
