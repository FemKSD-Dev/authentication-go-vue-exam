package inbound

import (
	"context"
	"testing"
)

func TestRegisterPayload(t *testing.T) {
	p := RegisterPayload{
		Username:        "u",
		Password:        "p",
		ConfirmPassword: "p",
	}
	if p.Username != "u" {
		t.Error("RegisterPayload.Username mismatch")
	}
}

func TestLoginPayload(t *testing.T) {
	p := LoginPayload{Username: "u", Password: "p"}
	if p.Username != "u" {
		t.Error("LoginPayload.Username mismatch")
	}
}

func TestMeQuery(t *testing.T) {
	q := MeQuery{UserID: "id-1"}
	if q.UserID != "id-1" {
		t.Error("MeQuery.UserID mismatch")
	}
}

func TestAuthPort_Interface(t *testing.T) {
	var _ AuthPort = (*mockAuthPort)(nil)
}

type mockAuthPort struct{}

func (mockAuthPort) Register(ctx context.Context, p RegisterPayload) (*RegisterResult, error) {
	return nil, nil
}
func (mockAuthPort) Login(ctx context.Context, p LoginPayload) (*LoginResult, error) {
	return nil, nil
}
func (mockAuthPort) Me(ctx context.Context, q MeQuery) (*MeResult, error) {
	return nil, nil
}
func (mockAuthPort) RefreshToken(ctx context.Context, p RefreshTokenPayload) (*RefreshTokenResult, error) {
	return nil, nil
}
