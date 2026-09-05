package dependency

import (
	"testing"
	"time"

	"github.com/gtantech/pdm/enums"
)

func TestType(t *testing.T) {
	kind := enums.FS
	d := New(kind)
	if got := d.Type(); got != kind {
		t.Errorf("got %v, want %v", got, kind)
	}
}

func TestLag(t *testing.T) {
	lag := 5 * time.Minute
	d := NewWithLag(enums.FS, lag)
	if got := d.lag; got != lag {
		t.Errorf("got %v, want %v", got, lag)
	}
}
