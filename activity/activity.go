package activity

import (
	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/toposort/graph/vertex"
)

type Activity[D Data] interface {
	vertex.Vertex[Activity[D]]
	Data() D
	Timestamps() timestamp.Timestamp //timestamp relative to the start (time at 0)
	UpdateTimestamps(timestamps timestamp.Timestamp)
}

type activity[D Data] struct {
	vertex.Vertex[Activity[D]]
	data       D
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

func New[D Data](data D) *activity[D] {
	a := &activity[D]{data: data}
	a.Vertex = vertex.New[Activity[D]](a)
	return a
}

// DisplayName implements [Activity].
func (a *activity[D]) Data() D {
	return a.data
}

var _ vertex.Vertex[Activity[Data]] = (*activity[Data])(nil) //ensures *activity implements vertex.Vertex[Activity] at compile time
var _ Activity[Data] = (*activity[Data])(nil)                //ensures *activity implements Activity at compile time
