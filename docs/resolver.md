# DID Resolver API

The Accumulate DID Resolver provides standards-compliant DID document resolution.

## Base URL
```
http://localhost:8080
```

## Endpoints

### GET /resolve

Resolves a DID to its DID document.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `did` | string | Yes | The DID to resolve |
| `versionTime` | string | No | ISO 8601 timestamp for historical resolution |
| `transform` | string | No | Response transformation (`jsonld`) |

#### Request Example

=== "Basic Resolution"

    ```bash
    curl -X GET "http://localhost:8080/resolve?did=did:acc:beastmode.acme" \
         -H "Accept: application/json"
    ```

=== "Historical Resolution"

    ```bash
    curl -X GET "http://localhost:8080/resolve?did=did:acc:beastmode.acme&versionTime=2024-01-01T00:00:00Z" \
         -H "Accept: application/json"
    ```

=== "JSON-LD Transform"

    ```bash
    curl -X GET "http://localhost:8080/resolve?did=did:acc:beastmode.acme&transform=jsonld" \
         -H "Accept: application/ld+json"
    ```

#### Response Example

```json
{
  "didDocument": {
    "@context": [
      "https://www.w3.org/ns/did/v1",
      "https://w3id.org/security/v1"
    ],
    "id": "did:acc:beastmode.acme",
    "verificationMethod": [
      {
        "id": "did:acc:beastmode.acme#key-1",
        "type": "AccumulateKeyPage",
        "controller": "did:acc:beastmode.acme",
        "keyPageUrl": "acc://beastmode.acme/book/1",
        "threshold": 1
      }
    ],
    "authentication": ["#key-1"],
    "assertionMethod": ["#key-1"],
    "service": [
      {
        "id": "did:acc:beastmode.acme#resolver",
        "type": "DIDResolver",
        "serviceEndpoint": "https://resolver.accumulate.defi"
      }
    ]
  },
  "didDocumentMetadata": {
    "created": "2024-01-01T12:00:00Z",
    "updated": "2024-01-15T14:30:00Z",
    "versionId": "1705329000-7c8d9e0f",
    "contentHash": "2c624232cdd221771294dfbb310aca000a0df6ac8b66b696d90ef06fdefb64a3",
    "txId": "0x1234567890abcdef1234567890abcdef12345678"
  }
}
```

### GET /health

Health check endpoint.

#### Request Example

```bash
curl -X GET "http://localhost:8080/health"
```

#### Response Example

```json
{
  "status": "healthy",
  "service": "accu-did-resolver",
  "timestamp": "2024-01-15T14:30:00Z"
}
```

## Error Responses

### DID Not Found (404)

```json
{
  "error": "notFound",
  "message": "DID document not found",
  "did": "did:acc:nonexistent"
}
```

### Invalid DID Format (400)

```json
{
  "error": "invalidDid",
  "message": "Invalid DID format",
  "details": {
    "expected": "did:acc:<adi>",
    "received": "did:key:invalid"
  }
}
```

### DID Deactivated (410 Gone)

When a DID has been deactivated, the resolver returns HTTP 410 Gone with deactivation metadata:

**Headers:**
```
HTTP/1.1 410 Gone
Content-Type: application/did+json
```

**Response Body:**
```json
{
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:beastmode.acme",
    "deactivated": true,
    "deactivatedAt": "2024-01-15T10:30:00Z"
  },
  "didDocumentMetadata": {
    "deactivated": true,
    "versionId": "1705329000-deactivated",
    "updated": "2024-01-15T10:30:00Z",
    "contentHash": "8f434346648f6b96df89dda901c5176b10e6d8b9b1ee1e6e6e5e8e3d5c7c2e1a"
  },
  "didResolutionMetadata": {
    "contentType": "application/did+json"
  }
}

## Testing with Fixtures

### Available Test DIDs

- `did:acc:alice` - Simple test identity
- `did:acc:beastmode.acme` - Complex hierarchical identity

### Quick Test Script

```bash
#!/bin/bash
echo "Testing DID Resolution..."

# Test 1: Basic resolution
echo "1. Basic resolution:"
curl -s "http://localhost:8080/resolve?did=did:acc:beastmode.acme" | jq .

# Test 2: Health check
echo "2. Health check:"
curl -s "http://localhost:8080/health" | jq .

# Test 3: Alice DID
echo "3. Alice DID:"
curl -s "http://localhost:8080/resolve?did=did:acc:alice" | jq .
```

## Performance

- **Resolution Time**: < 100ms typical
- **Throughput**: 1000+ requests/second
- **Cache TTL**: 5 minutes default