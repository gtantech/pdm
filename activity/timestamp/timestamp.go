package timestamp

import "github.com/gtantech/pdm/v2/interval"

type Timestamp interface {
	Early() interval.Interval
	Late() interval.Interval
}

type timestamp struct {
	early interval.Interval
	late  interval.Interval
}

func New(early interval.Interval, late interval.Interval) *timestamp {
	return &timestamp{early: early, late: late}
}

// Early implements [Timestamp].
func (d *timestamp) Early() interval.Interval {
	return d.early
}

// Late implements [Timestamp].
func (d *timestamp) Late() interval.Interval {
	return d.late
}

var _ Timestamp = (*timestamp)(nil) //ensures duration implements Duration at compile time
