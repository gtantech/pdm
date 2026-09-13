package models

import (
	"time"
	"uuid"

	"github.com/gtantech/pdm"
	_activity "github.com/gtantech/pdm/activity"
)

type Activity interface {
	DisplayName() string
	Delete() //Delete removes the activity from the associated pdm
	ID() uuid.UUID
	_activity.Data
}

type activity struct {
	_activity.Data
	id          uuid.UUID
	dispName    string
	pdmActivity _activity.Activity[Activity]
	project     Project
}

func (a *activity) ID() uuid.UUID {
	return a.id
}

func (a *activity) DisplayName() string {
	return a.dispName
}

// Delete removes the activity from the associated pdm
func (a *activity) Delete() {
	a.removeFromPDM(a.project)
}

func (a *activity) removeFromPDM(pdm pdm.PDM[Activity]) {
	pdm.RemoveActivity(a.pdmActivity)
	a.pdmActivity = nil
}

var _ Activity = (*activity)(nil) //ensures activity implements Activity at compile time

func NewActivity(dispName string, duration time.Duration, project Project) *activity {
	a := &activity{Data: _activity.NewData(duration), id: uuid.NewV7(), dispName: dispName, project: project}
	a.pdmActivity = project.AddActivity(project.AddActivity(_activity.New(Activity(a))))
	return a
}
