package table

import (
	"maps"

	"github.com/gtantech/pdm/activity"
)

type DependencyTable[D activity.Data] interface {
	AddRow(successor activity.Activity[D])
	GetRow(successor activity.Activity[D]) ([]PredecessorDependency[D], bool)
	GetActivities() func(yield func(activity.Activity[D]) bool)
	UpdatePredecessors(successor activity.Activity[D], predecessors []PredecessorDependency[D])
	DeleteRow(successor activity.Activity[D])
}

type depedencyTable[D activity.Data] struct {
	table map[activity.Activity[D]][]PredecessorDependency[D]
}

// GetActivities implements [DependencyTable]. GetActivities returns an iterator over activities in d. The iteration order is not specified and is not guaranteed to be the same from one call to the next.
//
// Added in pdm v1.0.0.
func (d *depedencyTable[D]) GetActivities() func(yield func(activity.Activity[D]) bool) {
	return maps.Keys(d.table)
}

// GetRow implements [DependencyTable]. GetRow returns a slice of [PredecessorDependency] and bool. The boolean value returns true if a successful result is returned, otherwise false.
//
// Added in pdm v1.0.0.
func (d *depedencyTable[D]) GetRow(successor activity.Activity[D]) ([]PredecessorDependency[D], bool) {
	predecessors, ok := d.table[successor]
	return predecessors, ok
}

// AddRow implements [DependencyTable]. AddRow adds an activity to the table and associates a empty slice of [PredecessorDependency] with it.
//
// Added in pdm v1.0.0.
func (d *depedencyTable[D]) AddRow(successor activity.Activity[D]) {
	d.table[successor] = []PredecessorDependency[D]{}
}

// DeleteRow implements [DependencyTable]. DeleteRow deletes a row from the table with the specificed successor activity.
//
// Added in pdm v1.0.0.
func (d *depedencyTable[D]) DeleteRow(successor activity.Activity[D]) {
	delete(d.table, successor)
}

// UpdatePredecessors implements [DependencyTable]. UpdatePredecessors replaces the predecessors in the table at the specified successor activity.
//
// Added in pdm v1.0.0.
func (d *depedencyTable[D]) UpdatePredecessors(successor activity.Activity[D], predecessors []PredecessorDependency[D]) {
	for _, a := range predecessors {
		if _, ok := d.table[a.Predecessor()]; !ok {
			d.AddRow(a.Predecessor())
		}
	}
	d.table[successor] = predecessors
}

// Returns a new *[depedencyTable]
//
// Added in pdm v1.0.0.
func New[D activity.Data]() *depedencyTable[D] {
	return &depedencyTable[D]{table: map[activity.Activity[D]][]PredecessorDependency[D]{}}
}

var _ DependencyTable[activity.Data] = (*depedencyTable[activity.Data])(nil) //ensures ExampleStruct implements ExampleInterface at compile time
