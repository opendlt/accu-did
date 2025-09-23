#!/usr/bin/env bash
# scripts/sdk-openapi-merge.sh - Merge OpenAPI specs for SDK generation
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
cd "$ROOT"

echo "🔧 Merging OpenAPI specifications for SDK..."

# Read version from VERSION file
if [ ! -f "VERSION" ]; then
    echo "ERROR: VERSION file not found"
    exit 1
fi

VER=$(cat VERSION)
echo "   Version: $VER"

# Ensure tools directory exists and is built
echo ""
echo "📦 Building sdkmerge tool..."
cd tools/sdkmerge
if [ ! -f "go.mod" ]; then
    go mod init github.com/opendlt/accu-did/tools/sdkmerge
    go mod tidy
fi

# Download dependencies if needed
go mod download
go mod tidy

# Build the merge tool
go build -o sdkmerge .

cd "$ROOT"

# Check input files exist
RESOLVER_SPEC="docs/spec/openapi/resolver.yaml"
REGISTRAR_SPEC="docs/spec/openapi/registrar.yaml"
OUTPUT_SPEC="sdks/spec/openapi/accdid-sdk.yaml"

if [ ! -f "$RESOLVER_SPEC" ]; then
    echo "   ⚠️  Resolver OpenAPI spec not found: $RESOLVER_SPEC"
    echo "   Creating placeholder..."
    mkdir -p "$(dirname "$RESOLVER_SPEC")"
    cat > "$RESOLVER_SPEC" <<'EOF'
openapi: 3.0.3
info:
  title: Accumulate DID Resolver
  version: 0.1.0
paths:
  /resolve:
    get:
      summary: Resolve a DID
      parameters:
        - name: did
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: DID resolution result
EOF
fi

if [ ! -f "$REGISTRAR_SPEC" ]; then
    echo "   ⚠️  Registrar OpenAPI spec not found: $REGISTRAR_SPEC"
    echo "   Creating placeholder..."
    mkdir -p "$(dirname "$REGISTRAR_SPEC")"
    cat > "$REGISTRAR_SPEC" <<'EOF'
openapi: 3.0.3
info:
  title: Accumulate DID Registrar
  version: 0.1.0
paths:
  /register:
    post:
      summary: Register a new DID
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        '200':
          description: Registration successful
EOF
fi

# Merge the specifications
echo ""
echo "🔄 Merging OpenAPI specifications..."
./tools/sdkmerge/sdkmerge "$RESOLVER_SPEC" "$REGISTRAR_SPEC" "$OUTPUT_SPEC" "$VER"

# Verify output
if [ -f "$OUTPUT_SPEC" ]; then
    echo ""
    echo "📄 Generated SDK specification:"
    echo "   File: $OUTPUT_SPEC"
    echo "   Size: $(wc -l < "$OUTPUT_SPEC") lines"

    # Show basic info
    if command -v grep >/dev/null 2>&1; then
        echo "   Title: $(grep -E '^\s*title:' "$OUTPUT_SPEC" | head -1 | sed 's/.*title: *//')"
        echo "   Version: $(grep -E '^\s*version:' "$OUTPUT_SPEC" | head -1 | sed 's/.*version: *//')"
        echo "   Paths: $(grep -E '^\s*/.*:$' "$OUTPUT_SPEC" | wc -l || echo "unknown")"
    fi
else
    echo "   ❌ Failed to generate merged specification"
    exit 1
fi

echo ""
echo "[OK] OpenAPI merge complete"