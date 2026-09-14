package models

import "uuid"

type ID interface {
	ID() uuid.UUID
}
