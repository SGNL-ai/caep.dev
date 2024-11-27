package ssf

import (
	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// SSF Event Types as defined in the OpenID SSF specification
const (
	// EventTypeVerification represents a verification event used to verify stream configuration
	EventTypeVerification event.EventType = "https://schemas.openid.net/secevent/ssf/event-type/verification"

	// EventTypeStreamUpdate represents a stream status update event
	EventTypeStreamUpdate event.EventType = "https://schemas.openid.net/secevent/ssf/event-type/stream-updated"
)
