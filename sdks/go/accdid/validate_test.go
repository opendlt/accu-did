package accdid

import (
	"errors"
	"testing"
)

func TestValidateDID(t *testing.T) {
	tests := []struct {
		name      string
		did       string
		expectErr bool
	}{
		{"valid simple DID", "did:acc:alice", false},
		{"valid DID with dots", "did:acc:beastmode.acme", false},
		{"valid DID with path", "did:acc:alice/documents", false},
		{"valid DID with subdomain", "did:acc:subdomain.example.com", false},
		{"empty DID", "", true},
		{"missing prefix", "alice", true},
		{"wrong method", "did:web:example.com", true},
		{"starts with slash", "did:acc:/invalid", true},
		{"starts with dot", "did:acc:.invalid", true},
		{"ends with slash", "did:acc:invalid/", true},
		{"ends with dot", "did:acc:invalid.", true},
		{"contains spaces", "did:acc:alice bob", true},
		{"contains tab", "did:acc:alice\tcharacter", true},
		{"contains newline", "did:acc:alice\ncharacter", true},
		{"missing identifier", "did:acc:", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDID(tt.did)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for DID %s, got nil", tt.did)
				}
				if !errors.Is(err, ErrInvalidDID) {
					t.Errorf("Expected ErrInvalidDID, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for DID %s: %v", tt.did, err)
				}
			}
		})
	}
}

func TestParseDID(t *testing.T) {
	tests := []struct {
		name         string
		did          string
		expectedADI  string
		expectedPath string
		expectErr    bool
	}{
		{
			name:         "simple DID",
			did:          "did:acc:alice",
			expectedADI:  "alice",
			expectedPath: "did",
			expectErr:    false,
		},
		{
			name:         "DID with path",
			did:          "did:acc:alice/documents",
			expectedADI:  "alice",
			expectedPath: "documents",
			expectErr:    false,
		},
		{
			name:         "DID with dots",
			did:          "did:acc:beastmode.acme",
			expectedADI:  "beastmode.acme",
			expectedPath: "did",
			expectErr:    false,
		},
		{
			name:         "DID with complex path",
			did:          "did:acc:company.com/services/auth",
			expectedADI:  "company.com",
			expectedPath: "services/auth",
			expectErr:    false,
		},
		{
			name:      "invalid DID",
			did:       "invalid",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adi, path, err := ParseDID(tt.did)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for DID %s, got nil", tt.did)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for DID %s: %v", tt.did, err)
				return
			}

			if adi != tt.expectedADI {
				t.Errorf("Expected ADI %s, got %s", tt.expectedADI, adi)
			}

			if path != tt.expectedPath {
				t.Errorf("Expected path %s, got %s", tt.expectedPath, path)
			}
		})
	}
}