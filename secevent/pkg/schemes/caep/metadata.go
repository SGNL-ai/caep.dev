package caep

import (
	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// InitiatingEntity represents who/what initiated an event
type InitiatingEntity uint8

const (
	InitiatingEntityAdmin InitiatingEntity = iota
	InitiatingEntityUser
	InitiatingEntityPolicy
	InitiatingEntitySystem
)

// String returns the string representation of the InitiatingEntity
func (i InitiatingEntity) String() string {
	switch i {
	case InitiatingEntityAdmin:
		return "admin"
	case InitiatingEntityUser:
		return "user"
	case InitiatingEntityPolicy:
		return "policy"
	case InitiatingEntitySystem:
		return "system"
	default:
		return "unknown"
	}
}

func (i InitiatingEntity) MarshalJSON() ([]byte, error) {
	return []byte(`"` + i.String() + `"`), nil
}

func (i *InitiatingEntity) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return event.NewError(event.ErrCodeInvalidFormat, "initiating entity must be a string", "initiating_entity")
	}

	// Remove quotes
	str = str[1 : len(str)-1]

	switch str {
	case "admin":
		*i = InitiatingEntityAdmin
	case "user":
		*i = InitiatingEntityUser
	case "policy":
		*i = InitiatingEntityPolicy
	case "system":
		*i = InitiatingEntitySystem
	default:
		return event.NewError(event.ErrCodeInvalidValue, "invalid initiating entity value", "initiating_entity")
	}

	return nil
}

// EventMetadata represents optional metadata for CAEP events
type EventMetadata struct {
	EventTimestamp   *int64            `json:"event_timestamp,omitempty"`
	InitiatingEntity *InitiatingEntity `json:"initiating_entity,omitempty"`
	ReasonAdmin      map[string]string `json:"reason_admin,omitempty"`
	ReasonUser       map[string]string `json:"reason_user,omitempty"`
}

func NewEventMetadata() *EventMetadata {
	return &EventMetadata{}
}

func (m *EventMetadata) WithEventTimestamp(timestamp int64) *EventMetadata {
	m.EventTimestamp = &timestamp

	return m
}

func (m *EventMetadata) WithInitiatingEntity(entity InitiatingEntity) *EventMetadata {
	m.InitiatingEntity = &entity

	return m
}

func (m *EventMetadata) WithReasonAdmin(language, reason string) *EventMetadata {
	if m.ReasonAdmin == nil {
		m.ReasonAdmin = make(map[string]string)
	}

	m.ReasonAdmin[language] = reason

	return m
}

func (m *EventMetadata) WithReasonUser(language, reason string) *EventMetadata {
	if m.ReasonUser == nil {
		m.ReasonUser = make(map[string]string)
	}

	m.ReasonUser[language] = reason

	return m
}

func (m *EventMetadata) GetEventTimestamp() (int64, bool) {
	if m.EventTimestamp == nil {
		return 0, false
	}

	return *m.EventTimestamp, true
}

func (m *EventMetadata) GetInitiatingEntity() (InitiatingEntity, bool) {
	if m.InitiatingEntity == nil {
		return InitiatingEntitySystem, false
	}

	return *m.InitiatingEntity, true
}

func (m *EventMetadata) GetReasonAdmin(language string) (string, bool) {
	if m.ReasonAdmin == nil {
		return "", false
	}

	reason, ok := m.ReasonAdmin[language]

	return reason, ok
}

func (m *EventMetadata) GetReasonUser(language string) (string, bool) {
	if m.ReasonUser == nil {
		return "", false
	}

	reason, ok := m.ReasonUser[language]

	return reason, ok
}

func (m *EventMetadata) GetAllReasonAdmin() map[string]string {
	if m.ReasonAdmin == nil {
		return nil
	}
	return m.ReasonAdmin
}

func (m *EventMetadata) GetAllReasonUser() map[string]string {
	if m.ReasonUser == nil {
		return nil
	}
	return m.ReasonUser
}

func (m *EventMetadata) Validate() error {
	if m.EventTimestamp != nil && *m.EventTimestamp < 0 {
		return event.NewError(event.ErrCodeInvalidValue,
			"event timestamp cannot be negative",
			"event_timestamp")
	}

	return nil
}
