package accdid

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRegistrarClient_Register(t *testing.T) {
	responseData, err := os.ReadFile("testdata/registrar_register_200.json")
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/register" {
			t.Errorf("Expected path /register, got %s", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		// Verify request body structure
		var req NativeRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		if req.DID != "did:acc:test" {
			t.Errorf("Expected DID 'did:acc:test', got %s", req.DID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(responseData)
	}))
	defer server.Close()

	client, err := NewRegistrarClient(ClientOptions{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	didDoc := json.RawMessage(`{
		"@context": ["https://www.w3.org/ns/did/v1"],
		"id": "did:acc:test"
	}`)

	req := NativeRegisterRequest{
		DID:         "did:acc:test",
		DIDDocument: didDoc,
	}

	txID, err := client.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTxID := "tx-abc123def456"
	if txID != expectedTxID {
		t.Errorf("Expected transaction ID %s, got %s", expectedTxID, txID)
	}
}

func TestRegistrarClient_Update(t *testing.T) {
	responseData, err := os.ReadFile("testdata/registrar_update_200.json")
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/native/update" {
			t.Errorf("Expected path /native/update, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(responseData)
	}))
	defer server.Close()

	client, err := NewRegistrarClient(ClientOptions{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	patch := json.RawMessage(`{
		"addService": {
			"id": "did:acc:test#website",
			"type": "LinkedDomains",
			"serviceEndpoint": "https://example.com"
		}
	}`)

	req := NativeUpdateRequest{
		DID:   "did:acc:test",
		Patch: patch,
	}

	txID, err := client.Update(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTxID := "tx-update456789"
	if txID != expectedTxID {
		t.Errorf("Expected transaction ID %s, got %s", expectedTxID, txID)
	}
}

func TestRegistrarClient_Deactivate(t *testing.T) {
	responseData, err := os.ReadFile("testdata/registrar_deactivate_200.json")
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/native/deactivate" {
			t.Errorf("Expected path /native/deactivate, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(responseData)
	}))
	defer server.Close()

	client, err := NewRegistrarClient(ClientOptions{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := NativeDeactivateRequest{
		DID:    "did:acc:test",
		Reason: "Test deactivation",
	}

	txID, err := client.Deactivate(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTxID := "tx-deactivate789xyz"
	if txID != expectedTxID {
		t.Errorf("Expected transaction ID %s, got %s", expectedTxID, txID)
	}
}

func TestRegistrarClient_IdempotencyHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for idempotency key header
		idempotencyKey := r.Header.Get("Idempotency-Key")
		if idempotencyKey != "test-key-123" {
			t.Errorf("Expected Idempotency-Key header 'test-key-123', got %s", idempotencyKey)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"transactionId":"tx-123"}`))
	}))
	defer server.Close()

	client, err := NewRegistrarClient(ClientOptions{
		BaseURL:        server.URL,
		IdempotencyKey: "test-key-123",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	didDoc := json.RawMessage(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:test"}`)
	req := NativeRegisterRequest{
		DID:         "did:acc:test",
		DIDDocument: didDoc,
	}

	_, err = client.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestRegistrarClient_InvalidDID(t *testing.T) {
	client, err := NewRegistrarClient(ClientOptions{
		BaseURL: "http://localhost:8081",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name string
		did  string
	}{
		{"empty", ""},
		{"invalid format", "invalid"},
		{"wrong method", "did:invalid:test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NativeRegisterRequest{
				DID:         tt.did,
				DIDDocument: json.RawMessage(`{}`),
			}

			_, err := client.Register(context.Background(), req)
			if err == nil {
				t.Error("Expected error for invalid DID, got nil")
			}
			if !errors.Is(err, ErrInvalidDID) {
				t.Errorf("Expected ErrInvalidDID, got %v", err)
			}
		})
	}
}

func TestRegistrarClient_UniversalMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var expectedPath string
		switch {
		case strings.HasPrefix(r.URL.Path, "/1.0/create"):
			expectedPath = "/1.0/create"
		case strings.HasPrefix(r.URL.Path, "/1.0/update"):
			expectedPath = "/1.0/update"
		case strings.HasPrefix(r.URL.Path, "/1.0/deactivate"):
			expectedPath = "/1.0/deactivate"
		default:
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"transactionId":"tx-universal-123"}`))
	}))
	defer server.Close()

	client, err := NewRegistrarClient(ClientOptions{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test Universal Create
	didDoc := map[string]interface{}{
		"@context": []string{"https://www.w3.org/ns/did/v1"},
		"id":       "did:acc:universal",
	}

	txID, err := client.UniversalCreate(context.Background(), didDoc)
	if err != nil {
		t.Fatalf("UniversalCreate failed: %v", err)
	}
	if txID != "tx-universal-123" {
		t.Errorf("Expected tx-universal-123, got %s", txID)
	}

	// Test Universal Update
	patch := map[string]interface{}{
		"addService": map[string]interface{}{
			"id":              "did:acc:universal#service1",
			"type":            "LinkedDomains",
			"serviceEndpoint": "https://example.com",
		},
	}

	txID, err = client.UniversalUpdate(context.Background(), "did:acc:universal", patch)
	if err != nil {
		t.Fatalf("UniversalUpdate failed: %v", err)
	}
	if txID != "tx-universal-123" {
		t.Errorf("Expected tx-universal-123, got %s", txID)
	}

	// Test Universal Deactivate
	txID, err = client.UniversalDeactivate(context.Background(), "did:acc:universal")
	if err != nil {
		t.Fatalf("UniversalDeactivate failed: %v", err)
	}
	if txID != "tx-universal-123" {
		t.Errorf("Expected tx-universal-123, got %s", txID)
	}
}