package activity

import (
	"testing"
	"time"

	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/pdm/interval"
)

func TestData(t *testing.T) {
	duration, _ := time.ParseDuration("3m")
	activityData := NewData(duration)
	a := New(activityData)

	if got := a.Data(); got != activityData {
		t.Errorf("got %v, want %v", got, activityData)
	}
}

func TestTimestamp(t *testing.T) {
	duration, _ := time.ParseDuration("3m")
	a := New(NewData(duration))
	ts := timestamp.New(interval.New(time.Duration(0), time.Duration(1)), interval.New(time.Duration(0), time.Duration(1)))
	a.UpdateTimestamps(ts)
	if got, want := a.Timestamps(), ts; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
