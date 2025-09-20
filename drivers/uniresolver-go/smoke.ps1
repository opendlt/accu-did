#!/usr/bin/env pwsh

Write-Host "Universal Resolver Driver Smoke Test" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan

$driverUrl = "http://localhost:8081"
$testDid = "did:acc:beastmode.acme"

# Wait for services to be ready
Write-Host "`nWaiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Test health endpoint
Write-Host "`nTesting health endpoint..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "$driverUrl/health" -Method Get
    Write-Host "✓ Health check passed" -ForegroundColor Green
    Write-Host "  Status: $($healthResponse.status)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Health check failed: $_" -ForegroundColor Red
    exit 1
}

# Test driver info endpoint
Write-Host "`nTesting driver info endpoint..." -ForegroundColor Yellow
try {
    $infoResponse = Invoke-RestMethod -Uri "$driverUrl/" -Method Get
    Write-Host "✓ Driver info retrieved" -ForegroundColor Green
    Write-Host "  Driver: $($infoResponse.driver)" -ForegroundColor Gray
    Write-Host "  Version: $($infoResponse.version)" -ForegroundColor Gray
    Write-Host "  Methods: $($infoResponse.methods -join ', ')" -ForegroundColor Gray
} catch {
    Write-Host "✗ Failed to get driver info: $_" -ForegroundColor Red
    exit 1
}

# Test DID resolution
Write-Host "`nTesting DID resolution for: $testDid" -ForegroundColor Yellow
try {
    $resolveUrl = "$driverUrl/1.0/identifiers/$testDid"
    Write-Host "  URL: $resolveUrl" -ForegroundColor Gray

    $response = Invoke-RestMethod -Uri $resolveUrl -Method Get

    if ($response.didDocument) {
        Write-Host "✓ DID resolved successfully" -ForegroundColor Green
        Write-Host "`nDID Document:" -ForegroundColor Cyan
        $response.didDocument | ConvertTo-Json -Depth 10 | Write-Host

        if ($response.didDocumentMetadata) {
            Write-Host "`nDID Document Metadata:" -ForegroundColor Cyan
            $response.didDocumentMetadata | ConvertTo-Json -Depth 10 | Write-Host
        }

        if ($response.didResolutionMetadata) {
            Write-Host "`nDID Resolution Metadata:" -ForegroundColor Cyan
            $response.didResolutionMetadata | ConvertTo-Json -Depth 10 | Write-Host
        }

        Write-Host "`n✓ All tests passed!" -ForegroundColor Green
    } else {
        Write-Host "✗ No DID document returned" -ForegroundColor Red
        $response | ConvertTo-Json -Depth 10 | Write-Host
        exit 1
    }
} catch {
    Write-Host "✗ DID resolution failed: $_" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response: $responseBody" -ForegroundColor Red
    }
    exit 1
}

Write-Host "`n====================================" -ForegroundColor Cyan
Write-Host "Smoke test completed successfully!" -ForegroundColor Green