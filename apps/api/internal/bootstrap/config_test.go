package bootstrap

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.HTTPPort == "" {
		t.Error("expected HTTPPort to be set")
	}
}
