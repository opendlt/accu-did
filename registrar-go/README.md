# Accumulate DID Registrar

W3C DID Registration compliant service for `did:acc` method with Universal Registrar compatibility.

## Quick Start

```bash
# Run with defaults (port 8081, FAKE mode)
go run cmd/server/main.go

# Run with custom port and REAL mode
go run cmd/server/main.go --addr :8082 --mode REAL

# Set Accumulate node URL (for REAL mode)
export ACC_NODE_URL=http://localhost:26657  # Local devnet
# or
export ACC_NODE_URL=https://testnet.accumulatenetwork.io  # Testnet
```

## Configuration

| Flag/Env | Default | Description |
|----------|---------|-------------|
| `--addr` | `:8081` | Server listen address |
| `--mode` | `FAKE` | Operation mode: `FAKE` (mock) or `REAL` (Accumulate) |
| `ACC_NODE_URL` | - | Accumulate JSON-RPC endpoint (required for REAL mode) |

## Endpoints

### Health Check
```http
GET /health
```

### Native API

#### Create DID
```http
POST /create
```

**Request:**
```json
{
  "did": "did:acc:alice.acme",
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:alice.acme",
    "verificationMethod": [{
      "id": "did:acc:alice.acme#key-1",
      "type": "Ed25519VerificationKey2020",
      "controller": "did:acc:alice.acme",
      "publicKeyMultibase": "z8c7e8b4f2d1a9c5e3b7f6a8d9e2c4f1b5a7c8e9f0d2b4c6e8a1d3f5c7e9b2a4d"
    }],
    "authentication": ["did:acc:alice.acme#key-1"]
  }
}
```

**Response:**
```json
{
  "jobId": "create-1234",
  "didState": {
    "did": "did:acc:alice.acme",
    "state": "finished",
    "action": "create"
  },
  "didRegistrationMetadata": {
    "txid": "0x1234567890abcdef",
    "versionId": "1",
    "created": "2024-01-01T00:00:00Z"
  }
}
```

#### Update DID
```http
POST /update
```

**Request:**
```json
{
  "did": "did:acc:alice.acme",
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
    "service": [{
      "id": "did:acc:alice.acme#messaging",
      "type": "MessagingService",
      "serviceEndpoint": "https://alice.example.com"
    }]
  }
}
```

#### Deactivate DID
```http
POST /deactivate
```

**Request:**
```json
{
  "did": "did:acc:alice.acme"
}
```

**Response:**
```json
{
  "jobId": "deactivate-5678",
  "didState": {
    "did": "did:acc:alice.acme",
    "state": "finished",
    "action": "deactivate"
  },
  "didRegistrationMetadata": {
    "txid": "0xabcdef1234567890",
    "deactivated": true
  }
}
```

### Universal Registrar 1.0 API

#### Create (Universal)
```http
POST /1.0/create
```

**Request:**
```json
{
  "jobId": "test-create",
  "options": {
    "network": "testnet"
  },
  "secret": {},
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:alice.acme",
    "verificationMethod": [{
      "id": "did:acc:alice.acme#key-1",
      "type": "Ed25519VerificationKey2020",
      "controller": "did:acc:alice.acme",
      "publicKeyMultibase": "z8c7e8b4f2d1a9c5e3b7f6a8d9e2c4f1b5a7c8e9f0d2b4c6e8a1d3f5c7e9b2a4d"
    }]
  }
}
```

#### Update (Universal)
```http
POST /1.0/update
```

**Request with Patch:**
```json
{
  "jobId": "test-update",
  "did": "did:acc:alice.acme",
  "options": {
    "network": "testnet"
  },
  "secret": {},
  "didDocumentOperation": [
    {
      "op": "add",
      "path": "/service",
      "value": [{
        "id": "did:acc:alice.acme#messaging",
        "type": "MessagingService",
        "serviceEndpoint": "https://alice.example.com"
      }]
    }
  ]
}
```

