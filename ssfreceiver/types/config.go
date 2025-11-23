package types

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

// StringOrStringArray represents a value that can be either a single string or an array of strings
type StringOrStringArray []string

func (s *StringOrStringArray) UnmarshalJSON(data []byte) error {
	// Try as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = []string{str}

		return nil
	}

	// If that fails, try as string array
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("audience must be either a string or an array of strings: %w", err)
	}

	*s = arr

	return nil
}

func (s StringOrStringArray) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		// If there's only one item, marshal as a string
		return json.Marshal(s[0])
	}

	// Otherwise marshal as an array
	return json.Marshal([]string(s))
}

// StreamConfiguration represents the complete configuration for an SSF stream
type StreamConfiguration struct {
	// streamID uniquely identifies the stream
	streamID string

	// issuer is the URL using the HTTPS scheme that the Transmitter asserts as its issuer Identifier
	issuer       *url.URL
	issuerString string

	// audience contains audience claim that identifies the Event Receiver(s)
	// Can be either a single string or an array of strings
	audience StringOrStringArray

	// delivery contains configuration parameters for the SET delivery method
	delivery DeliveryConfig

	// eventsSupported is the set of events supported by the Transmitter
	eventsSupported []event.EventType

	// eventsRequested is the set of events requested by the Receiver
	eventsRequested []event.EventType

	// eventsDelivered is the set of events that the Transmitter will include in the stream
	eventsDelivered []event.EventType

	// minVerificationInterval is the minimum amount of time in seconds between verification requests
	minVerificationInterval int

	// description describes the properties of the stream
	description string
}

// StreamConfigurationRequest represents the request body for creating or updating a stream
type StreamConfigurationRequest struct {
	// StreamID uniquely identifies the stream (required for updates)
	StreamID string `json:"stream_id,omitempty"`

	// Delivery contains configuration parameters for the SET delivery method
	Delivery *DeliveryConfig `json:"delivery,omitempty"`

	// EventsRequested is the set of events requested by the Receiver
	EventsRequested []event.EventType `json:"events_requested,omitempty"`

	// Description describes the properties of the stream
	Description string `json:"description,omitempty"`
}

// DeliveryConfig represents the delivery configuration for a stream
type DeliveryConfig struct {
	// Method is the specific delivery method to be used (push or poll)
	Method DeliveryMethod `json:"method"`

	// EndpointURL is the location at which the push or poll delivery will take place
	EndpointURL       *url.URL `json:"-"`
	EndpointURLString string   `json:"endpoint_url"`
}

func (d *DeliveryConfig) UnmarshalJSON(data []byte) error {
	type TempDelivery DeliveryConfig

	var temp TempDelivery
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*d = DeliveryConfig(temp)

	var err error
	if d.EndpointURLString != "" {
		if d.EndpointURL, err = url.Parse(d.EndpointURLString); err != nil {
			return fmt.Errorf("invalid endpoint URL: %w", err)
		}
	}

	return nil
}

func (d *DeliveryConfig) MarshalJSON() ([]byte, error) {
	type TempDelivery DeliveryConfig

	temp := TempDelivery(*d)

	if d.EndpointURL != nil {
		temp.EndpointURLString = d.EndpointURL.String()
	}

	return json.Marshal(temp)
}

