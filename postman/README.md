# Accumulate DID Postman Collection

Complete API testing suite for the Accumulate DID method implementation.

## Import Instructions

1. **Import Environment:**
   - Open Postman
   - Click "Import" → "Upload Files"
   - Select `local.postman_environment.json`
   - Activate the "accu-did-local" environment

2. **Import Collection:**
   - Click "Import" → "Upload Files"
   - Select `accu-did.postman_collection.json`

## Setup

### 1. Start Services (REAL Mode)

```bash
# Set environment variable
export ACC_NODE_URL=http://localhost:26660

# Start resolver (terminal 1)
cd resolver-go && go run cmd/server/main.go --addr :8080 --mode REAL

# Start registrar (terminal 2)
cd registrar-go && go run cmd/server/main.go --addr :8081 --mode REAL
```

### 2. Replace Key Placeholder

Before running create requests:
1. Generate an Ed25519 public key or use existing one
2. Replace `z6Mk...REPLACE...with-real-key` in request bodies with actual multibase-encoded public key

**Example valid key:** `z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK`

## Request List

### Health Checks
- **Resolver Health** - `GET /healthz` - Verify resolver service
- **Registrar Health** - `GET /healthz` - Verify registrar service

### Resolution
- **Resolve DID** - `GET /resolve?did={{DID}}` - Get DID document

### Native API
- **Create DID (Native)** - `POST /register` - Create new DID with full document
- **Update DID (Native)** - `POST /update` - Add LinkedDomains service
- **Deactivate DID (Native)** - `POST /deactivate` - Tombstone the DID

### Universal Registrar API
- **Create DID (Universal)** - `POST /1.0/create` - DIF-compatible create
- **Update DID (Universal)** - `POST /1.0/update` - Patch-based service addition
- **Deactivate DID (Universal)** - `POST /1.0/deactivate` - DIF-compatible deactivate

## Suggested Test Flow

1. **Health Checks** - Verify both services are running
2. **Create DID** - Use either Native or Universal endpoint
3. **Resolve DID** - Verify creation worked (status 200, active)
4. **Update DID** - Add a service to the DID document
5. **Resolve DID** - Verify update worked (service present)
6. **Deactivate DID** - Tombstone the DID
7. **Resolve DID** - Verify deactivation (status 410 or deactivated:true)

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `RESOLVER_URL` | `http://127.0.0.1:8080` | Resolver service endpoint |
| `REGISTRAR_URL` | `http://127.0.0.1:8081` | Registrar service endpoint |
| `DID` | `did:acc:beastmode.acme` | Test DID identifier |
| `ACC_NODE_URL` | `http://127.0.0.1:26660` | Accumulate node endpoint |

## Notes

- All requests include basic test assertions
- DID variable is automatically updated from successful responses
- Replace the placeholder public key before running create operations
- Ensure Accumulate devnet is running for REAL mode testing
- Use FAKE mode for offline testing without blockchain dependency