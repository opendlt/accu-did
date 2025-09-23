package accdid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Sentinel errors for common conditions
var (
	ErrNotFound        = errors.New("not found")
	ErrGoneDeactivated = errors.New("gone: DID has been deactivated")
	ErrBadRequest      = errors.New("bad request")
	ErrServer          = errors.New("server error")
	ErrTimeout         = errors.New("timeout")
	ErrNetwork         = errors.New("network error")
)

// HTTPError represents a structured HTTP error
type HTTPError struct {
	StatusCode int
	Status     string
	Envelope   *ErrorEnvelope
	Body       []byte
	Err        error
}

func (e *HTTPError) Error() string {
	if e.Envelope != nil && e.Envelope.Message != "" {
		return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Envelope.Message)
	}
	if e.Envelope != nil && e.Envelope.Error != "" {
		return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Envelope.Error)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Status)
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

// decodeHTTPError converts an HTTP response to an appropriate error
func decodeHTTPError(resp *http.Response, body []byte) error {
	httpErr := &HTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
	}

	// Try to decode error envelope from response body
	if len(body) > 0 {
		var envelope ErrorEnvelope
		if err := json.Unmarshal(body, &envelope); err == nil {
			httpErr.Envelope = &envelope
		}
	}

	// Map to sentinel errors
	switch resp.StatusCode {
	case 404:
		httpErr.Err = ErrNotFound
	case 410:
		httpErr.Err = ErrGoneDeactivated
	case 400, 401, 403, 422, 429:
		httpErr.Err = ErrBadRequest
	default:
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			httpErr.Err = ErrBadRequest
		} else if resp.StatusCode >= 500 {
			httpErr.Err = ErrServer
		}
	}

	return httpErr
}

// classifyError determines the type of error and whether it's retryable
func classifyError(err error) (retryable bool, classified error) {
	if err == nil {
		return false, nil
	}

	// Check for timeout errors
	if isTimeoutError(err) {
		return false, fmt.Errorf("%w: %v", ErrTimeout, err)
	}

	// Check for network errors
	if isNetworkError(err) {
		return true, fmt.Errorf("%w: %v", ErrNetwork, err)
	}

	// Check for HTTP errors
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		// Retryable status codes
		switch httpErr.StatusCode {
		case 429, 502, 503, 504:
			return true, httpErr
		default:
			return false, httpErr
		}
	}

	// Unknown error type
	return false, err
}

// isTimeoutError checks if an error is a timeout
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	// Check for context deadline exceeded
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	// Check for net.Error timeout
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	// Check error message for timeout indicators
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "deadline exceeded") ||
		strings.Contains(errStr, "context canceled")
}

// isNetworkError checks if an error is a network connectivity issue
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common network error types
	var netErr net.Error
	if errors.As(err, &netErr) && !netErr.Timeout() {
		return true
	}

	// Check for DNS errors
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}

	// Check for connection errors
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "network unreachable") ||
		strings.Contains(errStr, "host unreachable")
}