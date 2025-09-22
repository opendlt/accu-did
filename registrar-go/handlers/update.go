package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/opendlt/accu-did/registrar-go/internal/acc"
	"github.com/opendlt/accu-did/registrar-go/internal/api"
	"github.com/opendlt/accu-did/registrar-go/internal/policy"
)

// UpdateHandler handles DID update requests
type UpdateHandler struct {
	accClient  acc.Submitter
	authPolicy policy.AuthPolicy
}

// LegacyUpdateRequest represents a DID update request for legacy compatibility
type LegacyUpdateRequest struct {
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
func NewUpdateHandler(accClient acc.Submitter, authPolicy policy.AuthPolicy) *UpdateHandler {
	return &UpdateHandler{
		accClient:  accClient,
		authPolicy: authPolicy,
	}
}

// Update handles POST /update requests
func (h *UpdateHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req LegacyUpdateRequest
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
	_, err := h.authPolicy.GetRequiredKeyPage(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Optional version bump: add versionId if missing
	didDoc := req.DIDDocument
	if didDoc == nil {
		didDoc = make(map[string]interface{})
	}

	// Add versionId if not present
	if _, hasVersion := didDoc["versionId"]; !hasVersion {
		didDoc["versionId"] = fmt.Sprintf("%d-update", time.Now().Unix())
	}

	// Ensure id field matches the DID
	didDoc["id"] = req.DID

	// Get data account URL using safe helper
	dataAccountURL, err := policy.DIDToDataAccountURL(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Convert DID document to JSON for writing
	didDocData, err := json.Marshal(didDoc)
	if err != nil {
		h.writeError(w, "internalError", "Failed to marshal DID document", http.StatusInternalServerError, nil)
		return
	}

	// Submit to Accumulate
	txID, err := h.accClient.WriteDataEntry(dataAccountURL, didDocData)
	if err != nil {
		h.writeError(w, "internalError", "Failed to submit update transaction", http.StatusInternalServerError, nil)
		return
	}

	// Generate version ID for this update
	versionID := fmt.Sprintf("%d-update", time.Now().Unix())

	// Build response
	response := UpdateResponse{
		JobID: h.generateJobID(),
		DIDState: DIDState{
			DID:    req.DID,
			State:  "finished",
			Action: "update",
		},
		DIDRegistrationMetadata: DIDRegistrationMetadata{
			VersionID:   versionID,
			ContentHash: h.calculateContentHash(didDocData),
			TxID:        txID,
		},
		DIDDocumentMetadata: DIDDocumentMetadata{
			Created:   time.Now().UTC(), // In real implementation, this would be the original creation time
			VersionID: versionID,
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
func (h *UpdateHandler) validateUpdateRequest(req *LegacyUpdateRequest) error {
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

// generateJobID generates a job ID for tracking the operation
func (h *UpdateHandler) generateJobID() string {
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}

// calculateContentHash calculates the SHA256 hash of content
func (h *UpdateHandler) calculateContentHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// writeError writes an error response
func (h *UpdateHandler) writeError(w http.ResponseWriter, errorCode, message string, status int, details map[string]string) {
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
