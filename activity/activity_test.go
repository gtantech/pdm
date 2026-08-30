package activity

import (
	"testing"
	"time"
)

func TestDisplayName(t *testing.T) {
	dispName := "test_name"
	a := New(dispName, time.Duration(0))

	if got := a.DisplayName(); got != dispName {
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
