# scripts/check-imports.ps1 - Check for forbidden imports
param(
    [switch]$Verbose = $false
)

$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot

Write-Host "🔍 Checking for forbidden imports..." -ForegroundColor Cyan

$ForbiddenFound = $false

# Check for accumulate/internal imports in our services
$ResolverPath = Join-Path $Root "resolver-go"
$RegistrarPath = Join-Path $Root "registrar-go"

if ($Verbose) {
    Write-Host "[INFO] Checking resolver-go at: $ResolverPath"
    Write-Host "[INFO] Checking registrar-go at: $RegistrarPath"
}

# Check for forbidden patterns
$patterns = @(
    "accumulate/internal/",
    "gitlab.com/accumulatenetwork/accumulate/internal"
)

foreach ($pattern in $patterns) {
    $found = $false

    # Check resolver-go
    if (Test-Path $ResolverPath) {
        $resolverMatches = Select-String -Path "$ResolverPath\*.go" -Recurse -Pattern $pattern -ErrorAction SilentlyContinue
        if ($resolverMatches) {
            Write-Host "❌ Forbidden import found in resolver-go: $pattern" -ForegroundColor Red
            if ($Verbose) {
                $resolverMatches | ForEach-Object { Write-Host "   $($_.Filename):$($_.LineNumber): $($_.Line.Trim())" }
            }
            $found = $true
        }
    }

    # Check registrar-go
    if (Test-Path $RegistrarPath) {
        $registrarMatches = Select-String -Path "$RegistrarPath\*.go" -Recurse -Pattern $pattern -ErrorAction SilentlyContinue
        if ($registrarMatches) {
            Write-Host "❌ Forbidden import found in registrar-go: $pattern" -ForegroundColor Red
            if ($Verbose) {
                $registrarMatches | ForEach-Object { Write-Host "   $($_.Filename):$($_.LineNumber): $($_.Line.Trim())" }
            }
            $found = $true
        }
    }

    if ($found) {
        $ForbiddenFound = $true
        Write-Host "   These are internal packages and should not be imported" -ForegroundColor Yellow
    }
}

# Check for local replace directives (except in go.work)
$modules = @($ResolverPath, $RegistrarPath)
foreach ($module in $modules) {
    $goModPath = Join-Path $module "go.mod"
    if (Test-Path $goModPath) {
        $replaceDirectives = Select-String -Path $goModPath -Pattern "^replace.*=>\s*\.\." -ErrorAction SilentlyContinue
        if ($replaceDirectives) {
            Write-Host "⚠️  Warning: Local replace directive found in $goModPath" -ForegroundColor Yellow
            Write-Host "   This may cause issues in production builds" -ForegroundColor Yellow
            if ($Verbose) {
                $replaceDirectives | ForEach-Object { Write-Host "   Line $($_.LineNumber): $($_.Line.Trim())" }
            }
        }
    }
}

Write-Host ""
if ($ForbiddenFound) {
    Write-Host "❌ Import guard check FAILED" -ForegroundColor Red
    exit 1
} else {
    Write-Host "[OK] Import guard passed - no forbidden imports found" -ForegroundColor Green
    exit 0
}