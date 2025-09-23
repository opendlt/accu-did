#!/usr/bin/env bash
# scripts/sbom.local.sh - Generate Software Bill of Materials locally
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
cd "$ROOT"

echo "🔍 Generating Software Bill of Materials (SBOM)..."

# Create SBOM directory
mkdir -p dist/sbom

# Check if syft is available
if ! command -v syft >/dev/null 2>&1; then
    echo "⚠️  syft not found, skipping SBOM generation"
    echo "   Install with: curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin"
    echo "   Or download from: https://github.com/anchore/syft/releases"

    # Create placeholder file
    echo "# SBOM generation skipped - syft not available" > dist/sbom/README.txt
    echo "# Install syft to generate SBOMs" >> dist/sbom/README.txt
    echo ""
    echo "[SKIP] SBOM generation - syft not available"
    exit 0
fi

echo "📋 syft found: $(syft version --output json 2>/dev/null | grep -o '"version":"[^"]*"' | cut -d'"' -f4 || syft version 2>/dev/null | head -1)"

# Read version
VER=$(cat VERSION)

# Generate SBOM for resolver image
echo ""
echo "📦 Generating SBOM for resolver..."
if docker images --format "{{.Repository}}:{{.Tag}}" | grep -q "accu-did/resolver:latest"; then
    syft packages "accu-did/resolver:latest" -o json > dist/sbom/resolver-latest.syft.json
    syft packages "accu-did/resolver:latest" -o spdx-json > dist/sbom/resolver-latest.spdx.json
    syft packages "accu-did/resolver:latest" -o table > dist/sbom/resolver-latest.txt
    echo "   ✅ SBOM saved: dist/sbom/resolver-latest.{syft.json,spdx.json,txt}"
else
    echo "   ⚠️  Image accu-did/resolver:latest not found, skipping"
fi

# Generate SBOM for registrar image
echo ""
echo "📦 Generating SBOM for registrar..."
if docker images --format "{{.Repository}}:{{.Tag}}" | grep -q "accu-did/registrar:latest"; then
    syft packages "accu-did/registrar:latest" -o json > dist/sbom/registrar-latest.syft.json
    syft packages "accu-did/registrar:latest" -o spdx-json > dist/sbom/registrar-latest.spdx.json
    syft packages "accu-did/registrar:latest" -o table > dist/sbom/registrar-latest.txt
    echo "   ✅ SBOM saved: dist/sbom/registrar-latest.{syft.json,spdx.json,txt}"
else
    echo "   ⚠️  Image accu-did/registrar:latest not found, skipping"
fi

# Generate SBOM for source code
echo ""
echo "📦 Generating SBOM for source code..."
syft packages . -o json > dist/sbom/source-$VER.syft.json
syft packages . -o spdx-json > dist/sbom/source-$VER.spdx.json
syft packages . -o table > dist/sbom/source-$VER.txt
echo "   ✅ Source SBOM saved: dist/sbom/source-$VER.{syft.json,spdx.json,txt}"

echo ""
echo "📊 SBOM Summary:"
find dist/sbom -name "*.json" -exec wc -l {} \; | awk '{print "  " $2 ": " $1 " lines"}'

echo ""
echo "📁 Generated files:"
find dist/sbom -type f | sort | sed 's|^|  |'

echo ""
echo "[OK] SBOM generation complete"