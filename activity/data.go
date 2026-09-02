package activity

import "time"

type Data interface {
	Duration() time.Duration
}

type data struct {
	duration time.Duration
}

func NewData(duration time.Duration) *data {
	return &data{duration: duration}
}

func (d *data) Duration() time.Duration {
	return d.duration
}
