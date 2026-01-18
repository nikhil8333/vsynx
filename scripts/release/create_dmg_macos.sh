#!/usr/bin/env bash
set -euo pipefail

APP_NAME="Vsynx Manager.app"
DMG_NAME=${1:?"output dmg base name required (e.g., vsynx-manager-macos-arm64)"}

cd build/bin

echo "=== DEBUG: build/bin contents ==="
ls -R

echo "=== Preparing ${APP_NAME} ==="

# Find existing .app bundle first (before any deletion)

APP_BUNDLE=$(ls -1d *.app 2>/dev/null | head -n 1 || true)
if [[ -z "${APP_BUNDLE}" ]]; then
  echo "ERROR: No .app bundle found in build/bin" >&2
  exit 1
fi

echo "Found app bundle: ${APP_BUNDLE}"
if [[ "${APP_BUNDLE}" != "${APP_NAME}" ]]; then
  # Only remove target if we need to rename a different bundle into it
  rm -rf "${APP_NAME}"

  mv "${APP_BUNDLE}" "${APP_NAME}"
fi

echo "=== Code signing ==="
# Use proper signing if available, otherwise ad-hoc
if [[ -n "${APPLE_CERTIFICATE:-}" ]]; then
    echo "Running Developer ID signing script..."
    bash "$(dirname "$0")/sign_macos.sh" "${APP_NAME}"
else
    echo "No APPLE_CERTIFICATE - using ad-hoc signing"
    codesign --force --deep --sign - "${APP_NAME}"
fi

echo "=== Creating DMG ==="

DMG_ARGS=(
  --volname "Vsynx Manager"
  --window-pos 200 120
  --window-size 800 400
  --icon-size 100
  --icon "${APP_NAME}" 200 190
  --hide-extension "${APP_NAME}"
  --app-drop-link 600 185
)

if [[ -f "../../build/appicon.icns" ]]; then
  DMG_ARGS+=(--volicon "../../build/appicon.icns")
fi

create-dmg "${DMG_ARGS[@]}" "${DMG_NAME}.dmg" "${APP_NAME}"

