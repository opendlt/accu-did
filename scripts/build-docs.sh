#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
SPEC_DIR="$ROOT/docs/spec/openapi"
SITE_DIR="$ROOT/docs/site"
mkdir -p "$SITE_DIR"

build_one () {
  local y="$1"; local h="$2"
  if command -v docker >/dev/null 2>&1; then
    docker run --rm -v "$ROOT":/work -w /work redocly/redoc build -o "$h" "$y"
  else
    npx --yes redoc-cli@0.15.1 build "$y" -o "$h"
  fi
}

build_one "$SPEC_DIR/resolver.yaml"  "$SITE_DIR/resolver.html"
build_one "$SPEC_DIR/registrar.yaml" "$SITE_DIR/registrar.html"

cp -f "$ROOT/docs/spec/diagrams/"*.mmd "$SITE_DIR" 2>/dev/null || true
cp -f "$ROOT/docs/spec/method.md" "$SITE_DIR" 2>/dev/null || true
cp -f "$ROOT/docs/site/index.template.html" "$ROOT/docs/site/index.html"
echo "✅ Docs built at $SITE_DIR"