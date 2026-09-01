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

func New(start time.Duration, finish time.Duration) *interval {
	return &interval{start: start, finish: finish}
}

// Finish implements [ActivityInterval].
func (i *interval) Finish() time.Duration {
	return i.finish
}

// Start implements [ActivityInterval].
func (i *interval) Start() time.Duration {
	return i.start
}

var _ Interval = (*interval)(nil) //ensures interval implements Interval at compile time
