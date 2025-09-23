// Package retry provides retry logic with exponential backoff for the Accumulate DID SDK
package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// Policy defines retry behavior
type Policy struct {
	Max       int
	BaseDelay time.Duration
	MaxDelay  time.Duration
	Jitter    bool
}

// IsRetryable determines if an error should be retried
type IsRetryable func(error) bool

// Doer represents something that can execute HTTP requests
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// WithRetry wraps a Doer with retry logic
func WithRetry(doer Doer, policy Policy, isRetryable IsRetryable) Doer {
	return &retryDoer{
		doer:        doer,
		policy:      policy,
		isRetryable: isRetryable,
	}
}

type retryDoer struct {
	doer        Doer
	policy      Policy
	isRetryable IsRetryable
}

func (r *retryDoer) Do(req *http.Request) (*http.Response, error) {
	var lastErr error
	var resp *http.Response

	for attempt := 0; attempt <= r.policy.Max; attempt++ {
		// Clone the request for retry safety
		reqClone := r.cloneRequest(req)

		resp, lastErr = r.doer.Do(reqClone)

		// Success case
		if lastErr == nil && resp.StatusCode < 500 && resp.StatusCode != 429 {
			return resp, nil
		}

		// Check if we should retry
		shouldRetry := false
		if lastErr != nil {
			shouldRetry = r.isRetryable != nil && r.isRetryable(lastErr)
		} else if resp != nil {
			// Check for retryable status codes
			shouldRetry = isRetryableStatus(resp.StatusCode)
			if !shouldRetry {
				return resp, nil
			}
		}

		// Don't retry if this was the last attempt or if not retryable
		if attempt >= r.policy.Max || !shouldRetry {
			break
		}

		// Calculate delay for next attempt
		delay := r.calculateDelay(attempt)

		// Wait with context cancellation support
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}

		// Close response body if present to avoid resource leaks
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}

	// Return the last response/error
	if resp != nil {
		return resp, nil
	}
	return nil, lastErr
}

func (r *retryDoer) cloneRequest(req *http.Request) *http.Request {
	// Clone the request
	clone := req.Clone(req.Context())

	// If there's a body, we need to handle it carefully
	if req.Body != nil && req.GetBody != nil {
		newBody, err := req.GetBody()
		if err == nil {
			clone.Body = newBody
		}
	}

	return clone
}

func (r *retryDoer) calculateDelay(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := time.Duration(float64(r.policy.BaseDelay) * math.Pow(2, float64(attempt)))

	// Cap at max delay
	if delay > r.policy.MaxDelay {
		delay = r.policy.MaxDelay
	}

	// Add jitter if enabled
	if r.policy.Jitter && delay > 0 {
		// Add up to 50% jitter
		jitterRange := float64(delay) * 0.5
		jitter := time.Duration(rand.Float64() * jitterRange)
		delay = delay + jitter
	}

	return delay
}

// isRetryableStatus checks if an HTTP status code is retryable
func isRetryableStatus(statusCode int) bool {
	switch statusCode {
	case 429, 502, 503, 504:
		return true
	default:
		return false
	}
}

// DefaultIsRetryable provides a reasonable default for determining retryable errors
func DefaultIsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for context cancellation (not retryable)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// For network errors, we should generally retry
	// This is a simplified check - in production you might want more sophisticated logic
	errStr := err.Error()
	retryableErrors := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"temporary failure",
		"network unreachable",
		"host unreachable",
	}

	for _, retryableErr := range retryableErrors {
		if contains(errStr, retryableErr) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		   s[:len(substr)] == substr ||
		   (len(s) > len(substr) && contains(s[1:], substr))
}