# accu-DID Operations Guide

Complete operator's guide for running and managing the Accumulate DID resolver and registrar services.

## 1. Overview

### What These Services Do

- **Resolver** (`:8080`) - Resolves `did:acc:*` identifiers to DID documents following W3C DID Core spec
- **Registrar** (`:8081`) - Creates, updates, and deactivates DID documents on Accumulate Protocol

### Operating Modes

- **FAKE Mode** - Uses static fixtures, no blockchain dependency (testing/development)
- **REAL Mode** - Connects to live Accumulate network, requires `ACC_NODE_URL`

### Port Table

| Service | Port | Purpose |
|---------|------|---------|
| Resolver | 8080 | DID resolution service |
| Registrar | 8081 | DID lifecycle management |
| UniResolver Driver | 8090 | Universal Resolver proxy (optional) |
| UniRegistrar Driver | 8091 | Universal Registrar proxy (optional) |

### Environment Variables

- **`ACC_NODE_URL`** - Accumulate JSON-RPC endpoint (required for REAL mode)
  - Local devnet: `http://127.0.0.1:26660`
  - MainNet: `https://mainnet.accumulate.defidevs.io/v2`

### Credit Requirements (REAL Mode)

**⚠️ IMPORTANT:** DID operations on Accumulate require credits for blockchain transactions.

**Credit Costs:**
- **ADI Creation:** ~10 credits per ADI (one-time for each unique `did:acc:<adi>` namespace)
- **Data Account Creation:** ~5 credits per data account
- **DID Document Write:** ~2-5 credits per write operation (create/update/deactivate)
- **Total for new DID:** ~17-20 credits for complete `did:acc:<adi>/<path>` setup

**Credit Management:**
- Credits must be pre-funded to the signer's lite account before operations
- Monitor credit balance: insufficient credits will cause transaction failures
- Each service operation may require multiple blockchain transactions
- Budget approximately 20-50 credits per DID for full lifecycle testing

**Funding Commands:**
```bash
# Check balance (requires Accumulate CLI)
accumulate account get <lite-account-url>

# Add credits to lite account
accumulate credits <lite-account-url> <amount>
```

**Development Recommendations:**
- Use FAKE mode for development to avoid credit consumption
- For REAL mode testing, start with a funded devnet lite account
- Consider credit costs when designing production deployment strategies

### Accumulate Credits Quickstart

**Development/Testing (Devnet):**
```bash
# Get devnet faucet credits (if available)
curl -X POST https://devnet-faucet.accumulate.io/get -d '{"url":"<lite-account-url>"}'

# Check balance
accumulate account get <lite-account-url>

# Manual credit transfer from funded account
accumulate tx create <source-lite-account> --to <destination-lite-account> --amount 1000.00
```

**Mainnet Production:**
1. **Acquire ACME tokens** via exchange or transfer
2. **Convert ACME to credits** using oracle rate:
   ```bash
   accumulate oracle price  # Check current ACME:credit rate
   accumulate tx create <acme-account> --to <lite-account> --amount <acme-amount>
   ```
3. **Fund ADI Key Page** for transaction authorization:
   ```bash
   accumulate credits <adi-key-page-url> <credit-amount>
   ```

**Registrar Preflight Check:**
- **Operator must ensure sufficient credits** before running registrar operations
- **Registrar returns error** if submitted transactions fail due to insufficient credits
- **Monitor credit balances** and set up alerts for low balance conditions
- **Budget 20-50 credits** per DID for full lifecycle testing scenarios

## 2. Container-First Workflow (RECOMMENDED)

### Prerequisites

- Docker and Docker Compose
- Optional: Accumulate devnet running for REAL mode testing

### Quick Start with Containers

**One-liner development shell:**
```bash
make dev-shell
```

**Complete development workflow:**
```bash
# 1. Interactive development environment
make dev-shell

# 2. Run all tests
make test-all

# 3. Build documentation
make docs-container

# 4. Run complete local CI
make ci-local

# 5. Start services for testing
make dev-up

# 6. Stop services
make dev-down
```

