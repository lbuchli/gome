package common

import "gitlocal/gome"

/*
	LightComponent
*/

type LightComponent struct {
	Strength float32
	Color    gome.FloatVector3
}

func (lc *LightComponent) Name() string { return "Light" }

/*
	LightSystem
*/

type lightSource struct {
	Position gome.FloatVector3
	Strength float32
	Color    gome.FloatVector3
}

type LightSystem struct {
	gome.MultiSystem
}

// getLightSources returns all the registered light sources:
// their poition, strength and color
func (ls *LightSystem) getLightSources() []lightSource {
	sources := make([]lightSource, len(ls.MultiSystem.Entities))
	index := 0 // cant use the range index, because the for iterates over a map
	for _, components := range ls.MultiSystem.Entities {
		lightComponent := components[0].(*LightComponent)
		spaceComponent := components[1].(*SpaceComponent)

		sources[index] = lightSource{
			Position: spaceComponent.GetPosition(),
			Strength: lightComponent.Strength,
			Color:    lightComponent.Color,
		}

		index++
	}

	return sources
}

func (ls *LightSystem) Name() string { return "Light" }

func (ls *LightSystem) RequiredComponents() []string { return []string{"Light", "Space"} }
