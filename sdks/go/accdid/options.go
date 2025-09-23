package accdid

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

// Logger defines the logging interface for the SDK
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// noopLogger is a logger that discards all messages
type noopLogger struct{}

func (noopLogger) Debugf(format string, args ...interface{}) {}
func (noopLogger) Infof(format string, args ...interface{})  {}
func (noopLogger) Warnf(format string, args ...interface{})  {}
func (noopLogger) Errorf(format string, args ...interface{}) {}

// RetryPolicy defines the retry behavior for failed requests
type RetryPolicy struct {
	// Max is the maximum number of retry attempts
	Max int

	// BaseDelay is the initial delay between retries
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration

	// Jitter adds randomization to retry delays to avoid thundering herd
	Jitter bool

	// Backoff strategy: "exp" for exponential (default)
	Backoff string
}

// ClientOptions configures the SDK clients
type ClientOptions struct {
	// BaseURL is the base URL of the service (required)
	BaseURL string

	// APIKey for authentication (optional)
	APIKey string

	// Timeout for individual requests
	Timeout time.Duration

	// HTTP client to use (optional, defaults to http.DefaultClient)
	HTTP *http.Client

	// Retries configuration for failed requests
	Retries RetryPolicy

	// Logger for debug output (optional)
	Logger Logger

	// IdempotencyKey for safe retries (optional)
	IdempotencyKey string

	// RequestID generator function (optional)
	RequestID func() string
}

// defaultOptions returns ClientOptions with sensible defaults
func defaultOptions() ClientOptions {
	return ClientOptions{
		Timeout: 10 * time.Second,
		HTTP:    http.DefaultClient,
		Retries: RetryPolicy{
			Max:       3,
			BaseDelay: 250 * time.Millisecond,
			MaxDelay:  4 * time.Second,
			Jitter:    true,
			Backoff:   "exp",
		},
		Logger: noopLogger{},
		RequestID: defaultRequestID,
	}
}

// defaultRequestID generates a random request ID
func defaultRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// applyDefaults fills in default values for unset options
func applyDefaults(opts ClientOptions) ClientOptions {
	defaults := defaultOptions()

	if opts.Timeout == 0 {
		opts.Timeout = defaults.Timeout
	}
	if opts.HTTP == nil {
		opts.HTTP = defaults.HTTP
	}
	if opts.Retries.Max == 0 {
		opts.Retries = defaults.Retries
	}
	if opts.Retries.BaseDelay == 0 {
		opts.Retries.BaseDelay = defaults.Retries.BaseDelay
	}
	if opts.Retries.MaxDelay == 0 {
		opts.Retries.MaxDelay = defaults.Retries.MaxDelay
	}
	if opts.Retries.Backoff == "" {
		opts.Retries.Backoff = defaults.Retries.Backoff
	}
	if opts.Logger == nil {
		opts.Logger = defaults.Logger
	}
	if opts.RequestID == nil {
		opts.RequestID = defaults.RequestID
	}

	return opts
}