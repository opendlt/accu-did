# scripts/sdk-openapi-merge.ps1 - Merge OpenAPI specs for SDK generation
param(
    [switch]$Force = $false
)

$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

Write-Host "🔧 Merging OpenAPI specifications for SDK..." -ForegroundColor Cyan

# Read version from VERSION file
if (-not (Test-Path "VERSION")) {
    Write-Error "VERSION file not found"
    exit 1
}

$Ver = Get-Content "VERSION" -Raw | ForEach-Object { $_.Trim() }
Write-Host "   Version: $Ver" -ForegroundColor Green

# Ensure tools directory exists and is built
Write-Host ""
Write-Host "📦 Building sdkmerge tool..." -ForegroundColor Cyan
Set-Location "tools\sdkmerge"

if (-not (Test-Path "go.mod")) {
    go mod init github.com/opendlt/accu-did/tools/sdkmerge
    go mod tidy
}

# Download dependencies if needed
go mod download
go mod tidy

# Build the merge tool
go build -o sdkmerge.exe .

Set-Location $Root

# Check input files exist
$ResolverSpec = "docs\spec\openapi\resolver.yaml"
$RegistrarSpec = "docs\spec\openapi\registrar.yaml"
$OutputSpec = "sdks\spec\openapi\accdid-sdk.yaml"

if (-not (Test-Path $ResolverSpec)) {
    Write-Host "   ⚠️  Resolver OpenAPI spec not found: $ResolverSpec" -ForegroundColor Yellow
    Write-Host "   Creating placeholder..."
    $ResolverDir = Split-Path -Parent $ResolverSpec
    if (-not (Test-Path $ResolverDir)) {
        New-Item -ItemType Directory -Path $ResolverDir -Force | Out-Null
    }

    @"
openapi: 3.0.3
info:
  title: Accumulate DID Resolver
  version: 0.1.0
paths:
  /resolve:
    get:
      summary: Resolve a DID
      parameters:
        - name: did
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: DID resolution result
"@ | Out-File -FilePath $ResolverSpec -Encoding UTF8
}

if (-not (Test-Path $RegistrarSpec)) {
    Write-Host "   ⚠️  Registrar OpenAPI spec not found: $RegistrarSpec" -ForegroundColor Yellow
    Write-Host "   Creating placeholder..."
    $RegistrarDir = Split-Path -Parent $RegistrarSpec
    if (-not (Test-Path $RegistrarDir)) {
        New-Item -ItemType Directory -Path $RegistrarDir -Force | Out-Null
    }

    @"
openapi: 3.0.3
info:
  title: Accumulate DID Registrar
  version: 0.1.0
paths:
  /register:
    post:
      summary: Register a new DID
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        '200':
          description: Registration successful
"@ | Out-File -FilePath $RegistrarSpec -Encoding UTF8
}

# Merge the specifications
Write-Host ""
Write-Host "🔄 Merging OpenAPI specifications..." -ForegroundColor Cyan
& ".\tools\sdkmerge\sdkmerge.exe" $ResolverSpec $RegistrarSpec $OutputSpec $Ver

# Verify output
if (Test-Path $OutputSpec) {
    Write-Host ""
    Write-Host "📄 Generated SDK specification:" -ForegroundColor Cyan
    Write-Host "   File: $OutputSpec"
    $LineCount = (Get-Content $OutputSpec | Measure-Object -Line).Lines
    Write-Host "   Size: $LineCount lines"

    # Show basic info
    $Content = Get-Content $OutputSpec
    $Title = ($Content | Select-String "^\s*title:" | Select-Object -First 1) -replace '.*title:\s*', ''
    $Version = ($Content | Select-String "^\s*version:" | Select-Object -First 1) -replace '.*version:\s*', ''
    $PathCount = ($Content | Select-String "^\s*/.*:$" | Measure-Object).Count

    Write-Host "   Title: $Title"
    Write-Host "   Version: $Version"
    Write-Host "   Paths: $PathCount"
} else {
    Write-Host "   ❌ Failed to generate merged specification" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "[OK] OpenAPI merge complete" -ForegroundColor Green