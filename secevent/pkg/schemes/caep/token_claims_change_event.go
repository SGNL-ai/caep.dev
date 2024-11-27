package caep

import (
	"encoding/json"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

type TokenClaimsChangePayload struct {
	Claims map[string]interface{} `json:"claims"` // REQUIRED
}

type TokenClaimsChangeEvent struct {
	BaseCAEPEvent
	TokenClaimsChangePayload
}

func NewTokenClaimsChangeEvent() *TokenClaimsChangeEvent {
	e := &TokenClaimsChangeEvent{
		TokenClaimsChangePayload: TokenClaimsChangePayload{
			Claims: make(map[string]interface{}),
		},
	}

	e.SetType(TokenClaimsChange)

	return e
}

func (e *TokenClaimsChangeEvent) WithClaim(name string, value interface{}) *TokenClaimsChangeEvent {
	e.Claims[name] = value

	return e
}

func (e *TokenClaimsChangeEvent) Validate() error {
	if err := e.ValidateMetadata(); err != nil {
		return err
	}

	if len(e.Claims) == 0 {
		return event.NewError(event.ErrCodeMissingValue,
			"claims are required and must not be empty",
			"claims")
	}

	return nil
}

func (e *TokenClaimsChangeEvent) Payload() interface{} {
	payload := e.TokenClaimsChangePayload

	if e.Metadata != nil {
		return struct {
			TokenClaimsChangePayload
			*EventMetadata
		}{
			TokenClaimsChangePayload: payload,
			EventMetadata:            e.Metadata,
		}
	}

	return payload
}

func (e *TokenClaimsChangeEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Payload())
}

func (e *TokenClaimsChangeEvent) UnmarshalJSON(data []byte) error {
	var payload struct {
		TokenClaimsChangePayload
		*EventMetadata
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return event.NewError(event.ErrCodeParseError,
			"failed to parse token claims change event data", "")
	}

	e.SetType(TokenClaimsChange)

	e.TokenClaimsChangePayload = payload.TokenClaimsChangePayload
	e.Metadata = payload.EventMetadata

	return e.Validate()
}

func (e *TokenClaimsChangeEvent) WithEventTimestamp(timestamp int64) *TokenClaimsChangeEvent {
	e.BaseCAEPEvent.WithEventTimestamp(timestamp)

	return e
}

func (e *TokenClaimsChangeEvent) WithInitiatingEntity(entity InitiatingEntity) *TokenClaimsChangeEvent {
	e.BaseCAEPEvent.WithInitiatingEntity(entity)

	return e
}

func (e *TokenClaimsChangeEvent) WithReasonAdmin(language, reason string) *TokenClaimsChangeEvent {
	e.BaseCAEPEvent.WithReasonAdmin(language, reason)

	return e
}

func (e *TokenClaimsChangeEvent) WithReasonUser(language, reason string) *TokenClaimsChangeEvent {
	e.BaseCAEPEvent.WithReasonUser(language, reason)

	return e
}

func ParseTokenClaimsChangeEvent(data []byte) (event.Event, error) {
	var e TokenClaimsChangeEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, event.NewError(event.ErrCodeParseError,
			"failed to parse token claims change event", "")
	}

	return &e, nil
}

func init() {
	event.RegisterEventParser(TokenClaimsChange, ParseTokenClaimsChangeEvent)
}
