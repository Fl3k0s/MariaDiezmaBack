package domain

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input data")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden: insufficient permissions")
	ErrInternal          = errors.New("internal error")
	ErrInvalidCredentials= errors.New("invalid email or password")
)
