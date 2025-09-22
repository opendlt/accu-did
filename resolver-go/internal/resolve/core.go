package resolve

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"gitlab.com/accumulatenetwork/accumulate/pkg/url"

	"github.com/opendlt/accu-did/resolver-go/internal/acc"
	"github.com/opendlt/accu-did/shared/did"
)

// ResolveOrder defines the ordering strategy for selecting the latest DID document
type ResolveOrder string

const (
	ResolveOrderSequence  ResolveOrder = "sequence"
	ResolveOrderTimestamp ResolveOrder = "timestamp"
)

// DIDResolutionResult represents a W3C DID Resolution result
type DIDResolutionResult struct {
	DIDDocument           interface{}           `json:"didDocument"`
	DIDDocumentMetadata   DIDDocumentMetadata   `json:"didDocumentMetadata"`
	DIDResolutionMetadata DIDResolutionMetadata `json:"didResolutionMetadata"`
}

// DIDDocumentMetadata represents DID document metadata
type DIDDocumentMetadata struct {
	Updated       time.Time `json:"updated"`
	Deactivated   bool      `json:"deactivated,omitempty"`
	CanonicalID   string    `json:"canonicalId"`
	EquivalentID  []string  `json:"equivalentId,omitempty"`
	ContentHash   string    `json:"contentHash"`
	Sequence      *uint64   `json:"sequence,omitempty"`
	VersionID     *string   `json:"versionId,omitempty"`
}

// DIDResolutionMetadata represents DID resolution metadata
type DIDResolutionMetadata struct {
	ContentType string    `json:"contentType"`
	Retrieved   time.Time `json:"retrieved"`
	Resolver    string    `json:"resolver"`
	VersionID   *string   `json:"versionId,omitempty"`
}

