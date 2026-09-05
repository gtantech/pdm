package pdm

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/gtantech/pdm/activity"
	"github.com/gtantech/pdm/dependency"
	"github.com/gtantech/pdm/enums"
	"github.com/gtantech/pdm/interval"
	"github.com/gtantech/toposort/v2/graph"
)

type mockActivityData struct {
	activity.Data
	name string
}

func NewMockActivityData(name string, duration time.Duration) *mockActivityData {
	return &mockActivityData{name: name, Data: activity.NewData(duration)}
}

type mockGraph struct {
	graph.Graph[activity.Activity[*mockActivityData], dependency.Dependency]
	isAddVertexCalled    bool
	isRemoveVertexCalled bool
	isAddEdgeCalled      bool
	isRemoveEdgeCalled   bool
	isVerticesCalled     bool
}

func (g *mockGraph) AddVertex(v activity.Activity[*mockActivityData]) {
	g.isAddVertexCalled = true
}
func (g *mockGraph) RemoveVertex(v activity.Activity[*mockActivityData]) {
	g.isRemoveVertexCalled = true
}
func (g *mockGraph) AddEdge(value dependency.Dependency, origin activity.Activity[*mockActivityData], destination activity.Activity[*mockActivityData]) {
	g.isAddEdgeCalled = true
}
func (g *mockGraph) RemoveEdge(origin activity.Activity[*mockActivityData], destination activity.Activity[*mockActivityData]) {
	g.isRemoveEdgeCalled = true
}
func (g *mockGraph) Vertices() func(yield func(activity.Activity[*mockActivityData]) bool) {
	g.isVerticesCalled = true
	return g.Graph.Vertices()
}

type mockGraphGetEdgePanic struct {
	graph.Graph[activity.Activity[*mockActivityData], dependency.Dependency]
}

func (g *mockGraphGetEdgePanic) GetEdgeValue(origin activity.Activity[*mockActivityData], destination activity.Activity[*mockActivityData]) (dependency.Dependency, bool) {
	return nil, false
}

func TestFields(t *testing.T) {
	p := New[*mockActivityData]()
	mockGraph := &mockGraph{Graph: graph.New[activity.Activity[*mockActivityData], dependency.Dependency]()}
	p.graph = mockGraph
	a := activity.New(NewMockActivityData("a", time.Duration(0)))
	b := activity.New(NewMockActivityData("b", time.Duration(0)))
	added := p.AddActivity(a)
	if added != a {
		t.Errorf("expected return value to be %v, got %v", a, added)
	}
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
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*5))
	B := activity.New(NewMockActivityData("B", time.Minute*4))
	C := activity.New(NewMockActivityData("C", time.Minute*5))
	D := activity.New(NewMockActivityData("D", time.Minute*6))
	E := activity.New(NewMockActivityData("E", time.Minute*3))
	F := activity.New(NewMockActivityData("F", time.Minute*4))

	p.AddDependency(A, B, dependency.New(enums.FS))
	p.AddDependency(A, C, dependency.New(enums.FS))
	p.AddDependency(B, D, dependency.New(enums.FS))
	p.AddDependency(C, E, dependency.New(enums.FS))
	p.AddDependency(D, F, dependency.New(enums.FS))
	p.AddDependency(E, F, dependency.New(enums.FS))

	for v := range p.graph.IncomingVertices(A) {
		t.Errorf("got %v but %v does not have any incoming vertices", v.Data(), A.Data())
	}
	if v, want := slices.Collect(p.graph.IncomingVertices(B)), []activity.Activity[*mockActivityData]{A}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.IncomingVertices(C)), []activity.Activity[*mockActivityData]{A}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.IncomingVertices(D)), []activity.Activity[*mockActivityData]{B}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.IncomingVertices(E)), []activity.Activity[*mockActivityData]{C}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want, want2 := slices.Collect(p.graph.IncomingVertices(F)), []activity.Activity[*mockActivityData]{D, E}, []activity.Activity[*mockActivityData]{E, D}; (len(v) != 2) && (slices.Equal(v, want) || slices.Equal(v, want2)) {
		t.Errorf("got %v want %v or %v", v, want, want2)
	}

	if v, want, want2 := slices.Collect(p.graph.OutgoingVertices(A)), []activity.Activity[*mockActivityData]{B, C}, []activity.Activity[*mockActivityData]{C, B}; (len(v) != 2) && (slices.Equal(v, want) || slices.Equal(v, want2)) {
		t.Errorf("got %v want %v or %v", v, want, want2)
	}
	if v, want := slices.Collect(p.graph.OutgoingVertices(B)), []activity.Activity[*mockActivityData]{D}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.OutgoingVertices(C)), []activity.Activity[*mockActivityData]{E}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.OutgoingVertices(D)), []activity.Activity[*mockActivityData]{F}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	if v, want := slices.Collect(p.graph.OutgoingVertices(E)), []activity.Activity[*mockActivityData]{F}; (len(v) != 1) && slices.Equal(v, want) {
		t.Errorf("got %v want %v", v, want)
	}
	for v := range p.graph.OutgoingVertices(F) {
		t.Errorf("got %v but %v does not have any outgoing vertices", v.Data(), F.Data())
	}
}

