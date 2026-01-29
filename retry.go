package filemaker

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// RetryConfig defines the configuration for retry behavior.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts.
	// Default: 3
	MaxRetries int

	// MinWaitTime is the minimum wait time between retries.
	// Default: 1 second
	MinWaitTime time.Duration

	// MaxWaitTime is the maximum wait time between retries.
	// Default: 30 seconds
	MaxWaitTime time.Duration

	// EnableJitter adds randomization to wait times to prevent thundering herd.
	// Default: true
	EnableJitter bool

	// RetryableStatusCodes defines which HTTP status codes should trigger a retry.
	// Default: 429, 500, 502, 503, 504
	RetryableStatusCodes []int

	// OnRetry is an optional callback function called before each retry attempt.
	// It receives the attempt number and the error that caused the retry.
	OnRetry func(attempt int, err error)
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:   3,
		MinWaitTime:  1 * time.Second,
		MaxWaitTime:  30 * time.Second,
		EnableJitter: true,
		RetryableStatusCodes: []int{
			429, // Too Many Requests
			500, // Internal Server Error
			502, // Bad Gateway
			503, // Service Unavailable
			504, // Gateway Timeout
		},
	}
}

// shouldRetry determines if an error is retryable based on the retry configuration.
func (rc *RetryConfig) shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	// Use the IsRetryable helper from errors.go
	return IsRetryable(err)
}

// calculateBackoff calculates the wait time for a given attempt using exponential backoff.
// Formula: min(minWait * 2^attempt, maxWait) with optional jitter.
func (rc *RetryConfig) calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: 2^attempt
	exp := math.Pow(2, float64(attempt))
	backoff := time.Duration(float64(rc.MinWaitTime) * exp)

	// Cap at maximum wait time
	if backoff > rc.MaxWaitTime {
		backoff = rc.MaxWaitTime
	}

	// Add jitter if enabled to prevent thundering herd
	if rc.EnableJitter {
		// Add random jitter between 0-25% of backoff time
		jitter := time.Duration(rand.Float64() * float64(backoff) * 0.25)
		backoff = backoff + jitter
	}

	return backoff
}

// executeWithRetry wraps a function with retry logic.
// It will retry the function up to MaxRetries times if it returns a retryable error.
func (rc *RetryConfig) executeWithRetry(ctx context.Context, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= rc.MaxRetries; attempt++ {
		// Check if context is cancelled before attempting
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute the function
		err := fn()

		// Success - no retry needed
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should retry this error
		if !rc.shouldRetry(err) {
			return err
		}

		// Don't retry after the last attempt
		if attempt >= rc.MaxRetries {
			return err
		}

		// Calculate backoff time
		backoffTime := rc.calculateBackoff(attempt)

		// Call OnRetry callback if provided
		if rc.OnRetry != nil {
			rc.OnRetry(attempt+1, err)
		}

		// Wait with context cancellation support
		timer := time.NewTimer(backoffTime)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			// Continue to next attempt
		}
	}

	return lastErr
}
