# Integration Test Runner for Accumulate DID SDK
# Runs end-to-end integration tests against local devnet

param(
    [Parameter()]
    [string]$ResolverURL = "http://127.0.0.1:8080",

    [Parameter()]
    [string]$RegistrarURL = "http://127.0.0.1:8081",

    [Parameter()]
    [string]$AccNodeURL = "http://127.0.0.1:26656",

    [Parameter()]
    [string]$FaucetURL = "",

    [Parameter()]
    [string]$LiteAccountURL = "",

    [Parameter()]
    [string]$APIKey = "",

    [Parameter()]
    [string]$IdempotencyKey = "",

    [Parameter()]
    [switch]$VerboseOutput = $false
)

function Write-Message($Text, $Type = "Info") {
    switch ($Type) {
        "Success" { Write-Host "[OK] $Text" -ForegroundColor Green }
        "Error"   { Write-Host "[ERROR] $Text" -ForegroundColor Red }
        "Info"    { Write-Host "[INFO] $Text" -ForegroundColor Blue }
        "Status"  { Write-Host "[STATUS] $Text" -ForegroundColor Cyan }
        default   { Write-Host $Text }
    }
}

Write-Message "Accumulate DID SDK Integration Test Runner" "Status"
Write-Message "===========================================" "Status"

# Display configuration
Write-Message "Test configuration:" "Info"
Write-Message "  Resolver URL: $ResolverURL" "Info"
Write-Message "  Registrar URL: $RegistrarURL" "Info"
Write-Message "  Accumulate Node: $AccNodeURL" "Info"

if ($FaucetURL) {
    Write-Message "  Faucet URL: $FaucetURL" "Info"
} else {
    Write-Message "  Faucet URL: (not set)" "Info"
}

if ($LiteAccountURL) {
    Write-Message "  Lite Account: $LiteAccountURL" "Info"
} else {
    Write-Message "  Lite Account: (not set)" "Info"
}

if ($APIKey) {
    Write-Message "  API Key: configured" "Info"
} else {
    Write-Message "  API Key: (not set)" "Info"
}

# Set environment variables
$env:RESOLVER_URL = $ResolverURL
$env:REGISTRAR_URL = $RegistrarURL
$env:ACC_NODE_URL = $AccNodeURL

if ($FaucetURL) {
    $env:ACC_FAUCET_URL = $FaucetURL
}

if ($LiteAccountURL) {
    $env:LITE_ACCOUNT_URL = $LiteAccountURL
}

if ($APIKey) {
    $env:ACCU_API_KEY = $APIKey
}

if ($IdempotencyKey) {
    $env:IDEMPOTENCY_KEY = $IdempotencyKey
}

Write-Message ""
Write-Message "Prerequisites check:" "Status"
Write-Message "  Make sure devnet is running: scripts\devnet.ps1 status" "Info"
Write-Message "  Make sure services are running: resolver :8080, registrar :8081" "Info"
Write-Message ""

# Change to SDK directory
$OriginalPath = Get-Location
try {
    Set-Location "sdks\go\accdid"

    if (-not (Test-Path "integration\integration_test.go")) {
        Write-Message "Integration test file not found" "Error"
        exit 1
    }

    Write-Message "Running integration tests..." "Status"

    # Build test command
    $testCmd = @("go", "test")

    if ($VerboseOutput) {
        $testCmd += "-v"
    }

    $testCmd += @(
        "-run", "TestAccuEndToEnd",
        "-tags=integration",
        "./integration"
    )

    Write-Message "Executing: $($testCmd -join ' ')" "Status"

    # Run the test
    & $testCmd[0] $testCmd[1..($testCmd.Length-1)]

    if ($LASTEXITCODE -eq 0) {
        Write-Message "Integration tests passed!" "Success"
    } else {
        Write-Message "Integration tests failed with exit code: $LASTEXITCODE" "Error"
        exit $LASTEXITCODE
    }

} finally {
    Set-Location $OriginalPath
}

Write-Message ""
Write-Message "Integration test completed successfully" "Success"