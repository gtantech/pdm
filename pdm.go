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

func (p *pdm) AddActivity(activity activity.Activity) {
	p.graph.AddVertex(activity)
}

func (p *pdm) RemoveActivity(activity activity.Activity) {
	p.graph.RemoveVertex(activity)
}

func (p *pdm) AddDependency(predecessor activity.Activity, successor activity.Activity, dependsVia enums.Dependency) {
	p.graph.AddEdge(dependsVia, predecessor, successor)
}

func (p *pdm) RemoveDependency(predecessor activity.Activity, successor activity.Activity) {
	p.graph.RemoveEdge(predecessor, successor)
}
