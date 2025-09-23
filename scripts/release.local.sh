#!/usr/bin/env bash
# scripts/release.local.sh - Create local release with version tagging
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
cd "$ROOT"

echo "🚀 Creating local release..."

# Read version from VERSION file
if [ ! -f "VERSION" ]; then
    echo "ERROR: VERSION file not found"
    exit 1
fi

VER=$(cat VERSION)
echo "   Version: $VER"

# Verify working tree is clean
echo ""
echo "🔍 Checking git status..."
if [ -n "$(git status --porcelain)" ]; then
    echo "❌ Working tree is not clean. Commit or stash changes first."
    echo ""
    git status --short
    exit 1
fi

echo "   ✅ Working tree is clean"

# Check if version is consistent across OpenAPI specs
echo ""
echo "📝 Verifying version consistency..."

# Check resolver OpenAPI spec
RESOLVER_SPEC="docs/spec/openapi/resolver.yaml"
if [ -f "$RESOLVER_SPEC" ]; then
    RESOLVER_VER=$(grep -E "^\s*version:" "$RESOLVER_SPEC" | awk '{print $2}' | tr -d '"' || echo "not found")
    if [ "$RESOLVER_VER" != "$VER" ]; then
        echo "⚠️  Version mismatch in $RESOLVER_SPEC: found '$RESOLVER_VER', expected '$VER'"
        echo "   Update the version field in $RESOLVER_SPEC"
        # Don't exit, just warn
    else
        echo "   ✅ Resolver OpenAPI version matches: $VER"
    fi
else
    echo "   ⚠️  Resolver OpenAPI spec not found: $RESOLVER_SPEC"
fi

# Check registrar OpenAPI spec
REGISTRAR_SPEC="docs/spec/openapi/registrar.yaml"
if [ -f "$REGISTRAR_SPEC" ]; then
    REGISTRAR_VER=$(grep -E "^\s*version:" "$REGISTRAR_SPEC" | awk '{print $2}' | tr -d '"' || echo "not found")
    if [ "$REGISTRAR_VER" != "$VER" ]; then
        echo "⚠️  Version mismatch in $REGISTRAR_SPEC: found '$REGISTRAR_VER', expected '$VER'"
        echo "   Update the version field in $REGISTRAR_SPEC"
        # Don't exit, just warn
    else
        echo "   ✅ Registrar OpenAPI version matches: $VER"
    fi
else
    echo "   ⚠️  Registrar OpenAPI spec not found: $REGISTRAR_SPEC"
fi

# Check if tag already exists
TAG_NAME="v$VER"
if git tag -l | grep -q "^$TAG_NAME$"; then
    echo ""
    echo "⚠️  Tag $TAG_NAME already exists"
    echo "   Existing tags:"
    git tag -l | grep "^v" | tail -5 | sed 's|^|     |'
    echo ""
    echo "   To create a new release:"
    echo "   1. Update VERSION file with new version"
    echo "   2. Commit the version change"
    echo "   3. Run this script again"
    echo ""
    echo "   To delete existing tag (if needed):"
    echo "   git tag -d $TAG_NAME"
    exit 1
fi

# Get current commit hash
COMMIT_HASH=$(git rev-parse --short HEAD)
echo "   📍 Current commit: $COMMIT_HASH"

# Get last tag for release notes
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -n "$LAST_TAG" ]; then
    echo "   📜 Last tag: $LAST_TAG"
    COMMIT_COUNT=$(git rev-list --count "$LAST_TAG..HEAD" 2>/dev/null || echo "0")
    echo "   📊 Commits since last tag: $COMMIT_COUNT"
fi

# Create the tag
echo ""
echo "🏷️  Creating git tag..."

# Create annotated tag with release information
TAG_MESSAGE=$(cat <<EOF
accu-did v$VER

Release Information:
- Version: $VER
- Commit: $COMMIT_HASH
- Date: $(date -u +"%Y-%m-%d %H:%M:%S UTC")
$([ -n "$LAST_TAG" ] && echo "- Commits since $LAST_TAG: $COMMIT_COUNT")

Distribution:
- Binaries: dist/bin/
- Docker Images: accu-did/resolver:$VER, accu-did/registrar:$VER
- Documentation: dist/docs/docs-$VER.zip

To push this release:
  git push origin main
  git push origin $TAG_NAME

EOF
)

git tag -a "$TAG_NAME" -m "$TAG_MESSAGE"

if [ $? -eq 0 ]; then
    echo "   ✅ Tag created: $TAG_NAME"
else
    echo "   ❌ Failed to create tag"
    exit 1
fi

# Show tag information
echo ""
echo "📋 Tag information:"
git show --stat "$TAG_NAME" | head -20

echo ""
echo "🎉 Local release complete!"
echo ""
echo "📦 Release artifacts:"
echo "   Tag: $TAG_NAME"
echo "   Binaries: dist/bin/ (if built)"
echo "   Images: accu-did/resolver:$VER, accu-did/registrar:$VER (if built)"
echo "   Docs: dist/docs/docs-$VER.zip (if built)"
echo ""
echo "📤 Next steps (optional):"
echo "   git push origin main           # Push commits"
echo "   git push origin $TAG_NAME      # Push tag"
echo ""
echo "💡 To build all release artifacts:"
echo "   make release-local"
echo ""
echo "[OK] Tag created locally. Push is optional; not required."