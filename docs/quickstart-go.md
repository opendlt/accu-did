# Go Development Quickstart

Get started developing with the Accumulate DID Go SDK and services.

## Prerequisites

- Go 1.21+
- Docker (for local testing)
- Git

## Setup

### 1. Clone Repository

```bash
git clone https://github.com/opendlt/accu-did.git
cd accu-did
```

### 2. Initialize Go Workspace

```bash
go work sync
```

This synchronizes all modules in the workspace:
- `resolver-go`
- `registrar-go`
- `drivers/uniresolver-go`
- `drivers/uniregistrar-go`

### 3. Install Dependencies

```bash
go mod download
```

## Running Services

### Option 1: Go Run (Development)

=== "Resolver"

    ```bash
    cd resolver-go
    go run cmd/server/main.go
    ```

    Service available at: `http://localhost:8080`

=== "Registrar"

    ```bash
    cd registrar-go
    go run cmd/server/main.go
    ```

    Service available at: `http://localhost:8082`

=== "Universal Drivers"

    ```bash
    # Resolver Driver
    cd drivers/uniresolver-go
    go run cmd/driver/main.go

    # Registrar Driver
    cd drivers/uniregistrar-go
    go run cmd/driver/main.go
    ```

### Option 2: Docker Compose (Production-like)

```bash
docker compose up --build
```

## Testing the Setup

### 1. Health Checks

```bash
# Resolver health
curl http://localhost:8080/health

# Registrar health
curl http://localhost:8082/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "accu-did-resolver"
}
```

### 2. DID Resolution

```bash
curl "http://localhost:8080/resolve?did=did:acc:beastmode.acme" | jq .
```

Expected response structure:
```json
{
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:beastmode.acme",
    "verificationMethod": [...],
    "authentication": [...],
    "assertionMethod": [...]
  },
  "didDocumentMetadata": {
    "created": "2024-01-01T12:00:00Z",
    "versionId": "...",
    "contentHash": "..."
  }
}
```

### 3. DID Registration

```bash
curl -X POST "http://localhost:8082/create" \
     -H "Content-Type: application/json" \
     -d '{
       "did": "did:acc:test123",
       "didDocument": {
         "@context": ["https://www.w3.org/ns/did/v1"],
         "id": "did:acc:test123",
         "verificationMethod": [
           {
             "id": "did:acc:test123#key-1",
             "type": "AccumulateKeyPage",
             "controller": "did:acc:test123",
             "keyPageUrl": "acc://test123/book/1",
             "threshold": 1
           }
         ]
       }
     }' | jq .
```

## Development Workflow

### 1. Code Organization

```
accu-did/
├── resolver-go/          # DID Resolution service
│   ├── cmd/server/       # Main server entry point
│   ├── handlers/         # HTTP request handlers
│   ├── internal/         # Internal packages
│   └── tests/           # Integration tests
├── registrar-go/         # DID Registration service
│   ├── cmd/server/       # Main server entry point
│   ├── handlers/         # HTTP request handlers
│   ├── internal/         # Internal packages
│   └── tests/           # Integration tests
└── drivers/             # Universal drivers
    ├── uniresolver-go/  # Universal Resolver driver
    └── uniregistrar-go/ # Universal Registrar driver
```

### 2. Running Tests

```bash
# Test resolver
cd resolver-go
go test ./...

# Test registrar
cd registrar-go
go test ./...

# Test with race detection
go test -race ./...

# Test with coverage
go test -cover ./...
```

### 3. Building Binaries

```bash
# Build resolver
cd resolver-go
go build -o bin/resolver cmd/server/main.go

# Build registrar
cd registrar-go
go build -o bin/registrar cmd/server/main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o bin/resolver-linux cmd/server/main.go
```

## SDK Usage Examples

### Basic Resolution

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/opendlt/accu-did/resolver-go/client"
)

func main() {
    // Create resolver client
    resolver := client.New("http://localhost:8080")

    // Resolve DID
    doc, err := resolver.Resolve(context.Background(), "did:acc:beastmode.acme")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Resolved DID: %s\n", doc.ID)
    fmt.Printf("Verification Methods: %d\n", len(doc.VerificationMethod))
}
```

### DID Registration

```go
package main

import (
    "context"
    "log"

    "github.com/opendlt/accu-did/registrar-go/client"
    "github.com/opendlt/accu-did/registrar-go/types"
)

func main() {
    // Create registrar client
    registrar := client.New("http://localhost:8082")

    // Create DID document
    doc := &types.DIDDocument{
        Context: []string{"https://www.w3.org/ns/did/v1"},
        ID:      "did:acc:example",
        VerificationMethod: []types.VerificationMethod{
            {
                ID:         "did:acc:example#key-1",
                Type:       "AccumulateKeyPage",
                Controller: "did:acc:example",
                KeyPageURL: "acc://example/book/1",
                Threshold:  1,
            },
        },
    }

    // Register DID
    result, err := registrar.Create(context.Background(), "did:acc:example", doc)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created DID: %s\n", result.DIDState.DID)
    fmt.Printf("Transaction ID: %s\n", result.DIDRegistrationMetadata.TxID)
}
```

## Debugging

### Enable Debug Logging

```bash
export LOG_LEVEL=debug
go run cmd/server/main.go
```

### Profile Performance

```bash
go run -pprof=:6060 cmd/server/main.go
```

Access profiler at: `http://localhost:6060/debug/pprof/`

### Mock Accumulate Network

For testing without a real Accumulate network:

```bash
export ACCUMULATE_URL=mock
go run cmd/server/main.go
```

## Next Steps

- [API Reference](resolver.md) - Detailed API documentation
- [Universal Drivers](universal.md) - Standards-compliant interfaces
- [Interoperability](interop/didcomm.md) - DIDComm, SD-JWT, and BBS+ integration

## Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Kill process using port 8080
lsof -ti:8080 | xargs kill -9
```

#### Module Download Fails
```bash
# Clean module cache
go clean -modcache
go mod download
```

#### Build Fails
```bash
# Ensure Go workspace is synced
go work sync

# Update dependencies
go get -u ./...
go mod tidy
```