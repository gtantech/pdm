package relationship

import (
	"fmt"
	"time"

	"github.com/gtantech/pdm/activity/timestamp"
	"github.com/gtantech/pdm/enums"
)

type Relationship interface {
	Type() enums.RelationshipType
	ForwardPassValue(predecessorTimestamp timestamp.Timestamp) time.Duration
	BackwardPassValue(successorTimestamp timestamp.Timestamp) time.Duration
}

type relationship struct {
	kind enums.RelationshipType
	lag  time.Duration
}

func New(kind enums.RelationshipType) *relationship {
	return &relationship{kind: kind, lag: time.Duration(0)}
}

func NewWithLag(kind enums.RelationshipType, lag time.Duration) *relationship {
	return &relationship{kind: kind, lag: lag}
}

func (r *relationship) Type() enums.RelationshipType {
	return r.kind
}

// BackwardPassValue implements [Relationship].
func (r *relationship) BackwardPassValue(successorTimestamp timestamp.Timestamp) time.Duration {
	interval := successorTimestamp.Late()
	switch r.kind {
	case enums.FF:
		return interval.Finish() - r.lag
	case enums.FS:
		return interval.Start() - r.lag
	case enums.SF:
		return interval.Finish() - r.lag
	case enums.SS:
		return interval.Start() - r.lag
	default:
		panic(fmt.Sprintf("unknown enum used: %v", r.kind))
	}
}

// ForwardPassValue implements [Relationship].
func (r *relationship) ForwardPassValue(predecessorTimestamp timestamp.Timestamp) time.Duration {
	interval := predecessorTimestamp.Early()
	switch r.kind {
	case enums.FF:
		return interval.Finish() + r.lag
	case enums.FS:
		return interval.Finish() + r.lag
	case enums.SF:
		return interval.Start() + r.lag
	case enums.SS:
		return interval.Start() + r.lag
	default:
		panic(fmt.Sprintf("unknown enum used: %v", r.kind))
	}
}
