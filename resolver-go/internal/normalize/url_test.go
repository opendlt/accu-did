package normalize

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type URLNormalizationVector struct {
	Description string           `json:"description"`
	Input       string           `json:"input"`
	Expected    NormalizedDIDURL `json:"expected"`
}

type URLNormalizationVectors struct {
	Description string                   `json:"description"`
	Version     string                   `json:"version"`
	Vectors     []URLNormalizationVector `json:"vectors"`
}

func TestURLNormalizationVectors(t *testing.T) {
	// Load test vectors
	vectorPath := filepath.Join("..", "..", "..", "spec", "vectors", "url-normalization.json")
	vectorData, err := os.ReadFile(vectorPath)
	require.NoError(t, err, "Failed to read URL normalization vectors")

	var vectors URLNormalizationVectors
	err = json.Unmarshal(vectorData, &vectors)
	require.NoError(t, err, "Failed to parse URL normalization vectors")

	for _, vector := range vectors.Vectors {
		t.Run(vector.Description, func(t *testing.T) {
			result, err := NormalizeDIDURL(vector.Input)
			require.NoError(t, err, "URL normalization should not fail for valid input")

			assert.Equal(t, vector.Expected.Scheme, result.Scheme, "Scheme mismatch")
			assert.Equal(t, vector.Expected.Method, result.Method, "Method mismatch")
			assert.Equal(t, vector.Expected.MethodSpecificID, result.MethodSpecificID, "Method-specific ID mismatch")
			assert.Equal(t, vector.Expected.Path, result.Path, "Path mismatch")
			assert.Equal(t, vector.Expected.Query, result.Query, "Query mismatch")
			assert.Equal(t, vector.Expected.Fragment, result.Fragment, "Fragment mismatch")
		})
	}
}

func TestNormalizeDIDURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    NormalizedDIDURL
		expectError bool
	}{
		{
			name:  "simple DID",
			input: "did:acc:alice",
			expected: NormalizedDIDURL{
				Scheme:           "did",
				Method:           "acc",
				MethodSpecificID: "alice",
				Path:             "",
				Query:            map[string]string{},
				Fragment:         "",
			},
		},
		{
			name:        "invalid scheme",
			input:       "foo:acc:alice",
			expectError: true,
		},
		{
			name:        "invalid method",
			input:       "did:key:alice",
			expectError: true,
		},
		{
			name:        "empty method-specific-id",
			input:       "did:acc:",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeDIDURL(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
