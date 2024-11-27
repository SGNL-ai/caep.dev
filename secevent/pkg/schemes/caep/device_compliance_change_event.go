package caep

import (
	"encoding/json"
	"fmt"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// ComplianceStatus represents the compliance status of a device
type ComplianceStatus string

const (
	ComplianceStatusCompliant    ComplianceStatus = "compliant"
	ComplianceStatusNotCompliant ComplianceStatus = "not-compliant"
)

type DeviceComplianceChangePayload struct {
	CurrentStatus  ComplianceStatus `json:"current_status"`
	PreviousStatus ComplianceStatus `json:"previous_status"`
}

type DeviceComplianceChangeEvent struct {
	BaseCAEPEvent
	DeviceComplianceChangePayload
}

func NewDeviceComplianceChangeEvent(
	currentStatus, previousStatus ComplianceStatus,
) *DeviceComplianceChangeEvent {
	e := &DeviceComplianceChangeEvent{
		DeviceComplianceChangePayload: DeviceComplianceChangePayload{
			CurrentStatus:  currentStatus,
			PreviousStatus: previousStatus,
		},
	}

	e.SetType(DeviceComplianceChange)

	return e
}

func (e *DeviceComplianceChangeEvent) Validate() error {
	if err := e.ValidateMetadata(); err != nil {
		return err
	}

	validStatuses := map[ComplianceStatus]bool{
		ComplianceStatusCompliant:    true,
		ComplianceStatusNotCompliant: true,
	}

	if !validStatuses[e.CurrentStatus] {
		return event.NewError(event.ErrCodeInvalidValue,
			fmt.Sprintf("invalid current status: %s", e.CurrentStatus),
			"current_status")
	}

	if !validStatuses[e.PreviousStatus] {
		return event.NewError(event.ErrCodeInvalidValue,
			fmt.Sprintf("invalid previous status: %s", e.PreviousStatus),
			"previous_status")
	}

	if e.CurrentStatus == e.PreviousStatus {
		return event.NewError(event.ErrCodeInvalidValue,
			"current and previous status must be different",
			"status")
	}

	return nil
}

func (e *DeviceComplianceChangeEvent) Payload() interface{} {
	payload := e.DeviceComplianceChangePayload

	if e.Metadata != nil {
		return struct {
			DeviceComplianceChangePayload
			*EventMetadata
		}{
			DeviceComplianceChangePayload: payload,
			EventMetadata:                 e.Metadata,
		}
	}

	return payload
}

func (e *DeviceComplianceChangeEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Payload())
}

func (e *DeviceComplianceChangeEvent) UnmarshalJSON(data []byte) error {
	var payload struct {
		DeviceComplianceChangePayload
		*EventMetadata
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return event.NewError(event.ErrCodeParseError,
			"failed to parse device compliance change event data", "")
	}

	e.SetType(DeviceComplianceChange)

	e.DeviceComplianceChangePayload = payload.DeviceComplianceChangePayload
	e.Metadata = payload.EventMetadata

	return e.Validate()
}

func (e *DeviceComplianceChangeEvent) WithEventTimestamp(timestamp int64) *DeviceComplianceChangeEvent {
	e.BaseCAEPEvent.WithEventTimestamp(timestamp)

	return e
}

func (e *DeviceComplianceChangeEvent) WithInitiatingEntity(entity InitiatingEntity) *DeviceComplianceChangeEvent {
	e.BaseCAEPEvent.WithInitiatingEntity(entity)

	return e
}

func (e *DeviceComplianceChangeEvent) WithReasonAdmin(language, reason string) *DeviceComplianceChangeEvent {
	e.BaseCAEPEvent.WithReasonAdmin(language, reason)

	return e
}

func (e *DeviceComplianceChangeEvent) WithReasonUser(language, reason string) *DeviceComplianceChangeEvent {
	e.BaseCAEPEvent.WithReasonUser(language, reason)

	return e
}

func ParseDeviceComplianceChangeEvent(data []byte) (event.Event, error) {
	var e DeviceComplianceChangeEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, event.NewError(event.ErrCodeParseError,
			"failed to parse device compliance change event", "")
	}

	return &e, nil
}

func init() {
	event.RegisterEventParser(DeviceComplianceChange, ParseDeviceComplianceChangeEvent)
}
