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

func (p *pdm[D]) forwardPassEarlyInterval(successor activity.Activity[D]) {
	maxValue := time.Duration(-1)
	for predecessor := range p.graph.IncomingVertices(successor) {
		dependency, ok := p.graph.GetEdgeValue(predecessor, successor)
		if !ok {
			panic("expected successful depedency value assignment")
		}
		if value := dependency.ForwardPassValue(predecessor); value > maxValue {
			maxValue = value
		}
	}
	if maxValue < 0 {
		return //max value was not updated because there was no predecessor
	}
	successor.UpdateEarly(interval.FromStart(maxValue, successor.Data().Duration()))
}
func (p *pdm[D]) backwardPassLateInterval(predecessor activity.Activity[D]) {
	minValue := time.Duration(math.MaxInt64)
	for successor := range p.graph.OutgoingVertices(predecessor) {
		dependency, ok := p.graph.GetEdgeValue(predecessor, successor)
		if !ok {
			panic("expected successful depedency value assignment")
		}
		if value := dependency.BackwardPassValue(successor); value < minValue {
			minValue = value
		}
	}
	if minValue == time.Duration(math.MaxInt64) {
		return //min value was not updated because there was no successor
	}
	predecessor.UpdateLate(interval.FromFinish(minValue, predecessor.Data().Duration()))
}

// returns a map of activities to their early start and early finish values
func (p *pdm[D]) forwardPass(topologicalSortedOrder []activity.Activity[D]) {
	//initialize all start nodes
	countInitial := 0
	for a := range p.InitialPredecessorActivities() {
		a.UpdateEarly(interval.New(time.Duration(0), a.Data().Duration()))
		countInitial++
	}
	for _, v := range topologicalSortedOrder[countInitial:] { //skip initial activities
		p.forwardPassEarlyInterval(v)
	}
}

func (p *pdm[D]) backwardPass(topologicalSortedOrder []activity.Activity[D]) {
	//initialize all final nodes
	order := topologicalSortedOrder
	slices.Reverse(order)
	countFinal := 0
	for a := range p.FinalSuccessorActivities() {
		a.UpdateLate(interval.FromFinish(a.Early().Finish(), a.Data().Duration()))
		countFinal++
	}
	for _, v := range order[countFinal:] { //skip final activities
		p.backwardPassLateInterval(v)
	}
}

func (p *pdm[D]) updateActivityTimestamp(topologicalSorter func(g graph.Graph[activity.Activity[D], dependency.Dependency]) ([]activity.Activity[D], error)) error {
	order, err := topologicalSorter(p.graph)
	if err != nil {
		return err
	}

	// forwards pass
	p.forwardPass(order)

	earlyIntervals := make(map[activity.Activity[D]]interval.Interval)
	for _, a := range order {
		earlyIntervals[a] = a.Early()
	}

	// backwards pass
	p.backwardPass(order)

	//update with lone activities
	for a := range p.LoneActivities() {
		i := interval.New(time.Duration(0), a.Data().Duration())
		a.UpdateEarly(i)
		a.UpdateLate(i)
	}
	return nil
}

func (p *pdm[D]) UpdateActivityTimestamps() error {
	return p.updateActivityTimestamp(toposort.TopologicalSort)
}
