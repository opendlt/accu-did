package accdid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/opendlt/accu-did/sdks/go/accdid/retry"
)

func TestRetryLogic(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := atomic.AddInt32(&attempts, 1)

		// Fail first two attempts with 503, succeed on third
		if attempt <= 2 {
			w.WriteHeader(503)
			w.Write([]byte(`{"error":"Service temporarily unavailable"}`))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()

	client, err := NewResolverClient(ClientOptions{
		BaseURL: server.URL,
		Retries: RetryPolicy{
			Max:       3,
			BaseDelay: 10 * time.Millisecond, // Short delay for testing
			MaxDelay:  100 * time.Millisecond,
			Jitter:    false, // Disable jitter for predictable timing
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	start := time.Now()
	err = client.Health(context.Background())
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Expected success after retries, got error: %v", err)
	}

	// Should have made 3 attempts (2 failures + 1 success)
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts, got %d", atomic.LoadInt32(&attempts))
	}

	// Should have taken at least 2 retry delays (10ms + 20ms = 30ms minimum)
	minExpectedDuration := 30 * time.Millisecond
	if duration < minExpectedDuration {
		t.Errorf("Expected duration >= %v, got %v", minExpectedDuration, duration)
	}
}

func TestRetryGivesUpAfterMaxAttempts(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(503) // Always fail
		w.Write([]byte(`{"error":"Service unavailable"}`))
	}))
	defer server.Close()

	client, err := NewResolverClient(ClientOptions{
		BaseURL: server.URL,
		Retries: RetryPolicy{
			Max:       2, // Only 2 retries
			BaseDelay: 5 * time.Millisecond,
			MaxDelay:  50 * time.Millisecond,
			Jitter:    false,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Health(context.Background())
	if err == nil {
		t.Fatal("Expected error after max retries, got nil")
	}

	// Should have made max attempts + 1 (initial + 2 retries = 3 total)
	expectedAttempts := int32(3)
	if atomic.LoadInt32(&attempts) != expectedAttempts {
		t.Errorf("Expected %d attempts, got %d", expectedAttempts, atomic.LoadInt32(&attempts))
	}
}

func TestNoRetryOnNonRetryableStatus(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(404) // Not retryable
		w.Write([]byte(`{"error":"Not found"}`))
	}))
	defer server.Close()

	client, err := NewResolverClient(ClientOptions{
		BaseURL: server.URL,
		Retries: RetryPolicy{
			Max:       3,
			BaseDelay: 10 * time.Millisecond,
			MaxDelay:  100 * time.Millisecond,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Health(context.Background())
	if err == nil {
		t.Fatal("Expected error for 404, got nil")
	}

	// Should have made only 1 attempt (404 is not retryable)
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("Expected 1 attempt, got %d", atomic.LoadInt32(&attempts))
	}
}

func TestRetryableStatusCodes(t *testing.T) {
	retryableCodes := []int{429, 502, 503, 504}

	for _, code := range retryableCodes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			var attempts int32

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempt := atomic.AddInt32(&attempts, 1)
				if attempt == 1 {
					w.WriteHeader(code)
					w.Write([]byte(`{"error":"Retryable error"}`))
					return
				}
				w.WriteHeader(200)
				w.Write([]byte(`{"result":"success"}`))
			}))
			defer server.Close()

			client, err := NewResolverClient(ClientOptions{
				BaseURL: server.URL,
				Retries: RetryPolicy{
					Max:       2,
					BaseDelay: 5 * time.Millisecond,
					MaxDelay:  50 * time.Millisecond,
				},
			})
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			err = client.Health(context.Background())
			if err != nil {
				t.Fatalf("Expected success after retry, got error: %v", err)
			}

			// Should have retried once
			if atomic.LoadInt32(&attempts) != 2 {
				t.Errorf("Expected 2 attempts, got %d", atomic.LoadInt32(&attempts))
			}
		})
	}
}

func TestExponentialBackoff(t *testing.T) {
	policy := retry.Policy{
		Max:       3,
		BaseDelay: 10 * time.Millisecond,
		MaxDelay:  1 * time.Second,
		Jitter:    false,
	}

	// Create a mock doer to capture timing
	var delays []time.Duration
	var lastCall time.Time
	var callCount int

	mockDoer := func(req *http.Request) (*http.Response, error) {
		now := time.Now()
		callCount++

		if callCount > 1 { // Skip first call timing
			delays = append(delays, now.Sub(lastCall))
		}
		lastCall = now

		// Always return retryable error
		resp := &http.Response{
			StatusCode: 503,
			Body:       http.NoBody,
		}
		return resp, nil
	}

	doer := retry.WithRetry(doerFunc(mockDoer), policy, retry.DefaultIsRetryable)

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	doer.Do(req)

	// Should have 3 delays (between 4 attempts)
	if len(delays) != 3 {
		t.Fatalf("Expected 3 delays, got %d", len(delays))
	}

	// Verify exponential growth (within tolerance for timing variations)
	expectedDelays := []time.Duration{
		10 * time.Millisecond,  // First retry
		20 * time.Millisecond,  // Second retry
		40 * time.Millisecond,  // Third retry
	}

	for i, expected := range expectedDelays {
		if delays[i] < expected {
			t.Errorf("Delay %d too short: expected >= %v, got %v", i, expected, delays[i])
		}
		// Allow some tolerance for timing variations
		maxExpected := expected + 50*time.Millisecond
		if delays[i] > maxExpected {
			t.Errorf("Delay %d too long: expected <= %v, got %v", i, maxExpected, delays[i])
		}
	}
}