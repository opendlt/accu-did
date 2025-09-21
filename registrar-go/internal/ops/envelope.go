package ops

import (
	"fmt"
	"time"

	"github.com/opendlt/accu-did/registrar-go/internal/canon"
)

// Envelope represents a DID document entry envelope
type Envelope struct {
	ContentType string                 `json:"contentType"`
	Document    map[string]interface{} `json:"document"`
	Meta        EnvelopeMeta           `json:"meta"`
}

// EnvelopeMeta represents envelope metadata
type EnvelopeMeta struct {
	VersionID         string    `json:"versionId"`
	PreviousVersionID string    `json:"previousVersionId,omitempty"`
	Timestamp         time.Time `json:"timestamp"`
	AuthorKeyPage     string    `json:"authorKeyPage"`
	Proof             Proof     `json:"proof"`
}

// Proof represents envelope proof data
type Proof struct {
	Type        string `json:"type,omitempty"`
	TxID        string `json:"txid"`
	ContentHash string `json:"contentHash"`
}

// BuildEnvelope creates a new envelope for a DID document
func BuildEnvelope(document map[string]interface{}, authorKeyPage string, previousVersionID string) (*Envelope, error) {
	// Generate version ID
	timestamp := time.Now().UTC()
	versionID := generateVersionID(timestamp)

	// Canonicalize document and compute content hash
	canonical, err := canon.Canonicalize(document)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize document: %w", err)
	}

	contentHash := canon.SHA256(canonical)

	// Create envelope
	envelope := &Envelope{
		ContentType: "application/did+json",
		Document:    document,
		Meta: EnvelopeMeta{
			VersionID:         versionID,
			PreviousVersionID: previousVersionID,
			Timestamp:         timestamp,
			AuthorKeyPage:     authorKeyPage,
			Proof: Proof{
				Type:        "accumulate",
				ContentHash: contentHash,
				// TxID will be set after submission
			},
		},
	}

	return envelope, nil
}

// SetTransactionID sets the transaction ID in the envelope proof
func (e *Envelope) SetTransactionID(txID string) {
	e.Meta.Proof.TxID = txID
}

// GetContentHash returns the content hash from the envelope
func (e *Envelope) GetContentHash() string {
	return e.Meta.Proof.ContentHash
}

// ValidateContentHash verifies that the content hash matches the document
func (e *Envelope) ValidateContentHash() error {
	canonical, err := canon.Canonicalize(e.Document)
	if err != nil {
		return fmt.Errorf("failed to canonicalize document: %w", err)
	}

	expectedHash := canon.SHA256(canonical)
	if expectedHash != e.Meta.Proof.ContentHash {
		return fmt.Errorf("content hash mismatch: expected %s, got %s", expectedHash, e.Meta.Proof.ContentHash)
	}

	return nil
}

// generateVersionID creates a unique version ID combining timestamp and hash prefix
func generateVersionID(timestamp time.Time) string {
	// Format: <unix-timestamp>-<random-suffix>
	unix := timestamp.Unix()

	// Use last 8 characters of timestamp as suffix for uniqueness
	suffix := fmt.Sprintf("%08x", unix)[0:8]

	return fmt.Sprintf("%d-%s", unix, suffix)
}
