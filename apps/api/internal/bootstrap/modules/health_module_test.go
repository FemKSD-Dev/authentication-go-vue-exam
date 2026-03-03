package modules

import (
	"testing"
)

func TestNewHealthModule(t *testing.T) {
	mod := NewHealthModule()
	if mod == nil {
		t.Fatal("expected non-nil HealthModule")
	}
	if mod.Handler == nil {
		t.Fatal("expected non-nil Handler")
	}
}
