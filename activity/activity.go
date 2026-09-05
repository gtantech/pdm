package activity

import (
	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/pdm/interval"
)

type Activity[D Data] interface {
	Data() D
	timestamp.Timestamp //timestamp relative to the start (time at 0)
}

type activity[D Data] struct {
	data                D
	timestamp.Timestamp //timestamp relative to the start (time at 0)
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
