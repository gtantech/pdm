package pdm

import (
	"iter"
	"time"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/dependency"
	"github.com/gtantech/pdm/interval"
	"github.com/gtantech/toposort"
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

func (p *pdm) InitialPredecessorActivities() func(yield func(vertex.Vertex[activity.Activity]) bool) {
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

func (p *pdm) FinalSuccessorActivities() func(yield func(vertex.Vertex[activity.Activity]) bool) {
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

// returns a map of activities to their early start and early finish values
func (p *pdm) earlyIntervals(targetStart time.Time) (map[activity.Activity]interval.Interval, error) {
	earlyInterval := make(map[activity.Activity]interval.Interval)
	//initialize all start nodes
	for a := range p.InitialPredecessorActivities() {
		earlyInterval[a.Value()] = interval.New(targetStart, targetStart.Add(a.Value().Duration()))
	}
	order, err := toposort.TopologicalSort(p.graph)
	if err != nil {
		return nil, err
	}
	for _, v := range order {
		if _, ok := earlyInterval[v.Value()]; !ok {
			//!ok -> no early start/finish value for this activity
			//get predecessor
			maxEarlyFinishPredecessor := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
			for predecessor := range p.graph.IncomingVertices(v) {
				if earlyFinish := earlyInterval[predecessor.Value()].Finish(); earlyFinish.Compare(maxEarlyFinishPredecessor) > 0 {
					maxEarlyFinishPredecessor = earlyFinish
				}
			}
			earlyInterval[v.Value()] = interval.New(maxEarlyFinishPredecessor, maxEarlyFinishPredecessor.Add(v.Value().Duration()))
		}
	}
	return earlyInterval, nil
}