### Container Environment Details

The development container provides:
- **Repository mounted at:** `/workspace`
- **Accumulate repo (read-only):** `/workspace/.accumulate-ro` (if `../accumulate` exists)
- **Go 1.23** with modules and workspace support
- **Node.js LTS** with npm for documentation builds
- **Non-root user:** `dev` (uid 1000) for file permission compatibility

**Environment Variables in Container:**
- `ACC_NODE_URL=http://host.docker.internal:26660` (for REAL mode)
- `GO111MODULE=on`
- `CGO_ENABLED=0`

### Available Make Targets

| Target | Description |
|--------|-------------|
| `make help` | Show all available targets |
| `make dev-shell` | Launch interactive development container |
| `make test-all` | Run all Go tests in container |
| `make docs-container` | Build API documentation in container |
| `make ci-local` | Run complete CI pipeline in container |
| `make dev-up` | Start resolver and registrar services |
| `make dev-down` | Stop all development services |
| `make check-imports` | Verify no forbidden imports |
| `make conformance` | Run conformance tests |
| `make perf` | Run performance tests with k6 |

## 3. Legacy Workflow (Local Tools)

### Prerequisites

- Go 1.23+
- Node.js LTS (for documentation)
- Optional: Accumulate devnet running

### Option A: Go Services (Development)

**Windows PowerShell:**
```powershell
# Set environment
$env:ACC_NODE_URL = "http://127.0.0.1:26660"

# Terminal 1 - Resolver
cd resolver-go
go run cmd/server/main.go --addr :8080 --mode REAL

# Terminal 2 - Registrar
cd registrar-go
go run cmd/server/main.go --addr :8081 --mode REAL
```

**Unix/Mac Bash:**
```bash
# Set environment
export ACC_NODE_URL=http://127.0.0.1:26660

# Terminal 1 - Resolver
cd resolver-go && go run cmd/server/main.go --addr :8080 --mode REAL

# Terminal 2 - Registrar
cd registrar-go && go run cmd/server/main.go --addr :8081 --mode REAL
```

### Option B: Docker Compose (Production-like)

**Core services only:**
```powershell
# Windows
docker-compose up -d

# Check status
docker-compose ps
```

```bash
# Unix/Mac
docker-compose up -d

# Check status
docker-compose ps
```

**With Universal drivers:**
```powershell
# Windows
docker-compose --profile uni up -d

# Check all services
docker-compose --profile uni ps
```

```bash
# Unix/Mac
docker-compose --profile uni up -d

# Check all services
docker-compose --profile uni ps
```

### Health Checks

```powershell
# Windows PowerShell
Invoke-RestMethod http://127.0.0.1:8080/health
Invoke-RestMethod http://127.0.0.1:8081/health

# With Universal drivers
Invoke-RestMethod http://127.0.0.1:8090/health
Invoke-RestMethod http://127.0.0.1:8091/health
```

```bash
# Unix/Mac
curl http://127.0.0.1:8080/health
curl http://127.0.0.1:8081/health

# With Universal drivers
curl http://127.0.0.1:8090/health
curl http://127.0.0.1:8091/health
```

### First Smoke Test

```powershell
# Windows PowerShell - Create DID
$body = @{
    did = "did:acc:smoketest.acme"
    document = @{
        "@context" = @("https://www.w3.org/ns/did/v1")
        id = "did:acc:smoketest.acme"
        verificationMethod = @(@{
            id = "did:acc:smoketest.acme#key1"
            type = "Ed25519VerificationKey2020"
            controller = "did:acc:smoketest.acme"
            publicKeyMultibase = "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK"
        })
    }
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri "http://127.0.0.1:8081/register" -Method POST -Body $body -ContentType "application/json"

# Resolve
Invoke-RestMethod "http://127.0.0.1:8080/resolve?did=did:acc:smoketest.acme"
```

