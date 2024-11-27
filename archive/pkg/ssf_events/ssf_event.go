package ssf_events

import (
	"errors"
	"fmt"
)

type EventType int

const (
	SessionRevoked EventType = iota
	CredentialChange
	DeviceCompliance
	AssuranceLevelChange
	TokenClaimsChange
	VerificationEventType
	StreamUpdatedEventType
)

type SubjectFormat int

const (
	Account SubjectFormat = iota
	Email
	IssuerAndSubject
	Opaque
	PhoneNumber
	DecentralizedIdentifier
	UniqueResourceIdentifier
	Aliases
	ComplexSubject
)

const AccountSubjectFormat = "account"
const EmailSubjectFormat = "email"
const IssuerAndSubjectFormat = "iss_sub"
const OpaqueSubjectFormat = "opaque"
const PhoneNumberSubjectFormat = "phone_number"
const DecentralizedIdentifierSubjectFormat = "did"
const UniqueResourceIdentifierSubjectFormat = "uri"
const AliasesSubjectFormat = "aliases"

// Represents the interface that all SSF Events should implement
//
// See the SessionRevokedEvent (./events/session_revoked_event.go)
// for an example
type SsfEvent interface {
	// Returns the Event URI for the given event
	GetEventUri() string

	// Returns the format of the event's subject
	GetSubjectFormat() SubjectFormat

	// Returns the subject of the event
	GetSubject() map[string]interface{}

	// Returns the Unix timestamp of the event
	GetTimestamp() int64

	// Return the type of event
	GetType() EventType
}

var EventUri = map[EventType]string{
	SessionRevoked:         "https://schemas.openid.net/secevent/caep/event-type/session-revoked",
	CredentialChange:       "https://schemas.openid.net/secevent/caep/event-type/credential-change",
	DeviceCompliance:       "https://schemas.openid.net/secevent/caep/event-type/device-compliance-change",
	AssuranceLevelChange:   "https://schemas.openid.net/secevent/caep/event-type/assurance-level-change",
	TokenClaimsChange:      "https://schemas.openid.net/secevent/caep/event-type/token-claims-change",
	VerificationEventType:  "https://schemas.openid.net/secevent/ssf/event-type/verification",
	StreamUpdatedEventType: "https://schemas.openid.net/secevent/caep/event-type/stream-updated",
}

var EventEnum = map[string]EventType{
	"https://schemas.openid.net/secevent/caep/event-type/session-revoked":          SessionRevoked,
	"https://schemas.openid.net/secevent/caep/event-type/credential-change":        CredentialChange,
	"https://schemas.openid.net/secevent/caep/event-type/device-compliance-change": DeviceCompliance,
	"https://schemas.openid.net/secevent/caep/event-type/assurance-level-change":   AssuranceLevelChange,
	"https://schemas.openid.net/secevent/caep/event-type/token-claims-change":      TokenClaimsChange,
	"https://schemas.openid.net/secevent/ssf/event-type/verification":              VerificationEventType,
	"https://schemas.openid.net/secevent/caep/event-type/stream-updated":           StreamUpdatedEventType,
}

func extractSubject(claimsJson, subjectAttributes map[string]interface{}) (map[string]interface{}, error) {
	if subId, found := claimsJson["sub_id"]; found {
		if mapSubID, ok := subId.(map[string]interface{}); ok {
			return mapSubID, nil
		}
	}

	if subject, found := subjectAttributes["subject"]; found {
		if mapSubject, ok := subject.(map[string]interface{}); ok {
			return mapSubject, nil
		}
	}

	return nil, errors.New("cannot retrieve subject of an event")
}

