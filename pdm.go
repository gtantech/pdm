package pdm

import (
	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/toposort/graph"
)

type pdm struct {
	graph graph.Graph[activity.Activity, enums.Dependency]
	Name  string
}
