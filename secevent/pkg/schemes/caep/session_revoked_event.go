package caep

import (
	"encoding/json"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

type SessionRevokedEvent struct {
	BaseCAEPEvent
	// No additional fields; uses only CAEP metadata
}

func NewSessionRevokedEvent() *SessionRevokedEvent {
	e := &SessionRevokedEvent{}

	e.SetType(EventTypeSessionRevoked)

	return e
}

func (e *SessionRevokedEvent) WithEventTimestamp(timestamp int64) *SessionRevokedEvent {
	e.BaseCAEPEvent.WithEventTimestamp(timestamp)

	return e
}

func (e *SessionRevokedEvent) WithInitiatingEntity(entity InitiatingEntity) *SessionRevokedEvent {
	e.BaseCAEPEvent.WithInitiatingEntity(entity)

	return e
}

func (e *SessionRevokedEvent) WithReasonAdmin(language, reason string) *SessionRevokedEvent {
	e.BaseCAEPEvent.WithReasonAdmin(language, reason)

	return e
}

func (e *SessionRevokedEvent) WithReasonUser(language, reason string) *SessionRevokedEvent {
	e.BaseCAEPEvent.WithReasonUser(language, reason)

	return e
}

func (e *SessionRevokedEvent) Validate() error {
	return e.ValidateMetadata()
}

func (e *SessionRevokedEvent) Payload() interface{} {
	if e.Metadata != nil {
		return struct {
			*EventMetadata
		}{
			EventMetadata: e.Metadata,
		}
	}

	// If no metadata, return empty struct as payload
	return struct{}{}
}

func (e *SessionRevokedEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Payload())
}

func (e *SessionRevokedEvent) UnmarshalJSON(data []byte) error {
	var payload struct {
		*EventMetadata
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return event.NewError(event.ErrCodeParseError,
			"failed to parse session revoked event data", "", err.Error())
	}

	e.SetType(EventTypeSessionRevoked)

	e.Metadata = payload.EventMetadata

	return e.Validate()
}

func ParseSessionRevokedEvent(data []byte) (event.Event, error) {
	var e SessionRevokedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, event.NewError(event.ErrCodeParseError,
			"failed to parse session revoked event", "", err.Error())
	}

	return &e, nil
}

func init() {
	event.RegisterEventParser(EventTypeSessionRevoked, ParseSessionRevokedEvent)
}
