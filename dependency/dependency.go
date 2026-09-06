package dependency

import (
	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/relationship"
)

type Dependency[D activity.Data] interface {
	Predecessor() activity.Activity[D]
	Successor() activity.Activity[D]
	DependsVia() relationship.Relationship
}

type dependency[D activity.Data] struct {
	predecessor activity.Activity[D]
	successor   activity.Activity[D]
	dependsVia  relationship.Relationship
}

// DependsVia implements [Dependency].
func (d *dependency[D]) DependsVia() relationship.Relationship {
	return d.dependsVia
}

// Predecessor implements [Dependency].
func (d *dependency[D]) Predecessor() activity.Activity[D] {
	return d.predecessor
}

// Successor implements [Dependency].
func (d *dependency[D]) Successor() activity.Activity[D] {
	return d.successor
}

var _ Dependency[activity.Data] = (*dependency[activity.Data])(nil) //ensures ExampleStruct implements ExampleInterface at compile time

func New[D activity.Data](predecessor activity.Activity[D], successor activity.Activity[D], dependsVia relationship.Relationship) *dependency[D] {
	return &dependency[D]{predecessor: predecessor, successor: successor, dependsVia: dependsVia}
}
