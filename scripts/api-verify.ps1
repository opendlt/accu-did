# scripts/api-verify.ps1
# Validates OpenAPI specifications using redoc-cli
# Usage: powershell -ExecutionPolicy Bypass -File .\scripts\api-verify.ps1

param(
  [switch]$Verbose = $false
)

$ErrorActionPreference = "Stop"

$Root    = Split-Path -Parent $PSScriptRoot
$SpecDir = Join-Path $Root "docs/spec/openapi"

$ResolverYaml  = Join-Path $SpecDir "resolver.yaml"
$RegistrarYaml = Join-Path $SpecDir "registrar.yaml"

if ($Verbose) {
  Write-Host "[INFO] Root directory: $Root"
  Write-Host "[INFO] Spec directory: $SpecDir"
}

function Test-File {
  param([string]$Path, [string]$Name)
  if (-not (Test-Path $Path)) {
    Write-Error "[ERROR] $Name not found: $Path"
    exit 1
  }
  if ($Verbose) {
    Write-Host "[INFO] Found $Name`: $Path"
  }
}

function Lint-OpenAPI {
  param([string]$YamlPath, [string]$Name)

  Write-Host "[INFO] Validating $Name..."

  try {
    # Use npx to run swagger-cli validate command
    # --yes flag automatically installs if not present
    # Exit code 0 = valid, non-zero = validation errors
    $result = & npx --yes swagger-cli validate "$YamlPath" 2>&1

    if ($LASTEXITCODE -eq 0) {
      Write-Host "[OK] $Name is valid" -ForegroundColor Green
      if ($Verbose -and $result) {
        Write-Host $result
      }
    } else {
      Write-Host "[ERROR] $Name validation failed:" -ForegroundColor Red
      Write-Host $result -ForegroundColor Red
      exit $LASTEXITCODE
    }
  } catch {
    Write-Host "[ERROR] Failed to run redoc-cli lint for $Name`: $_" -ForegroundColor Red
    exit 1
  }
}

# Verify files exist
Test-File $ResolverYaml "Resolver OpenAPI spec"
Test-File $RegistrarYaml "Registrar OpenAPI spec"

# Check for npx availability
try {
  $npmVersion = & npx --version 2>&1
  if ($Verbose) {
    Write-Host "[INFO] Found npx version: $npmVersion"
  }
} catch {
  Write-Error "[ERROR] npx not found. Please install Node.js and npm."
  exit 1
}

Write-Host ""
Write-Host "=== OpenAPI Specification Validation ===" -ForegroundColor Cyan
Write-Host ""

# Validate both OpenAPI specifications
Lint-OpenAPI $ResolverYaml "Resolver API (resolver.yaml)"
Lint-OpenAPI $RegistrarYaml "Registrar API (registrar.yaml)"

Write-Host ""
Write-Host "=== Validation Summary ===" -ForegroundColor Cyan
Write-Host "[OK] All OpenAPI specifications are valid" -ForegroundColor Green
Write-Host ""

# Additional checks
$resolverContent = Get-Content -Raw $ResolverYaml
$registrarContent = Get-Content -Raw $RegistrarYaml

# Check for API freeze markers
if ($resolverContent -match "x-api-freeze:\s*true") {
  Write-Host "[OK] Resolver API freeze marker found" -ForegroundColor Green
} else {
  Write-Host "[WARN] Resolver API freeze marker missing" -ForegroundColor Yellow
}

if ($registrarContent -match "x-api-freeze:\s*true") {
  Write-Host "[OK] Registrar API freeze marker found" -ForegroundColor Green
} else {
  Write-Host "[WARN] Registrar API freeze marker missing" -ForegroundColor Yellow
}

# Check version consistency
$resolverVersion = if ($resolverContent -match "version:\s*[`"`']?([^`"`'\r\n]+)[`"`']?") { $matches[1].Trim() } else { "unknown" }
$registrarVersion = if ($registrarContent -match "version:\s*[`"`']?([^`"`'\r\n]+)[`"`']?") { $matches[1].Trim() } else { "unknown" }

Write-Host ""
Write-Host "=== Version Information ===" -ForegroundColor Cyan
Write-Host "Resolver API version:  $resolverVersion"
Write-Host "Registrar API version: $registrarVersion"

if ($resolverVersion -eq $registrarVersion) {
  Write-Host "[OK] API versions are synchronized" -ForegroundColor Green
} else {
  Write-Host "[WARN] API versions are not synchronized" -ForegroundColor Yellow
  Write-Host "       Consider updating to maintain version consistency"
}

Write-Host ""
Write-Host "[OK] API verification completed successfully" -ForegroundColor Green