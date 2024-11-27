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
	// StreamID uniquely identifies the stream
	StreamID string `json:"stream_id"`

	// Issuer is the URL using the HTTPS scheme that the Transmitter asserts as its Issuer Identifier
	Issuer       *url.URL `json:"-"`
	IssuerString string   `json:"iss"`

	// Audience contains audience claim that identifies the Event Receiver(s)
	// Can be either a single string or an array of strings
	Audience StringOrStringArray `json:"aud"`

	// Delivery contains configuration parameters for the SET delivery method
	Delivery DeliveryConfig `json:"delivery"`

	// EventsSupported is the set of events supported by the Transmitter
	EventsSupported []event.EventType `json:"events_supported,omitempty"`

	// EventsRequested is the set of events requested by the Receiver
	EventsRequested []event.EventType `json:"events_requested,omitempty"`

	// EventsDelivered is the set of events that the Transmitter will include in the stream
	EventsDelivered []event.EventType `json:"events_delivered"`

	// MinVerificationInterval is the minimum amount of time in seconds between verification requests
	MinVerificationInterval int `json:"min_verification_interval,omitempty"`

	// Description describes the properties of the stream
	Description string `json:"description,omitempty"`
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
	type TempConfig StreamConfiguration

	var temp TempConfig
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*c = StreamConfiguration(temp)

	var err error
	if c.IssuerString != "" {
		if c.Issuer, err = url.Parse(c.IssuerString); err != nil {
			return fmt.Errorf("invalid issuer URL: %w", err)
		}
	}

	return nil
}

func (c *StreamConfiguration) MarshalJSON() ([]byte, error) {
	type TempConfig StreamConfiguration

	temp := TempConfig(*c)

	if c.Issuer != nil {
		temp.IssuerString = c.Issuer.String()
	}

	return json.Marshal(temp)
}

func (c *StreamConfiguration) Validate() error {
	if c.StreamID == "" {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateConfig",
			"stream_id is required",
		)
	}

	if c.Issuer == nil {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateConfig",
			"issuer is required",
		)
	}

	if len(c.Audience) == 0 {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateConfig",
			"at least one audience is required",
		)
	}

	if err := c.Delivery.Validate(); err != nil {
		return NewError(
			err,
			"ValidateConfig",
			"invalid delivery configuration",
		)
	}

	if len(c.EventsDelivered) == 0 {
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

func (d *DeliveryConfig) IsPollDelivery() bool {
	return d.Method == DeliveryMethodPoll
}

func (d *DeliveryConfig) IsPushDelivery() bool {
	return d.Method == DeliveryMethodPush
}
