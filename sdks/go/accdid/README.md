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

## Integration Tests (Devnet)

The SDK includes comprehensive integration tests that run against a live Accumulate devnet, providing end-to-end validation of the complete DID lifecycle.

### Prerequisites

1. **Start Local Devnet:**
   ```bash
   # From repository root
   make devnet-up
   ```

2. **Start DID Services in REAL Mode:**
   ```bash
   # From repository root
   make services-up
   ```

3. **Verify Services are Running:**
   ```bash
   # Check devnet status
   make devnet-status

   # Check service health
   curl http://localhost:8080/healthz  # Resolver
   curl http://localhost:8081/healthz  # Registrar
   ```

### Running Integration Tests

**Via Make (Unix/Linux):**
```bash
make test-int          # Basic integration tests
make test-int-verbose  # Verbose output
```

**Via PowerShell (Windows):**
```bash
scripts\integration.ps1           # Basic tests
scripts\integration.ps1 -Verbose  # Verbose output
```

**Manual Execution:**
```bash
cd sdks/go/accdid
export ACC_NODE_URL=http://127.0.0.1:26656
export RESOLVER_URL=http://127.0.0.1:8080
export REGISTRAR_URL=http://127.0.0.1:8081
go test -v -tags=integration ./integration
```

### What the Tests Do

The integration tests execute the complete DID lifecycle:

1. **Health Checks** - Verify resolver and registrar services are responding
2. **404 Verification** - Confirm non-existent DID returns proper error
3. **Register DID** - Create a new DID with complete document
4. **Resolve (200)** - Verify successful resolution with valid document
5. **Update (Patch)** - Add service endpoint via JSON patch
6. **Resolve Again** - Verify update was applied
7. **Deactivate** - Mark DID as deactivated with tombstone
8. **Resolve (410 Gone)** - Verify deactivation returns proper error
9. **Idempotency Test** - Test duplicate operations with same key

### Environment Variables

The integration tests respect these environment variables:

| Variable | Default | Purpose |
|----------|---------|---------|
| `ACC_NODE_URL` | *(required)* | Accumulate devnet RPC endpoint |
| `RESOLVER_URL` | `http://127.0.0.1:8080` | DID resolver service |
| `REGISTRAR_URL` | `http://127.0.0.1:8081` | DID registrar service |
| `ACC_FAUCET_URL` | *(optional)* | Devnet faucet endpoint for funding |
| `LITE_ACCOUNT_URL` | *(optional)* | Lite account for faucet funding |
| `ACCU_API_KEY` | *(optional)* | API key for authentication |
| `IDEMPOTENCY_KEY` | *(optional)* | Fixed key for idempotency tests |

### Optional Faucet Funding

If both `ACC_FAUCET_URL` and `LITE_ACCOUNT_URL` are set, the tests will attempt to fund the lite account before running:

```bash
export ACC_FAUCET_URL="http://127.0.0.1:26659/get"
export LITE_ACCOUNT_URL="acc://your-lite-account-hex"
scripts\integration.ps1
```

### Troubleshooting

**Services Not Healthy:**
- Verify devnet is running: `make devnet-status`
- Check service logs for errors
- Ensure ports 8080, 8081, 26656 are not blocked

**Insufficient Credits:**
- Use faucet funding (devnet/testnet)
- Verify lite account has sufficient credits
- Check ADI creation requirements (~10 credits)

**Wrong Ports/URLs:**
- Verify `ACC_NODE_URL` matches devnet output
- Check resolver/registrar are in REAL mode
- Confirm port forwarding if using containers

**Tests Skip or Fail:**
- Tests skip if `ACC_NODE_URL` not set (indicates FAKE mode)
- Services must be healthy before tests run
- Each test creates unique DIDs to avoid conflicts

### Sample Output

```
[INFO] Integration test configuration:
[INFO]   Resolver URL: http://127.0.0.1:8080
[INFO]   Registrar URL: http://127.0.0.1:8081
[INFO]   Accumulate Node: http://127.0.0.1:26656

[INFO] Waiting for services to be healthy...
[INFO] Both services are healthy

[INFO] Testing DID: did:acc:it1640995200000000000

[INFO] Step 1: Verify DID does not exist (expect 404)
[INFO] ✓ DID does not exist yet

[INFO] Step 2: Register new DID
[INFO] ✓ DID registered successfully (txID: a1b2c3...)

[INFO] Step 3: Resolve DID (expect 200)
[INFO] ✓ DID resolved successfully (id: did:acc:it1640995200000000000)

[INFO] Step 4: Update DID (add service)
[INFO] ✓ DID updated successfully (txID: d4e5f6...)

[INFO] Step 5: Resolve DID again (verify update)
[INFO] ✓ Service found in updated DID document: 1 services

[INFO] Step 6: Deactivate DID
[INFO] ✓ DID deactivated successfully (txID: g7h8i9...)

[INFO] Step 7: Resolve deactivated DID (expect 410 Gone)
[INFO] ✓ Deactivated DID returns expected error

[INFO] === Integration Test Summary ===
[INFO] DID: did:acc:it1640995200000000000
[INFO] Register txID: a1b2c3d4e5f6789...
[INFO] Update txID: d4e5f6g7h8i9012...
[INFO] Deactivate txID: g7h8i9j0k1l2345...
[INFO] ✓ All integration test steps completed successfully

PASS
```

## Examples

See the [examples](examples/) directory for complete working examples:
- [basic](examples/basic/) - Simple resolve, register, update, deactivate flow
- More examples coming soon

## License

See the repository [LICENSE](../../../LICENSE) for details.