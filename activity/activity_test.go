package activity

import (
	"testing"
	"time"
)

func TestData(t *testing.T) {
	duration, _ := time.ParseDuration("3m")
	activityData := NewData(duration)
	a := New(activityData)

	if got := a.Data(); got != activityData {
		t.Errorf("got %v, want %v", got, activityData)
	}
}
