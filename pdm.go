package pdm

import (
	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/toposort/graph"
)

type pdm struct {
	graph.Graph[activity.Activity, enums.Dependency]
	Name string
}

var _ graph.Graph[activity.Activity, enums.Dependency] = (*pdm)(nil) //ensures PDM implements graph.Graph[activity.Activity, enums.Dependency] at compile time
