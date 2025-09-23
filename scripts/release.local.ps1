# scripts/release.local.ps1 - Create local release with version tagging
param(
    [switch]$Force = $false
)

$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

Write-Host "🚀 Creating local release..." -ForegroundColor Cyan

# Read version from VERSION file
if (-not (Test-Path "VERSION")) {
    Write-Error "VERSION file not found"
    exit 1
}

$Version = Get-Content "VERSION" -Raw | ForEach-Object { $_.Trim() }
Write-Host "   Version: $Version" -ForegroundColor Green

# Verify working tree is clean
Write-Host ""
Write-Host "🔍 Checking git status..." -ForegroundColor Cyan

$GitStatus = git status --porcelain
if ($GitStatus) {
    Write-Host "❌ Working tree is not clean. Commit or stash changes first." -ForegroundColor Red
    Write-Host ""
    git status --short
    exit 1
}

Write-Host "   ✅ Working tree is clean" -ForegroundColor Green

# Check if version is consistent across OpenAPI specs
Write-Host ""
Write-Host "📝 Verifying version consistency..." -ForegroundColor Cyan

# Check resolver OpenAPI spec
$ResolverSpec = "docs\spec\openapi\resolver.yaml"
if (Test-Path $ResolverSpec) {
    $ResolverContent = Get-Content $ResolverSpec -Raw
    if ($ResolverContent -match "version:\s*[`"']?([^`"'\s]+)[`"']?") {
        $ResolverVer = $Matches[1]
        if ($ResolverVer -ne $Version) {
            Write-Host "⚠️  Version mismatch in $ResolverSpec`: found '$ResolverVer', expected '$Version'" -ForegroundColor Yellow
            Write-Host "   Update the version field in $ResolverSpec"
        } else {
            Write-Host "   ✅ Resolver OpenAPI version matches: $Version" -ForegroundColor Green
        }
    } else {
        Write-Host "   ⚠️  Could not find version in $ResolverSpec" -ForegroundColor Yellow
    }
} else {
    Write-Host "   ⚠️  Resolver OpenAPI spec not found: $ResolverSpec" -ForegroundColor Yellow
}

# Check registrar OpenAPI spec
$RegistrarSpec = "docs\spec\openapi\registrar.yaml"
if (Test-Path $RegistrarSpec) {
    $RegistrarContent = Get-Content $RegistrarSpec -Raw
    if ($RegistrarContent -match "version:\s*[`"']?([^`"'\s]+)[`"']?") {
        $RegistrarVer = $Matches[1]
        if ($RegistrarVer -ne $Version) {
            Write-Host "⚠️  Version mismatch in $RegistrarSpec`: found '$RegistrarVer', expected '$Version'" -ForegroundColor Yellow
            Write-Host "   Update the version field in $RegistrarSpec"
        } else {
            Write-Host "   ✅ Registrar OpenAPI version matches: $Version" -ForegroundColor Green
        }
    } else {
        Write-Host "   ⚠️  Could not find version in $RegistrarSpec" -ForegroundColor Yellow
    }
} else {
    Write-Host "   ⚠️  Registrar OpenAPI spec not found: $RegistrarSpec" -ForegroundColor Yellow
}

# Check if tag already exists
$TagName = "v$Version"
$ExistingTag = git tag -l | Where-Object { $_ -eq $TagName }

if ($ExistingTag -and -not $Force) {
    Write-Host ""
    Write-Host "⚠️  Tag $TagName already exists" -ForegroundColor Yellow
    Write-Host "   Existing tags:" -ForegroundColor Yellow
    git tag -l | Where-Object { $_ -match "^v" } | Select-Object -Last 5 | ForEach-Object { "     $_" }
    Write-Host ""
    Write-Host "   To create a new release:" -ForegroundColor Yellow
    Write-Host "   1. Update VERSION file with new version"
    Write-Host "   2. Commit the version change"
    Write-Host "   3. Run this script again"
    Write-Host ""
    Write-Host "   To delete existing tag (if needed):" -ForegroundColor Yellow
    Write-Host "   git tag -d $TagName"
    exit 1
}

# Get current commit hash
$CommitHash = git rev-parse --short HEAD
Write-Host "   📍 Current commit: $CommitHash" -ForegroundColor Green

# Get last tag for release notes
try {
    $LastTag = git describe --tags --abbrev=0 2>$null
    if ($LastTag) {
        Write-Host "   📜 Last tag: $LastTag" -ForegroundColor Green
        $CommitCount = git rev-list --count "$LastTag..HEAD" 2>$null
        Write-Host "   📊 Commits since last tag: $CommitCount" -ForegroundColor Green
    }
} catch {
    $LastTag = ""
    $CommitCount = 0
}

# Create the tag
Write-Host ""
Write-Host "🏷️  Creating git tag..." -ForegroundColor Cyan

# Create tag message
$TagMessage = @"
accu-did v$Version

Release Information:
- Version: $Version
- Commit: $CommitHash
- Date: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss UTC")
$(if ($LastTag) { "- Commits since $LastTag`: $CommitCount" })

Distribution:
- Binaries: dist/bin/
- Docker Images: accu-did/resolver:$Version, accu-did/registrar:$Version
- Documentation: dist/docs/docs-$Version.zip

To push this release:
  git push origin main
  git push origin $TagName
"@

try {
    git tag -a $TagName -m $TagMessage
    Write-Host "   ✅ Tag created: $TagName" -ForegroundColor Green
}
catch {
    Write-Host "   ❌ Failed to create tag: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Show tag information
Write-Host ""
Write-Host "📋 Tag information:" -ForegroundColor Cyan
git show --stat $TagName | Select-Object -First 20

Write-Host ""
Write-Host "🎉 Local release complete!" -ForegroundColor Green
Write-Host ""
Write-Host "📦 Release artifacts:" -ForegroundColor Cyan
Write-Host "   Tag: $TagName"
Write-Host "   Binaries: dist\bin\ (if built)"
Write-Host "   Images: accu-did/resolver:$Version, accu-did/registrar:$Version (if built)"
Write-Host "   Docs: dist\docs\docs-$Version.zip (if built)"
Write-Host ""
Write-Host "📤 Next steps (optional):" -ForegroundColor Yellow
Write-Host "   git push origin main           # Push commits"
Write-Host "   git push origin $TagName      # Push tag"
Write-Host ""
Write-Host "💡 To build all release artifacts:" -ForegroundColor Yellow
Write-Host "   make release-local"
Write-Host ""
Write-Host "[OK] Tag created locally. Push is optional; not required." -ForegroundColor Green