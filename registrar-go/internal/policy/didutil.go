package policy

import (
	"fmt"
	"strings"
)

func DIDToDataAccountURL(did string) (string, error) {
	const prefix = "did:acc:"
	if !strings.HasPrefix(did, prefix) {
		return "", fmt.Errorf("invalid DID (missing did:acc: prefix)")
	}
	rest := strings.TrimPrefix(did, prefix)
	for _, sep := range []rune{'/', '?', '#', ';'} {
		if i := strings.IndexRune(rest, sep); i >= 0 {
			rest = rest[:i]
		}
	}
	if rest == "" {
		return "", fmt.Errorf("invalid DID (empty identifier)")
	}
	return fmt.Sprintf("acc://%s/data/did", rest), nil
}
