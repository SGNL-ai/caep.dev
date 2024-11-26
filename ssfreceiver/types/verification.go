package types

// StreamVerificationRequest represents a request to verify a stream
type StreamVerificationRequest struct {
	// StreamID uniquely identifies the stream
	StreamID string `json:"stream_id"`

	// State is an arbitrary string that the Event Transmitter must echo back
	State string `json:"state,omitempty"`
}
