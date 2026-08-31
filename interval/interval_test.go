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
