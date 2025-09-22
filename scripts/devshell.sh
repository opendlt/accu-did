#!/usr/bin/env bash
# scripts/devshell.sh - Launch interactive development container shell
set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}🐳 Launching development container shell...${NC}"

# Run the development container with an interactive bash shell
docker compose -f docker-compose.dev.yml run --rm dev "bash -l"

echo -e "${GREEN}✅ Exited development container${NC}"