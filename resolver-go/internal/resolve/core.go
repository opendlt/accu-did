package resolve

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/opendlt/accu-did/resolver-go/internal/acc"
	"github.com/opendlt/accu-did/resolver-go/internal/canon"
	"github.com/opendlt/accu-did/resolver-go/internal/normalize"
	"github.com/opendlt/accu-did/shared/did"
)

// DIDResolutionResult represents a W3C DID Resolution result
type DIDResolutionResult struct {
	DIDDocument           interface{}           `json:"didDocument"`
	DIDDocumentMetadata   DIDDocumentMetadata   `json:"didDocumentMetadata"`
	DIDResolutionMetadata DIDResolutionMetadata `json:"didResolutionMetadata"`
}

// DIDDocumentMetadata represents DID document metadata
type DIDDocumentMetadata struct {
	VersionID     string     `json:"versionId"`
	Created       time.Time  `json:"created"`
	Updated       time.Time  `json:"updated"`
	Deactivated   bool       `json:"deactivated"`
	NextUpdate    *time.Time `json:"nextUpdate,omitempty"`
	NextVersionID *string    `json:"nextVersionId,omitempty"`
}

// DIDResolutionMetadata represents DID resolution metadata
type DIDResolutionMetadata struct {
	ContentType string    `json:"contentType"`
	Retrieved   time.Time `json:"retrieved"`
	Pattern     string    `json:"pattern,omitempty"`
	DriverURL   string    `json:"driverUrl,omitempty"`
	Duration    int       `json:"duration,omitempty"`
}

// Custom error types
type NotFoundError struct {
	DID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("DID not found: %s", e.DID)
}

type InvalidDIDError struct {
	DID    string
	Reason string
}

func (e *InvalidDIDError) Error() string {
	return fmt.Sprintf("Invalid DID %s: %s", e.DID, e.Reason)
}

type DeactivatedError struct {
	DID string
}

func (e *DeactivatedError) Error() string {
	return fmt.Sprintf("DID is deactivated: %s", e.DID)
}

// ResolveDID resolves a DID according to the W3C DID Core specification
func ResolveDID(client acc.Client, didStr string, versionTime *time.Time) (*DIDResolutionResult, error) {
	start := time.Now()

	// Step 1: Parse DID into Accumulate URLs
	adiURL, dataAccountURL, err := did.ParseDID(didStr)
	if err != nil {
		return nil, &InvalidDIDError{DID: didStr, Reason: err.Error()}
	}

	// Step 2: Read DID document from data account
	var didDocBytes []byte
	if versionTime != nil {
		// TODO: Implement versioned data entry reading
		didDocBytes, err = client.GetDataAccountEntry(dataAccountURL)
	} else {
		didDocBytes, err = client.GetDataAccountEntry(dataAccountURL)
	}
	if err != nil {
		return nil, &NotFoundError{DID: didStr}
	}

	// Step 3: Parse DID document
	var didDoc map[string]interface{}
	if err := json.Unmarshal(didDocBytes, &didDoc); err != nil {
		return nil, fmt.Errorf("invalid DID document JSON: %w", err)
	}

	// Step 4: Check if deactivated
	if deactivated, exists := didDoc["deactivated"].(bool); exists && deactivated {
		return nil, &DeactivatedError{DID: didStr}
	}

	// Step 5: Build resolution result
	result := &DIDResolutionResult{
		DIDDocument: didDoc,
		DIDDocumentMetadata: DIDDocumentMetadata{
			VersionID:   fmt.Sprintf("%d", time.Now().Unix()), // TODO: Real version from metadata
			Created:     time.Now().UTC(),                     // TODO: Real timestamp from metadata
			Updated:     time.Now().UTC(),                     // TODO: Real timestamp from metadata
			Deactivated: false,
		},
		DIDResolutionMetadata: DIDResolutionMetadata{
			ContentType: "application/did+json",
			Retrieved:   time.Now().UTC(),
			Pattern:     "^did:acc:",
			Duration:    int(time.Since(start).Milliseconds()),
		},
	}

	return result, nil
}