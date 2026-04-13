package apperrors

import "fmt"

type ErrorCode string

const (
	CodeNotFound     ErrorCode = "NOT_FOUND"
	CodeInvalidInput ErrorCode = "INVALID_INPUT"
	CodeInternal     ErrorCode = "INTERNAL_ERROR"
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeDbError      ErrorCode = "DB_ERROR"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(code ErrorCode, message string, details ...string) *AppError {
	d := ""
	if len(details) > 0 {
		d = details[0]
	}
	return &AppError{
		Code:    code,
		Message: message,
		Details: d,
	}
}
