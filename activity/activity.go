package activity

import (
	"github.com/gtantech/toposort/graph/vertex"
)

type Activity[A any] interface {
	vertex.Vertex[A]
}

type activity struct {
	vertex.Vertex[*activity]
}

var _ vertex.Vertex[*activity] = (*activity)(nil) //ensures *activity implements vertex.Vertex[*activity] at compile time
var _ Activity[*activity] = (*activity)(nil)      //ensures *activity implements Activity[*activity] at compile time
