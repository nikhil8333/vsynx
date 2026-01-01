# Vsynx Manager

A cross-platform desktop application for securely managing VS Code extensions. Built with Wails (Go + React), Vsynx validates extension identities, classifies trust levels, and ensures you're using official, verified extensions from the Microsoft Marketplace.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)
![Wails](https://img.shields.io/badge/Wails-v2-DF3C46)
![Version](https://img.shields.io/badge/version-1.0.0-green.svg)

## Features

- 🔒 **Security Validation**: Compare extensions against Microsoft Marketplace
- 🎯 **Trust Classification**: Automatically classify extensions as Legitimate, Suspicious, or Malicious
- 📊 **Audit Reports**: Comprehensive security audits of all installed extensions
- 🔄 **Extension Sync**: Sync extensions between multiple editors (VS Code, Windsurf, Cursor, etc.)
- 🖥️ **Multi-Editor Support**: Manage extensions across VS Code, Windsurf, Cursor, VSCodium, and more
- 🎨 **Modern UI**: Beautiful React interface with Tailwind CSS
- ⚡ **CLI & GUI**: Use from command line (`vsynx`) or desktop application
- 🔍 **Marketplace Search**: Search and install extensions directly

## Quick Start

### Download

Download the latest release for your platform from the [Releases page](https://github.com/yourusername/vsynx/releases).

### Build from Source

See [BUILD.md](BUILD.md) for detailed build instructions.

```bash
git clone https://github.com/yourusername/vsynx.git
cd vsynx
go mod download
cd frontend && npm install && cd ..
wails build
```

## Usage

### GUI Application

Launch the desktop application:

```bash
# Run the built binary
./vsynx-manager      # macOS/Linux
vsynx-manager.exe    # Windows

# Or in development mode
wails dev
```

### CLI Commands

The `vsynx` CLI provides powerful command-line tools:

```bash
# List all supported editors
vsynx editors list

# Check editor status
vsynx editors status vscode

# Validate an extension
vsynx validate ms-python.python

# Audit all extensions
vsynx audit --path ~/.vscode/extensions

# Search marketplace
vsynx marketplace search python

# Sync extensions between editors
vsynx sync preview --from vscode --to windsurf --all
vsynx sync run --from vscode --to cursor --all

# Install extensions
vsynx install ms-python.python github.copilot
```

## Documentation

| Document | Description |
|----------|-------------|
| [BUILD.md](BUILD.md) | Build from source instructions |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contribution guidelines & versioning |
| [CHANGELOG.md](CHANGELOG.md) | Version history |
| [docs/](docs/index.md) | Full documentation |

### Documentation Index

- **Usage**: [GUI Guide](docs/usage/gui.md) · [CLI Testing](docs/usage/cli-testing.md)
- **Development**: [Debugging](docs/development/debugging.md) · [Testing](docs/development/testing.md)
- **Reference**: [Extension Paths](docs/reference/paths.md) · [Manual Paths](docs/reference/manual-paths.md)
- **Project**: [Features](docs/features.md) · [Roadmap](docs/roadmap.md)

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Vsynx Manager                         │
├──────────────────────┬──────────────────────────────────┤
│   CLI Interface      │      GUI Interface (Wails)       │
│   (vsynx)            │      (React + Tailwind CSS)      │
├──────────────────────┴──────────────────────────────────┤
│              Shared Go Packages (Core Logic)            │
│  ┌──────────────┬──────────────┬─────────────────┐      │
│  │ Validation   │ Marketplace  │ Editor Manager  │      │
│  │ Engine       │ Client       │                 │      │
│  └──────────────┴──────────────┴─────────────────┘      │
├─────────────────────────────────────────────────────────┤
│              Microsoft Marketplace API                   │
└─────────────────────────────────────────────────────────┘
```

## Trust Classification

| Level | Icon | Description |
|-------|------|-------------|
| **Legitimate** | 🟢 | Metadata matches Microsoft Marketplace - safe to use |
| **Suspicious** | 🟡 | Minor differences detected - review manually |
| **Malicious** | 🔴 | Critical mismatches - do NOT use |
| **Unknown** | ⚪ | Not found or validation failed - investigate |

## Supported Editors

- VS Code
- VS Code Insiders
- VSCodium
- Windsurf
- Cursor
- Kiro

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines, including:
- Conventional Commits format
- Semantic Versioning strategy
- Development setup
- Pull request process

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/vsynx/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/vsynx/discussions)

---

Made with ❤️ for the VS Code community
