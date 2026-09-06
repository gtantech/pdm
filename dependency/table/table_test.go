package table

import (
	"slices"
	"testing"

	"github.com/gtantech/pdm/activity"
)

func TestAddRow(t *testing.T) {
	tbl := New[activity.Data]()
	A := activity.New(activity.Data(activity.NewData(0)))
	tbl.AddRow(A)
	row, ok := tbl.table[A]
	if !ok {
		t.Errorf("expected get value to succeed, got false")
	}
	if got, want := len(row), 0; got != want {
		t.Errorf("expected empty slice, got %v", row)
	}
}

func TestGetRow(t *testing.T) {
	tbl := New[activity.Data]()
	A := activity.New(activity.Data(activity.NewData(0)))
	B := activity.New(activity.Data(activity.NewData(0)))
	C := activity.New(activity.Data(activity.NewData(0)))
	inputRow := []activity.Activity[activity.Data]{B, C}
	tbl.table[A] = inputRow
	row, ok := tbl.GetRow(A)
	if !ok {
		t.Errorf("expected GetRow to succeed, got false")
	}
	if got, want := row, inputRow; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", row, inputRow)
	}
}

func TestDeleteRow(t *testing.T) {
	tbl := New[activity.Data]()
	A := activity.New(activity.Data(activity.NewData(0)))
	B := activity.New(activity.Data(activity.NewData(0)))
	C := activity.New(activity.Data(activity.NewData(0)))
	inputRow := []activity.Activity[activity.Data]{B, C}
	tbl.table[A] = inputRow
	_, ok := tbl.GetRow(A)
	if !ok {
		t.Errorf("expected GetRow to succeed, got false")
	}
	tbl.DeleteRow(A)
	_, ok = tbl.GetRow(A)
	if ok {
		t.Errorf("expected GetRow to not succeed, got true")
	}
}

func TestUpdatePredecessors(t *testing.T) {
	tbl := New[activity.Data]()
	A := activity.New(activity.Data(activity.NewData(0)))
	B := activity.New(activity.Data(activity.NewData(0)))
	C := activity.New(activity.Data(activity.NewData(0)))
	D := activity.New(activity.Data(activity.NewData(0)))
	inputRow := []activity.Activity[activity.Data]{B, C}
	tbl.table[A] = inputRow
	row, ok := tbl.GetRow(A)
	if !ok {
		t.Errorf("expected GetRow to succeed, got false")
	}
	if got, want := row, inputRow; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", row, inputRow)
	}

	newInputRow := []activity.Activity[activity.Data]{B, C, D}

	tbl.UpdatePredecessors(A, newInputRow)

	row, ok = tbl.GetRow(A)
	if !ok {
		t.Errorf("expected GetRow to succeed, got false")
	}
	if got, want := row, newInputRow; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", row, inputRow)
	}
}

func TestGetActivities(t *testing.T) {
	tbl := New[activity.Data]()
	A := activity.New(activity.Data(activity.NewData(0)))
	B := activity.New(activity.Data(activity.NewData(0)))
	C := activity.New(activity.Data(activity.NewData(0)))
	inputRow := []activity.Activity[activity.Data]{B, C}
	activities := []activity.Activity[activity.Data]{A, B, C}
	tbl.UpdatePredecessors(A, inputRow)
	if got, want := len(slices.Collect(tbl.GetActivities())), 3; got != want {
		t.Errorf("expected slice of length %v, got %v", want, got)
	}
	for a := range tbl.GetActivities() {
		if !slices.Contains(activities, a) {
			t.Errorf("missing activity %v", a)
		}
	}
}
