#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-npx}"

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
SPEC_DIR="$ROOT/docs/spec/openapi"
SITE_DIR="$ROOT/docs/site"
mkdir -p "$SITE_DIR"

if [ "$MODE" = "npx" ]; then
  echo "Building docs with npx..."
  npx --yes @redocly/cli build-docs "$SPEC_DIR/resolver.yaml" --output "$SITE_DIR/resolver.html"
  npx --yes @redocly/cli build-docs "$SPEC_DIR/registrar.yaml" --output "$SITE_DIR/registrar.html"
elif [ "$MODE" = "docker" ]; then
  echo "Building docs with Docker..."
  docker run --rm -v "$ROOT":/work -w /work redocly/redoc build -o "$SITE_DIR/resolver.html" "$SPEC_DIR/resolver.yaml"
  docker run --rm -v "$ROOT":/work -w /work redocly/redoc build -o "$SITE_DIR/registrar.html" "$SPEC_DIR/registrar.yaml"
else
  echo "Error: Only 'npx' and 'docker' modes are supported"
  exit 1
fi

# Copy supporting files
cp -f "$ROOT/docs/spec/diagrams/"*.mmd "$SITE_DIR" 2>/dev/null || true
cp -f "$ROOT/docs/spec/method.md" "$SITE_DIR" 2>/dev/null || true
cp -f "$ROOT/docs/site/index.template.html" "$ROOT/docs/site/index.html" 2>/dev/null || true

echo "[OK] Docs built at $SITE_DIR"