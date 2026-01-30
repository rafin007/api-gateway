package errors

import "net/http"

// AppError is responsible for client facing errors
type AppError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func (e *AppError) Error() string {
	return e.Message
}

func BadRequest(msg string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func InternalServerError(msg string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: msg,
	}
}

func Conflict(msg string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: msg,
	}
}

func InvalidCredentials(msg string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: msg,
	}
}

func ValidationError(errors map[string]string) *AppError {
	return &AppError{
		Code:   http.StatusBadRequest,
		Errors: errors,
	}
}
