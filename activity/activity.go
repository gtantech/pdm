package activity

import (
	"time"

	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/toposort/graph/vertex"
)

type Activity interface {
	vertex.Vertex[Activity]
	Duration() time.Duration
	DisplayName() string
	Timestamps() timestamp.Timestamp //timestamp relative to the start (time at 0)
	UpdateTimestamps(timestamps timestamp.Timestamp)
}

type activity struct {
	vertex.Vertex[Activity]
	name       string
	duration   time.Duration
	timestamps timestamp.Timestamp
}

// UpdateTimestamps implements [Activity].
func (a *activity) UpdateTimestamps(timestamps timestamp.Timestamp) {
	a.timestamps = timestamps
}

// Timestamps implements [Activity].
func (a *activity) Timestamps() timestamp.Timestamp {
	return a.timestamps
}

func New(name string, duration time.Duration) *activity {
	a := &activity{name: name, duration: duration}
	a.Vertex = vertex.New[Activity](a)
	return a
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
