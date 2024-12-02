package ssf

import (
	"encoding/json"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

type VerificationPayload struct {
	State *string `json:"state,omitempty"`
}

type VerificationEvent struct {
	BaseSSFEvent
	VerificationPayload
}

func NewVerificationEvent() *VerificationEvent {
	e := &VerificationEvent{}

	e.SetType(EventTypeVerification)

	return e
}

func (e *VerificationEvent) WithState(state string) *VerificationEvent {
	e.State = &state

	return e
}

func (e *VerificationEvent) GetState() (string, bool) {
	if e.State == nil {
		return "", false
	}

	return *e.State, true
}

func (e *VerificationEvent) Validate() error {
	// State is optional, so no validation needed
	return nil
}

func (e *VerificationEvent) Payload() interface{} {
	return e.VerificationPayload
}

func (e *VerificationEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Payload())
}

func (e *VerificationEvent) UnmarshalJSON(data []byte) error {
	var payload VerificationPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return event.NewError(event.ErrCodeParseError,
			"failed to parse verification event data", "", err.Error())
	}

	e.SetType(EventTypeVerification)

	e.VerificationPayload = payload

	return e.Validate()
}

func ParseVerificationEvent(data []byte) (event.Event, error) {
	var e VerificationEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, event.NewError(event.ErrCodeParseError,
			"failed to parse verification event", "", err.Error())
	}

	return &e, nil
}

func init() {
	event.RegisterEventParser(EventTypeVerification, ParseVerificationEvent)
}
