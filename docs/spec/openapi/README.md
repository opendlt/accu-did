# OpenAPI Specification for Accu-DID Resolver

Complete OpenAPI 3.1.0 specification for the Accumulate DID Resolver service.

## Preview the Specification

### Online Editors
- **Swagger Editor**: https://editor.swagger.io/
  1. Copy contents of `resolver.yaml`
  2. Paste into the editor
  3. View interactive documentation

- **Redocly**: https://redocly.github.io/redoc/
  1. Upload `resolver.yaml` file
  2. Generate beautiful documentation

### Local Preview
```bash
# Using Swagger UI (Docker)
docker run -p 8081:8080 -v $(pwd):/tmp -e SWAGGER_JSON=/tmp/resolver.yaml swaggerapi/swagger-ui

# Using Redoc CLI
npm install -g redoc-cli
redoc-cli serve resolver.yaml --watch --port 8081
```

## Example Usage

### Health Check
```bash
curl http://127.0.0.1:8080/healthz
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-21T12:00:00Z"
}
```

### Resolve DID
```bash
# Basic resolution
curl "http://127.0.0.1:8080/resolve?did=did:acc:beastmode.acme"

# Universal Resolver endpoint
curl "http://127.0.0.1:8080/1.0/identifiers/did:acc:beastmode.acme"

# Historical resolution (optional)
curl "http://127.0.0.1:8080/resolve?did=did:acc:alice&versionTime=2024-01-01T00:00:00Z"
```

**Success Response (200):**
```json
{
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:beastmode.acme",
    "verificationMethod": [{
      "id": "did:acc:beastmode.acme#key-1",
      "type": "Ed25519VerificationKey2020",
      "controller": "did:acc:beastmode.acme",
      "publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK"
    }],
    "authentication": ["did:acc:beastmode.acme#key-1"]
  },
  "didDocumentMetadata": {
    "versionId": "2",
    "created": "2024-01-01T00:00:00Z",
    "updated": "2024-01-02T00:00:00Z"
  },
  "didResolutionMetadata": {
    "contentType": "application/did+json"
  }
}
```

**Error Response (404):**
```json
{
  "error": "notFound",
  "errorMessage": "DID not found"
}
```

## Service Configuration

### Ports and Modes

| Mode | Description | Port | Dependencies |
|------|-------------|------|--------------|
| **FAKE** | Offline testing with fixtures | 8080 | None |
| **REAL** | Live Accumulate blockchain | 8080 | `ACC_NODE_URL` |

### Environment Variables

- **`ACC_NODE_URL`**: Accumulate JSON-RPC endpoint (required for REAL mode)
  - Local devnet: `http://localhost:26657`
  - Testnet: `https://testnet.accumulatenetwork.io`

### Starting the Service

```bash
# FAKE mode (default)
cd resolver-go
go run cmd/server/main.go --addr :8080

# REAL mode
export ACC_NODE_URL=http://localhost:26657
go run cmd/server/main.go --addr :8080 --mode REAL
```

## API Features

### W3C DID Core Compliance
- Standard DID resolution result format
- Proper error codes and messages
- Support for DID Document metadata
- Historical resolution (versionTime parameter)

### Universal Resolver Compatibility
- `/1.0/identifiers/{did}` endpoint
- DIF Universal Resolver patterns
- Compatible response formats

### Accumulate-Specific Features
- Maps `did:acc:<adi>` to `acc://<adi>/did`
- Supports custom paths: `did:acc:company/docs`
- Handles deactivated DIDs properly
- Returns Accumulate transaction metadata

## Testing

Use the [Postman collection](../../../postman/) for comprehensive API testing:

1. Import `postman/accu-did.postman_collection.json`
2. Import `postman/local.postman_environment.json`
3. Run the "Health" and "Resolve" request groups

## Registrar API

### Import Registrar Specification

Follow the same process as the resolver:

1. **Import into Swagger Editor**: Copy `registrar.yaml` contents
2. **Import into Postman**:
   ```bash
   # Generate Postman collection from OpenAPI
   openapi2postman -s registrar.yaml -o registrar-postman.json
   ```

### Example Usage - Native API

#### Create DID
```bash
curl -X POST http://127.0.0.1:8081/register \
  -H "Content-Type: application/json" \
  -d '{
    "didDocument": {
      "@context": ["https://www.w3.org/ns/did/v1"],
      "id": "did:acc:beastmode.acme",
      "verificationMethod": [{
        "id": "did:acc:beastmode.acme#key-1",
        "type": "Ed25519VerificationKey2020",
        "controller": "did:acc:beastmode.acme",
        "publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK"
      }],
      "authentication": ["did:acc:beastmode.acme#key-1"]
    }
  }'
```

#### Update DID (Add Service)
```bash
curl -X POST http://127.0.0.1:8081/update \
  -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:beastmode.acme",
    "patch": {
      "addService": [{
        "id": "did:acc:beastmode.acme#messaging",
        "type": "MessagingService",
        "serviceEndpoint": "https://messaging.example.com"
      }]
    }
  }'
```

#### Deactivate DID
```bash
curl -X POST http://127.0.0.1:8081/deactivate \
  -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:beastmode.acme",
    "deactivate": true
  }'
```

### Example Usage - Universal Registrar API

#### Create DID (Universal)
```bash
curl -X POST "http://127.0.0.1:8081/1.0/create?method=acc" \
  -H "Content-Type: application/json" \
  -d '{
    "jobId": "test-create-001",
    "options": {
      "network": "devnet"
    },
    "secret": {},
    "registration": {
      "didDocument": {
        "@context": ["https://www.w3.org/ns/did/v1"],
        "id": "did:acc:beastmode.acme",
        "verificationMethod": [{
          "id": "did:acc:beastmode.acme#key-1",
          "type": "Ed25519VerificationKey2020",
          "controller": "did:acc:beastmode.acme",
          "publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK"
        }]
      }
    }
  }'
```

#### Update DID (Universal)
```bash
curl -X POST "http://127.0.0.1:8081/1.0/update?method=acc" \
  -H "Content-Type: application/json" \
  -d '{
    "jobId": "test-update-001",
    "registration": {
      "did": "did:acc:beastmode.acme",
      "patch": {
        "addService": {
          "id": "did:acc:beastmode.acme#messaging",
          "type": "MessagingService",
          "serviceEndpoint": "https://messaging.example.com"
        }
      }
    }
  }'
```

### Registrar Configuration

**Important**: REAL mode requires `ACC_NODE_URL` environment variable:

```bash
# FAKE mode (default)
cd registrar-go
go run cmd/server/main.go --addr :8081

# REAL mode
export ACC_NODE_URL=http://localhost:26657
go run cmd/server/main.go --addr :8081 --mode REAL
```

**Ports:** Registrar runs on port 8081 by default (resolver uses 8080).

## See Also

- [DID Method Specification](../method.md) - Complete `did:acc` method definition
- [Registrar OpenAPI](./registrar.yaml) - DID registration service API
- [W3C DID Core](https://www.w3.org/TR/did-core/) - DID specification
- [Universal Resolver](https://github.com/decentralized-identity/universal-resolver) - Resolution standard