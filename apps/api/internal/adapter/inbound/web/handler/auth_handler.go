package handler

import (
	webRequest "authentication-project-exam/internal/adapter/inbound/web/request"
	webResponse "authentication-project-exam/internal/adapter/inbound/web/response"
	"authentication-project-exam/internal/config"
	errorx "authentication-project-exam/internal/core/error"
	"authentication-project-exam/internal/core/port/inbound"
	"authentication-project-exam/internal/core/port/outbound"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	auth         inbound.AuthPort
	config       *config.Config
	tokenManager outbound.TokenManager
}

func NewAuthHandler(auth inbound.AuthPort, cfg *config.Config) *AuthHandler {
	return &AuthHandler{auth: auth, config: cfg}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req webRequest.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(webResponse.Fail[any](webResponse.CodeBadRequest, "invalid request body"))
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" || req.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(webResponse.Fail[any](webResponse.CodeBadRequest, "username, password, and confirm password are required"))
	}

	_, err := h.auth.Register(c.Context(), inbound.RegisterPayload{
		Username:        req.Username,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		return h.handleRegisterError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(webResponse.Success(webResponse.RegisterResponse{}, "User registered successfully"))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req webRequest.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(webResponse.Fail[any](webResponse.CodeBadRequest, "invalid request body"))
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(webResponse.Fail[any](webResponse.CodeBadRequest, "username and password are required"))
	}

	result, err := h.auth.Login(c.Context(), inbound.LoginPayload{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return h.handleLoginError(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Minute * time.Duration(h.config.JWTExpireMinutes)),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Minute * time.Duration(h.config.RefreshTokenExpireMinutes)),
	})

	return c.Status(fiber.StatusOK).JSON(webResponse.Success(webResponse.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, "Login successful"))
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(webResponse.Fail[any](webResponse.CodeUnauthorized, "unauthorized"))
	}

	result, err := h.auth.Me(c.Context(), inbound.MeQuery{
		UserID: userID,
	})
	if err != nil {
		return h.handleMeError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(webResponse.Success(webResponse.MeResponse{
		UserID:   result.UserID,
		Username: result.Username,
	}, "User information retrieved successfully"))
}

func (h *AuthHandler) handleRegisterError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, errorx.ErrorPasswordMismatch):
		return c.Status(fiber.StatusBadRequest).JSON(webResponse.Fail[any](webResponse.CodePasswordMismatch, err.Error()))
	case errors.Is(err, errorx.ErrorUserAlreadyExists):
		return c.Status(fiber.StatusBadRequest).JSON(webResponse.Fail[any](webResponse.CodeUsernameExists, err.Error()))
	case errors.Is(err, errorx.ErrorUserNotFound):
		return c.Status(fiber.StatusBadRequest).JSON(webResponse.Fail[any](webResponse.CodeUserNotFound, err.Error()))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(webResponse.Fail[any](webResponse.CodeInternalError, err.Error()))
	}
}

func (h *AuthHandler) handleLoginError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, errorx.ErrorInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(webResponse.Fail[any](webResponse.CodeUnauthorized, err.Error()))
	case errors.Is(err, errorx.ErrorUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(webResponse.Fail[any](webResponse.CodeUserNotFound, err.Error()))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(webResponse.Fail[any](webResponse.CodeInternalError, err.Error()))
	}
}

func (h *AuthHandler) handleMeError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, errorx.ErrorUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(webResponse.Fail[any](webResponse.CodeUserNotFound, err.Error()))
	case errors.Is(err, errorx.ErrorUnauthorized):
		return c.Status(fiber.StatusUnauthorized).JSON(webResponse.Fail[any](webResponse.CodeUnauthorized, err.Error()))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(webResponse.Fail[any](webResponse.CodeInternalError, err.Error()))
	}
}
