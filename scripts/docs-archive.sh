#!/usr/bin/env bash
# scripts/docs-archive.sh - Package documentation for distribution
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
cd "$ROOT"

echo "📚 Packaging documentation for distribution..."

# Read version from VERSION file
if [ ! -f "VERSION" ]; then
    echo "ERROR: VERSION file not found"
    exit 1
fi

VER=$(cat VERSION)
echo "   Version: $VER"

# Create docs dist directory
mkdir -p dist/docs

# Build docs first
echo ""
echo "🔨 Building documentation..."
if [ -f "scripts/build-docs.sh" ]; then
    chmod +x scripts/build-docs.sh
    ./scripts/build-docs.sh npx
else
    echo "   ⚠️  scripts/build-docs.sh not found, building manually..."

    # Create docs/site if it doesn't exist
    mkdir -p docs/site

    # Try to build with available tools
    if command -v npx >/dev/null 2>&1; then
        echo "   Using npx to build docs..."
        npx redoc-cli build docs/spec/openapi/resolver.yaml --output docs/site/resolver.html || echo "   Warning: resolver docs build failed"
        npx redoc-cli build docs/spec/openapi/registrar.yaml --output docs/site/registrar.html || echo "   Warning: registrar docs build failed"
    else
        echo "   No doc builder available, copying markdown files..."
        cp docs/spec/method.md docs/site/ 2>/dev/null || true
        cp docs/resolver.md docs/site/ 2>/dev/null || true
        cp docs/registrar.md docs/site/ 2>/dev/null || true
    fi
fi

# Verify docs/site exists
if [ ! -d "docs/site" ]; then
    echo "   Creating minimal docs/site directory..."
    mkdir -p docs/site
    echo "# Accumulate DID Documentation" > docs/site/index.md
    echo "See docs/ directory for full documentation." >> docs/site/index.md
fi

# Create archive filename
ARCHIVE_NAME="docs-$VER.zip"
ARCHIVE_PATH="dist/docs/$ARCHIVE_NAME"

echo ""
echo "📦 Creating documentation archive..."

# Check for zip command
if command -v zip >/dev/null 2>&1; then
    echo "   Using zip to create archive..."

    # Create zip archive including:
    # - Built documentation from docs/site/
    # - OpenAPI specs
    # - Method specification
    # - README and other markdown docs

    (cd docs && zip -r "../$ARCHIVE_PATH" \
        site/ \
        spec/ \
        *.md \
        2>/dev/null) || true

    # Add top-level docs
    zip -u "$ARCHIVE_PATH" \
        README.md \
        CHANGELOG.md \
        2>/dev/null || true

elif command -v 7z >/dev/null 2>&1; then
    echo "   Using 7z to create archive..."
    7z a "$ARCHIVE_PATH" \
        docs/site/ \
        docs/spec/ \
        docs/*.md \
        README.md \
        CHANGELOG.md \
        >/dev/null 2>&1 || true

else
    echo "   Neither zip nor 7z found, using tar..."
    tar -czf "${ARCHIVE_PATH%.zip}.tar.gz" \
        docs/site \
        docs/spec \
        docs/*.md \
        README.md \
        CHANGELOG.md \
        2>/dev/null || true
    ARCHIVE_PATH="${ARCHIVE_PATH%.zip}.tar.gz"
    ARCHIVE_NAME="${ARCHIVE_NAME%.zip}.tar.gz"
fi

# Verify archive was created
if [ -f "$ARCHIVE_PATH" ]; then
    ARCHIVE_SIZE=$(du -h "$ARCHIVE_PATH" | cut -f1)
    echo "   ✅ Archive created: $ARCHIVE_NAME ($ARCHIVE_SIZE)"

    # List archive contents
    echo ""
    echo "📋 Archive contents:"
    if [[ "$ARCHIVE_PATH" == *.zip ]]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -l "$ARCHIVE_PATH" 2>/dev/null | tail -n +4 | head -n -2 | awk '{print "   " $4}' || true
        fi
    elif [[ "$ARCHIVE_PATH" == *.tar.gz ]]; then
        tar -tzf "$ARCHIVE_PATH" 2>/dev/null | sed 's|^|   |' || true
    fi

    # Create checksum
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$ARCHIVE_PATH" > "$ARCHIVE_PATH.sha256"
        echo "   🔐 Checksum: $(basename "$ARCHIVE_PATH").sha256"
    fi

else
    echo "   ❌ Failed to create documentation archive"
    exit 1
fi

echo ""
echo "📁 Documentation distribution files:"
find dist/docs -type f | sort | sed 's|^|  |'

echo ""
echo "[OK] Documentation packaging complete"