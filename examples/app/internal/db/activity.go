package db

import (
	"time"
	"uuid"
)

type Activity struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}
