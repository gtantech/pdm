package timestamp

import (
	"testing"
	"time"

	"github.com/gtantech/pdm/v2/interval"
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