**Request with addService/removeService:**
```json
{
  "jobId": "test-update",
  "did": "did:acc:alice.acme",
  "options": {
    "network": "testnet"
  },
  "secret": {},
  "didDocumentOperation": {
    "addService": {
      "id": "did:acc:alice.acme#messaging",
      "type": "MessagingService",
      "serviceEndpoint": "https://alice.example.com"
    }
  }
}
```

#### Deactivate (Universal)
```http
POST /1.0/deactivate
```

**Request:**
```json
{
  "jobId": "test-deactivate",
  "did": "did:acc:beastmode.acme",
  "options": {
    "network": "testnet"
  },
  "secret": {}
}
```

## FAKE vs REAL Mode

### FAKE Mode (Development)
- Uses mock submitter for instant responses
- No blockchain connection required
- Stores DIDs in memory
- Default mode for testing

### REAL Mode (Production)
- Connects to Accumulate blockchain via JSON-RPC
- Requires `ACC_NODE_URL` environment variable
- Creates actual blockchain transactions
- Maps `did:acc:name` → `acc://name/did`

## Example Workflows

### Complete DID Lifecycle
```bash
# 1. Create DID
curl -X POST http://localhost:8081/create \
  -H "Content-Type: application/json" \
  -d '{"did":"did:acc:alice.acme","didDocument":{...}}'

# 2. Update DID (add service)
curl -X POST http://localhost:8081/1.0/update \
  -H "Content-Type: application/json" \
  -d '{"did":"did:acc:alice.acme","didDocumentOperation":{"addService":{...}}}'

# 3. Resolve DID (using resolver)
curl "http://localhost:8080/resolve?did=did:acc:alice.acme"

# 4. Deactivate DID
curl -X POST http://localhost:8081/deactivate \
  -H "Content-Type: application/json" \
  -d '{"did":"did:acc:alice.acme"}'

# 5. Resolve deactivated DID
curl "http://localhost:8080/resolve?did=did:acc:alice.acme"
# Returns: {"deactivated": true}
```

## Troubleshooting

### Windows Port Issues
```powershell
# Check if port is in use
netstat -an | findstr :8081

# Kill process using port
Get-Process -Id (Get-NetTCPConnection -LocalPort 8081).OwningProcess | Stop-Process -Force

# Use alternate port
./registrar.exe --addr :8082
```

### Firewall Rules
```powershell
# Allow registrar through Windows Firewall
New-NetFirewallRule -DisplayName "DID Registrar" -Direction Inbound -LocalPort 8081 -Protocol TCP -Action Allow
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
| `invalidDid` | Wrong DID format | Use `did:acc:name` format |
| `alreadyExists` | DID already registered | Use update instead of create |
| `notFound` | DID doesn't exist | Create DID before update/deactivate |
| `unauthorized` | Missing credentials | Provide proper authentication |
| Connection refused | Node unreachable | Check `ACC_NODE_URL` |
| Port already in use | Another process on port | Use different port with `--addr` |

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│   Client    │────▶│  Registrar   │────▶│ Accumulate │
│ (curl/http) │     │   :8081      │     │    Node    │
└─────────────┘     └──────────────┘     └────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │     Mock     │
                    │ (FAKE mode)  │
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
docker build -t accu-registrar .

# Run container (FAKE mode)
docker run -p 8081:8081 accu-registrar

# Run container (REAL mode)
docker run -p 8081:8081 -e ACC_NODE_URL=http://host.docker.internal:26657 accu-registrar --mode REAL
```

## Authorization Policy

The registrar enforces the following authorization rules:
- **Create**: Requires signature from the ADI's key page (`acc://name/book/1`)
- **Update**: Requires signature from the ADI's key page
- **Deactivate**: Requires signature from the ADI's key page
- Only the ADI owner can manage their DID document

## See Also

- [Resolver Service](../resolver-go/README.md) - DID document resolution
- [Universal Registrar](https://github.com/decentralized-identity/universal-registrar) - DID registration standard
- [W3C DID Core](https://www.w3.org/TR/did-core/) - DID specification
- [Accumulate Protocol](https://accumulatenetwork.io) - Blockchain platform