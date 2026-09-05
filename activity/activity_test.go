package activity

import (
	"testing"
	"time"

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

func TestTotalFloat(t *testing.T) {
	duration, _ := time.ParseDuration("3m")
	activityData := NewData(duration)
	a := New(activityData)
	a.UpdateEarly(interval.New(5*time.Minute, 8*time.Minute))
	a.UpdateLate(interval.New(6*time.Minute, 9*time.Minute))

	if got, want := a.TotalFloat(), 1*time.Minute; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
