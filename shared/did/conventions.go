package did

import (
	"fmt"
	"strings"

	"gitlab.com/accumulatenetwork/accumulate/pkg/url"
)

// ParseDID converts did:acc:label[/path] to Accumulate URLs
func ParseDID(did string) (adiURL, dataAccountURL *url.URL, err error) {
	if !strings.HasPrefix(did, "did:acc:") {
		return nil, nil, fmt.Errorf("invalid DID method: %s", did)
	}

	identifier := strings.TrimPrefix(did, "did:acc:")
	parts := strings.SplitN(identifier, "/", 2)

	adiLabel := parts[0]
	if adiLabel == "" {
		return nil, nil, fmt.Errorf("empty ADI label")
	}

	adiURL, err = url.Parse(fmt.Sprintf("acc://%s", adiLabel))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid ADI URL: %w", err)
	}

	// Default data account path is /did unless path specified
	dataPath := "did"
	if len(parts) > 1 && parts[1] != "" {
		dataPath = parts[1]
	}

	dataAccountURL, err = url.Parse(fmt.Sprintf("acc://%s/%s", adiLabel, dataPath))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid data account URL: %w", err)
	}

	return adiURL, dataAccountURL, nil
}

// FormatDID creates did:acc:label[/path] from ADI and optional path
func FormatDID(adiLabel, path string) string {
	if path == "" || path == "did" {
		return fmt.Sprintf("did:acc:%s", adiLabel)
	}
	return fmt.Sprintf("did:acc:%s/%s", adiLabel, path)
}

// ExtractADILabel extracts the ADI label from a DID
func ExtractADILabel(did string) (string, error) {
	if !strings.HasPrefix(did, "did:acc:") {
		return "", fmt.Errorf("invalid DID method: %s", did)
	}

	identifier := strings.TrimPrefix(did, "did:acc:")
	parts := strings.SplitN(identifier, "/", 2)

	adiLabel := parts[0]
	if adiLabel == "" {
		return "", fmt.Errorf("empty ADI label")
	}

	return adiLabel, nil
}