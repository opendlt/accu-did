# scripts/conformance.ps1 - Run conformance tests against resolver and registrar
param(
    [string]$ResolverURL = "http://127.0.0.1:8080",
    [string]$RegistrarURL = "http://127.0.0.1:8081",
    [string]$TestDID = "did:acc:conformance-test",
    [string]$APIKey = ""
)

$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot

Write-Host "🔍 Running conformance tests..." -ForegroundColor Cyan
Write-Host "   Resolver:  $ResolverURL"
Write-Host "   Registrar: $RegistrarURL"
Write-Host "   Test DID:  $TestDID"

# Set environment variables
$env:RESOLVER_URL = $ResolverURL
$env:REGISTRAR_URL = $RegistrarURL
$env:TEST_DID = $TestDID
if ($APIKey) {
    $env:REGISTRAR_API_KEY = $APIKey
}

# Check if conformance tool exists
$ConformancePath = Join-Path $Root "tools\conformance\conformance.go"
if (Test-Path $ConformancePath) {
    Write-Host "Running conformance tool..." -ForegroundColor Green

    # Build and run the conformance tool
    Push-Location (Join-Path $Root "tools\conformance")
    try {
        go run conformance.go
        $exitCode = $LASTEXITCODE
    } finally {
        Pop-Location
    }

    exit $exitCode
} else {
    Write-Host "⚠️  Conformance tool not found at $ConformancePath" -ForegroundColor Yellow
    Write-Host "   Running basic health checks instead..." -ForegroundColor Yellow

    # Run basic health checks instead
    Write-Host ""
    Write-Host "Running basic health checks..." -ForegroundColor Cyan

    $allPassed = $true

    try {
        $resolverHealth = Invoke-RestMethod -Uri "$ResolverURL/healthz" -Method Get -TimeoutSec 10
        Write-Host "✅ Resolver health check passed" -ForegroundColor Green
    } catch {
        Write-Host "❌ Resolver health check failed: $($_.Exception.Message)" -ForegroundColor Red
        $allPassed = $false
    }

    try {
        $registrarHealth = Invoke-RestMethod -Uri "$RegistrarURL/healthz" -Method Get -TimeoutSec 10
        Write-Host "✅ Registrar health check passed" -ForegroundColor Green
    } catch {
        Write-Host "❌ Registrar health check failed: $($_.Exception.Message)" -ForegroundColor Red
        $allPassed = $false
    }

    Write-Host ""
    if ($allPassed) {
        Write-Host "[OK] Basic health checks passed" -ForegroundColor Green
        exit 0
    } else {
        Write-Host "[FAIL] Some health checks failed" -ForegroundColor Red
        exit 1
    }
}