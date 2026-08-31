package pdm

import (
	"iter"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/dependency"
	"github.com/gtantech/toposort/graph"
	"github.com/gtantech/toposort/graph/vertex"
)

type pdm struct {
	graph       graph.Graph[activity.Activity, dependency.Dependency]
	displayName string
}

func New(name string) *pdm {
	return &pdm{graph: graph.New[activity.Activity, dependency.Dependency](), displayName: name}
}

func (p *pdm) DisplayName() string {
	return p.displayName
}

func (p *pdm) AddActivity(activity activity.Activity) {
	p.graph.AddVertex(activity)
}

func (p *pdm) RemoveActivity(activity activity.Activity) {
	p.graph.RemoveVertex(activity)
}

func (p *pdm) AddDependency(predecessor activity.Activity, successor activity.Activity, dependsVia dependency.Dependency) {
	p.graph.AddEdge(dependsVia, predecessor, successor)
}

func (p *pdm) RemoveDependency(predecessor activity.Activity, successor activity.Activity) {
	p.graph.RemoveEdge(predecessor, successor)
}

func (p *pdm) Activities() func(yield func(vertex.Vertex[activity.Activity]) bool) {
	return p.graph.Vertices()
}

func (p *pdm) filter(seq iter.Seq[vertex.Vertex[activity.Activity]], predicate func(vertex.Vertex[activity.Activity]) bool) iter.Seq[vertex.Vertex[activity.Activity]] {
	return func(yield func(vertex.Vertex[activity.Activity]) bool) {
		seq(func(value vertex.Vertex[activity.Activity]) bool {
			if predicate(value) {
				return yield(value)
			}
			return true
		})
	}
}

func (p *pdm) StartingPredecessorActivities() func(yield func(vertex.Vertex[activity.Activity]) bool) {
	return p.filter(p.Activities(), func(v vertex.Vertex[activity.Activity]) bool {
		for range p.graph.IncomingVertices(v) {
			// starting predecessor should not have any incoming vertex, return false if it does
			return false
		}
		hasOutgoingVertices := false

		for range p.graph.OutgoingVertices(v) {
			// starting predecessor should have at least 1 outgoing vertex
			if hasOutgoingVertices {
				break
			}
			hasOutgoingVertices = true
		}
		return hasOutgoingVertices
	})
}

func (p *pdm) FinishingSuccessorActivities() func(yield func(vertex.Vertex[activity.Activity]) bool) {
	return p.filter(p.Activities(), func(v vertex.Vertex[activity.Activity]) bool {
		for range p.graph.OutgoingVertices(v) {
			// starting successor should not have any outgoing vertex, return false if it does
			return false
		}
		hasIncomingVertices := false

		for range p.graph.IncomingVertices(v) {
			// starting predecessor should have at least 1 outgoing vertex
			if hasIncomingVertices {
				break
			}
			hasIncomingVertices = true
		}
		return hasIncomingVertices
	})
}
