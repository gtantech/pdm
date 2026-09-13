package models

import (
	"time"
	"uuid"

	_activity "github.com/gtantech/pdm/activity"
)

type Activity interface {
	DisplayName() string
	ID() uuid.UUID
	_activity.Data
}

type activity struct {
	_activity.Data
	id       uuid.UUID
	dispName string
}

func (a *activity) ID() uuid.UUID {
	return a.id
}

func (a *activity) DisplayName() string {
	return a.dispName
}

var _ Activity = (*activity)(nil) //ensures activity implements Activity at compile time

func NewActivity(dispName string, duration time.Duration) *activity {
	return &activity{Data: _activity.NewData(duration), id: uuid.NewV7(), dispName: dispName}
}
