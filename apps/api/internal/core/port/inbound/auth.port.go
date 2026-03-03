package inbound

import (
	"context"
)

type RegisterPayload struct {
	Username        string
	Password        string
	ConfirmPassword string
}

type RegisterResult struct {
	UserID   string
	Username string
}

type LoginPayload struct {
	Username string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

type MeQuery struct {
	UserID string
}

type MeResult struct {
	UserID   string
	Username string
}

type RefreshTokenPayload struct {
	RefreshToken string
}

type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
}

type AuthPort interface {
	Register(ctx context.Context, payload RegisterPayload) (*RegisterResult, error)
	Login(ctx context.Context, payload LoginPayload) (*LoginResult, error)
	Me(ctx context.Context, query MeQuery) (*MeResult, error)
	RefreshToken(ctx context.Context, payload RefreshTokenPayload) (*RefreshTokenResult, error)
}
