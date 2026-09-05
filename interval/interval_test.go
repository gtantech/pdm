package interval

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	start := time.Duration(0 * time.Minute)
	finish := time.Duration(1 * time.Minute)
	i := New(start, finish)

	if got, want := i.Start(), start; got != want {
		t.Errorf("got %v want %v", got, want)
	}

	if got, want := i.Finish(), finish; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestFromStart(t *testing.T) {
	start := time.Duration(1 * time.Minute)
	duration := time.Duration(3 * time.Minute)
	finish := time.Duration(4 * time.Minute)
	i := FromStart(start, duration)

	if got, want := i.Start(), start; got != want {
		t.Errorf("got %v want %v", got, want)
	}

	if got, want := i.Finish(), finish; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestFromFinish(t *testing.T) {
	start := time.Duration(1 * time.Minute)
	duration := time.Duration(3 * time.Minute)
	finish := time.Duration(4 * time.Minute)
	i := FromFinish(finish, duration)

	if got, want := i.Start(), start; got != want {
		t.Errorf("got %v want %v", got, want)
	}

	if got, want := i.Finish(), finish; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
