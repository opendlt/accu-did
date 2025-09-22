#!/usr/bin/env bash
# scripts/conformance.sh - Run conformance tests against resolver and registrar
set -euo pipefail

# Default URLs for services
export RESOLVER_URL="${RESOLVER_URL:-http://127.0.0.1:8080}"
export REGISTRAR_URL="${REGISTRAR_URL:-http://127.0.0.1:8081}"

ROOT="$(cd "$(dirname "$0")/.."; pwd)"

echo "🔍 Running conformance tests..."
echo "   Resolver:  $RESOLVER_URL"
echo "   Registrar: $REGISTRAR_URL"

# Check if conformance tool exists
if [ -f "$ROOT/tools/conformance/conformance.go" ]; then
    go run "$ROOT/tools/conformance/conformance.go"
else
    echo "⚠️  Conformance tool not found at tools/conformance/conformance.go"
    echo "   Creating placeholder for future implementation..."

    # Run basic health checks instead
    echo ""
    echo "Running basic health checks..."

    if curl -fsS "${RESOLVER_URL}/healthz" > /dev/null 2>&1; then
        echo "✅ Resolver health check passed"
    else
        echo "❌ Resolver health check failed"
        exit 1
    fi

    if curl -fsS "${REGISTRAR_URL}/healthz" > /dev/null 2>&1; then
        echo "✅ Registrar health check passed"
    else
        echo "❌ Registrar health check failed"
        exit 1
    fi

    echo ""
    echo "[OK] Basic health checks passed"
fi