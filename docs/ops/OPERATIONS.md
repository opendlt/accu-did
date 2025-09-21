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

## 2. Quick Start (Local Devnet)

### Prerequisites

Ensure Accumulate devnet is running:
```powershell
# Windows (if not already running)
accumulated run devnet --reset
```

```bash
# Unix/Mac (if not already running)
accumulated run devnet --reset
```

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

## 7. Releases (Local)

### Version Management

```powershell
# Windows PowerShell - Bump version
.\scripts\bump-version.ps1 patch   # 0.1.0 → 0.1.1
.\scripts\bump-version.ps1 minor   # 0.1.0 → 0.2.0
.\scripts\bump-version.ps1 major   # 0.1.0 → 1.0.0
```

```bash
# Unix/Mac - Bump version
bash scripts/bump-version.sh patch   # 0.1.0 → 0.1.1
bash scripts/bump-version.sh minor   # 0.1.0 → 0.2.0
bash scripts/bump-version.sh major   # 0.1.0 → 1.0.0
```

### Release Process

```powershell
# Windows PowerShell - Complete release
.\scripts\release.ps1

# This will:
# - Build documentation (docs/site/)
# - Build Docker images tagged with vX.Y.Z + latest
# - Create git tag vX.Y.Z
# - Print manual push commands
```

```bash
# Unix/Mac - Complete release
bash scripts/release.sh

# This will:
# - Build documentation (docs/site/)
# - Build Docker images tagged with vX.Y.Z + latest
# - Create git tag vX.Y.Z
# - Print manual push commands
```

### Manual Push Commands

After release script completes:
```bash
# Push to git
git push origin main
git push origin v0.1.0

# Push Docker images (optional)
docker push accu-did/resolver:v0.1.0
docker push accu-did/resolver:latest
docker push accu-did/registrar:v0.1.0
docker push accu-did/registrar:latest
```

### Changelog Updates

1. Edit `CHANGELOG.md` before release
2. Add entries under `## [Unreleased]`
3. Release script will commit VERSION and CHANGELOG
4. Manually update links at bottom of CHANGELOG

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
| `docker-compose.yml` | Container orchestration |