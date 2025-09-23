# TODO Scanner - Windows PowerShell Wrapper
# Scans the repository for TODO, FIXME, XXX, HACK, and other markers
# Generates reports in JSON, Markdown, and CSV formats

[CmdletBinding()]
param(
    [Parameter(Position=0)]
    [string]$RepoPath = (Get-Location).Path,

    [Parameter()]
    [ValidateSet('auto', 'yes', 'no')]
    [string]$UseDocker = 'auto',

    [Parameter()]
    [string]$DockerImage = 'golang:1.25-alpine',

    [Parameter()]
    [switch]$Help
)

function Show-Usage {
    Write-Host @"
TODO Scanner for Accumulate DID Repository

USAGE:
    .\todo-scan.ps1 [RepoPath] [OPTIONS]

PARAMETERS:
    -RepoPath       Path to repository (default: current directory)
    -UseDocker      Force docker usage: 'auto', 'yes', or 'no' (default: auto)
    -DockerImage    Docker image to use (default: golang:1.25-alpine)
    -Help           Show this help message

EXAMPLES:
    .\todo-scan.ps1                                    # Scan current directory
    .\todo-scan.ps1 C:\path\to\repo                   # Scan specific repository
    .\todo-scan.ps1 -UseDocker yes                    # Force Docker usage
    .\todo-scan.ps1 -UseDocker no                     # Force local Go usage

REPORTS:
    Reports are generated in: .\reports\
    - todo-report.json     # Machine-readable JSON
    - todo-report.md       # Human-readable Markdown
    - todo-report.csv      # Spreadsheet-compatible CSV

PREREQUISITES:
    Either Go (local) or Docker (containerized) must be installed.
"@
}

function Write-ColoredText {
    param(
        [string]$Text,
        [string]$Color = 'White'
    )
    Write-Host $Text -ForegroundColor $Color
}

function Write-StatusMessage {
    param(
        [string]$Icon,
        [string]$Message,
        [string]$Color = 'White'
    )
    Write-Host "$Icon " -ForegroundColor $Color -NoNewline
    Write-Host $Message
}

function Test-Prerequisites {
    Write-StatusMessage "🔍" "Checking prerequisites..."

    $hasGo = $null -ne (Get-Command go -ErrorAction SilentlyContinue)
    $hasDocker = $null -ne (Get-Command docker -ErrorAction SilentlyContinue)

    if ($UseDocker -eq 'auto') {
        if ($hasGo) {
            Write-StatusMessage "✓" "Found local Go installation" "Green"
            $script:UseDocker = 'no'
        } elseif ($hasDocker) {
            Write-StatusMessage "⚠" "No local Go found, using Docker" "Yellow"
            $script:UseDocker = 'yes'
        } else {
            Write-StatusMessage "✗" "Neither Go nor Docker found. Please install one of them." "Red"
            exit 1
        }
    }

    if ($UseDocker -eq 'yes' -and -not $hasDocker) {
        Write-StatusMessage "✗" "Docker not found but UseDocker is 'yes'" "Red"
        exit 1
    }

    if ($UseDocker -eq 'no' -and -not $hasGo) {
        Write-StatusMessage "✗" "Go not found but UseDocker is 'no'" "Red"
        exit 1
    }
}

