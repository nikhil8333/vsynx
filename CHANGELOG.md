# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-12-31

### Added

- **GUI Application (Vsynx Manager)**
  - Modern React + Tailwind CSS desktop interface
  - Extension browser with search and filtering
  - Detailed extension metadata view
  - Security validation with trust level classification
  - Audit view with comprehensive security reports
  - Marketplace search integration
  - Extension sync between editors
  - Settings view with CLI installation management
  - Multi-editor support (VS Code, Windsurf, Cursor, VSCodium, Kiro)

- **CLI Tool (vsynx)**
  - `vsynx validate <extension-id>` - Validate individual extensions
  - `vsynx audit` - Audit all installed extensions
  - `vsynx list` - List installed extensions
  - `vsynx download` - Download official extensions
  - `vsynx editors list` - List supported editors
  - `vsynx editors status` - Check editor status
  - `vsynx editors extensions` - List extensions for an editor
  - `vsynx marketplace search` - Search Microsoft Marketplace
  - `vsynx marketplace open` - Open extension in browser
  - `vsynx sync preview` - Preview sync operation
  - `vsynx sync run` - Execute extension sync
  - `vsynx sync conflicts` - Detect sync conflicts
  - `vsynx install` - Install extensions via editor CLI
  - JSON output support for all commands

- **Core Features**
  - Extension validation against Microsoft Marketplace
  - Trust level classification (Legitimate, Suspicious, Malicious, Unknown)
  - Cross-platform support (Windows, macOS, Linux)
  - Automatic editor detection
  - Extension metadata comparison
  - Sync extensions between multiple editors

- **Documentation**
  - Comprehensive README with quick start guide
  - BUILD.md with platform-specific build instructions
  - CONTRIBUTING.md with versioning strategy and conventional commits
  - Full documentation in /docs folder

- **Build & Release**
  - GitHub Actions CI/CD pipeline
  - Multi-platform release builds (Windows, macOS, Linux)
  - Support for multiple architectures (amd64, arm64, 386, arm)
  - Portable archives and installers

### Security

- All API requests use HTTPS
- No telemetry or data collection
- Local-first processing
- Strict Content Security Policy in GUI

---

## Version History

For changes prior to v1.0.0, see [docs/history.md](docs/history.md).
