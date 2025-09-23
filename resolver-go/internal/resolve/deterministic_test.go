package resolve

import (
	"testing"
	"time"

	"github.com/opendlt/accu-did/resolver-go/internal/acc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	accurl "gitlab.com/accumulatenetwork/accumulate/pkg/url"
)

func TestDeterministicResolver_SequenceTiebreaking(t *testing.T) {
	mockClient := &DeterministicMockClient{entries: []*DataEntry{
		{
			Sequence:  &[]uint64{100}[0], // Same sequence
			Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			Data:      []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:test","version":"v1"}`),
		},
		{
			Sequence:  &[]uint64{100}[0],                            // Same sequence
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
	mockClient := &DeterministicMockClient{entries: []*DataEntry{
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
	mockClient := &DeterministicMockClient{entries: []*DataEntry{
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
	mockClient := &DeterministicMockClient{entries: []*DataEntry{
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
	assert.NotEmpty(t, result.DIDDocumentMetadata.ContentHash)

	// Document should contain deactivation fields
	doc := result.DIDDocument.(map[string]interface{})
	assert.Equal(t, true, doc["deactivated"])
	assert.Equal(t, "2024-01-01T12:00:00Z", doc["deactivatedAt"])
}

// DeterministicMockClient is test-only and avoids colliding with resolve.MockClient
type DeterministicMockClient struct {
	entries []*DataEntry
}

// Optional helper used only by these tests
func (m *DeterministicMockClient) GetDataEntries(adi string) ([]*DataEntry, error) {
	return m.entries, nil
}

// Implement acc.Client exactly:
func (m *DeterministicMockClient) GetLatestDIDEntry(adi string) (acc.Envelope, error) {
	// Not used by these tests; return zero-value envelope.
	return acc.Envelope{}, nil
}

func (m *DeterministicMockClient) GetEntryAtTime(adi string, t time.Time) (acc.Envelope, error) {
	// Not used by these tests; return zero-value envelope.
	return acc.Envelope{}, nil
}

func (m *DeterministicMockClient) GetKeyPageState(u string) (acc.KeyPageState, error) {
	// Minimal stub for tests
	return acc.KeyPageState{URL: u, Threshold: 1}, nil
}

func (m *DeterministicMockClient) GetDataAccountEntry(dataAccountURL *accurl.URL) ([]byte, error) {
	if len(m.entries) > 0 {
		return m.entries[len(m.entries)-1].Data, nil
	}
	return nil, &NotFoundError{DID: "did:acc:test"}
}
