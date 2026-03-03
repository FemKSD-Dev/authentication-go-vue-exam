package error

import "errors"

var (
	ErrorInvalidCredentials  = errors.New("invalid credentials")
	ErrorUserAlreadyExists   = errors.New("username already exists")
	ErrorPasswordMismatch    = errors.New("password and confirm password do not match")
	ErrorUserNotFound        = errors.New("user not found")
	ErrorInternalServerError = errors.New("internal server error")
	ErrorUnauthorized        = errors.New("unauthorized")
)
