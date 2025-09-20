package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opendlt/accu-did/registrar-go/internal/acc"
	"github.com/opendlt/accu-did/registrar-go/internal/policy"
)

func TestCreateHandler_Create(t *testing.T) {
	// Setup
	accClient := acc.NewMockClient()
	authPolicy := policy.NewPolicyV1()
	handler := NewCreateHandler(accClient, authPolicy)

	// Test valid create request
	t.Run("valid create request", func(t *testing.T) {
		request := CreateRequest{
			DID: "did:acc:alice",
			DIDDocument: map[string]interface{}{
				"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
				"id":       "did:acc:alice",
				"verificationMethod": []interface{}{
					map[string]interface{}{
						"id":         "did:acc:alice#key-1",
						"type":       "AccumulateKeyPage",
						"controller": "did:acc:alice",
						"keyPageUrl": "acc://alice/book/1",
						"threshold":  1,
					},
				},
			},
		}

		requestBody, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/create", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response CreateResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.JobID)
		assert.Equal(t, "did:acc:alice", response.DIDState.DID)
		assert.Equal(t, "finished", response.DIDState.State)
		assert.Equal(t, "create", response.DIDState.Action)
		assert.NotEmpty(t, response.DIDRegistrationMetadata.VersionID)
		assert.NotEmpty(t, response.DIDRegistrationMetadata.ContentHash)
		assert.NotEmpty(t, response.DIDRegistrationMetadata.TxID)
	})

	// Test invalid requests
	t.Run("missing DID", func(t *testing.T) {
		request := CreateRequest{
			DIDDocument: map[string]interface{}{
				"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
				"id":       "did:acc:alice",
			},
		}

		requestBody, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/create", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "invalidRequest", response.Error)
		assert.Contains(t, response.Message, "DID is required")
	})

	t.Run("DID mismatch", func(t *testing.T) {
		request := CreateRequest{
			DID: "did:acc:alice",
			DIDDocument: map[string]interface{}{
				"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
				"id":       "did:acc:bob", // Different DID
			},
		}

		requestBody, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/create", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "invalidRequest", response.Error)
		assert.Contains(t, response.Message, "DID mismatch")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/create", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "invalidRequest", response.Error)
		assert.Contains(t, response.Message, "Invalid JSON")
	})
}