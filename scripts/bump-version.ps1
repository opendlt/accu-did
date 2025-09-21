#!/usr/bin/env pwsh
Param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("patch", "minor", "major")]
    [string]$Bump
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path | Split-Path -Parent
$VersionFile = Join-Path $Root "VERSION"

# Read current version
if (!(Test-Path $VersionFile)) {
    Write-Error "VERSION file not found"
    exit 1
}

$CurrentVersion = (Get-Content $VersionFile).Trim()
if ($CurrentVersion -notmatch '^(\d+)\.(\d+)\.(\d+)$') {
    Write-Error "Invalid version format: $CurrentVersion"
    exit 1
}

$Major = [int]$Matches[1]
$Minor = [int]$Matches[2]
$Patch = [int]$Matches[3]

# Bump version
switch ($Bump) {
    "major" {
        $Major++
        $Minor = 0
        $Patch = 0
    }
    "minor" {
        $Minor++
        $Patch = 0
    }
    "patch" {
        $Patch++
    }
}

$NewVersion = "$Major.$Minor.$Patch"

# Write new version
Set-Content -Path $VersionFile -Value $NewVersion -NoNewline

# Output new version
Write-Host $NewVersion