# Accumulate Devnet Management Script

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
$NodeRpcUrl = "http://127.0.0.1:$BasePort"

function Write-Message($Text, $Type = "Info") {
    switch ($Type) {
        "Success" { Write-Host "[OK] $Text" -ForegroundColor Green }
        "Error"   { Write-Host "[ERROR] $Text" -ForegroundColor Red }
        "Info"    { Write-Host "[INFO] $Text" -ForegroundColor Blue }
        "Status"  { Write-Host "[STATUS] $Text" -ForegroundColor Cyan }
        default   { Write-Host $Text }
    }
}

function Test-AccumulateRepo {
    if (-not (Test-Path $AccumulateRepoPath)) {
        Write-Message "Accumulate repo not found at: $AccumulateRepoPath" "Error"
        Write-Message "Expected sibling repo structure with accumulate/ directory" "Info"
        return $false
    }

    $accumulatedCmd = Join-Path $AccumulateRepoPath "cmd\accumulated"
    if (-not (Test-Path $accumulatedCmd)) {
        Write-Message "Accumulated command directory not found in repo" "Error"
        return $false
    }

    return $true
}

function Test-DevnetRunning {
    try {
        $response = Invoke-WebRequest -Uri "$NodeRpcUrl/status" -TimeoutSec 5 -ErrorAction SilentlyContinue
        return ($response.StatusCode -eq 200)
    } catch {
        return $false
    }
}

function Get-AccumulatedProcesses {
    return Get-Process -Name "accumulated" -ErrorAction SilentlyContinue
}

function Start-Devnet {
    Write-Message "Starting Accumulate devnet..." "Status"

    if (-not (Test-AccumulateRepo)) {
        return $false
    }

    # Check if already running
    if (Test-DevnetRunning) {
        Write-Message "Devnet already running" "Info"
        Write-Message "Node RPC URL: $NodeRpcUrl" "Success"
        $env:ACC_NODE_URL = $NodeRpcUrl
        return $true
    }

    Push-Location $AccumulateRepoPath
    try {
        # Initialize devnet
        Write-Message "Initializing devnet (this may take a moment)..." "Status"
        $initOutput = & go run ./cmd/accumulated init devnet -w $DevnetWorkDir --reset --faucet-seed $FaucetSeed 2>&1

        if ($LASTEXITCODE -ne 0) {
            Write-Message "Failed to initialize devnet" "Error"
            Write-Host $initOutput
            return $false
        }

        # Start devnet in background
        Write-Message "Starting devnet processes..." "Status"
        $job = Start-Job -ScriptBlock {
            param($RepoPath, $WorkDir, $Port, $Seed)
            Set-Location $RepoPath
            & go run ./cmd/accumulated run devnet -w $WorkDir --port $Port --faucet-seed $Seed
        } -ArgumentList (Get-Location), $DevnetWorkDir, $BasePort, $FaucetSeed

        # Wait for startup
        Write-Message "Waiting for devnet to become ready..." "Status"
        $timeout = 45
        $elapsed = 0

        while ($elapsed -lt $timeout) {
            Start-Sleep -Seconds 3
            $elapsed += 3

            if (Test-DevnetRunning) {
                Write-Message "Devnet started successfully!" "Success"
                Write-Message "Node RPC URL: $NodeRpcUrl" "Success"

                # Set environment variable
                $env:ACC_NODE_URL = $NodeRpcUrl
                Write-Message "Environment variable set: ACC_NODE_URL=$NodeRpcUrl" "Info"
                Write-Message "Use 'scripts\devnet.ps1 status' to check devnet health" "Info"

                return $true
            }

            # Check if job failed
            if ($job.State -eq "Failed") {
                Write-Message "Devnet startup failed" "Error"
                Remove-Job -Job $job -Force
                return $false
            }

            Write-Host "." -NoNewline
        }

        Write-Host ""
        Write-Message "Timeout waiting for devnet to start (${timeout}s)" "Error"
        Stop-Job -Job $job -ErrorAction SilentlyContinue
        Remove-Job -Job $job -Force -ErrorAction SilentlyContinue
        return $false

    } finally {
        Pop-Location
    }
}

function Stop-Devnet {
    Write-Message "Stopping Accumulate devnet..." "Status"

    # Kill accumulated processes
    $processes = Get-AccumulatedProcesses
    if ($processes) {
        Write-Message "Terminating devnet processes..." "Status"
        $processes | Stop-Process -Force -ErrorAction SilentlyContinue
        Start-Sleep -Seconds 2
    }

    # Clean up background jobs
    $jobs = Get-Job | Where-Object { $_.Command -like "*accumulated*" }
    if ($jobs) {
        Write-Message "Cleaning up background jobs..." "Status"
        $jobs | Stop-Job -ErrorAction SilentlyContinue
        $jobs | Remove-Job -Force -ErrorAction SilentlyContinue
    }

    Write-Message "Devnet stopped" "Success"

    # Clear environment variable
    if ($env:ACC_NODE_URL) {
        Remove-Item Env:ACC_NODE_URL -ErrorAction SilentlyContinue
        Write-Message "Cleared environment variable: ACC_NODE_URL" "Info"
    }

    return $true
}

function Show-DevnetStatus {
    Write-Message "Checking Accumulate devnet status..." "Status"

    $processes = Get-AccumulatedProcesses
    $isRunning = Test-DevnetRunning

    if ($isRunning -and $processes) {
        Write-Message "Devnet is running (PID: $($processes[0].Id))" "Success"
        Write-Message "Node RPC URL: $NodeRpcUrl" "Info"
        Write-Message "Environment: ACC_NODE_URL=$env:ACC_NODE_URL" "Info"

        # Test RPC health
        Write-Message "Testing RPC endpoint..." "Status"
        if (Test-DevnetRunning) {
            Write-Message "RPC endpoint is healthy" "Success"
        } else {
            Write-Message "RPC endpoint not responding" "Error"
        }

    } elseif ($processes) {
        Write-Message "Accumulated processes found but devnet not responding" "Error"
        Write-Message "Process PID: $($processes[0].Id)" "Info"
    } else {
        Write-Message "Devnet is not running" "Info"
        Write-Message "Use 'scripts\devnet.ps1 up' to start devnet" "Info"
    }
}

# Main execution
switch ($Command) {
    "up" {
        $success = Start-Devnet
        if ($success) { exit 0 } else { exit 1 }
    }
    "down" {
        $success = Stop-Devnet
        if ($success) { exit 0 } else { exit 1 }
    }
    "status" {
        Show-DevnetStatus
        exit 0
    }
}