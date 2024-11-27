package options

import (
	"github.com/sgnl-ai/caep.dev/ssfreceiver/auth"
)

// OperationOptions holds common options for stream operations
type OperationOptions struct {
	// Auth overrides the default authorizer for this operation
	Auth auth.Authorizer

	// MaxEvents sets the maximum number of events to retrieve in a poll operation
	MaxEvents int

	// AutoAck enables automatic acknowledgment of events in a poll operation
	AutoAck bool

	// AckJTIs sets specific JTIs to acknowledge in a poll operation
	AckJTIs []string

	// State sets the state for stream verification
	State string

	// SubjectVerified indicates whether a subject has been verified
	SubjectVerified bool

	// Headers specifies additional HTTP headers for the request
	Headers map[string]string

	// StatusReason provides a reason for status change
	StatusReason string
}

// Option represents a function that modifies OperationOptions
type Option func(*OperationOptions)

func DefaultOptions() *OperationOptions {
	return &OperationOptions{
		MaxEvents: 100,
		AutoAck:   false,
	}
}

func WithAuth(auth auth.Authorizer) Option {
	return func(o *OperationOptions) {
		o.Auth = auth
	}
}

func WithMaxEvents(max int) Option {
	return func(o *OperationOptions) {
		o.MaxEvents = max
	}
}

func WithAutoAck(autoAck bool) Option {
	return func(o *OperationOptions) {
		o.AutoAck = autoAck
	}
}

func WithAckJTIs(jtis []string) Option {
	return func(o *OperationOptions) {
		o.AckJTIs = jtis
	}
}

func WithState(state string) Option {
	return func(o *OperationOptions) {
		o.State = state
	}
}

func WithSubjectVerification(verified bool) Option {
	return func(o *OperationOptions) {
		o.SubjectVerified = verified
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(o *OperationOptions) {
		o.Headers = headers
	}
}

func WithStatusReason(reason string) Option {
	return func(o *OperationOptions) {
		o.StatusReason = reason
	}
}

func Apply(opts ...Option) *OperationOptions {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}
