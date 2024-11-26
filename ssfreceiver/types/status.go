package types

import (
	"fmt"
)

// StreamStatus represents the current status of an SSF stream
type StreamStatus struct {
	// StreamID uniquely identifies the stream
	StreamID string `json:"stream_id"`

	// Status represents the current state of the stream
	Status StreamStatusType `json:"status"`

	// Reason optionally explains why the stream's status is set to the current value
	Reason string `json:"reason,omitempty"`
}

// StreamStatusType represents the possible status values for a stream
type StreamStatusType string

const (
	// StatusEnabled indicates the Transmitter must transmit events over the stream
	StatusEnabled StreamStatusType = "enabled"

	// StatusPaused indicates the Transmitter must not transmit events but will hold them
	StatusPaused StreamStatusType = "paused"

	// StatusDisabled indicates the Transmitter must not transmit events and will not hold them
	StatusDisabled StreamStatusType = "disabled"
)

// StreamStatusRequest represents a request to update a stream's status
type StreamStatusRequest struct {
	// StreamID uniquely identifies the stream
	StreamID string `json:"stream_id"`

	// Status represents the desired state of the stream
	Status StreamStatusType `json:"status"`

	// Reason optionally explains why the stream's status is being changed
	Reason string `json:"reason,omitempty"`
}

func (s *StreamStatus) Validate() error {
	if s.StreamID == "" {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateStatus",
			"stream_id is required",
		)
	}

	if !s.Status.IsValid() {
		return NewError(
			ErrInvalidStatus,
			"ValidateStatus",
			fmt.Sprintf("invalid status: %s", s.Status),
		)
	}

	return nil
}

func (r *StreamStatusRequest) ValidateRequest() error {
	if r.StreamID == "" {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateStatusRequest",
			"stream_id is required",
		)
	}

	if !r.Status.IsValid() {
		return NewError(
			ErrInvalidStatus,
			"ValidateStatusRequest",
			fmt.Sprintf("invalid status: %s", r.Status),
		)
	}

	return nil
}

func (s StreamStatusType) IsValid() bool {
	switch s {
	case StatusEnabled, StatusPaused, StatusDisabled:
		return true
	default:
		return false
	}
}

func (s StreamStatusType) String() string {
	return string(s)
}

func (s StreamStatusType) IsEnabled() bool {
	return s == StatusEnabled
}

func (s StreamStatusType) IsPaused() bool {
	return s == StatusPaused
}

func (s StreamStatusType) IsDisabled() bool {
	return s == StatusDisabled
}
