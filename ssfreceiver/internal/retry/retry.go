package retry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sgnl-ai/caep.dev-receiver/ssfreceiver/internal/config"
	"github.com/sgnl-ai/caep.dev-receiver/ssfreceiver/types"
)

// Operation represents a retryable operation
type Operation func(context.Context) (*http.Response, error)

// Do executes the given operation with retry logic based on the provided configuration
func Do(ctx context.Context, op Operation, cfg config.RetryConfig) (*http.Response, error) {
	var lastResp *http.Response
	var lastErr error

	backoff := newBackoff(cfg.InitialBackoff, cfg.MaxBackoff, cfg.BackoffMultiplier)

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			if lastErr != nil {
				return nil, fmt.Errorf("context cancelled after %d attempts: %w", attempt, lastErr)
			}

			return nil, ctx.Err()
		default:
			if attempt > 0 {
				// Add delay
				delay := backoff.next()
				timer := time.NewTimer(delay)

				select {
				case <-ctx.Done():
					timer.Stop()

					if lastErr != nil {
						return nil, fmt.Errorf("context cancelled after %d attempts: %w", attempt, lastErr)
					}

					return nil, ctx.Err()
				case <-timer.C:
				}
			}

			var resp *http.Response

			resp, lastErr = op(ctx)
			if lastErr != nil {
				continue
			}

			if lastResp != nil {
				lastResp.Body.Close()
			}

			lastResp = resp

			// Check if status code is retryable
			if lastResp != nil {
				if !cfg.RetryableStatus[lastResp.StatusCode] {
					return lastResp, nil
				}

				lastErr = fmt.Errorf("received retryable status %d", lastResp.StatusCode)

				continue
			}

			return lastResp, nil
		}
	}

	return lastResp, types.NewError(
		types.ErrMaxRetriesExceeded,
		"Retry",
		fmt.Sprintf("operation failed after %d attempts: %v", cfg.MaxRetries, lastErr),
	)
}
