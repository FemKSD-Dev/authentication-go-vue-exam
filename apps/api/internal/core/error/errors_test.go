package error

import (
	"errors"
	"testing"
)

func TestErrorVariables(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"InvalidCredentials", ErrorInvalidCredentials, "invalid credentials"},
		{"UserAlreadyExists", ErrorUserAlreadyExists, "username already exists"},
		{"PasswordMismatch", ErrorPasswordMismatch, "password and confirm password do not match"},
		{"UserNotFound", ErrorUserNotFound, "user not found"},
		{"InternalServerError", ErrorInternalServerError, "internal server error"},
		{"Unauthorized", ErrorUnauthorized, "unauthorized"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatal("expected non-nil error")
			}
			if tt.err.Error() != tt.msg {
				t.Errorf("expected message %q, got %q", tt.msg, tt.err.Error())
			}
			if !errors.Is(tt.err, tt.err) {
				t.Error("errors.Is should match self")
			}
		})
	}
}
