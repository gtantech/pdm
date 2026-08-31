package interval

import "time"

type Interval interface {
	Start() time.Time
	Finish() time.Time
}

type interval struct {
	start  time.Time
	finish time.Time
}

func New(start time.Time, finish time.Time) *interval {
	return &interval{start: start, finish: finish}
}

// Finish implements [ActivityInterval].
func (i *interval) Finish() time.Time {
	return i.finish
}

// Start implements [ActivityInterval].
func (i *interval) Start() time.Time {
	return i.start
}

var _ Interval = (*interval)(nil) //ensures interval implements Interval at compile time
