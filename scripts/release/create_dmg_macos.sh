#!/usr/bin/env bash
set -euo pipefail

APP_NAME="Vsynx Manager.app"
DMG_NAME=${1:?"output dmg base name required (e.g., vsynx-manager-macos-arm64)"}

cd build/bin

echo "=== DEBUG: build/bin contents ==="
ls -R

echo "=== Preparing ${APP_NAME} ==="
rm -rf "${APP_NAME}"

APP_BUNDLE=$(ls -1d *.app 2>/dev/null | head -n 1 || true)
if [[ -z "${APP_BUNDLE}" ]]; then
  echo "ERROR: No .app bundle found in build/bin" >&2
  exit 1
fi

echo "Found app bundle: ${APP_BUNDLE}"
if [[ "${APP_BUNDLE}" != "${APP_NAME}" ]]; then
  mv "${APP_BUNDLE}" "${APP_NAME}"
fi

echo "=== Ad-hoc signing ==="
codesign --force --deep --sign - "${APP_NAME}"

echo "=== Creating DMG ==="
ICON_ARG=""
if [[ -f "../../build/appicon.icns" ]]; then
  ICON_ARG=(--volicon "../../build/appicon.icns")
fi

create-dmg \
  --volname "Vsynx Manager" \
  "${ICON_ARG[@]}" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "${APP_NAME}" 200 190 \
  --hide-extension "${APP_NAME}" \
  --app-drop-link 600 185 \
  "${DMG_NAME}.dmg" \
  "${APP_NAME}"
