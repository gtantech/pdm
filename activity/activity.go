package activity

import (
	"time"

	"github.com/gtantech/toposort/graph/vertex"
)

type Activity interface {
	vertex.Vertex[Activity]
	Duration() time.Duration
	DisplayName() string
}

type activity struct {
	vertex.Vertex[Activity]
	name     string
	duration time.Duration
}

// DisplayName implements [Activity].
func (a *activity) DisplayName() string {
	return a.name
}

// Duration implements [Activity].
func (a *activity) Duration() time.Duration {
	return a.duration
}

var _ vertex.Vertex[Activity] = (*activity)(nil) //ensures *activity implements vertex.Vertex[Activity] at compile time
var _ Activity = (*activity)(nil)                //ensures *activity implements Activity at compile time
