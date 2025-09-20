package resolve

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opendlt/accu-did/resolver-go/internal/acc"
	"github.com/opendlt/accu-did/resolver-go/internal/normalize"
)

func TestNormalizeDID(t *testing.T) {
	// Load URL normalization test vectors
	vectorsPath := filepath.Join("..", "..", "testdata", "vectors", "url-normalization.json")
	data, err := os.ReadFile(vectorsPath)
	require.NoError(t, err)

	var vectors struct {
		Vectors []struct {
			Name     string `json:"name"`
			Input    string `json:"input"`
			Expected string `json:"expected"`
		} `json:"vectors"`
	}
	require.NoError(t, json.Unmarshal(data, &vectors))

	for _, vector := range vectors.Vectors {
		t.Run(vector.Name, func(t *testing.T) {
			normalized, _, err := normalize.NormalizeDID(vector.Input)
			require.NoError(t, err)
			assert.Equal(t, vector.Expected, normalized)
		})
	}
}

func TestResolveDID_Latest(t *testing.T) {
	client := acc.NewFakeClient("../../testdata")

	result, err := ResolveDID(client, "did:acc:alice", nil)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Should return the latest version (entry.update.service.json)
	doc := result.DIDDocument.(map[string]interface{})
	assert.Equal(t, "did:acc:alice", doc["id"])

	// Should have services (indicating it's the updated version)
	services, ok := doc["service"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, services, 4) // Updated version has 4 services
}

func TestResolveDID_VersionTime(t *testing.T) {
	client := acc.NewFakeClient("../../testdata")

	// Request version before cutoff (should get v1)
	versionTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	result, err := ResolveDID(client, "did:acc:alice", &versionTime)
	require.NoError(t, err)
	assert.NotNil(t, result)

	doc := result.DIDDocument.(map[string]interface{})
	services, ok := doc["service"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, services, 2) // v1 has 2 services
}

func TestResolveDID_CaseNormalization(t *testing.T) {
	client := acc.NewFakeClient("../../testdata")

	// Test case normalization
	result, err := ResolveDID(client, "did:acc:ALICE", nil)
	require.NoError(t, err)
	assert.NotNil(t, result)

	doc := result.DIDDocument.(map[string]interface{})
	assert.Equal(t, "did:acc:alice", doc["id"])
}

func TestResolveDID_InvalidDID(t *testing.T) {
	client := acc.NewFakeClient("../../testdata")

	tests := []struct {
		name string
		did  string
	}{
		{"empty DID", ""},
		{"wrong method", "did:key:abc"},
		{"malformed", "did:acc:"},
		{"invalid chars", "did:acc:alice@example"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolveDID(client, tt.did, nil)
			assert.Error(t, err)
			assert.IsType(t, &InvalidDIDError{}, err)
		})
	}
}

func TestResolveDID_Deactivated(t *testing.T) {
	// Create a mock client that returns deactivated document
	client := &mockDeactivatedClient{}

	_, err := ResolveDID(client, "did:acc:alice", nil)
	assert.Error(t, err)
	assert.IsType(t, &DeactivatedError{}, err)
}

// mockDeactivatedClient returns a deactivated DID document
type mockDeactivatedClient struct{}

func (c *mockDeactivatedClient) GetLatestDIDEntry(adi string) (acc.Envelope, error) {
	doc := map[string]interface{}{
		"@context":    []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":          "did:acc:alice",
		"deactivated": true,
	}

	return acc.Envelope{
		ContentType: "application/did+json",
		Document:    doc,
		Meta: acc.EnvelopeMeta{
			VersionID:     "1704326400-final",
			Timestamp:     time.Now(),
			AuthorKeyPage: "acc://alice/book/1",
			Proof: acc.Proof{
				TxID: "0xdeactivated",
			},
		},
	}, nil
}

func (c *mockDeactivatedClient) GetEntryAtTime(adi string, t time.Time) (acc.Envelope, error) {
	return c.GetLatestDIDEntry(adi)
}

func (c *mockDeactivatedClient) GetKeyPageState(url string) (acc.KeyPageState, error) {
	return acc.KeyPageState{}, nil
}