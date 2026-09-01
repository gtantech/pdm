package data

import (
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	duration := time.Minute * 3

	d := New(duration)

	if got, want := d.Duration(), duration; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
