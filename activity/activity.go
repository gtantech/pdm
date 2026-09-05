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

// TotalFloat implements [Activity].
func (a *activity[D]) TotalFloat() time.Duration {
	return a.Late().Start() - a.Early().Start()
}

func New[D Data](data D) *activity[D] {
	interval := interval.FromStart(0, data.Duration())
	return &activity[D]{data: data, Timestamp: timestamp.New(interval, interval)}
}

// DisplayName implements [Activity].
func (a *activity[D]) Data() D {
	return a.data
}

var _ Activity[Data] = (*activity[Data])(nil) //ensures *activity implements Activity at compile time
