package errors

import (
	"errors"
)

// Converts service error into AppError
func MapServiceError(err error) *AppError {
	switch {
	case errors.Is(err, ErrInternalServerError):
		return InternalServerError(err.Error())
	case errors.Is(err, ErrBadRequest):
		return BadRequest(err.Error())
	case errors.Is(err, ErrUserAlreadyExists):
		return Conflict(err.Error())
	case errors.Is(err, ErrInvalidCredentials):
		return InvalidCredentials(err.Error())
	default:
		return InternalServerError(ErrInternalServerError.Error())
	}
}
