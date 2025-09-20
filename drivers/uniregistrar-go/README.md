# Universal Registrar Driver for Accumulate DID

A proxy driver that implements the Universal Registrar Driver API and forwards requests to the Accumulate DID Registrar.

## Overview

This driver implements the [Universal Registrar Driver API](https://github.com/decentralized-identity/universal-registrar) specification for `did:acc` method DIDs.

## API Endpoints

- `POST /1.0/create?method=acc` - Create a new DID
- `POST /1.0/update?method=acc` - Update an existing DID
- `POST /1.0/deactivate?method=acc` - Deactivate a DID

## Configuration

Environment variables:
- `REGISTRAR_URL`: URL of the registrar-go service (default: http://registrar:8082)
- `PORT`: Port to listen on (default: 8083)

## Running

### Docker Compose
```bash
docker compose up --build
```

### Testing
```powershell
./smoke.ps1
```

## Request/Response Format

### Create Request
```json
{
  "did": "did:acc:alice",
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:alice",
    "verificationMethod": [...]
  },
  "options": {},
  "secret": {}
}
```

### Response
```json
{
  "jobId": "...",
  "didState": {
    "did": "did:acc:alice",
    "state": "finished",
    "action": "create"
  },
  "didRegistrationMetadata": {
    "versionId": "...",
    "contentHash": "...",
    "txId": "..."
  },
  "didDocumentMetadata": {
    "created": "...",
    "versionId": "..."
  }
}
```