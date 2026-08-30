package enums

type DependencyType string

const (
	SS DependencyType = "SS" // start to start relationship
	FS DependencyType = "FS" // finish to start relationship
	SF DependencyType = "SF" // start to finish relationship
	FF DependencyType = "FF" //finish to finish relationship
)

type Dependency interface {
	Type() DependencyType
}

type relationship struct {
	kind DependencyType
}

func New(kind DependencyType) *relationship {
	return &relationship{kind: kind}
}

func (r *relationship) Type() DependencyType {
	return r.kind
}