```bash
# Unix/Mac - Create DID
curl -X POST http://127.0.0.1:8081/register \
  -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:smoketest.acme",
    "document": {
      "@context": ["https://www.w3.org/ns/did/v1"],
      "id": "did:acc:smoketest.acme",
      "verificationMethod": [{
        "id": "did:acc:smoketest.acme#key1",
        "type": "Ed25519VerificationKey2020",
        "controller": "did:acc:smoketest.acme",
        "publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK"
      }]
    }
  }'

# Resolve
curl "http://127.0.0.1:8080/resolve?did=did:acc:smoketest.acme"
```

## 3. Service Management

### Docker Compose Operations

**Start services:**
```powershell
# Core services only
docker-compose up -d

# With Universal drivers
docker-compose --profile uni up -d

# Rebuild and start
docker-compose up -d --build
```

**Stop services:**
```powershell
# Stop all
docker-compose down

# Stop with cleanup
docker-compose down --volumes --remove-orphans
```

**Restart services:**
```powershell
# Restart specific service
docker-compose restart resolver
docker-compose restart registrar

# Restart all
docker-compose restart
```

### Logs

```powershell
# Windows PowerShell
docker-compose logs -f resolver
docker-compose logs -f registrar
docker-compose logs -f uniresolver uniregistrar

# All services
docker-compose logs -f
```

```bash
# Unix/Mac
docker-compose logs -f resolver
docker-compose logs -f registrar
docker-compose logs -f uniresolver uniregistrar

# All services
docker-compose logs -f
```

### Windows Firewall

If accessing from other machines, allow ports:
```powershell
# Windows PowerShell (as Administrator)
New-NetFirewallRule -DisplayName "accu-DID Resolver" -Direction Inbound -Protocol TCP -LocalPort 8080 -Action Allow
New-NetFirewallRule -DisplayName "accu-DID Registrar" -Direction Inbound -Protocol TCP -LocalPort 8081 -Action Allow
New-NetFirewallRule -DisplayName "accu-DID UniResolver" -Direction Inbound -Protocol TCP -LocalPort 8090 -Action Allow
New-NetFirewallRule -DisplayName "accu-DID UniRegistrar" -Direction Inbound -Protocol TCP -LocalPort 8091 -Action Allow
```

## 4. End-to-End Smoke Tests

### FAKE Mode (No Network)

```powershell
# Windows PowerShell
cd resolver-go
go test ./... -v

cd ..\registrar-go
go test ./... -v

# Start in FAKE mode (separate terminals)
go run cmd/server/main.go --addr :8080 --mode FAKE
go run cmd/server/main.go --addr :8081 --mode FAKE

# Health checks should work
Invoke-RestMethod http://127.0.0.1:8080/health
Invoke-RestMethod http://127.0.0.1:8081/health
```

```bash
# Unix/Mac
cd resolver-go && go test ./... -v
cd ../registrar-go && go test ./... -v

# Start in FAKE mode
cd resolver-go && go run cmd/server/main.go --addr :8080 --mode FAKE &
cd registrar-go && go run cmd/server/main.go --addr :8081 --mode FAKE &

# Health checks
curl http://127.0.0.1:8080/health
curl http://127.0.0.1:8081/health
```

### REAL Mode (Full Stack)

**Complete lifecycle test:**
```bash
# 1. Create DID
curl -X POST http://127.0.0.1:8081/register \
  -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:e2etest.acme",
    "document": {
      "@context": ["https://www.w3.org/ns/did/v1"],
      "id": "did:acc:e2etest.acme",
      "verificationMethod": [{
        "id": "did:acc:e2etest.acme#key1",
        "type": "Ed25519VerificationKey2020",
        "controller": "did:acc:e2etest.acme",
        "publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK"
      }]
    }
  }'

# 2. Resolve (should return status 200)
curl "http://127.0.0.1:8080/resolve?did=did:acc:e2etest.acme"

# 3. Update (add service)
curl -X POST http://127.0.0.1:8081/update \
  -H "Content-Type: application/json" \
  -d '{
    "did": "did:acc:e2etest.acme",
    "patch": {
      "addService": {
        "id": "did:acc:e2etest.acme#website",
        "type": "LinkedDomains",
        "serviceEndpoint": "https://example.com"
      }
    }
  }'

# 4. Resolve (should show new service)
curl "http://127.0.0.1:8080/resolve?did=did:acc:e2etest.acme"

# 5. Deactivate
curl -X POST http://127.0.0.1:8081/deactivate \
  -H "Content-Type: application/json" \
  -d '{"did": "did:acc:e2etest.acme"}'

# 6. Resolve (should return 410 or deactivated:true)
curl "http://127.0.0.1:8080/resolve?did=did:acc:e2etest.acme"
```

