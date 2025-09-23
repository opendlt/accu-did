#!/usr/bin/env bash
# scripts/build-binaries.sh - Cross-compile binaries for multiple platforms
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
cd "$ROOT"

# Read version from VERSION file
if [ ! -f "VERSION" ]; then
    echo "ERROR: VERSION file not found"
    exit 1
fi

VER=$(cat VERSION)
echo "🔨 Building binaries for version: $VER"

# Export build settings
export CGO_ENABLED=0
export GO111MODULE=on

# Define platforms to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/arm64"
)

# Create dist directories
mkdir -p dist/bin

# Build function
build_binary() {
    local platform=$1
    local service=$2
    local source_dir=$3

    IFS='/' read -r os arch <<< "$platform"

    echo "  Building $service for $os/$arch..."

    # Create platform directory
    local output_dir="dist/bin/${os}-${arch}"
    mkdir -p "$output_dir"

    # Set platform-specific variables
    export GOOS="$os"
    export GOARCH="$arch"

    # Determine binary name
    local binary_name="$service"
    if [ "$os" = "windows" ]; then
        binary_name="${service}.exe"
    fi

    # Build with version embedded
    local ldflags="-X main.version=$VER -s -w"

    # Build the binary
    if ! go build -ldflags "$ldflags" -o "$output_dir/$binary_name" "./$source_dir/cmd/$service"; then
        echo "    ❌ Failed to build $service for $os/$arch"
        return 1
    fi

    echo "    ✅ Built: $output_dir/$binary_name"
}

# Build all binaries
echo ""
echo "Building binaries..."

for platform in "${PLATFORMS[@]}"; do
    # Skip darwin if not supported (some CI environments)
    if [[ "$platform" == darwin/* ]] && ! command -v go >/dev/null 2>&1; then
        echo "  Skipping $platform (go not available or unsupported)"
        continue
    fi

    echo ""
    echo "📦 Platform: $platform"

    # Build resolver
    build_binary "$platform" "resolver" "resolver-go"

    # Build registrar
    build_binary "$platform" "registrar" "registrar-go"
done

echo ""
echo "🎉 Binary compilation complete!"
echo ""
echo "📁 Distribution tree:"
find dist/bin -type f | sort | sed 's|^|  |'

echo ""
echo "📊 Binary sizes:"
find dist/bin -type f -exec ls -lh {} \; | awk '{print $5 "\t" $9}' | sort

# Create checksums
echo ""
echo "🔐 Generating checksums..."
find dist/bin -type f -exec sha256sum {} \; > dist/bin/checksums.sha256
echo "  ✅ Checksums saved to dist/bin/checksums.sha256"

echo ""
echo "[OK] Build complete - binaries ready in dist/bin/"