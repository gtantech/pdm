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

// PDM stores and manages all activites and dependencies. Provides functions to manage a project via Precedence Diagram Method.
type PDM[D activity.Data] interface {
	// AddActivity adds an activity into [PDM] and returns the same activity.
	AddActivity(activity activity.Activity[D]) activity.Activity[D]

	// RemoveActivity removes the specified activity from [PDM].
	RemoveActivity(activity activity.Activity[D])

	// AddDependency adds a dependency from a predecessor activity to the successor activity, joined by the dependsVia relationship.
	AddDependency(predecessor activity.Activity[D], successor activity.Activity[D], dependsVia relationship.Relationship)

	// AddDependencies will add each element in the slice of [dependency.Dependency]
	AddDependencies(dependencies []dependency.Dependency[D]) error

	// AddDependenciesFromTable adds all [table.PredecessorDependency]
	AddDependenciesFromTable(table table.DependencyTable[D]) error

	// RemoveDependency removes a dependency from a predecessor activity to the successor activity.
	RemoveDependency(predecessor activity.Activity[D], successor activity.Activity[D])

	// GetRelationship returns the [relationship.Relationship] between the predecessor activity and successor activity.
	GetRelationship(predecessor activity.Activity[D], successor activity.Activity[D]) (relationship.Relationship, bool)

	// Activities returns an iterator over all activities in [PDM]. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
	Activities() func(yield func(activity.Activity[D]) bool)

	// Successors returns an iterator over succeeding activities to the activity specified. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
	Successors(activity activity.Activity[D]) func(yield func(activity.Activity[D]) bool)

	// Predecessors returns an iterator over preceeding activities to the activity specified. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
	Predecessors(activity activity.Activity[D]) func(yield func(activity.Activity[D]) bool)

	// InitialPredecessorActivities returns an iterator of all starting activities in [PDM]. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
	InitialPredecessorActivities() func(yield func(activity.Activity[D]) bool)

	// LoneActivities returns an iterator of all activities in [PDM] with no predecessor or successor in [PDM]. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
	LoneActivities() func(yield func(activity.Activity[D]) bool)

	// IntermediaryActivities returns an iterator over activities between the start and final activities. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
	IntermediaryActivities() func(yield func(activity.Activity[D]) bool)

	// FinalSuccessorActivities returns an iterator of all final activities in [PDM]. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
	FinalSuccessorActivities() func(yield func(activity.Activity[D]) bool)

	// CriticalActivities returns an iterator over activities that are considered critical (any delay will affect the project end date). The iteration order is not specified and is not guaranteed to be the same from one call to the next.
	//
	// The threshold parameter is the value at which if the total float of an activity is less than equal to the threshold value, it is considered critical.
	CriticalActivities(threshold time.Duration) func(yield func(activity.Activity[D]) bool)

	// UpdateActivityTimestamps updates all start/finish intervals for all activities in p. Returns CycleDetectedError if a cycle is detected in [PDM].
	UpdateActivityTimestamps() error

	// FreeFloat returns the maximum amount of time the specified activity can be delayed before the early start of any succeeding activity is delayed. (free float <= total float).
	FreeFloat(activity activity.Activity[D]) time.Duration

	// TotalFloat returns the float in activity [activity.Activity]. This is the amount of time the activity can be delayed without delaying the project end date
	TotalFloat(activity activity.Activity[D]) time.Duration

	// Duration returns the duration of the project [PDM] , ie. the project end date.
	Duration() time.Duration
}

var _ PDM[activity.Data] = (*pdm[activity.Data])(nil) //ensures pdm implements PDM at compile time

type pdm[D activity.Data] struct {
	graph graph.Graph[activity.Activity[D], relationship.Relationship]
}

// Duration implements [PDM]. Duration returns the duration of the project [PDM] , ie. the project end date.
func (p *pdm[D]) Duration() time.Duration {
	maxLateFinish := time.Duration(0)
	for a := range p.Activities() {
		if a.Late().Finish() > maxLateFinish {
			maxLateFinish = a.Late().Finish()
		}
	}
	return maxLateFinish
}

// GetRelationship implements [PDM]. GetRelationship returns the [relationship.Relationship] between the predecessor activity and successor activity. Will return false if failed to return.
//
// Added in pdm v1.1.0.
func (p *pdm[D]) GetRelationship(predecessor activity.Activity[D], successor activity.Activity[D]) (relationship.Relationship, bool) {
	d, ok := p.graph.GetEdgeValue(predecessor, successor)
	return d, ok
}

// AddDependenciesFromTable implements [PDM]. AddDependenciesFromTable adds all [table.PredecessorDependency] into p.
//
// Calls UpdateActivityTimestamps, and returns a CycleDetectedError if a cycle is encountered in p.
//
// Added in pdm v1.0.0.
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

