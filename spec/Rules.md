# Accumulate DID Method Specification Rules

This document defines the normalization, canonicalization, and policy rules for the `did:acc` method implementation in Go.

## DID URL Normalization Rules

### 1. Method Validation
- DID method MUST be `acc`
- DID method-specific-id MUST be a valid Accumulate ADI name

### 2. ADI Name Normalization
```go
// ADI names are normalized to lowercase
func normalizeADI(adi string) string {
    return strings.ToLower(strings.TrimSuffix(adi, "."))
}
```

### 3. URL Component Handling
- **Scheme**: Always `did`
- **Method**: Always `acc`
- **Method-specific-id**: Normalized ADI name
- **Path**: Preserved as-is after first `/`
- **Query**: Key-value pairs, preserved with original casing
- **Fragment**: Preserved as-is after `#`

### 4. Normalization Algorithm
```go
type NormalizedDID struct {
    Scheme           string            // "did"
    Method           string            // "acc"
    MethodSpecificID string            // normalized ADI
    Path             string            // optional path component
    Query            map[string]string // optional query parameters
    Fragment         string            // optional fragment
}
```

## Canonical JSON + SHA-256 Algorithm

### 1. JSON Canonicalization
Following RFC 8785 (JSON Canonicalization Scheme):

```go
func CanonicalizeJSON(data interface{}) ([]byte, error) {
    // 1. Remove whitespace
    // 2. Sort object keys lexicographically
    // 3. No trailing commas
    // 4. Escape sequences normalized
    // 5. Number format normalized
    return json.Marshal(normalizeForCanonical(data))
}
```

### 2. Content Hash Algorithm
```go
func ComputeContentHash(document interface{}) (string, error) {
    canonical, err := CanonicalizeJSON(document)
    if err != nil {
        return "", err
    }

    hash := sha256.Sum256(canonical)
    return hex.EncodeToString(hash[:]), nil
}
```

### 3. Hash Properties
- **Algorithm**: SHA-256
- **Input**: Canonical JSON bytes
- **Output**: Lowercase hexadecimal string (64 characters)
- **Deterministic**: Same document always produces same hash

## Policy V1 (book/1)

### 1. Authorization Rules
Policy V1 enforces that operations on `did:acc:<adi>` must be authorized by `acc://<adi>/book/1`:

```go
func (p *PolicyV1) GetRequiredKeyPage(did string) (string, error) {
    adi, err := extractADI(did)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("acc://%s/book/1", adi), nil
}
```

### 2. Validation
```go
func (p *PolicyV1) ValidateAuthorization(did, authorKeyPage string) error {
    required, err := p.GetRequiredKeyPage(did)
    if err != nil {
        return err
    }

    if authorKeyPage != required {
        return fmt.Errorf("unauthorized: expected %s, got %s", required, authorKeyPage)
    }

    return nil
}
```

### 3. Policy Properties
- **Single key page**: Only `book/1` can authorize operations
- **Strict matching**: Exact URL match required
- **Case sensitive**: Authorization URLs are case-sensitive
- **No delegation**: No cross-ADI authorization allowed

## Version Ordering & Replay Protection

### 1. Version ID Format
```go
// Format: <unix-timestamp>-<random-suffix>
// Example: "1704067200-a1b2c3d4"
func generateVersionID(timestamp time.Time) string {
    suffix := generateRandomSuffix(8)
    return fmt.Sprintf("%d-%s", timestamp.Unix(), suffix)
}
```

### 2. Ordering Rules
- **Primary**: Unix timestamp (ascending)
- **Secondary**: Lexicographic comparison of full version ID
- **Uniqueness**: Each operation gets unique version ID

### 3. Replay Protection
```go
type EnvelopeMetadata struct {
    VersionID         string    `json:"versionId"`
    PreviousVersionID string    `json:"previousVersionId,omitempty"`
    Timestamp         time.Time `json:"timestamp"`
    AuthorKeyPage     string    `json:"authorKeyPage"`
    Proof             Proof     `json:"proof"`
}
```

### 4. Validation Algorithm
```go
func ValidateVersionOrder(current, previous string) error {
    currentTime, err := extractTimestamp(current)
    if err != nil {
        return err
    }

    if previous != "" {
        prevTime, err := extractTimestamp(previous)
        if err != nil {
            return err
        }

        if currentTime <= prevTime {
            return errors.New("version timestamp must be greater than previous")
        }
    }

    return nil
}
```

## Implementation Notes

