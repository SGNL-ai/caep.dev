package caep

import (
	"encoding/json"
	"fmt"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// CredentialType represents the type of credential
type CredentialType string

const (
	CredentialTypePassword      CredentialType = "password"
	CredentialTypePin           CredentialType = "pin"
	CredentialTypeX509          CredentialType = "x509"
	CredentialTypeFIDO2Platform CredentialType = "fido2-platform"
	CredentialTypeFIDO2Roaming  CredentialType = "fido2-roaming"
	CredentialTypeFIDOU2F       CredentialType = "fido-u2f"
	CredentialTypeVerifiable    CredentialType = "verifiable-credential"
	CredentialTypePhoneVoice    CredentialType = "phone-voice"
	CredentialTypePhoneSMS      CredentialType = "phone-sms"
	CredentialTypeApp           CredentialType = "app"
)

// ChangeType represents the type of credential change
type ChangeType string

const (
	ChangeTypeCreate ChangeType = "create"
	ChangeTypeRevoke ChangeType = "revoke"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
)

type CredentialChangePayload struct {
	CredentialType CredentialType `json:"credential_type"`         // REQUIRED
	ChangeType     ChangeType     `json:"change_type"`             // REQUIRED
	FriendlyName   *string        `json:"friendly_name,omitempty"` // OPTIONAL
	X509Issuer     *string        `json:"x509_issuer,omitempty"`   // OPTIONAL
	X509Serial     *string        `json:"x509_serial,omitempty"`   // OPTIONAL
	FIDO2AAGUID    *string        `json:"fido2_aaguid,omitempty"`  // OPTIONAL
}

type CredentialChangeEvent struct {
	BaseCAEPEvent
	CredentialChangePayload
}

func NewCredentialChangeEvent(credType CredentialType, changeType ChangeType) *CredentialChangeEvent {
	e := &CredentialChangeEvent{
		CredentialChangePayload: CredentialChangePayload{
			CredentialType: credType,
			ChangeType:     changeType,
		},
	}

	e.SetType(AssuranceLevelChange)

	return e
}

func (e *CredentialChangeEvent) WithFriendlyName(name string) *CredentialChangeEvent {
	e.FriendlyName = &name

	return e
}

func (e *CredentialChangeEvent) WithX509Details(issuer, serial string) *CredentialChangeEvent {
	e.X509Issuer = &issuer
	e.X509Serial = &serial

	return e
}

func (e *CredentialChangeEvent) WithFIDO2AAGUID(aaguid string) *CredentialChangeEvent {
	e.FIDO2AAGUID = &aaguid

	return e
}

func (e *CredentialChangeEvent) WithEventTimestamp(timestamp int64) *CredentialChangeEvent {
	e.BaseCAEPEvent.WithEventTimestamp(timestamp)

	return e
}

func (e *CredentialChangeEvent) WithInitiatingEntity(entity InitiatingEntity) *CredentialChangeEvent {
	e.BaseCAEPEvent.WithInitiatingEntity(entity)

	return e
}

func (e *CredentialChangeEvent) WithReasonAdmin(language, reason string) *CredentialChangeEvent {
	e.BaseCAEPEvent.WithReasonAdmin(language, reason)

	return e
}

func (e *CredentialChangeEvent) WithReasonUser(language, reason string) *CredentialChangeEvent {
	e.BaseCAEPEvent.WithReasonUser(language, reason)

	return e
}

func (e *CredentialChangeEvent) Validate() error {
	if err := e.ValidateMetadata(); err != nil {
		return err
	}

	if !IsValidCredentialType(e.CredentialType) {
		return event.NewError(event.ErrCodeInvalidValue,
			fmt.Sprintf("invalid credential type: %s", e.CredentialType),
			"credential_type")
	}

	if !IsValidChangeType(e.ChangeType) {
		return event.NewError(event.ErrCodeInvalidValue,
			fmt.Sprintf("invalid change type: %s", e.ChangeType),
			"change_type")
	}

	if (e.X509Issuer != nil && e.X509Serial == nil) ||
		(e.X509Issuer == nil && e.X509Serial != nil) {
		return event.NewError(event.ErrCodeInvalidValue,
			"both x509_issuer and x509_serial must be provided together",
			"x509")
	}

	return nil
}

func (e *CredentialChangeEvent) Payload() interface{} {
	payload := e.CredentialChangePayload

	if e.Metadata != nil {
		return struct {
			CredentialChangePayload
			*EventMetadata
		}{
			CredentialChangePayload: payload,
			EventMetadata:           e.Metadata,
		}
	}

	return payload
}

func (e *CredentialChangeEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Payload())
}

func (e *CredentialChangeEvent) UnmarshalJSON(data []byte) error {
	var payload struct {
		CredentialChangePayload
		*EventMetadata
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return event.NewError(event.ErrCodeParseError,
			"failed to parse credential change event data", "")
	}

	e.SetType(AssuranceLevelChange)

	e.CredentialChangePayload = payload.CredentialChangePayload
	e.Metadata = payload.EventMetadata

	return e.Validate()
}

func IsValidCredentialType(ct CredentialType) bool {
	switch ct {
	case CredentialTypePassword, CredentialTypePin, CredentialTypeX509,
		CredentialTypeFIDO2Platform, CredentialTypeFIDO2Roaming,
		CredentialTypeFIDOU2F, CredentialTypeVerifiable,
		CredentialTypePhoneVoice, CredentialTypePhoneSMS,
		CredentialTypeApp:
		return true
	default:
		return false
	}
}

func IsValidChangeType(ct ChangeType) bool {
	switch ct {
	case ChangeTypeCreate, ChangeTypeRevoke,
		ChangeTypeUpdate, ChangeTypeDelete:
		return true
	default:
		return false
	}
}

func ParseCredentialChangeEvent(data []byte) (event.Event, error) {
	var e CredentialChangeEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, event.NewError(event.ErrCodeParseError,
			"failed to parse credential change event", "")
	}

	return &e, nil
}

func init() {
	event.RegisterEventParser(CredentialChange, ParseCredentialChangeEvent)
}
