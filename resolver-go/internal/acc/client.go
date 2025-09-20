package acc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
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

// KeyPageState represents the state of an Accumulate Key Page
type KeyPageState struct {
	URL       string `json:"url"`
	Threshold int    `json:"threshold"`
	Keys      []Key  `json:"keys"`
}

// Key represents a key in a Key Page
type Key struct {
	PublicKey string `json:"publicKey"`
	KeyType   string `json:"keyType"`
}

// Client interface for Accumulate operations
type Client interface {
	GetLatestDIDEntry(adi string) (Envelope, error)
	GetEntryAtTime(adi string, t time.Time) (Envelope, error)
	GetKeyPageState(url string) (KeyPageState, error)
}

// FakeClient implements Client interface using golden files
type FakeClient struct {
	testdataDir string
}

// NewFakeClient creates a new fake client that reads from testdata
func NewFakeClient(testdataDir string) *FakeClient {
	return &FakeClient{
		testdataDir: testdataDir,
	}
}

// GetLatestDIDEntry returns the latest DID entry for an ADI
func (c *FakeClient) GetLatestDIDEntry(adi string) (Envelope, error) {
	// For testing, return the update.service version as "latest"
	return c.loadEnvelope("entry.update.service.json")
}

// GetEntryAtTime returns a DID entry at a specific time
func (c *FakeClient) GetEntryAtTime(adi string, t time.Time) (Envelope, error) {
	// Simple logic: before 2024-01-02 returns v1, after returns v2
	cutoff := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	if t.Before(cutoff) {
		return c.loadEnvelope("entry.v1.json")
	}

	return c.loadEnvelope("entry.update.service.json")
}

// GetKeyPageState returns the state of a Key Page
func (c *FakeClient) GetKeyPageState(url string) (KeyPageState, error) {
	// Return a mock key page state
	return KeyPageState{
		URL:       url,
		Threshold: 1,
		Keys: []Key{
			{
				PublicKey: "ed25519:abc123...",
				KeyType:   "ed25519",
			},
		},
	}, nil
}

// loadEnvelope loads an envelope from testdata
func (c *FakeClient) loadEnvelope(filename string) (Envelope, error) {
	// Try examples directory first
	path := filepath.Join(c.testdataDir, "examples", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return Envelope{}, fmt.Errorf("failed to read %s: %w", path, err)
	}

	// Parse as DID document and wrap in envelope
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return Envelope{}, fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	// Create envelope wrapper
	envelope := Envelope{
		ContentType: "application/did+json",
		Document:    doc,
		Meta: EnvelopeMeta{
			VersionID:     "1704067200-8b4c4f7b",
			Timestamp:     time.Now().UTC(),
			AuthorKeyPage: "acc://alice/book/1",
			Proof: Proof{
				TxID:        "0x1234567890abcdef",
				ContentHash: "", // Will be computed if needed
			},
		},
	}

	return envelope, nil
}