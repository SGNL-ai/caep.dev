package caep

import (
	"encoding/json"
	"fmt"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// AssuranceLevel represents NIST Authenticator Assurance Level (AAL)
type AssuranceLevel string

const (
	AssuranceLevelAAL1 AssuranceLevel = "nist-aal1"
	AssuranceLevelAAL2 AssuranceLevel = "nist-aal2"
	AssuranceLevelAAL3 AssuranceLevel = "nist-aal3"
)

// ChangeDirection represents the direction of assurance level change
type ChangeDirection string

const (
	ChangeDirectionIncrease ChangeDirection = "increase"
	ChangeDirectionDecrease ChangeDirection = "decrease"
)

type AssuranceLevelChangePayload struct {
	CurrentLevel    AssuranceLevel  `json:"current_level"`
	PreviousLevel   AssuranceLevel  `json:"previous_level"`
	ChangeDirection ChangeDirection `json:"change_direction"`
}

type AssuranceLevelChangeEvent struct {
	BaseCAEPEvent
	AssuranceLevelChangePayload
}

func NewAssuranceLevelChangeEvent(currentLevel, previousLevel AssuranceLevel, direction ChangeDirection) *AssuranceLevelChangeEvent {
	e := &AssuranceLevelChangeEvent{
		AssuranceLevelChangePayload: AssuranceLevelChangePayload{
			CurrentLevel:    currentLevel,
			PreviousLevel:   previousLevel,
			ChangeDirection: direction,
		},
	}

	e.SetType(EventTypeAssuranceLevelChange)

	return e
}

func (e *AssuranceLevelChangeEvent) WithEventTimestamp(timestamp int64) *AssuranceLevelChangeEvent {
	e.BaseCAEPEvent.WithEventTimestamp(timestamp)

	return e
}

func (e *AssuranceLevelChangeEvent) WithInitiatingEntity(entity InitiatingEntity) *AssuranceLevelChangeEvent {
	e.BaseCAEPEvent.WithInitiatingEntity(entity)

	return e
}

func (e *AssuranceLevelChangeEvent) WithReasonAdmin(language, reason string) *AssuranceLevelChangeEvent {
	e.BaseCAEPEvent.WithReasonAdmin(language, reason)

	return e
}

func (e *AssuranceLevelChangeEvent) WithReasonUser(language, reason string) *AssuranceLevelChangeEvent {
	e.BaseCAEPEvent.WithReasonUser(language, reason)

	return e
}

func (e *AssuranceLevelChangeEvent) GetCurrentLevel() AssuranceLevel {
    return e.CurrentLevel
}

func (e *AssuranceLevelChangeEvent) GetPreviousLevel() AssuranceLevel {
    return e.PreviousLevel
}

func (e *AssuranceLevelChangeEvent) GetChangeDirection() ChangeDirection {
    return e.ChangeDirection
}

func (e *AssuranceLevelChangeEvent) Validate() error {
	if err := e.ValidateMetadata(); err != nil {
		return err
	}

	validLevels := map[AssuranceLevel]bool{
		AssuranceLevelAAL1: true,
		AssuranceLevelAAL2: true,
		AssuranceLevelAAL3: true,
	}

	if !validLevels[e.CurrentLevel] {
		return event.NewError(event.ErrCodeInvalidValue,
			fmt.Sprintf("invalid current level: %s", e.CurrentLevel),
			"current_level", "")
	}

	if !validLevels[e.PreviousLevel] {
		return event.NewError(event.ErrCodeInvalidValue,
			fmt.Sprintf("invalid previous level: %s", e.PreviousLevel),
			"previous_level", "")
	}

	if e.ChangeDirection != ChangeDirectionIncrease && e.ChangeDirection != ChangeDirectionDecrease {
		return event.NewError(event.ErrCodeInvalidValue,
			fmt.Sprintf("invalid change direction: %s", e.ChangeDirection),
			"change_direction", "")
	}

	if e.CurrentLevel == e.PreviousLevel {
		return event.NewError(event.ErrCodeInvalidValue,
			"current and previous levels must be different",
			"levels", "")
	}

	return nil
}

func (e *AssuranceLevelChangeEvent) Payload() interface{} {
	payload := e.AssuranceLevelChangePayload

	if e.Metadata != nil {
		return struct {
			AssuranceLevelChangePayload
			*EventMetadata
		}{
			AssuranceLevelChangePayload: payload,
			EventMetadata:               e.Metadata,
		}
	}

	return payload
}

func (e *AssuranceLevelChangeEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Payload())
}

func (e *AssuranceLevelChangeEvent) UnmarshalJSON(data []byte) error {
	var payload struct {
		AssuranceLevelChangePayload
		*EventMetadata
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return event.NewError(event.ErrCodeParseError,
			"failed to parse assurance level change event data", "", err.Error())
	}

	e.SetType(EventTypeAssuranceLevelChange)

	e.AssuranceLevelChangePayload = payload.AssuranceLevelChangePayload
	e.Metadata = payload.EventMetadata

	return e.Validate()
}

func ParseAssuranceLevelChangeEvent(data []byte) (event.Event, error) {
	var e AssuranceLevelChangeEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, event.NewError(event.ErrCodeParseError,
			"failed to parse assurance level change event", "", err.Error())
	}

	return &e, nil
}

func init() {
	event.RegisterEventParser(EventTypeAssuranceLevelChange, ParseAssuranceLevelChangeEvent)
}
