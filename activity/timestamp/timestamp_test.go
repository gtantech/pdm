package timestamp

import (
	"testing"
	"time"

	"github.com/gtantech/pdm/interval"
)

func TestNew(t *testing.T) {
	early := interval.New(time.Duration(0), time.Duration(1))
	late := interval.New(time.Duration(0), time.Duration(1))
	ts := New(early, late)

	if got, want := ts.Early(), early; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := ts.Late(), late; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTimestampUpdate(t *testing.T) {
	early := interval.New(time.Duration(4), time.Duration(6))
	late := interval.New(time.Duration(5), time.Duration(8))
	ts := New(interval.New(time.Duration(0), time.Duration(1)), interval.New(time.Duration(0), time.Duration(1)))
	ts.UpdateEarly(early)
	ts.UpdateLate(late)
	if got, want := ts.Early(), early; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := ts.Late(), late; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
