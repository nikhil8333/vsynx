# Build Instructions

This document provides detailed instructions for building Vsynx Manager on different platforms.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
- [Development Build](#development-build)
- [Production Build](#production-build)
- [Platform-Specific Builds](#platform-specific-builds)
- [Code Signing](#code-signing)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### All Platforms

1. **Go 1.21 or higher**
   ```bash
   go version  # Should show 1.21 or higher
   ```
   Download from: https://golang.org/dl/

2. **Node.js 18 or higher**
   ```bash
   node --version  # Should show v18.0.0 or higher
   npm --version
   ```
   Download from: https://nodejs.org/

3. **Wails CLI**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   wails doctor  # Check installation
   ```

### Windows

- **WebView2 Runtime** (usually pre-installed on Windows 10/11)
  - Download from: https://developer.microsoft.com/microsoft-edge/webview2/
- **GCC Compiler** (for CGo, optional):
  - MinGW-w64: https://www.mingw-w64.org/
  - Or use TDM-GCC: https://jmeubank.github.io/tdm-gcc/

### macOS

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p
```

### Linux (Debian/Ubuntu)

```bash
sudo apt update
sudo apt install -y build-essential libgtk-3-dev libwebkit2gtk-4.0-dev
```

### Linux (Fedora/RHEL)

```bash
sudo dnf install -y gtk3-devel webkit2gtk3-devel
```

### Linux (Arch)

```bash
sudo pacman -S gtk3 webkit2gtk
```

## Initial Setup

1. **Clone the repository**:
```bash
git clone https://github.com/yourusername/secureopenvsx.git
cd secureopenvsx
```

2. **Install Go dependencies**:
```bash
go mod download
go mod verify
```

3. **Install frontend dependencies**:
```bash
cd frontend
npm install
cd ..
```

4. **Verify setup**:
```bash
wails doctor
```

## Development Build

### GUI Development Mode (Recommended)

Start the application in development mode with hot reload:

```bash
wails dev
```

This will:
- Start the Vite dev server on http://localhost:3000
- Compile and run the Go backend
- Open the application window
- Watch for file changes and auto-reload

**Features in Dev Mode:**
- Hot reload for React components
- Go code changes trigger rebuild
- Dev tools available (F12)
- Source maps enabled

### CLI Development

Run CLI commands directly with Go:

```bash
go run main.go validate ms-python.python
go run main.go audit
go run main.go list
```

### Building CLI Only

```bash
# Build CLI executable
go build -o vsynx-manager.exe .

# Run CLI
./securevsx validate ms-python.python
```

## Production Build

### Build GUI Application

```bash
# Build for current platform
wails build

# Build with custom output directory
wails build -o ./dist

# Build with optimizations
wails build -clean -trimpath
```

**Build flags:**
- `-clean`: Clean build directory before building
- `-trimpath`: Remove file paths from executable
- `-ldflags`: Pass flags to Go linker
- `-tags`: Build tags (e.g., production)
- `-upx`: Compress binary with UPX (requires UPX installed)
- `-webview2`: Embed WebView2 (Windows only)

**Output location:**
- Windows: `build/bin/vsynx-manager.exe`
- macOS: `build/bin/Vsynx Manager.app`
- Linux: `build/bin/securevsx`

### Build CLI Application

```bash
# Standard build
go build -o vsynx-manager.exe .

# Optimized build
go build -ldflags="-s -w" -trimpath -o vsynx-manager.exe .
```

**Optimization flags:**
- `-ldflags="-s -w"`: Strip debug information
- `-trimpath`: Remove file system paths
- Reduces binary size by ~30%

### Build Both CLI and GUI

```bash
# Build GUI
wails build -clean

# Build CLI
go build -ldflags="-s -w" -trimpath -o ./build/bin/securevsx-cli.exe .
```

## Platform-Specific Builds

### Cross-Compilation

Wails supports cross-platform builds with the `-platform` flag.

#### Build for Windows

```bash
# From any platform
wails build -platform windows/amd64

# With WebView2 embedded (increases size)
wails build -platform windows/amd64 -webview2 embed
```

#### Build for macOS

```bash
# Intel Macs
wails build -platform darwin/amd64

# Apple Silicon (M1/M2)
wails build -platform darwin/arm64

# Universal binary (both architectures)
wails build -platform darwin/universal
```

#### Build for Linux

```bash
# 64-bit Linux
wails build -platform linux/amd64

# ARM64 Linux (Raspberry Pi, etc.)
wails build -platform linux/arm64
```

### Build All Platforms

```bash
# Build for all major platforms
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64
```

Or use a script:

```bash
#!/bin/bash
platforms=("windows/amd64" "darwin/universal" "linux/amd64")

for platform in "${platforms[@]}"; do
    echo "Building for $platform..."
    wails build -platform "$platform" -clean -trimpath
done

echo "All builds completed!"
```

## Code Signing

### macOS Code Signing

Required for distribution outside the App Store.

#### Prerequisites

- Apple Developer Account
- Developer ID Application certificate installed in Keychain

#### Sign the Application

```bash
# Build first
wails build -platform darwin/universal -clean

# Sign
codesign --force --deep --sign "Developer ID Application: Your Name (TEAMID)" \
  ./build/bin/Vsynx Manager.app

# Verify signature
codesign --verify --verbose=4 ./build/bin/Vsynx Manager.app
spctl --assess --verbose=4 ./build/bin/Vsynx Manager.app
```

#### Notarize for Gatekeeper

```bash
# Create a ZIP archive
cd build/bin
ditto -c -k --keepParent Vsynx Manager.app SecureVSX.zip

# Submit for notarization (requires notarytool profile configured)
xcrun notarytool submit SecureVSX.zip \
  --keychain-profile "notarytool-profile" \
  --wait

# Staple the notarization ticket
xcrun stapler staple Vsynx Manager.app

# Verify notarization
xcrun stapler validate Vsynx Manager.app
```

#### Configure notarytool Profile

```bash
xcrun notarytool store-credentials "notarytool-profile" \
  --apple-id "your-apple-id@example.com" \
  --team-id "TEAMID" \
  --password "<APP_SPECIFIC_PASSWORD>"
```

### Windows Code Signing

Required for avoiding SmartScreen warnings.

#### Prerequisites

- Code Signing Certificate (.pfx or .p12 file)
- Windows SDK (for signtool.exe)

#### Sign the Executable

```bash
# Using signtool (Windows SDK)
signtool sign /f certificate.pfx /p <CERT_PASSWORD> \
  /t http://timestamp.digicert.com \
  /fd SHA256 \
  /d "Vsynx Manager" \
  ./build/bin/vsynx-manager.exe

# Verify signature
signtool verify /pa ./build/bin/vsynx-manager.exe
```

#### Alternative: osslsigncode (Cross-platform)

```bash
# Install osslsigncode
# Linux: sudo apt install osslsigncode
# macOS: brew install osslsigncode

osslsigncode sign \
  -pkcs12 certificate.pfx \
  -pass <CERT_PASSWORD> \
  -t http://timestamp.digicert.com \
  -in ./build/bin/vsynx-manager.exe \
  -out ./build/bin/securevsx-signed.exe
```

### Linux Code Signing

Linux doesn't require code signing, but you can create GPG-signed checksums:

```bash
# Create checksums
cd build/bin
sha256sum securevsx > SHA256SUMS

# Sign checksums with GPG
gpg --clearsign SHA256SUMS

# Verify
gpg --verify SHA256SUMS.asc
```

## Advanced Build Options

### Custom Version Information

```bash
# Set version via ldflags
VERSION="1.0.0"
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

wails build -ldflags "\
  -X main.Version=$VERSION \
  -X main.Commit=$COMMIT \
  -X main.BuildTime=$BUILD_TIME"
```

### Reduce Binary Size

```bash
# Use UPX compression (requires UPX installed)
wails build -upx

# Manual UPX compression
upx --best --lzma ./build/bin/vsynx-manager.exe
```

**Warning:** UPX may cause antivirus false positives.

### Custom Build Directory

```bash
wails build -o ./releases/v1.0.0
```

### Skip Frontend Build

```bash
# If frontend is already built
wails build -skipbindings -s
```

## Creating Installers

### Windows Installer (NSIS)

Create `installer.nsi`:

```nsis
!include "MUI2.nsh"

Name "Vsynx Manager"
OutFile "SecureVSX-Setup.exe"
InstallDir "$PROGRAMFILES64\SecureVSX"

!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

Section "Install"
  SetOutPath "$INSTDIR"
  File "build\bin\vsynx-manager.exe"
  CreateShortcut "$DESKTOP\SecureVSX.lnk" "$INSTDIR\vsynx-manager.exe"
SectionEnd
```

Build:
```bash
makensis installer.nsi
```

### macOS DMG

```bash
# Create DMG with applications folder link
create-dmg \
  --volname "SecureVSX Installer" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "Vsynx Manager.app" 200 190 \
  --hide-extension "Vsynx Manager.app" \
  --app-drop-link 600 185 \
  "SecureVSX-Installer.dmg" \
  "build/bin/"
```

### Linux Package (Debian)

Create directory structure:
```bash
mkdir -p securevsx-deb/usr/local/bin
mkdir -p securevsx-deb/DEBIAN

# Copy binary
cp build/bin/securevsx securevsx-deb/usr/local/bin/

# Create control file
cat > securevsx-deb/DEBIAN/control << EOF
Package: securevsx
Version: 1.0.0
Architecture: amd64
Maintainer: Your Name <you@example.com>
Description: Vsynx Manager
 Secure VS Code extension management tool
EOF

# Build package
dpkg-deb --build securevsx-deb
```

## Troubleshooting

### Build Fails with "wails: command not found"

```bash
# Ensure $GOPATH/bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Reinstall Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Frontend Build Errors

```bash
# Clean and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
cd ..

# Rebuild
wails build -clean
```

### WebView2 Issues (Windows)

Install WebView2 Runtime or embed it:
```bash
wails build -webview2 embed
```

### macOS Gatekeeper Issues

```bash
# Remove quarantine attribute
xattr -cr ./build/bin/Vsynx Manager.app

# Or sign and notarize properly (see Code Signing section)
```

### Large Binary Size

```bash
# Use all size optimizations
wails build -clean -trimpath -ldflags "-s -w" -upx

# Typical sizes:
# - Without optimization: 50-80 MB
# - With optimization: 20-30 MB
# - With UPX: 10-15 MB
```

### CGo Compilation Errors

```bash
# Windows: Ensure MinGW-w64 is installed and in PATH
# macOS: Install Xcode Command Line Tools
# Linux: Install build-essential

# Disable CGo if not needed
CGO_ENABLED=0 go build
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Build

on: [push, pull_request]

jobs:
  build:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    
    steps:
      - uses: actions/checkout@v3
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
      
      - name: Install dependencies
        run: |
          cd frontend
          npm install
          cd ..
      
      - name: Build
        run: wails build -clean
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: securevsx-${{ matrix.os }}
          path: build/bin/*
```

## Additional Resources

- [Wails Documentation](https://wails.io/docs/introduction)
- [Wails Build Flags](https://wails.io/docs/reference/cli#build)
- [Go Build Flags](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)
- [Vite Build Options](https://vitejs.dev/guide/build.html)

---

For more information, see the main [README.md](README.md).
