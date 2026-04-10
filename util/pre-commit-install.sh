#!/bin/bash
# QMK Pre-commit Hook Installer
# Run this script to install the pre-commit hook for automatic keyboard compilation

set -e

HOOK_SOURCE="$(cd "$(dirname "$0")/.." && pwd)/hooks/pre-commit"
HOOK_TARGET="$(pwd)/.git/hooks/pre-commit"
HOOK_DIR="$(pwd)/.git/hooks"

echo "=========================================="
echo "QMK Pre-commit Hook Installer"
echo "=========================================="
echo ""

if [ ! -d ".git" ]; then
    echo "[ERROR] This is not a QMK firmware repository (no .git directory)"
    echo "Run this script from the root of your QMK firmware clone"
    exit 1
fi

if [ ! -f "$HOOK_SOURCE" ]; then
    echo "[ERROR] Hook source file not found: $HOOK_SOURCE"
    exit 1
fi

mkdir -p "$HOOK_DIR"

if [ -f "$HOOK_TARGET" ]; then
    echo "[INFO] A pre-commit hook already exists"
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "[ABORTED] Installation cancelled"
        exit 0
    fi
fi

cp "$HOOK_SOURCE" "$HOOK_TARGET"
chmod +x "$HOOK_TARGET"

echo "[OK] Pre-commit hook installed successfully!"
echo ""
echo "The following keyboards will be compiled before each commit:"
echo "  - crkbd:zonekeyboards"
echo "  - sofle:zonekeyboards"
echo "  - lily58:zonekeyboards"
echo "  - lily58zk:zonekeyboards"
echo "  - reviung41:zonekeyboards"
echo "  - zkpad:default"
echo ""
echo "To skip the hook temporarily, use: git commit --no-verify"