func TestInitialPredecessorActivities(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*5))
	B := activity.New(NewMockActivityData("B", time.Minute*4))
	C := activity.New(NewMockActivityData("C", time.Minute*5))
	D := activity.New(NewMockActivityData("D", time.Minute*6))
	E := activity.New(NewMockActivityData("E", time.Minute*3))
	F := activity.New(NewMockActivityData("F", time.Minute*4))

	p.AddDependency(A, B, dependency.New(enums.FS))
	p.AddDependency(A, C, dependency.New(enums.FS))
	p.AddDependency(B, D, dependency.New(enums.FS))
	p.AddDependency(C, E, dependency.New(enums.FS))
	p.AddDependency(D, F, dependency.New(enums.FS))
	p.AddDependency(E, F, dependency.New(enums.FS))

	spa := slices.Collect(p.InitialPredecessorActivities())

	if want := []activity.Activity[*mockActivityData]{A}; !slices.Equal(spa, want) {
		t.Errorf("got %v, want %v", spa, want)
	}
}

func TestFinalSuccessorActivities(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*5))
	B := activity.New(NewMockActivityData("B", time.Minute*4))
	C := activity.New(NewMockActivityData("C", time.Minute*5))
	D := activity.New(NewMockActivityData("D", time.Minute*6))
	E := activity.New(NewMockActivityData("E", time.Minute*3))
	F := activity.New(NewMockActivityData("F", time.Minute*4))

	p.AddDependency(A, B, dependency.New(enums.FS))
	p.AddDependency(A, C, dependency.New(enums.FS))
	p.AddDependency(B, D, dependency.New(enums.FS))
	p.AddDependency(C, E, dependency.New(enums.FS))
	p.AddDependency(D, F, dependency.New(enums.FS))
	p.AddDependency(E, F, dependency.New(enums.FS))

	spa := slices.Collect(p.FinalSuccessorActivities())

	if want := []activity.Activity[*mockActivityData]{F}; !slices.Equal(spa, want) {
		t.Errorf("got %v, want %v", spa, want)
	}
}

