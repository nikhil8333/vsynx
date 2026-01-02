# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1](https://github.com/nikhil8333/vsynx/compare/vsynx-v1.0.0...vsynx-v1.0.1) (2026-01-02)


### Bug Fixes

* **release:** update WebKit dependency for Ubuntu 24.04 ([#2](https://github.com/nikhil8333/vsynx/issues/2)) ([b32a86b](https://github.com/nikhil8333/vsynx/commit/b32a86ba5fad29286cccd791e7e81ef4002cd732))
* **release:** use ubuntu-22.04 for Linux builds, remove cross-compile targets ([#3](https://github.com/nikhil8333/vsynx/issues/3)) ([8e67af6](https://github.com/nikhil8333/vsynx/commit/8e67af6c06793e52ad52547beb54158518c9a483))

## [Unreleased]

### Added

- Comprehensive automated testing across backend and frontend
  - Go unit tests for `internal/editor` (including sync logic)
  - Frontend unit tests using Vitest + Testing Library
  - End-to-end tests using Playwright covering core user flows
- New test commands for frontend unit tests, coverage, and E2E runs

### Changed

- Updated documentation to include test commands and watch-mode notes

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
