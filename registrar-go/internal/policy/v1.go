package policy

import (
	"fmt"
	"strings"
)

// AuthPolicy defines the interface for authorization policies
type AuthPolicy interface {
	ValidateAuthorization(did string, authorKeyPage string) error
	GetRequiredKeyPage(did string) (string, error)
}

// PolicyV1 implements the default authorization policy:
// Only acc://<adi>/book/1 may author DID document updates
type PolicyV1 struct {
	// Future: could add configuration for different book numbers
}

// NewPolicyV1 creates a new PolicyV1 instance
func NewPolicyV1() *PolicyV1 {
	return &PolicyV1{}
}

// ValidateAuthorization checks if the given authorKeyPage is authorized for the DID
func (p *PolicyV1) ValidateAuthorization(did string, authorKeyPage string) error {
	requiredKeyPage, err := p.GetRequiredKeyPage(did)
	if err != nil {
		return err
	}

	if authorKeyPage != requiredKeyPage {
		return fmt.Errorf("unauthorized: expected %s, got %s", requiredKeyPage, authorKeyPage)
	}

	return nil
}

// GetRequiredKeyPage returns the required key page URL for a given DID
func (p *PolicyV1) GetRequiredKeyPage(did string) (string, error) {
	// Extract ADI from DID
	adi, err := extractADI(did)
	if err != nil {
		return "", err
	}

	// Policy v1: always use book/1
	return fmt.Sprintf("acc://%s/book/1", adi), nil
}

// extractADI extracts the ADI name from a DID
func extractADI(did string) (string, error) {
	if !strings.HasPrefix(did, "did:acc:") {
		return "", fmt.Errorf("not a did:acc DID")
	}

	// Extract the part after "did:acc:"
	remainder := strings.TrimPrefix(did, "did:acc:")

	// Split on first occurrence of path, query, or fragment to get just the ADI
	adi := remainder
	for _, separator := range []string{"/", "?", "#", ";"} {
		if idx := strings.Index(remainder, separator); idx != -1 {
			adi = remainder[:idx]
			break
		}
	}

	// Validate ADI is not empty
	if adi == "" {
		return "", fmt.Errorf("empty ADI name")
	}

	// Normalize ADI (lowercase, remove trailing dots)
	adi = strings.ToLower(adi)
	adi = strings.TrimSuffix(adi, ".")

	return adi, nil
}

// ValidateDID validates that a DID is properly formatted
func ValidateDID(did string) error {
	if did == "" {
		return fmt.Errorf("DID cannot be empty")
	}

	if !strings.HasPrefix(did, "did:acc:") {
		return fmt.Errorf("DID must start with 'did:acc:'")
	}

	_, err := extractADI(did)
	return err
}