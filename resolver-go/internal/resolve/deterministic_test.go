package resolve

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeterministicResolver_SequenceTiebreaking(t *testing.T) {
	mockClient := &MockClient{
		entries: []*DataEntry{
			{
				Sequence:  &[]uint64{100}[0], // Same sequence
				Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Data:      []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:test","version":"v1"}`),
			},
			{
				Sequence:  &[]uint64{100}[0], // Same sequence
				Timestamp: time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC), // Later timestamp
				Data:      []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:test","version":"v2"}`),
			},
		},
	}

	resolver := NewDeterministicResolver(mockClient, ResolveOrderSequence)
	result, err := resolver.ResolveDID("did:acc:test", nil)
	require.NoError(t, err)

	// Should pick later timestamp when sequences are equal
	doc := result.DIDDocument.(map[string]interface{})
	assert.Equal(t, "v2", doc["version"])
}

func TestDeterministicResolver_TimestampTiebreaking(t *testing.T) {
	mockClient := &MockClient{
		entries: []*DataEntry{
			{
				Sequence:  &[]uint64{100}[0],
				Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), // Same timestamp
				Data:      []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:test","version":"hash1"}`),
			},
			{
				Sequence:  &[]uint64{100}[0],
				Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), // Same timestamp
				Data:      []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:test","version":"hash2"}`),
			},
		},
	}

	resolver := NewDeterministicResolver(mockClient, ResolveOrderSequence)
	result, err := resolver.ResolveDID("did:acc:test", nil)
	require.NoError(t, err)

	// Should consistently pick same result based on content hash tiebreaking
	doc := result.DIDDocument.(map[string]interface{})
	version := doc["version"].(string)
	assert.True(t, version == "hash1" || version == "hash2", "Should pick one consistently")

	// Run again to verify consistency
	result2, err := resolver.ResolveDID("did:acc:test", nil)
	require.NoError(t, err)
	doc2 := result2.DIDDocument.(map[string]interface{})
	assert.Equal(t, version, doc2["version"], "Should be deterministic")
}

func TestDeterministicResolver_MalformedFiltering(t *testing.T) {
	mockClient := &MockClient{
		entries: []*DataEntry{
			{
				Sequence:  &[]uint64{99}[0],
				Timestamp: time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
				Data:      []byte(`{invalid json`), // Malformed
			},
			{
				Sequence:  &[]uint64{100}[0],
				Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Data:      []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:test","version":"valid"}`),
			},
		},
	}

	resolver := NewDeterministicResolver(mockClient, ResolveOrderSequence)
	result, err := resolver.ResolveDID("did:acc:test", nil)
	require.NoError(t, err)

	// Should ignore malformed entry and use valid one
	doc := result.DIDDocument.(map[string]interface{})
	assert.Equal(t, "valid", doc["version"])
}

func TestDeterministicResolver_DeactivatedMetadata(t *testing.T) {
	mockClient := &MockClient{
		entries: []*DataEntry{
			{
				Sequence:  &[]uint64{100}[0],
				Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Data:      []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:test","deactivated":true,"deactivatedAt":"2024-01-01T12:00:00Z"}`),
			},
		},
	}

	resolver := NewDeterministicResolver(mockClient, ResolveOrderSequence)
	result, err := resolver.ResolveDID("did:acc:test", nil)
	require.NoError(t, err)

	// Should mark as deactivated in metadata
	assert.True(t, result.DIDDocumentMetadata.Deactivated)
	assert.NotEmpty(t, result.DIDResolutionMetadata.ContentHash)

	// Document should contain deactivation fields
	doc := result.DIDDocument.(map[string]interface{})
	assert.Equal(t, true, doc["deactivated"])
	assert.Equal(t, "2024-01-01T12:00:00Z", doc["deactivatedAt"])
}

// MockClient for deterministic testing
type MockClient struct {
	entries []*DataEntry
}

func (m *MockClient) GetDataEntries(adi string) ([]*DataEntry, error) {
	return m.entries, nil
}

func (m *MockClient) GetDataAccountEntry(dataAccountURL interface{}) ([]byte, error) {
	if len(m.entries) > 0 {
		return m.entries[len(m.entries)-1].Data, nil
	}
	return nil, &NotFoundError{DID: "did:acc:test"}
}

func (m *MockClient) GetLatestDIDEntry(adi string) (interface{}, error) {
	return nil, nil
}

func (m *MockClient) GetEntryAtTime(adi string, t time.Time) (interface{}, error) {
	return nil, nil
}

func (m *MockClient) GetKeyPageState(url string) (interface{}, error) {
	return nil, nil
}