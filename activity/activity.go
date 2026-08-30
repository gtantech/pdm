package activity

import (
	"github.com/gtantech/toposort/graph/vertex"
)

type Activity interface {
	vertex.Vertex[Activity]
	DisplayName() string
}

type activity struct {
	vertex.Vertex[Activity]
	name string
}

// DisplayName implements [Activity].
func (a *activity) DisplayName() string {
	return a.name
}

var _ vertex.Vertex[Activity] = (*activity)(nil) //ensures *activity implements vertex.Vertex[Activity] at compile time
var _ Activity = (*activity)(nil)                //ensures *activity implements Activity at compile time
