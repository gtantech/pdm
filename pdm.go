package pdm

import (
	"errors"
	"fmt"
	"iter"
	"math"
	"slices"
	"time"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/dependency"
	"github.com/gtantech/pdm/dependency/table"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/pdm/interval"
	"github.com/gtantech/pdm/relationship"
	"github.com/gtantech/toposort/v2"
	"github.com/gtantech/toposort/v2/graph"
)

type PDM[D activity.Data] interface {
	AddActivity(activity activity.Activity[D]) activity.Activity[D]
	RemoveActivity(activity activity.Activity[D])
	AddDependency(predecessor activity.Activity[D], successor activity.Activity[D], dependsVia relationship.Relationship)
	AddDependencies(dependencies []dependency.Dependency[D]) error
	AddDependenciesFromTable(table table.DependencyTable[D]) error
	RemoveDependency(predecessor activity.Activity[D], successor activity.Activity[D])
	Activities() func(yield func(activity.Activity[D]) bool)
	Successors(activity activity.Activity[D]) func(yield func(activity.Activity[D]) bool)
	Predecessors(activity activity.Activity[D]) func(yield func(activity.Activity[D]) bool)
	InitialPredecessorActivities() func(yield func(activity.Activity[D]) bool)
	LoneActivities() func(yield func(activity.Activity[D]) bool)
	IntermediaryActivities() func(yield func(activity.Activity[D]) bool)
	FinalSuccessorActivities() func(yield func(activity.Activity[D]) bool)
	CriticalActivities(threshold time.Duration) func(yield func(activity.Activity[D]) bool)
	UpdateActivityTimestamps() error
	FreeFloat(activity activity.Activity[D]) time.Duration
	TotalFloat(activity activity.Activity[D]) time.Duration
}

var _ PDM[activity.Data] = (*pdm[activity.Data])(nil) //ensures pdm implements PDM at compile time

type pdm[D activity.Data] struct {
	graph graph.Graph[activity.Activity[D], relationship.Relationship]
}

// AddDependenciesFromTable implements [PDM]. Calls UpdateActivityTimestamps.
func (p *pdm[D]) AddDependenciesFromTable(table table.DependencyTable[D]) error {
	for successor := range table.GetActivities() {
		row, ok := table.GetRow(successor)
		if !ok {
			panic(fmt.Sprintf("failed to get predecessors from dependency table %v", table))
		}
		for _, predecessor := range row {
			p.AddDependency(predecessor.Predecessor(), successor, predecessor.DependsVia())
		}
	}
	return p.UpdateActivityTimestamps()
}

// Predecessors implements [PDM].
func (p *pdm[D]) Predecessors(activity activity.Activity[D]) func(yield func(activity.Activity[D]) bool) {
	return p.graph.IncomingVertices(activity)
}

// Successors implements [PDM].
func (p *pdm[D]) Successors(activity activity.Activity[D]) func(yield func(activity.Activity[D]) bool) {
	return p.graph.OutgoingVertices(activity)
}

// IntermediaryActivities implements [PDM].
func (p *pdm[D]) IntermediaryActivities() func(yield func(activity.Activity[D]) bool) {
	return p.filter(p.Activities(), func(v activity.Activity[D]) bool {
		hasIncoming := false
		for range p.graph.IncomingVertices(v) {
			if hasIncoming {
				break
			}
			hasIncoming = true
		}
		if !hasIncoming {
			return false
		}
		for range p.graph.OutgoingVertices(v) {
			return true
		}
		return false
	})
}

// AddDependencies implements [PDM]. Will call UpdateActivityTimestamps().
func (p *pdm[D]) AddDependencies(dependencies []dependency.Dependency[D]) error {
	for _, d := range dependencies {
		p.AddDependency(d.Predecessor(), d.Successor(), d.DependsVia())
	}
	return p.UpdateActivityTimestamps()
}

// CriticalActivities implements [PDM].
// threshold is the value at which if the total float of an activity is less than equal to the threshold value, it is considered critical.
func (p *pdm[D]) CriticalActivities(threshold time.Duration) func(yield func(activity.Activity[D]) bool) {
	return p.filter(p.Activities(), func(v activity.Activity[D]) bool {
		return v.TotalFloat() <= threshold
	})
}

// TotalFloat implements [PDM].
func (p *pdm[D]) TotalFloat(activity activity.Activity[D]) time.Duration {
	return activity.TotalFloat()
}

