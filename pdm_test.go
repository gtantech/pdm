package pdm

import (
	"testing"
)

func TestPDMDisplayName(t *testing.T) {
	dispName := "test_name"
	p := New(dispName)

	if got := p.DisplayName(); got != dispName {
		t.Errorf("got %v, want %v", got, dispName)
	}
}
