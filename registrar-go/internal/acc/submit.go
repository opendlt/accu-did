package acc

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"time"

	"gitlab.com/accumulatenetwork/accumulate/pkg/api/v3"
	"gitlab.com/accumulatenetwork/accumulate/pkg/api/v3/jsonrpc"
	"gitlab.com/accumulatenetwork/accumulate/pkg/build"
	"gitlab.com/accumulatenetwork/accumulate/pkg/types/messaging"
	"gitlab.com/accumulatenetwork/accumulate/pkg/url"
	"gitlab.com/accumulatenetwork/accumulate/protocol"

	"github.com/opendlt/accu-did/registrar-go/internal/ops"
)

// Submitter interface for Accumulate operations
type Submitter interface {
	CreateIdentity(adiLabel string, keyPageURL string) (string, error)
	CreateDataAccount(adiURL, dataAccountLabel string) (string, error)
	WriteDataEntry(dataAccountURL string, data []byte) (string, error)
	SubmitWriteData(dataAccountURL string, envelope *ops.Envelope) (string, error)
	UpdateKeyPage(keyPageURL string, operations []KeyPageOperation) (string, error)
	GetKeyPageState(keyPageURL string) (*KeyPageState, error)
}

// KeyPageOperation represents a key page operation
type KeyPageOperation struct {
	Type      string `json:"type"` // "add", "remove", "update"
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

// FakeSubmitter implements Submitter interface for testing and development
type FakeSubmitter struct {
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

// NewSubmitter creates a new submitter based on mode
func NewSubmitter(realMode bool, nodeURL string) Submitter {
	if realMode {
		return NewRealSubmitter(nodeURL)
	}
	return NewFakeSubmitter()
}

// NewFakeSubmitter creates a new fake submitter
func NewFakeSubmitter() *FakeSubmitter {
	return &FakeSubmitter{
		transactions: make(map[string]*MockTransaction),
		keyPages:     make(map[string]*KeyPageState),
	}
}

// CreateIdentity creates a new ADI (fake implementation)
func (c *FakeSubmitter) CreateIdentity(adiLabel string, keyPageURL string) (string, error) {
	txID := c.generateTxID()

	transaction := &MockTransaction{
		ID:        txID,
		Status:    "committed",
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"type":       "createIdentity",
			"adiLabel":   adiLabel,
			"keyPageURL": keyPageURL,
		},
	}

	c.transactions[txID] = transaction
	return txID, nil
}

// CreateDataAccount creates a new data account (fake implementation)
func (c *FakeSubmitter) CreateDataAccount(adiURL, dataAccountLabel string) (string, error) {
	txID := c.generateTxID()

	transaction := &MockTransaction{
		ID:        txID,
		Status:    "committed",
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"type":             "createDataAccount",
			"adiURL":           adiURL,
			"dataAccountLabel": dataAccountLabel,
		},
	}

	c.transactions[txID] = transaction
	return txID, nil
}

// WriteDataEntry writes data to a data account (fake implementation)
func (c *FakeSubmitter) WriteDataEntry(dataAccountURL string, data []byte) (string, error) {
	txID := c.generateTxID()

	transaction := &MockTransaction{
		ID:        txID,
		Status:    "committed",
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"type":           "writeData",
			"dataAccountURL": dataAccountURL,
			"data":           string(data),
		},
	}

	c.transactions[txID] = transaction
	return txID, nil
}

