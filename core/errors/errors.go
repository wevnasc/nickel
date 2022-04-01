package errors

import "fmt"

type ErrorType string

const (
	NotFound      ErrorType = "NotFound"
	InsertData    ErrorType = "InsertData"
	FindData      ErrorType = "FindData"
	Serialization ErrorType = "Serialization"
)

type AppError struct {
	Type    ErrorType
	Message string
	Reason  error
}

func New(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
	}
}

func Wrap(errorType ErrorType, message string, reason error) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Reason:  reason,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
