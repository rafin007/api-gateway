package errors

import "errors"

var (
	ErrInternalServerError = errors.New("Something went wrong")
	ErrBadRequest          = errors.New("Error in request body")
	ErrUserAlreadyExists   = errors.New("User already exists")
	ErrUserNotFound        = errors.New("User not found")
	ErrInvalidCredentials  = errors.New("Invalid credentials")
)
