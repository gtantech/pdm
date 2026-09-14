package models

import (
	"uuid"
)

type ProjectID interface {
	ID
}

type projectId struct {
	id uuid.UUID
}

func NewProjectID(id uuid.UUID) *projectId {
	return &projectId{id: id}
}

// ID implements [ProjectID].
func (a *projectId) ID() uuid.UUID {
	return a.id
}

type Project struct {
	ProjectID
	DispName string
}

func NewProject(dispName string) *Project {
	return &Project{ProjectID: NewProjectID(uuid.NewV7()), DispName: dispName}
}

var _ ProjectID = (*projectId)(nil) //ensures projectId implements ProjectID at compile time
var _ ProjectID = (*Project)(nil)   //ensures Project implements ProjectID at compile time
