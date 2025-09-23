# Accumulate DID SDK for Go

Production-grade Go SDK for interacting with Accumulate DID resolver and registrar services, providing W3C-compliant DID operations with robust error handling, retries, and idempotency support.

## Overview

This SDK provides a clean, portable HTTP client for Accumulate DID operations:
- Pure HTTP implementation against resolver/registrar APIs (OpenAPI-aligned)
- Automatic retries with exponential backoff and jitter
- Comprehensive error mapping (404→NotFound, 410→Deactivated, etc.)
- Request tracking with X-Request-Id headers
- Optional API key authentication
- Idempotency key support for safe retries
- Deterministic resolution ordering
- Canonical deactivation (410 Gone) handling

## Quick Start

### Installation

```bash
go get github.com/opendlt/accu-did/sdks/go/accdid
```

### Resolve a DID

```go
package main

import (
    "context"
    "fmt"
    "log"
    "errors"

    "github.com/opendlt/accu-did/sdks/go/accdid"
)

func main() {
    // Create resolver client
    resolver, err := accdid.NewResolverClient(accdid.ClientOptions{
        BaseURL: "http://localhost:8080",
        // Optional: APIKey: "your-api-key",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Resolve a DID
    result, err := resolver.Resolve(context.Background(), "did:acc:alice")
    if err != nil {
        // Handle specific errors
        if errors.Is(err, accdid.ErrNotFound) {
            fmt.Println("DID not found")
        } else if errors.Is(err, accdid.ErrGoneDeactivated) {
            fmt.Println("DID has been deactivated (410 Gone)")
        } else {
            log.Fatal(err)
        }
        return
    }

    fmt.Printf("DID Document: %+v\n", result.DIDDocument)
    fmt.Printf("Metadata: %+v\n", result.Metadata)
}
```

### Register a New DID

```go
registrar, err := accdid.NewRegistrarClient(accdid.ClientOptions{
    BaseURL: "http://localhost:8081",
    IdempotencyKey: "unique-request-123", // Prevent duplicate registrations
})
if err != nil {
    log.Fatal(err)
}

doc := json.RawMessage(`{
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:mycompany",
    "verificationMethod": [{
        "id": "did:acc:mycompany#key1",
        "type": "Ed25519VerificationKey2020",
        "controller": "did:acc:mycompany",
        "publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK"
    }]
}`)

txID, err := registrar.Register(context.Background(), accdid.NativeRegisterRequest{
    DID:         "did:acc:mycompany",
    DIDDocument: doc,
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Transaction ID: %s\n", txID)
```

### Update an Existing DID

```go
patch := json.RawMessage(`{
    "addService": {
        "id": "did:acc:mycompany#website",
        "type": "LinkedDomains",
        "serviceEndpoint": "https://example.com"
    }
}`)

txID, err := registrar.Update(context.Background(), accdid.NativeUpdateRequest{
    DID:   "did:acc:mycompany",
    Patch: patch,
})
```

### Deactivate a DID

```go
txID, err := registrar.Deactivate(context.Background(), accdid.NativeDeactivateRequest{
    DID:    "did:acc:mycompany",
    Reason: "Company dissolved",
})
```

## FAKE vs REAL Modes

The underlying resolver and registrar services support two modes:

**FAKE Mode** (Development/Testing):
- Services use static fixtures and mock responses
- No blockchain interaction or credit consumption
- Perfect for development and testing

**REAL Mode** (Production):
- Services connect to live Accumulate network
- Requires funded lite account with credits
- Actual blockchain transactions

The SDK works identically with both modes - the mode is determined by how the services are configured, not the SDK.

## Credits and Faucets

**Understanding Accumulate Credits:**

Accumulate uses a credit system for blockchain operations:

1. **Lite Token Account** → Holds ACME tokens
2. **ACME Tokens** → Convert to credits at oracle rate
3. **Credits** → Non-transferable, used for transactions
4. **ADI Creation** → Requires ~10 credits
5. **Data Writes** → Require ~2-5 credits per operation

**Testnet/Devnet Faucet:**
```bash
# Get free credits for testing
curl -X POST https://devnet-faucet.accumulate.io/get \
  -d '{"url":"<your-lite-account-url>"}'
```

