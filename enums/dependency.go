package enums

type DependencyType string

const (
	StartToStart   DependencyType = "SS"
	FinishToStart  DependencyType = "FS"
	StartToFinish  DependencyType = "SF"
	FinishToFinish DependencyType = "FF"
)

type Dependency interface {
	Type() DependencyType
}
