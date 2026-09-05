package dependency

import (
	"time"

	"github.com/gtantech/pdm/enums"
)

type Dependency interface {
	Type() enums.DependencyType
}

type relationship struct {
	kind enums.DependencyType
	lag  time.Duration
}

func New(kind enums.DependencyType) *relationship {
	return &relationship{kind: kind, lag: time.Duration(0)}
}

func NewWithLag(kind enums.DependencyType, lag time.Duration) *relationship {
	return &relationship{kind: kind, lag: lag}
}

func (r *relationship) Type() enums.DependencyType {
	return r.kind
}
