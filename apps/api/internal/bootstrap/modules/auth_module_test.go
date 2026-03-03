package modules

import (
	"testing"

	"authentication-project-exam/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewAuthModule(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	cfg := &config.Config{
		JWTSecret:                 "test-secret",
		JWTExpireMinutes:          30,
		RefreshTokenExpireMinutes:  60,
	}

	mod, err := NewAuthModule(db, cfg)
	if err != nil {
		t.Fatalf("NewAuthModule failed: %v", err)
	}
	if mod == nil {
		t.Fatal("expected non-nil AuthModule")
	}
	if mod.Handler == nil {
		t.Fatal("expected non-nil Handler")
	}
	if mod.AuthMiddleware == nil {
		t.Fatal("expected non-nil AuthMiddleware")
	}
}
