package models

import (
	"time"
	"uuid"

	"github.com/gtantech/pdm/activity"
)

type ActivityID interface {
	ID
}

type Activity struct {
	activity.Data
	ActivityID
	ProjectID ProjectID
	DispName  string
}

type activityId struct {
	id uuid.UUID
}

// ID implements [ActivityID].
func (a *activityId) ID() uuid.UUID {
	return a.id
}

func NewActivityID(id uuid.UUID) *activityId {
	return &activityId{id: id}
}

var _ ActivityID = (*activityId)(nil) //ensures activityId implements ActivityID at compile time
var _ ActivityID = (*Activity)(nil)   //ensures activity implements ActivityID at compile time

func NewActivity(dispName string, duration time.Duration, projectID ProjectID) *Activity {
	return &Activity{Data: activity.NewData(duration), ActivityID: NewActivityID(uuid.NewV7()), ProjectID: projectID, DispName: dispName}
}
