package accdid

import (
	"context"
	"net/http"
	"time"
)

// Doer is the interface for HTTP request execution
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// doerFunc allows a function to implement the Doer interface
type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Middleware modifies HTTP requests or responses
type Middleware func(Doer) Doer

// WithMiddleware wraps a Doer with a chain of middleware
func WithMiddleware(d Doer, middlewares ...Middleware) Doer {
	for i := len(middlewares) - 1; i >= 0; i-- {
		d = middlewares[i](d)
	}
	return d
}

// headerMiddleware adds common headers to requests
func headerMiddleware(apiKey, idempotencyKey string, requestIDGen func() string) Middleware {
	return func(next Doer) Doer {
		return doerFunc(func(req *http.Request) (*http.Response, error) {
			// Always set request ID
			if requestIDGen != nil {
				req.Header.Set("X-Request-Id", requestIDGen())
			}

			// Set API key if provided
			if apiKey != "" {
				req.Header.Set("X-API-Key", apiKey)
			}

			// Set idempotency key if provided
			if idempotencyKey != "" {
				req.Header.Set("Idempotency-Key", idempotencyKey)
			}

			// Ensure JSON content type and accept headers
			if req.Body != nil && req.Header.Get("Content-Type") == "" {
				req.Header.Set("Content-Type", "application/json")
			}
			if req.Header.Get("Accept") == "" {
				req.Header.Set("Accept", "application/json")
			}

			return next.Do(req)
		})
	}
}

// timeoutMiddleware adds a timeout context to requests
func timeoutMiddleware(timeout time.Duration) Middleware {
	return func(next Doer) Doer {
		return doerFunc(func(req *http.Request) (*http.Response, error) {
			if timeout <= 0 {
				return next.Do(req)
			}

			ctx, cancel := context.WithTimeout(req.Context(), timeout)
			defer cancel()

			return next.Do(req.WithContext(ctx))
		})
	}
}