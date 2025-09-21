# scripts/build-docs.ps1
param(
  [ValidateSet("auto","docker","npx")]
  [string]$Mode = "auto"
)

$ErrorActionPreference = "Stop"

$Root    = Split-Path -Parent $PSScriptRoot | Split-Path -Parent
$SpecDir = Join-Path $Root "docs/spec/openapi"
$SiteDir = Join-Path $Root "docs/site"
New-Item -ItemType Directory -Force -Path $SiteDir | Out-Null

$ResolverYaml  = Join-Path $SpecDir "resolver.yaml"
$RegistrarYaml = Join-Path $SpecDir "registrar.yaml"
$ResolverHtml  = Join-Path $SiteDir "resolver.html"
$RegistrarHtml = Join-Path $SiteDir "registrar.html"

function Build-One {
  param([string]$yaml,[string]$html)
  if ($Mode -eq "docker" -or ($Mode -eq "auto" -and (Get-Command docker -ErrorAction SilentlyContinue))) {
    # Use ${Root} to avoid PowerShell drive parsing on ":" in C:\...
    docker run --rm -v "${Root}:/work" -w /work redocly/redoc build -o "$html" "$yaml"
  } else {
    npx --yes redoc-cli@0.15.1 build "$yaml" -o "$html"
  }
}

Build-One $ResolverYaml  $ResolverHtml
Build-One $RegistrarYaml $RegistrarHtml

# Copy diagrams and method spec if present
Copy-Item "$Root\docs\spec\diagrams\*.mmd" $SiteDir -ErrorAction SilentlyContinue
Copy-Item "$Root\docs\spec\method.md"      $SiteDir -ErrorAction SilentlyContinue

# Create index.html from template
$tpl = Get-Content -Raw -Path "$Root\docs\site\index.template.html"
Set-Content -Path "$Root\docs\site\index.html" -Value $tpl

Write-Host "[OK] Docs built at $SiteDir"