package stream

import (
	"context"
	"net/http"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/auth"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/internal/config"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/options"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/types"
)

// Stream represents an SSF stream
type Stream interface {
	// Get stream ID
	GetStreamID() string

	// GetMetadata returns the transmitter metadata
	GetMetadata() *types.TransmitterMetadata

	// GetConfiguration returns the current stream configuration
	GetConfiguration() *types.StreamConfiguration

	// UpdateConfiguration updates the stream configuration
	UpdateConfiguration(ctx context.Context, config *types.StreamConfigurationRequest, opts ...options.Option) (*types.StreamConfiguration, error)

	// GetStatus returns the current stream status
	GetStatus(ctx context.Context, opts ...options.Option) (*types.StreamStatus, error)

	// UpdateStatus updates the stream status
	UpdateStatus(ctx context.Context, status types.StreamStatusType, opts ...options.Option) error

	// AddSubject adds a subject to the stream
	AddSubject(ctx context.Context, sub subject.Subject, opts ...options.Option) error

	// RemoveSubject removes a subject from the stream
	RemoveSubject(ctx context.Context, sub subject.Subject, opts ...options.Option) error

	// Verify triggers stream verification
	Verify(ctx context.Context, opts ...options.Option) error

	// Poll retrieves events from the stream (only valid for poll-based streams)
	Poll(ctx context.Context, opts ...options.Option) (map[string]string, error)

	// Acknowledge acknowledges events (only valid for poll-based streams)
	Acknowledge(ctx context.Context, jtis []string, opts ...options.Option) error

	// Delete deletes the stream
	Delete(ctx context.Context, opts ...options.Option) error

	Pause(ctx context.Context, opts ...options.Option) error

	// Resume resumes a paused stream
	Resume(ctx context.Context, opts ...options.Option) error

	// Disable disables the stream
	Disable(ctx context.Context, opts ...options.Option) error
}

// stream implements the Stream interface
type stream struct {
	streamID        string
	metadata        *types.TransmitterMetadata
	config          *types.StreamConfiguration
	authorizer      auth.Authorizer
	retryConfig     config.RetryConfig
	httpClient      *http.Client
	endpointHeaders map[string]map[string]string
}

func NewStream(
	streamID string,
	metadata *types.TransmitterMetadata,
	config *types.StreamConfiguration,
	authorizer auth.Authorizer,
	retryConfig config.RetryConfig,
	httpClient *http.Client,
	endpointHeaders map[string]map[string]string,
) Stream {
	return &stream{
		streamID:        streamID,
		metadata:        metadata,
		config:          config,
		authorizer:      authorizer,
		httpClient:      httpClient,
		endpointHeaders: endpointHeaders,
	}
}

func (s *stream) Poll(ctx context.Context, opts ...options.Option) (map[string]string, error) {
	if !s.config.IsPollDelivery() {
		return nil, types.NewError(
			types.ErrOperationNotSupported,
			"Poll",
			"operation only supported for poll-based streams",
		)
	}

	return s.doPoll(ctx, opts...)
}

func (s *stream) Acknowledge(ctx context.Context, jtis []string, opts ...options.Option) error {
	if !s.config.IsPollDelivery() {
		return types.NewError(
			types.ErrOperationNotSupported,
			"Acknowledge",
			"operation only supported for poll-based streams",
		)
	}

	return s.doAcknowledge(ctx, jtis, opts...)
}

func (s *stream) GetStreamID() string {
	return s.streamID
}

func (s *stream) GetMetadata() *types.TransmitterMetadata {
	return s.metadata
}

func (s *stream) GetConfiguration() *types.StreamConfiguration {
	return s.config
}

// getEndpointHeaders returns the headers for a specific endpoint
func (s *stream) getEndpointHeaders(endpoint string, operationHeaders map[string]string) map[string]string {
	// Start with endpoint-specific headers from stream configuration
	headers := make(map[string]string)
	if endpointHeaders, ok := s.endpointHeaders[endpoint]; ok {
		for k, v := range endpointHeaders {
			headers[k] = v
		}
	}

	// Override with operation-specific headers if provided
	for k, v := range operationHeaders {
		headers[k] = v
	}

	return headers
}

// getAuthorizer returns the authorizer to use for an operation
func (s *stream) getAuthorizer(opts *options.OperationOptions) auth.Authorizer {
	if opts.Auth != nil {
		return opts.Auth
	}

	return s.authorizer
}
