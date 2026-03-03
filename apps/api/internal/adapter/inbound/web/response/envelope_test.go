package response

import (
	"encoding/json"
	"testing"
)

func TestSuccess_WithMessage(t *testing.T) {
	resp := Success("data", "Operation completed")
	if resp == nil {
		t.Fatal("Success returned nil")
	}
	if !resp.Success {
		t.Error("expected Success true")
	}
	if resp.Code != "OK" {
		t.Errorf("expected Code OK, got %s", resp.Code)
	}
	if resp.Message != "Operation completed" {
		t.Errorf("expected Message 'Operation completed', got %s", resp.Message)
	}
	if resp.Data == nil {
		t.Error("expected Data non-nil")
	}
}

func TestSuccess_EmptyMessage(t *testing.T) {
	resp := Success[any](nil, "")
	if resp == nil {
		t.Fatal("Success returned nil")
	}
	if resp.Message != "Operation completed successfully" {
		t.Errorf("expected default message, got %s", resp.Message)
	}
}

func TestSuccess_WithStruct(t *testing.T) {
	type Data struct {
		ID string `json:"id"`
	}
	resp := Success(Data{ID: "123"}, "Created")
	if resp == nil {
		t.Fatal("Success returned nil")
	}
	_, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
}

func TestFail_WithCodeAndMessage(t *testing.T) {
	resp := Fail[any](CodeBadRequest, "invalid input")
	if resp == nil {
		t.Fatal("Fail returned nil")
	}
	if resp.Success {
		t.Error("expected Success false")
	}
	if resp.Code != CodeBadRequest {
		t.Errorf("expected Code %s, got %s", CodeBadRequest, resp.Code)
	}
	if resp.Message != "invalid input" {
		t.Errorf("expected Message 'invalid input', got %s", resp.Message)
	}
}

func TestFail_EmptyMessage(t *testing.T) {
	resp := Fail[any]("ERR", "")
	if resp.Message != "An error occurred" {
		t.Errorf("expected default message, got %s", resp.Message)
	}
}

func TestFail_EmptyCode(t *testing.T) {
	resp := Fail[any]("", "something went wrong")
	if resp.Code != CodeInternalError {
		t.Errorf("expected Code %s when empty, got %s", CodeInternalError, resp.Code)
	}
}
