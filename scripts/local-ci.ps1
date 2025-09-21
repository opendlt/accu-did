#!/usr/bin/env pwsh
# Local CI script for Windows
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

Write-Host "🚀 Starting local CI build..." -ForegroundColor Cyan

$failed = $false
$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path | Split-Path -Parent

# Test resolver-go
Write-Host "`n📋 Testing resolver-go..." -ForegroundColor Yellow
Push-Location "$rootDir/resolver-go"
try {
    & go test ./... -v
    if ($LASTEXITCODE -ne 0) { throw "Resolver tests failed" }
    Write-Host "✅ Resolver tests passed" -ForegroundColor Green
} catch {
    Write-Host "❌ Resolver tests failed: $_" -ForegroundColor Red
    $failed = $true
} finally {
    Pop-Location
}

# Test registrar-go
Write-Host "`n📋 Testing registrar-go..." -ForegroundColor Yellow
Push-Location "$rootDir/registrar-go"
try {
    & go test ./... -v
    if ($LASTEXITCODE -ne 0) { throw "Registrar tests failed" }
    Write-Host "✅ Registrar tests passed" -ForegroundColor Green
} catch {
    Write-Host "❌ Registrar tests failed: $_" -ForegroundColor Red
    $failed = $true
} finally {
    Pop-Location
}

# Build documentation
Write-Host "`n📚 Building documentation..." -ForegroundColor Yellow
try {
    & "$rootDir/scripts/build-docs.ps1"
    if ($LASTEXITCODE -ne 0) { throw "Docs build failed" }
    Write-Host "✅ Documentation built" -ForegroundColor Green
} catch {
    Write-Host "❌ Documentation build failed: $_" -ForegroundColor Red
    $failed = $true
}

# Build Docker images (optional)
if (Get-Command docker -ErrorAction SilentlyContinue) {
    Write-Host "`n🐳 Building Docker images..." -ForegroundColor Yellow

    # Build resolver image
    try {
        docker build -t accu-did/resolver:latest -f "$rootDir/drivers/resolver/Dockerfile" "$rootDir"
        if ($LASTEXITCODE -ne 0) { throw "Resolver Docker build failed" }
        Write-Host "✅ Resolver Docker image built" -ForegroundColor Green
    } catch {
        Write-Host "❌ Resolver Docker build failed: $_" -ForegroundColor Red
        $failed = $true
    }

    # Build registrar image
    try {
        docker build -t accu-did/registrar:latest -f "$rootDir/drivers/registrar/Dockerfile" "$rootDir"
        if ($LASTEXITCODE -ne 0) { throw "Registrar Docker build failed" }
        Write-Host "✅ Registrar Docker image built" -ForegroundColor Green
    } catch {
        Write-Host "❌ Registrar Docker build failed: $_" -ForegroundColor Red
        $failed = $true
    }
} else {
    Write-Host "`n⚠️  Docker not found, skipping image builds" -ForegroundColor Yellow
}

# Summary
Write-Host "`n" + ("=" * 60) -ForegroundColor Cyan
if ($failed) {
    Write-Host "❌ LOCAL CI FAILED" -ForegroundColor Red
    exit 1
} else {
    Write-Host "✅ LOCAL CI PASSED" -ForegroundColor Green
    exit 0
}