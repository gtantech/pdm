package models

import (
	"uuid"

	"github.com/gtantech/pdm"
)

type Project interface {
	pdm.PDM[Activity]
	ID() uuid.UUID
	DisplayName() string
}

type project struct {
	pdm.PDM[Activity]
	id       uuid.UUID
	dispName string
}

// DisplayName implements [Project].
func (p *project) DisplayName() string {
	return p.dispName
}

// ID implements [Project].
func (p *project) ID() uuid.UUID {
	return p.id
}

var _ Project = (*project)(nil) //ensures project implements Project at compile time
func NewProject(dispName string) *project {
	return &project{PDM: pdm.New[Activity](), id: uuid.NewV7(), dispName: dispName}
}
