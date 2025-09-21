package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Proxy handles proxying requests to the resolver
type Proxy struct {
	resolverURL string
	client      *http.Client
}

// ResolverResponse represents the response from resolver-go
type ResolverResponse struct {
	DIDDocument         interface{} `json:"didDocument"`
	DIDDocumentMetadata interface{} `json:"didDocumentMetadata"`
}

// UniversalResolverResponse represents the Universal Resolver API response
type UniversalResolverResponse struct {
	DIDDocument           interface{} `json:"didDocument"`
	DIDDocumentMetadata   interface{} `json:"didDocumentMetadata"`
	DIDResolutionMetadata interface{} `json:"didResolutionMetadata"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	DIDResolutionMetadata map[string]interface{} `json:"didResolutionMetadata"`
}

// New creates a new proxy instance
func New(resolverURL string) *Proxy {
	return &Proxy{
		resolverURL: strings.TrimSuffix(resolverURL, "/"),
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ResolveHandler handles DID resolution requests
func (p *Proxy) ResolveHandler(w http.ResponseWriter, r *http.Request) {
	// Extract DID from path
	vars := mux.Vars(r)
	did := vars["did"]

	// Validate DID format
	if !strings.HasPrefix(did, "did:acc:") {
		p.writeError(w, "invalidDid", "Only did:acc method is supported", http.StatusBadRequest)
		return
	}

	// Build resolver URL
	resolverURL := fmt.Sprintf("%s/resolve?did=%s", p.resolverURL, url.QueryEscape(did))

	// Forward query parameters if any
	if r.URL.RawQuery != "" {
		resolverURL = fmt.Sprintf("%s&%s", resolverURL, r.URL.RawQuery)
	}

	// Make request to resolver
	resp, err := p.client.Get(resolverURL)
	if err != nil {
		p.writeError(w, "internalError", fmt.Sprintf("Failed to contact resolver: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.writeError(w, "internalError", "Failed to read resolver response", http.StatusInternalServerError)
		return
	}

	// Handle non-200 responses from resolver
	if resp.StatusCode != http.StatusOK {
		// Try to parse error response
		var errorData map[string]interface{}
		if err := json.Unmarshal(body, &errorData); err == nil {
			// Forward error as resolution metadata
			p.writeError(w, "resolutionError", fmt.Sprintf("Resolver error: %v", errorData), resp.StatusCode)
			return
		}
		p.writeError(w, "resolutionError", string(body), resp.StatusCode)
		return
	}

	// Parse resolver response
	var resolverResp ResolverResponse
	if err := json.Unmarshal(body, &resolverResp); err != nil {
		p.writeError(w, "internalError", "Invalid resolver response format", http.StatusInternalServerError)
		return
	}

	// Build Universal Resolver response
	universalResp := UniversalResolverResponse{
		DIDDocument:         resolverResp.DIDDocument,
		DIDDocumentMetadata: resolverResp.DIDDocumentMetadata,
		DIDResolutionMetadata: map[string]interface{}{
			"contentType": "application/did+ld+json",
			"duration":    0,
			"did": map[string]interface{}{
				"didString":        did,
				"methodSpecificId": strings.TrimPrefix(did, "did:acc:"),
				"method":           "acc",
			},
		},
	}

	// Write response
	w.Header().Set("Content-Type", "application/did+ld+json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(universalResp); err != nil {
		// Log error but response is already committed
		fmt.Printf("Failed to encode response: %v\n", err)
	}
}

// writeError writes an error response in Universal Resolver format
func (p *Proxy) writeError(w http.ResponseWriter, error string, message string, status int) {
	response := ErrorResponse{
		DIDResolutionMetadata: map[string]interface{}{
			"error":   error,
			"message": message,
		},
	}

	w.Header().Set("Content-Type", "application/did+ld+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
