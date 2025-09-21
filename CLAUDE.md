# Accumulate DID Implementation

Decentralized identifiers on Accumulate Protocol with Go services and Universal drivers.

## Quick Start

**Prerequisites:** Go 1.22+, Docker

**Key Services:**
- Resolver: `:8080` - DID document resolution with proper did:acc → acc:// URL mapping
- Registrar: `:8081` - DID lifecycle management with native and Universal Registrar endpoints

**Most Common Commands:**
```bash
# Build & test everything
make build test

# FAKE Mode (default) - offline development
cd resolver-go && go run cmd/resolver/main.go --addr :8080
cd registrar-go && go run cmd/registrar/main.go --addr :8081

# REAL Mode - connect to Accumulate network
export ACC_NODE_URL=http://localhost:26657
cd resolver-go && go run cmd/resolver/main.go --real --addr :8080
cd registrar-go && go run cmd/registrar/main.go --real --addr :8081
```

## Tech Stack

- **Framework:** chi router, golangci-lint
- **Module Base:** `github.com/opendlt/accu-did`
- **Standards:** W3C DID Core, Universal Resolver/Registrar
- **DID Method:** `did:acc:<adi-label>[/<path>]` maps to `acc://<adi-label>/<path|did>`

## DID Conventions

The Accumulate DID method follows the pattern `did:acc:<adi-label>[/<path>]`:

- **Simple DID:** `did:acc:alice` → ADI: `acc://alice`, Data Account: `acc://alice/did`
- **DID with dots:** `did:acc:beastmode.acme` → ADI: `acc://beastmode.acme`, Data Account: `acc://beastmode.acme/did`
- **Custom path:** `did:acc:alice/documents` → ADI: `acc://alice`, Data Account: `acc://alice/documents`

**Registration Flow (Native endpoint):**
1. Create ADI: `acc://<adi-label>` with specified key page
2. Create data account: `acc://<adi-label>/<path>` (default: "did")
3. Write DID document as JSON data to the data account

**Resolution Flow:**
1. Parse DID to extract ADI and data account URLs
2. Read JSON data from data account: `acc://<adi-label>/<path>`
3. Return W3C DID Resolution Result with metadata

## Development

Use `make help` for all targets. Key services run on localhost with health checks at `/healthz`.

**Operation Modes:**
- **FAKE Mode (default):** Uses golden files (resolver) and mock submitter (registrar) for offline development
- **REAL Mode:** Connects to live Accumulate network via `ACC_NODE_URL` when using `--real` flag

**Universal Drivers:** Standards-compliant proxies in `drivers/` directory.

## Documentation

- **API Docs:** `docs/` (MkDocs)
- **Deep Context:** [.claude/memory/CLAUDE.md](.claude/memory/CLAUDE.md)
- **Interop Plans:** `docs/interop/`

## API Endpoints

### Resolver (Port 8080)
- `GET /resolve?did={did}` - Resolve DID documents using proper DID→URL mapping
- `GET /healthz` - Health check endpoint

### Registrar (Port 8081)

**Native Endpoints (Clean Internal API):**
- `POST /register` - Register new DID with full ADI/data account creation
- `POST /native/update` - Update existing DID document
- `POST /native/deactivate` - Deactivate DID

**Universal Registrar v1.0 Compatibility:**
- `POST /1.0/create` - Universal Registrar create endpoint
- `POST /1.0/update` - Universal Registrar update endpoint
- `POST /1.0/deactivate` - Universal Registrar deactivate endpoint

**Legacy Universal Registrar v0.x:**
- `POST /create` - Legacy create endpoint
- `POST /update` - Legacy update endpoint
- `POST /deactivate` - Legacy deactivate endpoint

## Quick Test

### FAKE Mode Smoke Tests
```bash
# Start services in FAKE mode
./resolver --addr :8080 &
./registrar --addr :8081 &

# Health checks
curl http://localhost:8080/healthz
curl http://localhost:8081/healthz

# DID resolution (uses testdata files)
curl "http://localhost:8080/resolve?did=did:acc:alice"
curl "http://localhost:8080/resolve?did=did:acc:beastmode.acme"

# Native DID registration (creates ADI + data account + writes DID doc)
curl -X POST -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:testuser",
    "didDocument": {
      "@context": ["https://www.w3.org/ns/did/v1"],
      "id": "did:acc:testuser",
      "verificationMethod": [{
        "id": "did:acc:testuser#key1",
        "type": "Ed25519VerificationKey2020",
        "controller": "did:acc:testuser",
        "publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK"
      }]
    }
  }' \
  http://localhost:8081/register

# Universal Registrar compatibility
curl -X POST -H "Content-Type: application/json" \
  -d '{
    "didDocument": {
      "@context": ["https://www.w3.org/ns/did/v1"],
      "id": "did:acc:universaltest"
    }
  }' \
  http://localhost:8081/1.0/create
```

### REAL Mode Smoke Tests
```bash
# Set environment and start services in REAL mode
export ACC_NODE_URL=http://localhost:26657
./resolver --addr :8080 --real &
./registrar --addr :8081 --real &

# Health checks (same as FAKE mode)
curl http://localhost:8080/healthz
curl http://localhost:8081/healthz

# DID resolution (reads from live Accumulate network)
curl "http://localhost:8080/resolve?did=did:acc:alice"

# DID registration (submits to live Accumulate network)
curl -X POST -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:mycompany",
    "didDocument": {
      "@context": ["https://www.w3.org/ns/did/v1"],
      "id": "did:acc:mycompany"
    }
  }' \
  http://localhost:8081/register
```

### Environment Variables
- **ACC_NODE_URL**: Accumulate node endpoint (required for REAL mode)
  - Example: `http://localhost:26657` (local devnet)
  - Example: `https://mainnet.accumulatenetwork.io` (mainnet)

## Next Steps

### Recently Completed ✅
- **End-to-End DID Flows**: Complete create/resolve/update/deactivate workflows working
- **Universal Registrar 1.0**: Patch-based updates with addService/removeService support
- **REAL Mode Integration**: JSON-RPC v3 client with proper Accumulate build API usage
- **Hello Accu Example**: Complete `examples/hello_accu/` with ADI creation, data account setup, and DID document writing
- **Repository Surgery**: Fixed compilation issues, consolidated duplicate types, added proper workspace setup
- **Enhanced Documentation**: Comprehensive README files with troubleshooting, examples, and API reference

### Smoke Tests Verified ✅
- **Deactivation Flow**: PowerShell test confirmed Universal Registrar 1.0 `/1.0/deactivate` endpoint works correctly
- **Service Updates**: PowerShell test confirmed patch-based service addition via `/1.0/update` endpoint
- **Hello Accu Compilation**: Complete Accumulate DID lifecycle example compiles and includes smoke.ps1 script
- **Services Running**: Both resolver (:8084) and registrar (:8082) services deployed and responding

### Current Status
- **Resolver Service**: Running on `:8084` with native and Universal endpoints
- **Registrar Service**: Running on `:8082` with Universal Registrar 1.0 compatibility
- **Hello Accu Example**: Ready for testing with local or remote Accumulate nodes
- **Documentation**: Up-to-date with port configurations, endpoint details, and troubleshooting guides

### Immediate Next Steps
- **Integration Testing**: Run hello_accu example against live Accumulate devnet
- **Driver Integration**: Connect services to Universal Resolver/Registrar infrastructure
- **Performance Testing**: Benchmark resolution and registration throughput
- **Security Review**: Audit key management and authorization flows