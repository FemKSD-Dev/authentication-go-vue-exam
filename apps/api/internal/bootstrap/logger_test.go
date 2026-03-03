package bootstrap

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}
