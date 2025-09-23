# Accumulate DID Method Implementation

[![License](https://img.shields.io/badge/license-Proprietary-red.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org/dl/)
[![Docker](https://img.shields.io/badge/docker-supported-green.svg)](https://docs.docker.com/)

Complete implementation of the `did:acc` method for Accumulate Protocol, providing decentralized identifiers with W3C DID Core compliance and Universal Resolver/Registrar compatibility.

## 🚀 Project Synopsis

**Accumulate DID Services:**
- **Resolver** (`:8080`) - DID document resolution service with proper `did:acc` → `acc://` URL mapping
- **Registrar** (`:8081`) - DID lifecycle management with native and Universal endpoints
- **Universal Drivers** - Standards-compliant proxy drivers for ecosystem integration
- **Complete Documentation** - OpenAPI specs, method specification, and operator guides

**Key Features:**
- **Native Accumulate Integration** - Direct JSON-RPC v3 client integration with Accumulate nodes
- **W3C DID Core Compliance** - Full standards compliance with resolution metadata
- **Universal Compatibility** - Both Universal Resolver/Registrar v1.0 and legacy v0.x support
- **Production Ready** - Static analysis, race testing, performance benchmarks, security scanning
- **Local-First Development** - Complete local release workflow with no remote CI dependency

## ⚡ Quick Start

### FAKE vs REAL Modes

**FAKE Mode (Default)** - Offline development with static fixtures:
```bash
# No network dependency - uses golden files and mock submitter
./resolver --addr :8080
./registrar --addr :8081

# Test immediately
curl "http://localhost:8080/resolve?did=did:acc:alice"
```

**REAL Mode** - Live Accumulate network integration:
```bash
# Requires running Accumulate node
export ACC_NODE_URL=http://localhost:26657

./resolver --addr :8080 --real
./registrar --addr :8081 --real

# Requires funded lite account for registrar operations
curl -X POST "http://localhost:8081/register" -d '{...}'
```

### 🐳 Linux-First Development (RECOMMENDED)

**Container-First Workflow:**
```bash
# Interactive development shell
make dev-shell

# Complete CI pipeline
make ci-local

# Start services for testing
make dev-up
make dev-down
```

**Development Environment:**
- Uses `docker-compose.dev.yml` for standardized Linux environment
- Go 1.25 + Node.js LTS + development tools pre-installed
- Repository mounted at `/workspace` with proper permissions
- Host Docker integration for image builds

## 💰 Fees at a Glance

**⚠️ REAL Mode Credit Requirements:**

DID operations on Accumulate require credits for blockchain transactions. USD values shown are estimates for illustration; credits are the unit of account.

| DID Operation | Credits | USD (Approx) |
|---------------|---------|--------------|
| Create ADI (8+ chars) | 500.00 | $5.0000 |
| Create Data Account | 25.00 | $0.2500 |
| Write DID Document | 0.10 | $0.0010 |
| **Complete DID Setup** | **~525.10** | **~$5.25** |

**Short ADI Names Cost More:** 1-char names cost 4.8M credits, 7-char names cost 1,800 credits.

👉 **[Complete fee schedule and details](docs/ops/OPERATIONS.md#fees--credits)**

**Credit Management:**

**Development/Testing (Devnet):**
```bash
# Get devnet faucet credits
curl -X POST https://devnet-faucet.accumulate.io/get \
  -d '{"url":"<lite-account-url>"}'

# Check balance
accumulate account get <lite-account-url>
```

**Production (MainNet):**
```bash
# 1. Acquire ACME tokens via exchange
# 2. Convert ACME to credits using oracle rate
accumulate oracle price
accumulate tx create <acme-account> --to <lite-account> --amount <acme-amount>

# 3. Fund ADI key pages for authorization
accumulate credits <adi-key-page-url> <credit-amount>
```

**Recommended Funding:**
- **Development:** 100-500 credits per lite account
- **Production:** Monitor balances and set up alerts for low credit conditions
- **Budget:** 20-50 credits per DID for full lifecycle testing

## 🌐 Ports & Endpoints

### Core Services

| Service | Port | Purpose | Health Check |
|---------|------|---------|--------------|
| **Resolver** | 8080 | DID document resolution | `/healthz` |
| **Registrar** | 8081 | DID lifecycle management | `/healthz` |

### Universal Drivers (Optional)

| Service | Port | Purpose | Health Check |
|---------|------|---------|--------------|
| **UniResolver** | 8090 | Universal Resolver proxy | `/health` |
| **UniRegistrar** | 8091 | Universal Registrar proxy | `/health` |

### Key Endpoints

**Resolver:**
- `GET /resolve?did={did}` - W3C DID Resolution
- `GET /healthz` - Service health

**Registrar (Native):**
- `POST /register` - Create new DID (ADI + data account + document)
- `POST /native/update` - Update existing DID document
- `POST /native/deactivate` - Deactivate DID

**Registrar (Universal v1.0):**
- `POST /1.0/create` - Universal Registrar create
- `POST /1.0/update` - Universal Registrar update (patch-based)
- `POST /1.0/deactivate` - Universal Registrar deactivate

**Environment Variables:**
- `ACC_NODE_URL` - Accumulate JSON-RPC endpoint (required for REAL mode)
  - Local devnet: `http://localhost:26657`
  - MainNet: `https://mainnet.accumulatenetwork.io`

## 🏗️ Local Release Workflow

**Complete Local Release Pipeline** (no remote CI dependency):

### Single Command Release
```bash
# Complete release pipeline
make release-local
```

This runs:
1. **QA Checks** - static analysis, race tests, imports verification
2. **Cross-Platform Binaries** - linux/amd64, linux/arm64, windows/amd64, darwin/arm64
3. **Multi-Arch Docker Images** - linux/amd64,linux/arm64 with embedded versions
4. **Documentation Packaging** - complete docs archive with checksums
5. **Security Scanning** - SBOM generation and vulnerability scans
6. **Git Tagging** - annotated tags with version consistency checks

### Individual Components
```bash
# Quality assurance
make qa                 # All QA checks (lint, vet, race, imports)

# Build artifacts
make binaries-local     # Cross-platform binaries
make images-local       # Multi-arch Docker images
make docs-archive       # Documentation package

# Security compliance
make sbom-local         # Software Bill of Materials
make scan-local         # Vulnerability scanning

# Maintenance
make dist-clean         # Clean distribution directory
```

### Prerequisites
**Required:**
- Docker with buildx (multi-arch support)
- Git (version tagging)

**Optional (graceful degradation):**
- `syft` for SBOM generation
- `trivy` for vulnerability scanning
- `zip`/`7z` for documentation packaging (fallback to tar)

### Distribution Artifacts

After `make release-local`:
```
dist/bin/           # Cross-platform binaries
dist/images/        # Docker images (local registry)
dist/docs/          # Documentation archive with checksums
dist/sbom/          # Software Bill of Materials
dist/scan/          # Vulnerability reports
```

## 📚 Documentation

### Quick Links
- **[Operations Guide](docs/ops/OPERATIONS.md)** - Complete deployment and management guide
- **[Method Specification](docs/spec/method.md)** - DID method detailed specification
- **[API Documentation](docs/site/)** - Generated OpenAPI documentation

### OpenAPI Specifications
- **[Resolver API](docs/spec/openapi/resolver.yaml)** - DID resolution service
- **[Registrar API](docs/spec/openapi/registrar.yaml)** - DID registration service

### Development Resources
- **[Hello Accumulate Example](examples/hello_accu/)** - Complete DID lifecycle example
- **[Performance Tests](perf/)** - k6 load testing scenarios
- **[Postman Collection](postman/)** - API testing collection

### Standards References
- [W3C DID Core 1.0](https://www.w3.org/TR/did-core/)
- [DID Resolution](https://w3c-ccg.github.io/did-resolution/)
- [Universal Resolver](https://github.com/decentralized-identity/universal-resolver)
- [Universal Registrar](https://github.com/decentralized-identity/universal-registrar)

## 🔧 Development Workflow

### Container-First Development
```bash
# Launch development environment
make dev-shell

# Run tests
make test-all

# Build documentation
make docs-container

# Complete CI validation
make ci-local
```

### Local Tools (Alternative)
```bash
# Prerequisites: Go 1.25+, Node.js LTS, Docker
make test           # Run tests
make docs           # Build documentation
make docker-build   # Build images
make all            # Complete build
```

### Available Make Targets
```bash
make help           # Show all available targets
```

**Development:**
- `dev-shell` - Interactive development container
- `test-all` - Run all tests in container
- `ci-local` - Complete CI pipeline

**Quality Assurance:**
- `qa` - All QA checks (imports, vet, lint, race)
- `lint` - golangci-lint static analysis
- `test-race` - Race condition testing

**Release Pipeline:**
- `release-local` - Complete local release
- `binaries-local` - Cross-platform binaries
- `images-local` - Multi-arch Docker images

## 🏛️ Architecture

### DID Method Format
```
did:acc:<adi-label>[/<path>]
```

**Examples:**
- `did:acc:alice` → ADI: `acc://alice`, Data Account: `acc://alice/did`
- `did:acc:beastmode.acme` → ADI: `acc://beastmode.acme`, Data Account: `acc://beastmode.acme/did`
- `did:acc:alice/documents` → ADI: `acc://alice`, Data Account: `acc://alice/documents`

### Implementation Stack
- **Language:** Go 1.25 with modules and workspace
- **Router:** chi router with middleware
- **Client:** Accumulate JSON-RPC v3 client
- **Standards:** W3C DID Core, Universal Resolver/Registrar v1.0
- **Testing:** Go testing, race detection, k6 performance testing
- **Quality:** golangci-lint, static analysis, import verification

### Service Architecture
```
┌─────────────────┐    ┌─────────────────┐
│   Universal     │    │   Universal     │
│   Resolver      │    │   Registrar     │
│   Driver        │    │   Driver        │
│   (Port 8090)   │    │   (Port 8091)   │
└─────────────────┘    └─────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌─────────────────┐
│    Resolver     │    │    Registrar    │
│   (Port 8080)   │    │   (Port 8081)   │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────────┬───────────┘
                     ▼
              ┌─────────────────┐
              │   Accumulate    │
              │     Node        │
              │ (JSON-RPC v3)   │
              └─────────────────┘
```

## 📜 License

**License:** See [LICENSE](LICENSE) file for details.

**Note:** This implementation follows Accumulate Protocol licensing terms and W3C DID specification compatibility requirements.

---

## 🚀 Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/opendlt/accu-did.git
   cd accu-did
   ```

2. **Start development environment:**
   ```bash
   make dev-shell
   ```

3. **Run smoke test:**
   ```bash
   make dev-up
   curl http://localhost:8080/healthz
   curl http://localhost:8081/healthz
   ```

4. **Create your first DID:**
   ```bash
   curl -X POST http://localhost:8081/register \
     -H "Content-Type: application/json" \
     -d '{
       "did": "did:acc:mytest",
       "didDocument": {
         "@context": ["https://www.w3.org/ns/did/v1"],
         "id": "did:acc:mytest"
       }
     }'
   ```

5. **Resolve the DID:**
   ```bash
   curl "http://localhost:8080/resolve?did=did:acc:mytest"
   ```

**Need help?** Check the [Operations Guide](docs/ops/OPERATIONS.md) for detailed deployment instructions and troubleshooting.