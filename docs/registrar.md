# DID Registrar API

The Accumulate DID Registrar enables DID lifecycle management operations.

## Base URL
```
http://localhost:8082
```

## Endpoints

### POST /create

Creates a new DID document.

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `did` | string | Yes | The DID to create |
| `didDocument` | object | Yes | The DID document |
| `options` | object | No | Creation options |
| `secret` | object | No | Authentication secrets |

#### Request Example

```bash
curl -X POST "http://localhost:8082/create" \
     -H "Content-Type: application/json" \
     -d '{
       "did": "did:acc:alice",
       "didDocument": {
         "@context": ["https://www.w3.org/ns/did/v1"],
         "id": "did:acc:alice",
         "verificationMethod": [
           {
             "id": "did:acc:alice#key-1",
             "type": "AccumulateKeyPage",
             "controller": "did:acc:alice",
             "keyPageUrl": "acc://alice/book/1",
             "threshold": 1
           }
         ],
         "authentication": ["#key-1"],
         "assertionMethod": ["#key-1"]
       }
     }'
```

#### Response Example

```json
{
  "jobId": "job-1705329000123456789",
  "didState": {
    "did": "did:acc:alice",
    "state": "finished",
    "action": "create"
  },
  "didRegistrationMetadata": {
    "versionId": "1705329000-abc12345",
    "contentHash": "8f434346648f6b96df89dda901c5176b10e6d8b9b1ee1e6e6e5e8e3d5c7c2e1a",
    "txId": "0x9876543210fedcba9876543210fedcba98765432"
  },
  "didDocumentMetadata": {
    "created": "2024-01-15T14:30:00Z",
    "versionId": "1705329000-abc12345"
  }
}
```

### POST /update

Updates an existing DID document.

#### Request Example

```bash
curl -X POST "http://localhost:8082/update" \
     -H "Content-Type: application/json" \
     -d '{
       "did": "did:acc:alice",
       "didDocument": {
         "@context": ["https://www.w3.org/ns/did/v1"],
         "id": "did:acc:alice",
         "verificationMethod": [
           {
             "id": "did:acc:alice#key-1",
             "type": "AccumulateKeyPage",
             "controller": "did:acc:alice",
             "keyPageUrl": "acc://alice/book/1",
             "threshold": 1
           },
           {
             "id": "did:acc:alice#key-2",
             "type": "AccumulateKeyPage",
             "controller": "did:acc:alice",
             "keyPageUrl": "acc://alice/book/2",
             "threshold": 1
           }
         ],
         "authentication": ["#key-1", "#key-2"],
         "assertionMethod": ["#key-1", "#key-2"]
       }
     }'
```

#### Response Example

```json
{
  "jobId": "job-1705329060123456789",
  "didState": {
    "did": "did:acc:alice",
    "state": "finished",
    "action": "update"
  },
  "didRegistrationMetadata": {
    "versionId": "1705329060-def67890",
    "contentHash": "7d865e959b2466918c9863afca942d0fb89d7c9ac0c99bafc3749504ded97730",
    "txId": "0x5432109876fedcba5432109876fedcba54321098"
  },
  "didDocumentMetadata": {
    "created": "2024-01-15T14:30:00Z",
    "updated": "2024-01-15T14:31:00Z",
    "versionId": "1705329060-def67890"
  }
}
```

### POST /deactivate

Deactivates a DID.

#### Request Example

```bash
curl -X POST "http://localhost:8082/deactivate" \
     -H "Content-Type: application/json" \
     -d '{
       "did": "did:acc:alice"
     }'
```

#### Response Example

```json
{
  "jobId": "job-1705329120123456789",
  "didState": {
    "did": "did:acc:alice",
    "state": "finished",
    "action": "deactivate"
  },
  "didRegistrationMetadata": {
    "versionId": "1705329120-ghi01234",
    "contentHash": "f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9",
    "txId": "0x1098765432fedcba1098765432fedcba10987654"
  },
  "didDocumentMetadata": {
    "created": "2024-01-15T14:30:00Z",
    "updated": "2024-01-15T14:32:00Z",
    "deactivated": "2024-01-15T14:32:00Z",
    "versionId": "1705329120-ghi01234"
  }
}
```

### GET /health

Health check endpoint.

#### Request Example

```bash
curl -X GET "http://localhost:8082/health"
```

#### Response Example

```json
{
  "status": "healthy",
  "service": "accu-did-registrar"
}
```

## Error Responses

### Invalid Request (400)

```json
{
  "error": "invalidRequest",
  "message": "DID is required",
  "timestamp": "2024-01-15T14:30:00Z"
}
```

### DID Mismatch (400)

```json
{
  "error": "invalidRequest",
  "message": "DID mismatch: request DID did:acc:alice does not match document ID did:acc:bob",
  "timestamp": "2024-01-15T14:30:00Z"
}
```

## Complete Test Flow

```bash
#!/bin/bash

# Variables
DID="did:acc:test$RANDOM"
BASE_URL="http://localhost:8082"

echo "Testing complete DID lifecycle for: $DID"

# 1. Create DID
echo "1. Creating DID..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/create" \
  -H "Content-Type: application/json" \
  -d "{
    \"did\": \"$DID\",
    \"didDocument\": {
      \"@context\": [\"https://www.w3.org/ns/did/v1\"],
      \"id\": \"$DID\",
      \"verificationMethod\": [
        {
          \"id\": \"$DID#key-1\",
          \"type\": \"AccumulateKeyPage\",
          \"controller\": \"$DID\",
          \"keyPageUrl\": \"acc://test$RANDOM/book/1\",
          \"threshold\": 1
        }
      ]
    }
  }")

echo "Create response: $CREATE_RESPONSE"

# 2. Update DID
echo "2. Updating DID..."
UPDATE_RESPONSE=$(curl -s -X POST "$BASE_URL/update" \
  -H "Content-Type: application/json" \
  -d "{
    \"did\": \"$DID\",
    \"didDocument\": {
      \"@context\": [\"https://www.w3.org/ns/did/v1\"],
      \"id\": \"$DID\",
      \"verificationMethod\": [
        {
          \"id\": \"$DID#key-1\",
          \"type\": \"AccumulateKeyPage\",
          \"controller\": \"$DID\",
          \"keyPageUrl\": \"acc://test$RANDOM/book/1\",
          \"threshold\": 2
        }
      ]
    }
  }")

echo "Update response: $UPDATE_RESPONSE"

# 3. Deactivate DID
echo "3. Deactivating DID..."
DEACTIVATE_RESPONSE=$(curl -s -X POST "$BASE_URL/deactivate" \
  -H "Content-Type: application/json" \
  -d "{
    \"did\": \"$DID\"
  }")

echo "Deactivate response: $DEACTIVATE_RESPONSE"
```