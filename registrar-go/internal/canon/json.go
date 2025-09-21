package canon

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
)

// Canonicalize returns the canonical JSON representation of a document
// following RFC8785-style canonicalization rules
func Canonicalize(v interface{}) ([]byte, error) {
	// Step 1: Marshal to JSON to validate structure
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	// Step 2: Unmarshal to generic interface for processing
	var generic interface{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	// Step 3: Apply canonicalization rules
	canonical := canonicalize(generic)

	// Step 4: Marshal with no HTML escaping and no indentation
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")

	if err := encoder.Encode(canonical); err != nil {
		return nil, fmt.Errorf("failed to encode canonical: %w", err)
	}

	// Step 5: Remove trailing newline added by encoder
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result, nil
}

// canonicalize recursively applies canonicalization rules
func canonicalize(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		// Sort keys lexicographically and recurse on values
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		result := make(map[string]interface{})
		for _, k := range keys {
			result[k] = canonicalize(val[k])
		}
		return result

	case []interface{}:
		// Preserve array order but recurse on elements
		result := make([]interface{}, len(val))
		for i, elem := range val {
			result[i] = canonicalize(elem)
		}
		return result

	default:
		// Primitive values are returned as-is
		return v
	}
}

// SHA256 computes the SHA-256 hash of data and returns it in the format "sha256:hex"
func SHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", hash)
}
