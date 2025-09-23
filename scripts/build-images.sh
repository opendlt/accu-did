#!/usr/bin/env bash
# scripts/build-images.sh - Build multi-arch Docker images locally
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
cd "$ROOT"

# Read version from VERSION file
if [ ! -f "VERSION" ]; then
    echo "ERROR: VERSION file not found"
    exit 1
fi

VER=$(cat VERSION)
echo "🐳 Building Docker images for version: $VER"

# Ensure docker buildx is available
if ! docker buildx version >/dev/null 2>&1; then
    echo "❌ docker buildx is required but not available"
    echo "   Install Docker Desktop or docker buildx plugin"
    exit 1
fi

# Create/use buildx builder
echo "🔧 Setting up buildx builder..."
docker buildx create --use --name accu-did-builder >/dev/null 2>&1 || {
    echo "   Using existing builder or default"
    docker buildx use accu-did-builder 2>/dev/null || docker buildx use default
}

# Create images directory
mkdir -p dist/images/resolver dist/images/registrar

# Build function
build_image() {
    local service=$1
    local dockerfile_path="${service}-go/Dockerfile"
    local context_path="${service}-go"

    echo ""
    echo "📦 Building $service image..."

    # Verify Dockerfile exists
    if [ ! -f "$dockerfile_path" ]; then
        echo "❌ Dockerfile not found: $dockerfile_path"
        return 1
    fi

    # Build multi-arch image
    echo "   Building for linux/amd64,linux/arm64..."

    # Build and load into local Docker (amd64 only for local use)
    docker buildx build \
        --platform linux/amd64,linux/arm64 \
        --build-arg VERSION="$VER" \
        -t "accu-did/$service:$VER" \
        -t "accu-did/$service:latest" \
        -f "$dockerfile_path" \
        "$context_path" \
        --provenance=false \
        --sbom=false \
        --output=type=docker

    if [ $? -eq 0 ]; then
        echo "   ✅ Built: accu-did/$service:$VER"
        echo "   ✅ Tagged: accu-did/$service:latest"

        # Save manifest info
        {
            echo "# Built: $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
            echo "# Version: $VER"
            echo "accu-did/$service:$VER"
            echo "accu-did/$service:latest"
            echo ""
            echo "# Image details:"
            docker images --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" \
                | grep "accu-did/$service" || true
        } > "dist/images/$service/manifests.txt"

        echo "   📄 Manifest saved to dist/images/$service/manifests.txt"
    else
        echo "   ❌ Failed to build $service image"
        return 1
    fi
}

# Build resolver image
build_image "resolver"

# Build registrar image
build_image "registrar"

echo ""
echo "🎉 Docker image build complete!"
echo ""
echo "📊 Local images:"
docker images --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" | grep accu-did || echo "   No accu-did images found"

echo ""
echo "📁 Image manifests:"
find dist/images -name "manifests.txt" | sed 's|^|  |'

echo ""
echo "💡 To export images for distribution:"
echo "   docker save accu-did/resolver:$VER | gzip > dist/images/resolver/resolver-$VER.tar.gz"
echo "   docker save accu-did/registrar:$VER | gzip > dist/images/registrar/registrar-$VER.tar.gz"

echo ""
echo "[OK] Images built successfully"