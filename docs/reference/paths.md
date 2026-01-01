# Common VS Code Extensions Paths

If the app cannot auto-detect your VS Code extensions directory, use the "Change Path" button and navigate to one of these locations:

## WSL (Windows Subsystem for Linux)

### Option 1: WSL Extensions (Remote-WSL)
```bash
~/.vscode-server/extensions
```
**Full path example**: `/home/nikhil/.vscode-server/extensions`

### Option 2: Windows Extensions (from WSL)
```bash
/mnt/c/Users/YOUR_USERNAME/.vscode/extensions
```
**Example**: `/mnt/c/Users/nikhil/.vscode/extensions`

### Option 3: Windows AppData
```bash
/mnt/c/Users/YOUR_USERNAME/AppData/Roaming/Code/User/extensions
```

## Linux (Native)

### Standard Location
```bash
~/.vscode/extensions
```

### Alternative (Snap install)
```bash
~/snap/code/current/.config/Code/User/extensions
```

### System-wide
```bash
/usr/share/code/resources/app/extensions
```

### VS Code OSS
```bash
~/.vscode-oss/extensions
```

### VSCodium
```bash
~/.vscodium/extensions
```

## Windows (Native)

### User Extensions
```
%USERPROFILE%\.vscode\extensions
```
**Example**: `C:\Users\nikhil\.vscode\extensions`

### AppData Location
```
%APPDATA%\Code\User\extensions
```
**Example**: `C:\Users\nikhil\AppData\Roaming\Code\User\extensions`

### Local AppData
```
%LOCALAPPDATA%\Programs\Microsoft VS Code\resources\app\extensions
```

## macOS

### User Extensions
```bash
~/.vscode/extensions
```

### Application Support
```bash
~/Library/Application Support/Code/User/extensions
```

## How to Find Your Extensions Path

### Method 1: VS Code Command
1. Open VS Code
2. Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on macOS)
3. Type "Extensions: Show Installed Extensions"
4. Right-click any extension
5. Select "Copy Extension ID"
6. Check your file system for that extension

### Method 2: VS Code Settings
1. Open VS Code
2. Go to Settings (`Ctrl+,`)
3. Search for "extensions"
4. Look for extension installation path

### Method 3: Terminal/Command Line

#### Linux/macOS:
```bash
# Find all .vscode* directories
find ~ -name ".vscode*" -type d 2>/dev/null

# Check common locations
ls -la ~/.vscode/extensions
ls -la ~/.vscode-server/extensions
ls -la ~/.config/Code/User/extensions
```

#### Windows (PowerShell):
```powershell
# Check common locations
Get-ChildItem -Path "$env:USERPROFILE\.vscode\extensions" -ErrorAction SilentlyContinue
Get-ChildItem -Path "$env:APPDATA\Code\User\extensions" -ErrorAction SilentlyContinue
```

#### WSL:
```bash
# Check WSL locations
ls -la ~/.vscode-server/extensions

# Check Windows locations from WSL
ls -la /mnt/c/Users/$USER/.vscode/extensions
```

## Troubleshooting

### "Directory not found" error
- Verify VS Code is installed
- Check if you're using VS Code Insiders (path might be `.vscode-insiders`)
- Try alternative locations listed above

### Empty directory
- Extensions might not be installed yet
- Check if you're looking at the right VS Code variant (Code vs Code-OSS vs VSCodium)
- Some extensions might be in a different location

### Permission issues
- On Linux/macOS, ensure you have read permissions: `ls -la ~/.vscode/extensions`
- On Windows, run app as your user (not administrator)

### WSL Specific
- If using Remote-WSL extension, extensions are in `~/.vscode-server/extensions`
- Windows native extensions are in `/mnt/c/Users/USERNAME/.vscode/extensions`
- They are separate installations

## Quick Check Command

Run this in your terminal to find extension directories:

### Linux/WSL:
```bash
find ~ /mnt/c/Users/$USER -name "extensions" -type d 2>/dev/null | grep -E "(vscode|Code)"
```

### macOS:
```bash
find ~ -name "extensions" -type d 2>/dev/null | grep -E "vscode|Code"
```

### Windows (PowerShell):
```powershell
Get-ChildItem -Path $env:USERPROFILE -Recurse -Directory -Filter "extensions" -ErrorAction SilentlyContinue | Where-Object {$_.FullName -match "vscode|Code"}
```

## Current User

To find your username:

### Linux/WSL:
```bash
whoami
echo $USER
```

### Windows:
```cmd
echo %USERNAME%
```

**Your username**: Based on the logs, you're user `nikhil`

## Recommended Path for You (nikhil@hostofnight)

Since you're on WSL, try these in order:

1. **VSCode Server** (most likely for Remote-WSL):
   ```
   /home/nikhil/.vscode-server/extensions
   ```

2. **Windows Extensions from WSL**:
   ```
   /mnt/c/Users/nikhil/.vscode/extensions
   ```

3. **Local WSL Extensions**:
   ```
   /home/nikhil/.vscode/extensions
   ```

## Verification

To verify the correct path, it should contain directories like:
- `publisher.extension-name-version/`
- Example: `ms-python.python-2024.0.1/`

Each directory should have:
- `package.json` file
- Extension code and resources
