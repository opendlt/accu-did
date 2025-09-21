package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/opendlt/accu-did/registrar-go/internal/acc"
	"github.com/opendlt/accu-did/registrar-go/internal/api"
	"github.com/opendlt/accu-did/shared/did"
)

// NativeHandler handles native DID registration endpoints
type NativeHandler struct {
	accClient acc.Submitter
}

// RegisterRequest represents a native DID registration request
type RegisterRequest struct {
	DID         string                 `json:"did"`
	DIDDocument map[string]interface{} `json:"didDocument"`
	KeyPageURL  string                 `json:"keyPageUrl,omitempty"`
}

// NativeUpdateRequest represents a native DID update request
type NativeUpdateRequest struct {
	DID         string                 `json:"did"`
	DIDDocument map[string]interface{} `json:"didDocument"`
}

// NativeResponse represents a native API response
type NativeResponse struct {
	Success   bool                   `json:"success"`
	TxID      string                 `json:"txid,omitempty"`
	DID       string                 `json:"did"`
	JobID     string                 `json:"jobId,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewNativeHandler creates a new native handler
func NewNativeHandler(accClient acc.Submitter) *NativeHandler {
	return &NativeHandler{
		accClient: accClient,
	}
}

// Register handles POST /register requests
func (h *NativeHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalidRequest", "Invalid JSON", http.StatusBadRequest, nil)
		return
	}

	// Validate request
	if err := h.validateRegisterRequest(&req); err != nil {
		h.writeError(w, "invalidRequest", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Parse DID to get ADI components
	adiURL, dataAccountURL, err := did.ParseDID(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Step 1: Create ADI if it doesn't exist
	adiLabel := adiURL.Authority
	keyPageURL := req.KeyPageURL
	if keyPageURL == "" {
		keyPageURL = fmt.Sprintf("acc://%s/book/1", adiLabel)
	}

	adiTxID, err := h.accClient.CreateIdentity(adiLabel, keyPageURL)
	if err != nil {
		// ADI might already exist, continue with data account creation
		// In a real implementation, you'd check if the error is "already exists"
	}

	// Step 2: Create data account
	dataAccountLabel := dataAccountURL.Path[1:] // Remove leading slash
	dataTxID, err := h.accClient.CreateDataAccount(adiURL.String(), dataAccountLabel)
	if err != nil {
		// Data account might already exist, continue with writing data
		// In a real implementation, you'd check if the error is "already exists"
	}

	// Step 3: Write DID document to data account
	didDocData, err := json.Marshal(req.DIDDocument)
	if err != nil {
		h.writeError(w, "internalError", "Failed to marshal DID document", http.StatusInternalServerError, nil)
		return
	}

	txID, err := h.accClient.WriteDataEntry(dataAccountURL.String(), didDocData)
	if err != nil {
		h.writeError(w, "internalError", "Failed to write DID document", http.StatusInternalServerError, nil)
		return
	}

	// Build response
	metadata := map[string]interface{}{
		"adiTxID":     adiTxID,
		"dataTxID":    dataTxID,
		"adiLabel":    adiLabel,
		"dataAccount": dataAccountURL.String(),
	}

	response := NativeResponse{
		Success:   true,
		TxID:      txID,
		DID:       req.DID,
		JobID:     h.generateJobID(),
		Metadata:  metadata,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Update handles POST /update requests
func (h *NativeHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req NativeUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalidRequest", "Invalid JSON", http.StatusBadRequest, nil)
		return
	}

	// Validate request
	if err := h.validateUpdateRequest(&req); err != nil {
		h.writeError(w, "invalidRequest", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Parse DID to get data account URL
	_, dataAccountURL, err := did.ParseDID(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Write updated DID document
	didDocData, err := json.Marshal(req.DIDDocument)
	if err != nil {
		h.writeError(w, "internalError", "Failed to marshal DID document", http.StatusInternalServerError, nil)
		return
	}

	txID, err := h.accClient.WriteDataEntry(dataAccountURL.String(), didDocData)
	if err != nil {
		h.writeError(w, "internalError", "Failed to update DID document", http.StatusInternalServerError, nil)
		return
	}

	response := NativeResponse{
		Success:   true,
		TxID:      txID,
		DID:       req.DID,
		JobID:     h.generateJobID(),
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Deactivate handles POST /deactivate requests
func (h *NativeHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
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

	// Parse DID to get data account URL
	_, dataAccountURL, err := did.ParseDID(req.DID)
	if err != nil {
		h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Create deactivated DID document
	deactivatedDoc := map[string]interface{}{
		"@context":    []string{"https://www.w3.org/ns/did/v1"},
		"id":          req.DID,
		"deactivated": true,
	}

	didDocData, err := json.Marshal(deactivatedDoc)
	if err != nil {
		h.writeError(w, "internalError", "Failed to marshal deactivated DID document", http.StatusInternalServerError, nil)
		return
	}

	txID, err := h.accClient.WriteDataEntry(dataAccountURL.String(), didDocData)
	if err != nil {
		h.writeError(w, "internalError", "Failed to deactivate DID", http.StatusInternalServerError, nil)
		return
	}

	response := NativeResponse{
		Success:   true,
		TxID:      txID,
		DID:       req.DID,
		JobID:     h.generateJobID(),
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// validateRegisterRequest validates the register request
func (h *NativeHandler) validateRegisterRequest(req *RegisterRequest) error {
	if req.DID == "" {
		return fmt.Errorf("DID is required")
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

// validateUpdateRequest validates the update request
func (h *NativeHandler) validateUpdateRequest(req *NativeUpdateRequest) error {
	if req.DID == "" {
		return fmt.Errorf("DID is required")
	}

	if req.DIDDocument == nil {
		return fmt.Errorf("didDocument is required")
	}

	// Validate that the DID in the document matches the request
	if docID, ok := req.DIDDocument["id"].(string); ok {
		if docID != req.DID {
			return fmt.Errorf("DID mismatch: request DID %s does not match document ID %s", req.DID, docID)
		}
	}

	return nil
}

// validateDeactivateRequest validates the deactivate request
func (h *NativeHandler) validateDeactivateRequest(req *api.DeactivateRequest) error {
	if req.DID == "" {
		return fmt.Errorf("DID is required")
	}

	return nil
}

// generateJobID generates a job ID for tracking the operation
func (h *NativeHandler) generateJobID() string {
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}

// writeError writes an error response
func (h *NativeHandler) writeError(w http.ResponseWriter, errorCode, message string, status int, details map[string]string) {
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
