# Go SDK Quick Start

Get started with the Accumulate DID Go SDK in minutes.

## Installation

```bash
go get github.com/opendlt/accu-did/sdks/go/accdid
```

## Basic Usage

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
    })
    if err != nil {
        log.Fatal(err)
    }

    // Resolve a DID
    result, err := resolver.Resolve(context.Background(), "did:acc:alice")
    if err != nil {
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
}
```

### Handle 410 Deactivated DIDs

The SDK automatically handles deactivated DIDs that return HTTP 410 Gone:

```go
result, err := resolver.Resolve(ctx, "did:acc:deactivated")
if err != nil {
    if errors.Is(err, accdid.ErrGoneDeactivated) {
        fmt.Println("This DID has been deactivated")

        // You might still get metadata about the deactivation
        if result != nil && result.DIDDocument != nil {
            fmt.Println("Deactivation tombstone:", result.DIDDocument)
        }
        return
    }
    // Handle other errors...
}
```

The canonical deactivation tombstone format is:

```json
{
  "@context": ["https://www.w3.org/ns/did/v1"],
  "id": "did:acc:example",
  "deactivated": true,
  "deactivatedAt": "2024-01-01T00:00:00Z"
}
```

### Register a New DID

```go
// Create registrar client
registrar, err := accdid.NewRegistrarClient(accdid.ClientOptions{
    BaseURL: "http://localhost:8081",
    IdempotencyKey: "unique-request-123", // Prevent duplicate registrations
})
if err != nil {
    log.Fatal(err)
}

// Create DID document
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

// Register the DID
txID, err := registrar.Register(context.Background(), accdid.NativeRegisterRequest{
    DID:         "did:acc:mycompany",
    DIDDocument: doc,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Registration successful, Transaction ID: %s\n", txID)
```

## Error Handling

The SDK provides structured error handling:

```go
result, err := resolver.Resolve(ctx, did)
if err != nil {
    switch {
    case errors.Is(err, accdid.ErrNotFound):
        // 404 - DID not found
        fmt.Println("DID does not exist")

    case errors.Is(err, accdid.ErrGoneDeactivated):
        // 410 - DID has been deactivated
        fmt.Println("DID has been deactivated")

    case errors.Is(err, accdid.ErrBadRequest):
        // 4xx - Client error (invalid DID format, etc.)
        fmt.Printf("Invalid request: %v\n", err)

    case errors.Is(err, accdid.ErrServer):
        // 5xx - Server error
        fmt.Printf("Server error: %v\n", err)

    case errors.Is(err, accdid.ErrTimeout):
        // Request timeout
        fmt.Printf("Request timeout: %v\n", err)

    case errors.Is(err, accdid.ErrNetwork):
        // Network connectivity issue
        fmt.Printf("Network error: %v\n", err)

    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
}
```

## Configuration

### Timeouts and Retries

```go
client, err := accdid.NewResolverClient(accdid.ClientOptions{
    BaseURL: "http://localhost:8080",
    Timeout: 30 * time.Second,
    Retries: accdid.RetryPolicy{
        Max:       5,
        BaseDelay: 500 * time.Millisecond,
        MaxDelay:  10 * time.Second,
        Jitter:    true,
    },
})
```

### API Key Authentication

```go
client, err := accdid.NewResolverClient(accdid.ClientOptions{
    BaseURL: "http://localhost:8080",
    APIKey:  "your-api-key-here",
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

## Complete Example

See [`examples/basic/main.go`](../sdks/go/accdid/examples/basic/main.go) for a complete working example that demonstrates:

- Service health checks
- DID resolution with error handling
- Complete DID lifecycle (register → update → deactivate)
- 410 Gone handling for deactivated DIDs
- Environment variable configuration

## Testing

```bash
# Run all SDK tests
make sdk-test

# Run example
make example-sdk
```

## Next Steps

- Read the [complete SDK documentation](../sdks/go/accdid/README.md)
- Explore [Universal Resolver compatibility](universal.md)
- Learn about [credit requirements](ops/OPERATIONS.md#credits-required)