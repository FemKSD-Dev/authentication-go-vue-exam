package service

import (
	errorx "authentication-project-exam/internal/core/error"
	"authentication-project-exam/internal/core/model"
	"authentication-project-exam/internal/core/port/inbound"
	"authentication-project-exam/internal/core/port/outbound"
	"context"
	"errors"
	"testing"
)

type fakeUsersRepository struct {
	saveFn           func(ctx context.Context, user *model.User) error
	findByUsernameFn func(ctx context.Context, username string) (*model.User, error)
	findByIDFn       func(ctx context.Context, id string) (*model.User, error)
}

func (fr *fakeUsersRepository) Save(ctx context.Context, user *model.User) error {
	if fr.saveFn != nil {
		return fr.saveFn(ctx, user)
	}
	return nil
}

func (fr *fakeUsersRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	if fr.findByUsernameFn != nil {
		return fr.findByUsernameFn(ctx, username)
	}
	return nil, nil
}

func (fr *fakeUsersRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	if fr.findByIDFn != nil {
		return fr.findByIDFn(ctx, id)
	}
	return nil, nil
}

type fakePasswordManager struct {
	encodeFn  func(raw string) (string, error)
	matchesFn func(encoded, raw string) error
}

func (fp *fakePasswordManager) Encode(raw string) (string, error) {
	if fp.encodeFn != nil {
		return fp.encodeFn(raw)
	}
	return "hashed-password", nil
}

func (fp *fakePasswordManager) Matches(encoded, raw string) error {
	if fp.matchesFn != nil {
		return fp.matchesFn(encoded, raw)
	}
	return nil
}

type fakeTokenManager struct {
	issueFn  func(userID, username string) (string, string, error)
	verifyFn func(token string) (*outbound.TokenClaims, error)
}

func (ft *fakeTokenManager) Issue(userID, username string) (string, string, error) {
	if ft.issueFn != nil {
		return ft.issueFn(userID, username)
	}
	return "fake-access-token", "fake-refresh-token", nil
}

func (ft *fakeTokenManager) Verify(token string) (*outbound.TokenClaims, error) {
	if ft.verifyFn != nil {
		return ft.verifyFn(token)
	}
	return &outbound.TokenClaims{
		UserID:   "fake-user-id",
		Username: "fake-username",
	}, nil
}

