package errors

import (
	"fmt"
)

type AppError struct {
	Status  int    `json:"-"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewAppError(status int, code int, msg string) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: msg,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Status: %d, Code: %d, Message: %s", e.Status, e.Code, e.Message)
}

// Predefined errors
var (
	ErrNotFound             = NewAppError(404, 10111, "Resource not found")
	ErrInternalServer       = NewAppError(500, 10222, "Internal server error")
	ErrBadRequest           = NewAppError(400, 10333, "Bad request")
	ErrUnauthorized         = NewAppError(401, 10444, "Unauthorized")
	ErrSeeksterUnauthorize  = NewAppError(401, 10445, "Seekster unauthorize")
	ErrSeeksterUserExist    = NewAppError(400, 10555, "Seekster user already exists")
	ErrRedis                = NewAppError(500, 10666, "Internal server error")
	ErrParseJSON            = NewAppError(500, 10777, "Internal server error")
	ErrDatabase             = NewAppError(500, 10888, "Internal server error")
	ErrExtractJWTTrueID     = NewAppError(500, 10999, "Internal server error")
	ErrValidationInput      = NewAppError(400, 11000, "Validation Input error")
	ErrValidationInputSSOID = NewAppError(400, 11001, "Validation Input error")
	ErrValidationModel      = NewAppError(400, 11002, "Validation Model error")
	ErrIoReader             = NewAppError(500, 11003, "Internal server error")
	// ... and so on for other common errors
)
