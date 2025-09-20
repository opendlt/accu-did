package acc

import (
	"fmt"
	"time"

	"github.com/opendlt/accu-did/registrar-go/internal/ops"
)

// Client interface for Accumulate operations
type Client interface {
	SubmitWriteData(dataAccountURL string, envelope *ops.Envelope) (string, error)
	UpdateKeyPage(keyPageURL string, operations []KeyPageOperation) (string, error)
	GetKeyPageState(keyPageURL string) (*KeyPageState, error)
}

// KeyPageOperation represents a key page operation
type KeyPageOperation struct {
	Type      string `json:"type"`      // "add", "remove", "update"
	PublicKey string `json:"publicKey"`
	KeyType   string `json:"keyType"`
}

// KeyPageState represents the current state of a key page
type KeyPageState struct {
	URL       string    `json:"url"`
	Threshold int       `json:"threshold"`
	Keys      []KeyInfo `json:"keys"`
	Height    uint64    `json:"height"`
}

// KeyInfo represents a key in a key page
type KeyInfo struct {
	PublicKey string `json:"publicKey"`
	KeyType   string `json:"keyType"`
}

// MockClient implements Client interface for testing and development
type MockClient struct {
	transactions map[string]*MockTransaction
	keyPages     map[string]*KeyPageState
}

// MockTransaction represents a mock transaction
type MockTransaction struct {
	ID        string
	Status    string
	Timestamp time.Time
	Data      interface{}
}

// NewMockClient creates a new mock client
func NewMockClient() *MockClient {
	return &MockClient{
		transactions: make(map[string]*MockTransaction),
		keyPages:     make(map[string]*KeyPageState),
	}
}

// SubmitWriteData submits a writeData transaction to Accumulate (mock implementation)
func (c *MockClient) SubmitWriteData(dataAccountURL string, envelope *ops.Envelope) (string, error) {
	// Generate mock transaction ID
	txID := c.generateTxID()

	// Create mock transaction
	transaction := &MockTransaction{
		ID:        txID,
		Status:    "committed",
		Timestamp: time.Now().UTC(),
		Data:      envelope,
	}

	// Store transaction
	c.transactions[txID] = transaction

	// Set transaction ID in envelope
	envelope.SetTransactionID(txID)

	return txID, nil
}

// UpdateKeyPage updates a key page (mock implementation)
func (c *MockClient) UpdateKeyPage(keyPageURL string, operations []KeyPageOperation) (string, error) {
	// Generate mock transaction ID
	txID := c.generateTxID()

	// Create mock transaction
	transaction := &MockTransaction{
		ID:        txID,
		Status:    "committed",
		Timestamp: time.Now().UTC(),
		Data:      operations,
	}

	// Store transaction
	c.transactions[txID] = transaction

	// Mock: apply operations to key page state
	if _, exists := c.keyPages[keyPageURL]; !exists {
		c.keyPages[keyPageURL] = &KeyPageState{
			URL:       keyPageURL,
			Threshold: 1,
			Keys:      []KeyInfo{},
			Height:    0,
		}
	}

	keyPage := c.keyPages[keyPageURL]
	keyPage.Height++

	// Apply operations (simplified)
	for _, op := range operations {
		switch op.Type {
		case "add":
			keyPage.Keys = append(keyPage.Keys, KeyInfo{
				PublicKey: op.PublicKey,
				KeyType:   op.KeyType,
			})
		case "remove":
			// Remove key by public key
			for i, key := range keyPage.Keys {
				if key.PublicKey == op.PublicKey {
					keyPage.Keys = append(keyPage.Keys[:i], keyPage.Keys[i+1:]...)
					break
				}
			}
		}
	}

	return txID, nil
}

// GetKeyPageState returns the current state of a key page
func (c *MockClient) GetKeyPageState(keyPageURL string) (*KeyPageState, error) {
	if keyPage, exists := c.keyPages[keyPageURL]; exists {
		return keyPage, nil
	}

	// Return default key page for testing
	return &KeyPageState{
		URL:       keyPageURL,
		Threshold: 1,
		Keys: []KeyInfo{
			{
				PublicKey: "ed25519:mockkey123...",
				KeyType:   "ed25519",
			},
		},
		Height: 1,
	}, nil
}

// GetTransaction returns a transaction by ID (for testing)
func (c *MockClient) GetTransaction(txID string) (*MockTransaction, error) {
	if tx, exists := c.transactions[txID]; exists {
		return tx, nil
	}
	return nil, fmt.Errorf("transaction not found: %s", txID)
}

// generateTxID generates a mock transaction ID
func (c *MockClient) generateTxID() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("0x%016x", timestamp)
}

// ListTransactions returns all transactions (for testing)
func (c *MockClient) ListTransactions() map[string]*MockTransaction {
	return c.transactions
}