### 1. Go Type Safety
All operations use strongly-typed structs with JSON tags for serialization.

### 2. Error Handling
Validation functions return specific error types for different failure modes.

### 3. Performance
- URL normalization is O(n) where n is URL length
- JSON canonicalization is O(m log m) where m is number of object keys
- Hash computation is O(k) where k is canonical JSON size

### 4. Security
- All inputs validated before processing
- Cryptographic operations use Go's crypto/sha256 package
- No user input directly used in system operations

## Envelope Structure

### Purpose
Envelopes wrap DID documents with metadata for registration operations.

### Structure

```go
type Envelope struct {
    ContentType string          `json:"contentType"`
    Document    json.RawMessage `json:"document"`
    Meta        EnvelopeMeta    `json:"meta"`
}

type EnvelopeMeta struct {
    VersionID         string    `json:"versionId"`
    PreviousVersionID string    `json:"previousVersionId,omitempty"`
    Timestamp         time.Time `json:"timestamp"`
    AuthorKeyPage     string    `json:"authorKeyPage"`
    Proof             Proof     `json:"proof"`
}

type Proof struct {
    Type        string `json:"type"`
    TxID        string `json:"txid"`
    ContentHash string `json:"contentHash"`
}
```

### Envelope Hash Computation

```go
func ComputeEnvelopeHash(envelope *Envelope) (string, error) {
    // Create deterministic representation
    toHash := map[string]interface{}{
        "contentType": envelope.ContentType,
        "document":    envelope.Document,
        "meta": map[string]interface{}{
            "versionId":         envelope.Meta.VersionID,
            "previousVersionId": envelope.Meta.PreviousVersionID,
            "timestamp":         envelope.Meta.Timestamp.Format(time.RFC3339),
            "authorKeyPage":     envelope.Meta.AuthorKeyPage,
        },
    }

    return ComputeHash(toHash)
}
```

## Authentication Patterns

### Policy v1: ADI Book Authorization

#### Rule
Only the Key Page at `acc://<adi>/book/1` can authorize DID operations for `did:acc:<adi>`.

#### Implementation

```go
func ValidateAuthorization(did string, authorKeyPage string) error {
    // Extract ADI from DID
    adi := extractADI(did) // "did:acc:alice" -> "alice"

    // Expected key page
    expected := fmt.Sprintf("acc://%s/book/1", adi)

    // Validate
    if authorKeyPage != expected {
        return fmt.Errorf("unauthorized: expected %s, got %s", expected, authorKeyPage)
    }

    return nil
}
```

### Signature Verification

#### Process
1. Verify transaction signatures against Key Page
2. Check threshold requirements are met
3. Validate signer permissions

```go
type AuthValidator struct {
    accumulate AccumulateClient
}

func (v *AuthValidator) ValidateSignatures(txID string, keyPageURL string) error {
    // Get transaction from Accumulate
    tx, err := v.accumulate.GetTransaction(txID)
    if err != nil {
        return err
    }

    // Get key page
    keyPage, err := v.accumulate.GetKeyPage(keyPageURL)
    if err != nil {
        return err
    }

    // Verify signatures
    validSigs := 0
    for _, sig := range tx.Signatures {
        if keyPage.ContainsKey(sig.PublicKey) {
            if verifySignature(sig, tx.Data) {
                validSigs++
            }
        }
    }

    // Check threshold
    if validSigs < keyPage.Threshold {
        return fmt.Errorf("insufficient signatures: %d < %d", validSigs, keyPage.Threshold)
    }

    return nil
}
```

## Version ID Generation

### Format
Version IDs combine timestamp and content hash for uniqueness and verifiability.

```go
func GenerateVersionID(timestamp time.Time, contentHash string) string {
    // Format: <unix-timestamp>-<hash-prefix>
    // Example: "1704067200-8b4c4f7b"

    unix := timestamp.Unix()
    hashPrefix := strings.TrimPrefix(contentHash, "sha256:")[:8]

    return fmt.Sprintf("%d-%s", unix, hashPrefix)
}
```

## URL Normalization

### DID URL Normalization Rules

1. **ADI names**: Convert to lowercase
2. **Remove trailing dots**: `alice.` → `alice`
3. **Preserve query parameters**: Order maintained
4. **Decode percent-encoding**: When safe
5. **Fragment handling**: Preserve as-is

