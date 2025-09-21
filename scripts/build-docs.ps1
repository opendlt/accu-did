Param(
  [ValidateSet("auto","docker","npx")]
  [string]$Mode = "auto"
)

$Root = Split-Path -Parent $MyInvocation.MyCommand.Path | Split-Path -Parent
$SpecDir = Join-Path $Root "docs/spec/openapi"
$SiteDir = Join-Path $Root "docs/site"
New-Item -ItemType Directory -Force -Path $SiteDir | Out-Null

$ResolverYaml  = Join-Path $SpecDir "resolver.yaml"
$RegistrarYaml = Join-Path $SpecDir "registrar.yaml"
$ResolverHtml  = Join-Path $SiteDir "resolver.html"
$RegistrarHtml = Join-Path $SiteDir "registrar.html"

function Build-One($yaml,$html) {
  if ($Mode -eq "docker" -or ($Mode -eq "auto" -and (Get-Command docker -ErrorAction SilentlyContinue))) {
    docker run --rm -v "$Root:/work" -w /work redocly/redoc build -o "$html" "$yaml"
  } else {
    npx --yes redoc-cli@0.15.1 build "$yaml" -o "$html"
  }
}

Build-One $ResolverYaml $ResolverHtml
Build-One $RegistrarYaml $RegistrarHtml

# copy diagrams and method.md
Copy-Item "$Root/docs/spec/diagrams/*.mmd" $SiteDir -ErrorAction SilentlyContinue
Copy-Item "$Root/docs/spec/method.md" $SiteDir -ErrorAction SilentlyContinue

# generate index.html from template
Copy-Item "$Root/docs/site/index.template.html" "$Root/docs/site/index.html" -Force
Write-Host "✅ Docs built at $SiteDir"