package db

import "uuid"

type Project struct {
	ID          uuid.UUID
	DisplayName string
}