func (c *StreamConfiguration) UnmarshalJSON(data []byte) error {
	// Create a temporary struct with exported fields for JSON unmarshaling
	var temp struct {
		StreamID                string              `json:"stream_id"`
		Issuer                  string              `json:"iss"`
		Audience                StringOrStringArray `json:"aud"`
		Delivery                DeliveryConfig      `json:"delivery"`
		EventsSupported         []event.EventType   `json:"events_supported,omitempty"`
		EventsRequested         []event.EventType   `json:"events_requested,omitempty"`
		EventsDelivered         []event.EventType   `json:"events_delivered"`
		MinVerificationInterval int                 `json:"min_verification_interval,omitempty"`
		Description             string              `json:"description,omitempty"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Copy values to the actual struct
	c.streamID = temp.StreamID
	c.issuerString = temp.Issuer
	c.audience = temp.Audience
	c.delivery = temp.Delivery
	c.eventsSupported = temp.EventsSupported
	c.eventsRequested = temp.EventsRequested
	c.eventsDelivered = temp.EventsDelivered
	c.minVerificationInterval = temp.MinVerificationInterval
	c.description = temp.Description

	// Parse URL
	var err error
	if temp.Issuer != "" {
		if c.issuer, err = url.Parse(temp.Issuer); err != nil {
			return fmt.Errorf("invalid issuer URL: %w", err)
		}
	}

	return nil
}

func (c *StreamConfiguration) MarshalJSON() ([]byte, error) {
	// Create a temporary struct with exported fields for JSON marshaling
	temp := struct {
		StreamID                string              `json:"stream_id"`
		Issuer                  string              `json:"iss"`
		Audience                StringOrStringArray `json:"aud"`
		Delivery                DeliveryConfig      `json:"delivery"`
		EventsSupported         []event.EventType   `json:"events_supported,omitempty"`
		EventsRequested         []event.EventType   `json:"events_requested,omitempty"`
		EventsDelivered         []event.EventType   `json:"events_delivered"`
		MinVerificationInterval int                 `json:"min_verification_interval,omitempty"`
		Description             string              `json:"description,omitempty"`
	}{
		StreamID:                c.streamID,
		Audience:                c.audience,
		Delivery:                c.delivery,
		EventsSupported:         c.eventsSupported,
		EventsRequested:         c.eventsRequested,
		EventsDelivered:         c.eventsDelivered,
		MinVerificationInterval: c.minVerificationInterval,
		Description:             c.description,
	}

	if c.issuer != nil {
		temp.Issuer = c.issuer.String()
	}

	return json.Marshal(temp)
}

func (c *StreamConfiguration) Validate() error {
	if c.streamID == "" {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateConfig",
			"stream_id is required",
		)
	}

	if c.issuer == nil {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateConfig",
			"issuer is required",
		)
	}

	if len(c.audience) == 0 {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateConfig",
			"at least one audience is required",
		)
	}

	if err := c.delivery.Validate(); err != nil {
		return NewError(
			err,
			"ValidateConfig",
			"invalid delivery configuration",
		)
	}

	if len(c.eventsDelivered) == 0 {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateConfig",
			"at least one delivered event type is required",
		)
	}

	return nil
}

func (d *DeliveryConfig) Validate() error {
	if !IsValidDeliveryMethod(d.Method) {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateDelivery",
			fmt.Sprintf("invalid delivery method: %s", d.Method),
		)
	}

	if d.EndpointURL == nil {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateDelivery",
			"endpoint_url is required",
		)
	}

	return nil
}

func (r *StreamConfigurationRequest) ValidateRequest() error {
	if r.Delivery == nil {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateRequest",
			"delivery configuration is required",
		)
	}

	if err := r.Delivery.Validate(); err != nil {
		return NewError(
			err,
			"ValidateRequest",
			"invalid delivery configuration",
		)
	}

	if len(r.EventsRequested) == 0 {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateRequest",
			"at least one requested event type is required",
		)
	}

	return nil
}

func (c *StreamConfiguration) GetStreamID() string {
	return c.streamID
}

func (c *StreamConfiguration) GetIssuer() *url.URL {
	if c.issuer == nil {
		return nil
	}

	clone := *c.issuer

	return &clone
}

func (c *StreamConfiguration) GetAudience() []string {
	if c.audience == nil {
		return nil
	}

	audience := make([]string, len(c.audience))

	copy(audience, c.audience)

	return audience
}

func (c *StreamConfiguration) GetDeliveryMethod() DeliveryMethod {
	return c.delivery.Method
}

func (c *StreamConfiguration) GetDeliveryEndpoint() *url.URL {
	return c.delivery.EndpointURL
}

func (c *StreamConfiguration) GetEventsSupported() []event.EventType {
	if c.eventsSupported == nil {
		return nil
	}

	events := make([]event.EventType, len(c.eventsSupported))

	copy(events, c.eventsSupported)

	return events
}

func (c *StreamConfiguration) GetEventsRequested() []event.EventType {
	if c.eventsRequested == nil {
		return nil
	}

	events := make([]event.EventType, len(c.eventsRequested))

	copy(events, c.eventsRequested)

	return events
}

func (c *StreamConfiguration) GetEventsDelivered() []event.EventType {
	if c.eventsDelivered == nil {
		return nil
	}

	events := make([]event.EventType, len(c.eventsDelivered))

	copy(events, c.eventsDelivered)

	return events
}

func (c *StreamConfiguration) GetMinVerificationInterval() int {
	return c.minVerificationInterval
}

func (c *StreamConfiguration) GetDescription() string {
	return c.description
}

func (c *StreamConfiguration) IsPollDelivery() bool {
	return c.GetDeliveryMethod() == DeliveryMethodPoll
}

func (c *StreamConfiguration) IsPushDelivery() bool {
	return c.GetDeliveryMethod() == DeliveryMethodPush
}
