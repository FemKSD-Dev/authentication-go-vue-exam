package outbound

import (
	"context"
	"testing"

	"authentication-project-exam/internal/core/model"
)

func TestTokenClaims(t *testing.T) {
	c := TokenClaims{UserID: "u1", Username: "john"}
	if c.UserID != "u1" {
		t.Error("TokenClaims.UserID mismatch")
	}
}

func TestTokenManager_Interface(t *testing.T) {
	var _ TokenManager = (*mockTokenManager)(nil)
}

type mockTokenManager struct{}

func (mockTokenManager) Issue(userID, username string) (string, string, error) {
	return "", "", nil
}
func (mockTokenManager) Verify(token string) (*TokenClaims, error) {
	return nil, nil
}

func TestPasswordManager_Interface(t *testing.T) {
	var _ PasswordManager = (*mockPasswordManager)(nil)
}

type mockPasswordManager struct{}

func (mockPasswordManager) Encode(raw string) (string, error) {
	return "", nil
}
func (mockPasswordManager) Matches(encoded, raw string) error {
	return nil
}

func TestUsersRepository_Interface(t *testing.T) {
	var _ UsersRepository = (*mockUsersRepository)(nil)
}

type mockUsersRepository struct{}

func (mockUsersRepository) Save(ctx context.Context, user *model.User) error {
	return nil
}
func (mockUsersRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return nil, nil
}
func (mockUsersRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	return nil, nil
}
