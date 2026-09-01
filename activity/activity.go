package activity

import (
	"time"

	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/toposort/graph/vertex"
)

type Activity[D any] interface {
	vertex.Vertex[Activity[D]]
	Duration() time.Duration
	DisplayName() string
	Timestamps() timestamp.Timestamp //timestamp relative to the start (time at 0)
	UpdateTimestamps(timestamps timestamp.Timestamp)
}

type activity[D any] struct {
	vertex.Vertex[Activity[D]]
	name       string
	duration   time.Duration
	timestamps timestamp.Timestamp
}

// UpdateTimestamps implements [Activity].
func (a *activity[D]) UpdateTimestamps(timestamps timestamp.Timestamp) {
	a.timestamps = timestamps
}

// Timestamps implements [Activity].
func (a *activity[D]) Timestamps() timestamp.Timestamp {
	return a.timestamps
}

func New[D any](name string, duration time.Duration) *activity[D] {
	a := &activity[D]{name: name, duration: duration}
	a.Vertex = vertex.New[Activity[D]](a)
	return a
}

// DisplayName implements [Activity].
func (a *activity[D]) DisplayName() string {
	return a.name
}

// Duration implements [Activity].
func (a *activity[D]) Duration() time.Duration {
	return a.duration
}

var _ vertex.Vertex[Activity[int]] = (*activity[int])(nil) //ensures *activity implements vertex.Vertex[Activity] at compile time
var _ Activity[int] = (*activity[int])(nil)                //ensures *activity implements Activity at compile time
