#!/usr/bin/env bash
set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m'

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
VERSION_FILE="$ROOT/VERSION"

echo -e "${CYAN}🚀 Starting release process...${NC}"

# Read version
if [ ! -f "$VERSION_FILE" ]; then
    echo -e "${RED}ERROR: VERSION file not found${NC}" >&2
    exit 1
fi

VERSION=$(cat "$VERSION_FILE")
TAG="v$VERSION"

echo -e "${YELLOW}📦 Releasing version: $TAG${NC}"

# Check for uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${YELLOW}⚠️  Uncommitted changes detected. Committing VERSION and CHANGELOG...${NC}"
    git add VERSION CHANGELOG.md
    git commit -m "chore(release): prepare $TAG" 2>/dev/null || echo "Nothing to commit"
fi

# Build documentation
echo -e "\n${YELLOW}📚 Building documentation...${NC}"
if ! bash "$ROOT/scripts/build-docs.sh"; then
    echo -e "${RED}Documentation build failed${NC}" >&2
    exit 1
fi

# Build Docker images
if command -v docker >/dev/null 2>&1; then
    echo -e "\n${YELLOW}🐳 Building Docker images...${NC}"

    # Build resolver
    if ! docker build -t "accu-did/resolver:$TAG" -t "accu-did/resolver:latest" \
        -f "$ROOT/drivers/resolver/Dockerfile" "$ROOT"; then
        echo -e "${RED}Resolver Docker build failed${NC}" >&2
        exit 1
    fi

    # Build registrar
    if ! docker build -t "accu-did/registrar:$TAG" -t "accu-did/registrar:latest" \
        -f "$ROOT/drivers/registrar/Dockerfile" "$ROOT"; then
        echo -e "${RED}Registrar Docker build failed${NC}" >&2
        exit 1
    fi

    echo -e "${GREEN}✅ Docker images built and tagged${NC}"
fi

# Create git tag
echo -e "\n${YELLOW}🏷️  Creating git tag: $TAG${NC}"
if ! git tag -a "$TAG" -m "Release $TAG"; then
    echo -e "${RED}Failed to create git tag${NC}" >&2
    exit 1
fi

# Success message
echo -e "\n${CYAN}$(printf '=%.0s' {1..60})${NC}"
echo -e "${GREEN}✅ Release $TAG prepared successfully!${NC}"
echo -e "\n${YELLOW}Next steps:${NC}"
echo -e "  ${WHITE}git push origin main${NC}"
echo -e "  ${WHITE}git push origin $TAG${NC}"
if command -v docker >/dev/null 2>&1; then
    echo -e "\n${YELLOW}Optional - Push Docker images:${NC}"
    echo -e "  ${WHITE}docker push accu-did/resolver:$TAG${NC}"
    echo -e "  ${WHITE}docker push accu-did/resolver:latest${NC}"
    echo -e "  ${WHITE}docker push accu-did/registrar:$TAG${NC}"
    echo -e "  ${WHITE}docker push accu-did/registrar:latest${NC}"
fi