func New[D activity.Data]() *pdm[D] {
	return &pdm[D]{graph: graph.New[activity.Activity[D], relationship.Relationship]()}
}

func (p *pdm[D]) AddActivity(activity activity.Activity[D]) activity.Activity[D] {
	p.graph.AddVertex(activity)
	return activity
}

func (p *pdm[D]) RemoveActivity(activity activity.Activity[D]) {
	p.graph.RemoveVertex(activity)
}

func (p *pdm[D]) AddDependency(predecessor activity.Activity[D], successor activity.Activity[D], dependsVia relationship.Relationship) {
	p.graph.AddEdge(dependsVia, predecessor, successor)
}

func (p *pdm[D]) RemoveDependency(predecessor activity.Activity[D], successor activity.Activity[D]) {
	p.graph.RemoveEdge(predecessor, successor)
}

func (p *pdm[D]) Activities() func(yield func(activity.Activity[D]) bool) {
	return p.graph.Vertices()
}

func (p *pdm[D]) FreeFloat(activity activity.Activity[D]) time.Duration {
	successorMinimumEarlyStart := time.Duration(math.MaxInt64)
	for successor := range p.graph.OutgoingVertices(activity) {
		if successor.Early().Start() < successorMinimumEarlyStart {
			successorMinimumEarlyStart = successor.Early().Start()
		}
	}
	if successorMinimumEarlyStart == time.Duration(math.MaxInt64) {
		//no outgoing vertices
		return time.Duration(0) //return 0 to satisfy (free float <= total float)
	}
	return successorMinimumEarlyStart - activity.Early().Start() - activity.Data().Duration()
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

func getDependency[D activity.Data](g graph.Graph[activity.Activity[D], relationship.Relationship], predecessor, successor activity.Activity[D]) relationship.Relationship {
	dependency, ok := g.GetEdgeValue(predecessor, successor)
	if !ok {
		panic("expected successful depedency value assignment")
	}
	return dependency
}

func (p *pdm[D]) forwardPassEarlyInterval(successor activity.Activity[D]) {
	maxStartValue := time.Duration(-1)
	//PICK THE LARGEST VALUE FOR THE SUCCESSOR START VALUE
	for predecessor := range p.graph.IncomingVertices(successor) {
		dependency := getDependency(p.graph, predecessor, successor)
		var successorCandidateStart time.Duration
		if dependency.Type() == enums.FS || dependency.Type() == enums.SS {
			successorCandidateStart = dependency.ForwardPassValue(predecessor)
		} else {
			successorCandidateInterval := interval.FromFinish(dependency.ForwardPassValue(predecessor), successor.Data().Duration())
			successorCandidateStart = successorCandidateInterval.Start()
		}
		if successorCandidateStart > maxStartValue {
			maxStartValue = successorCandidateStart
		}
	}
	if maxStartValue < 0 {
		return //max value was not updated because there was no predecessor
	}
	successor.UpdateEarly(interval.FromStart(maxStartValue, successor.Data().Duration()))
}
func (p *pdm[D]) backwardPassLateInterval(predecessor activity.Activity[D]) {
	minValue := time.Duration(math.MaxInt64)
	//PICK THE SMALLEST VALUE FOR THE PREDECESSOR FINISH VALUE
	for successor := range p.graph.OutgoingVertices(predecessor) {
		dependency := getDependency(p.graph, predecessor, successor)
		var predecessorCandidateFinish time.Duration
		if dependency.Type() == enums.FS || dependency.Type() == enums.FF {
			predecessorCandidateFinish = dependency.BackwardPassValue(successor)
		} else {
			predecessorCandidateInterval := interval.FromStart(dependency.BackwardPassValue(successor), predecessor.Data().Duration())
			predecessorCandidateFinish = predecessorCandidateInterval.Finish()
		}
		if predecessorCandidateFinish < minValue {
			minValue = predecessorCandidateFinish
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

func (p *pdm[D]) updateActivityTimestamp(topologicalSorter func(g graph.Graph[activity.Activity[D], relationship.Relationship]) ([]activity.Activity[D], error)) error {
	order, err := topologicalSorter(p.graph)
	if err != nil {
		var e *toposort.CycleDetectedError[activity.Activity[D], relationship.Relationship]
		if errors.As(err, &e) {
			return &CycleDetectedError[activity.Activity[D], relationship.Relationship]{Predecessor: e.Origin, Successor: e.Destination, Relationship: e.EdgeValue}
		}
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
