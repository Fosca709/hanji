#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_SCRIPT="$ROOT_DIR/scripts/build.sh"
BUILD_DIR="$ROOT_DIR/build"
DIST_DIR="$ROOT_DIR/dist"
STAGING_DIR="$ROOT_DIR/.pkgroot"
DESKTOP_FILE_SOURCE="$ROOT_DIR/packaging/deb/hanji.desktop"
ICON_SOURCE="$ROOT_DIR/assets/icon.svg"
BINARY_SOURCE="$BUILD_DIR/hanji.bin"

PACKAGE_NAME="${PACKAGE_NAME:-hanji}"
PACKAGE_VERSION="${PACKAGE_VERSION:-0.1.0}"
PACKAGE_ARCH="${PACKAGE_ARCH:-$(dpkg --print-architecture)}"
PACKAGE_SECTION="${PACKAGE_SECTION:-utils}"
PACKAGE_PRIORITY="${PACKAGE_PRIORITY:-optional}"
PACKAGE_MAINTAINER="${PACKAGE_MAINTAINER:-Hanji <noreply@example.com>}"
PACKAGE_DESCRIPTION="${PACKAGE_DESCRIPTION:-Hanji sticky note app}"
PACKAGE_LONG_DESCRIPTION="${PACKAGE_LONG_DESCRIPTION:-A tiny sticky note desktop app.}"
PACKAGE_DEPENDS="${PACKAGE_DEPENDS:-libgl1, libxkbcommon0}"
PACKAGE_FILENAME="${PACKAGE_NAME}_${PACKAGE_VERSION}_${PACKAGE_ARCH}.deb"

if ! command -v dpkg-deb >/dev/null 2>&1; then
    echo "dpkg-deb is required but was not found in PATH." >&2
    exit 1
fi

if ! command -v dpkg >/dev/null 2>&1; then
    echo "dpkg is required but was not found in PATH." >&2
    exit 1
fi

bash "$BUILD_SCRIPT"

if [[ ! -f "$BINARY_SOURCE" ]]; then
    echo "Expected build output was not found: $BINARY_SOURCE" >&2
    exit 1
fi

if [[ ! -f "$DESKTOP_FILE_SOURCE" ]]; then
    echo "Desktop entry template was not found: $DESKTOP_FILE_SOURCE" >&2
    exit 1
fi

if [[ ! -f "$ICON_SOURCE" ]]; then
    echo "Icon asset was not found: $ICON_SOURCE" >&2
    exit 1
fi

rm -rf "$STAGING_DIR"
mkdir -p \
    "$STAGING_DIR/DEBIAN" \
    "$STAGING_DIR/opt/$PACKAGE_NAME" \
    "$STAGING_DIR/usr/share/applications" \
    "$STAGING_DIR/usr/share/icons/hicolor/scalable/apps" \
    "$DIST_DIR"

install -Dm755 "$BINARY_SOURCE" "$STAGING_DIR/opt/$PACKAGE_NAME/$PACKAGE_NAME"
install -Dm644 "$ICON_SOURCE" "$STAGING_DIR/usr/share/icons/hicolor/scalable/apps/$PACKAGE_NAME.svg"
install -Dm644 "$DESKTOP_FILE_SOURCE" "$STAGING_DIR/usr/share/applications/$PACKAGE_NAME.desktop"

cat > "$STAGING_DIR/DEBIAN/control" <<EOF
Package: $PACKAGE_NAME
Version: $PACKAGE_VERSION
Section: $PACKAGE_SECTION
Priority: $PACKAGE_PRIORITY
Architecture: $PACKAGE_ARCH
Maintainer: $PACKAGE_MAINTAINER
Depends: $PACKAGE_DEPENDS
Description: $PACKAGE_DESCRIPTION
 $PACKAGE_LONG_DESCRIPTION
EOF

dpkg-deb --build "$STAGING_DIR" "$DIST_DIR/$PACKAGE_FILENAME"

echo
echo "Debian package complete."
echo "Output file: $DIST_DIR/$PACKAGE_FILENAME"
