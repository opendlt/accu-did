package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/opendlt/accu-did/registrar-go/internal/acc"
	"github.com/opendlt/accu-did/registrar-go/internal/api"
	"github.com/opendlt/accu-did/registrar-go/internal/ops"
	"github.com/opendlt/accu-did/registrar-go/internal/policy"
)

// CreateHandler handles DID creation requests
type CreateHandler struct {
	accClient  acc.Submitter
	authPolicy policy.AuthPolicy
}

// CreateRequest represents a DID creation request
type CreateRequest struct {
	DID         string                 `json:"did"`
	DIDDocument map[string]interface{} `json:"didDocument"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Secret      map[string]interface{} `json:"secret,omitempty"`
}

// CreateResponse represents a DID creation response
type CreateResponse struct {
	JobID                   string                  `json:"jobId"`
	DIDState                DIDState                `json:"didState"`
	DIDRegistrationMetadata DIDRegistrationMetadata `json:"didRegistrationMetadata"`
	DIDDocumentMetadata     DIDDocumentMetadata     `json:"didDocumentMetadata"`
}

// DIDState represents the state of a DID operation
type DIDState struct {
	DID    string `json:"did"`
	State  string `json:"state"`  // "finished", "failed", "action"
	Action string `json:"action"` // "create", "update", "deactivate"
	Reason string `json:"reason,omitempty"`
}

// DIDRegistrationMetadata represents registration-specific metadata
type DIDRegistrationMetadata struct {
	VersionID   string `json:"versionId"`
	ContentHash string `json:"contentHash"`
	TxID        string `json:"txid"`
}

// DIDDocumentMetadata represents document metadata
type DIDDocumentMetadata struct {
	Created   time.Time `json:"created"`
	VersionID string    `json:"versionId"`
}

// NewCreateHandler creates a new create handler
func NewCreateHandler(accClient acc.Submitter, authPolicy policy.AuthPolicy) *CreateHandler {
	return &CreateHandler{
		accClient:  accClient,
		authPolicy: authPolicy,
	}
}

// Create handles POST /create requests
func (h *CreateHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalidRequest", "Invalid JSON", http.StatusBadRequest, nil)
		return
	}

	// Validate request
	if err := h.validateCreateRequest(&req); err != nil {
		h.writeError(w, "invalidRequest", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Get required key page for authorization
	requiredKeyPage, err := h.authPolicy.GetRequiredKeyPage(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Build envelope
	envelope, err := ops.BuildEnvelope(req.DIDDocument, requiredKeyPage, "")
	if err != nil {
		h.writeError(w, "internalError", "Failed to build envelope", http.StatusInternalServerError, nil)
		return
	}

	// Get data account URL using safe helper
	dataAccountURL, err := policy.DIDToDataAccountURL(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Submit to Accumulate
	txID, err := h.accClient.SubmitWriteData(dataAccountURL, envelope)
	if err != nil {
		h.writeError(w, "internalError", "Failed to submit transaction", http.StatusInternalServerError, nil)
		return
	}

	// Build response
	response := CreateResponse{
		JobID: h.generateJobID(),
		DIDState: DIDState{
			DID:    req.DID,
			State:  "finished",
			Action: "create",
		},
		DIDRegistrationMetadata: DIDRegistrationMetadata{
			VersionID:   envelope.Meta.VersionID,
			ContentHash: envelope.GetContentHash(),
			TxID:        txID,
		},
		DIDDocumentMetadata: DIDDocumentMetadata{
			Created:   envelope.Meta.Timestamp,
			VersionID: envelope.Meta.VersionID,
		},
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.writeError(w, "internalError", "Failed to encode response", http.StatusInternalServerError, nil)
		return
	}
}

// validateCreateRequest validates the create request
func (h *CreateHandler) validateCreateRequest(req *CreateRequest) error {
	if req.DID == "" {
		return fmt.Errorf("DID is required")
	}

	if err := policy.ValidateDID(req.DID); err != nil {
		return fmt.Errorf("invalid DID: %w", err)
	}

	if req.DIDDocument == nil {
		return fmt.Errorf("didDocument is required")
	}

	// Validate that the DID in the document matches the request
	if docID, ok := req.DIDDocument["id"].(string); ok {
		if docID != req.DID {
			return fmt.Errorf("DID mismatch: request DID %s does not match document ID %s", req.DID, docID)
		}
	} else {
		return fmt.Errorf("didDocument must have an 'id' field")
	}

	// Validate required fields
	if _, ok := req.DIDDocument["@context"]; !ok {
		return fmt.Errorf("didDocument must have '@context' field")
	}

	return nil
}

// generateJobID generates a job ID for tracking the operation
func (h *CreateHandler) generateJobID() string {
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}

// writeError writes an error response
func (h *CreateHandler) writeError(w http.ResponseWriter, errorCode, message string, status int, details map[string]string) {
	response := api.ErrorResponse{
		Error:     errorCode,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(response)
}
