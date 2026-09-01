package data

import "time"

type ActivityData interface {
	Duration() time.Duration
}

type activityData struct {
	duration time.Duration
}

func New(duration time.Duration) *activityData {
	return &activityData{duration: duration}
}

func (d *activityData) Duration() time.Duration {
	return d.duration
}
