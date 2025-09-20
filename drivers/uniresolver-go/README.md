# Universal Resolver Driver for Accumulate DID

A thin proxy driver that implements the Universal Resolver Driver API and forwards requests to the Accumulate DID Resolver.

## Overview

This driver implements the [Universal Resolver Driver API](https://github.com/decentralized-identity/universal-resolver) specification for `did:acc` method DIDs.

## API

### GET /1.0/identifiers/{did}
Resolves a DID and returns the DID document along with resolution metadata.

## Configuration

Environment variables:
- `RESOLVER_URL`: URL of the resolver-go service (default: http://resolver:8080)
- `PORT`: Port to listen on (default: 8081)

## Running

### Docker Compose
```bash
docker compose up --build
```

### Testing
```powershell
./smoke.ps1
```

## Response Format

The driver returns responses in the DID Core format:
```json
{
  "didDocument": { ... },
  "didDocumentMetadata": { ... },
  "didResolutionMetadata": { ... }
}
```