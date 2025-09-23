package accdid

import (
	"encoding/json"
)

// DIDDocument represents a W3C DID Document
type DIDDocument map[string]interface{}

// Service represents a service in a DID Document
type Service struct {
	ID              string      `json:"id"`
	Type            string      `json:"type"`
	ServiceEndpoint interface{} `json:"serviceEndpoint"`
}

// VerificationMethod represents a verification method in a DID Document
type VerificationMethod struct {
	ID                 string      `json:"id"`
	Type               string      `json:"type"`
	Controller         string      `json:"controller"`
	PublicKeyJwk       interface{} `json:"publicKeyJwk,omitempty"`
	PublicKeyMultibase string      `json:"publicKeyMultibase,omitempty"`
	PublicKeyBase58    string      `json:"publicKeyBase58,omitempty"`
}

// ResolutionResult represents the result of DID resolution
type ResolutionResult struct {
	// DIDDocument contains the resolved DID document
	DIDDocument interface{} `json:"didDocument"`

	// Metadata contains resolution metadata
	Metadata map[string]interface{} `json:"resolutionMetadata,omitempty"`

	// DocumentMetadata contains document metadata
	DocumentMetadata map[string]interface{} `json:"didDocumentMetadata,omitempty"`
}

// ErrorEnvelope represents a structured API error response
type ErrorEnvelope struct {
	Code    string      `json:"code"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// NativeRegisterRequest represents a request to register a new DID
type NativeRegisterRequest struct {
	DID         string          `json:"did"`
	DIDDocument json.RawMessage `json:"didDocument"`
}

// NativeUpdateRequest represents a request to update an existing DID
type NativeUpdateRequest struct {
	DID   string          `json:"did"`
	Patch json.RawMessage `json:"patch"`
}

// NativeDeactivateRequest represents a request to deactivate a DID
type NativeDeactivateRequest struct {
	DID    string `json:"did"`
	Reason string `json:"reason,omitempty"`
}

// RegistrarResponse represents a successful registrar operation response
type RegistrarResponse struct {
	TransactionID string `json:"transactionId"`
	JobID         string `json:"jobId,omitempty"`
}

// UniversalResolveResponse represents Universal Resolver response format
type UniversalResolveResponse struct {
	DIDResolutionMetadata map[string]interface{} `json:"didResolutionMetadata,omitempty"`
	DIDDocument           interface{}            `json:"didDocument"`
	DIDDocumentMetadata   map[string]interface{} `json:"didDocumentMetadata,omitempty"`
}