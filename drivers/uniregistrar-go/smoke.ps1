#!/usr/bin/env pwsh

Write-Host "Universal Registrar Driver Smoke Test" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan

$driverUrl = "http://localhost:8083"

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
    Write-Host "  Operations: $($infoResponse.operations -join ', ')" -ForegroundColor Gray
} catch {
    Write-Host "✗ Failed to get driver info: $_" -ForegroundColor Red
    exit 1
}

# Test DID creation
Write-Host "`nTesting DID creation..." -ForegroundColor Yellow

# Generate unique DID for test
$timestamp = [DateTimeOffset]::Now.ToUnixTimeSeconds()
$testDid = "did:acc:test$timestamp"

$createRequest = @{
    did = $testDid
    didDocument = @{
        "@context" = @("https://www.w3.org/ns/did/v1")
        id = $testDid
        verificationMethod = @(
            @{
                id = "$testDid#key-1"
                type = "AccumulateKeyPage"
                controller = $testDid
                keyPageUrl = "acc://test$timestamp/book/1"
                threshold = 1
            }
        )
    }
    options = @{}
    secret = @{}
} | ConvertTo-Json -Depth 10

try {
    $createUrl = "$driverUrl/1.0/create?method=acc"
    Write-Host "  URL: $createUrl" -ForegroundColor Gray
    Write-Host "  DID: $testDid" -ForegroundColor Gray

    $headers = @{
        "Content-Type" = "application/json"
    }

    $response = Invoke-RestMethod -Uri $createUrl -Method Post -Body $createRequest -Headers $headers

    if ($response.didState -and $response.didState.state -eq "finished") {
        Write-Host "✓ DID created successfully" -ForegroundColor Green

        Write-Host "`nDID State:" -ForegroundColor Cyan
        $response.didState | ConvertTo-Json -Depth 10 | Write-Host

        if ($response.didRegistrationMetadata) {
            Write-Host "`nDID Registration Metadata:" -ForegroundColor Cyan
            Write-Host "  Version ID: $($response.didRegistrationMetadata.versionId)" -ForegroundColor Gray
            Write-Host "  Content Hash: $($response.didRegistrationMetadata.contentHash)" -ForegroundColor Gray
            Write-Host "  Transaction ID: $($response.didRegistrationMetadata.txId)" -ForegroundColor Gray
        }

        if ($response.jobId) {
            Write-Host "`nJob ID: $($response.jobId)" -ForegroundColor Cyan
        }

        Write-Host "`n✓ All tests passed!" -ForegroundColor Green
    } else {
        Write-Host "✗ DID creation failed" -ForegroundColor Red
        $response | ConvertTo-Json -Depth 10 | Write-Host
        exit 1
    }
} catch {
    Write-Host "✗ DID creation request failed: $_" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response: $responseBody" -ForegroundColor Red
    }
    exit 1
}

# Optional: Test update and deactivate with stub requests
Write-Host "`n=====================================" -ForegroundColor Cyan
Write-Host "Testing additional operations..." -ForegroundColor Yellow

# Test Update endpoint availability
Write-Host "`nTesting update endpoint..." -ForegroundColor Yellow
$updateRequest = @{
    did = "did:acc:alice"
    didDocument = @{
        "@context" = @("https://www.w3.org/ns/did/v1")
        id = "did:acc:alice"
    }
} | ConvertTo-Json -Depth 10

try {
    $updateUrl = "$driverUrl/1.0/update?method=acc"
    $response = Invoke-RestMethod -Uri $updateUrl -Method Post -Body $updateRequest -Headers @{"Content-Type"="application/json"}
    Write-Host "✓ Update endpoint is responsive" -ForegroundColor Green
} catch {
    Write-Host "⚠ Update endpoint returned error (expected in test environment)" -ForegroundColor Yellow
}

# Test Deactivate endpoint availability
Write-Host "`nTesting deactivate endpoint..." -ForegroundColor Yellow
$deactivateRequest = @{
    did = "did:acc:alice"
} | ConvertTo-Json

try {
    $deactivateUrl = "$driverUrl/1.0/deactivate?method=acc"
    $response = Invoke-RestMethod -Uri $deactivateUrl -Method Post -Body $deactivateRequest -Headers @{"Content-Type"="application/json"}
    Write-Host "✓ Deactivate endpoint is responsive" -ForegroundColor Green
} catch {
    Write-Host "⚠ Deactivate endpoint returned error (expected in test environment)" -ForegroundColor Yellow
}

Write-Host "`n=====================================" -ForegroundColor Cyan
Write-Host "Smoke test completed successfully!" -ForegroundColor Green