```go
func NormalizeDIDURL(didURL string) (string, error) {
    parsed, err := url.Parse(didURL)
    if err != nil {
        return "", err
    }

    // Extract DID parts
    parts := strings.Split(parsed.Opaque, ":")
    if len(parts) != 2 || parts[0] != "acc" {
        return "", fmt.Errorf("invalid DID")
    }

    // Normalize ADI name
    adi := strings.ToLower(strings.TrimSuffix(parts[1], "."))

    // Reconstruct
    normalized := fmt.Sprintf("did:acc:%s", adi)

    // Add path, query, fragment if present
    if parsed.Path != "" {
        normalized += parsed.Path
    }
    if parsed.RawQuery != "" {
        normalized += "?" + parsed.RawQuery
    }
    if parsed.Fragment != "" {
        normalized += "#" + parsed.Fragment
    }

    return normalized, nil
}
```

## Error Response Format

### Standard Error Structure

```go
type ErrorResponse struct {
    Error      string            `json:"error"`
    Message    string            `json:"message"`
    Details    map[string]string `json:"details,omitempty"`
    RequestID  string            `json:"requestId"`
    Timestamp  time.Time         `json:"timestamp"`
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `invalid_did` | 400 | DID syntax invalid |
| `not_found` | 404 | DID does not exist |
| `unauthorized` | 403 | Insufficient permissions |
| `conflict` | 409 | Resource already exists |
| `invalid_document` | 400 | Document validation failed |
| `threshold_not_met` | 403 | Signature threshold not met |
| `internal_error` | 500 | Server error |

## Validation Rules

### DID Document Validation

```go
func ValidateDIDDocument(doc map[string]interface{}) error {
    // Required fields
    if _, ok := doc["@context"]; !ok {
        return fmt.Errorf("missing @context")
    }

    id, ok := doc["id"].(string)
    if !ok {
        return fmt.Errorf("missing or invalid id")
    }

    // Validate DID format
    if !strings.HasPrefix(id, "did:acc:") {
        return fmt.Errorf("invalid DID format")
    }

    // Validate verification methods if present
    if vms, ok := doc["verificationMethod"].([]interface{}); ok {
        for _, vm := range vms {
            if err := validateVerificationMethod(vm); err != nil {
                return err
            }
        }
    }

    return nil
}
```

### AccumulateKeyPage Validation

```go
func validateAccumulateKeyPage(vm map[string]interface{}) error {
    // Required fields
    required := []string{"id", "type", "controller", "keyPageUrl", "threshold"}

    for _, field := range required {
        if _, ok := vm[field]; !ok {
            return fmt.Errorf("missing field: %s", field)
        }
    }

    // Type must be AccumulateKeyPage
    if vm["type"] != "AccumulateKeyPage" {
        return fmt.Errorf("invalid type")
    }

    // Validate keyPageUrl format
    keyPageUrl, ok := vm["keyPageUrl"].(string)
    if !ok || !strings.HasPrefix(keyPageUrl, "acc://") {
        return fmt.Errorf("invalid keyPageUrl")
    }

    // Threshold must be positive
    threshold, ok := vm["threshold"].(float64)
    if !ok || threshold < 1 {
        return fmt.Errorf("invalid threshold")
    }

    return nil
}
```

## Test Vector Categories

### 1. Canonical JSON Tests
- Object key ordering
- Array preservation
- Number formatting
- String escaping
- Nested structures

### 2. Hash Computation Tests
- Empty documents
- Complex nested structures
- Special characters
- Large documents

### 3. Authorization Tests
- Valid authorization
- Invalid key page
- Threshold not met
- Missing signatures

### 4. Envelope Tests
- Complete envelopes
- Missing metadata
- Invalid proofs
- Hash verification

## Security Considerations

### Hash Collision Resistance
- SHA-256 provides 2^128 collision resistance
- Monitor for SHA-256 deprecation advisories

### Canonicalization Attacks
- Reject duplicate keys
- Validate JSON structure before canonicalization
- Use strict parsing

### Timing Attacks
- Use constant-time comparison for hashes
- Avoid early returns in validation

### Replay Protection
- Include timestamp in envelopes
- Track processed version IDs
- Reject old timestamps

## Implementation Notes

### Performance Optimization
1. Cache canonical representations
2. Pre-compute hashes where possible
3. Use streaming for large documents

### Compatibility
1. Support both compressed and pretty-printed input
2. Handle different timestamp formats
3. Accept various hash representations

### Testing
1. Use table-driven tests
2. Include edge cases
3. Validate against test vectors
4. Fuzz testing for parsers

---

*This document defines the authoritative encoding and authentication rules for the Accumulate DID implementation.*