// errors.go
package errors

import "fmt"

// AppError represents a structured error in the application.
type AppError struct {
	Code    string // Unique error code for identification
	Message string // Human-readable message for clients
	Err     error  // Underlying error (optional)
}

// New creates a new AppError.
func New(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
/*
// Unwrap allows access to the underlying error for use with errors.Is and errors.As.
func (e *AppError) Unwrap() error {
	return e.Err
}
*/

// Predefined application errors
var (
	ErrUserNotFound     = New("USER_NOT_FOUND", "The requested user does not exist", nil)
	ErrInvalidCredentials = New("INVALID_CREDENTIALS", "Invalid username or password", nil)
	ErrInvalidInput      = New("INVALID_INPUT", "The provided input is invalid", nil)
	ErrInternal          = New("INTERNAL_ERROR", "An internal error occurred", nil)
	ErrExternalAPI       = New("EXTERNAL_API_ERROR", "An error occurred while communicating with an external service", nil)
    ErrDatabase         = New("DATABASE_ERROR", "Database operation failed", nil)
    ErrInvalidOTP       = New("INVALID_OTP", "Invalid OTP Code", nil)
    ErrOTPNotFound      = New("OTP_NOT_FOUND", "OTP not found", nil)
    ErrInvalidMobile    = New("INVALID_MOBILE", "invalid mobile number", nil)
)
/*
// Wrap wraps an existing error with additional context.
func Wrap(err error, code, message string) *AppError {
	return New(code, message, err)
}

// Example usage:
// return errors.Wrap(err, "DB_CONNECTION_ERROR", "Failed to connect to the database")
*/


