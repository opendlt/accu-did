package resolve

import (
	"fmt"
	"time"

	"github.com/opendlt/accu-did/resolver-go/internal/acc"
	"github.com/opendlt/accu-did/resolver-go/internal/canon"
	"github.com/opendlt/accu-did/resolver-go/internal/normalize"
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
func ResolveDID(client acc.Client, did string, versionTime *time.Time) (*DIDResolutionResult, error) {
	start := time.Now()

	// Step 1: Normalize DID
	normalizedDID, adi, err := normalize.NormalizeDID(did)
	if err != nil {
		return nil, &InvalidDIDError{DID: did, Reason: err.Error()}
	}

	// Step 2: Get DID entry from Accumulate
	var envelope acc.Envelope
	if versionTime != nil {
		envelope, err = client.GetEntryAtTime(adi, *versionTime)
	} else {
		envelope, err = client.GetLatestDIDEntry(adi)
	}
	if err != nil {
		return nil, &NotFoundError{DID: normalizedDID}
	}

	// Step 3: Validate content hash
	if envelope.Meta.Proof.ContentHash != "" {
		canonical, err := canon.Canonicalize(envelope.Document)
		if err != nil {
			return nil, fmt.Errorf("failed to canonicalize document: %w", err)
		}

		expectedHash := canon.SHA256(canonical)
		if expectedHash != envelope.Meta.Proof.ContentHash {
			return nil, fmt.Errorf("content hash mismatch")
		}
	}

	// Step 4: Check if deactivated
	if deactivated, exists := envelope.Document["deactivated"].(bool); exists && deactivated {
		return nil, &DeactivatedError{DID: normalizedDID}
	}

	// Step 5: Build resolution result
	result := &DIDResolutionResult{
		DIDDocument: envelope.Document,
		DIDDocumentMetadata: DIDDocumentMetadata{
			VersionID:   envelope.Meta.VersionID,
			Created:     envelope.Meta.Timestamp, // For now, use same timestamp
			Updated:     envelope.Meta.Timestamp,
			Deactivated: false,
		},
		DIDResolutionMetadata: DIDResolutionMetadata{
			ContentType: "application/did+ld+json",
			Retrieved:   time.Now().UTC(),
			Pattern:     "^did:acc:",
			Duration:    int(time.Since(start).Milliseconds()),
		},
	}

	return result, nil
}