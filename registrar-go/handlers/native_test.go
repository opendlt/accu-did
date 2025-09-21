package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opendlt/accu-did/registrar-go/internal/acc"
	"github.com/opendlt/accu-did/registrar-go/internal/api"
)

func TestNativeRegister(t *testing.T) {
	// Create handler with fake client
	client := acc.NewFakeSubmitter()
	handler := NewNativeHandler(client)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *NativeResponse)
	}{
		{
			name: "successful registration",
			requestBody: RegisterRequest{
				DID: "did:acc:testuser",
				DIDDocument: map[string]interface{}{
					"@context": []string{"https://www.w3.org/ns/did/v1"},
					"id":       "did:acc:testuser",
					"verificationMethod": []interface{}{
						map[string]interface{}{
							"id":                 "did:acc:testuser#key1",
							"type":               "Ed25519VerificationKey2020",
							"controller":         "did:acc:testuser",
							"publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK",
						},
					},
				},
				KeyPageURL: "acc://testuser/book/1",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *NativeResponse) {
				if !resp.Success {
					t.Error("expected success to be true")
				}
				if resp.DID != "did:acc:testuser" {
					t.Errorf("expected DID 'did:acc:testuser', got %s", resp.DID)
				}
				if resp.TxID == "" {
					t.Error("expected non-empty transaction ID")
				}
				if resp.JobID == "" {
					t.Error("expected non-empty job ID")
				}
			},
		},
		{
			name: "missing DID",
			requestBody: RegisterRequest{
				DIDDocument: map[string]interface{}{
					"@context": []string{"https://www.w3.org/ns/did/v1"},
					"id":       "did:acc:testuser",
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing DID document",
			requestBody: RegisterRequest{
				DID: "did:acc:testuser",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "DID mismatch",
			requestBody: RegisterRequest{
				DID: "did:acc:testuser",
				DIDDocument: map[string]interface{}{
					"@context": []string{"https://www.w3.org/ns/did/v1"},
					"id":       "did:acc:different",
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Create recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.Register(w, req)

			// Check status
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response for successful cases
			if tt.expectedStatus == http.StatusOK && tt.checkResponse != nil {
				var resp NativeResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				tt.checkResponse(t, &resp)
			}
		})
	}
}

func TestNativeUpdate(t *testing.T) {
	// Create handler with fake client
	client := acc.NewFakeSubmitter()
	handler := NewNativeHandler(client)

	requestBody := api.UpdateRequest{
		DID: "did:acc:testuser",
		Document: map[string]interface{}{
			"@context": []string{"https://www.w3.org/ns/did/v1"},
			"id":       "did:acc:testuser",
			"verificationMethod": []interface{}{
				map[string]interface{}{
					"id":                 "did:acc:testuser#key1",
					"type":               "Ed25519VerificationKey2020",
					"controller":         "did:acc:testuser",
					"publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK",
				},
				map[string]interface{}{
					"id":                 "did:acc:testuser#key2",
					"type":               "Ed25519VerificationKey2020",
					"controller":         "did:acc:testuser",
					"publicKeyMultibase": "z6MkgoLTnTypo3tDRwCkZXSccTPHRLhF4ZnjhueYAFpEX6vg",
				},
			},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest("POST", "/native/update", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp NativeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}

	if resp.DID != "did:acc:testuser" {
		t.Errorf("expected DID 'did:acc:testuser', got %s", resp.DID)
	}
}

func TestNativeDeactivate(t *testing.T) {
	// Create handler with fake client
	client := acc.NewFakeSubmitter()
	handler := NewNativeHandler(client)

	requestBody := api.DeactivateRequest{
		DID: "did:acc:testuser",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest("POST", "/native/deactivate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Deactivate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp NativeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}

	if resp.DID != "did:acc:testuser" {
		t.Errorf("expected DID 'did:acc:testuser', got %s", resp.DID)
	}
}
