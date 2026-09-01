package activity

import (
	"testing"
	"time"

	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/pdm/interval"
)

func TestData(t *testing.T) {
	dispName := "test_name"
	a := New(dispName, time.Duration(0))

	if got := a.Data(); got != dispName {
		t.Errorf("got %v, want %v", got, dispName)
	}
}

func TestDuration(t *testing.T) {
	dispName := "test_name"
	duration, _ := time.ParseDuration("3m")
	a := New(dispName, duration)

	if got := a.Duration(); got != duration {
		t.Errorf("got %v, want %v", got, duration)
	}
}

func TestTimestamp(t *testing.T) {
	dispName := "test_name"
	duration, _ := time.ParseDuration("3m")
	a := New(dispName, duration)
	ts := timestamp.New(interval.New(time.Duration(0), time.Duration(1)), interval.New(time.Duration(0), time.Duration(1)))
	a.UpdateTimestamps(ts)
	if got, want := a.Timestamps(), ts; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
