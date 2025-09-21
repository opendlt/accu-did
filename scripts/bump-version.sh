#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 1 ] || ! [[ "$1" =~ ^(patch|minor|major)$ ]]; then
    echo "Usage: $0 <patch|minor|major>" >&2
    exit 1
fi

BUMP="$1"
ROOT="$(cd "$(dirname "$0")/.."; pwd)"
VERSION_FILE="$ROOT/VERSION"

# Read current version
if [ ! -f "$VERSION_FILE" ]; then
    echo "ERROR: VERSION file not found" >&2
    exit 1
fi

CURRENT_VERSION=$(cat "$VERSION_FILE")
if ! [[ "$CURRENT_VERSION" =~ ^([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
    echo "ERROR: Invalid version format: $CURRENT_VERSION" >&2
    exit 1
fi

MAJOR="${BASH_REMATCH[1]}"
MINOR="${BASH_REMATCH[2]}"
PATCH="${BASH_REMATCH[3]}"

# Bump version
case "$BUMP" in
    major)
        ((MAJOR++))
        MINOR=0
        PATCH=0
        ;;
    minor)
        ((MINOR++))
        PATCH=0
        ;;
    patch)
        ((PATCH++))
        ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"

# Write new version
echo -n "$NEW_VERSION" > "$VERSION_FILE"

# Output new version
echo "$NEW_VERSION"