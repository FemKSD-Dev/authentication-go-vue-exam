package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultValues(t *testing.T) {
	keys := []string{"HTTP_PORT", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD",
		"POSTGRES_DB", "JWT_SECRET", "JWT_EXPIRE_MINUTES", "REFRESH_TOKEN_EXPIRE_MINUTES"}
	restore := make(map[string]string)
	for _, k := range keys {
		if v, ok := os.LookupEnv(k); ok {
			restore[k] = v
			os.Unsetenv(k)
		}
	}
	defer func() {
		for k, v := range restore {
			os.Setenv(k, v)
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load returned nil config")
	}
	if cfg.HTTPPort != "8080" {
		t.Errorf("expected HTTP_PORT 8080, got %s", cfg.HTTPPort)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected DBHost localhost, got %s", cfg.DBHost)
	}
	if cfg.JWTExpireMinutes != 30 {
		t.Errorf("expected JWTExpireMinutes 30, got %d", cfg.JWTExpireMinutes)
	}
	if cfg.RefreshTokenExpireMinutes != 60 {
		t.Errorf("expected RefreshTokenExpireMinutes 60, got %d", cfg.RefreshTokenExpireMinutes)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("HTTP_PORT", "3000")
	os.Setenv("POSTGRES_HOST", "db.example.com")
	os.Setenv("JWT_EXPIRE_MINUTES", "45")
	os.Setenv("REFRESH_TOKEN_EXPIRE_MINUTES", "120")
	defer func() {
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("JWT_EXPIRE_MINUTES")
		os.Unsetenv("REFRESH_TOKEN_EXPIRE_MINUTES")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.HTTPPort != "3000" {
		t.Errorf("expected HTTP_PORT 3000, got %s", cfg.HTTPPort)
	}
	if cfg.DBHost != "db.example.com" {
		t.Errorf("expected DBHost db.example.com, got %s", cfg.DBHost)
	}
	if cfg.JWTExpireMinutes != 45 {
		t.Errorf("expected JWTExpireMinutes 45, got %d", cfg.JWTExpireMinutes)
	}
	if cfg.RefreshTokenExpireMinutes != 120 {
		t.Errorf("expected RefreshTokenExpireMinutes 120, got %d", cfg.RefreshTokenExpireMinutes)
	}
}

func TestLoad_InvalidIntEnv(t *testing.T) {
	os.Setenv("JWT_EXPIRE_MINUTES", "not-a-number")
	defer os.Unsetenv("JWT_EXPIRE_MINUTES")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.JWTExpireMinutes != 30 {
		t.Errorf("expected fallback 30, got %d", cfg.JWTExpireMinutes)
	}
}
