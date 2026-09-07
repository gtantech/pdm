package activity

import (
	"time"

	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/pdm/interval"
)

type Activity[D Data] interface {
	Data() D
	timestamp.Timestamp //timestamp relative to the start (time at 0)
	TotalFloat() time.Duration
}

type activity[D Data] struct {
	data                D
	timestamp.Timestamp //timestamp relative to the start (time at 0)
}

// TotalFloat implements [Activity]. TotalFloat returns the float in a *[activity]. This is the amount of time the activity can be delayed without delaying the project end date
//
// Added in pdm v1.0.0.
func (a *activity[D]) TotalFloat() time.Duration {
	return a.Late().Start() - a.Early().Start()
}

// New returns a new *[activity].
//
// Added in pdm v1.0.0.
func New[D Data](data D) *activity[D] {
	interval := interval.FromStart(0, data.Duration())
	return &activity[D]{data: data, Timestamp: timestamp.New(interval, interval)}
}

// Data implements [Activity]. Data returns the associated data in a *[activity].
//
// Added in pdm v1.0.0.
func (a *activity[D]) Data() D {
	return a.data
}

var _ Activity[Data] = (*activity[Data])(nil) //ensures *activity implements Activity at compile time
