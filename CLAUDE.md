# Accumulate DID Implementation

Decentralized identifiers on Accumulate Protocol with Go services and Universal drivers.

## Quick Start

**Prerequisites:** Go 1.22+, Docker

**Key Services:**
- Resolver: `:8080` - DID document resolution (FAKE/REAL modes)
- Registrar: `:8081` - DID lifecycle management (FAKE/REAL modes)

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

## Quick Test

### FAKE Mode Smoke Tests
```bash
# Start services in FAKE mode
./resolver --addr :8080 &
./registrar --addr :8081 &

# Health checks
curl http://localhost:8080/healthz  # Should return {"status":"ok","timestamp":"..."}
curl http://localhost:8081/healthz  # Should return {"status":"ok","timestamp":"..."}

# DID resolution (uses golden files)
curl "http://localhost:8080/resolve?did=did:acc:beastmode.acme"

# DID registration (uses mock submitter)
curl -X POST -H "Content-Type: application/json" \
  -d '{"did":"did:acc:test","didDocument":{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:test"}}' \
  http://localhost:8081/create
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

# DID operations (connects to live Accumulate network)
curl "http://localhost:8080/resolve?did=did:acc:alice"
curl -X POST -H "Content-Type: application/json" \
  -d '{"did":"did:acc:alice","didDocument":{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:alice"}}' \
  http://localhost:8081/create
```

### Environment Variables
- **ACC_NODE_URL**: Accumulate node endpoint (required for REAL mode)
  - Example: `http://localhost:26657` (local devnet)
  - Example: `https://mainnet.accumulatenetwork.io` (mainnet)