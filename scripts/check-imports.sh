#!/usr/bin/env bash
# scripts/check-imports.sh - Check for forbidden imports
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"

echo "🔍 Checking for forbidden imports..."

# Check for accumulate/internal imports in our services
FORBIDDEN_FOUND=false

if grep -r "accumulate/internal/" "$ROOT/registrar-go" "$ROOT/resolver-go" 2>/dev/null | grep -v "^Binary file"; then
    echo "❌ Forbidden import found: accumulate/internal/"
    echo "   These are internal packages and should not be imported"
    FORBIDDEN_FOUND=true
fi

# Check for other potentially problematic imports
if grep -r "gitlab.com/accumulatenetwork/accumulate/internal" "$ROOT/registrar-go" "$ROOT/resolver-go" 2>/dev/null | grep -v "^Binary file"; then
    echo "❌ Forbidden import found: gitlab accumulate/internal"
    echo "   Use public API packages instead"
    FORBIDDEN_FOUND=true
fi

# Check for replace directives pointing to local paths (except in go.work)
for module in "$ROOT/registrar-go" "$ROOT/resolver-go"; do
    if [ -f "$module/go.mod" ]; then
        if grep "^replace.*=>\s*\.\." "$module/go.mod" 2>/dev/null; then
            echo "⚠️  Warning: Local replace directive found in $module/go.mod"
            echo "   This may cause issues in production builds"
        fi
    fi
done

if [ "$FORBIDDEN_FOUND" = true ]; then
    echo ""
    echo "❌ Import guard check FAILED"
    exit 1
else
    echo "[OK] Import guard passed - no forbidden imports found"
    exit 0
fi