### Postman Collection

```powershell
# Windows PowerShell (requires Newman)
cd postman
newman run accu-did.postman_collection.json -e local.postman_environment.json --verbose
```

```bash
# Unix/Mac (requires Newman)
cd postman && newman run accu-did.postman_collection.json -e local.postman_environment.json --verbose
```

## 5. Windows Services (Optional)

Using NSSM (Non-Sucking Service Manager) to run as Windows services.

### Install NSSM

Download from https://nssm.cc/ and add to PATH.

### Install Services

```powershell
# Windows PowerShell (as Administrator)

# Resolver service
nssm install "accu-did-resolver" "go.exe"
nssm set "accu-did-resolver" Application "C:\path\to\go.exe"
nssm set "accu-did-resolver" AppDirectory "C:\Accumulate_Stuff\accu-did\resolver-go"
nssm set "accu-did-resolver" AppParameters "run cmd/server/main.go --addr :8080 --mode REAL"
nssm set "accu-did-resolver" AppEnvironmentExtra "ACC_NODE_URL=http://127.0.0.1:26660"
nssm set "accu-did-resolver" DisplayName "Accumulate DID Resolver"
nssm set "accu-did-resolver" Description "Accumulate DID Resolution Service"

# Registrar service
nssm install "accu-did-registrar" "go.exe"
nssm set "accu-did-registrar" Application "C:\path\to\go.exe"
nssm set "accu-did-registrar" AppDirectory "C:\Accumulate_Stuff\accu-did\registrar-go"
nssm set "accu-did-registrar" AppParameters "run cmd/server/main.go --addr :8081 --mode REAL"
nssm set "accu-did-registrar" AppEnvironmentExtra "ACC_NODE_URL=http://127.0.0.1:26660"
nssm set "accu-did-registrar" DisplayName "Accumulate DID Registrar"
nssm set "accu-did-registrar" Description "Accumulate DID Registration Service"
```

### Manage Services

```powershell
# Start services
Start-Service "accu-did-resolver"
Start-Service "accu-did-registrar"

# Stop services
Stop-Service "accu-did-resolver"
Stop-Service "accu-did-registrar"

# Check status
Get-Service "accu-did-*"

# Uninstall
nssm remove "accu-did-resolver" confirm
nssm remove "accu-did-registrar" confirm
```

## 6. Local Automation (No Remote CI)

### Local CI Scripts

```powershell
# Windows PowerShell - Full CI pipeline
.\scripts\local-ci.ps1

# Individual components
.\scripts\build-docs.ps1
go test .\resolver-go\... -v
go test .\registrar-go\... -v
```

```bash
# Unix/Mac - Full CI pipeline
bash scripts/local-ci.sh

# Individual components
bash scripts/build-docs.sh
cd resolver-go && go test ./... -v
cd registrar-go && go test ./... -v
```

### Makefile Targets

```bash
# Available targets
make help

# Common operations
make test          # Run all Go tests
make docs          # Build Redoc documentation
make docker-build  # Build Docker images
make all           # test + docs + docker-build
make ci            # Run local CI checks
make clean-all     # Clean all artifacts
```

## 7. Release Workflow (Complete Local Pipeline)

### Overview

The accu-did project provides a complete local release workflow that eliminates dependency on remote CI/CD systems. The workflow includes:

