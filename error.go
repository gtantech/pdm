package pdm

import (
	"fmt"
)

type CycleDetectedError[V comparable, E any] struct {
	Relationship E
	Predecessor  V
	Successor    V
}

func (e *CycleDetectedError[V, E]) Error() string {
	return fmt.Sprintf("encountered cycle in graph for edge: %v from %v to %v", e.Relationship, e.Predecessor, e.Successor)
}
