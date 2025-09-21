#!/usr/bin/env bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}🚀 Starting local CI build...${NC}"

ROOT_DIR="$(cd "$(dirname "$0")/.."; pwd)"
FAILED=0

# Test resolver-go
echo -e "\n${YELLOW}📋 Testing resolver-go...${NC}"
(cd "$ROOT_DIR/resolver-go" && go test ./... -v) || {
    echo -e "${RED}❌ Resolver tests failed${NC}"
    FAILED=1
}
[ $FAILED -eq 0 ] && echo -e "${GREEN}✅ Resolver tests passed${NC}"

# Test registrar-go
echo -e "\n${YELLOW}📋 Testing registrar-go...${NC}"
(cd "$ROOT_DIR/registrar-go" && go test ./... -v) || {
    echo -e "${RED}❌ Registrar tests failed${NC}"
    FAILED=1
}
[ $FAILED -eq 0 ] && echo -e "${GREEN}✅ Registrar tests passed${NC}"

# Build documentation
echo -e "\n${YELLOW}📚 Building documentation...${NC}"
if bash "$ROOT_DIR/scripts/build-docs.sh"; then
    echo -e "${GREEN}✅ Documentation built${NC}"
else
    echo -e "${RED}❌ Documentation build failed${NC}"
    FAILED=1
fi

# Build Docker images (optional)
if command -v docker >/dev/null 2>&1; then
    echo -e "\n${YELLOW}🐳 Building Docker images...${NC}"

    # Build resolver image
    if docker build -t accu-did/resolver:latest -f "$ROOT_DIR/drivers/resolver/Dockerfile" "$ROOT_DIR"; then
        echo -e "${GREEN}✅ Resolver Docker image built${NC}"
    else
        echo -e "${RED}❌ Resolver Docker build failed${NC}"
        FAILED=1
    fi

    # Build registrar image
    if docker build -t accu-did/registrar:latest -f "$ROOT_DIR/drivers/registrar/Dockerfile" "$ROOT_DIR"; then
        echo -e "${GREEN}✅ Registrar Docker image built${NC}"
    else
        echo -e "${RED}❌ Registrar Docker build failed${NC}"
        FAILED=1
    fi
else
    echo -e "\n${YELLOW}⚠️  Docker not found, skipping image builds${NC}"
fi

# Summary
echo -e "\n${CYAN}$(printf '=%.0s' {1..60})${NC}"
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ LOCAL CI PASSED${NC}"
    exit 0
else
    echo -e "${RED}❌ LOCAL CI FAILED${NC}"
    exit 1
fi