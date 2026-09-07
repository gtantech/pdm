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

// New returns a new *[relationship]
//
// Added in pdm v1.0.0.
func New(kind enums.RelationshipType) *relationship {
	return &relationship{kind: kind, lag: time.Duration(0)}
}

// NewWithLag returns a new *[relationship] with a specified lag. Lead times can be specified by a negative lag value.
//
// Added in pdm v1.0.0.
func NewWithLag(kind enums.RelationshipType, lag time.Duration) *relationship {
	return &relationship{kind: kind, lag: lag}
}

// Type implements [Relationship]. Type returns the logical relationship type specified in [enums.RelationshipType].
// Four relationship types exists : Finish to Start (FS), Start to Start (SS), Start to Finish (SF), Finish to Finish (FF).
//
// Added in pdm v1.0.0.
func (r *relationship) Type() enums.RelationshipType {
	return r.kind
}

// BackwardPassValue implements [Relationship]. BackwardPassValue will return either the late finish or late start value of the successor, depending on the relationship type in r *[relationship].
//
// If the relationship is FF or SF, it will return the finish value minus lag, and if it is FS or SS, it will return the start value minus lag.
//
// Added in pdm v1.0.0.
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

// ForwardPassValue implements [Relationship]. ForwardPassValue will return either the early finish or early start value of the successor, depending on the relationship type in r *[relationship].
//
// If the relationship is FF or FS, it will return the finish value plus lag, and if it is SF or SS, it will return the start value plus lag.
//
// Added in pdm v1.0.0.
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
