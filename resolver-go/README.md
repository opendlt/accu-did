# Accumulate DID Resolver

W3C DID Core compliant resolver for the Accumulate blockchain.

## Quick Start

### Prerequisites
- Go 1.22+
- golangci-lint (for development)

### Install Dependencies
```bash
make deps
```

### Run Resolver
```bash
make run
```

The resolver will start on port 8080.

### Test Endpoints

#### Health Check
```bash
curl http://localhost:8080/healthz
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

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

Environment variables:
- `RESOLVER_PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (default: info)
- `ACCUMULATE_API_URL`: Accumulate API endpoint (for production)

## Testing

The resolver includes comprehensive tests using golden files:
- URL normalization tests
- DID resolution scenarios
- Error handling
- Content hash validation

All tests run offline using mock data from `testdata/`.