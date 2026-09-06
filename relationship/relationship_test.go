package relationship

import (
	"testing"
	"time"

	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/pdm/interval"
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

func TestForwardsPassValue(t *testing.T) {
	early := interval.New(time.Duration(0), time.Duration(1))
	late := interval.New(time.Duration(0), time.Duration(1))
	ts := timestamp.New(early, late)

	if got, want := New(enums.FS).ForwardPassValue(ts), early.Finish(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := New(enums.SS).ForwardPassValue(ts), early.Start(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := New(enums.FF).ForwardPassValue(ts), early.Finish(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := New(enums.SF).ForwardPassValue(ts), early.Start(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBackwardsPassValue(t *testing.T) {
	early := interval.New(time.Duration(0), time.Duration(1))
	late := interval.New(time.Duration(0), time.Duration(1))
	ts := timestamp.New(early, late)

	if got, want := New(enums.FS).BackwardPassValue(ts), early.Start(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := New(enums.SS).BackwardPassValue(ts), early.Start(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := New(enums.FF).BackwardPassValue(ts), early.Finish(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := New(enums.SF).BackwardPassValue(ts), early.Finish(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPanicUnknownEnumForwardsPass(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("failed to get edge did not panic")
		}
	}()
	early := interval.New(time.Duration(0), time.Duration(1))
	late := interval.New(time.Duration(0), time.Duration(1))
	ts := timestamp.New(early, late)
	New("test_wrong_enum").ForwardPassValue(ts) //this should panic
}

func TestPanicUnknownEnumBackwardPass(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("failed to get edge did not panic")
		}
	}()
	early := interval.New(time.Duration(0), time.Duration(1))
	late := interval.New(time.Duration(0), time.Duration(1))
	ts := timestamp.New(early, late)
	New("test_wrong_enum").BackwardPassValue(ts) //this should panic
}