**Mainnet:**
1. Acquire ACME tokens from exchanges
2. Convert ACME to credits using oracle rate
3. Fund your ADI key pages with credits

## Deterministic Resolution

The SDK handles deterministic resolution automatically. When multiple entries exist, they are ordered by:
1. **Sequence number** (highest wins)
2. **Timestamp** (latest wins if sequences match)
3. **SHA-256 hash** (lexicographic tiebreaker)

This ensures consistent resolution results across different clients and nodes.

## 410 Deactivation Semantics

When a DID is deactivated, the resolver returns HTTP 410 Gone with a canonical tombstone:

```json
{
  "@context": ["https://www.w3.org/ns/did/v1"],
  "id": "did:acc:example",
  "deactivated": true,
  "deactivatedAt": "2024-01-01T00:00:00Z"
}
```

The SDK maps this to `ErrGoneDeactivated` for easy handling:

```go
result, err := resolver.Resolve(ctx, did)
if errors.Is(err, accdid.ErrGoneDeactivated) {
    // Handle deactivated DID
}
```

## Advanced Configuration

### Custom Timeouts and Retries

```go
client, err := accdid.NewResolverClient(accdid.ClientOptions{
    BaseURL: "http://localhost:8080",
    Timeout: 30 * time.Second,
    Retries: accdid.RetryPolicy{
        Max:       5,
        BaseDelay: 500 * time.Millisecond,
        MaxDelay:  10 * time.Second,
        Jitter:    true,
        Backoff:   "exp",
    },
})
```

### Custom HTTP Client

```go
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:    100,
        IdleConnTimeout: 90 * time.Second,
    },
}

client, err := accdid.NewResolverClient(accdid.ClientOptions{
    BaseURL: "http://localhost:8080",
    HTTP:    httpClient,
})
```

### Request ID Generation

```go
client, err := accdid.NewResolverClient(accdid.ClientOptions{
    BaseURL: "http://localhost:8080",
    RequestID: func() string {
        return uuid.New().String()
    },
})
```

### Custom Logger

```go
type MyLogger struct{}

func (l MyLogger) Debugf(format string, args ...interface{}) {
    log.Printf("[DEBUG] "+format, args...)
}
func (l MyLogger) Infof(format string, args ...interface{}) {
    log.Printf("[INFO] "+format, args...)
}
func (l MyLogger) Warnf(format string, args ...interface{}) {
    log.Printf("[WARN] "+format, args...)
}
func (l MyLogger) Errorf(format string, args ...interface{}) {
    log.Printf("[ERROR] "+format, args...)
}

client, err := accdid.NewResolverClient(accdid.ClientOptions{
    BaseURL: "http://localhost:8080",
    Logger:  MyLogger{},
})
```

## Headers

The SDK automatically manages these headers:

| Header | Purpose | When Set |
|--------|---------|----------|
| `X-Request-Id` | Request tracking | Always (generated or custom) |
| `X-API-Key` | Authentication | When APIKey option provided |
| `Idempotency-Key` | Safe retries | When IdempotencyKey provided |
| `Content-Type` | JSON payload | On POST/PUT requests |
| `Accept` | Response format | Always (application/json) |

## Error Handling

The SDK provides consistent error mapping:

| HTTP Status | SDK Error | Description |
|-------------|-----------|-------------|
| 404 | `ErrNotFound` | DID or resource not found |
| 410 | `ErrGoneDeactivated` | DID has been deactivated |
| 400-499 | `ErrBadRequest` | Client error |
| 500-599 | `ErrServer` | Server error |
| Timeout | `ErrTimeout` | Request timeout |
| Network | `ErrNetwork` | Connection failure |

Retryable errors (429, 502, 503, 504) are automatically retried with exponential backoff.

## Versioning Policy

The SDK follows semantic versioning:
- **Major**: Breaking API changes
- **Minor**: New features, backwards compatible
- **Patch**: Bug fixes

The SDK version can be checked at runtime:
```go
fmt.Printf("SDK Version: %s\n", accdid.Version)
```

## Testing

Run tests:
```bash
go test ./...
```

With coverage:
```bash
go test -cover ./...
```

## Examples

See the [examples](examples/) directory for complete working examples:
- [basic](examples/basic/) - Simple resolve, register, update, deactivate flow
- More examples coming soon

## License

See the repository [LICENSE](../../../LICENSE) for details.