function Invoke-LocalScan {
    Write-StatusMessage "🔍" "Running TODO scanner locally..." "Blue"

    Push-Location $RepoPath
    try {
        # Ensure reports directory exists
        $outputDir = Join-Path $RepoPath "reports"
        if (-not (Test-Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory -Force | Out-Null
        }

        # Check if scanner exists
        $scannerPath = Join-Path $RepoPath "tools\todoscan\main.go"
        if (-not (Test-Path $scannerPath)) {
            Write-StatusMessage "✗" "tools\todoscan\main.go not found in repository" "Red"
            Write-Host "Please ensure you're running this from the repository root."
            exit 1
        }

        # Run the scanner
        Write-Host "Running: go run tools\todoscan\main.go ."
        & go run "tools\todoscan\main.go" "."

        if ($LASTEXITCODE -ne 0) {
            Write-StatusMessage "✗" "Scanner failed with exit code $LASTEXITCODE" "Red"
            exit $LASTEXITCODE
        }
    }
    finally {
        Pop-Location
    }
}

function Invoke-DockerScan {
    Write-StatusMessage "🐳" "Running TODO scanner in Docker..." "Blue"

    # Ensure reports directory exists
    $outputDir = Join-Path $RepoPath "reports"
    if (-not (Test-Path $outputDir)) {
        New-Item -Path $outputDir -ItemType Directory -Force | Out-Null
    }

    # Check if we have a dev container setup
    $dockerComposePath = Join-Path $RepoPath "docker-compose.dev.yml"
    if (Test-Path $dockerComposePath) {
        Write-StatusMessage "📦" "Using docker-compose.dev.yml" "Blue"
        Push-Location $RepoPath
        try {
            $command = @"
echo 'Running TODO scanner...'
if [[ -f tools/todoscan/main.go ]]; then
    go run tools/todoscan/main.go .
else
    echo 'Error: tools/todoscan/main.go not found'
    exit 1
fi
"@
            & docker-compose -f docker-compose.dev.yml run --rm dev bash -c $command

            if ($LASTEXITCODE -ne 0) {
                Write-StatusMessage "✗" "Docker scan failed with exit code $LASTEXITCODE" "Red"
                exit $LASTEXITCODE
            }
        }
        finally {
            Pop-Location
        }
    } else {
        Write-StatusMessage "🚀" "Using standalone Docker container" "Blue"

        # Convert Windows path to Unix path for Docker
        $unixPath = $RepoPath -replace '\\', '/' -replace '^C:', '/c'

        $command = @"
apk add --no-cache git >/dev/null 2>&1 || true
if [[ -f tools/todoscan/main.go ]]; then
    echo 'Running TODO scanner...'
    go run tools/todoscan/main.go .
else
    echo 'Error: tools/todoscan/main.go not found'
    exit 1
fi
"@

        & docker run --rm -v "${RepoPath}:/workspace" -w /workspace $DockerImage sh -c $command

        if ($LASTEXITCODE -ne 0) {
            Write-StatusMessage "✗" "Docker scan failed with exit code $LASTEXITCODE" "Red"
            exit $LASTEXITCODE
        }
    }
}

function Show-Results {
    $outputDir = Join-Path $RepoPath "reports"
    $jsonFile = Join-Path $outputDir "todo-report.json"
    $mdFile = Join-Path $outputDir "todo-report.md"
    $csvFile = Join-Path $outputDir "todo-report.csv"

    Write-Host ""
    Write-StatusMessage "✓" "Scan completed successfully!" "Green"
    Write-Host ""

    if (Test-Path $jsonFile) {
        try {
            $report = Get-Content $jsonFile -Raw | ConvertFrom-Json
            $totalCount = $report.totalCount
            Write-StatusMessage "📊" "Found $totalCount TODO items" "Blue"

            if ($report.summary.countsByTag) {
                Write-StatusMessage "📋" "Summary by tag:" "Blue"
                $report.summary.countsByTag.PSObject.Properties | ForEach-Object {
                    Write-Host "  - $($_.Name): $($_.Value)"
                }
            }
        }
        catch {
            Write-StatusMessage "⚠" "Could not parse JSON report for summary" "Yellow"
        }
    }

    Write-Host ""
    Write-StatusMessage "📁" "Reports generated:" "Green"

    @($jsonFile, $mdFile, $csvFile) | ForEach-Object {
        if (Test-Path $_) {
            $size = [math]::Round((Get-Item $_).Length / 1KB, 1)
            $filename = Split-Path $_ -Leaf
            Write-StatusMessage "✓" "$filename (${size} KB)" "Green"
        } else {
            $filename = Split-Path $_ -Leaf
            Write-StatusMessage "✗" "$filename (missing)" "Red"
        }
    }

    Write-Host ""
    Write-StatusMessage "💡" "Next steps:" "Blue"
    Write-Host "  - Review the Markdown report: $mdFile"
    Write-Host "  - Import CSV data: $csvFile"
    Write-Host "  - Process JSON programmatically: $jsonFile"
    Write-Host "  - Filter by tag: Select-String 'TODO' $mdFile"
    Write-Host "  - Filter by directory: Select-String 'resolver-go' $mdFile"
}

function Main {
    if ($Help) {
        Show-Usage
        return
    }

    Write-StatusMessage "🔍" "TODO Scanner for Accumulate DID Repository" "Blue"
    Write-StatusMessage "📂" "Repository: $RepoPath" "Blue"
    Write-Host ""

    # Verify repository path exists
    if (-not (Test-Path $RepoPath)) {
        Write-StatusMessage "✗" "Repository path does not exist: $RepoPath" "Red"
        exit 1
    }

    Test-Prerequisites

    # Run the scanner
    if ($UseDocker -eq 'yes') {
        Invoke-DockerScan
    } else {
        Invoke-LocalScan
    }

    # Print results summary
    Show-Results
}

# Only run main if script is executed directly (not dot-sourced)
if ($MyInvocation.InvocationName -eq $MyInvocation.MyCommand.Name) {
    Main
}