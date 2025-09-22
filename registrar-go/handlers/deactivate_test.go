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
	"github.com/opendlt/accu-did/registrar-go/internal/api"
	"github.com/opendlt/accu-did/registrar-go/internal/policy"
)

func TestDeactivateHandler_Deactivate(t *testing.T) {
	// Setup
	accClient := acc.NewMockClient()
	authPolicy := policy.NewPolicyV1()
	handler := NewDeactivateHandler(accClient, authPolicy)

	// Test valid deactivate request
	t.Run("valid deactivate request", func(t *testing.T) {
		request := api.DeactivateRequest{
			DID: "did:acc:alice",
		}

		requestBody, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/deactivate", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Deactivate(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response DeactivateResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.JobID)
		assert.Equal(t, "did:acc:alice", response.DIDState.DID)
		assert.Equal(t, "finished", response.DIDState.State)
		assert.Equal(t, "deactivate", response.DIDState.Action)
		assert.NotEmpty(t, response.DIDRegistrationMetadata.VersionID)
		assert.Contains(t, response.DIDRegistrationMetadata.VersionID, "-deactivated")
		assert.NotEmpty(t, response.DIDRegistrationMetadata.ContentHash)
		assert.NotEmpty(t, response.DIDRegistrationMetadata.TxID)

		// Verify canonical tombstone structure was created
		mockClient := accClient.(*acc.MockClient)
		require.NotNil(t, mockClient.LastWriteData)

		var tombstone map[string]interface{}
		err = json.Unmarshal(mockClient.LastWriteData, &tombstone)
		require.NoError(t, err)

		// Check canonical tombstone fields
		assert.Equal(t, []interface{}{"https://www.w3.org/ns/did/v1"}, tombstone["@context"])
		assert.Equal(t, "did:acc:alice", tombstone["id"])
		assert.Equal(t, true, tombstone["deactivated"])
		assert.Contains(t, tombstone, "deactivatedAt")
	})

	// Test invalid requests
	t.Run("missing DID", func(t *testing.T) {
		request := api.DeactivateRequest{}

		requestBody, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/deactivate", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Deactivate(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response api.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "invalidRequest", response.Error)
		assert.Contains(t, response.Message, "DID is required")
	})

	t.Run("invalid DID format", func(t *testing.T) {
		request := api.DeactivateRequest{
			DID: "did:key:invalid",
		}

		requestBody, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/deactivate", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Deactivate(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response api.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "invalidRequest", response.Error)
		assert.Contains(t, response.Message, "invalid DID")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/deactivate", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Deactivate(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response api.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "invalidRequest", response.Error)
		assert.Contains(t, response.Message, "Invalid JSON")
	})
}
