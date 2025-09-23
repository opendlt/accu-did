package accdid

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestResolverClient_Resolve(t *testing.T) {
	tests := []struct {
		name           string
		did            string
		responseFile   string
		statusCode     int
		expectedErr    error
		expectMetadata bool
	}{
		{
			name:           "successful resolution",
			did:            "did:acc:alice",
			responseFile:   "testdata/resolver_200.json",
			statusCode:     200,
			expectedErr:    nil,
			expectMetadata: true,
		},
		{
			name:         "not found",
			did:          "did:acc:nonexistent",
			responseFile: "testdata/resolver_404.json",
			statusCode:   404,
			expectedErr:  ErrNotFound,
		},
		{
			name:           "deactivated DID",
			did:            "did:acc:deactivated",
			responseFile:   "testdata/resolver_410.json",
			statusCode:     410,
			expectedErr:    ErrGoneDeactivated,
			expectMetadata: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load test response
			responseData, err := os.ReadFile(tt.responseFile)
			if err != nil {
				t.Fatalf("Failed to load test data: %v", err)
			}

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.HasPrefix(r.URL.Path, "/resolve") {
					t.Errorf("Expected path to start with /resolve, got %s", r.URL.Path)
				}

				if r.URL.Query().Get("did") != tt.did {
					t.Errorf("Expected DID %s, got %s", tt.did, r.URL.Query().Get("did"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write(responseData)
			}))
			defer server.Close()

			// Create client
			client, err := NewResolverClient(ClientOptions{
				BaseURL: server.URL,
			})
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Test resolution
			result, err := client.Resolve(context.Background(), tt.did)

			if tt.expectedErr != nil {
				if err == nil {
					t.Fatalf("Expected error %v, got nil", tt.expectedErr)
				}
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result, got nil")
			}

			if tt.expectMetadata && result.Metadata == nil {
				t.Error("Expected metadata, got nil")
			}

			// Verify DID document structure
			if result.DIDDocument == nil {
				t.Error("Expected DID document, got nil")
			}
		})
	}
}

func TestResolverClient_UniversalResolve(t *testing.T) {
	responseData := `{
		"didDocument": {
			"@context": ["https://www.w3.org/ns/did/v1"],
			"id": "did:acc:alice"
		},
		"didResolutionMetadata": {
			"contentType": "application/did+ld+json"
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/1.0/identifiers/did%3Aacc%3Aalice"
		if r.URL.Path != expectedPath {
			t.Logf("Path encoding test: expected %s, got %s", expectedPath, r.URL.Path)
			// Accept both encoded and non-encoded for flexibility
			if r.URL.Path != "/1.0/identifiers/did:acc:alice" {
				t.Errorf("Expected path %s or /1.0/identifiers/did:acc:alice, got %s", expectedPath, r.URL.Path)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(responseData))
	}))
	defer server.Close()

	client, err := NewResolverClient(ClientOptions{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.UniversalResolve(context.Background(), "did:acc:alice")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.DIDDocument == nil {
		t.Error("Expected DID document, got nil")
	}
}

func TestResolverClient_InvalidDID(t *testing.T) {
	client, err := NewResolverClient(ClientOptions{
		BaseURL: "http://localhost:8080",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []string{
		"",
		"invalid",
		"did:invalid:test",
		"not:a:did",
	}

	for _, invalidDID := range tests {
		t.Run("invalid_"+invalidDID, func(t *testing.T) {
			_, err := client.Resolve(context.Background(), invalidDID)
			if err == nil {
				t.Error("Expected error for invalid DID, got nil")
			}
			if !errors.Is(err, ErrInvalidDID) {
				t.Errorf("Expected ErrInvalidDID, got %v", err)
			}
		})
	}
}

func TestResolverClient_Health(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Errorf("Expected path /healthz, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client, err := NewResolverClient(ClientOptions{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Health(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestResolverClient_ErrorMapping(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		expectedErr error
	}{
		{"bad request", 400, ErrBadRequest},
		{"unauthorized", 401, ErrBadRequest},
		{"forbidden", 403, ErrBadRequest},
		{"not found", 404, ErrNotFound},
		{"gone", 410, ErrGoneDeactivated},
		{"unprocessable", 422, ErrBadRequest},
		{"too many requests", 429, ErrBadRequest},
		{"internal server error", 500, ErrServer},
		{"bad gateway", 502, ErrServer},
		{"service unavailable", 503, ErrServer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(`{"code":"error","error":"Error","message":"Test error"}`))
			}))
			defer server.Close()

			client, err := NewResolverClient(ClientOptions{
				BaseURL: server.URL,
			})
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			_, err = client.Resolve(context.Background(), "did:acc:test")
			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}