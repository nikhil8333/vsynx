# Contributing to Vsynx Manager

Thank you for your interest in contributing to Vsynx Manager! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful, inclusive, and constructive in all interactions.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/vsynx.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test thoroughly
6. Commit using Conventional Commits format
7. Push and open a Pull Request

## Development Setup

See [BUILD.md](BUILD.md) for detailed setup instructions.

Quick start:
```bash
# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Run in development mode
wails dev
```

## Project Structure

- `cmd/` - CLI commands (Cobra)
- `internal/` - Shared Go packages
- `frontend/src/` - React application
- `main.go`, `app.go`, `gui.go` - Application entry points
- `docs/` - Documentation

## Versioning Strategy (Semantic Versioning)

We use [Semantic Versioning](https://semver.org/) (SemVer) for all releases: `MAJOR.MINOR.PATCH`

### Version Bump Rules

| Change Type | Version Bump | Examples |
|-------------|--------------|----------|
| **MAJOR** | Breaking changes | Removing CLI commands, changing output formats, incompatible config changes, breaking API changes |
| **MINOR** | New features (backwards-compatible) | New CLI commands, new GUI views, new output options, new editor support |
| **PATCH** | Bug fixes & improvements | Bug fixes, performance improvements, documentation updates, internal refactors |

### Pre-release Versions

- **Alpha**: `1.1.0-alpha.1` - Early development, unstable
- **Beta**: `1.1.0-beta.1` - Feature complete, testing phase
- **RC**: `1.1.0-rc.1` - Release candidate, final testing

## Conventional Commits

All commits **MUST** follow the [Conventional Commits](https://www.conventionalcommits.org/) specification. This enables automatic version bumping and changelog generation.

### Commit Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Commit Types

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New feature | MINOR |
| `fix` | Bug fix | PATCH |
| `docs` | Documentation only | PATCH |
| `style` | Code style (formatting, semicolons) | PATCH |
| `refactor` | Code refactor (no feature/fix) | PATCH |
| `perf` | Performance improvement | PATCH |
| `test` | Adding/updating tests | PATCH |
| `build` | Build system changes | PATCH |
| `ci` | CI configuration changes | PATCH |
| `chore` | Maintenance tasks | PATCH |

### Breaking Changes

For breaking changes, add `!` after the type or add `BREAKING CHANGE:` in the footer:

```
feat(cli)!: rename --output flag to --format

BREAKING CHANGE: The --output flag has been renamed to --format for consistency.
```

This triggers a **MAJOR** version bump.

### Examples

```bash
# Feature (MINOR bump)
feat(cli): add sync command for multi-editor support

# Bug fix (PATCH bump)
fix(validation): handle extensions with missing metadata

# Breaking change (MAJOR bump)
feat(api)!: change audit output format to include timestamps

# Documentation (PATCH bump)
docs(readme): update installation instructions for v1.0

# Multiple scopes
feat(gui,cli): add marketplace search functionality
```

### Scopes

Common scopes for this project:
- `cli` - CLI commands and interface
- `gui` - Desktop GUI application
- `api` - Backend Go APIs
- `validation` - Extension validation logic
- `marketplace` - Marketplace client
- `sync` - Extension sync functionality
- `editor` - Editor detection and management
- `build` - Build and release process
- `deps` - Dependencies

## Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `go vet` before committing
- Add tests for new functionality
- Keep functions focused and small
- Use meaningful variable names

### TypeScript/React Code

- Use TypeScript for type safety
- Follow React best practices
- Use functional components and hooks
- Keep components focused and reusable
- Add prop types/interfaces

## Testing

### Go Tests

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

### Frontend Tests

```bash
cd frontend
npm test
npm run test:coverage
```

## Pull Request Process

1. Ensure your branch is up to date with `main`
2. Follow Conventional Commits for all commits
3. Update documentation if needed
4. Add tests for new features
5. Ensure all tests pass
6. Update CHANGELOG.md (or let release-please handle it)
7. Request review from maintainers

## Release Process

Releases are automated via GitHub Actions:

1. Commits to `main` are analyzed by release-please
2. A release PR is created/updated automatically
3. When merged, a new release is published with:
   - Git tag (e.g., `v1.1.0`)
   - GitHub Release with binaries
   - Updated CHANGELOG.md

## Areas for Contribution

- **Bug fixes**: Check open issues
- **New features**: Discuss in an issue first
- **Documentation**: Always welcome
- **Tests**: Improve coverage
- **Performance**: Optimize slow operations
- **UI/UX**: Enhance user experience

## Questions?

Open an issue or start a discussion on GitHub.
