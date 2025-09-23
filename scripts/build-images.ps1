# scripts/build-images.ps1 - Build Docker images locally
param(
    [switch]$MultiArch = $false
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
Write-Host "🐳 Building Docker images for version: $Version" -ForegroundColor Cyan

# Check if docker is available
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Error "docker is required but not available"
    exit 1
}

# Check buildx availability for multi-arch
$HasBuildx = $false
try {
    docker buildx version | Out-Null
    $HasBuildx = $true
    Write-Host "🔧 docker buildx available" -ForegroundColor Green
}
catch {
    Write-Host "⚠️  docker buildx not available, building single-arch only" -ForegroundColor Yellow
}

# Create images directory
New-Item -ItemType Directory -Path "dist\images\resolver" -Force | Out-Null
New-Item -ItemType Directory -Path "dist\images\registrar" -Force | Out-Null

# Build function
function Build-Image {
    param(
        [string]$Service
    )

    $DockerfilePath = "$Service-go\Dockerfile"
    $ContextPath = "$Service-go"

    Write-Host ""
    Write-Host "📦 Building $Service image..." -ForegroundColor Cyan

    # Verify Dockerfile exists
    if (-not (Test-Path $DockerfilePath)) {
        Write-Host "❌ Dockerfile not found: $DockerfilePath" -ForegroundColor Red
        return $false
    }

    try {
        if ($HasBuildx -and $MultiArch) {
            # Multi-arch build with buildx
            Write-Host "   Building for linux/amd64,linux/arm64..." -ForegroundColor Green

            & docker buildx build `
                --platform linux/amd64,linux/arm64 `
                --build-arg VERSION=$Version `
                -t "accu-did/$Service`:$Version" `
                -t "accu-did/$Service`:latest" `
                -f $DockerfilePath `
                $ContextPath `
                --provenance=false `
                --sbom=false `
                --output=type=docker

        } else {
            # Single-arch build
            Write-Host "   Building for current platform..." -ForegroundColor Green

            & docker build `
                --build-arg VERSION=$Version `
                -t "accu-did/$Service`:$Version" `
                -t "accu-did/$Service`:latest" `
                -f $DockerfilePath `
                $ContextPath
        }

        if ($LASTEXITCODE -eq 0) {
            Write-Host "   ✅ Built: accu-did/$Service`:$Version" -ForegroundColor Green
            Write-Host "   ✅ Tagged: accu-did/$Service`:latest" -ForegroundColor Green

            # Save manifest info
            $ManifestContent = @"
# Built: $(Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")
# Version: $Version
accu-did/$Service`:$Version
accu-did/$Service`:latest

# Image details:
$((docker images --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" | Select-String "accu-did/$Service") -join "`n")
"@
            $ManifestContent | Out-File -FilePath "dist\images\$Service\manifests.txt" -Encoding UTF8

            Write-Host "   📄 Manifest saved to dist\images\$Service\manifests.txt" -ForegroundColor Green
            return $true
        } else {
            Write-Host "   ❌ Failed to build $Service image" -ForegroundColor Red
            return $false
        }
    }
    catch {
        Write-Host "   ❌ Failed to build $Service image: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Build resolver image
$ResolverSuccess = Build-Image -Service "resolver"

# Build registrar image
$RegistrarSuccess = Build-Image -Service "registrar"

Write-Host ""
if ($ResolverSuccess -and $RegistrarSuccess) {
    Write-Host "🎉 Docker image build complete!" -ForegroundColor Green
} else {
    Write-Host "⚠️  Some images failed to build" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "📊 Local images:" -ForegroundColor Cyan
try {
    docker images --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" | Select-String "accu-did"
} catch {
    Write-Host "   No accu-did images found" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "📁 Image manifests:" -ForegroundColor Cyan
Get-ChildItem -Path "dist\images" -Filter "manifests.txt" -Recurse | ForEach-Object { "  $($_.FullName.Replace($PWD.Path + '\', ''))" }

Write-Host ""
Write-Host "💡 To export images for distribution:" -ForegroundColor Yellow
Write-Host "   docker save accu-did/resolver:$Version | gzip > dist\images\resolver\resolver-$Version.tar.gz"
Write-Host "   docker save accu-did/registrar:$Version | gzip > dist\images\registrar\registrar-$Version.tar.gz"

Write-Host ""
Write-Host "[OK] Images built successfully" -ForegroundColor Green