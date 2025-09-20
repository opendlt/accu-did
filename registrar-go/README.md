# Accumulate DID Registrar

DID Registration service implementing DID Registration protocol for the Accumulate blockchain.

## Service Ports & Flags

**Default Port:** 8081
**Listen Address:** `--addr` (default ":8081")
**Mode:** `--real` (default: FAKE mode using mock submitter)
**Environment:** `ACC_NODE_URL` (required when using `--real`)

**Examples:**
```bash
# FAKE mode (default) - uses mock submitter
./registrar --addr :8081

# REAL mode - connects to Accumulate network
export ACC_NODE_URL=http://localhost:26657
./registrar --real --addr :8081
```

## Quick Start

### Prerequisites
- Go 1.22+
- golangci-lint (for development)

### Install Dependencies
```bash
make deps
```

### Run Registrar

**FAKE Mode (Development):**
```bash
make run
# OR directly:
go run cmd/registrar/main.go --addr :8081
```

**REAL Mode (Production):**
```bash
export ACC_NODE_URL=http://localhost:26657
go run cmd/registrar/main.go --real --addr :8081
```

The registrar will start on port 8081 by default.

### Test Endpoints

#### Health Check (Both FAKE and REAL modes)
```bash
curl http://localhost:8081/healthz
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2025-09-20T14:26:41.6380802Z"
}
```

**Note:** Health endpoint is always at `/healthz` (not `/health`)

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

**Command Line Flags:**
- `--addr`: Listen address (default: ":8081")
- `--real`: Enable real mode to connect to Accumulate network

**Environment Variables:**
- `ACC_NODE_URL`: Accumulate node URL (required when using `--real`)
- `LOG_LEVEL`: Logging level (default: info)

**Legacy Environment Variables (deprecated):**
- `REGISTRAR_PORT`: Use `--addr` flag instead

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

**Unit Tests:** All tests run with FakeSubmitter for offline development.
**Integration Tests:** Optional `*_integration_test.go` files can test against real Accumulate networks.