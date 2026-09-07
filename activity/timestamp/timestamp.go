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

// New returns a new *[timestamp].
//
// Added in pdm v1.0.0.
func New(early interval.Interval, late interval.Interval) *timestamp {
	return &timestamp{early: early, late: late}
}

// Early implements [Timestamp]. Early returns the early value [interval.Interval] in d *[timestamp].
//
// Added in pdm v1.0.0.
func (d *timestamp) Early() interval.Interval {
	return d.early
}

// Late implements [Timestamp]. Late returns the late value [interval.Interval] in d *[timestamp].
//
// Added in pdm v1.0.0.
func (d *timestamp) Late() interval.Interval {
	return d.late
}

// UpdateEarly implements [Timestamp]. UpdateEarly replaces the early interval in d *[timestamp] with the one specified.
//
// Added in pdm v1.0.0.
func (d *timestamp) UpdateEarly(interval interval.Interval) {
	d.early = interval
}

// UpdateLate implements [Timestamp]. UpdateLate replaces the late interval in d *[timestamp] with the one specified.
//
// Added in pdm v1.0.0.
func (d *timestamp) UpdateLate(interval interval.Interval) {
	d.late = interval
}

var _ Timestamp = (*timestamp)(nil) //ensures duration implements Duration at compile time
