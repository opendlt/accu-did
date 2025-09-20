package did

import (
	"testing"
)

func TestParseDID(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedADI        string
		expectedDataPath   string
		expectedError      bool
	}{
		{
			name:             "simple DID",
			input:            "did:acc:alice",
			expectedADI:      "acc://alice",
			expectedDataPath: "acc://alice/did",
			expectedError:    false,
		},
		{
			name:             "DID with dots",
			input:            "did:acc:beastmode.acme",
			expectedADI:      "acc://beastmode.acme",
			expectedDataPath: "acc://beastmode.acme/did",
			expectedError:    false,
		},
		{
			name:             "DID with custom path",
			input:            "did:acc:alice/documents",
			expectedADI:      "acc://alice",
			expectedDataPath: "acc://alice/documents",
			expectedError:    false,
		},
		{
			name:             "DID with nested path",
			input:            "did:acc:alice/data/personal",
			expectedADI:      "acc://alice",
			expectedDataPath: "acc://alice/data/personal",
			expectedError:    false,
		},
		{
			name:          "invalid method",
			input:         "did:web:example.com",
			expectedError: true,
		},
		{
			name:          "empty identifier",
			input:         "did:acc:",
			expectedError: true,
		},
		{
			name:          "missing did prefix",
			input:         "acc:alice",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adiURL, dataAccountURL, err := ParseDID(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error for input %s, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error for input %s: %v", tt.input, err)
				return
			}

			if adiURL.String() != tt.expectedADI {
				t.Errorf("expected ADI URL %s, got %s", tt.expectedADI, adiURL.String())
			}

			if dataAccountURL.String() != tt.expectedDataPath {
				t.Errorf("expected data account URL %s, got %s", tt.expectedDataPath, dataAccountURL.String())
			}
		})
	}
}

func TestFormatDID(t *testing.T) {
	tests := []struct {
		name     string
		adiLabel string
		path     string
		expected string
	}{
		{
			name:     "simple ADI",
			adiLabel: "alice",
			path:     "",
			expected: "did:acc:alice",
		},
		{
			name:     "ADI with default path",
			adiLabel: "alice",
			path:     "did",
			expected: "did:acc:alice",
		},
		{
			name:     "ADI with custom path",
			adiLabel: "alice",
			path:     "documents",
			expected: "did:acc:alice/documents",
		},
		{
			name:     "ADI with dots and path",
			adiLabel: "beastmode.acme",
			path:     "credentials",
			expected: "did:acc:beastmode.acme/credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDID(tt.adiLabel, tt.path)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestExtractADILabel(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError bool
	}{
		{
			name:     "simple DID",
			input:    "did:acc:alice",
			expected: "alice",
		},
		{
			name:     "DID with dots",
			input:    "did:acc:beastmode.acme",
			expected: "beastmode.acme",
		},
		{
			name:     "DID with path",
			input:    "did:acc:alice/documents",
			expected: "alice",
		},
		{
			name:          "invalid method",
			input:         "did:web:example.com",
			expectedError: true,
		},
		{
			name:          "empty identifier",
			input:         "did:acc:",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractADILabel(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error for input %s, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error for input %s: %v", tt.input, err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}