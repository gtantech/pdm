package pdm

import (
	"slices"
	"testing"
	"time"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/toposort/graph/vertex"
)

func TestPDMDisplayName(t *testing.T) {
	dispName := "test_name"
	p := New(dispName)

	if got := p.DisplayName(); got != dispName {
		t.Errorf("got %v, want %v", got, dispName)
	}
}

func TestPDMCreation(t *testing.T) {
	dispName := "test_name"
	p := New(dispName)

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New("A", time.Minute*5)
	B := activity.New("B", time.Minute*4)
	C := activity.New("C", time.Minute*5)
	D := activity.New("D", time.Minute*6)
	E := activity.New("E", time.Minute*3)
	F := activity.New("F", time.Minute*4)

	p.AddDependency(A, B, enums.New(enums.FS))
	p.AddDependency(A, C, enums.New(enums.FS))
	p.AddDependency(B, D, enums.New(enums.FS))
	p.AddDependency(C, E, enums.New(enums.FS))
	p.AddDependency(D, F, enums.New(enums.FS))
	p.AddDependency(E, F, enums.New(enums.FS))

	for v := range p.graph.IncomingVertices(A) {
		t.Errorf("got %v but %v does not have any incoming vertices", v.Value().DisplayName(), A.DisplayName())
	}
	if v, want := slices.Collect(p.graph.IncomingVertices(B)), []vertex.Vertex[activity.Activity]{A}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.IncomingVertices(C)), []vertex.Vertex[activity.Activity]{A}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.IncomingVertices(D)), []vertex.Vertex[activity.Activity]{B}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.IncomingVertices(E)), []vertex.Vertex[activity.Activity]{C}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want, want2 := slices.Collect(p.graph.IncomingVertices(F)), []vertex.Vertex[activity.Activity]{D, E}, []vertex.Vertex[activity.Activity]{E, D}; (len(v) != 2) && (slices.Equal(v, want) || slices.Equal(v, want2)) {
		t.Errorf("got %v want %v or %v", v, want, want2)
	}

	if v, want, want2 := slices.Collect(p.graph.OutgoingVertices(A)), []vertex.Vertex[activity.Activity]{B, C}, []vertex.Vertex[activity.Activity]{C, B}; (len(v) != 2) && (slices.Equal(v, want) || slices.Equal(v, want2)) {
		t.Errorf("got %v want %v or %v", v, want, want2)
	}
	if v, want := slices.Collect(p.graph.OutgoingVertices(B)), []vertex.Vertex[activity.Activity]{D}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.OutgoingVertices(C)), []vertex.Vertex[activity.Activity]{E}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.OutgoingVertices(D)), []vertex.Vertex[activity.Activity]{F}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.OutgoingVertices(E)), []vertex.Vertex[activity.Activity]{F}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	for v := range p.graph.OutgoingVertices(F) {
		t.Errorf("got %v but %v does not have any outgoing vertices", v.Value().DisplayName(), F.DisplayName())
	}
}
