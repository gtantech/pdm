package activity

import (
	"github.com/gtantech/pdm/activity/timestamp"
)

type Activity[D Data] interface {
	Data() D
	timestamp.Timestamp //timestamp relative to the start (time at 0)
	UpdateTimestamps(timestamps timestamp.Timestamp)
}

type activity[D Data] struct {
	data                D
	timestamp.Timestamp //timestamp relative to the start (time at 0)
}

// UpdateTimestamps implements [Activity].
func (a *activity[D]) UpdateTimestamps(timestamps timestamp.Timestamp) {
	a.Timestamp = timestamps
}

func New[D Data](data D) *activity[D] {
	return &activity[D]{data: data}
}

// DisplayName implements [Activity].
func (a *activity[D]) Data() D {
	return a.data
}

var _ Activity[Data] = (*activity[Data])(nil) //ensures *activity implements Activity at compile time
