package caep

import (
	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// CAEP Event Types as defined in the OpenID Connect CAEP specification
const (
	// TokenClaimsChange event type indicates that the claims associated with a token have changed
	TokenClaimsChange event.EventType = "https://schemas.openid.net/secevent/caep/event-type/token-claims-change"

	// SessionRevoked event type indicates that a session has been terminated
	SessionRevoked event.EventType = "https://schemas.openid.net/secevent/caep/event-type/session-revoked"

	// CredentialChange event type indicates that credentials have been changed or reset
	CredentialChange event.EventType = "https://schemas.openid.net/secevent/caep/event-type/credential-change"

	// AssuranceLevelChange event type indicates a change in authentication assurance level
	AssuranceLevelChange event.EventType = "https://schemas.openid.net/secevent/caep/event-type/assurance-level-change"

	// DeviceComplianceChange event type indicates a change in device compliance status
	DeviceComplianceChange event.EventType = "https://schemas.openid.net/secevent/caep/event-type/device-compliance-change"
)

// IsCaepEventType checks if the given event type is a valid CAEP event type
func IsCaepEventType(eventType event.EventType) bool {
	switch eventType {
	case TokenClaimsChange,
		SessionRevoked,
		CredentialChange,
		AssuranceLevelChange,
		DeviceComplianceChange:
		return true
	default:
		return false
	}
}
