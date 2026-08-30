package dependency

import (
	"github.com/gtantech/pdm/enums"
)

type Dependency interface {
	Type() enums.DependencyType
}

type relationship struct {
	kind enums.DependencyType
}

func New(kind enums.DependencyType) *relationship {
	return &relationship{kind: kind}
}

func (r *relationship) Type() enums.DependencyType {
	return r.kind
}
