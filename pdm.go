package pdm

import (
	"iter"
	"math"
	"slices"
	"time"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/dependency"
	"github.com/gtantech/pdm/interval"
	"github.com/gtantech/toposort/v2"
	"github.com/gtantech/toposort/v2/graph"
)

type PDM[D activity.Data] interface {
	AddActivity(activity activity.Activity[D]) activity.Activity[D]
	RemoveActivity(activity activity.Activity[D])
	AddDependency(predecessor activity.Activity[D], successor activity.Activity[D], dependsVia dependency.Dependency)
	RemoveDependency(predecessor activity.Activity[D], successor activity.Activity[D])
	Activities() func(yield func(activity.Activity[D]) bool)
	InitialPredecessorActivities() func(yield func(activity.Activity[D]) bool)
	LoneActivities() func(yield func(activity.Activity[D]) bool)
	FinalSuccessorActivities() func(yield func(activity.Activity[D]) bool)
	UpdateActivityTimestamps() error
}

var _ PDM[activity.Data] = (*pdm[activity.Data])(nil) //ensures pdm implements PDM at compile time

type pdm[D activity.Data] struct {
	graph graph.Graph[activity.Activity[D], dependency.Dependency]
}

func New[D activity.Data]() *pdm[D] {
	return &pdm[D]{graph: graph.New[activity.Activity[D], dependency.Dependency]()}
}

func (p *pdm[D]) AddActivity(activity activity.Activity[D]) activity.Activity[D] {
	p.graph.AddVertex(activity)
	return activity
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

func (p *pdm[D]) Activities() func(yield func(activity.Activity[D]) bool) {
	return p.graph.Vertices()
}

func (p *pdm[D]) filter(seq iter.Seq[activity.Activity[D]], predicate func(activity.Activity[D]) bool) iter.Seq[activity.Activity[D]] {
	return func(yield func(activity.Activity[D]) bool) {
		seq(func(value activity.Activity[D]) bool {
			if predicate(value) {
				return yield(value)
			}
			return true
		})
	}
}

func (p *pdm[D]) InitialPredecessorActivities() func(yield func(activity.Activity[D]) bool) {
	return p.filter(p.Activities(), func(v activity.Activity[D]) bool {
		for range p.graph.IncomingVertices(v) {
			// starting predecessor should not have any incoming vertex, return false if it does
			return false
		}

		for range p.graph.OutgoingVertices(v) {
			// starting predecessor should have at least 1 outgoing vertex
			return true
		}
		return false
	})
}

func (p *pdm[D]) LoneActivities() func(yield func(activity.Activity[D]) bool) {
	return p.filter(p.Activities(), func(v activity.Activity[D]) bool {
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

func (p *pdm[D]) FinalSuccessorActivities() func(yield func(activity.Activity[D]) bool) {
	return p.filter(p.Activities(), func(v activity.Activity[D]) bool {
		for range p.graph.OutgoingVertices(v) {
			// starting successor should not have any outgoing vertex, return false if it does
			return false
		}

		for range p.graph.IncomingVertices(v) {
			// starting predecessor should have at least 1 incoming vertex
			return true
		}
		return false
	})
}

// returns a map of activities to their early start and early finish values
func (p *pdm[D]) earlyIntervals(topologicalSortedOrder []activity.Activity[D]) map[activity.Activity[D]]interval.Interval {
	earlyInterval := make(map[activity.Activity[D]]interval.Interval)
	//initialize all start nodes
	for a := range p.InitialPredecessorActivities() {
		earlyInterval[a] = interval.New(time.Duration(0), a.Data().Duration())
	}
	for _, v := range topologicalSortedOrder {
		if _, ok := earlyInterval[v]; !ok {
			//!ok -> no early start/finish value for this activity
			//get predecessor
			maxEarlyFinishPredecessor := time.Duration(0)
			for predecessor := range p.graph.IncomingVertices(v) {
				if earlyFinish := earlyInterval[predecessor].Finish(); earlyFinish > maxEarlyFinishPredecessor {
					maxEarlyFinishPredecessor = earlyFinish
				}
			}
			earlyInterval[v] = interval.New(maxEarlyFinishPredecessor, maxEarlyFinishPredecessor+v.Data().Duration())
		}
	}
	return earlyInterval
}

func (p *pdm[D]) updateActivityTimestamp(topologicalSorter func(g graph.Graph[activity.Activity[D], dependency.Dependency]) ([]activity.Activity[D], error)) error {
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
		i := interval.New(time.Duration(0), a.Data().Duration())
		earlyIntervals[a] = i
		lateIntervals[a] = i
	}

	// initialize all end nodes
	for a := range p.FinalSuccessorActivities() {
		earlyInterval := earlyIntervals[a]
		lateFinish := earlyInterval.Finish()
		lateStart := lateFinish - a.Data().Duration()
		lateIntervals[a] = interval.New(lateStart, lateFinish)
	}

	slices.Reverse(order)
	for _, v := range order {
		if _, ok := lateIntervals[v]; !ok {
			//!ok -> no late start/finish value for this activity
			//get successor
			minLateStartSuccessor := time.Duration(math.MaxInt64)
			for successor := range p.graph.OutgoingVertices(v) {
				if lateStart := lateIntervals[successor].Start(); lateStart < minLateStartSuccessor {
					minLateStartSuccessor = lateStart
				}
			}
			lateIntervals[v] = interval.New(minLateStartSuccessor-v.Data().Duration(), minLateStartSuccessor)
		}
	}

	//update activities
	for v := range p.graph.Vertices() {
		v.UpdateEarly(earlyIntervals[v])
		v.UpdateLate(lateIntervals[v])
	}
	return nil
}

func (p *pdm[D]) UpdateActivityTimestamps() error {
	return p.updateActivityTimestamp(toposort.TopologicalSort)
}