- **Quality Assurance**: Static analysis, race testing, imports checking
- **Cross-platform Builds**: Binaries for linux/amd64, linux/arm64, windows/amd64, darwin/arm64
- **Container Images**: Multi-arch Docker images with embedded version info
- **Security Scanning**: Vulnerability scans with trivy and SBOM generation with syft
- **Documentation Packaging**: Complete docs archive for distribution
- **Version Management**: Git tagging with consistency checks across OpenAPI specs

### Prerequisites

**Required Tools for Complete Workflow:**
- Docker with buildx (multi-arch support)
- Git (version tagging)
- zip/7z/tar (documentation packaging)

**Optional Security Tools:**
- `syft` for SBOM generation: https://github.com/anchore/syft/releases
- `trivy` for vulnerability scanning: https://github.com/aquasecurity/trivy/releases

**Tool Installation (Optional):**
```bash
# Install syft (SBOM generation)
curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin

# Install trivy (vulnerability scanning)
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin
```

### Complete Release Pipeline

**Single Command Release:**
```bash
# Complete local release (all steps)
make release-local
```

This runs the full pipeline:
1. **QA Checks** (`make qa`) - imports, vet, lint, race tests
2. **Binary Compilation** (`make binaries-local`) - cross-platform binaries
3. **Image Building** (`make images-local`) - multi-arch Docker images
4. **Documentation** (`make docs-archive`) - packaged docs archive
5. **Security Scanning** (`make sbom-local scan-local`) - SBOM + vulnerability scans
6. **Version Tagging** (`scripts/release.local.sh`) - git tag with metadata

### Individual Components

**Quality Assurance:**
```bash
# Run all QA checks
make qa

# Individual checks
make lint           # golangci-lint static analysis
make test-race      # race condition testing
make vet           # go vet analysis
./scripts/check-imports.sh  # forbidden imports check
```

**Cross-Platform Binaries:**
```bash
# Build all platform binaries
make binaries-local

# Outputs to:
# dist/bin/linux-amd64/resolver, dist/bin/linux-amd64/registrar
# dist/bin/linux-arm64/resolver, dist/bin/linux-arm64/registrar
# dist/bin/windows-amd64/resolver.exe, dist/bin/windows-amd64/registrar.exe
# dist/bin/darwin-arm64/resolver, dist/bin/darwin-arm64/registrar
```

**Multi-Arch Docker Images:**
```bash
# Build Docker images for linux/amd64,linux/arm64
make images-local

# Creates images:
# accu-did/resolver:v<version>, accu-did/resolver:latest
# accu-did/registrar:v<version>, accu-did/registrar:latest
```

**Documentation Archive:**
```bash
# Package documentation for distribution
make docs-archive

# Creates: dist/docs/docs-<version>.zip
# Contains: API docs, OpenAPI specs, method specs, README
```

**Security Scanning:**
```bash
# Generate Software Bill of Materials
make sbom-local

# Run vulnerability scans
make scan-local

# Outputs:
# dist/sbom/ - SBOM files in JSON, SPDX, and table formats
# dist/scan/ - Vulnerability reports and critical findings
```

**Version Management:**
```bash
# Create git tag with version consistency checks
./scripts/release.local.sh   # Unix/Linux
./scripts/release.local.ps1  # Windows PowerShell

# Checks:
# - Working tree is clean
# - VERSION file exists
# - OpenAPI spec versions match VERSION
# - Tag doesn't already exist
# Creates annotated tag with release metadata
```

### Version Lifecycle

**Update Version:**
```bash
# Edit VERSION file manually
echo "0.2.0" > VERSION

# Update OpenAPI spec versions to match
# docs/spec/openapi/resolver.yaml
# docs/spec/openapi/registrar.yaml
```

**Create Release:**
```bash
# Commit version change
git add VERSION docs/spec/openapi/*.yaml
git commit -m "bump: release version 0.2.0"

# Run release pipeline
make release-local

# Optionally push to remote (manual step)
git push origin main
git push origin v0.2.0
```

