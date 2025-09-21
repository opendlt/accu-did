#!/usr/bin/env pwsh
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $MyInvocation.MyCommand.Path | Split-Path -Parent
$VersionFile = Join-Path $Root "VERSION"

Write-Host "🚀 Starting release process..." -ForegroundColor Cyan

# Read version
if (!(Test-Path $VersionFile)) {
    Write-Error "VERSION file not found"
    exit 1
}

$Version = (Get-Content $VersionFile).Trim()
$Tag = "v$Version"

Write-Host "📦 Releasing version: $Tag" -ForegroundColor Yellow

# Check for uncommitted changes
$GitStatus = git status --porcelain
if ($GitStatus) {
    Write-Host "⚠️  Uncommitted changes detected. Committing VERSION and CHANGELOG..." -ForegroundColor Yellow
    git add VERSION CHANGELOG.md
    git commit -m "chore(release): prepare $Tag" 2>$null || Write-Host "Nothing to commit"
}

# Build documentation
Write-Host "`n📚 Building documentation..." -ForegroundColor Yellow
& "$Root/scripts/build-docs.ps1"
if ($LASTEXITCODE -ne 0) {
    Write-Error "Documentation build failed"
    exit 1
}

# Build Docker images
if (Get-Command docker -ErrorAction SilentlyContinue) {
    Write-Host "`n🐳 Building Docker images..." -ForegroundColor Yellow

    # Build resolver
    docker build -t "accu-did/resolver:$Tag" -t "accu-did/resolver:latest" `
        -f "$Root/drivers/resolver/Dockerfile" "$Root"
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Resolver Docker build failed"
        exit 1
    }

    # Build registrar
    docker build -t "accu-did/registrar:$Tag" -t "accu-did/registrar:latest" `
        -f "$Root/drivers/registrar/Dockerfile" "$Root"
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Registrar Docker build failed"
        exit 1
    }

    Write-Host "✅ Docker images built and tagged" -ForegroundColor Green
}

# Create git tag
Write-Host "`n🏷️  Creating git tag: $Tag" -ForegroundColor Yellow
git tag -a $Tag -m "Release $Tag"
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to create git tag"
    exit 1
}

# Success message
Write-Host "`n" + ("=" * 60) -ForegroundColor Cyan
Write-Host "✅ Release $Tag prepared successfully!" -ForegroundColor Green
Write-Host "`nNext steps:" -ForegroundColor Yellow
Write-Host "  git push origin main" -ForegroundColor White
Write-Host "  git push origin $Tag" -ForegroundColor White
if (Get-Command docker -ErrorAction SilentlyContinue) {
    Write-Host "`nOptional - Push Docker images:" -ForegroundColor Yellow
    Write-Host "  docker push accu-did/resolver:$Tag" -ForegroundColor White
    Write-Host "  docker push accu-did/resolver:latest" -ForegroundColor White
    Write-Host "  docker push accu-did/registrar:$Tag" -ForegroundColor White
    Write-Host "  docker push accu-did/registrar:latest" -ForegroundColor White
}