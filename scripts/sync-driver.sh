#!/bin/bash
# Syncs driver files from upstream git submodule to internal/driver/files/ for //go:embed
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🔄 Updating upstream submodule..."
git submodule update --init --recursive --remote "$REPO_ROOT/internal/driver/upstream"

echo "📂 Copying driver files to internal/driver/files..."
mkdir -p "$REPO_ROOT/internal/driver/files/src"
cp "$REPO_ROOT/internal/driver/upstream/src/linuwu_sense.c" "$REPO_ROOT/internal/driver/files/src/"
cp "$REPO_ROOT/internal/driver/upstream/Makefile" "$REPO_ROOT/internal/driver/files/"
cp "$REPO_ROOT/internal/driver/upstream/linuwu_sense.service" "$REPO_ROOT/internal/driver/files/"
cp "$REPO_ROOT/internal/driver/upstream/module_signing_readme" "$REPO_ROOT/internal/driver/files/"

echo "✓ Driver files synced successfully!"
