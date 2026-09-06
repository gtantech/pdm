package dependency

import (
	"testing"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/pdm/relationship"
)

func TestNew(t *testing.T) {
	A := activity.New(activity.NewData(0))
	B := activity.New(activity.NewData(0))
	relationship := relationship.New(enums.FS)
	d := New(A, B, relationship)

	if got, want := d.Predecessor(), A; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := d.Successor(), B; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := d.DependsVia(), relationship; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

}
