package request

import (
	"encoding/json"
	"testing"
)

func TestRegisterRequest_JSON(t *testing.T) {
	body := `{"username":"user1","password":"pass1234","confirm_password":"pass1234"}`
	var req RegisterRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if req.Username != "user1" {
		t.Errorf("expected username user1, got %s", req.Username)
	}
	if req.Password != "pass1234" {
		t.Errorf("expected password, got %s", req.Password)
	}
}

func TestLoginRequest_JSON(t *testing.T) {
	body := `{"username":"user1","password":"pass1234"}`
	var req LoginRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if req.Username != "user1" {
		t.Errorf("expected username user1, got %s", req.Username)
	}
}
