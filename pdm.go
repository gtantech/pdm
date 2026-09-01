package pdm

import (
	"iter"
	"math"
	"slices"
	"time"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/activity/timestamp"
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

func (p *pdm) LoneActivities() func(yield func(vertex.Vertex[activity.Activity]) bool) {
	return p.filter(p.Activities(), func(v vertex.Vertex[activity.Activity]) bool {
		for range p.graph.OutgoingVertices(v) {
			// lone activity should not have any outgoing vertex, return false if it does
			return false
		}
		for range p.graph.IncomingVertices(v) {
			// lone activity should not have any incoming vertex, return false if it does
			return false
		}
		return true
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
func (p *pdm) earlyIntervals(topologicalSortedOrder []vertex.Vertex[activity.Activity]) map[activity.Activity]interval.Interval {
	earlyInterval := make(map[activity.Activity]interval.Interval)
	//initialize all start nodes
	for a := range p.InitialPredecessorActivities() {
		earlyInterval[a.Value()] = interval.New(time.Duration(0), a.Value().Duration())
	}
	for _, v := range topologicalSortedOrder {
		if _, ok := earlyInterval[v.Value()]; !ok {
			//!ok -> no early start/finish value for this activity
			//get predecessor
			maxEarlyFinishPredecessor := time.Duration(0)
			for predecessor := range p.graph.IncomingVertices(v) {
				if earlyFinish := earlyInterval[predecessor.Value()].Finish(); earlyFinish > maxEarlyFinishPredecessor {
					maxEarlyFinishPredecessor = earlyFinish
				}
			}
			earlyInterval[v.Value()] = interval.New(maxEarlyFinishPredecessor, maxEarlyFinishPredecessor+v.Value().Duration())
		}
	}
	return earlyInterval
}

func (p *pdm) UpdateActivityTimestamps() error {
	order, err := toposort.TopologicalSort(p.graph)
	if err != nil {
		return err
	}

	// forwards pass
	earlyIntervals := p.earlyIntervals(order)

	// backwards pass
	lateIntervals := make(map[activity.Activity]interval.Interval)

	//update with lone activities
	for a := range p.LoneActivities() {
		i := interval.New(time.Duration(0), a.Value().Duration())
		earlyIntervals[a.Value()] = i
		lateIntervals[a.Value()] = i
	}

	// initialize all end nodes
	for a := range p.FinalSuccessorActivities() {
		earlyInterval := earlyIntervals[a.Value()]
		lateFinish := earlyInterval.Finish()
		lateStart := lateFinish - a.Value().Duration()
		lateIntervals[a.Value()] = interval.New(lateStart, lateFinish)
	}

	slices.Reverse(order)
	for _, v := range order {
		if _, ok := lateIntervals[v.Value()]; !ok {
			//!ok -> no late start/finish value for this activity
			//get successor
			minLateStartSuccessor := time.Duration(math.MaxInt64)
			for successor := range p.graph.OutgoingVertices(v) {
				if lateStart := lateIntervals[successor.Value()].Start(); lateStart < minLateStartSuccessor {
					minLateStartSuccessor = lateStart
				}
			}
			lateIntervals[v.Value()] = interval.New(minLateStartSuccessor-v.Value().Duration(), minLateStartSuccessor)
		}
	}

	//update activities
	for v := range p.graph.Vertices() {
		v.Value().UpdateTimestamps(timestamp.New(earlyIntervals[v.Value()], lateIntervals[v.Value()]))
	}
	return nil
}
