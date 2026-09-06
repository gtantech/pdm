package table

import (
	"testing"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/pdm/relationship"
)

func TestFields(t *testing.T) {
	A := activity.New(activity.Data(activity.NewData(0)))
	dependsVia := relationship.New(enums.FS)
	pd := NewPredecessorDependency(A, dependsVia)
	if got, want := pd.DependsVia(), dependsVia; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := pd.Predecessor(), A; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
