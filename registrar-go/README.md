# Accumulate DID Registrar

DID Registration service implementing DID Registration protocol for the Accumulate blockchain.

## Quick Start

### Prerequisites
- Go 1.22+
- golangci-lint (for development)

### Install Dependencies
```bash
make deps
```

### Run Registrar
```bash
make run
```

The registrar will start on port 8081.

### Test Endpoints

#### Health Check
```bash
curl http://localhost:8081/healthz
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### DID Creation
```bash
curl -X POST http://localhost:8081/create \
  -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:alice",
    "didDocument": {
      "@context": ["https://www.w3.org/ns/did/v1"],
      "id": "did:acc:alice",
      "verificationMethod": [{
        "id": "did:acc:alice#key-1",
        "type": "AccumulateKeyPage",
        "controller": "did:acc:alice",
        "keyPageUrl": "acc://alice/book/1",
        "threshold": 1
      }],
      "authentication": ["did:acc:alice#key-1"]
    }
  }'
```

#### DID Update
```bash
curl -X POST http://localhost:8081/update \
  -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:alice",
    "didDocument": {
      "@context": ["https://www.w3.org/ns/did/v1"],
      "id": "did:acc:alice",
      "verificationMethod": [{
        "id": "did:acc:alice#key-1",
        "type": "AccumulateKeyPage",
        "controller": "did:acc:alice",
        "keyPageUrl": "acc://alice/book/1",
        "threshold": 1
      }],
      "service": [{
        "id": "did:acc:alice#messaging",
        "type": "MessagingService",
        "serviceEndpoint": "https://messaging.alice.example.com"
      }]
    }
  }'
```

#### DID Deactivation
```bash
curl -X POST http://localhost:8081/deactivate \
  -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:alice"
  }'
```

Expected response format:
```json
{
  "jobId": "uuid-1234-5678-9abc",
  "didState": {
    "did": "did:acc:alice",
    "state": "finished",
    "action": "create"
  },
  "didRegistrationMetadata": {
    "versionId": "1704067200-8b4c4f7b",
    "contentHash": "sha256:abc123...",
    "txid": "0x1234567890abcdef"
  },
  "didDocumentMetadata": {
    "created": "2024-01-01T00:00:00Z",
    "versionId": "1704067200-8b4c4f7b"
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

### POST /create

Creates a new DID document.

**Request Body:**
- `did` (required): The DID to create (e.g., `did:acc:alice`)
- `didDocument` (required): The DID document to create
- `options` (optional): Registration options

**Response:** DID Registration Result

### POST /update

Updates an existing DID document.

**Request Body:**
- `did` (required): The DID to update
- `didDocument` (required): The updated DID document
- `options` (optional): Registration options

**Response:** DID Registration Result

### POST /deactivate

Deactivates a DID.

**Request Body:**
- `did` (required): The DID to deactivate
- `options` (optional): Registration options

**Response:** DID Registration Result

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
- `REGISTRAR_PORT`: Server port (default: 8081)
- `LOG_LEVEL`: Logging level (default: info)
- `ACCUMULATE_API_URL`: Accumulate API endpoint (for production)

## Authorization Policy

By default, the registrar enforces Policy v1:
- Only the Key Page at `acc://<adi>/book/1` can authorize DID operations
- This ensures that only the ADI owner can create/update/deactivate their DID

## Testing

The registrar includes comprehensive tests:
- DID document validation
- Envelope generation and signing
- Authorization policy enforcement
- Error handling scenarios

All tests run with mock Accumulate client for offline development.