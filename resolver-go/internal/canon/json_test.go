package canon

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type EnvelopeHashVector struct {
	Description  string      `json:"description"`
	Document     interface{} `json:"document"`
	ExpectedHash string      `json:"expectedHash"`
}

type EnvelopeHashVectors struct {
	Description           string               `json:"description"`
	Version               string               `json:"version"`
	Algorithm             string               `json:"algorithm"`
	DocumentVectors       []EnvelopeHashVector `json:"document_vectors"`
	CanonicalizationNotes []string             `json:"canonicalization_notes"`
}

func TestEnvelopeHashVectors(t *testing.T) {
	// Load test vectors
	vectorPath := filepath.Join("..", "..", "..", "spec", "vectors", "envelope-hash.json")
	vectorData, err := os.ReadFile(vectorPath)
	require.NoError(t, err, "Failed to read envelope hash vectors")

	var vectors EnvelopeHashVectors
	err = json.Unmarshal(vectorData, &vectors)
	require.NoError(t, err, "Failed to parse envelope hash vectors")

	for _, vector := range vectors.DocumentVectors {
		t.Run(vector.Description, func(t *testing.T) {
			hash, err := ComputeContentHash(vector.Document)
			require.NoError(t, err, "Hash computation should not fail")

			assert.Equal(t, vector.ExpectedHash, hash, "Content hash mismatch")
		})
	}
}

func TestCanonicalizeJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "simple object",
			input: map[string]interface{}{
				"b": 2,
				"a": 1,
			},
			expected: `{"a":1,"b":2}`,
		},
		{
			name: "nested object",
			input: map[string]interface{}{
				"z": map[string]interface{}{
					"y": 2,
					"x": 1,
				},
				"a": 1,
			},
			expected: `{"a":1,"z":{"x":1,"y":2}}`,
		},
		{
			name:     "array",
			input:    []interface{}{3, 1, 2},
			expected: `[3,1,2]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CanonicalizeJSON(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestComputeContentHash(t *testing.T) {
	document := map[string]interface{}{
		"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":       "did:acc:alice",
	}

	hash1, err := ComputeContentHash(document)
	require.NoError(t, err)

	// Hash should be deterministic
	hash2, err := ComputeContentHash(document)
	require.NoError(t, err)
	assert.Equal(t, hash1, hash2)

	// Hash should be 64 character hex string (SHA-256)
	assert.Len(t, hash1, 64)
	assert.Regexp(t, "^[a-f0-9]{64}$", hash1)
}
