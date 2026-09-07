package pdm

import (
	"fmt"
)

// CycleDetectedError indicates that the function called detected a cycle within a graph
//
// Added in pdm v1.0.0.
type CycleDetectedError[V comparable, E any] struct {
	Relationship E
	Predecessor  V
	Successor    V
}

// Error returns the error message in e *[CycleDetectedError]
//
// Added in pdm v1.0.0.
func (e *CycleDetectedError[V, E]) Error() string {
	return fmt.Sprintf("encountered cycle in graph for edge: %v from %v to %v", e.Relationship, e.Predecessor, e.Successor)
}
