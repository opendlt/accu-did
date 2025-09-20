package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicyV1_ValidateAuthorization(t *testing.T) {
	policy := NewPolicyV1()

	tests := []struct {
		name            string
		did             string
		authorKeyPage   string
		expectError     bool
		expectedMessage string
	}{
		{
			name:          "valid authorization",
			did:           "did:acc:alice",
			authorKeyPage: "acc://alice/book/1",
			expectError:   false,
		},
		{
			name:            "wrong book number",
			did:             "did:acc:alice",
			authorKeyPage:   "acc://alice/book/2",
			expectError:     true,
			expectedMessage: "unauthorized: expected acc://alice/book/1, got acc://alice/book/2",
		},
		{
			name:            "wrong ADI",
			did:             "did:acc:alice",
			authorKeyPage:   "acc://bob/book/1",
			expectError:     true,
			expectedMessage: "unauthorized: expected acc://alice/book/1, got acc://bob/book/1",
		},
		{
			name:          "complex ADI",
			did:           "did:acc:beastmode.acme",
			authorKeyPage: "acc://beastmode.acme/book/1",
			expectError:   false,
		},
		{
			name:          "case normalization",
			did:           "did:acc:ALICE",
			authorKeyPage: "acc://alice/book/1",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := policy.ValidateAuthorization(tt.did, tt.authorKeyPage)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedMessage != "" {
					assert.Contains(t, err.Error(), tt.expectedMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPolicyV1_GetRequiredKeyPage(t *testing.T) {
	policy := NewPolicyV1()

	tests := []struct {
		name         string
		did          string
		expected     string
		expectError  bool
	}{
		{
			name:     "simple ADI",
			did:      "did:acc:alice",
			expected: "acc://alice/book/1",
		},
		{
			name:     "complex ADI",
			did:      "did:acc:beastmode.acme",
			expected: "acc://beastmode.acme/book/1",
		},
		{
			name:     "case normalization",
			did:      "did:acc:ALICE",
			expected: "acc://alice/book/1",
		},
		{
			name:     "trailing dot",
			did:      "did:acc:alice.",
			expected: "acc://alice/book/1",
		},
		{
			name:        "invalid DID method",
			did:         "did:key:abc",
			expectError: true,
		},
		{
			name:        "empty DID",
			did:         "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := policy.GetRequiredKeyPage(tt.did)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestExtractADI(t *testing.T) {
	tests := []struct {
		name        string
		did         string
		expected    string
		expectError bool
	}{
		{
			name:     "simple ADI",
			did:      "did:acc:alice",
			expected: "alice",
		},
		{
			name:     "complex ADI",
			did:      "did:acc:beastmode.acme",
			expected: "beastmode.acme",
		},
		{
			name:     "with fragment",
			did:      "did:acc:alice#key-1",
			expected: "alice",
		},
		{
			name:     "with query",
			did:      "did:acc:alice?versionTime=123",
			expected: "alice",
		},
		{
			name:     "with path",
			did:      "did:acc:alice/path",
			expected: "alice",
		},
		{
			name:     "case normalization",
			did:      "did:acc:ALICE",
			expected: "alice",
		},
		{
			name:     "trailing dot removal",
			did:      "did:acc:alice.",
			expected: "alice",
		},
		{
			name:        "wrong method",
			did:         "did:key:abc",
			expectError: true,
		},
		{
			name:        "empty ADI",
			did:         "did:acc:",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractADI(tt.did)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateDID(t *testing.T) {
	tests := []struct {
		name        string
		did         string
		expectError bool
	}{
		{
			name: "valid DID",
			did:  "did:acc:alice",
		},
		{
			name: "valid complex DID",
			did:  "did:acc:beastmode.acme",
		},
		{
			name:        "empty DID",
			did:         "",
			expectError: true,
		},
		{
			name:        "wrong method",
			did:         "did:key:abc",
			expectError: true,
		},
		{
			name:        "malformed",
			did:         "did:acc:",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDID(tt.did)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}