// SubmitWriteData submits a writeData transaction to Accumulate (fake implementation)
func (c *FakeSubmitter) SubmitWriteData(dataAccountURL string, envelope *ops.Envelope) (string, error) {
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

// UpdateKeyPage updates a key page (fake implementation)
func (c *FakeSubmitter) UpdateKeyPage(keyPageURL string, operations []KeyPageOperation) (string, error) {
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
func (c *FakeSubmitter) GetKeyPageState(keyPageURL string) (*KeyPageState, error) {
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
func (c *FakeSubmitter) GetTransaction(txID string) (*MockTransaction, error) {
	if tx, exists := c.transactions[txID]; exists {
		return tx, nil
	}
	return nil, fmt.Errorf("transaction not found: %s", txID)
}

// generateTxID generates a mock transaction ID
func (c *FakeSubmitter) generateTxID() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("0x%016x", timestamp)
}

// ListTransactions returns all transactions (for testing)
func (c *FakeSubmitter) ListTransactions() map[string]*MockTransaction {
	return c.transactions
}

// RealSubmitter implements Submitter interface using JSON-RPC v3
type RealSubmitter struct {
	client     *jsonrpc.Client
	signerHook SignerHook
}

// SignerHook provides cryptographic signing capability
type SignerHook interface {
	// Sign signs the given message with the specified private key
	Sign(privateKey []byte, message []byte) ([]byte, error)

	// GetPrivateKey retrieves the private key for a given key page URL
	GetPrivateKey(keyPageURL string) ([]byte, error)

	// GetPublicKey retrieves the public key for a given key page URL
	GetPublicKey(keyPageURL string) ([]byte, error)
}

// NewRealSubmitter creates a new real submitter that connects to Accumulate network
func NewRealSubmitter(nodeURL string) *RealSubmitter {
	return &RealSubmitter{
		client:     jsonrpc.NewClient(nodeURL),
		signerHook: NewDefaultSignerHook(),
	}
}

// NewRealSubmitterWithSigner creates a new real submitter with custom signer
func NewRealSubmitterWithSigner(nodeURL string, signerHook SignerHook) *RealSubmitter {
	return &RealSubmitter{
		client:     jsonrpc.NewClient(nodeURL),
		signerHook: signerHook,
	}
}

// DefaultSignerHook implements SignerHook with in-memory key management
type DefaultSignerHook struct {
	keys map[string]ed25519.PrivateKey
}

// NewDefaultSignerHook creates a new default signer hook
func NewDefaultSignerHook() *DefaultSignerHook {
	return &DefaultSignerHook{
		keys: make(map[string]ed25519.PrivateKey),
	}
}

// RegisterKey registers a private key for a given key page URL
func (s *DefaultSignerHook) RegisterKey(keyPageURL string, privateKey ed25519.PrivateKey) {
	s.keys[keyPageURL] = privateKey
}

// GenerateKey generates and registers a new Ed25519 key pair for a key page URL
func (s *DefaultSignerHook) GenerateKey(keyPageURL string) (ed25519.PublicKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	s.keys[keyPageURL] = privateKey
	return publicKey, nil
}

// Sign signs the given message with the specified private key
func (s *DefaultSignerHook) Sign(privateKey []byte, message []byte) ([]byte, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key size: expected %d, got %d", ed25519.PrivateKeySize, len(privateKey))
	}

	signature := ed25519.Sign(privateKey, message)
	return signature, nil
}

// GetPrivateKey retrieves the private key for a given key page URL
func (s *DefaultSignerHook) GetPrivateKey(keyPageURL string) ([]byte, error) {
	privateKey, exists := s.keys[keyPageURL]
	if !exists {
		return nil, fmt.Errorf("private key not found for key page: %s", keyPageURL)
	}

	return privateKey, nil
}

// GetPublicKey retrieves the public key for a given key page URL
func (s *DefaultSignerHook) GetPublicKey(keyPageURL string) ([]byte, error) {
	privateKey, exists := s.keys[keyPageURL]
	if !exists {
		return nil, fmt.Errorf("private key not found for key page: %s", keyPageURL)
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)
	return publicKey, nil
}

// CreateIdentity creates a new ADI using Accumulate API
func (c *RealSubmitter) CreateIdentity(adiLabel string, keyPageURL string) (string, error) {
	// Parse URLs
	keyPageParsed, err := url.Parse(keyPageURL)
	if err != nil {
		return "", fmt.Errorf("invalid key page URL %s: %w", keyPageURL, err)
	}

	adiURL := fmt.Sprintf("acc://%s", adiLabel)
	adiParsed, err := url.Parse(adiURL)
	if err != nil {
		return "", fmt.Errorf("invalid ADI URL %s: %w", adiURL, err)
	}

	// Get or generate keys for this key page
	privateKey, err := c.signerHook.GetPrivateKey(keyPageURL)
	if err != nil {
		// Generate new key if not found
		if defaultSigner, ok := c.signerHook.(*DefaultSignerHook); ok {
			publicKey, genErr := defaultSigner.GenerateKey(keyPageURL)
			if genErr != nil {
				return "", fmt.Errorf("failed to generate key for %s: %w", keyPageURL, genErr)
			}
			privateKey, err = c.signerHook.GetPrivateKey(keyPageURL)
			if err != nil {
				return "", fmt.Errorf("failed to retrieve generated key: %w", err)
			}

			// Log the public key for reference
			fmt.Printf("Generated new key for %s: %x\n", keyPageURL, publicKey)
		} else {
			return "", fmt.Errorf("key not found for %s and cannot generate with custom signer", keyPageURL)
		}
	}

	publicKey, err := c.signerHook.GetPublicKey(keyPageURL)
	if err != nil {
		return "", fmt.Errorf("failed to get public key for %s: %w", keyPageURL, err)
	}

	// Build CreateIdentity transaction using pkg/build
	envelope, err := build.Transaction().
		For(adiParsed).
		CreateIdentity(adiURL).
		WithKeyBook(keyPageURL).
		WithKey(publicKey, protocol.SignatureTypeED25519).
		SignWith(keyPageParsed).
		Version(1).
		Timestamp(build.UnixTimeNow()).
		PrivateKey(privateKey).
		Done()
	if err != nil {
		return "", fmt.Errorf("failed to build CreateIdentity transaction: %w", err)
	}

	// Submit to network
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	submissions, err := c.client.Submit(ctx, envelope, api.SubmitOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to submit CreateIdentity transaction: %w", err)
	}

	if len(submissions) == 0 {
		return "", fmt.Errorf("no submissions returned")
	}

	if !submissions[0].Success {
		return "", fmt.Errorf("CreateIdentity submission failed: %v", submissions[0].Message)
	}

	// Extract transaction ID from envelope
	txID := extractTxID(envelope)
	return txID, nil
}

// CreateDataAccount creates a new data account using Accumulate API
func (c *RealSubmitter) CreateDataAccount(adiURL, dataAccountLabel string) (string, error) {
	// Parse the ADI URL
	adiParsed, err := url.Parse(adiURL)
	if err != nil {
		return "", fmt.Errorf("invalid ADI URL %s: %w", adiURL, err)
	}

	// Construct data account URL
	dataAccountURL := fmt.Sprintf("%s/%s", adiURL, dataAccountLabel)

	// Construct key page URL (assumes book/1 pattern)
	keyPageURL := fmt.Sprintf("%s/book/1", adiURL)
	keyPageParsed, err := url.Parse(keyPageURL)
	if err != nil {
		return "", fmt.Errorf("invalid key page URL %s: %w", keyPageURL, err)
	}

	// Get private key for signing
	privateKey, err := c.signerHook.GetPrivateKey(keyPageURL)
	if err != nil {
		return "", fmt.Errorf("failed to get private key for %s: %w", keyPageURL, err)
	}

	// Build CreateDataAccount transaction using pkg/build
	envelope, err := build.Transaction().
		For(adiParsed).
		CreateDataAccount(dataAccountURL).
		SignWith(keyPageParsed).
		Version(1).
		Timestamp(build.UnixTimeNow()).
		PrivateKey(privateKey).
		Done()
	if err != nil {
		return "", fmt.Errorf("failed to build CreateDataAccount transaction: %w", err)
	}

	// Submit to network
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	submissions, err := c.client.Submit(ctx, envelope, api.SubmitOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to submit CreateDataAccount transaction: %w", err)
	}

	if len(submissions) == 0 {
		return "", fmt.Errorf("no submissions returned")
	}

	if !submissions[0].Success {
		return "", fmt.Errorf("CreateDataAccount submission failed: %v", submissions[0].Message)
	}

	// Extract transaction ID from envelope
	txID := extractTxID(envelope)
	return txID, nil
}

// WriteDataEntry writes data to a data account using Accumulate API
func (c *RealSubmitter) WriteDataEntry(dataAccountURL string, data []byte) (string, error) {
	// Parse the data account URL
	dataAccountParsed, err := url.Parse(dataAccountURL)
	if err != nil {
		return "", fmt.Errorf("invalid data account URL %s: %w", dataAccountURL, err)
	}

	// Extract ADI from data account URL to construct key page URL
	// For acc://alice/did, ADI is acc://alice
	adiURL := fmt.Sprintf("acc://%s", dataAccountParsed.Authority)
	keyPageURL := fmt.Sprintf("%s/book/1", adiURL)
	keyPageParsed, err := url.Parse(keyPageURL)
	if err != nil {
		return "", fmt.Errorf("invalid key page URL %s: %w", keyPageURL, err)
	}

	// Get private key for signing
	privateKey, err := c.signerHook.GetPrivateKey(keyPageURL)
	if err != nil {
		return "", fmt.Errorf("failed to get private key for %s: %w", keyPageURL, err)
	}

	// Build WriteData transaction using pkg/build
	envelope, err := build.Transaction().
		For(dataAccountParsed).
		WriteData(data).
		SignWith(keyPageParsed).
		Version(1).
		Timestamp(build.UnixTimeNow()).
		PrivateKey(privateKey).
		Done()
	if err != nil {
		return "", fmt.Errorf("failed to build WriteData transaction: %w", err)
	}

	// Submit to network
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	submissions, err := c.client.Submit(ctx, envelope, api.SubmitOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to submit WriteData transaction: %w", err)
	}

	if len(submissions) == 0 {
		return "", fmt.Errorf("no submissions returned")
	}

	if !submissions[0].Success {
		return "", fmt.Errorf("WriteData submission failed: %v", submissions[0].Message)
	}

	// Extract transaction ID from envelope
	txID := extractTxID(envelope)
	return txID, nil
}

// SubmitWriteData submits a writeData transaction to Accumulate
func (c *RealSubmitter) SubmitWriteData(dataAccountURL string, envelope *ops.Envelope) (string, error) {
	// Parse the data account URL
	accountURL, err := url.Parse(dataAccountURL)
	if err != nil {
		return "", fmt.Errorf("invalid data account URL %s: %w", dataAccountURL, err)
	}

	// Convert ops.Envelope to messaging.Envelope
	msgEnvelope, err := c.convertToMessagingEnvelope(envelope, accountURL)
	if err != nil {
		return "", fmt.Errorf("failed to convert envelope: %w", err)
	}

	// Submit the envelope
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	submissions, err := c.client.Submit(ctx, msgEnvelope, api.SubmitOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction: %w", err)
	}

	// Return the transaction ID from the first submission
	if len(submissions) == 0 {
		return "", fmt.Errorf("no submissions returned")
	}

	// Extract transaction ID from submission
	// Note: The actual field name may vary based on the API structure
	// Generate a placeholder transaction ID since Success is boolean
	// In a real implementation, this should use the actual transaction hash from submissions
	if !submissions[0].Success {
		return "", fmt.Errorf("submission failed")
	}
	txID := fmt.Sprintf("0x%x", time.Now().UnixNano()) // Temporary solution

	return txID, nil
}

// UpdateKeyPage updates a key page
func (c *RealSubmitter) UpdateKeyPage(keyPageURL string, operations []KeyPageOperation) (string, error) {
	// Parse the key page URL
	keyPageParsed, err := url.Parse(keyPageURL)
	if err != nil {
		return "", fmt.Errorf("invalid key page URL %s: %w", keyPageURL, err)
	}

	// Get private key for signing
	privateKey, err := c.signerHook.GetPrivateKey(keyPageURL)
	if err != nil {
		return "", fmt.Errorf("failed to get private key for %s: %w", keyPageURL, err)
	}

	// Build UpdateKeyPage transaction using pkg/build
	// For now, we'll use a simple approach - in a full implementation,
	// we would need to convert operations to proper Accumulate key page operations
	builder := build.Transaction().
		For(keyPageParsed).
		SignWith(keyPageParsed).
		Version(1).
		Timestamp(build.UnixTimeNow()).
		PrivateKey(privateKey)

	// Apply operations (simplified - real implementation would use proper UpdateKeyPage transaction)
	// For now, we'll create a simple transaction that represents the key page update
	envelope, err := builder.Done()
	if err != nil {
		return "", fmt.Errorf("failed to build UpdateKeyPage transaction: %w", err)
	}

	// Submit to network
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	submissions, err := c.client.Submit(ctx, envelope, api.SubmitOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to submit UpdateKeyPage transaction: %w", err)
	}

	if len(submissions) == 0 {
		return "", fmt.Errorf("no submissions returned")
	}

	if !submissions[0].Success {
		return "", fmt.Errorf("UpdateKeyPage submission failed: %v", submissions[0].Message)
	}

	// Extract transaction ID from envelope
	txID := extractTxID(envelope)
	return txID, nil
}

// GetKeyPageState returns the current state of a key page
func (c *RealSubmitter) GetKeyPageState(keyPageURL string) (*KeyPageState, error) {
	// Parse the key page URL
	pageURL, err := url.Parse(keyPageURL)
	if err != nil {
		return nil, fmt.Errorf("invalid key page URL %s: %w", keyPageURL, err)
	}

	// Query the key page state using QueryAccount for proper AccountRecord return
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create a querier wrapper around the client
	querier := api.Querier2{Querier: c.client}

	// Query the account (key page)
	accountRecord, err := querier.QueryAccount(ctx, pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query key page %s: %w", keyPageURL, err)
	}

	// Convert AccountRecord to KeyPageState
	keyPageState := &KeyPageState{
		URL:       keyPageURL,
		Threshold: 1,           // Default threshold
		Keys:      []KeyInfo{}, // Will be populated below
		Height:    0,           // Will be populated below
	}

	// Extract information from AccountRecord
	if accountRecord != nil && accountRecord.Account != nil {
		// For key pages, the account should be a KeyPage type
		// This is a simplified conversion - actual implementation would need to
		// handle the specific protocol.KeyPage type
		keyPageState.Height = 1 // Placeholder - would extract from chain info

		// Note: In a complete implementation, we would:
		// 1. Cast accountRecord.Account to *protocol.KeyPage
		// 2. Extract the threshold and keys from the KeyPage
		// 3. Convert to our KeyInfo format
	}

	return keyPageState, nil
}

// convertToMessagingEnvelope converts ops.Envelope to messaging.Envelope
func (c *RealSubmitter) convertToMessagingEnvelope(opsEnv *ops.Envelope, accountURL *url.URL) (*messaging.Envelope, error) {
	// This is a placeholder implementation
	// The actual conversion would depend on the structure of ops.Envelope and
	// how it maps to the messaging.Envelope type

	// TODO: Implement proper conversion based on:
	// 1. ops.Envelope structure
	// 2. messaging.Envelope requirements
	// 3. The specific transaction type (WriteData)

	envelope := &messaging.Envelope{
		// Transaction: writeDataTx,
		// Signatures: signatures,
	}

	return envelope, nil
}

// extractTxID extracts transaction ID from envelope
func extractTxID(envelope *messaging.Envelope) string {
	if len(envelope.Transaction) > 0 {
		return fmt.Sprintf("%x", envelope.Transaction[0].GetHash())
	}
	return "unknown"
}
