package builder

import (
	"fmt"
	"net/http"

	"github.com/sgnl-ai/caep.dev-receiver/secevent/pkg/event"
	"github.com/sgnl-ai/caep.dev-receiver/ssfreceiver/auth"
	"github.com/sgnl-ai/caep.dev-receiver/ssfreceiver/internal/config"
	"github.com/sgnl-ai/caep.dev-receiver/ssfreceiver/internal/validation"
	"github.com/sgnl-ai/caep.dev-receiver/ssfreceiver/types"
)

// Option represents a builder option
type Option struct {
	apply func(*StreamBuilder) error
}

func WithPollDelivery() Option {
	return Option{func(b *StreamBuilder) error {
		b.deliveryMethod = types.DeliveryMethodPoll

		return nil
	}}
}

func WithPushDelivery(pushEndpoint string) Option {
	return Option{func(b *StreamBuilder) error {
		parsedURL, err := validation.ParseAndValidateURL(pushEndpoint)
		if err != nil {
			return fmt.Errorf("invalid push endpoint URL: %w", err)
		}

		b.deliveryMethod = types.DeliveryMethodPush
		b.pushEndpoint = parsedURL

		return nil
	}}
}

func WithEventTypes(eventTypes []event.EventType) Option {
	return Option{func(b *StreamBuilder) error {
		if err := validation.ValidateEventTypes(eventTypes); err != nil {
			return fmt.Errorf("invalid event types: %w", err)
		}

		b.eventTypes = eventTypes

		return nil
	}}
}

func WithDescription(description string) Option {
	return Option{func(b *StreamBuilder) error {
		b.description = description

		return nil
	}}
}

func WithAuth(authorizer auth.Authorizer) Option {
	return Option{func(b *StreamBuilder) error {
		if authorizer == nil {
			return fmt.Errorf("authorizer cannot be nil")
		}

		b.authorizer = authorizer

		return nil
	}}
}

func WithRetryConfig(cfg config.RetryConfig) Option {
	return Option{func(b *StreamBuilder) error {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid retry configuration: %w", err)
		}

		b.retryConfig = cfg

		return nil
	}}
}

// WithExistingCheck enables checking for existing streams
func WithExistingCheck() Option {
	return Option{func(b *StreamBuilder) error {
		b.checkExisting = true

		return nil
	}}
}

func WithHTTPClient(client *http.Client) Option {
	return Option{func(b *StreamBuilder) error {
		if client == nil {
			return fmt.Errorf("HTTP client cannot be nil")
		}

		b.httpClient = client

		return nil
	}}
}

// WithMetadataEndpointHeaders sets additional headers for metadata endpoint requests
func WithMetadataEndpointHeaders(headers map[string]string) Option {
	return Option{func(b *StreamBuilder) error {
		b.endpointHeaders["metadata"] = headers

		return nil
	}}
}

// WithConfigurationEndpointHeaders sets additional headers for configuration endpoint requests
func WithConfigurationEndpointHeaders(headers map[string]string) Option {
	return Option{func(b *StreamBuilder) error {
		b.endpointHeaders["configuration"] = headers

		return nil
	}}
}

// WithStatusEndpointHeaders sets additional headers for status endpoint requests
func WithStatusEndpointHeaders(headers map[string]string) Option {
	return Option{func(b *StreamBuilder) error {
		b.endpointHeaders["status"] = headers

		return nil
	}}
}

// WithAddSubjectEndpointHeaders sets additional headers for add subject endpoint requests
func WithAddSubjectEndpointHeaders(headers map[string]string) Option {
	return Option{func(b *StreamBuilder) error {
		b.endpointHeaders["add_subject"] = headers

		return nil
	}}
}

// WithRemoveSubjectEndpointHeaders sets additional headers for remove subject endpoint requests
func WithRemoveSubjectEndpointHeaders(headers map[string]string) Option {
	return Option{func(b *StreamBuilder) error {
		b.endpointHeaders["remove_subject"] = headers

		return nil
	}}
}

// WithVerificationEndpointHeaders sets additional headers for verification endpoint requests
func WithVerificationEndpointHeaders(headers map[string]string) Option {
	return Option{func(b *StreamBuilder) error {
		b.endpointHeaders["verification"] = headers

		return nil
	}}
}

// WithEndpointHeaders sets additional headers for all endpoint requests
func WithEndpointHeaders(headers map[string]string) Option {
	return Option{func(b *StreamBuilder) error {
		endpoints := []string{
			"metadata",
			"configuration",
			"status",
			"add_subject",
			"remove_subject",
			"verification",
			"poll",
		}

		for _, endpoint := range endpoints {
			if b.endpointHeaders[endpoint] == nil {
				b.endpointHeaders[endpoint] = make(map[string]string)
			}

			for k, v := range headers {
				b.endpointHeaders[endpoint][k] = v
			}
		}

		return nil
	}}
}
