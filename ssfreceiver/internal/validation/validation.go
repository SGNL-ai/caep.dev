package validation

import (
	"fmt"
	"log"
	"net/url"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/types"
)

func ParseAndValidateURL(rawURL string) (*url.URL, error) {
	if rawURL == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return nil, fmt.Errorf("URL must use HTTPS or HTTP scheme")
	}

	return u, nil
}

// EventTypesMatch checks if two sets of event types match exactly
func EventTypesMatch(a, b []event.EventType) bool {
	if len(a) != len(b) {
		return false
	}

	eventMap := make(map[event.EventType]bool, len(a))
	for _, evt := range a {
		eventMap[evt] = true
	}

	for _, evt := range b {
		if !eventMap[evt] {
			return false
		}
	}

	return true
}

// ValidateConfigurationMatch validates if a received stream configuration matches the requested configuration
func ValidateConfigurationMatch(received *types.StreamConfiguration, requested *types.StreamConfigurationRequest) error {
	// Delivery method is critical - return error if mismatched
	if received.GetDeliveryMethod() != requested.Delivery.Method {
		return fmt.Errorf("delivery method mismatch: received %s, requested %s",
			received.GetDeliveryMethod(), requested.Delivery.Method)
	}

	// For push delivery, endpoint URL is critical - return error if mismatched
	if requested.Delivery.Method == types.DeliveryMethodPush {
		if received.GetDeliveryEndpoint().String() != requested.Delivery.EndpointURL.String() {
			return fmt.Errorf("endpoint URL mismatch: received %s, requested %s",
				received.GetDeliveryEndpoint().String(), requested.Delivery.EndpointURL.String())
		}
	}

	// Log warning for event types mismatch
	if !EventTypesMatch(received.GetEventsRequested(), requested.EventsRequested) {
		log.Printf("WARNING: Event types mismatch - received: %v, requested: %v",
			received.GetEventsRequested(), requested.EventsRequested)
	}

	// Log warning for description mismatch
	if received.GetDescription() != requested.Description {
		log.Printf("WARNING: Description mismatch - received: %q, requested: %q",
			received.GetDescription(), requested.Description)
	}

	return nil
}

func ValidateEventTypes(eventTypes []event.EventType) error {
	if len(eventTypes) == 0 {
		return fmt.Errorf("at least one event type is required")
	}

	for _, et := range eventTypes {
		if et == "" {
			return fmt.Errorf("event type cannot be empty")
		}
	}

	return nil
}
