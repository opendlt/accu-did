package ops

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opendlt/accu-did/registrar-go/internal/policy"
)

func TestBuildEnvelope(t *testing.T) {
	// Test document
	document := map[string]interface{}{
		"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":       "did:acc:alice",
		"verificationMethod": []interface{}{
			map[string]interface{}{
				"id":         "did:acc:alice#key-1",
				"type":       "AccumulateKeyPage",
				"controller": "did:acc:alice",
				"keyPageUrl": "acc://alice/book/1",
				"threshold":  1,
			},
		},
	}

	authorKeyPage := "acc://alice/book/1"
	previousVersionID := ""

	envelope, err := BuildEnvelope(document, authorKeyPage, previousVersionID)
	require.NoError(t, err)
	assert.NotNil(t, envelope)

	// Validate envelope structure
	assert.Equal(t, "application/did+json", envelope.ContentType)
	assert.Equal(t, document, envelope.Document)
	assert.Equal(t, authorKeyPage, envelope.Meta.AuthorKeyPage)
	assert.NotEmpty(t, envelope.Meta.VersionID)
	assert.Empty(t, envelope.Meta.PreviousVersionID)
	assert.NotEmpty(t, envelope.Meta.Proof.ContentHash)
	assert.True(t, envelope.Meta.Timestamp.Before(time.Now().Add(time.Second)))

	// Validate content hash
	err = envelope.ValidateContentHash()
	assert.NoError(t, err)
}

func TestBuildEnvelopeWithPreviousVersion(t *testing.T) {
	document := map[string]interface{}{
		"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":       "did:acc:alice",
	}

	authorKeyPage := "acc://alice/book/1"
	previousVersionID := "1704067200-abc123"

	envelope, err := BuildEnvelope(document, authorKeyPage, previousVersionID)
	require.NoError(t, err)

	assert.Equal(t, previousVersionID, envelope.Meta.PreviousVersionID)
}

func TestEnvelopeSetTransactionID(t *testing.T) {
	document := map[string]interface{}{
		"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":       "did:acc:alice",
	}

	envelope, err := BuildEnvelope(document, "acc://alice/book/1", "")
	require.NoError(t, err)

	txID := "0x1234567890abcdef"
	envelope.SetTransactionID(txID)

	assert.Equal(t, txID, envelope.Meta.Proof.TxID)
}

func TestEnvelopeContentHashValidation(t *testing.T) {
	document := map[string]interface{}{
		"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":       "did:acc:alice",
	}

	envelope, err := BuildEnvelope(document, "acc://alice/book/1", "")
	require.NoError(t, err)

	// Valid hash should pass
	err = envelope.ValidateContentHash()
	assert.NoError(t, err)

	// Modify document to make hash invalid
	envelope.Document["id"] = "did:acc:bob"
	err = envelope.ValidateContentHash()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "content hash mismatch")
}

func TestVersionIDGeneration(t *testing.T) {
	timestamp := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	versionID := generateVersionID(timestamp)

	assert.NotEmpty(t, versionID)
	assert.Contains(t, versionID, "1704067200") // Unix timestamp
	assert.Contains(t, versionID, "-")
	assert.Len(t, versionID, 19) // timestamp(10) + "-" + suffix(8)
}

func TestEnvelopeHashVectors(t *testing.T) {
	// Load test vectors
	vectorPath := filepath.Join("..", "..", "..", "spec", "vectors", "envelope-hash.json")
	vectorData, err := os.ReadFile(vectorPath)
	require.NoError(t, err, "Failed to read envelope hash vectors")

	var vectors struct {
		Description     string `json:"description"`
		Version         string `json:"version"`
		Algorithm       string `json:"algorithm"`
		DocumentVectors []struct {
			Description  string      `json:"description"`
			Document     interface{} `json:"document"`
			ExpectedHash string      `json:"expectedHash"`
		} `json:"document_vectors"`
	}
	err = json.Unmarshal(vectorData, &vectors)
	require.NoError(t, err, "Failed to parse envelope hash vectors")

	for _, vector := range vectors.DocumentVectors {
		t.Run(vector.Description, func(t *testing.T) {
			envelope, err := BuildEnvelope(vector.Document, "acc://test/book/1", "")
			require.NoError(t, err)

			actualHash := envelope.GetContentHash()
			assert.Equal(t, vector.ExpectedHash, actualHash, "Content hash mismatch")
		})
	}
}

func TestAuthorizationVectors(t *testing.T) {
	// Load test vectors
	vectorPath := filepath.Join("..", "..", "..", "spec", "vectors", "auth-cases.json")
	vectorData, err := os.ReadFile(vectorPath)
	require.NoError(t, err, "Failed to read authorization vectors")

	var vectors struct {
		Description string `json:"description"`
		Version     string `json:"version"`
		Policy      string `json:"policy"`
		Vectors     []struct {
			Description   string `json:"description"`
			DID          string `json:"did"`
			AuthorKeyPage string `json:"authorKeyPage"`
			ExpectedResult string `json:"expectedResult"`
			ExpectedError string `json:"expectedError"`
		} `json:"vectors"`
	}
	err = json.Unmarshal(vectorData, &vectors)
	require.NoError(t, err, "Failed to parse authorization vectors")

	policy := policy.NewPolicyV1()

	for _, vector := range vectors.Vectors {
		t.Run(vector.Description, func(t *testing.T) {
			err := policy.ValidateAuthorization(vector.DID, vector.AuthorKeyPage)

			if vector.ExpectedResult == "success" {
				assert.NoError(t, err, "Authorization should succeed")
			} else {
				assert.Error(t, err, "Authorization should fail")
				if vector.ExpectedError != "" {
					assert.Contains(t, err.Error(), vector.ExpectedError, "Error message mismatch")
				}
			}
		})
	}
}

func TestEnvelopeReplayProtection(t *testing.T) {
	document := map[string]interface{}{
		"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":       "did:acc:alice",
	}

	// Create envelope with previous version
	previousVersionID := "1704063600-previous"
	envelope, err := BuildEnvelope(document, "acc://alice/book/1", previousVersionID)
	require.NoError(t, err)

	// Verify previous version is set
	assert.Equal(t, previousVersionID, envelope.Meta.PreviousVersionID)

	// Verify version ID is greater than previous
	currentTime := extractTimestampFromVersionID(envelope.Meta.VersionID)
	previousTime := extractTimestampFromVersionID(previousVersionID)
	assert.Greater(t, currentTime, previousTime, "Current version timestamp must be greater than previous")
}

func extractTimestampFromVersionID(versionID string) int64 {
	parts := strings.Split(versionID, "-")
	if len(parts) != 2 {
		return 0
	}

	timestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0
	}

	return timestamp
}