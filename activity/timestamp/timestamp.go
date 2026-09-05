package timestamp

import "github.com/gtantech/pdm/interval"

type Timestamp interface {
	Early() interval.Interval
	Late() interval.Interval
	UpdateEarly(interval interval.Interval)
	UpdateLate(interval interval.Interval)
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

// UpdateEarly implements [Activity].
func (d *timestamp) UpdateEarly(interval interval.Interval) {
	d.early = interval
}

// UpdateLate implements [Activity].
func (d *timestamp) UpdateLate(interval interval.Interval) {
	d.late = interval
}

var _ Timestamp = (*timestamp)(nil) //ensures duration implements Duration at compile time
