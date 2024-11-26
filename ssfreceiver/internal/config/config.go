package config

import (
	"fmt"
	"time"
)

// RetryConfig represents retry behavior configuration
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// InitialBackoff is the initial delay between retry attempts
	InitialBackoff time.Duration

	// MaxBackoff is the maximum delay between retry attempts
	MaxBackoff time.Duration

	// BackoffMultiplier is the factor by which the backoff increases
	BackoffMultiplier float64

	// RetryableStatus is a map of HTTP status codes that should trigger a retry
	RetryableStatus map[int]bool
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    time.Second,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
		RetryableStatus:   defaultRetryableStatus(),
	}
}

// defaultRetryableStatus returns the default set of retryable HTTP status codes
func defaultRetryableStatus() map[int]bool {
	return map[int]bool{
		408: true, // Request Timeout
		429: true, // Too Many Requests
		500: true, // Internal Server Error
		502: true, // Bad Gateway
		503: true, // Service Unavailable
		504: true, // Gateway Timeout
	}
}

func (r *RetryConfig) Validate() error {
	if r.MaxRetries < 0 {
		return fmt.Errorf("max retries must be non-negative")
	}

	if r.InitialBackoff <= 0 {
		return fmt.Errorf("initial backoff must be positive")
	}

	if r.MaxBackoff < r.InitialBackoff {
		return fmt.Errorf("max backoff must be greater than or equal to initial backoff")
	}

	if r.BackoffMultiplier <= 1.0 {
		return fmt.Errorf("backoff multiplier must be greater than 1.0")
	}

	if r.RetryableStatus == nil {
		return fmt.Errorf("retryable status map cannot be nil")
	}

	for code := range r.RetryableStatus {
		if code < 100 || code > 599 {
			return fmt.Errorf("invalid HTTP status code: %d", code)
		}
	}

	return nil
}