func TestAuthService_Register_Success(t *testing.T) {
	var savedUser *model.User

	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			if savedUser != nil && savedUser.Username == username {
				return savedUser, nil
			}
			return nil, nil
		},
		saveFn: func(ctx context.Context, user *model.User) error {
			savedUser = user
			user.ID = "generated-user-id"
			return nil
		},
	}

	passwordManager := &fakePasswordManager{
		encodeFn: func(raw string) (string, error) {
			return "hashed-secret123", nil
		},
	}

	token := &fakeTokenManager{}
	svc := NewAuthService(users, passwordManager, token)

	result, err := svc.Register(context.Background(), inbound.RegisterPayload{
		Username:        "john",
		Password:        "secret123",
		ConfirmPassword: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.UserID != "generated-user-id" {
		t.Fatalf("expected user id generated-user-id, got %s", result.UserID)
	}

	if result.Username != "john" {
		t.Fatalf("expected username john, got %s", result.Username)
	}

	if savedUser == nil {
		t.Fatal("expected user to be saved")
	}

	if savedUser.Username != "john" {
		t.Fatalf("expected saved username john, got %s", savedUser.Username)
	}

	if savedUser.Password != "hashed-secret123" {
		t.Fatalf("expected hashed password, got %s", savedUser.Password)
	}
}

func TestAuthService_Register_PasswordMismatch(t *testing.T) {
	svc := NewAuthService(
		&fakeUsersRepository{},
		&fakePasswordManager{},
		&fakeTokenManager{},
	)

	_, err := svc.Register(context.Background(), inbound.RegisterPayload{
		Username:        "john",
		Password:        "secret123",
		ConfirmPassword: "secret456",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errorx.ErrorPasswordMismatch) {
		t.Fatalf("expected password mismatch error, got %v", err)
	}
}

func TestAuthService_Register_UsernameAlreadyExists(t *testing.T) {
	svc := NewAuthService(
		&fakeUsersRepository{
			findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
				return &model.User{ID: "user-1", Username: username}, nil
			},
		},
		&fakePasswordManager{},
		&fakeTokenManager{},
	)

	_, err := svc.Register(context.Background(), inbound.RegisterPayload{
		Username:        "john",
		Password:        "secret123",
		ConfirmPassword: "secret123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errorx.ErrorUserAlreadyExists) {
		t.Fatalf("expected user already exists error, got %v", err)
	}
}

func TestAuthService_Register_FindByUsername_OtherError(t *testing.T) {
	dbErr := errors.New("database connection failed")
	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return nil, dbErr
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, &fakeTokenManager{})

	_, err := svc.Register(context.Background(), inbound.RegisterPayload{
		Username:        "john",
		Password:        "secret123",
		ConfirmPassword: "secret123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestAuthService_Register_FindByUsername_ErrorUserNotFound_Continues(t *testing.T) {
	var savedUser *model.User
	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			if savedUser != nil && savedUser.Username == username {
				return savedUser, nil
			}
			return nil, errorx.ErrorUserNotFound
		},
		saveFn: func(ctx context.Context, user *model.User) error {
			savedUser = user
			user.ID = "new-id"
			return nil
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, &fakeTokenManager{})

	result, err := svc.Register(context.Background(), inbound.RegisterPayload{
		Username:        "john",
		Password:        "secret123",
		ConfirmPassword: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Username != "john" {
		t.Fatalf("expected username john, got %s", result.Username)
	}
}

func TestAuthService_Register_EncodeError(t *testing.T) {
	encodeErr := errors.New("encoding failed")
	passwordManager := &fakePasswordManager{
		encodeFn: func(raw string) (string, error) {
			return "", encodeErr
		},
	}

	svc := NewAuthService(&fakeUsersRepository{}, passwordManager, &fakeTokenManager{})

	_, err := svc.Register(context.Background(), inbound.RegisterPayload{
		Username:        "john",
		Password:        "secret123",
		ConfirmPassword: "secret123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, encodeErr) {
		t.Fatalf("expected encode error, got %v", err)
	}
}

func TestAuthService_Register_SaveError(t *testing.T) {
	saveErr := errors.New("save failed")
	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return nil, nil
		},
		saveFn: func(ctx context.Context, user *model.User) error {
			return saveErr
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, &fakeTokenManager{})

	_, err := svc.Register(context.Background(), inbound.RegisterPayload{
		Username:        "john",
		Password:        "secret123",
		ConfirmPassword: "secret123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, saveErr) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestAuthService_Register_FindByUsernameAfterSave_Error(t *testing.T) {
	var savedUser *model.User
	findErr := errors.New("find failed")
	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			if savedUser != nil && savedUser.Username == username {
				return nil, findErr
			}
			return nil, nil
		},
		saveFn: func(ctx context.Context, user *model.User) error {
			savedUser = user
			return nil
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, &fakeTokenManager{})

	_, err := svc.Register(context.Background(), inbound.RegisterPayload{
		Username:        "john",
		Password:        "secret123",
		ConfirmPassword: "secret123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, findErr) {
		t.Fatalf("expected find error, got %v", err)
	}
}

// --- Login tests ---

func TestAuthService_Login_Success(t *testing.T) {
	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				ID:       "user-1",
				Username: "john",
				Password: "hashed-secret123",
			}, nil
		},
	}

	passwordManager := &fakePasswordManager{
		matchesFn: func(encoded, raw string) error {
			if encoded == "hashed-secret123" && raw == "secret123" {
				return nil
			}
			return errors.New("mismatch")
		},
	}

	token := &fakeTokenManager{}
	svc := NewAuthService(users, passwordManager, token)

	result, err := svc.Login(context.Background(), inbound.LoginPayload{
		Username: "john",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.AccessToken != "fake-access-token" {
		t.Fatalf("expected access token, got %s", result.AccessToken)
	}
	if result.RefreshToken != "fake-refresh-token" {
		t.Fatalf("expected refresh token, got %s", result.RefreshToken)
	}
}

func TestAuthService_Login_FindByUsername_Error(t *testing.T) {
	findErr := errors.New("find failed")
	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return nil, findErr
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, &fakeTokenManager{})

	_, err := svc.Login(context.Background(), inbound.LoginPayload{
		Username: "john",
		Password: "secret123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, findErr) {
		t.Fatalf("expected find error, got %v", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return nil, nil
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, &fakeTokenManager{})

	_, err := svc.Login(context.Background(), inbound.LoginPayload{
		Username: "john",
		Password: "secret123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errorx.ErrorUserNotFound) {
		t.Fatalf("expected user not found error, got %v", err)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				ID:       "user-1",
				Username: "john",
				Password: "hashed-secret123",
			}, nil
		},
	}

	passwordManager := &fakePasswordManager{
		matchesFn: func(encoded, raw string) error {
			return errors.New("password mismatch")
		},
	}

	svc := NewAuthService(users, passwordManager, &fakeTokenManager{})

	_, err := svc.Login(context.Background(), inbound.LoginPayload{
		Username: "john",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errorx.ErrorInvalidCredentials) {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestAuthService_Login_IssueError(t *testing.T) {
	users := &fakeUsersRepository{
		findByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				ID:       "user-1",
				Username: "john",
				Password: "hashed-secret123",
			}, nil
		},
	}

	issueErr := errors.New("issue failed")
	token := &fakeTokenManager{
		issueFn: func(userID, username string) (string, string, error) {
			return "", "", issueErr
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, token)

	_, err := svc.Login(context.Background(), inbound.LoginPayload{
		Username: "john",
		Password: "secret123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, issueErr) {
		t.Fatalf("expected issue error, got %v", err)
	}
}

// --- Me tests ---

func TestAuthService_Me_Success(t *testing.T) {
	users := &fakeUsersRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return &model.User{
				ID:       "user-1",
				Username: "john",
			}, nil
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, &fakeTokenManager{})

	result, err := svc.Me(context.Background(), inbound.MeQuery{UserID: "user-1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.UserID != "user-1" {
		t.Fatalf("expected user id user-1, got %s", result.UserID)
	}
	if result.Username != "john" {
		t.Fatalf("expected username john, got %s", result.Username)
	}
}

func TestAuthService_Me_FindByID_Error(t *testing.T) {
	findErr := errors.New("find failed")
	users := &fakeUsersRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, findErr
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, &fakeTokenManager{})

	_, err := svc.Me(context.Background(), inbound.MeQuery{UserID: "user-1"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, findErr) {
		t.Fatalf("expected find error, got %v", err)
	}
}

func TestAuthService_Me_UserNotFound(t *testing.T) {
	users := &fakeUsersRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, nil
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, &fakeTokenManager{})

	_, err := svc.Me(context.Background(), inbound.MeQuery{UserID: "user-1"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errorx.ErrorUserNotFound) {
		t.Fatalf("expected user not found error, got %v", err)
	}
}

// --- RefreshToken tests ---

func TestAuthService_RefreshToken_Success(t *testing.T) {
	users := &fakeUsersRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return &model.User{
				ID:       "user-1",
				Username: "john",
			}, nil
		},
	}

	token := &fakeTokenManager{
		verifyFn: func(token string) (*outbound.TokenClaims, error) {
			return &outbound.TokenClaims{UserID: "user-1", Username: "john"}, nil
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, token)

	result, err := svc.RefreshToken(context.Background(), inbound.RefreshTokenPayload{
		RefreshToken: "valid-refresh-token",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.AccessToken != "fake-access-token" {
		t.Fatalf("expected access token, got %s", result.AccessToken)
	}
	if result.RefreshToken != "fake-refresh-token" {
		t.Fatalf("expected refresh token, got %s", result.RefreshToken)
	}
}

func TestAuthService_RefreshToken_VerifyError(t *testing.T) {
	verifyErr := errors.New("invalid token")
	token := &fakeTokenManager{
		verifyFn: func(token string) (*outbound.TokenClaims, error) {
			return nil, verifyErr
		},
	}

	svc := NewAuthService(&fakeUsersRepository{}, &fakePasswordManager{}, token)

	_, err := svc.RefreshToken(context.Background(), inbound.RefreshTokenPayload{
		RefreshToken: "invalid-token",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, verifyErr) {
		t.Fatalf("expected verify error, got %v", err)
	}
}

func TestAuthService_RefreshToken_FindByID_Error(t *testing.T) {
	findErr := errors.New("find failed")
	users := &fakeUsersRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, findErr
		},
	}

	token := &fakeTokenManager{
		verifyFn: func(token string) (*outbound.TokenClaims, error) {
			return &outbound.TokenClaims{UserID: "user-1", Username: "john"}, nil
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, token)

	_, err := svc.RefreshToken(context.Background(), inbound.RefreshTokenPayload{
		RefreshToken: "valid-token",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, findErr) {
		t.Fatalf("expected find error, got %v", err)
	}
}

func TestAuthService_RefreshToken_UserNotFound(t *testing.T) {
	users := &fakeUsersRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, nil
		},
	}

	token := &fakeTokenManager{
		verifyFn: func(token string) (*outbound.TokenClaims, error) {
			return &outbound.TokenClaims{UserID: "user-1", Username: "john"}, nil
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, token)

	_, err := svc.RefreshToken(context.Background(), inbound.RefreshTokenPayload{
		RefreshToken: "valid-token",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errorx.ErrorUserNotFound) {
		t.Fatalf("expected user not found error, got %v", err)
	}
}

func TestAuthService_RefreshToken_IssueError(t *testing.T) {
	users := &fakeUsersRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return &model.User{ID: "user-1", Username: "john"}, nil
		},
	}

	issueErr := errors.New("issue failed")
	token := &fakeTokenManager{
		verifyFn: func(token string) (*outbound.TokenClaims, error) {
			return &outbound.TokenClaims{UserID: "user-1", Username: "john"}, nil
		},
		issueFn: func(userID, username string) (string, string, error) {
			return "", "", issueErr
		},
	}

	svc := NewAuthService(users, &fakePasswordManager{}, token)

	_, err := svc.RefreshToken(context.Background(), inbound.RefreshTokenPayload{
		RefreshToken: "valid-token",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, issueErr) {
		t.Fatalf("expected issue error, got %v", err)
	}
}

func TestNewAuthService(t *testing.T) {
	svc := NewAuthService(&fakeUsersRepository{}, &fakePasswordManager{}, &fakeTokenManager{})
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}
