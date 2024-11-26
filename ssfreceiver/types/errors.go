package types

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConfiguration indicates that the provided configuration is invalid
	ErrInvalidConfiguration = errors.New("invalid configuration")

	// ErrInvalidTransmitterMetadata indicates that the transmitter metadata is invalid
	ErrInvalidTransmitterMetadata = errors.New("invalid transmitter metadata")

	// ErrStreamNotFound indicates that the requested stream was not found
	ErrStreamNotFound = errors.New("stream not found")

	// ErrStreamAlreadyExists indicates that a stream with the same ID already exists
	ErrStreamAlreadyExists = errors.New("stream already exists")

	// ErrAuthorizationFailed indicates that the authorization attempt failed
	ErrAuthorizationFailed = errors.New("authorization failed")

	// ErrOperationNotSupported indicates that the requested operation is not supported
	ErrOperationNotSupported = errors.New("operation not supported")

	// ErrInvalidDeliveryMethod indicates that the delivery method is invalid
	ErrInvalidDeliveryMethod = errors.New("invalid delivery method")

	// ErrInvalidStatus indicates that the stream status is invalid
	ErrInvalidStatus = errors.New("invalid stream status")

	// ErrMaxRetriesExceeded indicates that the maximum number of retries was exceeded
	ErrMaxRetriesExceeded = errors.New("maximum retries exceeded")

	// ErrInvalidSubject indicates that the provided subject is invalid
	ErrInvalidSubject = errors.New("invalid subject")

	// ErrMultipleStreamsFound indicates that multiple streams were found when only one was expected
	ErrMultipleStreamsFound = errors.New("multiple streams found")

	// ErrConfigurationMismatch indicates that the stream configuration doesn't match the requested configuration
	ErrConfigurationMismatch = errors.New("configuration mismatch")

	// ErrInvalidVerificationState indicates that the verification state is invalid
	ErrInvalidVerificationState = errors.New("invalid verification state")
)

// SSFError represents a detailed error with context about what went wrong
type SSFError struct {
	// Err is the underlying error
	Err error

	// Operation is the operation that failed
	Operation string

	// Details provides additional context about the error
	Details string
}

// Error implements the error interface
func (e *SSFError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %v (%s)", e.Operation, e.Err, e.Details)
	}

	return fmt.Sprintf("%s: %v", e.Operation, e.Err)
}

// Unwrap returns the underlying error
func (e *SSFError) Unwrap() error {
	return e.Err
}

func NewError(err error, operation string, details string) error {
	return &SSFError{
		Err:       err,
		Operation: operation,
		Details:   details,
	}
}

func IsInvalidConfiguration(err error) bool {
	return errors.Is(err, ErrInvalidConfiguration)
}

func IsInvalidTransmitterMetadata(err error) bool {
	return errors.Is(err, ErrInvalidTransmitterMetadata)
}

func IsStreamNotFound(err error) bool {
	return errors.Is(err, ErrStreamNotFound)
}

func IsStreamAlreadyExists(err error) bool {
	return errors.Is(err, ErrStreamAlreadyExists)
}

func IsAuthorizationFailed(err error) bool {
	return errors.Is(err, ErrAuthorizationFailed)
}

func IsOperationNotSupported(err error) bool {
	return errors.Is(err, ErrOperationNotSupported)
}

func IsInvalidDeliveryMethod(err error) bool {
	return errors.Is(err, ErrInvalidDeliveryMethod)
}

func IsInvalidStatus(err error) bool {
	return errors.Is(err, ErrInvalidStatus)
}

func IsMaxRetriesExceeded(err error) bool {
	return errors.Is(err, ErrMaxRetriesExceeded)
}

func IsInvalidSubject(err error) bool {
	return errors.Is(err, ErrInvalidSubject)
}

func IsMultipleStreamsFound(err error) bool {
	return errors.Is(err, ErrMultipleStreamsFound)
}

func IsConfigurationMismatch(err error) bool {
	return errors.Is(err, ErrConfigurationMismatch)
}

func IsInvalidVerificationState(err error) bool {
	return errors.Is(err, ErrInvalidVerificationState)
}
