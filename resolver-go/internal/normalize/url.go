package normalize

import (
	"fmt"
	"strings"
)

// NormalizeDID normalizes a DID URL according to did:acc method rules
func NormalizeDID(did string) (normalizedDID, adi string, err error) {
	// Basic validation
	if !strings.HasPrefix(did, "did:acc:") {
		return "", "", fmt.Errorf("not a did:acc DID")
	}

	// Extract the part after "did:acc:"
	remainder := strings.TrimPrefix(did, "did:acc:")

	// Split on first occurrence of path, query, or fragment
	adi = remainder
	for _, separator := range []string{"/", "?", "#", ";"} {
		if idx := strings.Index(remainder, separator); idx != -1 {
			adi = remainder[:idx]
			break
		}
	}

	// Validate ADI is not empty
	if adi == "" {
		return "", "", fmt.Errorf("empty ADI name")
	}

	// Normalize ADI name first (case + trailing dots)
	normalizedADI := normalizeADI(adi)

	// Then validate the normalized ADI contains only allowed characters
	if err := ValidateADIName(normalizedADI); err != nil {
		return "", "", fmt.Errorf("invalid ADI name: %w", err)
	}

	// Reconstruct normalized DID
	normalizedDID = "did:acc:" + normalizedADI

	// Preserve the rest of the DID URL (path, query, fragment)
	if len(remainder) > len(adi) {
		normalizedDID += remainder[len(adi):]
	}

	return normalizedDID, normalizedADI, nil
}

// normalizeADI normalizes an ADI name according to the rules:
// - Convert to lowercase
// - Remove trailing dots
func normalizeADI(adi string) string {
	// Convert to lowercase
	normalized := strings.ToLower(adi)

	// Remove trailing dots
	normalized = strings.TrimSuffix(normalized, ".")

	return normalized
}

// ValidateADIName validates that an ADI name contains only allowed characters
func ValidateADIName(adi string) error {
	if adi == "" {
		return fmt.Errorf("ADI name cannot be empty")
	}

	// Check each character
	for i, r := range adi {
		if !isValidADIChar(r) {
			return fmt.Errorf("invalid character '%c' at position %d", r, i)
		}
	}

	// Check for invalid dot placement
	if strings.HasPrefix(adi, ".") || strings.HasSuffix(adi, ".") {
		return fmt.Errorf("ADI name cannot start or end with dots")
	}

	if strings.Contains(adi, "..") {
		return fmt.Errorf("ADI name cannot contain consecutive dots")
	}

	return nil
}

// isValidADIChar checks if a character is valid in an ADI name
func isValidADIChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_' || r == '.'
}