// Predecessors implements [PDM]. Predecessors returns an iterator over preceeding activities to the activity specified. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
//
// Added in pdm v1.0.0.
func (p *pdm[D]) Predecessors(activity activity.Activity[D]) func(yield func(activity.Activity[D]) bool) {
	return p.graph.IncomingVertices(activity)
}

// Successors implements [PDM]. Successors returns an iterator over succeeding activities to the activity specified. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
//
// Added in pdm v1.0.0.
func (p *pdm[D]) Successors(activity activity.Activity[D]) func(yield func(activity.Activity[D]) bool) {
	return p.graph.OutgoingVertices(activity)
}

// IntermediaryActivities implements [PDM]. IntermediaryActivities returns an iterator over activities between the start and final activities. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
//
// Added in pdm v1.0.0.
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

// AddDependencies implements [PDM]. AddDependencies will add each element in the slice of [dependency.Dependency] into p. Will call UpdateActivityTimestamps() at the end.
func (p *pdm[D]) AddDependencies(dependencies []dependency.Dependency[D]) error {
	for _, d := range dependencies {
		p.AddDependency(d.Predecessor(), d.Successor(), d.DependsVia())
	}
	return p.UpdateActivityTimestamps()
}

// CriticalActivities implements [PDM]. CriticalActivities returns an iterator over activities that are considered critical (any delay will affect the project end date). The iteration order is not specified and is not guaranteed to be the same from one call to the next.
//
// The threshold parameter is the value at which if the total float of an activity is less than equal to the threshold value, it is considered critical.
//
// Added in pdm v1.0.0.
func (p *pdm[D]) CriticalActivities(threshold time.Duration) func(yield func(activity.Activity[D]) bool) {
	return p.filter(p.Activities(), func(v activity.Activity[D]) bool {
		return v.TotalFloat() <= threshold
	})
}

// TotalFloat implements [PDM]. TotalFloat returns the float in activity [activity.Activity]. This is the amount of time the activity can be delayed without delaying the project end date
//
// Added in pdm v1.0.0.
func (p *pdm[D]) TotalFloat(activity activity.Activity[D]) time.Duration {
	return activity.TotalFloat()
}

// New returns a new *[pdm]
//
// Added in pdm v1.0.0.
func New[D activity.Data]() *pdm[D] {
	return &pdm[D]{graph: graph.New[activity.Activity[D], relationship.Relationship]()}
}

// AddActivity implements [PDM]. AddActivity adds an activity into p and returns the same activity.
//
// Added in pdm v1.0.0.
func (p *pdm[D]) AddActivity(activity activity.Activity[D]) activity.Activity[D] {
	p.graph.AddVertex(activity)
	return activity
}

// RemoveActivity implements [PDM]. RemoveActivity removes the specified activity from p.
//
// Added in pdm v1.0.0.
func (p *pdm[D]) RemoveActivity(activity activity.Activity[D]) {
	p.graph.RemoveVertex(activity)
}

// AddDependency implements [PDM]. AddDependency adds a dependency from a predecessor activity to the successor activity, joined by the dependsVia relationship.
//
// Remember to call UpdateActivityTimestamps once all dependencies are added.
//
// Added in pdm v1.0.0.
func (p *pdm[D]) AddDependency(predecessor activity.Activity[D], successor activity.Activity[D], dependsVia relationship.Relationship) {
	p.graph.AddEdge(dependsVia, predecessor, successor)
}

// RemoveDependency implements [PDM]. RemoveDependency removes a dependency from a predecessor activity to the successor activity.
//
// Added in pdm v1.0.0.
func (p *pdm[D]) RemoveDependency(predecessor activity.Activity[D], successor activity.Activity[D]) {
	p.graph.RemoveEdge(predecessor, successor)
}

// Activities implements [PDM]. Activities returns an iterator over all activities in p. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
//
// Added in pdm v1.0.0.
func (p *pdm[D]) Activities() func(yield func(activity.Activity[D]) bool) {
	return p.graph.Vertices()
}

// FreeFloat implements [PDM]. FreeFloat returns the maximum amount of time the specified activity can be delayed before the early start of any succeeding activity is delayed. (free float <= total float).
//
// Added in pdm v1.0.0.
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

// InitialPredecessorActivities implements [PDM]. InitialPredecessorActivities returns an iterator of all starting activities in p. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
//
// Added in pdm v1.0.0.
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

// LoneActivities implements [PDM]. LoneActivities returns an iterator of all activities in p with no predecessor or successor. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
//
// Added in pdm v1.0.0.
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

// FinalSuccessorActivities implements [PDM]. FinalSuccessorActivities returns an iterator of all final activities in p. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
//
// Added in pdm v1.0.0.
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

// UpdateActivityTimestamps implements [PDM]. UpdateActivityTimestamps updates all start/finish intervals for all activities in p. Returns CycleDetectedError if a cycle is detected in p.
//
// Added in pdm v1.0.0.
func (p *pdm[D]) UpdateActivityTimestamps() error {
	return p.updateActivityTimestamp(toposort.TopologicalSort)
}
