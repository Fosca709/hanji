#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$ROOT_DIR/build"
ENTRYPOINT="$ROOT_DIR/main.py"

if ! command -v uv >/dev/null 2>&1; then
    echo "uv is required but was not found in PATH." >&2
    exit 1
fi

if ! command -v patchelf >/dev/null 2>&1; then
    echo "patchelf is required for Nuitka standalone/onefile builds on Linux." >&2
    echo "Install it first, for example: sudo apt install patchelf" >&2
    exit 1
fi

uv run nuitka \
    --onefile \
    --assume-yes-for-downloads \
    --enable-plugin=pyside6 \
    --include-qt-plugins=qml,platforminputcontexts \
    --linux-onefile-icon="$ROOT_DIR/assets/icon.svg" \
    --include-data-dir="$ROOT_DIR/qml=qml" \
    --include-data-dir="$ROOT_DIR/assets=assets" \
    --output-dir="$BUILD_DIR" \
    --output-filename=hanji.bin \
    --remove-output \
    "$ENTRYPOINT"

echo
echo "Build complete."
echo "Output file: $BUILD_DIR/hanji.bin"
