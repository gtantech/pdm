package table

import (
	"maps"

	"github.com/gtantech/pdm/activity"
)

type DependencyTable[D activity.Data] interface {
	AddRow(successor activity.Activity[D])
	GetRow(successor activity.Activity[D]) ([]activity.Activity[D], bool)
	GetActivities() func(yield func(activity.Activity[D]) bool)
	UpdatePredecessors(successor activity.Activity[D], predecessors []activity.Activity[D])
	DeleteRow(successor activity.Activity[D])
}

type depedencyTable[D activity.Data] struct {
	table map[activity.Activity[D]][]activity.Activity[D]
}

// GetActivities implements [DependencyTable].
func (d *depedencyTable[D]) GetActivities() func(yield func(activity.Activity[D]) bool) {
	return maps.Keys(d.table)
}

// GetRow implements [DependencyTable].
func (d *depedencyTable[D]) GetRow(successor activity.Activity[D]) ([]activity.Activity[D], bool) {
	predecessors, ok := d.table[successor]
	return predecessors, ok
}

// AddRow implements [DependencyTable].
func (d *depedencyTable[D]) AddRow(successor activity.Activity[D]) {
	d.table[successor] = []activity.Activity[D]{}
}

// DeleteRow implements [DependencyTable].
func (d *depedencyTable[D]) DeleteRow(successor activity.Activity[D]) {
	delete(d.table, successor)
}

// UpdatePredecessors implements [DependencyTable].
func (d *depedencyTable[D]) UpdatePredecessors(successor activity.Activity[D], predecessors []activity.Activity[D]) {
	for _, a := range predecessors {
		if _, ok := d.table[a]; !ok {
			d.AddRow(a)
		}
	}
	d.table[successor] = predecessors
}

func New[D activity.Data]() *depedencyTable[D] {
	return &depedencyTable[D]{table: map[activity.Activity[D]][]activity.Activity[D]{}}
}

var _ DependencyTable[activity.Data] = (*depedencyTable[activity.Data])(nil) //ensures ExampleStruct implements ExampleInterface at compile time