## 8. Backlog Triage via TODO Scanner

The repository includes a comprehensive TODO scanner for tracking technical debt and work items across the codebase.

### Quick Scan

**Recommended command:**
```bash
make todo-scan
```

**Alternative execution methods:**
```bash
# Docker (cross-platform)
docker compose -f docker-compose.dev.yml run --rm dev 'go run tools/todoscan/main.go .'

# Linux/Unix shell script
./scripts/todo-scan.sh

# Windows PowerShell script
.\scripts\todo-scan.ps1

# Direct Go execution (requires local Go installation)
go run tools/todoscan/main.go .
```

### TODO Markers Detected

The scanner searches for these patterns (case-insensitive):

| Tag | Purpose | Priority |
|-----|---------|----------|
| `TODO` | General work items | Medium |
| `FIXME` | Known bugs/issues | High |
| `XXX` | Code requiring attention | High |
| `HACK` | Temporary workarounds | Medium |
| `STUB` | Placeholder implementations | Medium |
| `TBA/TBD` | Items to be added/determined | Low |
| `NOTIMPL` | Missing implementations | High |
| `PANIC("TODO")` | Critical unimplemented paths | Critical |

### Report Outputs

Reports are generated in `./reports/`:

- **`todo-report.json`** - Machine-readable data with full context
- **`todo-report.md`** - Human-readable report with code excerpts
- **`todo-report.csv`** - Spreadsheet-compatible format for analysis

### Operational Workflows

**Weekly Maintenance:**
```bash
# 1. Run scan during sprint planning
make todo-scan

# 2. Review critical items
grep -E "(FIXME|XXX|PANIC)" reports/todo-report.md

# 3. Check implementation gaps
grep -E "(NOTIMPL|STUB)" reports/todo-report.md

# 4. Monitor technical debt
grep -E "(HACK)" reports/todo-report.md
```

**Component-Specific Analysis:**
```bash
# Resolver issues
jq '.items[] | select(.path | startswith("resolver-go/"))' reports/todo-report.json

# Registrar issues
jq '.items[] | select(.path | startswith("registrar-go/"))' reports/todo-report.json

# Documentation tasks
jq '.items[] | select(.path | startswith("docs/"))' reports/todo-report.json

# Universal driver issues
jq '.items[] | select(.path | startswith("drivers/"))' reports/todo-report.json
```

**Trend Analysis:**
```bash
# Compare counts over time
echo "$(date): $(jq '.totalCount' reports/todo-report.json)" >> reports/todo-trend.log

# Tag distribution
jq '.summary.countsByTag' reports/todo-report.json

# Directory breakdown
jq '.summary.countsByDir' reports/todo-report.json
```

### Integration with Issue Tracking

**Escalation Criteria:**
- **FIXME/XXX items**: Convert to GitHub issues if affecting operations
- **NOTIMPL items**: Add to `spec/BACKLOG.md` if blocking features
- **PANIC items**: Address immediately - these indicate critical gaps
- **HACK items**: Schedule proper implementation in upcoming sprints

**Automation Integration:**
```yaml
# Example CI check
- name: TODO Scanner
  run: make todo-scan

- name: Alert on Critical TODOs
  run: |
    if jq -e '.items[] | select(.tag == "PANIC" or .tag == "FIXME")' reports/todo-report.json > /dev/null; then
      echo "::warning::Critical TODOs found - review required"
      jq '.items[] | select(.tag == "PANIC" or .tag == "FIXME")' reports/todo-report.json
    fi
```

### Best Practices

**TODO Creation Guidelines:**
- Include brief context: `// TODO: add rate limiting for /resolve endpoint`
- Reference issues when applicable: `// TODO(#123): implement batch resolution`
- Use appropriate tags: `FIXME` for bugs, `TODO` for features
- Avoid generic comments: prefer `TODO: validate DID format` over `TODO: fix this`

