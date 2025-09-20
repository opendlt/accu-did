# Accumulate DID Implementation

Decentralized identifiers on Accumulate Protocol with Go services and Universal drivers.

## Quick Start

**Prerequisites:** Go 1.22+, Docker

**Key Services:**
- Resolver: `:8080` - DID document resolution
- Registrar: `:8081` - DID lifecycle management

**Most Common Commands:**
```bash
# Build & test everything
make build test

# Run resolver (development)
cd resolver-go && go run cmd/server/main.go

# Run registrar (development)
cd registrar-go && go run cmd/server/main.go
```

## Tech Stack

- **Framework:** chi router, golangci-lint
- **Module Base:** `github.com/opendlt/accu-did`
- **Standards:** W3C DID Core, Universal Resolver/Registrar

## Development

Use `make help` for all targets. Key services run on localhost with health checks at `/health`.

**Universal Drivers:** Standards-compliant proxies in `drivers/` directory.

## Documentation

- **API Docs:** `docs/` (MkDocs)
- **Deep Context:** [.claude/memory/CLAUDE.md](.claude/memory/CLAUDE.md)
- **Interop Plans:** `docs/interop/`

## Quick Test

```bash
# Health check
curl http://localhost:8080/health

# Resolve fixture DID
curl "http://localhost:8080/resolve?did=did:acc:beastmode.acme"
```