package pdm

import (
	"slices"
	"testing"
	"time"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/dependency"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/toposort/graph"
	"github.com/gtantech/toposort/graph/vertex"
)

type mockGraph struct {
	graph.Graph[activity.Activity, dependency.Dependency]
	isAddVertexCalled    bool
	isRemoveVertexCalled bool
	isAddEdgeCalled      bool
	isRemoveEdgeCalled   bool
	isVerticesCalled     bool
}

func (g *mockGraph) AddVertex(v vertex.Vertex[activity.Activity]) {
	g.isAddVertexCalled = true
}
func (g *mockGraph) RemoveVertex(v vertex.Vertex[activity.Activity]) {
	g.isRemoveVertexCalled = true
}
func (g *mockGraph) AddEdge(value dependency.Dependency, origin vertex.Vertex[activity.Activity], destination vertex.Vertex[activity.Activity]) {
	g.isAddEdgeCalled = true
}
func (g *mockGraph) RemoveEdge(origin vertex.Vertex[activity.Activity], destination vertex.Vertex[activity.Activity]) {
	g.isRemoveEdgeCalled = true
}
func (g *mockGraph) Vertices() func(yield func(vertex.Vertex[activity.Activity]) bool) {
	g.isVerticesCalled = true
	return g.Graph.Vertices()
}

func TestPDMDisplayName(t *testing.T) {
	dispName := "test_name"
	p := New(dispName)

	if got := p.DisplayName(); got != dispName {
		t.Errorf("got %v, want %v", got, dispName)
	}
}

func TestFields(t *testing.T) {
	dispName := "test_name"
	p := New(dispName)
	mockGraph := &mockGraph{Graph: graph.New[activity.Activity, dependency.Dependency]()}
	p.graph = mockGraph
	a := activity.New("a", time.Duration(0))
	b := activity.New("b", time.Duration(0))
	p.AddActivity(a)
	if !mockGraph.isAddVertexCalled {
		t.Errorf("graph.AddVertex not called")
	}
	p.RemoveActivity(a)
	if !mockGraph.isRemoveVertexCalled {
		t.Errorf("graph.RemoveVertex not called")
	}
	p.AddDependency(a, b, dependency.New(enums.FS))
	if !mockGraph.isAddEdgeCalled {
		t.Errorf("graph.AddEdge not called")
	}
	p.RemoveDependency(a, b)
	if !mockGraph.isRemoveEdgeCalled {
		t.Errorf("graph.RemoveEdge not called")
	}
	p.Activities()
	if !mockGraph.isVerticesCalled {
		t.Errorf("graph.Vertices not called")
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

	p.AddDependency(A, B, dependency.New(enums.FS))
	p.AddDependency(A, C, dependency.New(enums.FS))
	p.AddDependency(B, D, dependency.New(enums.FS))
	p.AddDependency(C, E, dependency.New(enums.FS))
	p.AddDependency(D, F, dependency.New(enums.FS))
	p.AddDependency(E, F, dependency.New(enums.FS))

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

func TestStartingPredecessorActivities(t *testing.T) {
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

	p.AddDependency(A, B, dependency.New(enums.FS))
	p.AddDependency(A, C, dependency.New(enums.FS))
	p.AddDependency(B, D, dependency.New(enums.FS))
	p.AddDependency(C, E, dependency.New(enums.FS))
	p.AddDependency(D, F, dependency.New(enums.FS))
	p.AddDependency(E, F, dependency.New(enums.FS))

	spa := slices.Collect(p.StartingPredecessorActivities())

	if want := []vertex.Vertex[activity.Activity]{A}; !slices.Equal(spa, want) {
		t.Errorf("got %v, want %v", spa, want)
	}
}
