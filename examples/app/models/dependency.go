package models

import (
	"uuid"

	"github.com/gtantech/pdm/enums"
)

type DependencyID interface {
	ID
}

type dependencyId struct {
	id uuid.UUID
}

func NewDependencyID(id uuid.UUID) *dependencyId {
	return &dependencyId{id: id}
}

// ID implements [DependencyID].
func (d *dependencyId) ID() uuid.UUID {
	return d.id
}

type Dependency struct {
	DependencyID
	ProjectID ProjectID
	enums.RelationshipType
	Predecessor ActivityID
	Successor   ActivityID
}

func NewDependency(predecessor ActivityID, successor ActivityID, dependsVia enums.RelationshipType, inProject ProjectID) *Dependency {
	return &Dependency{DependencyID: NewDependencyID(uuid.NewV7()), RelationshipType: dependsVia, Predecessor: predecessor, Successor: successor, ProjectID: inProject}
}

var _ DependencyID = (*dependencyId)(nil) //ensures dependencyId implements DependencyID at compile time
var _ DependencyID = (*Dependency)(nil)   //ensures Dependency implements DependencyID at compile time
