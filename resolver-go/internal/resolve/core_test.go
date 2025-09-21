package resolve

import (
	"testing"
	"time"

	"github.com/opendlt/accu-did/resolver-go/internal/acc"
)

func TestResolveDID(t *testing.T) {
	// Create fake client for testing
	client := acc.NewFakeClient("../../testdata")

	tests := []struct {
		name          string
		did           string
		versionTime   *time.Time
		expectedError bool
		checkDoc      func(*testing.T, *DIDResolutionResult)
	}{
		{
			name:          "resolve alice DID",
			did:           "did:acc:alice",
			expectedError: false,
			checkDoc: func(t *testing.T, result *DIDResolutionResult) {
				if result.DIDDocument == nil {
					t.Error("expected DID document, got nil")
					return
				}

				doc, ok := result.DIDDocument.(map[string]interface{})
				if !ok {
					t.Error("expected DID document to be map[string]interface{}")
					return
				}

				if doc["id"] != "did:acc:alice" {
					t.Errorf("expected id 'did:acc:alice', got %v", doc["id"])
				}

				if result.DIDResolutionMetadata.ContentType != "application/did+json" {
					t.Errorf("expected content type 'application/did+json', got %s", result.DIDResolutionMetadata.ContentType)
				}
			},
		},
		{
			name:          "resolve beastmode.acme DID",
			did:           "did:acc:beastmode.acme",
			expectedError: false,
			checkDoc: func(t *testing.T, result *DIDResolutionResult) {
				if result.DIDDocument == nil {
					t.Error("expected DID document, got nil")
					return
				}

				doc, ok := result.DIDDocument.(map[string]interface{})
				if !ok {
					t.Error("expected DID document to be map[string]interface{}")
					return
				}

				if doc["id"] != "did:acc:beastmode.acme" {
					t.Errorf("expected id 'did:acc:beastmode.acme', got %v", doc["id"])
				}

				// Check for verification methods
				vm, ok := doc["verificationMethod"].([]interface{})
				if !ok || len(vm) < 2 {
					t.Error("expected at least 2 verification methods")
				}
			},
		},
		{
			name:          "resolve deactivated DID",
			did:           "did:acc:deactivated",
			expectedError: true,
		},
		{
			name:          "resolve non-existent DID",
			did:           "did:acc:nonexistent",
			expectedError: true,
		},
		{
			name:          "invalid DID format",
			did:           "did:web:example.com",
			expectedError: true,
		},
		{
			name:          "empty DID",
			did:           "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolveDID(client, tt.did, tt.versionTime)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error for DID %s, got nil", tt.did)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error for DID %s: %v", tt.did, err)
				return
			}

			if result == nil {
				t.Error("expected result, got nil")
				return
			}

			// Run custom check if provided
			if tt.checkDoc != nil {
				tt.checkDoc(t, result)
			}

			// Basic checks
			if result.DIDResolutionMetadata.Pattern != "^did:acc:" {
				t.Errorf("expected pattern '^did:acc:', got %s", result.DIDResolutionMetadata.Pattern)
			}

			if result.DIDResolutionMetadata.Duration <= 0 {
				t.Error("expected positive duration")
			}
		})
	}
}

func TestDIDResolutionErrors(t *testing.T) {
	client := acc.NewFakeClient("../../testdata")

	// Test error types
	t.Run("NotFoundError", func(t *testing.T) {
		_, err := ResolveDID(client, "did:acc:notfound", nil)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}

		if _, ok := err.(*NotFoundError); !ok {
			t.Errorf("expected NotFoundError, got %T", err)
		}
	})

	t.Run("InvalidDIDError", func(t *testing.T) {
		_, err := ResolveDID(client, "invalid-did", nil)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}

		if _, ok := err.(*InvalidDIDError); !ok {
			t.Errorf("expected InvalidDIDError, got %T", err)
		}
	})

	t.Run("DeactivatedError", func(t *testing.T) {
		_, err := ResolveDID(client, "did:acc:deactivated", nil)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}

		if _, ok := err.(*DeactivatedError); !ok {
			t.Errorf("expected DeactivatedError, got %T", err)
		}
	})
}

func TestDIDDocumentMetadata(t *testing.T) {
	client := acc.NewFakeClient("../../testdata")

	result, err := ResolveDID(client, "did:acc:alice", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metadata := result.DIDDocumentMetadata

	if metadata.VersionID == "" {
		t.Error("expected version ID, got empty string")
	}

	if metadata.Created.IsZero() {
		t.Error("expected created timestamp, got zero time")
	}

	if metadata.Updated.IsZero() {
		t.Error("expected updated timestamp, got zero time")
	}

	if metadata.Deactivated {
		t.Error("expected deactivated to be false")
	}
}
