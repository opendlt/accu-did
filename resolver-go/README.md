# Accumulate DID Resolver

W3C DID Core compliant resolver for the Accumulate blockchain.

## Service Ports & Flags

**Default Port:** 8080
**Listen Address:** `--addr` (default ":8080")
**Mode:** `--real` (default: FAKE mode using golden files)
**Environment:** `ACC_NODE_URL` (required when using `--real`)

**Examples:**
```bash
# FAKE mode (default) - uses golden files
./resolver --addr :8080

# REAL mode - connects to Accumulate network
export ACC_NODE_URL=http://localhost:26657
./resolver --real --addr :8080
```

## Quick Start

### Prerequisites
- Go 1.22+
- golangci-lint (for development)

### Install Dependencies
```bash
make deps
```

### Run Resolver

**FAKE Mode (Development):**
```bash
make run
# OR directly:
go run cmd/resolver/main.go --addr :8080
```

**REAL Mode (Production):**
```bash
export ACC_NODE_URL=http://localhost:26657
go run cmd/resolver/main.go --real --addr :8080
```

The resolver will start on port 8080 by default.

### Test Endpoints

#### Health Check (Both FAKE and REAL modes)
```bash
curl http://localhost:8080/healthz
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2025-09-20T14:26:41.5965281Z"
}
```

**Note:** Health endpoint is always at `/healthz` (not `/health`)

#### DID Resolution
```bash
# Basic resolution
curl "http://localhost:8080/resolve?did=did:acc:alice"

# Resolution with version time
curl "http://localhost:8080/resolve?did=did:acc:alice&versionTime=2024-01-01T00:00:00Z"

# Resolution with URL normalization
curl "http://localhost:8080/resolve?did=did:acc:ALICE"
```

Expected response format:
```json
{
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:alice",
    "verificationMethod": [...]
  },
  "didDocumentMetadata": {
    "versionId": "1704067200-8b4c4f7b",
    "created": "2024-01-01T00:00:00Z",
    "updated": "2024-01-02T00:00:00Z",
    "deactivated": false
  },
  "didResolutionMetadata": {
    "contentType": "application/did+ld+json"
  }
}
```

## Development

### Build
```bash
make build
```

### Test
```bash
make test
```

### Lint
```bash
make lint
```

### Format
```bash
make fmt
```

## API Reference

### GET /resolve

Resolves a DID to its DID document.

**Query Parameters:**
- `did` (required): The DID to resolve (e.g., `did:acc:alice`)
- `versionTime` (optional): ISO 8601 timestamp to resolve at specific time

**Response:** W3C DID Resolution Result

### GET /healthz

Health check endpoint.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Configuration

**Command Line Flags:**
- `--addr`: Listen address (default: ":8080")
- `--real`: Enable real mode to connect to Accumulate network

**Environment Variables:**
- `ACC_NODE_URL`: Accumulate node URL (required when using `--real`)
- `LOG_LEVEL`: Logging level (default: info)

**Legacy Environment Variables (deprecated):**
- `RESOLVER_PORT`: Use `--addr` flag instead

## Testing

The resolver includes comprehensive tests using golden files:
- URL normalization tests
- DID resolution scenarios
- Error handling
- Content hash validation

All tests run offline using mock data from `testdata/`.