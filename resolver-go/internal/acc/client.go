package acc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/accumulatenetwork/accumulate/pkg/api/v3"
	"gitlab.com/accumulatenetwork/accumulate/pkg/api/v3/jsonrpc"
	"gitlab.com/accumulatenetwork/accumulate/pkg/url"
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

// NewClient creates a new client based on mode
func NewClient(realMode bool, nodeURL string) Client {
	if realMode {
		return NewRealClient(nodeURL)
	}
	return NewFakeClient("testdata")
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

// RealClient implements Client interface using JSON-RPC v3
type RealClient struct {
	client *jsonrpc.Client
}

// NewRealClient creates a new real client that connects to Accumulate network
func NewRealClient(nodeURL string) *RealClient {
	return &RealClient{
		client: jsonrpc.NewClient(nodeURL),
	}
}

// GetLatestDIDEntry returns the latest DID entry for an ADI
func (c *RealClient) GetLatestDIDEntry(adi string) (Envelope, error) {
	// Build the account URL for the ADI
	accountURL, err := url.Parse(fmt.Sprintf("acc://%s", adi))
	if err != nil {
		return Envelope{}, fmt.Errorf("invalid ADI %s: %w", adi, err)
	}

	// Query the account to get the latest DID entry
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	record, err := c.client.Query(ctx, accountURL, nil)
	if err != nil {
		return Envelope{}, fmt.Errorf("failed to query account %s: %w", adi, err)
	}

	// Convert record to envelope format
	return c.recordToEnvelope(record)
}

// GetEntryAtTime returns a DID entry at a specific time
func (c *RealClient) GetEntryAtTime(adi string, t time.Time) (Envelope, error) {
	// Build the account URL for the ADI
	accountURL, err := url.Parse(fmt.Sprintf("acc://%s", adi))
	if err != nil {
		return Envelope{}, fmt.Errorf("invalid ADI %s: %w", adi, err)
	}

	// Create query with time constraint
	query := &api.DefaultQuery{
		// Add time-based query parameters if available in API
	}

	// Query the account at the specific time
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	record, err := c.client.Query(ctx, accountURL, query)
	if err != nil {
		return Envelope{}, fmt.Errorf("failed to query account %s at time %v: %w", adi, t, err)
	}

	// Convert record to envelope format
	return c.recordToEnvelope(record)
}

// GetKeyPageState returns the state of a Key Page
func (c *RealClient) GetKeyPageState(keyPageURLStr string) (KeyPageState, error) {
	// Parse the key page URL
	keyPageURL, err := url.Parse(keyPageURLStr)
	if err != nil {
		return KeyPageState{}, fmt.Errorf("invalid key page URL %s: %w", keyPageURLStr, err)
	}

	// Query the key page state
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = c.client.Query(ctx, keyPageURL, nil)
	if err != nil {
		return KeyPageState{}, fmt.Errorf("failed to query key page %s: %w", keyPageURLStr, err)
	}

	// Convert record to KeyPageState
	// Note: This conversion will depend on the actual record type returned by the API
	keyPageState := KeyPageState{
		URL:       keyPageURLStr,
		Threshold: 1, // Default value, should be extracted from record
		Keys:      []Key{},
	}

	// TODO: Extract actual key page data from record
	// This will depend on the specific record type returned by the API

	return keyPageState, nil
}

// recordToEnvelope converts an API record to our Envelope format
func (c *RealClient) recordToEnvelope(record api.Record) (Envelope, error) {
	// This is a simplified conversion. In practice, we would need to:
	// 1. Extract the DID document from the record
	// 2. Extract metadata like version ID, timestamp, etc.
	// 3. Build the proper envelope structure

	// For now, return a basic envelope structure
	// TODO: Implement proper record to envelope conversion based on actual API types
	envelope := Envelope{
		ContentType: "application/did+json",
		Document:    make(map[string]interface{}),
		Meta: EnvelopeMeta{
			VersionID:     fmt.Sprintf("%d", time.Now().Unix()),
			Timestamp:     time.Now().UTC(),
			AuthorKeyPage: "", // Extract from record
			Proof: Proof{
				TxID:        "", // Extract from record
				ContentHash: "",
			},
		},
	}

	return envelope, nil
}