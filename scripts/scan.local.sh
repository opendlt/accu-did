#!/usr/bin/env bash
# scripts/scan.local.sh - Run vulnerability scans locally
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
cd "$ROOT"

echo "🔐 Running local vulnerability scans..."

# Create scan directory
mkdir -p dist/scan

# Check if trivy is available
if ! command -v trivy >/dev/null 2>&1; then
    echo "⚠️  trivy not found, skipping vulnerability scanning"
    echo "   Install with: curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin"
    echo "   Or download from: https://github.com/aquasecurity/trivy/releases"

    # Create placeholder file
    echo "# Vulnerability scanning skipped - trivy not available" > dist/scan/README.txt
    echo "# Install trivy to run vulnerability scans" >> dist/scan/README.txt
    echo ""
    echo "[SKIP] Vulnerability scanning - trivy not available"
    exit 0
fi

echo "🔍 trivy found: $(trivy version --format json 2>/dev/null | grep -o '"Version":"[^"]*"' | cut -d'"' -f4 || trivy version 2>/dev/null | head -1)"

# Update trivy database
echo ""
echo "📊 Updating vulnerability database..."
trivy image --download-db-only >/dev/null 2>&1 || echo "   Database update may have failed, continuing..."

# Scan resolver image
echo ""
echo "🔍 Scanning resolver image..."
if docker images --format "{{.Repository}}:{{.Tag}}" | grep -q "accu-did/resolver:latest"; then
    # Full scan with soft-fail (exit-code 0)
    trivy image \
        --exit-code 0 \
        --format json \
        --output dist/scan/resolver-latest.trivy.json \
        accu-did/resolver:latest

    # Summary scan
    trivy image \
        --exit-code 0 \
        --format table \
        --output dist/scan/resolver-latest.txt \
        accu-did/resolver:latest

    # Critical vulnerabilities only (with hard-fail)
    if trivy image \
        --exit-code 1 \
        --severity CRITICAL \
        --format table \
        --output dist/scan/resolver-critical.txt \
        accu-did/resolver:latest; then
        echo "   ✅ No critical vulnerabilities found in resolver"
    else
        echo "   ⚠️  Critical vulnerabilities found in resolver - see dist/scan/resolver-critical.txt"
    fi

    echo "   📄 Resolver scan saved: dist/scan/resolver-latest.{trivy.json,txt}"
else
    echo "   ⚠️  Image accu-did/resolver:latest not found, skipping"
fi

# Scan registrar image
echo ""
echo "🔍 Scanning registrar image..."
if docker images --format "{{.Repository}}:{{.Tag}}" | grep -q "accu-did/registrar:latest"; then
    # Full scan with soft-fail (exit-code 0)
    trivy image \
        --exit-code 0 \
        --format json \
        --output dist/scan/registrar-latest.trivy.json \
        accu-did/registrar:latest

    # Summary scan
    trivy image \
        --exit-code 0 \
        --format table \
        --output dist/scan/registrar-latest.txt \
        accu-did/registrar:latest

    # Critical vulnerabilities only (with hard-fail)
    if trivy image \
        --exit-code 1 \
        --severity CRITICAL \
        --format table \
        --output dist/scan/registrar-critical.txt \
        accu-did/registrar:latest; then
        echo "   ✅ No critical vulnerabilities found in registrar"
    else
        echo "   ⚠️  Critical vulnerabilities found in registrar - see dist/scan/registrar-critical.txt"
    fi

    echo "   📄 Registrar scan saved: dist/scan/registrar-latest.{trivy.json,txt}"
else
    echo "   ⚠️  Image accu-did/registrar:latest not found, skipping"
fi

# Scan filesystem (source code and dependencies)
echo ""
echo "🔍 Scanning filesystem for vulnerabilities..."
trivy fs \
    --exit-code 0 \
    --format json \
    --output dist/scan/filesystem.trivy.json \
    .

trivy fs \
    --exit-code 0 \
    --format table \
    --output dist/scan/filesystem.txt \
    .

echo "   📄 Filesystem scan saved: dist/scan/filesystem.{trivy.json,txt}"

echo ""
echo "📊 Scan Summary:"
find dist/scan -name "*.json" -exec wc -l {} \; | awk '{print "  " $2 ": " $1 " lines"}'

echo ""
echo "📁 Generated files:"
find dist/scan -type f | sort | sed 's|^|  |'

echo ""
echo "💡 To fail on critical vulnerabilities, set exit-code 1 in scan commands"
echo ""
echo "[OK] Vulnerability scanning complete"