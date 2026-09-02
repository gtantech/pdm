package dependency

import (
	"testing"

	"github.com/gtantech/pdm/enums"
)

func TestType(t *testing.T) {
	kind := enums.FS
	d := New(kind)
	if got := d.Type(); got != kind {
		t.Errorf("got %v, want %v", got, kind)
	}
}