**Cleanup Workflow:**
1. **Monthly review**: Run scanner and categorize items by urgency
2. **Sprint planning**: Convert high-priority TODOs to formal tasks
3. **Refactoring sprints**: Dedicate time to address HACK items
4. **Release preparation**: Ensure no PANIC items in production code

**Reporting:**
- Include TODO count trends in sprint retrospectives
- Track resolution rate of TODO items over time
- Use TODO density (items per KLOC) as code quality metric

### Distribution Artifacts

After `make release-local` completes, these artifacts are available:

**Binaries (`dist/bin/`):**
```
linux-amd64/resolver, linux-amd64/registrar
linux-arm64/resolver, linux-arm64/registrar
windows-amd64/resolver.exe, windows-amd64/registrar.exe
darwin-arm64/resolver, darwin-arm64/registrar
```

**Docker Images:**
```
accu-did/resolver:v<version>, accu-did/resolver:latest
accu-did/registrar:v<version>, accu-did/registrar:latest
```

**Documentation (`dist/docs/`):**
```
docs-<version>.zip - Complete documentation archive
docs-<version>.zip.sha256 - Checksum file
```

**Security Reports (`dist/sbom/`, `dist/scan/`):**
```
source-<version>.{syft.json,spdx.json,txt} - Source code SBOM
resolver-latest.{syft.json,spdx.json,txt} - Resolver image SBOM
registrar-latest.{syft.json,spdx.json,txt} - Registrar image SBOM
resolver-latest.{trivy.json,txt} - Resolver vulnerability scan
registrar-latest.{trivy.json,txt} - Registrar vulnerability scan
filesystem.{trivy.json,txt} - Source code vulnerability scan
*-critical.txt - Critical vulnerability reports (if any)
```

**Git Tag:**
```
v<version> - Annotated tag with release metadata
```

### Clean-up and Maintenance

**Clean Distribution Directory:**
```bash
# Remove all distribution artifacts
make dist-clean

# Rebuilds dist/ structure:
# dist/bin/, dist/images/, dist/docs/, dist/sbom/, dist/scan/
```

**Debugging Failed Releases:**

**QA Failures:**
```bash
# Check specific failure
make lint 2>&1 | head -20
make test-race 2>&1 | head -20

# Fix issues and re-run
make qa
```

**Docker Build Failures:**
```bash
# Check Docker buildx availability
docker buildx version

# Enable buildx if needed
docker buildx create --name mybuilder --use
```

**Missing Security Tools:**
```bash
# Skip SBOM/scanning if tools unavailable
# Scripts automatically create placeholder files when tools missing

# Check tool availability
command -v syft && echo "syft available" || echo "syft missing"
command -v trivy && echo "trivy available" || echo "trivy missing"
```

### Continuous Integration Integration

**Local CI Validation:**
```bash
# Full CI pipeline in container (RECOMMENDED)
make ci-local

# Equivalent to running complete pipeline in standardized environment
```

**Pre-Release Checklist:**
1. ✅ All tests pass (`make test-all`)
2. ✅ QA checks pass (`make qa`)
3. ✅ Services start cleanly (`make dev-up`)
4. ✅ Health checks succeed (`curl localhost:8080/health`)
5. ✅ Smoke tests pass (basic create/resolve/update/deactivate)
6. ✅ VERSION file updated
7. ✅ OpenAPI spec versions updated
8. ✅ CHANGELOG.md entries added
9. ✅ Working tree clean for tagging

**Post-Release Actions (Optional):**
```bash
# Push to git remote
git push origin main
git push origin v<version>

# Push Docker images to registry
docker push accu-did/resolver:v<version>
docker push accu-did/resolver:latest
docker push accu-did/registrar:v<version>
docker push accu-did/registrar:latest

# Upload documentation and security reports to release assets
# (Manual upload to GitHub releases, internal artifact storage, etc.)
```

### Legacy Version Management (Deprecated)

The following scripts exist for compatibility but are **deprecated** in favor of the complete release workflow:

