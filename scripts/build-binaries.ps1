# scripts/build-binaries.ps1 - Cross-compile binaries for multiple platforms
param(
    [switch]$SkipDarwin = $false
)

$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

# Read version from VERSION file
if (-not (Test-Path "VERSION")) {
    Write-Error "VERSION file not found"
    exit 1
}

$Version = Get-Content "VERSION" -Raw | ForEach-Object { $_.Trim() }
Write-Host "🔨 Building binaries for version: $Version" -ForegroundColor Cyan

# Export build settings
$env:CGO_ENABLED = "0"
$env:GO111MODULE = "on"

# Define platforms to build for
$Platforms = @(
    @{ OS = "linux"; Arch = "amd64" },
    @{ OS = "linux"; Arch = "arm64" },
    @{ OS = "windows"; Arch = "amd64" },
    @{ OS = "darwin"; Arch = "arm64" }
)

# Create dist directories
New-Item -ItemType Directory -Path "dist\bin" -Force | Out-Null

# Build function
function Build-Binary {
    param(
        [string]$OS,
        [string]$Arch,
        [string]$Service,
        [string]$SourceDir
    )

    Write-Host "  Building $Service for $OS/$Arch..." -ForegroundColor Green

    # Create platform directory
    $OutputDir = "dist\bin\$OS-$Arch"
    New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null

    # Set platform-specific variables
    $env:GOOS = $OS
    $env:GOARCH = $Arch

    # Determine binary name
    $BinaryName = $Service
    if ($OS -eq "windows") {
        $BinaryName = "$Service.exe"
    }

    # Build with version embedded
    $LdFlags = "-X main.version=$Version -s -w"

    # Build the binary
    $BuildPath = "$SourceDir\cmd\$Service"
    $OutputPath = "$OutputDir\$BinaryName"

    try {
        & go build -ldflags $LdFlags -o $OutputPath ".\$BuildPath"
        Write-Host "    ✅ Built: $OutputPath" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "    ❌ Failed to build $Service for $OS/$Arch" -ForegroundColor Red
        return $false
    }
}

# Build all binaries
Write-Host ""
Write-Host "Building binaries..." -ForegroundColor Yellow

foreach ($Platform in $Platforms) {
    # Skip darwin if requested or not supported
    if ($Platform.OS -eq "darwin" -and ($SkipDarwin -or -not (Get-Command go -ErrorAction SilentlyContinue))) {
        Write-Host "  Skipping $($Platform.OS)/$($Platform.Arch) (unsupported or skipped)" -ForegroundColor Yellow
        continue
    }

    Write-Host ""
    Write-Host "📦 Platform: $($Platform.OS)/$($Platform.Arch)" -ForegroundColor Cyan

    # Build resolver
    Build-Binary -OS $Platform.OS -Arch $Platform.Arch -Service "resolver" -SourceDir "resolver-go"

    # Build registrar
    Build-Binary -OS $Platform.OS -Arch $Platform.Arch -Service "registrar" -SourceDir "registrar-go"
}

Write-Host ""
Write-Host "🎉 Binary compilation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "📁 Distribution tree:" -ForegroundColor Cyan
Get-ChildItem -Path "dist\bin" -Recurse -File | ForEach-Object { "  $($_.FullName.Replace($PWD.Path, '.'))" }

Write-Host ""
Write-Host "📊 Binary sizes:" -ForegroundColor Cyan
Get-ChildItem -Path "dist\bin" -Recurse -File | ForEach-Object {
    $SizeKB = [math]::Round($_.Length / 1KB, 2)
    "  $SizeKB KB`t$($_.FullName.Replace($PWD.Path, '.'))"
}

# Create checksums
Write-Host ""
Write-Host "🔐 Generating checksums..." -ForegroundColor Yellow
$ChecksumFile = "dist\bin\checksums.sha256"
Get-ChildItem -Path "dist\bin" -Recurse -File -Exclude "checksums.sha256" | ForEach-Object {
    $Hash = Get-FileHash -Path $_.FullName -Algorithm SHA256
    "$($Hash.Hash.ToLower())  $($_.FullName.Replace($PWD.Path + '\', ''))"
} | Out-File -FilePath $ChecksumFile -Encoding UTF8

Write-Host "  ✅ Checksums saved to $ChecksumFile" -ForegroundColor Green

Write-Host ""
Write-Host "[OK] Build complete - binaries ready in dist/bin/" -ForegroundColor Green