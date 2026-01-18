#!/usr/bin/env bash
set -euo pipefail

# macOS Code Signing and Notarization Script
# Requires secrets:
#   APPLE_CERTIFICATE - base64 encoded .p12 Developer ID Application certificate
#   APPLE_CERTIFICATE_PASSWORD - password for the .p12 file
#   APPLE_ID - Apple ID email for notarization
#   APPLE_ID_PASSWORD - App-specific password for notarization
#   APPLE_TEAM_ID - Apple Developer Team ID

APP_PATH=${1:?"App path required (e.g., build/bin/Vsynx Manager.app)"}
BUNDLE_ID=${2:-"dev.vsynx.manager"}

# Check if signing is configured
if [[ -z "${APPLE_CERTIFICATE:-}" ]]; then
    echo "APPLE_CERTIFICATE not set - skipping code signing"
    echo "To enable signing, add the following secrets:"
    echo "  - APPLE_CERTIFICATE (base64 .p12)"
    echo "  - APPLE_CERTIFICATE_PASSWORD"
    echo "  - APPLE_ID"
    echo "  - APPLE_ID_PASSWORD (app-specific password)"
    echo "  - APPLE_TEAM_ID"
    exit 0
fi

echo "=== macOS Code Signing and Notarization ==="
echo "Signing: $APP_PATH"

# Create temporary keychain
KEYCHAIN_NAME="build-$(date +%s).keychain"
KEYCHAIN_PASSWORD=$(openssl rand -base64 32)
KEYCHAIN_PATH="$HOME/Library/Keychains/$KEYCHAIN_NAME-db"

cleanup() {
    echo "Cleaning up..."
    security delete-keychain "$KEYCHAIN_PATH" 2>/dev/null || true
    rm -f /tmp/certificate.p12
}
trap cleanup EXIT

# Decode and import certificate
echo "Setting up keychain..."
security create-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN_PATH"
security set-keychain-settings -lut 21600 "$KEYCHAIN_PATH"
security unlock-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN_PATH"

echo "Importing certificate..."
echo "$APPLE_CERTIFICATE" | base64 --decode > /tmp/certificate.p12
security import /tmp/certificate.p12 -k "$KEYCHAIN_PATH" -P "${APPLE_CERTIFICATE_PASSWORD}" -T /usr/bin/codesign -T /usr/bin/security

# Set keychain search list
security list-keychains -d user -s "$KEYCHAIN_PATH" $(security list-keychains -d user | tr -d '"')
security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "$KEYCHAIN_PASSWORD" "$KEYCHAIN_PATH"

# Find the Developer ID Application certificate
CERT_IDENTITY=$(security find-identity -v -p codesigning "$KEYCHAIN_PATH" | grep "Developer ID Application" | head -1 | awk -F'"' '{print $2}')

if [[ -z "$CERT_IDENTITY" ]]; then
    echo "ERROR: Developer ID Application certificate not found"
    echo "Available identities:"
    security find-identity -v -p codesigning "$KEYCHAIN_PATH"
    exit 1
fi

echo "Using certificate: $CERT_IDENTITY"

# Sign the app bundle
echo "Signing app bundle..."
codesign --force --options runtime --sign "$CERT_IDENTITY" --timestamp --deep "$APP_PATH"

# Verify signature
echo "Verifying signature..."
codesign --verify --verbose=4 "$APP_PATH"
spctl --assess --verbose=4 --type execute "$APP_PATH" || echo "Note: spctl assessment may fail until notarized"

# Notarization (if credentials provided)
if [[ -n "${APPLE_ID:-}" ]] && [[ -n "${APPLE_ID_PASSWORD:-}" ]] && [[ -n "${APPLE_TEAM_ID:-}" ]]; then
    echo "=== Notarization ==="
    
    # Create zip for notarization
    NOTARIZE_ZIP="/tmp/notarize-$(date +%s).zip"
    ditto -c -k --keepParent "$APP_PATH" "$NOTARIZE_ZIP"
    
    echo "Submitting for notarization..."
    xcrun notarytool submit "$NOTARIZE_ZIP" \
        --apple-id "$APPLE_ID" \
        --password "$APPLE_ID_PASSWORD" \
        --team-id "$APPLE_TEAM_ID" \
        --wait \
        --timeout 30m
    
    # Staple the notarization ticket
    echo "Stapling notarization ticket..."
    xcrun stapler staple "$APP_PATH"
    
    # Verify stapling
    xcrun stapler validate "$APP_PATH"
    
    rm -f "$NOTARIZE_ZIP"
    echo "=== Notarization completed ==="
else
    echo "Skipping notarization - APPLE_ID, APPLE_ID_PASSWORD, or APPLE_TEAM_ID not set"
fi

echo "=== macOS signing completed successfully ==="
