package stream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/internal/retry"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/options"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/types"
)

func (s *stream) GetConfiguration(ctx context.Context, opts ...options.Option) (*types.StreamConfiguration, error) {
	operationOpts := options.Apply(opts...)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s?stream_id=%s", s.metadata.GetConfigurationEndpoint().String(), s.streamID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for k, v := range s.getEndpointHeaders("configuration", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("get-configuration request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var config types.StreamConfiguration
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &config, nil
}

func (s *stream) UpdateConfiguration(ctx context.Context, config *types.StreamConfigurationRequest, opts ...options.Option) (*types.StreamConfiguration, error) {
	if err := config.ValidateRequest(); err != nil {
		return nil, err
	}

	operationOpts := options.Apply(opts...)

	body, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.metadata.GetConfigurationEndpoint().String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for k, v := range s.getEndpointHeaders("configuration", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to update configuration: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("update-configuration request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var updatedConfig types.StreamConfiguration
	if err := json.NewDecoder(resp.Body).Decode(&updatedConfig); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.config = &updatedConfig

	return &updatedConfig, nil
}

func (s *stream) GetStatus(ctx context.Context, opts ...options.Option) (*types.StreamStatus, error) {
	operationOpts := options.Apply(opts...)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s?stream_id=%s", s.metadata.GetStatusEndpoint().String(), s.streamID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for k, v := range s.getEndpointHeaders("status", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("get-status request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var status types.StreamStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validate the received status
	if err := status.Validate(); err != nil {
		return nil, fmt.Errorf("invalid status received: %w", err)
	}

	// Verify the stream ID matches
	if status.StreamID != s.streamID {
		return nil, types.NewError(
			types.ErrInvalidConfiguration,
			"GetStatus",
			fmt.Sprintf("stream ID mismatch: expected %s, got %s", s.streamID, status.StreamID),
		)
	}

	return &status, nil
}

func (s *stream) UpdateStatus(ctx context.Context, status types.StreamStatusType, opts ...options.Option) error {
	operationOpts := options.Apply(opts...)

	request := &types.StreamStatusRequest{
		StreamID: s.streamID,
		Status:   status,
		Reason:   operationOpts.StatusReason,
	}

	if err := request.ValidateRequest(); err != nil {
		return err
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.metadata.GetStatusEndpoint().String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for k, v := range s.getEndpointHeaders("status", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("update-status request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *stream) AddSubject(ctx context.Context, sub subject.Subject, opts ...options.Option) error {
	if s.metadata.GetAddSubjectEndpoint() == nil {
		return fmt.Errorf("add subject endpoint is not configured")
	}

	operationOpts := options.Apply(opts...)

	request := &types.StreamSubjectRequest{
		StreamID: s.streamID,
		Subject:  sub,
		Verified: operationOpts.SubjectVerified,
	}

	if err := request.ValidateRequest(); err != nil {
		return err
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.metadata.GetAddSubjectEndpoint().String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for k, v := range s.getEndpointHeaders("add_subject", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return fmt.Errorf("failed to add subject: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("add-subject request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *stream) RemoveSubject(ctx context.Context, sub subject.Subject, opts ...options.Option) error {
	if s.metadata.GetRemoveSubjectEndpoint() == nil {
		return fmt.Errorf("add subject endpoint is not configured")
	}

	operationOpts := options.Apply(opts...)

	request := &types.StreamSubjectRequest{
		StreamID: s.streamID,
		Subject:  sub,
	}

	if err := request.ValidateRequest(); err != nil {
		return err
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.metadata.GetRemoveSubjectEndpoint().String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for k, v := range s.getEndpointHeaders("remove_subject", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return fmt.Errorf("failed to remove subject: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("remove-subject request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *stream) Verify(ctx context.Context, opts ...options.Option) error {
	operationOpts := options.Apply(opts...)

	request := &types.StreamVerificationRequest{
		StreamID: s.streamID,
		State:    operationOpts.State,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.metadata.GetVerificationEndpoint().String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for k, v := range s.getEndpointHeaders("verification", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return fmt.Errorf("failed to verify stream: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("verification-event-request request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *stream) Delete(ctx context.Context, opts ...options.Option) error {
	operationOpts := options.Apply(opts...)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s?stream_id=%s", s.metadata.GetConfigurationEndpoint().String(), s.streamID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for k, v := range s.getEndpointHeaders("configuration", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return fmt.Errorf("failed to delete stream: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("stream-deletion request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *stream) Pause(ctx context.Context, opts ...options.Option) error {
	return s.UpdateStatus(ctx, types.StatusPaused, opts...)
}

func (s *stream) Resume(ctx context.Context, opts ...options.Option) error {
	return s.UpdateStatus(ctx, types.StatusEnabled, opts...)
}

func (s *stream) Disable(ctx context.Context, opts ...options.Option) error {
	return s.UpdateStatus(ctx, types.StatusDisabled, opts...)
}

func (s *stream) doPoll(ctx context.Context, opts ...options.Option) (map[string]string, error) {
	operationOpts := options.Apply(opts...)

	requestBody := struct {
		StreamID          string   `json:"stream_id"`
		MaxEvents         int      `json:"max_events,omitempty"`
		AckIDs            []string `json:"ack,omitempty"`
		ReturnImmediately bool     `json:"returnImmediately"`
	}{
		StreamID:          s.streamID,
		MaxEvents:         operationOpts.MaxEvents,
		AckIDs:            operationOpts.AckJTIs,
		ReturnImmediately: true,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.Delivery.EndpointURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for k, v := range s.getEndpointHeaders("poll", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for events: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("pull-events request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response map[string]map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	secEvents := response["sets"]

	if operationOpts.AutoAck && len(secEvents) > 0 {
		jtis := make([]string, 0, len(secEvents))
		for jti := range secEvents {
			jtis = append(jtis, jti)
		}

		if err := s.doAcknowledge(ctx, jtis, opts...); err != nil {
			return secEvents, fmt.Errorf("auto-acknowledgment failed: %w", err)
		}
	}

	return secEvents, nil
}

func (s *stream) doAcknowledge(ctx context.Context, jtis []string, opts ...options.Option) error {
	operationOpts := options.Apply(opts...)

	requestBody := struct {
		StreamID          string   `json:"stream_id"`
		AckIDs            []string `json:"ack"`
		ReturnImmediately bool     `json:"returnImmediately"`
	}{
		StreamID:          s.streamID,
		AckIDs:            jtis,
		ReturnImmediately: true,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.Delivery.EndpointURL.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for k, v := range s.getEndpointHeaders("poll", operationOpts.Headers) {
		req.Header.Set(k, v)
	}

	// Add authorization
	if err := s.getAuthorizer(operationOpts).AddAuth(ctx, req); err != nil {
		return fmt.Errorf("failed to add authorization: %w", err)
	}

	operation := retry.Operation(func(ctx context.Context) (*http.Response, error) {
		return s.httpClient.Do(req)
	})

	resp, err := retry.Do(ctx, operation, s.retryConfig)
	if err != nil {
		return fmt.Errorf("failed to acknowledge events: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("events-acknowledgment request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
