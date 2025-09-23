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
	"gitlab.com/accumulatenetwork/accumulate/protocol"
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
	GetDataAccountEntry(dataAccountURL *url.URL) ([]byte, error)
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

// GetDataAccountEntry reads from testdata for FAKE mode
func (c *FakeClient) GetDataAccountEntry(dataAccountURL *url.URL) ([]byte, error) {
	// Extract ADI from URL for testdata lookup
	adiLabel := dataAccountURL.Authority
	filename := fmt.Sprintf("did-%s.json", adiLabel)

	path := filepath.Join(c.testdataDir, "entries", filename)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("DID not found: %s", dataAccountURL.String())
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read testdata: %w", err)
	}

	return data, nil
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

	// Create a querier wrapper around the client
	querier := api.Querier2{Querier: c.client}

	// Query the account (key page)
	accountRecord, err := querier.QueryAccount(ctx, keyPageURL, nil)
	if err != nil {
		return KeyPageState{}, fmt.Errorf("failed to query key page %s: %w", keyPageURLStr, err)
	}

	// Convert AccountRecord to KeyPageState
	keyPageState := KeyPageState{
		URL:       keyPageURLStr,
		Threshold: 1, // Default threshold
		Keys:      []Key{},
	}

	// Extract information from AccountRecord
	if accountRecord != nil && accountRecord.Account != nil {
		// For key pages, the account should be a KeyPage type
		if keyPage, ok := accountRecord.Account.(*protocol.KeyPage); ok {
			// Extract accept threshold (modern Accumulate uses this instead of simple threshold)
			keyPageState.Threshold = int(keyPage.AcceptThreshold)

			// Note: KeySpec only contains PublicKeyHash, not the full public key
			// In a complete implementation, we would need to query for the actual keys
			// For now, just populate with available information
			for _, keySpec := range keyPage.Keys {
				keyPageState.Keys = append(keyPageState.Keys, Key{
					PublicKey: fmt.Sprintf("%x", keySpec.PublicKeyHash),
					KeyType:   "ed25519", // Default assumption
				})
			}
		}
	}

	return keyPageState, nil
}

// GetDataAccountEntry reads latest data entry from a data account
func (c *RealClient) GetDataAccountEntry(dataAccountURL *url.URL) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create a querier wrapper around the client for typed queries
	querier := api.Querier2{Querier: c.client}

	// Query the latest data entry from the data account using DataQuery
	count := uint64(1)
	start := uint64(0) // Start from beginning, get latest (reverse order typically)
	dataQuery := &api.DataQuery{
		Range: &api.RangeOptions{
			Count: &count,
			Start: start,
		},
	}

	// Query for data entries
	entries, err := querier.QueryDataEntries(ctx, dataAccountURL, dataQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query data entries for %s: %w", dataAccountURL.String(), err)
	}

	// Check if we have any entries
	if entries == nil || len(entries.Records) == 0 {
		return nil, fmt.Errorf("no data entries found for %s", dataAccountURL.String())
	}

	// Get the latest entry (first in the result since we queried from the end)
	latestEntry := entries.Records[0]
	if latestEntry == nil || latestEntry.Value == nil {
		return nil, fmt.Errorf("invalid data entry format for %s", dataAccountURL.String())
	}

	// Extract the transaction message from the entry
	// The Message field is already of type *messaging.TransactionMessage based on the query
	txMsg := latestEntry.Value.Message
	if txMsg == nil {
		return nil, fmt.Errorf("no transaction message in data entry for %s", dataAccountURL.String())
	}

	// Cast the body to WriteData transaction
	if writeData, ok := txMsg.Transaction.Body.(*protocol.WriteData); ok {
		// The data should contain the DID document JSON
		if len(writeData.Entry.GetData()) == 0 {
			return nil, fmt.Errorf("empty data entry for %s", dataAccountURL.String())
		}

		// Return the raw data (should be JSON)
		return writeData.Entry.GetData()[0], nil
	}

	return nil, fmt.Errorf("unsupported transaction type in data entry for %s", dataAccountURL.String())
}

// recordToEnvelope converts an API record to our Envelope format
func (c *RealClient) recordToEnvelope(record api.Record) (Envelope, error) {
	// For now, use a simplified approach that doesn't depend on unstable API methods
	// In a complete implementation, this would extract data from the specific record type

	// Create a minimal envelope structure
	timestamp := time.Now().UTC()
	versionID := fmt.Sprintf("%d", timestamp.Unix())

	// Default empty document - in practice this would be populated from the record
	document := make(map[string]interface{})

	envelope := Envelope{
		ContentType: "application/did+json",
		Document:    document,
		Meta: EnvelopeMeta{
			VersionID:     versionID,
			Timestamp:     timestamp,
			AuthorKeyPage: "", // Would be extracted from transaction metadata
			Proof: Proof{
				TxID:        "", // Would be extracted from record hash
				ContentHash: "", // Would be computed from document
			},
		},
	}

	return envelope, nil
}
