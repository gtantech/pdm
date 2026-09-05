package dependency

import (
	"fmt"
	"time"

	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/pdm/enums"
)

type Dependency interface {
	Type() enums.DependencyType
	ForwardPassValue(predecessorTimestamp timestamp.Timestamp) time.Duration
	BackwardPassValue(successorTimestamp timestamp.Timestamp) time.Duration
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

// BackwardPassValue implements [Dependency].
func (r *relationship) BackwardPassValue(successorTimestamp timestamp.Timestamp) time.Duration {
	interval := successorTimestamp.Late()
	switch r.kind {
	case enums.FF:
		return interval.Finish()
	case enums.FS:
		return interval.Start()
	case enums.SF:
		return interval.Finish()
	case enums.SS:
		return interval.Start()
	default:
		panic(fmt.Sprintf("unknown enum used: %v", r.kind))
	}
}

// ForwardPassValue implements [Dependency].
func (r *relationship) ForwardPassValue(predecessorTimestamp timestamp.Timestamp) time.Duration {
	interval := predecessorTimestamp.Early()
	switch r.kind {
	case enums.FF:
		return interval.Finish()
	case enums.FS:
		return interval.Finish()
	case enums.SF:
		return interval.Start()
	case enums.SS:
		return interval.Start()
	default:
		panic(fmt.Sprintf("unknown enum used: %v", r.kind))
	}
}
