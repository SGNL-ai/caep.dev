package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/auth"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/internal/config"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/internal/retry"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/internal/validation"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/stream"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/types"
)

// StreamBuilder configures and creates SSF streams
type StreamBuilder struct {
	metadataURL     *url.URL
	deliveryMethod  types.DeliveryMethod
	pushEndpoint    *url.URL
	eventTypes      []event.EventType
	description     string
	authorizer      auth.Authorizer
	retryConfig     config.RetryConfig
	checkExisting   bool
	httpClient      *http.Client
	endpointHeaders map[string]map[string]string
}

func New(metadataURL string, opts ...Option) (*StreamBuilder, error) {
	parsedURL, err := validation.ParseAndValidateURL(metadataURL)
	if err != nil {
		return nil, types.NewError(
			types.ErrInvalidConfiguration,
			"NewBuilder",
			fmt.Sprintf("invalid metadata URL: %v", err),
		)
	}

	builder := &StreamBuilder{
		metadataURL:     parsedURL,
		retryConfig:     config.DefaultRetryConfig(),
		httpClient:      &http.Client{},
		endpointHeaders: make(map[string]map[string]string),
	}

	for _, opt := range opts {
		if err := opt.apply(builder); err != nil {
			return nil, err
		}
	}

	return builder, nil
}

func (b *StreamBuilder) Setup(ctx context.Context) (stream.Stream, error) {
	// Validate builder configuration
	if err := b.validate(); err != nil {
		return nil, types.NewError(
			types.ErrInvalidConfiguration,
			"Setup",
			err.Error(),
		)
	}

	// Fetch transmitter metadata
	metadata, err := b.fetchTransmitterMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transmitter metadata: %w", err)
	}

	if err := metadata.Validate(); err != nil {
		return nil, types.NewError(
			types.ErrInvalidTransmitterMetadata,
			"Setup",
			err.Error(),
		)
	}

	if !metadata.SupportsDeliveryMethod(b.deliveryMethod) {
		return nil, types.NewError(
			types.ErrInvalidDeliveryMethod,
			"Setup",
			fmt.Sprintf("transmitter does not support delivery method: %s", b.deliveryMethod),
		)
	}

	// Check for existing streams if enabled
	if b.checkExisting {
		stream, err := b.findExistingStream(ctx, metadata)
		if err != nil {
			return nil, err
		}

		if stream != nil {
			return stream, nil
		}
	}

	// Create new stream
	return b.createNewStream(ctx, metadata)
}

func (b *StreamBuilder) findExistingStream(ctx context.Context, metadata *types.TransmitterMetadata) (stream.Stream, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, metadata.GetConfigurationEndpoint().String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err := b.authorizer.AddAuth(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to add authorization: %w", err)
	}

	for k, v := range b.endpointHeaders["metadata"] {
		req.Header.Set(k, v)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return b.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, b.retryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get streams: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// First try to decode as array
	var configs []*types.StreamConfiguration

	err = json.Unmarshal(body, &configs)
	if err == nil {
		// Successfully decoded as array
		if len(configs) == 0 {
			return nil, nil
		}

		if len(configs) > 1 {
			return nil, types.NewError(
				types.ErrInvalidConfiguration,
				"FindExistingStream",
				"multiple streams found: cannot automatically select one",
			)
		}

		config := configs[0]

		return handleStreamConfig(config, b, metadata)
	}

	// If array decode failed, try single object
	var config types.StreamConfiguration
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, fmt.Errorf("failed to decode stream configuration: %w", err)
	}

	return handleStreamConfig(&config, b, metadata)
}

func handleStreamConfig(config *types.StreamConfiguration, b *StreamBuilder, metadata *types.TransmitterMetadata) (stream.Stream, error) {
	if err := config.Validate(); err != nil {
		return nil, types.NewError(
			types.ErrInvalidConfiguration,
			"FindExistingStream",
			fmt.Sprintf("invalid stream configuration: %v", err),
		)
	}

	request := &types.StreamConfigurationRequest{
		Delivery: &types.DeliveryConfig{
			Method:      b.deliveryMethod,
			EndpointURL: b.pushEndpoint,
		},
		EventsRequested: b.eventTypes,
		Description:     b.description,
	}

	if err := validation.ValidateConfigurationMatch(config, request); err != nil {
		return nil, types.NewError(
			types.ErrInvalidConfiguration,
			"FindExistingStream",
			fmt.Sprintf("existing stream configuration mismatch: %v", err),
		)
	}

	return stream.NewStream(
		config.StreamID,
		metadata,
		config,
		b.authorizer,
		b.retryConfig,
		b.httpClient,
		b.endpointHeaders,
	), nil
}

func (b *StreamBuilder) createNewStream(ctx context.Context, metadata *types.TransmitterMetadata) (stream.Stream, error) {
	config := &types.StreamConfigurationRequest{
		Delivery: &types.DeliveryConfig{
			Method:      b.deliveryMethod,
			EndpointURL: b.pushEndpoint,
		},
		EventsRequested: b.eventTypes,
		Description:     b.description,
	}

	body, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal configuration: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, metadata.GetConfigurationEndpoint().String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if err := b.authorizer.AddAuth(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to add authorization: %w", err)
	}

	for k, v := range b.endpointHeaders["configuration"] {
		req.Header.Set(k, v)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return b.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, b.retryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var streamConfig types.StreamConfiguration
	if err := json.NewDecoder(resp.Body).Decode(&streamConfig); err != nil {
		return nil, fmt.Errorf("failed to decode stream configuration: %w", err)
	}

	if err := streamConfig.Validate(); err != nil {
		return nil, types.NewError(
			types.ErrInvalidConfiguration,
			"CreateNewStream",
			fmt.Sprintf("invalid stream configuration received: %v", err),
		)
	}

	return stream.NewStream(
		streamConfig.StreamID,
		metadata,
		&streamConfig,
		b.authorizer,
		b.retryConfig,
		b.httpClient,
		b.endpointHeaders,
	), nil
}

func (b *StreamBuilder) fetchTransmitterMetadata(ctx context.Context) (*types.TransmitterMetadata, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, b.metadataURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range b.endpointHeaders["metadata"] {
		req.Header.Set(k, v)
	}

	if err := b.authorizer.AddAuth(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return b.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, b.retryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var metadata types.TransmitterMetadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return &metadata, nil
}

func (b *StreamBuilder) validate() error {
	if b.authorizer == nil {
		return fmt.Errorf("authorizer is required")
	}

	if b.deliveryMethod == "" {
		return fmt.Errorf("delivery method is required")
	}

	if !types.IsValidDeliveryMethod(b.deliveryMethod) {
		return fmt.Errorf("invalid delivery method: %s", b.deliveryMethod)
	}

	if b.deliveryMethod == types.DeliveryMethodPush && b.pushEndpoint == nil {
		return fmt.Errorf("push endpoint is required for push delivery")
	}

	if err := validation.ValidateEventTypes(b.eventTypes); err != nil {
		return fmt.Errorf("invalid event types: %w", err)
	}

	if err := b.retryConfig.Validate(); err != nil {
		return fmt.Errorf("invalid retry configuration: %w", err)
	}

	return nil
}
