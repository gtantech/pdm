package interval

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	start := time.Date(2026, 8, 30, 22, 21, 0, 0, time.Local)
	finish := time.Date(2026, 8, 31, 22, 21, 0, 0, time.Local)
	i := New(start, finish)

	if got, want := i.Start(), start; got != want {
		t.Errorf("got %v want %v", got, want)
	}

	if got, want := i.Finish(), finish; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