// Takes an event subject from the JSON of an SSF Event, and converts it into the matching struct for that event
func EventStructFromEvent(eventUri string, eventSubject interface{}, claimsJson map[string]interface{}) (SsfEvent, error) {
	eventEnum := EventEnum[eventUri]

	subjectAttributes, ok := eventSubject.(map[string]interface{})
	if !ok {
		return nil, errors.New("unable to parse event subject")
	}

	// Special Event Types
	if eventEnum == VerificationEventType {
		state, ok := subjectAttributes["state"].(string)
		if !ok {
			return nil, errors.New("unable to parse state")
		}

		event := VerificationEvent{
			Json:  claimsJson,
			State: state,
		}
		return &event, nil

	} else if eventEnum == StreamUpdatedEventType {
		status, ok := subjectAttributes["status"].(string)
		if !ok {
			return nil, errors.New("unable to parse state")
		}

		reason, _ := subjectAttributes["reason"].(string)

		event := StreamUpdatedEvent{
			Json:   claimsJson,
			Status: status,
			Reason: reason,
		}
		return &event, nil
	}

	floatTimestamp, ok := subjectAttributes["event_timestamp"].(float64)
	if !ok {
		return nil, errors.New("unable to parse event timestamp")
	}

	timestamp := int64(floatTimestamp)

	subject, err := extractSubject(claimsJson, subjectAttributes)
	if err != nil {
		return nil, err
	}

	format, err := GetSubjectFormat(subject)
	if err != nil {
		return nil, err
	}

	// Add more Ssf Events as desired
	switch eventEnum {
	case CredentialChange:
		rawCredentialType, ok := subjectAttributes["credential_type"].(string)
		if !ok {
			return nil, errors.New("unable to parse credential type of a credential change event")
		}

		credentialType, ok := CredentialTypesMap[rawCredentialType]
		if !ok {
			return nil, errors.New("received invalid credential type for a credential change event")
		}

		rawChangeType, ok := subjectAttributes["change_type"].(string)
		if !ok {
			return nil, errors.New("unable to parse change type of a credential change event")
		}

		changeType, ok := ChangeTypesMap[rawChangeType]
		if !ok {
			return nil, errors.New("received invalid change type for a credential change event")
		}

		event := CredentialChangeEvent{
			Json:           claimsJson,
			Format:         format,
			Subject:        subject,
			EventTimestamp: timestamp,
			CredentialType: credentialType,
			ChangeType:     changeType,
		}
		return &event, nil

	case SessionRevoked:
		event := SessionRevokedEvent{
			Json:           claimsJson,
			Format:         format,
			Subject:        subject,
			EventTimestamp: timestamp,
		}
		return &event, nil

	case DeviceCompliance:
		previousStatus, ok := subjectAttributes["previousStatus"].(string)
		if !ok {
			return nil, errors.New("unable to parse previous status")
		}

		currentStatus, ok := subjectAttributes["currentStatus"].(string)
		if !ok {
			return nil, errors.New("unable to parse current status")
		}

		event := DeviceComplianceEvent{
			Json:           claimsJson,
			Format:         format,
			Subject:        subject,
			EventTimestamp: timestamp,
			PreviousStatus: previousStatus,
			CurrentStatus:  currentStatus,
		}
		return &event, nil

	case AssuranceLevelChange:
		previousLevel, _ := subjectAttributes["previousLevel"].(string)
		changeDirection, _ := subjectAttributes["changeDirection"].(string)

		currentLevel, ok := subjectAttributes["currentLevel"].(string)
		if !ok {
			return nil, errors.New("unable to parse current level")
		}

		namespace, ok := subjectAttributes["namespace"].(string)
		if !ok {
			return nil, errors.New("unable to parse namespace")
		}

		event := AssuranceLevelChangeEvent{
			Json:            claimsJson,
			Format:          format,
			Subject:         subject,
			EventTimestamp:  timestamp,
			Namespace:       namespace,
			PreviousLevel:   &previousLevel,
			CurrentLevel:    currentLevel,
			ChangeDirection: &changeDirection,
		}
		return &event, nil

	case TokenClaimsChange:
		claims, ok := subjectAttributes["claims"].(map[string]interface{})
		if !ok {
			return nil, errors.New("unable to parse claims")
		}

		event := TokenClaimsChangeEvent{
			Json:           claimsJson,
			Format:         format,
			Subject:        subject,
			EventTimestamp: timestamp,
			Claims:         claims,
		}
		return &event, nil

	default:
		return nil, errors.New("no matching events")
	}
}

func GetSubjectFormat(subject map[string]interface{}) (SubjectFormat, error) {
	format, formatFound := subject["format"]
	formatString := fmt.Sprintf("%v", format)
	if !formatFound {
		return ComplexSubject, nil
	}

	switch formatString {
	case AccountSubjectFormat:
		return Account, nil
	case EmailSubjectFormat:
		return Email, nil
	case IssuerAndSubjectFormat:
		return IssuerAndSubject, nil
	case OpaqueSubjectFormat:
		return Opaque, nil
	case PhoneNumberSubjectFormat:
		return PhoneNumber, nil
	case DecentralizedIdentifierSubjectFormat:
		return DecentralizedIdentifier, nil
	case UniqueResourceIdentifierSubjectFormat:
		return UniqueResourceIdentifier, nil
	case AliasesSubjectFormat:
		return Aliases, nil
	default:
		return -1, errors.New("unable to determine subject format")
	}
}

// Converts a list of Ssf Events to a list of their corresponding Event URI's
func EventTypeArrayToEventUriArray(events []EventType) []string {
	var eventUriArr []string
	for i := 0; i < len(events); i++ {
		eventUriArr = append(eventUriArr, EventUri[events[i]])
	}
	return eventUriArr
}
