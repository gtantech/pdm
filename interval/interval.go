package interval

import "time"

type Interval interface {
	Start() time.Duration
	Finish() time.Duration
}

type interval struct {
	start  time.Duration
	finish time.Duration
}

// New returns a new *[interval]
func New(start time.Duration, finish time.Duration) *interval {
	return &interval{start: start, finish: finish}
}

// FromStart returns a new *[interval] whose start/finish values are calculated from start and duration.
//
// Added in pdm v1.0.0.
func FromStart(start time.Duration, duration time.Duration) *interval {
	return &interval{start: start, finish: start + duration}
}

// FromFinish returns a new *[interval] whose start/finish values are calculated from finish and duration.
//
// Added in pdm v1.0.0.
func FromFinish(finish time.Duration, duration time.Duration) *interval {
	return &interval{start: finish - duration, finish: finish}
}

// Finish implements [Interval]. Finish returns the finish value in i *[interval].
//
// Added in pdm v1.0.0.
func (i *interval) Finish() time.Duration {
	return i.finish
}

// Start implements [Interval]. Start returns the start value in i *[interval].
//
// Added in pdm v1.0.0.
func (i *interval) Start() time.Duration {
	return i.start
}

var _ Interval = (*interval)(nil) //ensures interval implements Interval at compile time
