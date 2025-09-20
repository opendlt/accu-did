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

// DeactivateHandler handles DID deactivation requests
type DeactivateHandler struct {
	accClient  acc.Client
	authPolicy policy.AuthPolicy
}

// DeactivateRequest represents a DID deactivation request
type DeactivateRequest struct {
	DID     string                 `json:"did"`
	Options map[string]interface{} `json:"options,omitempty"`
	Secret  map[string]interface{} `json:"secret,omitempty"`
}

// DeactivateResponse represents a DID deactivation response
type DeactivateResponse struct {
	JobID                   string                  `json:"jobId"`
	DIDState                DIDState                `json:"didState"`
	DIDRegistrationMetadata DIDRegistrationMetadata `json:"didRegistrationMetadata"`
	DIDDocumentMetadata     DIDDocumentMetadata     `json:"didDocumentMetadata"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string            `json:"error"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details,omitempty"`
	RequestID string            `json:"requestId,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// NewDeactivateHandler creates a new deactivate handler
func NewDeactivateHandler(accClient acc.Client, authPolicy policy.AuthPolicy) *DeactivateHandler {
	return &DeactivateHandler{
		accClient:  accClient,
		authPolicy: authPolicy,
	}
}

// Deactivate handles POST /deactivate requests
func (h *DeactivateHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req DeactivateRequest
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
	requiredKeyPage, err := h.authPolicy.GetRequiredKeyPage(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Create deactivation document
	deactivationDoc := map[string]interface{}{
		"@context":    []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":          req.DID,
		"deactivated": true,
	}

	// For deactivation, we should get the previous version ID
	// In a real implementation, this would query the current DID document
	previousVersionID := h.getPreviousVersionID(req.DID)

	// Build envelope
	envelope, err := ops.BuildEnvelope(deactivationDoc, requiredKeyPage, previousVersionID)
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
	response := DeactivateResponse{
		JobID: h.generateJobID(),
		DIDState: DIDState{
			DID:    req.DID,
			State:  "finished",
			Action: "deactivate",
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

// validateDeactivateRequest validates the deactivate request
func (h *DeactivateHandler) validateDeactivateRequest(req *DeactivateRequest) error {
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

// getDataAccountURL constructs the data account URL for a DID
func (h *DeactivateHandler) getDataAccountURL(did string) string {
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
func (h *DeactivateHandler) generateJobID() string {
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}

// writeError writes an error response
func (h *DeactivateHandler) writeError(w http.ResponseWriter, errorCode, message string, status int, details map[string]string) {
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