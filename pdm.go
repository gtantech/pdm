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

type PDM[D activity.Data] interface {
	AddActivity(activity activity.Activity[D])
	RemoveActivity(activity activity.Activity[D])
	AddDependency(predecessor activity.Activity[D], successor activity.Activity[D], dependsVia dependency.Dependency)
	RemoveDependency(predecessor activity.Activity[D], successor activity.Activity[D])
	Activities() func(yield func(vertex.Vertex[activity.Activity[D]]) bool)
	InitialPredecessorActivities() func(yield func(vertex.Vertex[activity.Activity[D]]) bool)
	LoneActivities() func(yield func(vertex.Vertex[activity.Activity[D]]) bool)
	FinalSuccessorActivities() func(yield func(vertex.Vertex[activity.Activity[D]]) bool)
	UpdateActivityTimestamps() error
}

var _ PDM[activity.Data] = (*pdm[activity.Data])(nil) //ensures pdm implements PDM at compile time

type pdm[D activity.Data] struct {
	graph graph.Graph[activity.Activity[D], dependency.Dependency]
}

func New[D activity.Data]() *pdm[D] {
	return &pdm[D]{graph: graph.New[activity.Activity[D], dependency.Dependency]()}
}

func (p *pdm[D]) AddActivity(activity activity.Activity[D]) {
	p.graph.AddVertex(activity)
}

func (p *pdm[D]) RemoveActivity(activity activity.Activity[D]) {
	p.graph.RemoveVertex(activity)
}

func (p *pdm[D]) AddDependency(predecessor activity.Activity[D], successor activity.Activity[D], dependsVia dependency.Dependency) {
	p.graph.AddEdge(dependsVia, predecessor, successor)
}

func (p *pdm[D]) RemoveDependency(predecessor activity.Activity[D], successor activity.Activity[D]) {
	p.graph.RemoveEdge(predecessor, successor)
}

func (p *pdm[D]) Activities() func(yield func(vertex.Vertex[activity.Activity[D]]) bool) {
	return p.graph.Vertices()
}

func (p *pdm[D]) filter(seq iter.Seq[vertex.Vertex[activity.Activity[D]]], predicate func(vertex.Vertex[activity.Activity[D]]) bool) iter.Seq[vertex.Vertex[activity.Activity[D]]] {
	return func(yield func(vertex.Vertex[activity.Activity[D]]) bool) {
		seq(func(value vertex.Vertex[activity.Activity[D]]) bool {
			if predicate(value) {
				return yield(value)
			}
			return true
		})
	}
}

func (p *pdm[D]) InitialPredecessorActivities() func(yield func(vertex.Vertex[activity.Activity[D]]) bool) {
	return p.filter(p.Activities(), func(v vertex.Vertex[activity.Activity[D]]) bool {
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

func (p *pdm[D]) LoneActivities() func(yield func(vertex.Vertex[activity.Activity[D]]) bool) {
	return p.filter(p.Activities(), func(v vertex.Vertex[activity.Activity[D]]) bool {
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

func (p *pdm[D]) FinalSuccessorActivities() func(yield func(vertex.Vertex[activity.Activity[D]]) bool) {
	return p.filter(p.Activities(), func(v vertex.Vertex[activity.Activity[D]]) bool {
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
func (p *pdm[D]) earlyIntervals(topologicalSortedOrder []vertex.Vertex[activity.Activity[D]]) map[activity.Activity[D]]interval.Interval {
	earlyInterval := make(map[activity.Activity[D]]interval.Interval)
	//initialize all start nodes
	for a := range p.InitialPredecessorActivities() {
		earlyInterval[a.Value()] = interval.New(time.Duration(0), a.Value().Data().Duration())
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
			earlyInterval[v.Value()] = interval.New(maxEarlyFinishPredecessor, maxEarlyFinishPredecessor+v.Value().Data().Duration())
		}
	}
	return earlyInterval
}

func (p *pdm[D]) updateActivityTimestamp(topologicalSorter func(g graph.Graph[activity.Activity[D], dependency.Dependency]) ([]vertex.Vertex[activity.Activity[D]], error)) error {
	order, err := topologicalSorter(p.graph)
	if err != nil {
		return err
	}

	// forwards pass
	earlyIntervals := p.earlyIntervals(order)

	// backwards pass
	lateIntervals := make(map[activity.Activity[D]]interval.Interval)

	//update with lone activities
	for a := range p.LoneActivities() {
		i := interval.New(time.Duration(0), a.Value().Data().Duration())
		earlyIntervals[a.Value()] = i
		lateIntervals[a.Value()] = i
	}

	// initialize all end nodes
	for a := range p.FinalSuccessorActivities() {
		earlyInterval := earlyIntervals[a.Value()]
		lateFinish := earlyInterval.Finish()
		lateStart := lateFinish - a.Value().Data().Duration()
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
			lateIntervals[v.Value()] = interval.New(minLateStartSuccessor-v.Value().Data().Duration(), minLateStartSuccessor)
		}
	}

	//update activities
	for v := range p.graph.Vertices() {
		v.Value().UpdateTimestamps(timestamp.New(earlyIntervals[v.Value()], lateIntervals[v.Value()]))
	}
	return nil
}

func (p *pdm[D]) UpdateActivityTimestamps() error {
	return p.updateActivityTimestamp(toposort.TopologicalSort)
}
