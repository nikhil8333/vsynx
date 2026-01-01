# Testing Vsynx CLI During Development

During development (`wails dev`), the CLI binary doesn't exist yet. Use `go run` to test commands:

## Basic Commands

```bash
# List all supported editors
go run . editors list

# Check editor status
go run . editors status
go run . editors status vscode
go run . editors status windsurf

# List extensions for an editor
go run . editors extensions vscode
```

## Marketplace Commands

```bash
# Search marketplace
go run . marketplace search python

# Get marketplace URL for an extension
go run . marketplace open ms-python.python
```

## Validation & Audit

```bash
# Validate a single extension
go run . validate ms-python.python

# List installed extensions
go run . list --path "C:\Users\reddy\.vscode\extensions"

# Audit all extensions
go run . audit --path "C:\Users\reddy\.vscode\extensions"
```

## Sync Commands

```bash
# Preview sync operation
go run . sync preview --from vscode --to windsurf --all

# Check for conflicts
go run . sync conflicts --from vscode --to windsurf --all

# Actually sync extensions
go run . sync run --from vscode --to windsurf --all
go run . sync run --from vscode --to windsurf --ext "ms-python.python,github.copilot"

# Overwrite conflicts
go run . sync run --from vscode --to cursor --all --overwrite
```

## Install Commands

```bash
# Install extension via VS Code CLI
go run . install ms-python.python
go run . install ms-python.python github.copilot

# Install from a file
go run . install --file extensions.txt

# Install to specific editor
go run . install ms-python.python --editor vscode-insiders
```

## Output Formats

Most commands support JSON output:

```bash
go run . editors list --output json
go run . validate ms-python.python --output json
go run . sync preview --from vscode --to windsurf --all --output json
```

## Production Build

To create the actual CLI binary for installation:

```bash
# Build for production
wails build

# The binaries will be in build/bin/
# On Windows: vsynx-manager.exe (GUI) and vsynx.exe (CLI)
```

## Notes

- All `go run .` commands work exactly like the final CLI
- Use this for development, testing, and scripting
- CLI installation is only available in production builds
