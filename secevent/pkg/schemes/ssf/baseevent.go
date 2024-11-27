package ssf

import (
	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// SSFEvent is the interface that all SSF events must implement
type SSFEvent interface {
	event.Event
}

// BaseSSFEvent provides common SSF event functionality
type BaseSSFEvent struct {
	event.BaseEvent
}

// ValidateSSFPayload is a helper function to validate common SSF payload fields
func ValidateSSFPayload(payload interface{}) error {
	// Currently no common validation for SSF events
	// This function is a placeholder for future common validation logic
	return nil
}
