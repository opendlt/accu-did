package accdid

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidDID = errors.New("invalid DID")
)

// ValidateDID validates that a DID string is properly formatted for the acc method
func ValidateDID(did string) error {
	if did == "" {
		return fmt.Errorf("%w: DID cannot be empty", ErrInvalidDID)
	}

	if !strings.HasPrefix(did, "did:acc:") {
		return fmt.Errorf("%w: missing 'did:acc:' prefix", ErrInvalidDID)
	}

	// Extract the identifier part after "did:acc:"
	identifier := strings.TrimPrefix(did, "did:acc:")
	if identifier == "" {
		return fmt.Errorf("%w: missing identifier after 'did:acc:'", ErrInvalidDID)
	}

	// Basic validation - identifiers should not start or end with special characters
	if strings.HasPrefix(identifier, "/") || strings.HasPrefix(identifier, ".") {
		return fmt.Errorf("%w: identifier cannot start with '/' or '.'", ErrInvalidDID)
	}

	if strings.HasSuffix(identifier, "/") || strings.HasSuffix(identifier, ".") {
		return fmt.Errorf("%w: identifier cannot end with '/' or '.'", ErrInvalidDID)
	}

	// Check for invalid characters
	for _, char := range identifier {
		if char < 32 || char == 127 { // Control characters
			return fmt.Errorf("%w: identifier contains invalid control character", ErrInvalidDID)
		}
		// Additional invalid characters for DID identifiers
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			return fmt.Errorf("%w: identifier contains invalid whitespace character", ErrInvalidDID)
		}
	}

	return nil
}

// ParseDID extracts components from a did:acc identifier
func ParseDID(did string) (adi string, path string, err error) {
	if err := ValidateDID(did); err != nil {
		return "", "", err
	}

	// Remove the "did:acc:" prefix
	identifier := strings.TrimPrefix(did, "did:acc:")

	// Split on the first slash to separate ADI from path
	parts := strings.SplitN(identifier, "/", 2)
	adi = parts[0]

	if len(parts) > 1 {
		path = parts[1]
	} else {
		path = "did" // Default path
	}

	return adi, path, nil
}