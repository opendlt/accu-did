#!/usr/bin/env pwsh

# Accumulate DID Hello World Smoke Test

Write-Host "=== Accumulate DID Hello World Smoke Test ===" -ForegroundColor Green

# Check prerequisites
Write-Host "`n1. Checking prerequisites..." -ForegroundColor Yellow

# Check Go installation
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Go is not installed or not in PATH" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Go is installed: $(go version)"

# Check ACC_NODE_URL
if (!$env:ACC_NODE_URL) {
    Write-Host "❌ ACC_NODE_URL environment variable is not set" -ForegroundColor Red
    Write-Host "   Set it with: `$env:ACC_NODE_URL='http://localhost:26657'" -ForegroundColor Cyan
    exit 1
}
Write-Host "✅ ACC_NODE_URL is set: $env:ACC_NODE_URL"

# Test node connectivity
Write-Host "`n2. Testing Accumulate node connectivity..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$env:ACC_NODE_URL/status" -Method GET -TimeoutSec 5
    Write-Host "✅ Accumulate node is reachable"
} catch {
    Write-Host "❌ Cannot reach Accumulate node at $env:ACC_NODE_URL" -ForegroundColor Red
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Initialize Go module if needed
Write-Host "`n3. Setting up Go module..." -ForegroundColor Yellow
if (!(Test-Path "go.mod")) {
    go mod init hello_accu
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to initialize Go module" -ForegroundColor Red
        exit 1
    }
    Write-Host "✅ Initialized Go module"
}

# Download dependencies
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to download dependencies" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Dependencies ready"

# Run the hello_accu example
Write-Host "`n4. Running hello_accu example..." -ForegroundColor Yellow
Write-Host "---" -ForegroundColor DarkGray

$output = go run main.go 2>&1
$exitCode = $LASTEXITCODE

Write-Host $output
Write-Host "---" -ForegroundColor DarkGray

if ($exitCode -eq 0) {
    Write-Host "`n✅ Hello Accu example completed successfully!" -ForegroundColor Green

    # Extract DID from output
    $didMatch = $output | Select-String "DID: (did:acc:[\w\.]+)"
    if ($didMatch) {
        $did = $didMatch.Matches[0].Groups[1].Value
        Write-Host "`n🎉 DID created: $did" -ForegroundColor Cyan
        Write-Host "   You can now resolve this DID using the Accumulate DID resolver" -ForegroundColor Gray
    }
} else {
    Write-Host "`n❌ Hello Accu example failed with exit code $exitCode" -ForegroundColor Red
    exit $exitCode
}

Write-Host "`n=== Smoke Test Complete ===" -ForegroundColor Green