package event

import (
	"fmt"
)

// ErrorCode represents the type of error that occurred
type ErrorCode string

const (
	// Standard error codes
	ErrCodeInvalidFormat    ErrorCode = "invalid_format"
	ErrCodeMissingValue     ErrorCode = "missing_value"
	ErrCodeInvalidValue     ErrorCode = "invalid_value"
	ErrCodeInvalidEventType ErrorCode = "invalid_event_type"
	ErrCodeMissingField     ErrorCode = "missing_field"
	ErrCodeParseError       ErrorCode = "parse_error"
)

// EventError represents an error that occurred during event operations
type EventError struct {
	Code    ErrorCode
	Message string
	Field   string
	Details string
}

// Error returns the string representation of the error
func (e *EventError) Error() string {
	if e.Field != "" && e.Details != "" {
		return fmt.Sprintf("%s: %s (field: %s, details: %s)", e.Code, e.Message, e.Field, e.Details)
	}

	if e.Field != "" {
		return fmt.Sprintf("%s: %s (field: %s)", e.Code, e.Message, e.Field)
	}

	if e.Details != "" {
		return fmt.Sprintf("%s: %s (details: %s)", e.Code, e.Message, e.Details)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code ErrorCode, msg string, field string, details string) error {
	err := &EventError{
		Code:    code,
		Message: msg,
		Field:   field,
		Details: details,
	}

	return err
}
