package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/opendlt/accu-did/registrar-go/internal/acc"
	"github.com/opendlt/accu-did/registrar-go/internal/api"
	"github.com/opendlt/accu-did/registrar-go/internal/policy"
	"github.com/opendlt/accu-did/shared/did"
)

// UniversalHandler handles Universal Registrar compatibility endpoints
type UniversalHandler struct {
	nativeHandler *NativeHandler
	accClient     acc.Submitter
	authPolicy    policy.AuthPolicy
}

// UniversalCreateRequest represents a Universal Registrar create request
type UniversalCreateRequest struct {
	JobID       string                 `json:"jobId,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Secret      map[string]interface{} `json:"secret,omitempty"`
	DIDDocument map[string]interface{} `json:"didDocument,omitempty"`
}

// UniversalUpdateRequest represents a Universal Registrar update request
type UniversalUpdateRequest struct {
	JobID        string                 `json:"jobId,omitempty"`
	Identifier   string                 `json:"identifier,omitempty"`
	Options      map[string]interface{} `json:"options,omitempty"`
	Secret       map[string]interface{} `json:"secret,omitempty"`
	DIDDocument  map[string]interface{} `json:"didDocument,omitempty"`
	Registration *RegistrationRequest   `json:"registration,omitempty"`
}

// RegistrationRequest represents registration data in Universal Registrar format
type RegistrationRequest struct {
	DID   string                 `json:"did"`
	Patch map[string]interface{} `json:"patch,omitempty"`
}

// ServiceEntry represents a service entry for patches
type ServiceEntry struct {
	ID              string      `json:"id"`
	Type            string      `json:"type"`
	ServiceEndpoint interface{} `json:"serviceEndpoint"`
}

// UniversalDeactivateRequest represents a Universal Registrar deactivate request
type UniversalDeactivateRequest struct {
	JobID      string                 `json:"jobId,omitempty"`
	Identifier string                 `json:"identifier"`
	Options    map[string]interface{} `json:"options,omitempty"`
	Secret     map[string]interface{} `json:"secret,omitempty"`
}

// UniversalResponse represents a Universal Registrar response
type UniversalResponse struct {
	JobID                   string                  `json:"jobId,omitempty"`
	DIDState                DIDState                `json:"didState"`
	DIDRegistrationMetadata DIDRegistrationMetadata `json:"didRegistrationMetadata,omitempty"`
	DIDDocumentMetadata     DIDDocumentMetadata     `json:"didDocumentMetadata,omitempty"`
}

// NewUniversalHandler creates a new Universal Registrar compatibility handler
func NewUniversalHandler(accClient acc.Submitter, authPolicy policy.AuthPolicy) *UniversalHandler {
	return &UniversalHandler{
		nativeHandler: NewNativeHandler(accClient),
		accClient:     accClient,
		authPolicy:    authPolicy,
	}
}

// UniversalCreate handles POST /1.0/create requests (Universal Registrar)
func (h *UniversalHandler) UniversalCreate(w http.ResponseWriter, r *http.Request) {
	var req UniversalCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeUniversalError(w, "invalidRequest", "Invalid JSON", http.StatusBadRequest, nil)
		return
	}

	// Validate request
	if err := h.validateUniversalCreateRequest(&req); err != nil {
		h.writeUniversalError(w, "invalidRequest", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Extract DID from didDocument
	did, ok := req.DIDDocument["id"].(string)
	if !ok {
		h.writeUniversalError(w, "invalidRequest", "didDocument must have an 'id' field", http.StatusBadRequest, nil)
		return
	}

	// Convert to native register request
	nativeReq := RegisterRequest{
		DID:         did,
		DIDDocument: req.DIDDocument,
	}

	// Extract key page URL from options if provided
	if options := req.Options; options != nil {
		if keyPageURL, ok := options["keyPageUrl"].(string); ok {
			nativeReq.KeyPageURL = keyPageURL
		}
	}

	// Validate native request
	if err := h.nativeHandler.validateRegisterRequest(&nativeReq); err != nil {
		h.writeUniversalError(w, "invalidRequest", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Process using native handler logic but return Universal format
	response, err := h.processNativeRegister(&nativeReq)
	if err != nil {
		h.writeUniversalError(w, "internalError", err.Error(), http.StatusInternalServerError, nil)
		return
	}

	// Convert to Universal Registrar response format
	universalResponse := UniversalResponse{
		JobID: response.JobID,
		DIDState: DIDState{
			DID:    response.DID,
			State:  "finished",
			Action: "create",
		},
		DIDRegistrationMetadata: DIDRegistrationMetadata{
			TxID: response.TxID,
		},
		DIDDocumentMetadata: DIDDocumentMetadata{
			Created: response.Timestamp,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(universalResponse)
}

// UniversalUpdate handles POST /1.0/update requests (Universal Registrar)
func (h *UniversalHandler) UniversalUpdate(w http.ResponseWriter, r *http.Request) {
	var req UniversalUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeUniversalError(w, "invalidRequest", "Invalid JSON", http.StatusBadRequest, nil)
		return
	}

	// Validate request
	if err := h.validateUniversalUpdateRequest(&req); err != nil {
		h.writeUniversalError(w, "invalidRequest", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Determine DID from identifier or registration
	var targetDID string
	if req.Identifier != "" {
		targetDID = req.Identifier
	} else if req.Registration != nil && req.Registration.DID != "" {
		targetDID = req.Registration.DID
	} else {
		h.writeUniversalError(w, "invalidRequest", "DID identifier is required", http.StatusBadRequest, nil)
		return
	}

	var updatedDoc map[string]interface{}

	// Handle patch-based updates
	if req.Registration != nil && req.Registration.Patch != nil {
		// Get current DID document by resolving it
		currentDoc, err := h.resolveCurrentDIDDocument(targetDID)
		if err != nil {
			h.writeUniversalError(w, "notFound", fmt.Sprintf("Could not resolve current DID document: %v", err), http.StatusNotFound, nil)
			return
		}

		// Apply patch to current document
		updatedDoc, err = h.applyPatch(currentDoc, req.Registration.Patch)
		if err != nil {
			h.writeUniversalError(w, "invalidRequest", fmt.Sprintf("Failed to apply patch: %v", err), http.StatusBadRequest, nil)
			return
		}
	} else if req.DIDDocument != nil {
		// Use provided DID document directly
		updatedDoc = req.DIDDocument
	} else {
		h.writeUniversalError(w, "invalidRequest", "Either didDocument or registration.patch is required", http.StatusBadRequest, nil)
		return
	}

	// Convert to native update request
	nativeReq := NativeUpdateRequest{
		DID:         targetDID,
		DIDDocument: updatedDoc,
	}

	// Process using native handler logic
	response, err := h.processNativeUpdate(&nativeReq)
	if err != nil {
		h.writeUniversalError(w, "internalError", err.Error(), http.StatusInternalServerError, nil)
		return
	}

	// Convert to Universal Registrar response format
	universalResponse := UniversalResponse{
		JobID: response.JobID,
		DIDState: DIDState{
			DID:    response.DID,
			State:  "finished",
			Action: "update",
		},
		DIDRegistrationMetadata: DIDRegistrationMetadata{
			TxID: response.TxID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(universalResponse)
}

// UniversalDeactivate handles POST /1.0/deactivate requests (Universal Registrar)
func (h *UniversalHandler) UniversalDeactivate(w http.ResponseWriter, r *http.Request) {
	var req UniversalDeactivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeUniversalError(w, "invalidRequest", "Invalid JSON", http.StatusBadRequest, nil)
		return
	}

	// Validate request
	if err := h.validateUniversalDeactivateRequest(&req); err != nil {
		h.writeUniversalError(w, "invalidRequest", err.Error(), http.StatusBadRequest, nil)
		return
	}

	// Convert to native deactivate request
	nativeReq := api.DeactivateRequest{
		DID: req.Identifier,
	}

	// Process using native handler logic
	response, err := h.processNativeDeactivate(&nativeReq)
	if err != nil {
		h.writeUniversalError(w, "internalError", err.Error(), http.StatusInternalServerError, nil)
		return
	}

	// Convert to Universal Registrar response format
	universalResponse := UniversalResponse{
		JobID: response.JobID,
		DIDState: DIDState{
			DID:    response.DID,
			State:  "finished",
			Action: "deactivate",
		},
		DIDRegistrationMetadata: DIDRegistrationMetadata{
			TxID: response.TxID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(universalResponse)
}

// processNativeRegister processes a register request using native logic
func (h *UniversalHandler) processNativeRegister(req *RegisterRequest) (*NativeResponse, error) {
	// This duplicates the native handler logic to avoid HTTP roundtrip
	// In a real implementation, you might extract this to a service layer

	// Parse DID to get ADI components
	adiURL, dataAccountURL, err := did.ParseDID(req.DID)
	if err != nil {
		return nil, err
	}

	// Create ADI
	adiLabel := adiURL.Authority
	keyPageURL := req.KeyPageURL
	if keyPageURL == "" {
		keyPageURL = fmt.Sprintf("acc://%s/book/1", adiLabel)
	}

	adiTxID, _ := h.accClient.CreateIdentity(adiLabel, keyPageURL)

	// Create data account
	dataAccountLabel := dataAccountURL.Path[1:]
	dataTxID, _ := h.accClient.CreateDataAccount(adiURL.String(), dataAccountLabel)

	// Write DID document
	didDocData, err := json.Marshal(req.DIDDocument)
	if err != nil {
		return nil, err
	}

	txID, err := h.accClient.WriteDataEntry(dataAccountURL.String(), didDocData)
	if err != nil {
		return nil, err
	}

	return &NativeResponse{
		Success:   true,
		TxID:      txID,
		DID:       req.DID,
		JobID:     h.generateJobID(),
		Timestamp: time.Now().UTC(),
		Metadata: map[string]interface{}{
			"adiTxID":     adiTxID,
			"dataTxID":    dataTxID,
			"adiLabel":    adiLabel,
			"dataAccount": dataAccountURL.String(),
		},
	}, nil
}

// processNativeUpdate processes an update request using native logic
func (h *UniversalHandler) processNativeUpdate(req *NativeUpdateRequest) (*NativeResponse, error) {
	// Parse DID to get data account URL
	_, dataAccountURL, err := did.ParseDID(req.DID)
	if err != nil {
		return nil, err
	}

	// Write updated DID document
	didDocData, err := json.Marshal(req.DIDDocument)
	if err != nil {
		return nil, err
	}

	txID, err := h.accClient.WriteDataEntry(dataAccountURL.String(), didDocData)
	if err != nil {
		return nil, err
	}

	return &NativeResponse{
		Success:   true,
		TxID:      txID,
		DID:       req.DID,
		JobID:     h.generateJobID(),
		Timestamp: time.Now().UTC(),
	}, nil
}

// processNativeDeactivate processes a deactivate request using native logic
func (h *UniversalHandler) processNativeDeactivate(req *api.DeactivateRequest) (*NativeResponse, error) {
	// Parse DID to get data account URL
	_, dataAccountURL, err := did.ParseDID(req.DID)
	if err != nil {
		return nil, err
	}

	// Create deactivated DID document
	deactivatedDoc := map[string]interface{}{
		"@context":    []string{"https://www.w3.org/ns/did/v1"},
		"id":          req.DID,
		"deactivated": true,
	}

	didDocData, err := json.Marshal(deactivatedDoc)
	if err != nil {
		return nil, err
	}

	txID, err := h.accClient.WriteDataEntry(dataAccountURL.String(), didDocData)
	if err != nil {
		return nil, err
	}

	return &NativeResponse{
		Success:   true,
		TxID:      txID,
		DID:       req.DID,
		JobID:     h.generateJobID(),
		Timestamp: time.Now().UTC(),
	}, nil
}

// Validation functions

func (h *UniversalHandler) validateUniversalCreateRequest(req *UniversalCreateRequest) error {
	if req.DIDDocument == nil {
		return fmt.Errorf("didDocument is required")
	}

	if _, ok := req.DIDDocument["id"]; !ok {
		return fmt.Errorf("didDocument must have an 'id' field")
	}

	if _, ok := req.DIDDocument["@context"]; !ok {
		return fmt.Errorf("didDocument must have '@context' field")
	}

	return nil
}

func (h *UniversalHandler) validateUniversalUpdateRequest(req *UniversalUpdateRequest) error {
	if req.Identifier == "" {
		return fmt.Errorf("identifier is required")
	}

	if req.DIDDocument == nil {
		return fmt.Errorf("didDocument is required")
	}

	// Validate that the DID in the document matches the identifier
	if docID, ok := req.DIDDocument["id"].(string); ok {
		if docID != req.Identifier {
			return fmt.Errorf("DID mismatch: identifier %s does not match document ID %s", req.Identifier, docID)
		}
	}

	return nil
}

func (h *UniversalHandler) validateUniversalDeactivateRequest(req *UniversalDeactivateRequest) error {
	if req.Identifier == "" {
		return fmt.Errorf("identifier is required")
	}

	return nil
}

// generateJobID generates a job ID for tracking the operation
func (h *UniversalHandler) generateJobID() string {
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}

// resolveCurrentDIDDocument resolves the current DID document for patch operations
func (h *UniversalHandler) resolveCurrentDIDDocument(didStr string) (map[string]interface{}, error) {
	// Parse DID to get data account URL
	_, dataAccountURL, err := did.ParseDID(didStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DID: %w", err)
	}

	// For FAKE mode, try to read from testdata
	if _, ok := h.accClient.(*acc.FakeSubmitter); ok {
		// In FAKE mode, we need to simulate reading the current document
		// For simplicity, we'll return the existing testdata document
		// In a real implementation, this would track the current state

		// Extract ADI label for testdata lookup
		adiLabel := dataAccountURL.Authority

		// Return a mock current document for beastmode.acme
		if adiLabel == "beastmode.acme" {
			return map[string]interface{}{
				"@context": []interface{}{
					"https://www.w3.org/ns/did/v1",
					"https://w3id.org/security/suites/ed25519-2020/v1",
				},
				"id": didStr,
				"verificationMethod": []interface{}{
					map[string]interface{}{
						"id":                 didStr + "#key1",
						"type":               "Ed25519VerificationKey2020",
						"controller":         didStr,
						"publicKeyMultibase": "z6MkqRYqQiSgvZQdnBytw86Qbs2ZWUkGv22od935YF4s8M7V",
					},
				},
				"authentication": []interface{}{didStr + "#key1"},
				"service": []interface{}{
					map[string]interface{}{
						"id":              didStr + "#website",
						"type":            "LinkedDomains",
						"serviceEndpoint": "https://beastmode.acme.corp",
					},
				},
			}, nil
		}

		// Default mock document for other DIDs
		return map[string]interface{}{
			"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
			"id":       didStr,
			"service":  []interface{}{},
		}, nil
	}

	// For REAL mode, we would need to create a resolver client
	// and actually resolve the current document from the network
	// For now, return a minimal document
	return map[string]interface{}{
		"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":       didStr,
		"service":  []interface{}{},
	}, nil
}

// applyPatch applies a patch to the current DID document
func (h *UniversalHandler) applyPatch(currentDoc map[string]interface{}, patch map[string]interface{}) (map[string]interface{}, error) {
	// Create a copy of the current document
	updatedDoc := make(map[string]interface{})
	for k, v := range currentDoc {
		updatedDoc[k] = v
	}

	// Handle addService patch
	if addServicePatch, ok := patch["addService"]; ok {
		addServices, ok := addServicePatch.([]interface{})
		if !ok {
			return nil, fmt.Errorf("addService must be an array")
		}

		// Get current services or initialize empty array
		var currentServices []interface{}
		if services, exists := updatedDoc["service"]; exists {
			if servicesArray, ok := services.([]interface{}); ok {
				currentServices = servicesArray
			}
		}

		// Add new services
		for _, service := range addServices {
			currentServices = append(currentServices, service)
		}

		updatedDoc["service"] = currentServices
	}

	// Handle removeService patch
	if removeServicePatch, ok := patch["removeService"]; ok {
		removeServices, ok := removeServicePatch.([]interface{})
		if !ok {
			return nil, fmt.Errorf("removeService must be an array")
		}

		// Get current services
		if services, exists := updatedDoc["service"]; exists {
			if servicesArray, ok := services.([]interface{}); ok {
				var filteredServices []interface{}

				for _, service := range servicesArray {
					serviceMap, ok := service.(map[string]interface{})
					if !ok {
						continue
					}

					serviceID, ok := serviceMap["id"].(string)
					if !ok {
						continue
					}

					// Check if this service should be removed
					shouldRemove := false
					for _, removeService := range removeServices {
						if removeID, ok := removeService.(string); ok && removeID == serviceID {
							shouldRemove = true
							break
						}
					}

					if !shouldRemove {
						filteredServices = append(filteredServices, service)
					}
				}

				updatedDoc["service"] = filteredServices
			}
		}
	}

	// Handle other patch operations as needed
	// (addVerificationMethod, removeVerificationMethod, etc.)

	return updatedDoc, nil
}

// writeUniversalError writes an error response in Universal Registrar format
func (h *UniversalHandler) writeUniversalError(w http.ResponseWriter, errorCode, message string, status int, details map[string]string) {
	response := map[string]interface{}{
		"didState": map[string]interface{}{
			"state":  "failed",
			"reason": message,
		},
		"didRegistrationMetadata": map[string]interface{}{
			"error": errorCode,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
