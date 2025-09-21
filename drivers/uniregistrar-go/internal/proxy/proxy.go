package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Proxy handles proxying requests to the registrar
type Proxy struct {
	registrarURL string
	client       *http.Client
}

// UniversalRegistrarRequest represents the Universal Registrar API request
type UniversalRegistrarRequest struct {
	DID         string                 `json:"did,omitempty"`
	DIDDocument map[string]interface{} `json:"didDocument,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Secret      map[string]interface{} `json:"secret,omitempty"`
}

// RegistrarRequest represents the request to registrar-go
type RegistrarRequest struct {
	DID         string                 `json:"did"`
	DIDDocument map[string]interface{} `json:"didDocument,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Secret      map[string]interface{} `json:"secret,omitempty"`
}

// RegistrarResponse represents the response from registrar-go
type RegistrarResponse struct {
	JobID                   string                 `json:"jobId"`
	DIDState                DIDState               `json:"didState"`
	DIDRegistrationMetadata map[string]interface{} `json:"didRegistrationMetadata"`
	DIDDocumentMetadata     map[string]interface{} `json:"didDocumentMetadata"`
}

// DIDState represents the state of a DID operation
type DIDState struct {
	DID    string `json:"did"`
	State  string `json:"state"`
	Action string `json:"action"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	JobID    string                 `json:"jobId,omitempty"`
	DIDState map[string]interface{} `json:"didState"`
}

// New creates a new proxy instance
func New(registrarURL string) *Proxy {
	return &Proxy{
		registrarURL: strings.TrimSuffix(registrarURL, "/"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateHandler handles DID creation requests
func (p *Proxy) CreateHandler(w http.ResponseWriter, r *http.Request) {
	// Check method parameter
	method := r.URL.Query().Get("method")
	if method != "acc" {
		p.writeError(w, "invalidMethod", fmt.Sprintf("Unsupported method: %s", method), http.StatusBadRequest)
		return
	}

	p.proxyRequest(w, r, "/create", "create")
}

// UpdateHandler handles DID update requests
func (p *Proxy) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Check method parameter
	method := r.URL.Query().Get("method")
	if method != "" && method != "acc" {
		p.writeError(w, "invalidMethod", fmt.Sprintf("Unsupported method: %s", method), http.StatusBadRequest)
		return
	}

	p.proxyRequest(w, r, "/update", "update")
}

// DeactivateHandler handles DID deactivation requests
func (p *Proxy) DeactivateHandler(w http.ResponseWriter, r *http.Request) {
	// Check method parameter
	method := r.URL.Query().Get("method")
	if method != "" && method != "acc" {
		p.writeError(w, "invalidMethod", fmt.Sprintf("Unsupported method: %s", method), http.StatusBadRequest)
		return
	}

	p.proxyRequest(w, r, "/deactivate", "deactivate")
}

// proxyRequest forwards a request to the registrar
func (p *Proxy) proxyRequest(w http.ResponseWriter, r *http.Request, endpoint string, action string) {
	// Parse request body
	var uniReq UniversalRegistrarRequest
	if err := json.NewDecoder(r.Body).Decode(&uniReq); err != nil {
		p.writeError(w, "invalidRequest", "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	// For create operations, validate DID method if DID is provided
	if action == "create" && uniReq.DID != "" && !strings.HasPrefix(uniReq.DID, "did:acc:") {
		p.writeError(w, "invalidDid", "Only did:acc method is supported", http.StatusBadRequest)
		return
	}

	// For update/deactivate, DID is required and must be did:acc
	if action != "create" {
		if uniReq.DID == "" {
			p.writeError(w, "invalidRequest", "DID is required", http.StatusBadRequest)
			return
		}
		if !strings.HasPrefix(uniReq.DID, "did:acc:") {
			p.writeError(w, "invalidDid", "Only did:acc method is supported", http.StatusBadRequest)
			return
		}
	}

	// Build registrar request
	regReq := RegistrarRequest{
		DID:         uniReq.DID,
		DIDDocument: uniReq.DIDDocument,
		Options:     uniReq.Options,
		Secret:      uniReq.Secret,
	}

	// Marshal request
	reqBody, err := json.Marshal(regReq)
	if err != nil {
		p.writeError(w, "internalError", "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	// Build registrar URL
	registrarURL := fmt.Sprintf("%s%s", p.registrarURL, endpoint)

	// Make request to registrar
	resp, err := p.client.Post(registrarURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		p.writeError(w, "internalError", fmt.Sprintf("Failed to contact registrar: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.writeError(w, "internalError", "Failed to read registrar response", http.StatusInternalServerError)
		return
	}

	// Handle non-200 responses from registrar
	if resp.StatusCode != http.StatusOK {
		// Try to parse error response
		var errorData map[string]interface{}
		if err := json.Unmarshal(body, &errorData); err == nil {
			// Forward error as DID state
			p.writeErrorWithState(w, action, fmt.Sprintf("Registrar error: %v", errorData), resp.StatusCode)
			return
		}
		p.writeErrorWithState(w, action, string(body), resp.StatusCode)
		return
	}

	// Parse registrar response
	var regResp RegistrarResponse
	if err := json.Unmarshal(body, &regResp); err != nil {
		p.writeError(w, "internalError", "Invalid registrar response format", http.StatusInternalServerError)
		return
	}

	// The registrar response is already in the correct format for Universal Registrar
	// Just forward it directly
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(regResp); err != nil {
		// Log error but response is already committed
		fmt.Printf("Failed to encode response: %v\n", err)
	}
}

// writeError writes an error response
func (p *Proxy) writeError(w http.ResponseWriter, error string, message string, status int) {
	response := ErrorResponse{
		JobID: fmt.Sprintf("error-%d", time.Now().UnixNano()),
		DIDState: map[string]interface{}{
			"state":  "failed",
			"reason": message,
			"error":  error,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// writeErrorWithState writes an error response with DID state
func (p *Proxy) writeErrorWithState(w http.ResponseWriter, action string, message string, status int) {
	response := ErrorResponse{
		JobID: fmt.Sprintf("error-%d", time.Now().UnixNano()),
		DIDState: map[string]interface{}{
			"state":  "failed",
			"action": action,
			"reason": message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
