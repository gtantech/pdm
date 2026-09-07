package activity

import "time"

type Data interface {
	Duration() time.Duration
}

type data struct {
	duration time.Duration
}

// NewData returns a new *[data].
//
// Added in pdm v1.0.0.
func NewData(duration time.Duration) *data {
	return &data{duration: duration}
}

// Duration implements [Data]. Duration returns the duration in a *[data].
//
// Added in pdm v1.0.0.
func (d *data) Duration() time.Duration {
	return d.duration
}
