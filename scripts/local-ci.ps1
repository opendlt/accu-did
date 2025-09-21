# scripts/local-ci.ps1
$ErrorActionPreference = "Stop"

Write-Host "==> Unit tests (FAKE mode)"
Push-Location resolver-go
go mod tidy
go test ./...
Pop-Location

Push-Location registrar-go
go mod tidy
go test ./...
Pop-Location

Write-Host "==> Build docs (Redoc)"
# You can comment this out if you only want tests
powershell -ExecutionPolicy Bypass -File "$PSScriptRoot\build-docs.ps1"

Write-Host "==> Docker build (optional)"
try {
  docker --version | Out-Null
  docker build -t accu-did/resolver:local ./resolver-go
  docker build -t accu-did/registrar:local ./registrar-go
} catch {
  Write-Warning "Docker not available; skipping image builds."
}

Write-Host "[OK] LOCAL CI PASSED"