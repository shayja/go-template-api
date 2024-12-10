package errors

import "fmt"

// AppError represents a structured application error.
type AppError struct {
    Code    string // Unique error code
    Message string // Human-readable error message
}

// Error implements the error interface.
func (e *AppError) Error() string {
    return fmt.Sprintf("Code: %s, Message: %s", e.Code, e.Message)
}

// New creates a new AppError.
func New(code, message string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
    }
}

// Predefined errors
var (
    ErrUserNotFound = New("USER_NOT_FOUND", "User not found")
    ErrInvalidInput = New("INVALID_INPUT", "Invalid input provided")
    ErrDatabase     = New("DATABASE_ERROR", "Database operation failed")
)