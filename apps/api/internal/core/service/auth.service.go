package service

import (
	errs "authentication-project-exam/internal/core/error"
	"authentication-project-exam/internal/core/model"
	"authentication-project-exam/internal/core/port/inbound"
	"authentication-project-exam/internal/core/port/outbound"
	"context"
	"errors"
)

type AuthService struct {
	usersRepository outbound.UsersRepository
	passwordManager outbound.PasswordManager
	tokenManager    outbound.TokenManager
}

func NewAuthService(usersRepository outbound.UsersRepository, passwordManager outbound.PasswordManager, tokenManager outbound.TokenManager) inbound.AuthPort {
	return &AuthService{
		usersRepository: usersRepository,
		passwordManager: passwordManager,
		tokenManager:    tokenManager,
	}
}

func (s *AuthService) Register(ctx context.Context, payload inbound.RegisterPayload) (*inbound.RegisterResult, error) {
	if payload.Password != payload.ConfirmPassword {
		return nil, errs.ErrorPasswordMismatch
	}

	existingUser, err := s.usersRepository.FindByUsername(ctx, payload.Username)
	if err != nil && !errors.Is(err, errs.ErrorUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errs.ErrorUserAlreadyExists
	}

	hashedPassword, err := s.passwordManager.Encode(payload.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: payload.Username,
		Password: hashedPassword,
	}

	if err := s.usersRepository.Save(ctx, user); err != nil {
		return nil, err
	}

	newUser, err := s.usersRepository.FindByUsername(ctx, payload.Username)
	if err != nil {
		return nil, err
	}

	return &inbound.RegisterResult{
		UserID:   newUser.ID,
		Username: newUser.Username,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, payload inbound.LoginPayload) (*inbound.LoginResult, error) {
	user, err := s.usersRepository.FindByUsername(ctx, payload.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.ErrorUserNotFound
	}

	if err := s.passwordManager.Matches(user.Password, payload.Password); err != nil {
		return nil, errs.ErrorInvalidCredentials
	}

	accessToken, refreshToken, err := s.tokenManager.Issue(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &inbound.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Me(ctx context.Context, query inbound.MeQuery) (*inbound.MeResult, error) {
	user, err := s.usersRepository.FindByID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.ErrorUserNotFound
	}
	return &inbound.MeResult{
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, payload inbound.RefreshTokenPayload) (*inbound.RefreshTokenResult, error) {
	claims, err := s.tokenManager.Verify(payload.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.usersRepository.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.ErrorUserNotFound
	}

	accessToken, refreshToken, err := s.tokenManager.Issue(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &inbound.RefreshTokenResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
