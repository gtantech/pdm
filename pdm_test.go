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
	"github.com/gtantech/pdm/relationship"
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
	graph.Graph[activity.Activity[*mockActivityData], relationship.Relationship]
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
func (g *mockGraph) AddEdge(value relationship.Relationship, origin activity.Activity[*mockActivityData], destination activity.Activity[*mockActivityData]) {
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
	graph.Graph[activity.Activity[*mockActivityData], relationship.Relationship]
}

func (g *mockGraphGetEdgePanic) GetEdgeValue(origin activity.Activity[*mockActivityData], destination activity.Activity[*mockActivityData]) (relationship.Relationship, bool) {
	return nil, false
}

func TestFields(t *testing.T) {
	p := New[*mockActivityData]()
	mockGraph := &mockGraph{Graph: graph.New[activity.Activity[*mockActivityData], relationship.Relationship]()}
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
	p.AddDependency(a, b, relationship.New(enums.FS))
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

	p.AddDependency(A, B, relationship.New(enums.FS))
	p.AddDependency(A, C, relationship.New(enums.FS))
	p.AddDependency(B, D, relationship.New(enums.FS))
	p.AddDependency(C, E, relationship.New(enums.FS))
	p.AddDependency(D, F, relationship.New(enums.FS))
	p.AddDependency(E, F, relationship.New(enums.FS))

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

	p.AddDependency(A, B, relationship.New(enums.FS))
	p.AddDependency(A, C, relationship.New(enums.FS))
	p.AddDependency(B, D, relationship.New(enums.FS))
	p.AddDependency(C, E, relationship.New(enums.FS))
	p.AddDependency(D, F, relationship.New(enums.FS))
	p.AddDependency(E, F, relationship.New(enums.FS))

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

	p.AddDependency(A, B, relationship.New(enums.FS))
	p.AddDependency(A, C, relationship.New(enums.FS))
	p.AddDependency(B, D, relationship.New(enums.FS))
	p.AddDependency(C, E, relationship.New(enums.FS))
	p.AddDependency(D, F, relationship.New(enums.FS))
	p.AddDependency(E, F, relationship.New(enums.FS))

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

	p.AddDependency(A, B, relationship.New(enums.FS))
	p.AddDependency(A, C, relationship.New(enums.FS))
	p.AddDependency(B, D, relationship.New(enums.FS))
	p.AddDependency(C, E, relationship.New(enums.FS))
	p.AddDependency(C, F, relationship.New(enums.FS))
	p.AddDependency(D, G, relationship.New(enums.FS))
	p.AddDependency(E, G, relationship.New(enums.FS))
	p.AddDependency(F, H, relationship.New(enums.FS))
	p.AddDependency(G, H, relationship.New(enums.FS))

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

	p.AddDependency(A, B, relationship.New(enums.FS))
	p.AddDependency(A, C, relationship.New(enums.FS))
	p.AddDependency(B, D, relationship.New(enums.FS))
	p.AddDependency(C, E, relationship.New(enums.FS))
	p.AddDependency(C, F, relationship.New(enums.FS))
	p.AddDependency(D, G, relationship.New(enums.FS))
	p.AddDependency(E, G, relationship.New(enums.FS))
	p.AddDependency(F, H, relationship.New(enums.FS))
	p.AddDependency(G, H, relationship.New(enums.FS))

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

	err := p.updateActivityTimestamp(func(g graph.Graph[activity.Activity[*mockActivityData], relationship.Relationship]) ([]activity.Activity[*mockActivityData], error) {
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

	p.AddDependency(A, B, relationship.New(enums.SS))
	p.AddDependency(B, C, relationship.New(enums.FS))
	p.AddDependency(C, D, relationship.NewWithLag(enums.SS, time.Minute*2))

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

	p.AddDependency(A, B, relationship.New(enums.FF))
	p.AddDependency(B, C, relationship.NewWithLag(enums.FF, time.Minute*2))
	p.AddDependency(C, D, relationship.New(enums.FS))

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

	p.AddDependency(A, B, relationship.NewWithLag(enums.SF, time.Minute*5))
	p.AddDependency(B, C, relationship.New(enums.FS))

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

	p.AddDependency(A, B, relationship.New(enums.FS))
	p.AddDependency(B, C, relationship.NewWithLag(enums.FS, time.Minute*2))

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

	p.AddDependency(A, B, relationship.NewWithLag(enums.FS, time.Minute*2))
	p.AddDependency(A, C, relationship.New(enums.SS))
	p.AddDependency(B, D, relationship.NewWithLag(enums.SS, time.Minute*1))
	p.AddDependency(C, E, relationship.NewWithLag(enums.SF, time.Minute*3))
	p.AddDependency(C, F, relationship.NewWithLag(enums.FF, time.Minute*3))
	p.AddDependency(D, G, relationship.NewWithLag(enums.SS, time.Minute*1))
	p.AddDependency(E, G, relationship.New(enums.FS))
	p.AddDependency(F, H, relationship.NewWithLag(enums.SF, time.Minute*2))
	p.AddDependency(G, H, relationship.New(enums.FS))

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
	mockGraph := &mockGraphGetEdgePanic{Graph: graph.New[activity.Activity[*mockActivityData], relationship.Relationship]()}
	getDependency(mockGraph, nil, nil)
}

func TestFreeFloat(t *testing.T) {
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

	p.AddDependency(A, B, relationship.New(enums.FS))
	p.AddDependency(A, C, relationship.New(enums.FS))
	p.AddDependency(B, D, relationship.New(enums.FS))
	p.AddDependency(C, E, relationship.New(enums.FS))
	p.AddDependency(C, F, relationship.New(enums.FS))
	p.AddDependency(D, G, relationship.New(enums.FS))
	p.AddDependency(E, G, relationship.New(enums.FS))
	p.AddDependency(F, H, relationship.New(enums.FS))
	p.AddDependency(G, H, relationship.New(enums.FS))

	p.UpdateActivityTimestamps()

	if got, want := p.FreeFloat(A), time.Duration(0*time.Minute); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := p.FreeFloat(B), time.Duration(0*time.Minute); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := p.FreeFloat(C), time.Duration(0*time.Minute); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := p.FreeFloat(D), time.Duration(0*time.Minute); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := p.FreeFloat(E), time.Duration(6*time.Minute); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := p.FreeFloat(F), time.Duration(9*time.Minute); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := p.FreeFloat(G), time.Duration(0*time.Minute); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := p.FreeFloat(H), time.Duration(0*time.Minute); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTotalFloat(t *testing.T) {
	p := New[*mockActivityData]()

	if p.graph == nil {
		t.Errorf("expected non-nil graph")
	}

	a := activity.New(NewMockActivityData("A", time.Minute*3))
	a.UpdateEarly(interval.New(5*time.Minute, 8*time.Minute))
	a.UpdateLate(interval.New(6*time.Minute, 9*time.Minute))

	if got, want := p.TotalFloat(a), 1*time.Minute; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCriticalActivities(t *testing.T) {
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

	p.AddDependency(A, B, relationship.New(enums.FS))
	p.AddDependency(A, C, relationship.New(enums.FS))
	p.AddDependency(B, D, relationship.New(enums.FS))
	p.AddDependency(C, E, relationship.New(enums.FS))
	p.AddDependency(C, F, relationship.New(enums.FS))
	p.AddDependency(D, G, relationship.New(enums.FS))
	p.AddDependency(E, G, relationship.New(enums.FS))
	p.AddDependency(F, H, relationship.New(enums.FS))
	p.AddDependency(G, H, relationship.New(enums.FS))

	p.UpdateActivityTimestamps()

	crticialActivities := slices.Collect(p.CriticalActivities(0))

	if got, want := len(crticialActivities), 5; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := slices.ContainsFunc(crticialActivities, func(a activity.Activity[*mockActivityData]) bool {
		return a == A
	}), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := slices.ContainsFunc(crticialActivities, func(a activity.Activity[*mockActivityData]) bool {
		return a == B
	}), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := slices.ContainsFunc(crticialActivities, func(a activity.Activity[*mockActivityData]) bool {
		return a == D
	}), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := slices.ContainsFunc(crticialActivities, func(a activity.Activity[*mockActivityData]) bool {
		return a == G
	}), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := slices.ContainsFunc(crticialActivities, func(a activity.Activity[*mockActivityData]) bool {
		return a == H
	}), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAddDependencies(t *testing.T) {
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

	d1 := dependency.New(A, B, relationship.New(enums.FS))
	d2 := dependency.New(A, C, relationship.New(enums.FS))
	d3 := dependency.New(B, D, relationship.New(enums.FS))
	d4 := dependency.New(C, E, relationship.New(enums.FS))
	d5 := dependency.New(C, F, relationship.New(enums.FS))
	d6 := dependency.New(D, G, relationship.New(enums.FS))
	d7 := dependency.New(E, G, relationship.New(enums.FS))
	d8 := dependency.New(F, H, relationship.New(enums.FS))
	d9 := dependency.New(G, H, relationship.New(enums.FS))

	err := p.AddDependencies([]dependency.Dependency[*mockActivityData]{d1, d2, d3, d4, d5, d6, d7, d8, d9})
	if err != nil {
		t.Errorf("unexpected error occurred: %v", err)
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