// DataEntry represents a single data entry from Accumulate
type DataEntry struct {
	Data        []byte
	Timestamp   time.Time
	Sequence    *uint64
	ContentHash string
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

// DeterministicResolver handles deterministic DID resolution with configurable ordering
type DeterministicResolver struct {
	client acc.Client
	order  ResolveOrder
}

// NewDeterministicResolver creates a new resolver with specified ordering strategy
func NewDeterministicResolver(client acc.Client, order ResolveOrder) *DeterministicResolver {
	return &DeterministicResolver{
		client: client,
		order:  order,
	}
}

// ResolveDID resolves a DID according to the deterministic algorithm
func (r *DeterministicResolver) ResolveDID(didStr string, versionTime *time.Time) (*DIDResolutionResult, error) {
	start := time.Now()

	// Step 1: Parse DID into Accumulate URLs
	_, dataAccountURL, err := did.ParseDID(didStr)
	if err != nil {
		return nil, &InvalidDIDError{DID: didStr, Reason: err.Error()}
	}

	// Step 2: Get all data entries from the data account
	entries, err := r.getAllDataEntries(dataAccountURL, versionTime)
	if err != nil {
		return nil, &NotFoundError{DID: didStr}
	}

	if len(entries) == 0 {
		return nil, &NotFoundError{DID: didStr}
	}

	// Step 3: Apply deterministic selection algorithm
	selectedEntry, validEntries := r.selectLatestEntry(entries, didStr)
	if selectedEntry == nil {
		return nil, &NotFoundError{DID: didStr}
	}

	// Step 4: Parse the selected DID document
	var didDoc map[string]interface{}
	if err := json.Unmarshal(selectedEntry.Data, &didDoc); err != nil {
		// This shouldn't happen since we filtered for valid JSON, but handle gracefully
		return nil, fmt.Errorf("selected entry has invalid JSON: %w", err)
	}

	// Step 5: Check if deactivated
	if deactivated, exists := didDoc["deactivated"].(bool); exists && deactivated {
		// Return deactivated DID with minimal document
		return r.buildDeactivatedResult(didStr, selectedEntry, start), nil
	}

	// Step 6: Extract version ID if present
	var versionID *string
	if vid, exists := didDoc["versionId"].(string); exists && vid != "" {
		versionID = &vid
	}

	// Step 7: Build successful resolution result
	result := &DIDResolutionResult{
		DIDDocument: didDoc,
		DIDDocumentMetadata: DIDDocumentMetadata{
			Updated:      selectedEntry.Timestamp,
			Deactivated:  false,
			CanonicalID:  didStr,
			ContentHash:  selectedEntry.ContentHash,
			Sequence:     selectedEntry.Sequence,
			VersionID:    versionID,
		},
		DIDResolutionMetadata: DIDResolutionMetadata{
			ContentType: "application/did+json",
			Retrieved:   time.Now().UTC(),
			Resolver:    "accu-did-resolver",
			VersionID:   versionID,
		},
	}

	// Log resolution details
	log.Printf("DID resolved: did=%s sequence=%v timestamp=%s hash=%s deactivated=false valid_entries=%d",
		didStr, selectedEntry.Sequence, selectedEntry.Timestamp.Format(time.RFC3339),
		selectedEntry.ContentHash[:8], validEntries)

	return result, nil
}

// buildDeactivatedResult builds a 410 Gone response for deactivated DIDs
func (r *DeterministicResolver) buildDeactivatedResult(didStr string, selectedEntry *DataEntry, start time.Time) *DIDResolutionResult {
	// Parse to extract deactivatedAt if present
	var didDoc map[string]interface{}
	json.Unmarshal(selectedEntry.Data, &didDoc)

	// Minimal DID document for deactivated state
	minimalDoc := map[string]interface{}{
		"@context": []string{"https://www.w3.org/ns/did/v1"},
		"id":       didStr,
	}

	var versionID *string
	if vid, exists := didDoc["versionId"].(string); exists && vid != "" {
		versionID = &vid
	}

	result := &DIDResolutionResult{
		DIDDocument: minimalDoc,
		DIDDocumentMetadata: DIDDocumentMetadata{
			Updated:     selectedEntry.Timestamp,
			Deactivated: true,
			CanonicalID: didStr,
			ContentHash: selectedEntry.ContentHash,
			Sequence:    selectedEntry.Sequence,
			VersionID:   versionID,
		},
		DIDResolutionMetadata: DIDResolutionMetadata{
			ContentType: "application/did+json",
			Retrieved:   time.Now().UTC(),
			Resolver:    "accu-did-resolver",
			VersionID:   versionID,
		},
	}

	// Log deactivation
	log.Printf("DID resolved: did=%s sequence=%v timestamp=%s hash=%s deactivated=true",
		didStr, selectedEntry.Sequence, selectedEntry.Timestamp.Format(time.RFC3339),
		selectedEntry.ContentHash[:8])

	return result
}

// getAllDataEntries retrieves all data entries from the data account
func (r *DeterministicResolver) getAllDataEntries(dataAccountURL *url.URL, versionTime *time.Time) ([]*DataEntry, error) {
	// For now, use the simplified client interface
	// In a full implementation, this would query all entries with pagination
	data, err := r.client.GetDataAccountEntry(dataAccountURL)
	if err != nil {
		return nil, err
	}

	// Calculate content hash
	hash := sha256.Sum256(data)
	contentHash := hex.EncodeToString(hash[:])

	// Create a single entry (simplified for current client interface)
	entry := &DataEntry{
		Data:        data,
		Timestamp:   time.Now().UTC(), // Would be extracted from transaction metadata
		Sequence:    nil,              // Would be extracted from chain metadata
		ContentHash: contentHash,
	}

	return []*DataEntry{entry}, nil
}

// selectLatestEntry applies the deterministic selection algorithm
func (r *DeterministicResolver) selectLatestEntry(entries []*DataEntry, didStr string) (*DataEntry, int) {
	var validEntries []*DataEntry

	// Filter out malformed JSON entries
	for _, entry := range entries {
		var doc map[string]interface{}
		if err := json.Unmarshal(entry.Data, &doc); err != nil {
			log.Printf("WARN: Skipping malformed JSON entry for DID %s: %v", didStr, err)
			continue
		}
		validEntries = append(validEntries, entry)
	}

	if len(validEntries) == 0 {
		return nil, 0
	}

	// Sort entries according to the specified ordering strategy
	sort.Slice(validEntries, func(i, j int) bool {
		return r.compareEntries(validEntries[i], validEntries[j])
	})

	// Return the latest (last in sorted order)
	return validEntries[len(validEntries)-1], len(validEntries)
}

// compareEntries implements the deterministic comparison algorithm
// Returns true if entry i should come before entry j in sorted order
func (r *DeterministicResolver) compareEntries(i, j *DataEntry) bool {
	switch r.order {
	case ResolveOrderSequence:
		// Primary: sequence height (highest = latest)
		if i.Sequence != nil && j.Sequence != nil {
			if *i.Sequence != *j.Sequence {
				return *i.Sequence < *j.Sequence // i comes before j if i has lower sequence
			}
		} else if i.Sequence != nil && j.Sequence == nil {
			return false // i (with sequence) comes after j (without sequence)
		} else if i.Sequence == nil && j.Sequence != nil {
			return true // i (without sequence) comes before j (with sequence)
		}

		// Fallback to timestamp when sequences are equal or unavailable
		fallthrough

	case ResolveOrderTimestamp:
		// Primary (or fallback): timestamp monotonic (latest = newest)
		if !i.Timestamp.Equal(j.Timestamp) {
			return i.Timestamp.Before(j.Timestamp) // i comes before j if i is earlier
		}

		// Final tiebreaker: lexicographically greatest content hash
		return i.ContentHash < j.ContentHash // i comes before j if i has smaller hash
	}

	// Default fallback
	return i.ContentHash < j.ContentHash
}

// Legacy function for backward compatibility
func ResolveDID(client acc.Client, didStr string, versionTime *time.Time) (*DIDResolutionResult, error) {
	resolver := NewDeterministicResolver(client, ResolveOrderSequence)
	return resolver.ResolveDID(didStr, versionTime)
}