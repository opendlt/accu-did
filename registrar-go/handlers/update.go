package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/opendlt/accu-did/registrar-go/internal/acc"
	"github.com/opendlt/accu-did/registrar-go/internal/ops"
	"github.com/opendlt/accu-did/registrar-go/internal/policy"
)

// UpdateHandler handles DID update requests
type UpdateHandler struct {
	accClient  acc.Client
	authPolicy policy.AuthPolicy
}

// UpdateRequest represents a DID update request
type UpdateRequest struct {
	DID         string                 `json:"did"`
	DIDDocument map[string]interface{} `json:"didDocument"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Secret      map[string]interface{} `json:"secret,omitempty"`
}

// UpdateResponse represents a DID update response
type UpdateResponse struct {
	JobID                   string                  `json:"jobId"`
	DIDState                DIDState                `json:"didState"`
	DIDRegistrationMetadata DIDRegistrationMetadata `json:"didRegistrationMetadata"`
	DIDDocumentMetadata     DIDDocumentMetadata     `json:"didDocumentMetadata"`
}

// NewUpdateHandler creates a new update handler
func NewUpdateHandler(accClient acc.Client, authPolicy policy.AuthPolicy) *UpdateHandler {
	return &UpdateHandler{
		accClient:  accClient,
		authPolicy: authPolicy,
	}
}

// Update handles POST /update requests
func (h *UpdateHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalidRequest", "Invalid JSON", http.StatusBadRequest, nil)
		return
	}

	// Validate request
	if err := h.validateUpdateRequest(&req); err != nil {
		h.writeError(w, "invalidRequest", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Get required key page for authorization
	requiredKeyPage, err := h.authPolicy.GetRequiredKeyPage(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// For updates, we should get the previous version ID
	// In a real implementation, this would query the current DID document
	previousVersionID := h.getPreviousVersionID(req.DID)

	// Build envelope
	envelope, err := ops.BuildEnvelope(req.DIDDocument, requiredKeyPage, previousVersionID)
	if err != nil {
		h.writeError(w, "internalError", "Failed to build envelope", http.StatusInternalServerError, nil)
		return
	}

	// Submit to Accumulate
	dataAccountURL := h.getDataAccountURL(req.DID)
	txID, err := h.accClient.SubmitWriteData(dataAccountURL, envelope)
	if err != nil {
		h.writeError(w, "internalError", "Failed to submit transaction", http.StatusInternalServerError, nil)
		return
	}

	// Build response
	response := UpdateResponse{
		JobID: h.generateJobID(),
		DIDState: DIDState{
			DID:    req.DID,
			State:  "finished",
			Action: "update",
		},
		DIDRegistrationMetadata: DIDRegistrationMetadata{
			VersionID:   envelope.Meta.VersionID,
			ContentHash: envelope.GetContentHash(),
			TxID:        txID,
		},
		DIDDocumentMetadata: DIDDocumentMetadata{
			Created:   envelope.Meta.Timestamp, // In real implementation, this would be the original creation time
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

// validateUpdateRequest validates the update request
func (h *UpdateHandler) validateUpdateRequest(req *UpdateRequest) error {
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

// getPreviousVersionID gets the previous version ID for an update
// In a real implementation, this would query the current DID document from Accumulate
func (h *UpdateHandler) getPreviousVersionID(did string) string {
	// Mock implementation - return a placeholder
	return fmt.Sprintf("%d-%s", time.Now().Unix()-3600, "previous")
}

// getDataAccountURL constructs the data account URL for a DID
func (h *UpdateHandler) getDataAccountURL(did string) string {
	// Extract ADI from DID (simplified)
	adi := did[8:] // Remove "did:acc:" prefix

	// Handle URL components
	for _, separator := range []string{"/", "?", "#", ";"} {
		if idx := len(adi); idx > 0 {
			for i, char := range adi {
				if string(char) == separator {
					adi = adi[:i]
					break
				}
			}
		}
	}

	return fmt.Sprintf("acc://%s/data/did", adi)
}

// generateJobID generates a job ID for tracking the operation
func (h *UpdateHandler) generateJobID() string {
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}

// writeError writes an error response
func (h *UpdateHandler) writeError(w http.ResponseWriter, errorCode, message string, status int, details map[string]string) {
	response := ErrorResponse{
		Error:     errorCode,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(response)
}