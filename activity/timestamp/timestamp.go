package duration

import "github.com/gtantech/pdm/interval"

type Duration interface {
	Early() interval.Interval
	Late() interval.Interval
}

type duration struct {
	early interval.Interval
	late  interval.Interval
}

func New(early interval.Interval, late interval.Interval) *duration {
	return &duration{early: early, late: late}
}

// Early implements [Duration].
func (d *duration) Early() interval.Interval {
	return d.early
}

// Late implements [Duration].
func (d *duration) Late() interval.Interval {
	return d.late
}

var _ Duration = (*duration)(nil) //ensures duration implements Duration at compile time
