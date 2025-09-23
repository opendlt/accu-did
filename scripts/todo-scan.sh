#!/bin/bash
#
# TODO Scanner - Linux/Docker Wrapper
# Scans the repository for TODO, FIXME, XXX, HACK, and other markers
# Generates reports in JSON, Markdown, and CSV formats
#

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
REPO_PATH="${1:-$(pwd)}"
OUTPUT_DIR="${REPO_PATH}/reports"
DOCKER_IMAGE="${DOCKER_IMAGE:-golang:1.25-alpine}"
USE_DOCKER="${USE_DOCKER:-auto}"

print_usage() {
    cat << EOF
Usage: $0 [REPO_PATH] [OPTIONS]

Scans a repository for TODO markers and generates reports.

Arguments:
    REPO_PATH       Path to repository (default: current directory)

Environment Variables:
    USE_DOCKER      Force docker usage: 'yes', 'no', or 'auto' (default)
    DOCKER_IMAGE    Docker image to use (default: golang:1.25-alpine)

Examples:
    $0                          # Scan current directory
    $0 /path/to/repo           # Scan specific repository
    USE_DOCKER=yes $0          # Force Docker usage
    USE_DOCKER=no $0           # Force local Go usage

Reports are generated in: ./reports/
    - todo-report.json         # Machine-readable JSON
    - todo-report.md          # Human-readable Markdown
    - todo-report.csv         # Spreadsheet-compatible CSV
EOF
}

check_prerequisites() {
    if [[ "$USE_DOCKER" == "auto" ]]; then
        if command -v go >/dev/null 2>&1; then
            echo -e "${GREEN}✓${NC} Found local Go installation"
            USE_DOCKER="no"
        elif command -v docker >/dev/null 2>&1; then
            echo -e "${YELLOW}⚠${NC} No local Go found, using Docker"
            USE_DOCKER="yes"
        else
            echo -e "${RED}✗${NC} Neither Go nor Docker found. Please install one of them."
            exit 1
        fi
    fi

    if [[ "$USE_DOCKER" == "yes" ]] && ! command -v docker >/dev/null 2>&1; then
        echo -e "${RED}✗${NC} Docker not found but USE_DOCKER=yes"
        exit 1
    fi

    if [[ "$USE_DOCKER" == "no" ]] && ! command -v go >/dev/null 2>&1; then
        echo -e "${RED}✗${NC} Go not found but USE_DOCKER=no"
        exit 1
    fi
}

run_local() {
    echo -e "${BLUE}🔍${NC} Running TODO scanner locally..."

    cd "$REPO_PATH"

    # Ensure reports directory exists
    mkdir -p "$OUTPUT_DIR"

    # Run the scanner
    if [[ -f "tools/todoscan/main.go" ]]; then
        go run tools/todoscan/main.go .
    else
        echo -e "${RED}✗${NC} tools/todoscan/main.go not found in repository"
        echo "Please ensure you're running this from the repository root."
        exit 1
    fi
}

run_docker() {
    echo -e "${BLUE}🐳${NC} Running TODO scanner in Docker..."

    # Ensure reports directory exists
    mkdir -p "$OUTPUT_DIR"

    # Check if we have a dev container setup
    if [[ -f "$REPO_PATH/docker-compose.dev.yml" ]]; then
        echo -e "${BLUE}📦${NC} Using docker-compose.dev.yml"
        cd "$REPO_PATH"
        docker-compose -f docker-compose.dev.yml run --rm dev bash -c "
            echo 'Running TODO scanner...'
            if [[ -f tools/todoscan/main.go ]]; then
                go run tools/todoscan/main.go .
            else
                echo 'Error: tools/todoscan/main.go not found'
                exit 1
            fi
        "
    else
        echo -e "${BLUE}🚀${NC} Using standalone Docker container"

        # Run in a standalone Go container
        docker run --rm \
            -v "$REPO_PATH:/workspace" \
            -w /workspace \
            "$DOCKER_IMAGE" \
            sh -c "
                apk add --no-cache git >/dev/null 2>&1 || true
                if [[ -f tools/todoscan/main.go ]]; then
                    echo 'Running TODO scanner...'
                    go run tools/todoscan/main.go .
                else
                    echo 'Error: tools/todoscan/main.go not found'
                    exit 1
                fi
            "
    fi
}

print_results() {
    local json_file="$OUTPUT_DIR/todo-report.json"
    local md_file="$OUTPUT_DIR/todo-report.md"
    local csv_file="$OUTPUT_DIR/todo-report.csv"

    echo
    echo -e "${GREEN}✓${NC} Scan completed successfully!"
    echo

    if [[ -f "$json_file" ]]; then
        local total_count=$(jq -r '.totalCount // 0' "$json_file" 2>/dev/null || echo "unknown")
        echo -e "${BLUE}📊${NC} Found ${YELLOW}$total_count${NC} TODO items"

        # Show summary by tag if jq is available
        if command -v jq >/dev/null 2>&1; then
            echo -e "${BLUE}📋${NC} Summary by tag:"
            jq -r '.summary.countsByTag | to_entries[] | "  - \(.key): \(.value)"' "$json_file" 2>/dev/null || true
        fi
    fi

    echo
    echo -e "${GREEN}📁${NC} Reports generated:"
    for file in "$json_file" "$md_file" "$csv_file"; do
        if [[ -f "$file" ]]; then
            local size=$(ls -lh "$file" | awk '{print $5}')
            echo -e "  ${GREEN}✓${NC} $(basename "$file") (${size})"
        else
            echo -e "  ${RED}✗${NC} $(basename "$file") (missing)"
        fi
    done

    echo
    echo -e "${BLUE}💡${NC} Next steps:"
    echo "  - Review the Markdown report: $md_file"
    echo "  - Import CSV data: $csv_file"
    echo "  - Process JSON programmatically: $json_file"
    echo "  - Filter by tag: grep 'TODO' $md_file"
    echo "  - Filter by directory: grep 'resolver-go' $md_file"
}

main() {
    # Handle help flags
    case "${1:-}" in
        -h|--help|help)
            print_usage
            exit 0
            ;;
    esac

    echo -e "${BLUE}🔍${NC} TODO Scanner for Accumulate DID Repository"
    echo -e "${BLUE}📂${NC} Repository: $REPO_PATH"
    echo

    check_prerequisites

    # Verify repository path exists
    if [[ ! -d "$REPO_PATH" ]]; then
        echo -e "${RED}✗${NC} Repository path does not exist: $REPO_PATH"
        exit 1
    fi

    # Run the scanner
    if [[ "$USE_DOCKER" == "yes" ]]; then
        run_docker
    else
        run_local
    fi

    # Print results summary
    print_results
}

# Only run main if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi