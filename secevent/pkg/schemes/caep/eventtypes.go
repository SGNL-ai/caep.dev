package caep

import (
	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// CAEP Event Types as defined in the OpenID Connect CAEP specification
const (
	// EventTypeTokenClaimsChange event type indicates that the claims associated with a token have changed
	EventTypeTokenClaimsChange event.EventType = "https://schemas.openid.net/secevent/caep/event-type/token-claims-change"

	// EventTypeSessionRevoked event type indicates that a session has been terminated
	EventTypeSessionRevoked event.EventType = "https://schemas.openid.net/secevent/caep/event-type/session-revoked"

	// EventTypeCredentialChange event type indicates that credentials have been changed or reset
	EventTypeCredentialChange event.EventType = "https://schemas.openid.net/secevent/caep/event-type/credential-change"

	// EventTypeAssuranceLevelChange event type indicates a change in authentication assurance level
	EventTypeAssuranceLevelChange event.EventType = "https://schemas.openid.net/secevent/caep/event-type/assurance-level-change"

	// EventTypeDeviceComplianceChange event type indicates a change in device compliance status
	EventTypeDeviceComplianceChange event.EventType = "https://schemas.openid.net/secevent/caep/event-type/device-compliance-change"
)

// IsCaepEventType checks if the given event type is a valid CAEP event type
func IsCaepEventType(eventType event.EventType) bool {
	switch eventType {
	case EventTypeTokenClaimsChange,
		EventTypeSessionRevoked,
		EventTypeCredentialChange,
		EventTypeAssuranceLevelChange,
		EventTypeDeviceComplianceChange:
		return true
	default:
		return false
	}
}
