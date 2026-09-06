package table

import (
	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/relationship"
)

type PredecessorDependency[D activity.Data] interface {
	Predecessor() activity.Activity[D]
	DependsVia() relationship.Relationship
}

type predecessorDependency[D activity.Data] struct {
	predecessor activity.Activity[D]
	dependsVia  relationship.Relationship
}

// DependsVia implements [PredecessorDependency].
func (p *predecessorDependency[D]) DependsVia() relationship.Relationship {
	return p.dependsVia
}

// Predecessor implements [PredecessorDependency].
func (p *predecessorDependency[D]) Predecessor() activity.Activity[D] {
	return p.predecessor
}

var _ PredecessorDependency[activity.Data] = (*predecessorDependency[activity.Data])(nil) //ensures ExampleStruct implements ExampleInterface at compile time
func NewPredecessorDependency[D activity.Data](predecessor activity.Activity[D], dependsVia relationship.Relationship) *predecessorDependency[D] {
	return &predecessorDependency[D]{predecessor: predecessor, dependsVia: dependsVia}
}
