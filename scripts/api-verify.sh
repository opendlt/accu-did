#!/usr/bin/env bash
# scripts/api-verify.sh
# Validates OpenAPI specifications using redoc-cli
# Usage: ./scripts/api-verify.sh [--verbose]

set -euo pipefail

VERBOSE=false
if [[ "${1:-}" == "--verbose" ]]; then
  VERBOSE=true
fi

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
SPEC_DIR="$ROOT/docs/spec/openapi"

RESOLVER_YAML="$SPEC_DIR/resolver.yaml"
REGISTRAR_YAML="$SPEC_DIR/registrar.yaml"

if [[ "$VERBOSE" == "true" ]]; then
  echo "[INFO] Root directory: $ROOT"
  echo "[INFO] Spec directory: $SPEC_DIR"
fi

test_file() {
  local path="$1"
  local name="$2"

  if [[ ! -f "$path" ]]; then
    echo "[ERROR] $name not found: $path" >&2
    exit 1
  fi

  if [[ "$VERBOSE" == "true" ]]; then
    echo "[INFO] Found $name: $path"
  fi
}

lint_openapi() {
  local yaml_path="$1"
  local name="$2"

  echo "[INFO] Validating $name..."

  # Use npx to run swagger-cli validate command
  # --yes flag automatically installs if not present
  # Capture both stdout and stderr, but preserve exit code
  local output
  local exit_code=0

  output=$(npx --yes swagger-cli validate "$yaml_path" 2>&1) || exit_code=$?

  if [[ $exit_code -eq 0 ]]; then
    echo -e "\033[32m[OK] $name is valid\033[0m"
    if [[ "$VERBOSE" == "true" && -n "$output" ]]; then
      echo "$output"
    fi
  else
    echo -e "\033[31m[ERROR] $name validation failed:\033[0m" >&2
    echo -e "\033[31m$output\033[0m" >&2
    exit $exit_code
  fi
}

# Verify files exist
test_file "$RESOLVER_YAML" "Resolver OpenAPI spec"
test_file "$REGISTRAR_YAML" "Registrar OpenAPI spec"

# Check for npx availability
if ! command -v npx >/dev/null 2>&1; then
  echo "[ERROR] npx not found. Please install Node.js and npm." >&2
  exit 1
fi

if [[ "$VERBOSE" == "true" ]]; then
  npm_version=$(npx --version 2>/dev/null || echo "unknown")
  echo "[INFO] Found npx version: $npm_version"
fi

echo ""
echo -e "\033[36m=== OpenAPI Specification Validation ===\033[0m"
echo ""

# Validate both OpenAPI specifications
lint_openapi "$RESOLVER_YAML" "Resolver API (resolver.yaml)"
lint_openapi "$REGISTRAR_YAML" "Registrar API (registrar.yaml)"

echo ""
echo -e "\033[36m=== Validation Summary ===\033[0m"
echo -e "\033[32m[OK] All OpenAPI specifications are valid\033[0m"
echo ""

# Additional checks
resolver_content=$(cat "$RESOLVER_YAML")
registrar_content=$(cat "$REGISTRAR_YAML")

# Check for API freeze markers
if echo "$resolver_content" | grep -q "x-api-freeze:.*true"; then
  echo -e "\033[32m[OK] Resolver API freeze marker found\033[0m"
else
  echo -e "\033[33m[WARN] Resolver API freeze marker missing\033[0m"
fi

if echo "$registrar_content" | grep -q "x-api-freeze:.*true"; then
  echo -e "\033[32m[OK] Registrar API freeze marker found\033[0m"
else
  echo -e "\033[33m[WARN] Registrar API freeze marker missing\033[0m"
fi

# Check version consistency
resolver_version=$(echo "$resolver_content" | grep -E "^\s*version:" | head -1 | sed -E 's/^\s*version:\s*["\'"'"']?([^"\'"'"'\r\n]+)["\'"'"']?.*/\1/' | tr -d ' ')
registrar_version=$(echo "$registrar_content" | grep -E "^\s*version:" | head -1 | sed -E 's/^\s*version:\s*["\'"'"']?([^"\'"'"'\r\n]+)["\'"'"']?.*/\1/' | tr -d ' ')

echo ""
echo -e "\033[36m=== Version Information ===\033[0m"
echo "Resolver API version:  ${resolver_version:-unknown}"
echo "Registrar API version: ${registrar_version:-unknown}"

if [[ "$resolver_version" == "$registrar_version" ]] && [[ -n "$resolver_version" ]]; then
  echo -e "\033[32m[OK] API versions are synchronized\033[0m"
else
  echo -e "\033[33m[WARN] API versions are not synchronized\033[0m"
  echo "       Consider updating to maintain version consistency"
fi

echo ""
echo -e "\033[32m[OK] API verification completed successfully\033[0m"