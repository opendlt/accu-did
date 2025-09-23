#!/usr/bin/env pwsh
# Accumulate Devnet Management Script
# Manages local Accumulate devnet for DID development

param(
    [Parameter(Mandatory = $true, Position = 0)]
    [ValidateSet("up", "down", "status")]
    [string]$Command,

    [Parameter()]
    [int]$BasePort = 26656,

    [Parameter()]
    [string]$FaucetSeed = "ci"
)

# Configuration
$AccumulateRepoPath = "..\accumulate"
$DevnetWorkDir = ".nodes"
$ProcessName = "accumulated"

# Derived URLs based on Accumulate devnet port conventions
$NodeRpcUrl = "http://127.0.0.1:$BasePort"
$FaucetUrl = "http://127.0.0.1:$($BasePort + 3)"  # Typically RPC + 3 for faucet

function Write-Status {
    param([string]$Message, [string]$Color = "White")
    Write-Host "🔧 " -NoNewline -ForegroundColor Cyan
    Write-Host $Message -ForegroundColor $Color
}

function Write-Success {
    param([string]$Message)
    Write-Host "✅ " -NoNewline -ForegroundColor Green
    Write-Host $Message -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host "❌ " -NoNewline -ForegroundColor Red
    Write-Host $Message -ForegroundColor Red
}

function Write-Info {
    param([string]$Message)
    Write-Host "ℹ️  " -NoNewline -ForegroundColor Blue
    Write-Host $Message -ForegroundColor Blue
}

function Test-AccumulateRepo {
    if (-not (Test-Path $AccumulateRepoPath)) {
        Write-Error "Accumulate repo not found at: $AccumulateRepoPath"
        Write-Info "Expected sibling repo at: $(Resolve-Path $AccumulateRepoPath -ErrorAction SilentlyContinue)"
        return $false
    }

    $accumulatedPath = Join-Path $AccumulateRepoPath "cmd\accumulated"
    if (-not (Test-Path $accumulatedPath)) {
        Write-Error "Accumulated command not found in repo"
        return $false
    }

    return $true
}

function Get-DevnetStatus {
    # Check if accumulated process is running with devnet
    $processes = Get-Process -Name $ProcessName -ErrorAction SilentlyContinue | Where-Object {
        $_.CommandLine -like "*run devnet*" -or $_.ProcessName -eq $ProcessName
    }

    if ($processes) {
        return @{
            Running = $true
            ProcessId = $processes[0].Id
            NodeRpcUrl = $NodeRpcUrl
            FaucetUrl = $FaucetUrl
        }
    }

    return @{
        Running = $false
        ProcessId = $null
        NodeRpcUrl = $null
        FaucetUrl = $null
    }
}

function Start-Devnet {
    Write-Status "Starting Accumulate devnet..."

    if (-not (Test-AccumulateRepo)) {
        return $false
    }

    # Check if already running
    $status = Get-DevnetStatus
    if ($status.Running) {
        Write-Info "Devnet already running (PID: $($status.ProcessId))"
        Write-Success "Node RPC URL: $($status.NodeRpcUrl)"
        Write-Success "Faucet URL: $($status.FaucetUrl)"
        return $true
    }

    Push-Location $AccumulateRepoPath
    try {
        # Initialize devnet (with reset to ensure clean state)
        Write-Status "Initializing devnet..."
        $initResult = & go run ./cmd/accumulated init devnet -w $DevnetWorkDir --reset --faucet-seed $FaucetSeed 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to initialize devnet: $initResult"
            return $false
        }

        # Start devnet in background
        Write-Status "Starting devnet processes..."
        $runJob = Start-Job -ScriptBlock {
            param($RepoPath, $WorkDir, $Port, $Seed)
            Set-Location $RepoPath
            & go run ./cmd/accumulated run devnet -w $WorkDir --port $Port --faucet-seed $Seed
        } -ArgumentList (Get-Location), $DevnetWorkDir, $BasePort, $FaucetSeed

        # Wait for startup with timeout
        Write-Status "Waiting for devnet to start..."
        $timeout = 30
        $elapsed = 0

        do {
            Start-Sleep -Seconds 2
            $elapsed += 2

            # Test if RPC endpoint is responding
            try {
                $response = Invoke-WebRequest -Uri "$NodeRpcUrl/status" -TimeoutSec 5 -ErrorAction SilentlyContinue
                if ($response.StatusCode -eq 200) {
                    Write-Success "Devnet started successfully!"
                    Write-Success "Node RPC URL: $NodeRpcUrl"
                    Write-Success "Faucet URL: $FaucetUrl"

                    # Set environment variable for other scripts
                    $env:ACC_NODE_URL = $NodeRpcUrl

                    Write-Info "Environment variable set: ACC_NODE_URL=$NodeRpcUrl"
                    Write-Info "Use 'scripts\devnet.ps1 status' to check devnet health"
                    return $true
                }
            }
            catch {
                # Continue waiting
            }

            # Check if job failed
            if ($runJob.State -eq "Failed") {
                $jobError = Receive-Job -Job $runJob 2>&1
                Write-Error "Devnet startup failed: $jobError"
                Remove-Job -Job $runJob -Force
                return $false
            }

        } while ($elapsed -lt $timeout)

        Write-Error "Timeout waiting for devnet to start (${timeout}s)"
        Stop-Job -Job $runJob -ErrorAction SilentlyContinue
        Remove-Job -Job $runJob -Force -ErrorAction SilentlyContinue
        return $false
    }
    finally {
        Pop-Location
    }
}

function Stop-Devnet {
    Write-Status "Stopping Accumulate devnet..."

    $status = Get-DevnetStatus
    if (-not $status.Running) {
        Write-Info "Devnet is not running"
        return $true
    }

    # Kill accumulated processes
    try {
        $processes = Get-Process -Name $ProcessName -ErrorAction SilentlyContinue
        if ($processes) {
            Write-Status "Terminating devnet processes..."
            $processes | Stop-Process -Force
            Start-Sleep -Seconds 2
        }

        # Clean up background jobs
        $jobs = Get-Job | Where-Object { $_.Command -like "*accumulated*" }
        if ($jobs) {
            Write-Status "Cleaning up background jobs..."
            $jobs | Stop-Job -ErrorAction SilentlyContinue
            $jobs | Remove-Job -Force -ErrorAction SilentlyContinue
        }

        Write-Success "Devnet stopped successfully"

        # Clear environment variable
        if ($env:ACC_NODE_URL) {
            Remove-Item Env:ACC_NODE_URL -ErrorAction SilentlyContinue
            Write-Info "Cleared environment variable: ACC_NODE_URL"
        }

        return $true
    }
    catch {
        Write-Error "Error stopping devnet: $($_.Exception.Message)"
        return $false
    }
}

function Show-DevnetStatus {
    Write-Status "Checking Accumulate devnet status..."

    $status = Get-DevnetStatus

    if ($status.Running) {
        Write-Success "Devnet is running (PID: $($status.ProcessId))"
        Write-Info "Node RPC URL: $($status.NodeRpcUrl)"
        Write-Info "Faucet URL: $($status.FaucetUrl)"
        Write-Info "Environment: ACC_NODE_URL=$env:ACC_NODE_URL"

        # Test RPC endpoint health
        try {
            Write-Status "Testing RPC endpoint..."
            $response = Invoke-WebRequest -Uri "$($status.NodeRpcUrl)/status" -TimeoutSec 10
            if ($response.StatusCode -eq 200) {
                Write-Success "RPC endpoint is healthy"
            }
        }
        catch {
            Write-Error "RPC endpoint not responding: $($_.Exception.Message)"
        }

        # Test faucet endpoint health
        try {
            Write-Status "Testing faucet endpoint..."
            $response = Invoke-WebRequest -Uri "$($status.FaucetUrl)/health" -TimeoutSec 10 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                Write-Success "Faucet endpoint is healthy"
            }
        }
        catch {
            Write-Info "Faucet endpoint check skipped (may not have /health endpoint)"
        }
    }
    else {
        Write-Info "Devnet is not running"
        Write-Info "Use 'scripts\devnet.ps1 up' to start devnet"
    }
}

# Main execution
switch ($Command) {
    "up" {
        $success = Start-Devnet
        exit ($success ? 0 : 1)
    }
    "down" {
        $success = Stop-Devnet
        exit ($success ? 0 : 1)
    }
    "status" {
        Show-DevnetStatus
        exit 0
    }
}