```powershell
# Windows PowerShell (DEPRECATED)
.\scripts\bump-version.ps1 patch   # Use VERSION file instead
.\scripts\release.ps1              # Use make release-local instead
```

```bash
# Unix/Mac (DEPRECATED)
bash scripts/bump-version.sh patch   # Use VERSION file instead
bash scripts/release.sh              # Use make release-local instead
```

**Migration:** Replace legacy scripts with:
1. Manual VERSION file editing
2. `make release-local` for complete pipeline
3. `git push` commands for remote distribution (optional)

## 8. Troubleshooting

### Common Errors

**Port already in use:**
```bash
# Find process using port
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Unix/Mac

# Kill process
taskkill /PID <pid> /F        # Windows
kill -9 <pid>                 # Unix/Mac
```

**ACC_NODE_URL not set:**
```
Error: REAL mode requires ACC_NODE_URL environment variable
```
Solution: Set `ACC_NODE_URL=http://127.0.0.1:26660` or desired endpoint.

**Health check returns 502/404:**
- Service not started or crashed
- Check logs: `docker-compose logs <service>`
- Verify port binding: `netstat -an | findstr :8080`

**Health check returns 410:**
- Normal for deactivated DIDs
- Use different DID identifier for testing

**P2P dial spam (devnet):**
```
ERR failed to dial ... connection refused
```
This is benign - devnet nodes trying to connect to each other.

### Development Issues

**Go workspace sync:**
```bash
go work sync
go mod tidy  # in each module directory
```

**CORS errors:**
Add CORS headers to service if accessing from browser:
```go
w.Header().Set("Access-Control-Allow-Origin", "*")
```

**Windows Antivirus:**
- Exclude accu-did directory from real-time scanning
- Allow `go.exe` and `docker.exe` through firewall

## 9. Appendices

### Port/Process Matrix

| Port | Service | Process | Docker Container |
|------|---------|---------|-----------------|
| 8080 | Resolver | `resolver-go` | `accu-did-resolver` |
| 8081 | Registrar | `registrar-go` | `accu-did-registrar` |
| 8090 | UniResolver | `uniresolver-go` | `accu-did-uniresolver` |
| 8091 | UniRegistrar | `uniregistrar-go` | `accu-did-uniregistrar` |
| 26660 | Accumulate | `accumulated` | External |

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ACC_NODE_URL` | REAL mode | none | Accumulate JSON-RPC endpoint |
| `PORT` | No | 8080/8081 | Service listen port |
| `RESOLVER_URL` | UniDrivers | `http://resolver:8080` | Internal resolver URL |
| `REGISTRAR_URL` | UniDrivers | `http://registrar:8081` | Internal registrar URL |

### File Map

| Path | Purpose |
|------|---------|
| `scripts/local-ci.{ps1,sh}` | Local CI automation |
| `scripts/build-docs.{ps1,sh}` | Documentation builder |
| `scripts/bump-version.{ps1,sh}` | Version management |
| `scripts/release.{ps1,sh}` | Release automation |
| `docs/spec/openapi/` | OpenAPI specifications |
| `docs/site/` | Generated documentation |
| `postman/` | API test collection |
| `VERSION` | Current version |
| `CHANGELOG.md` | Release notes |
| `Makefile` | Build automation |
| `docker-compose.yml` | Container orchestration (production) |
| `docker-compose.dev.yml` | Development container orchestration |
| `docker/dev/Dockerfile` | Development container image |

## Windows Development Notes

**Container-First Approach:** Windows users should primarily use the containerized development workflow (`make dev-shell`, etc.) for consistency and easier onboarding.

**PowerShell Scripts:** The `.ps1` scripts in `scripts/` are maintained as legacy fallbacks but are **secondary** to the Linux container approach. New team members should start with containers.

**WSL2 Recommended:** For Windows users, WSL2 with Docker Desktop provides the best container development experience.

**File Permissions:** The development container uses uid 1000 (`dev` user) which maps well to most Linux and WSL2 environments. Windows file permission issues are avoided by working inside the container.