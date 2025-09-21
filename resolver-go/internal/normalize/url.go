package normalize

import (
	"fmt"
	"net/url"
	"strings"
)

// NormalizedDIDURL represents a parsed and normalized DID URL
type NormalizedDIDURL struct {
	Scheme           string            `json:"scheme"`
	Method           string            `json:"method"`
	MethodSpecificID string            `json:"methodSpecificId"`
	Path             string            `json:"path"`
	Query            map[string]string `json:"query"`
	Fragment         string            `json:"fragment"`
}

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

// NormalizeDIDURL parses and normalizes a DID URL into structured components
func NormalizeDIDURL(didURL string) (NormalizedDIDURL, error) {
	result := NormalizedDIDURL{
		Query: make(map[string]string),
	}

	// Parse the URL
	u, err := url.Parse(didURL)
	if err != nil {
		return result, fmt.Errorf("invalid URL: %w", err)
	}

	// Validate scheme
	if u.Scheme != "did" {
		return result, fmt.Errorf("invalid scheme: expected 'did', got '%s'", u.Scheme)
	}
	result.Scheme = u.Scheme

	// Extract method and method-specific ID from opaque part
	// DID URLs have the form: did:method:method-specific-id
	parts := strings.SplitN(u.Opaque, ":", 2)
	if len(parts) < 2 {
		return result, fmt.Errorf("invalid DID format: missing method or method-specific-id")
	}

	method := parts[0]
	methodSpecificPart := parts[1]

	// Validate method
	if method != "acc" {
		return result, fmt.Errorf("invalid method: expected 'acc', got '%s'", method)
	}
	result.Method = method

	// Split method-specific part on first path/query/fragment separator
	methodSpecificID := methodSpecificPart
	pathQueryFragment := ""

	for _, sep := range []string{"/", "?", "#", ";"} {
		if idx := strings.Index(methodSpecificPart, sep); idx != -1 {
			methodSpecificID = methodSpecificPart[:idx]
			pathQueryFragment = methodSpecificPart[idx:]
			break
		}
	}

	// Validate method-specific ID is not empty
	if methodSpecificID == "" {
		return result, fmt.Errorf("empty method-specific-id")
	}

	// Normalize the ADI name
	normalizedADI := normalizeADI(methodSpecificID)
	if err := ValidateADIName(normalizedADI); err != nil {
		return result, fmt.Errorf("invalid method-specific-id: %w", err)
	}
	result.MethodSpecificID = normalizedADI

	// Parse path, query, and fragment from the remaining part
	if pathQueryFragment != "" {
		// Create a temporary URL to parse the path/query/fragment
		tempURL, err := url.Parse("did:acc:temp" + pathQueryFragment)
		if err != nil {
			return result, fmt.Errorf("invalid path/query/fragment: %w", err)
		}

		result.Path = tempURL.Path
		result.Fragment = tempURL.Fragment

		// Parse query parameters
		if tempURL.RawQuery != "" {
			queryParams, err := url.ParseQuery(tempURL.RawQuery)
			if err != nil {
				return result, fmt.Errorf("invalid query parameters: %w", err)
			}

			for key, values := range queryParams {
				if len(values) > 0 {
					result.Query[key] = values[0] // Take first value if multiple
				}
			}
		}
	}

	return result, nil
}
