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

// DeactivateHandler handles DID deactivation requests
type DeactivateHandler struct {
	accClient  acc.Submitter
	authPolicy policy.AuthPolicy
}

// DeactivateResponse represents a DID deactivation response
type DeactivateResponse struct {
	JobID                   string                  `json:"jobId"`
	DIDState                DIDState                `json:"didState"`
	DIDRegistrationMetadata DIDRegistrationMetadata `json:"didRegistrationMetadata"`
	DIDDocumentMetadata     DIDDocumentMetadata     `json:"didDocumentMetadata"`
}

// NewDeactivateHandler creates a new deactivate handler
func NewDeactivateHandler(accClient acc.Submitter, authPolicy policy.AuthPolicy) *DeactivateHandler {
	return &DeactivateHandler{
		accClient:  accClient,
		authPolicy: authPolicy,
	}
}

// Deactivate handles POST /deactivate requests
func (h *DeactivateHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req api.DeactivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalidRequest", "Invalid JSON", http.StatusBadRequest, nil)
		return
	}

	// Validate request
	if err := h.validateDeactivateRequest(&req); err != nil {
		h.writeError(w, "invalidRequest", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Get required key page for authorization
	_, err := h.authPolicy.GetRequiredKeyPage(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Create canonical deactivation tombstone
	deactivationDoc := map[string]interface{}{
		"@context":      []string{"https://www.w3.org/ns/did/v1"},
		"id":            req.DID,
		"deactivated":   true,
		"deactivatedAt": time.Now().UTC().Format(time.RFC3339),
	}

	// Add optional reason if provided in the request
	if req.Options != nil {
		if reason, ok := req.Options["reason"].(string); ok && reason != "" {
			deactivationDoc["reason"] = reason
		}
	}

	// Convert to JSON for writing
	deactivationData, err := json.Marshal(deactivationDoc)
	if err != nil {
		h.writeError(w, "internalError", "Failed to marshal deactivation document", http.StatusInternalServerError, nil)
		return
	}

	// Get data account URL using safe helper
	dataAccountURL, err := policy.DIDToDataAccountURL(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Submit deactivation tombstone to Accumulate
	txID, err := h.accClient.WriteDataEntry(dataAccountURL, deactivationData)
	if err != nil {
		h.writeError(w, "internalError", "Failed to submit deactivation transaction", http.StatusInternalServerError, nil)
		return
	}

	// Generate version ID for this deactivation
	versionID := fmt.Sprintf("%d-deactivated", time.Now().Unix())

	// Build response
	response := DeactivateResponse{
		JobID: h.generateJobID(),
		DIDState: DIDState{
			DID:    req.DID,
			State:  "finished",
			Action: "deactivate",
		},
		DIDRegistrationMetadata: DIDRegistrationMetadata{
			VersionID:   versionID,
			ContentHash: h.calculateContentHash(deactivationData),
			TxID:        txID,
		},
		DIDDocumentMetadata: DIDDocumentMetadata{
			Created:   time.Now().UTC(), // Use current time for deactivation
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

// validateDeactivateRequest validates the deactivate request
func (h *DeactivateHandler) validateDeactivateRequest(req *api.DeactivateRequest) error {
	if req.DID == "" {
		return fmt.Errorf("DID is required")
	}

	if err := policy.ValidateDID(req.DID); err != nil {
		return fmt.Errorf("invalid DID: %w", err)
	}

	return nil
}

// getPreviousVersionID gets the previous version ID for deactivation
// In a real implementation, this would query the current DID document from Accumulate
func (h *DeactivateHandler) getPreviousVersionID(did string) string {
	// Mock implementation - return a placeholder
	return fmt.Sprintf("%d-%s", time.Now().Unix()-3600, "current")
}

// generateJobID generates a job ID for tracking the operation
func (h *DeactivateHandler) generateJobID() string {
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}

// calculateContentHash calculates the SHA256 hash of content
func (h *DeactivateHandler) calculateContentHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// writeError writes an error response
func (h *DeactivateHandler) writeError(w http.ResponseWriter, errorCode, message string, status int, details map[string]string) {
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
