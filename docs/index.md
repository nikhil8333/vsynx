# Vsynx Manager Documentation

Welcome to the Vsynx Manager documentation. This guide covers all aspects of the application, from basic usage to advanced development topics.

## Quick Links

### Usage Guides
- **[GUI Guide](usage/gui.md)** - Using the desktop application
- **[CLI Testing](usage/cli-testing.md)** - Testing CLI commands during development

### Development
- **[Debugging Guide](development/debugging.md)** - Troubleshooting and debugging
- **[Testing Guide](development/testing.md)** - Running and writing tests

### Reference
- **[Extension Paths](reference/paths.md)** - Common extension directory locations
- **[Manual Path Configuration](reference/manual-paths.md)** - Configuring custom paths

### Project Information
- **[Features](features.md)** - Current features and capabilities
- **[Roadmap](roadmap.md)** - Planned improvements and future features
- **[History](history.md)** - Development history and past changes

## Getting Started

### For Users

1. Download the latest release from the [Releases page](https://github.com/yourusername/vsynx/releases)
2. Install the application for your platform
3. Launch Vsynx Manager and start managing your extensions

### For Developers

1. See [BUILD.md](../BUILD.md) for build instructions
2. See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines
3. Check the [CHANGELOG.md](../CHANGELOG.md) for version history

## Architecture Overview

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

## Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/vsynx/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/vsynx/discussions)