func TestUpdateActivityTimestamps(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*3))
	B := activity.New(NewMockActivityData("B", time.Minute*4))
	C := activity.New(NewMockActivityData("C", time.Minute*2))
	D := activity.New(NewMockActivityData("D", time.Minute*5))
	E := activity.New(NewMockActivityData("E", time.Minute*1))
	F := activity.New(NewMockActivityData("F", time.Minute*2))
	G := activity.New(NewMockActivityData("G", time.Minute*4))
	H := activity.New(NewMockActivityData("H", time.Minute*3))

	p.AddDependency(A, B, dependency.New(enums.FS))
	p.AddDependency(A, C, dependency.New(enums.FS))
	p.AddDependency(B, D, dependency.New(enums.FS))
	p.AddDependency(C, E, dependency.New(enums.FS))
	p.AddDependency(C, F, dependency.New(enums.FS))
	p.AddDependency(D, G, dependency.New(enums.FS))
	p.AddDependency(E, G, dependency.New(enums.FS))
	p.AddDependency(F, H, dependency.New(enums.FS))
	p.AddDependency(G, H, dependency.New(enums.FS))

	err := p.UpdateActivityTimestamps()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got, want := A.Early(), interval.New(time.Duration(0), time.Duration(time.Minute*3)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := A.Late(), interval.New(time.Duration(0), time.Duration(time.Minute*3)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := B.Early(), interval.New(time.Duration(3*time.Minute), time.Duration(time.Minute*7)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := B.Late(), interval.New(time.Duration(3*time.Minute), time.Duration(time.Minute*7)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := C.Early(), interval.New(time.Duration(3*time.Minute), time.Duration(5*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := C.Late(), interval.New(time.Duration(9*time.Minute), time.Duration(11*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := D.Early(), interval.New(time.Duration(7*time.Minute), time.Duration(12*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := D.Late(), interval.New(time.Duration(7*time.Minute), time.Duration(12*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := E.Early(), interval.New(time.Duration(5*time.Minute), time.Duration(6*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := E.Late(), interval.New(time.Duration(11*time.Minute), time.Duration(12*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := F.Early(), interval.New(time.Duration(5*time.Minute), time.Duration(7*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := F.Late(), interval.New(time.Duration(14*time.Minute), time.Duration(16*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := G.Early(), interval.New(time.Duration(12*time.Minute), time.Duration(16*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := G.Late(), interval.New(time.Duration(12*time.Minute), time.Duration(16*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := H.Early(), interval.New(time.Duration(16*time.Minute), time.Duration(19*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := H.Late(), interval.New(time.Duration(16*time.Minute), time.Duration(19*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
}

func TestLoneActivity(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*3))
	B := activity.New(NewMockActivityData("B", time.Minute*4))
	C := activity.New(NewMockActivityData("C", time.Minute*2))

	p.AddActivity(A)
	p.AddActivity(B)
	p.AddActivity(C)

	if got, want := len(slices.Collect(p.LoneActivities())), 3; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	foundA := false
	foundB := false
	foundC := false
	for a := range p.LoneActivities() {
		if a == A {
			foundA = true
		}
		if a == B {
			foundB = true
		}
		if a == C {
			foundC = true
		}
	}

	if !(foundA && foundB && foundC) {
		t.Errorf("missing activities")
	}
}

func TestUpdateActivityTimestampsLone(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*3))
	B := activity.New(NewMockActivityData("B", time.Minute*4))
	C := activity.New(NewMockActivityData("C", time.Minute*2))
	D := activity.New(NewMockActivityData("D", time.Minute*5))
	E := activity.New(NewMockActivityData("E", time.Minute*1))
	F := activity.New(NewMockActivityData("F", time.Minute*2))
	G := activity.New(NewMockActivityData("G", time.Minute*4))
	H := activity.New(NewMockActivityData("H", time.Minute*3))
	I := activity.New(NewMockActivityData("I", time.Minute*3))
	p.AddActivity(I)

	p.AddDependency(A, B, dependency.New(enums.FS))
	p.AddDependency(A, C, dependency.New(enums.FS))
	p.AddDependency(B, D, dependency.New(enums.FS))
	p.AddDependency(C, E, dependency.New(enums.FS))
	p.AddDependency(C, F, dependency.New(enums.FS))
	p.AddDependency(D, G, dependency.New(enums.FS))
	p.AddDependency(E, G, dependency.New(enums.FS))
	p.AddDependency(F, H, dependency.New(enums.FS))
	p.AddDependency(G, H, dependency.New(enums.FS))

	err := p.UpdateActivityTimestamps()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got, want := I.Early(), interval.New(time.Duration(0), time.Duration(time.Minute*3)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := I.Late(), interval.New(time.Duration(0), time.Duration(time.Minute*3)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
}

func TestTopologicalSorterError(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	err := p.updateActivityTimestamp(func(g graph.Graph[activity.Activity[*mockActivityData], dependency.Dependency]) ([]activity.Activity[*mockActivityData], error) {
		return nil, fmt.Errorf("test error")
	})

	if err == nil {
		t.Errorf("expected error")
	}
}

func TestStartToStartRelationship(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*2))
	B := activity.New(NewMockActivityData("B", time.Minute*3))
	C := activity.New(NewMockActivityData("C", time.Minute*4))
	D := activity.New(NewMockActivityData("D", time.Minute*3))

	p.AddDependency(A, B, dependency.New(enums.SS))
	p.AddDependency(B, C, dependency.New(enums.FS))
	p.AddDependency(C, D, dependency.NewWithLag(enums.SS, time.Minute*2))

	err := p.UpdateActivityTimestamps()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got, want := A.Early(), interval.New(time.Duration(0), time.Duration(time.Minute*2)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := A.Late(), interval.New(time.Duration(0), time.Duration(time.Minute*2)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := B.Early(), interval.New(time.Duration(0*time.Minute), time.Duration(time.Minute*3)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := B.Late(), interval.New(time.Duration(0*time.Minute), time.Duration(time.Minute*3)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := C.Early(), interval.New(time.Duration(3*time.Minute), time.Duration(7*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := C.Late(), interval.New(time.Duration(3*time.Minute), time.Duration(7*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := D.Early(), interval.New(time.Duration(5*time.Minute), time.Duration(8*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := D.Late(), interval.New(time.Duration(5*time.Minute), time.Duration(8*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

}

func TestFinishToFinishRelationship(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*3))
	B := activity.New(NewMockActivityData("B", time.Minute*2))
	C := activity.New(NewMockActivityData("C", time.Minute*3))
	D := activity.New(NewMockActivityData("D", time.Minute*3))

	p.AddDependency(A, B, dependency.New(enums.FF))
	p.AddDependency(B, C, dependency.NewWithLag(enums.FF, time.Minute*2))
	p.AddDependency(C, D, dependency.New(enums.FS))

	err := p.UpdateActivityTimestamps()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got, want := A.Early(), interval.New(time.Duration(0), time.Duration(3*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := A.Late(), interval.New(time.Duration(0), time.Duration(3*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := B.Early(), interval.New(time.Duration(1*time.Minute), time.Duration(3*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := B.Late(), interval.New(time.Duration(1*time.Minute), time.Duration(3*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := C.Early(), interval.New(time.Duration(2*time.Minute), time.Duration(5*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := C.Late(), interval.New(time.Duration(2*time.Minute), time.Duration(5*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := D.Early(), interval.New(time.Duration(5*time.Minute), time.Duration(8*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := D.Late(), interval.New(time.Duration(5*time.Minute), time.Duration(8*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

}

func TestStartToFinishRelationship(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*3))
	B := activity.New(NewMockActivityData("B", time.Minute*3))
	C := activity.New(NewMockActivityData("C", time.Minute*4))

	p.AddDependency(A, B, dependency.NewWithLag(enums.SF, time.Minute*5))
	p.AddDependency(B, C, dependency.New(enums.FS))

	err := p.UpdateActivityTimestamps()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got, want := A.Early(), interval.New(time.Duration(0), time.Duration(3*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := A.Late(), interval.New(time.Duration(0), time.Duration(3*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := B.Early(), interval.New(time.Duration(2*time.Minute), time.Duration(5*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := B.Late(), interval.New(time.Duration(2*time.Minute), time.Duration(5*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := C.Early(), interval.New(time.Duration(5*time.Minute), time.Duration(9*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := C.Late(), interval.New(time.Duration(5*time.Minute), time.Duration(9*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
}

func TestFinishToStartRelationship(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*3))
	B := activity.New(NewMockActivityData("B", time.Minute*4))
	C := activity.New(NewMockActivityData("C", time.Minute*3))

	p.AddDependency(A, B, dependency.New(enums.FS))
	p.AddDependency(B, C, dependency.NewWithLag(enums.FS, time.Minute*2))

	err := p.UpdateActivityTimestamps()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got, want := A.Early(), interval.New(time.Duration(0), time.Duration(3*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := A.Late(), interval.New(time.Duration(0), time.Duration(3*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := B.Early(), interval.New(time.Duration(3*time.Minute), time.Duration(7*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := B.Late(), interval.New(time.Duration(3*time.Minute), time.Duration(7*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := C.Early(), interval.New(time.Duration(9*time.Minute), time.Duration(12*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := C.Late(), interval.New(time.Duration(9*time.Minute), time.Duration(12*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
}

func TestMixedRelationships(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	A := activity.New(NewMockActivityData("A", time.Minute*3))
	B := activity.New(NewMockActivityData("B", time.Minute*2))
	C := activity.New(NewMockActivityData("C", time.Minute*2))
	D := activity.New(NewMockActivityData("D", time.Minute*4))
	E := activity.New(NewMockActivityData("E", time.Minute*1))
	F := activity.New(NewMockActivityData("F", time.Minute*2))
	G := activity.New(NewMockActivityData("G", time.Minute*4))
	H := activity.New(NewMockActivityData("H", time.Minute*3))

	p.AddDependency(A, B, dependency.NewWithLag(enums.FS, time.Minute*2))
	p.AddDependency(A, C, dependency.New(enums.SS))
	p.AddDependency(B, D, dependency.NewWithLag(enums.SS, time.Minute*1))
	p.AddDependency(C, E, dependency.NewWithLag(enums.SF, time.Minute*3))
	p.AddDependency(C, F, dependency.NewWithLag(enums.FF, time.Minute*3))
	p.AddDependency(D, G, dependency.NewWithLag(enums.SS, time.Minute*1))
	p.AddDependency(E, G, dependency.New(enums.FS))
	p.AddDependency(F, H, dependency.NewWithLag(enums.SF, time.Minute*2))
	p.AddDependency(G, H, dependency.New(enums.FS))

	err := p.UpdateActivityTimestamps()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got, want := A.Early(), interval.New(time.Duration(0), time.Duration(time.Minute*3)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := A.Late(), interval.New(time.Duration(0), time.Duration(time.Minute*3)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := B.Early(), interval.New(time.Duration(5*time.Minute), time.Duration(time.Minute*7)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := B.Late(), interval.New(time.Duration(5*time.Minute), time.Duration(time.Minute*7)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := C.Early(), interval.New(time.Duration(0*time.Minute), time.Duration(2*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := C.Late(), interval.New(time.Duration(4*time.Minute), time.Duration(6*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := D.Early(), interval.New(time.Duration(6*time.Minute), time.Duration(10*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := D.Late(), interval.New(time.Duration(6*time.Minute), time.Duration(10*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := E.Early(), interval.New(time.Duration(2*time.Minute), time.Duration(3*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := E.Late(), interval.New(time.Duration(6*time.Minute), time.Duration(7*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := F.Early(), interval.New(time.Duration(3*time.Minute), time.Duration(5*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := F.Late(), interval.New(time.Duration(12*time.Minute), time.Duration(14*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := G.Early(), interval.New(time.Duration(7*time.Minute), time.Duration(11*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := G.Late(), interval.New(time.Duration(7*time.Minute), time.Duration(11*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}

	if got, want := H.Early(), interval.New(time.Duration(11*time.Minute), time.Duration(14*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
	if got, want := H.Late(), interval.New(time.Duration(11*time.Minute), time.Duration(14*time.Minute)); !reflect.DeepEqual(got, want) {
		t.Errorf("got:  %T %#v", got, got)
		t.Errorf("want: %T %#v", want, want)
	}
}

func TestGetDependency(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("failed to GetDependency did not panic")
		}
	}()
	mockGraph := &mockGraphGetEdgePanic{Graph: graph.New[activity.Activity[*mockActivityData], dependency.Dependency]()}
	getDependency(mockGraph, nil, nil)
}
