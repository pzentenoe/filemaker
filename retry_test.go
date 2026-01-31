package filemaker

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("MaxRetries = %v, want 3", config.MaxRetries)
	}

	if config.MinWaitTime != 1*time.Second {
		t.Errorf("MinWaitTime = %v, want 1s", config.MinWaitTime)
	}

	if config.MaxWaitTime != 30*time.Second {
		t.Errorf("MaxWaitTime = %v, want 30s", config.MaxWaitTime)
	}

	if !config.EnableJitter {
		t.Error("EnableJitter should be true")
	}

	expectedCodes := []int{429, 500, 502, 503, 504}
	if len(config.RetryableStatusCodes) != len(expectedCodes) {
		t.Errorf("RetryableStatusCodes length = %v, want %v", len(config.RetryableStatusCodes), len(expectedCodes))
	}
}

func TestRetryConfig_shouldRetry(t *testing.T) {
	config := DefaultRetryConfig()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "retryable filemaker error",
			err:  &FileMakerError{Code: "952", HTTPStatus: http.StatusServiceUnavailable},
			want: true,
		},
		{
			name: "non-retryable filemaker error",
			err:  &FileMakerError{Code: "401", HTTPStatus: http.StatusUnauthorized},
			want: false,
		},
		{
			name: "network error",
			err:  &NetworkError{Message: "connection refused"},
			want: true,
		},
		{
			name: "timeout error",
			err:  &TimeoutError{Message: "timeout"},
			want: true,
		},
		{
			name: "generic error",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.shouldRetry(tt.err); got != tt.want {
				t.Errorf("shouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryConfig_calculateBackoff(t *testing.T) {
	tests := []struct {
		name    string
		config  *RetryConfig
		attempt int
		minWant time.Duration
		maxWant time.Duration
	}{
		{
			name: "first attempt no jitter",
			config: &RetryConfig{
				MinWaitTime:  1 * time.Second,
				MaxWaitTime:  30 * time.Second,
				EnableJitter: false,
			},
			attempt: 0,
			minWant: 1 * time.Second,
			maxWant: 1 * time.Second,
		},
		{
			name: "second attempt no jitter",
			config: &RetryConfig{
				MinWaitTime:  1 * time.Second,
				MaxWaitTime:  30 * time.Second,
				EnableJitter: false,
			},
			attempt: 1,
			minWant: 2 * time.Second,
			maxWant: 2 * time.Second,
		},
		{
			name: "third attempt no jitter",
			config: &RetryConfig{
				MinWaitTime:  1 * time.Second,
				MaxWaitTime:  30 * time.Second,
				EnableJitter: false,
			},
			attempt: 2,
			minWant: 4 * time.Second,
			maxWant: 4 * time.Second,
		},
		{
			name: "max wait time exceeded",
			config: &RetryConfig{
				MinWaitTime:  1 * time.Second,
				MaxWaitTime:  5 * time.Second,
				EnableJitter: false,
			},
			attempt: 10,
			minWant: 5 * time.Second,
			maxWant: 5 * time.Second,
		},
		{
			name: "with jitter",
			config: &RetryConfig{
				MinWaitTime:  1 * time.Second,
				MaxWaitTime:  30 * time.Second,
				EnableJitter: true,
			},
			attempt: 1,
			minWant: 2 * time.Second,
			maxWant: 3 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backoff := tt.config.calculateBackoff(tt.attempt)

			if backoff < tt.minWant {
				t.Errorf("calculateBackoff() = %v, want >= %v", backoff, tt.minWant)
			}

			if backoff > tt.maxWant {
				t.Errorf("calculateBackoff() = %v, want <= %v", backoff, tt.maxWant)
			}
		})
	}
}

func TestRetryConfig_executeWithRetry_Success(t *testing.T) {
	config := &RetryConfig{
		MaxRetries:   3,
		MinWaitTime:  10 * time.Millisecond,
		MaxWaitTime:  100 * time.Millisecond,
		EnableJitter: false,
	}

	ctx := context.Background()
	callCount := 0

	err := config.executeWithRetry(ctx, func() error {
		callCount++
		return nil
	})

	if err != nil {
		t.Errorf("executeWithRetry() returned error: %v", err)
	}

	if callCount != 1 {
		t.Errorf("function called %d times, want 1", callCount)
	}
}

func TestRetryConfig_executeWithRetry_NonRetryableError(t *testing.T) {
	config := &RetryConfig{
		MaxRetries:   3,
		MinWaitTime:  10 * time.Millisecond,
		MaxWaitTime:  100 * time.Millisecond,
		EnableJitter: false,
	}

	ctx := context.Background()
	callCount := 0
	expectedErr := &ValidationError{Field: "test", Message: "error"}

	err := config.executeWithRetry(ctx, func() error {
		callCount++
		return expectedErr
	})

	if err != expectedErr {
		t.Errorf("executeWithRetry() error = %v, want %v", err, expectedErr)
	}

	if callCount != 1 {
		t.Errorf("function called %d times, want 1", callCount)
	}
}

func TestRetryConfig_executeWithRetry_RetryableError(t *testing.T) {
	config := &RetryConfig{
		MaxRetries:   3,
		MinWaitTime:  10 * time.Millisecond,
		MaxWaitTime:  100 * time.Millisecond,
		EnableJitter: false,
	}

	ctx := context.Background()
	callCount := 0
	retryableErr := &NetworkError{Message: "connection refused"}

	err := config.executeWithRetry(ctx, func() error {
		callCount++
		return retryableErr
	})

	if err != retryableErr {
		t.Errorf("executeWithRetry() error = %v, want %v", err, retryableErr)
	}

	if callCount != 4 {
		t.Errorf("function called %d times, want 4 (1 initial + 3 retries)", callCount)
	}
}

func TestRetryConfig_executeWithRetry_EventualSuccess(t *testing.T) {
	config := &RetryConfig{
		MaxRetries:   3,
		MinWaitTime:  10 * time.Millisecond,
		MaxWaitTime:  100 * time.Millisecond,
		EnableJitter: false,
	}

	ctx := context.Background()
	callCount := 0

	err := config.executeWithRetry(ctx, func() error {
		callCount++
		if callCount < 3 {
			return &NetworkError{Message: "connection refused"}
		}
		return nil
	})

	if err != nil {
		t.Errorf("executeWithRetry() returned error: %v", err)
	}

	if callCount != 3 {
		t.Errorf("function called %d times, want 3", callCount)
	}
}

func TestRetryConfig_executeWithRetry_ContextCancelled(t *testing.T) {
	config := &RetryConfig{
		MaxRetries:   3,
		MinWaitTime:  100 * time.Millisecond,
		MaxWaitTime:  1 * time.Second,
		EnableJitter: false,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	callCount := 0

	err := config.executeWithRetry(ctx, func() error {
		callCount++
		return &NetworkError{Message: "connection refused"}
	})

	if err != context.Canceled {
		t.Errorf("executeWithRetry() error = %v, want %v", err, context.Canceled)
	}

	if callCount > 1 {
		t.Errorf("function called %d times, want 0 or 1", callCount)
	}
}

func TestRetryConfig_executeWithRetry_ContextTimeout(t *testing.T) {
	config := &RetryConfig{
		MaxRetries:   5,
		MinWaitTime:  200 * time.Millisecond,
		MaxWaitTime:  1 * time.Second,
		EnableJitter: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	callCount := 0

	err := config.executeWithRetry(ctx, func() error {
		callCount++
		return &NetworkError{Message: "connection refused"}
	})

	if err != context.DeadlineExceeded {
		t.Errorf("executeWithRetry() error = %v, want %v", err, context.DeadlineExceeded)
	}
}

func TestRetryConfig_executeWithRetry_OnRetryCallback(t *testing.T) {
	retryAttempts := []int{}
	retryErrors := []error{}

	config := &RetryConfig{
		MaxRetries:   3,
		MinWaitTime:  10 * time.Millisecond,
		MaxWaitTime:  100 * time.Millisecond,
		EnableJitter: false,
		OnRetry: func(attempt int, err error) {
			retryAttempts = append(retryAttempts, attempt)
			retryErrors = append(retryErrors, err)
		},
	}

	ctx := context.Background()
	expectedErr := &NetworkError{Message: "test"}

	_ = config.executeWithRetry(ctx, func() error {
		return expectedErr
	})

	if len(retryAttempts) != 3 {
		t.Errorf("OnRetry called %d times, want 3", len(retryAttempts))
	}

	for i, attempt := range retryAttempts {
		expectedAttempt := i + 1
		if attempt != expectedAttempt {
			t.Errorf("attempt %d = %d, want %d", i, attempt, expectedAttempt)
		}
	}

	for i, err := range retryErrors {
		if err != expectedErr {
			t.Errorf("error %d = %v, want %v", i, err, expectedErr)
		}
	}
}

func TestRetryConfig_executeWithRetry_MaxRetriesZero(t *testing.T) {
	config := &RetryConfig{
		MaxRetries:   0,
		MinWaitTime:  10 * time.Millisecond,
		MaxWaitTime:  100 * time.Millisecond,
		EnableJitter: false,
	}

	ctx := context.Background()
	callCount := 0

	err := config.executeWithRetry(ctx, func() error {
		callCount++
		return &NetworkError{Message: "error"}
	})

	if err == nil {
		t.Error("executeWithRetry() should return error with MaxRetries=0")
	}

	if callCount != 1 {
		t.Errorf("function called %d times, want 1", callCount)
	}
}

func BenchmarkRetryConfig_calculateBackoff(b *testing.B) {
	config := DefaultRetryConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.calculateBackoff(i % 5)
	}
}

func BenchmarkRetryConfig_executeWithRetry_Success(b *testing.B) {
	config := &RetryConfig{
		MaxRetries:   3,
		MinWaitTime:  1 * time.Millisecond,
		MaxWaitTime:  10 * time.Millisecond,
		EnableJitter: false,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.executeWithRetry(ctx, func() error {
			return nil
		})
	}
}
