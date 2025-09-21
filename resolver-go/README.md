# Accumulate DID Resolver

W3C DID Core compliant resolver for `did:acc` method with Universal Resolver compatibility.

## Quick Start

```bash
# Run with defaults (port 8080, FAKE mode)
go run cmd/server/main.go

# Run with custom port and REAL mode
go run cmd/server/main.go --addr :8084 --mode REAL

# Set Accumulate node URL (for REAL mode)
export ACC_NODE_URL=http://localhost:26657  # Local devnet
# or
export ACC_NODE_URL=https://testnet.accumulatenetwork.io  # Testnet
```

## Configuration

| Flag/Env | Default | Description |
|----------|---------|-------------|
| `--addr` | `:8080` | Server listen address |
| `--mode` | `FAKE` | Operation mode: `FAKE` (fixtures) or `REAL` (Accumulate) |
| `ACC_NODE_URL` | - | Accumulate JSON-RPC endpoint (required for REAL mode) |

## Endpoints

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Native Resolve
```http
GET /resolve?did=did:acc:beastmode.acme
```

**Response:**
```json
{
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:beastmode.acme",
    "verificationMethod": [{
      "id": "did:acc:beastmode.acme#key-1",
      "type": "Ed25519VerificationKey2020",
      "controller": "did:acc:beastmode.acme",
      "publicKeyMultibase": "z8c7e8b4f2d1a9c5e3b7f6a8d9e2c4f1b5a7c8e9f0d2b4c6e8a1d3f5c7e9b2a4d"
    }]
  },
  "didDocumentMetadata": {
    "versionId": "1",
    "created": "2024-01-01T00:00:00Z",
    "updated": "2024-01-01T00:00:00Z"
  },
  "didResolutionMetadata": {
    "contentType": "application/did+json"
  }
}
```

### Universal Resolver 1.0
```http
GET /1.0/identifiers/{did}
```

**Example:**
```bash
curl http://localhost:8080/1.0/identifiers/did:acc:beastmode.acme
```

## FAKE vs REAL Mode

### FAKE Mode (Development)
- Uses fixture files from `testdata/entries/`
- No blockchain connection required
- Instant responses for testing
- Default mode

### REAL Mode (Production)
- Connects to Accumulate blockchain via JSON-RPC
- Requires `ACC_NODE_URL` environment variable
- Reads actual DID documents from data accounts
- Maps `did:acc:name` → `acc://name/did`

## Example Requests

### Resolve Active DID
```bash
# Native endpoint
curl "http://localhost:8080/resolve?did=did:acc:alice.acme"

# Universal endpoint
curl "http://localhost:8080/1.0/identifiers/did:acc:alice.acme"

# Response:
{
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:alice.acme",
    "verificationMethod": [{
      "id": "did:acc:alice.acme#key-1",
      "type": "Ed25519VerificationKey2020",
      "controller": "did:acc:alice.acme",
      "publicKeyMultibase": "z8c7e8b4f2d1a9c5e3b7f6a8d9e2c4f1b5a7c8e9f0d2b4c6e8a1d3f5c7e9b2a4d"
    }],
    "authentication": ["did:acc:alice.acme#key-1"],
    "assertionMethod": ["did:acc:alice.acme#key-1"]
  },
  "didDocumentMetadata": {
    "versionId": "2",
    "created": "2024-01-01T00:00:00Z",
    "updated": "2024-01-02T00:00:00Z"
  }
}
```

### Resolve Deactivated DID
```bash
curl "http://localhost:8080/resolve?did=did:acc:beastmode.acme"

# Response:
{
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:beastmode.acme",
    "deactivated": true
  },
  "didDocumentMetadata": {
    "deactivated": true,
    "versionId": "3",
    "updated": "2024-01-03T00:00:00Z"
  },
  "didResolutionMetadata": {
    "contentType": "application/did+json"
  }
}
```

### Invalid DID Format
```bash
curl "http://localhost:8080/resolve?did=invalid:format"

# Response:
{
  "error": "invalidDid",
  "errorMessage": "DID must start with 'did:acc:'"
}
```

### DID Not Found
```bash
curl "http://localhost:8080/resolve?did=did:acc:nonexistent"

# Response:
{
  "error": "notFound",
  "errorMessage": "DID not found"
}
```

## Troubleshooting

### Windows Port Issues
```powershell
# Check if port is in use
netstat -an | findstr :8080

# Kill process using port
Get-Process -Id (Get-NetTCPConnection -LocalPort 8080).OwningProcess | Stop-Process -Force

# Use alternate port
./resolver.exe --addr :8084
```

### Firewall Rules
```powershell
# Allow resolver through Windows Firewall
New-NetFirewallRule -DisplayName "DID Resolver" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow
```

### Node Connection Issues
```bash
# Test Accumulate node
curl $ACC_NODE_URL/status

# Use local devnet
export ACC_NODE_URL=http://localhost:26657

# Use testnet
export ACC_NODE_URL=https://testnet.accumulatenetwork.io
```

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `notFound` | DID doesn't exist | Check DID format and registration |
| `invalidDid` | Wrong DID format | Use `did:acc:name` format |
| `deactivated` | DID was deactivated | Expected for tombstoned DIDs |
| Connection refused | Node unreachable | Check `ACC_NODE_URL` |
| Port already in use | Another process on port | Use different port with `--addr` |

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│   Client    │────▶│   Resolver   │────▶│ Accumulate │
│ (curl/http) │     │   :8080      │     │    Node    │
└─────────────┘     └──────────────┘     └────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Fixtures   │
                    │  (FAKE mode) │
                    └──────────────┘
```

## Testing

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Integration test (requires node)
ACC_NODE_URL=http://localhost:26657 go test -tags=integration
```

## Docker

```bash
# Build image
docker build -t accu-resolver .

# Run container (FAKE mode)
docker run -p 8080:8080 accu-resolver

# Run container (REAL mode)
docker run -p 8080:8080 -e ACC_NODE_URL=http://host.docker.internal:26657 accu-resolver --mode REAL
```

## See Also

- [Registrar Service](../registrar-go/README.md) - DID lifecycle management
- [Universal Resolver](https://github.com/decentralized-identity/universal-resolver) - DID resolution standard
- [W3C DID Core](https://www.w3.org/TR/did-core/) - DID specification
- [Accumulate Protocol](https://accumulatenetwork.io) - Blockchain platform