package ssf

import (
	"encoding/json"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// StreamStatus represents the possible states of a stream
type StreamStatus string

const (
	// StreamStatusEnabled indicates the stream is active and transmitting events
	StreamStatusEnabled StreamStatus = "enabled"

	// StreamStatusPaused indicates the stream is temporarily paused
	StreamStatusPaused StreamStatus = "paused"

	// StreamStatusDisabled indicates the stream is permanently disabled
	StreamStatusDisabled StreamStatus = "disabled"
)

func ValidateStreamStatus(status StreamStatus) bool {
	switch status {
	case StreamStatusEnabled,
		StreamStatusPaused,
		StreamStatusDisabled:
		return true
	default:
		return false
	}
}

type StreamUpdatePayload struct {
	Status StreamStatus `json:"status"`
	Reason *string      `json:"reason,omitempty"`
}

type StreamUpdateEvent struct {
	BaseSSFEvent
	StreamUpdatePayload
}

func NewStreamUpdateEvent(status StreamStatus) *StreamUpdateEvent {
	e := &StreamUpdateEvent{
		StreamUpdatePayload: StreamUpdatePayload{
			Status: status,
		},
	}

	e.SetType(EventTypeStreamUpdate)

	return e
}

func (e *StreamUpdateEvent) WithReason(reason string) *StreamUpdateEvent {
	e.Reason = &reason

	return e
}

func (e *StreamUpdateEvent) GetReason() (string, bool) {
	if e.Reason == nil {
		return "", false
	}

	return *e.Reason, true
}

func (e *StreamUpdateEvent) GetStatus() StreamStatus {
	return e.Status
}

func (e *StreamUpdateEvent) Validate() error {
	if !ValidateStreamStatus(e.Status) {
		return event.NewError(event.ErrCodeInvalidValue,
			"invalid stream status",
			"status")
	}

	return nil
}

func (e *StreamUpdateEvent) Payload() interface{} {
	return e.StreamUpdatePayload
}

func (e *StreamUpdateEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Payload())
}

func (e *StreamUpdateEvent) UnmarshalJSON(data []byte) error {
	var payload StreamUpdatePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return event.NewError(event.ErrCodeParseError,
			"failed to parse stream update event data", "")
	}

	e.SetType(EventTypeStreamUpdate)

	e.StreamUpdatePayload = payload

	return e.Validate()
}

func ParseStreamUpdateEvent(data []byte) (event.Event, error) {
	var e StreamUpdateEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, event.NewError(event.ErrCodeParseError,
			"failed to parse stream update event", "")
	}

	return &e, nil
}

func init() {
	event.RegisterEventParser(EventTypeStreamUpdate, ParseStreamUpdateEvent)
}
