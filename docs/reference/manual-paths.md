# Manual Path Entry Guide

## New Feature: Manual Path Input

You can now manually enter the path to your VS Code extensions directory!

## How to Use

### Step 1: Click "Edit" Button
In the toolbar, click the **"Edit"** button next to the path display.

### Step 2: Enter Your Path
An input field will appear. You can:

**Option A: Type manually**
```
~/.vscode-server/extensions
```

**Option B: Use quick path buttons**
Click one of the suggested paths:
- `~/.vscode-server/extensions` (Remote-WSL)
- `~/.vscode/extensions` (Local)
- `/mnt/c/Users/YOUR_USERNAME/.vscode/extensions` (Windows from WSL)

**Option C: Copy-paste**
Find your path using terminal and paste it:
```bash
ls -la ~/.vscode-server/extensions
# If this works, copy the path: ~/.vscode-server/extensions
```

### Step 3: Load Extensions
- Press **Enter** key, or
- Click **"Load"** button

The app will load extensions from your specified path.

## Supported Path Formats

### ✅ Tilde Notation
```
~/.vscode-server/extensions
~/.vscode/extensions
~/.config/Code/User/extensions
```
The `~` will automatically expand to your home directory.

### ✅ Absolute Paths
```
/home/nikhil/.vscode-server/extensions
/mnt/c/Users/nikhil/.vscode/extensions
```

### ✅ WSL Windows Paths
```
/mnt/c/Users/YOUR_USERNAME/.vscode/extensions
/mnt/d/VSCode/extensions
```

## Quick Reference: Finding Your Path

### For WSL Remote Development
```bash
# Most likely location
ls -la ~/.vscode-server/extensions

# If found, use this path:
~/.vscode-server/extensions
```

### For Local WSL VSCode
```bash
# Check local WSL installation
ls -la ~/.vscode/extensions

# If found, use:
~/.vscode/extensions
```

### For Windows VSCode (from WSL)
```bash
# Replace 'nikhil' with your Windows username
ls -la /mnt/c/Users/nikhil/.vscode/extensions

# If found, use:
/mnt/c/Users/YOUR_USERNAME/.vscode/extensions
```

### Find Your Windows Username
```bash
# From WSL
echo $USER

# Or list Windows users
ls -la /mnt/c/Users/
```

## Examples

### Example 1: Remote-WSL Setup
```
User: nikhil
Setup: VS Code on Windows, connected to WSL via Remote-WSL
Path: ~/.vscode-server/extensions
Full path: /home/nikhil/.vscode-server/extensions
```

**Steps:**
1. Click "Edit"
2. Click "~/.vscode-server/extensions" button
3. Click "Load"

### Example 2: Native WSL VSCode
```
User: nikhil
Setup: VS Code installed in WSL
Path: ~/.vscode/extensions
Full path: /home/nikhil/.vscode/extensions
```

**Steps:**
1. Click "Edit"
2. Click "~/.vscode/extensions" button
3. Click "Load"

### Example 3: Windows VSCode from WSL
```
User: nikhil
Setup: VS Code on Windows, accessing from WSL
Windows Username: nikhil
Path: /mnt/c/Users/nikhil/.vscode/extensions
```

**Steps:**
1. Click "Edit"
2. Click "Windows (WSL)" button
3. Edit "YOUR_USERNAME" to "nikhil"
4. Final: `/mnt/c/Users/nikhil/.vscode/extensions`
5. Click "Load"

## Keyboard Shortcuts

When the input field is active:
- **Enter** - Save and load extensions
- **Escape** - Cancel editing

## Troubleshooting

### "No extensions found"
- Verify the path exists: `ls -la YOUR_PATH`
- Check for typos in the path
- Ensure you have the right username

### "Permission denied"
- Check directory permissions: `ls -la YOUR_PATH`
- Ensure you can read the directory

### Path with spaces
If your path has spaces, you don't need quotes in the input field:
```
✅ Correct: /mnt/c/My Folder/extensions
❌ Wrong: "/mnt/c/My Folder/extensions"
```

### Tilde not expanding
If `~` doesn't work, use the full absolute path:
```bash
# Find your home directory
echo $HOME

# Use the full path
/home/nikhil/.vscode-server/extensions
```

## Visual Guide

### Before Editing
```
[Path: ~/.vscode-server/extensions] [Edit] [Browse]
```

### While Editing
```
[Input field: ~/.vscode-server/extensions      ] [Load] [Cancel]
Common paths: [~/.vscode-server/extensions] [~/.vscode/extensions] [Windows (WSL)]
```

### After Loading
```
[Path: ~/.vscode-server/extensions] [Edit] [Browse]
Showing: 42 extensions
```

## Pro Tips

### Tip 1: Try Auto-Detection First
Before manually entering, restart the app to see if improved auto-detection works:
```bash
# Restart app
wails dev
```

### Tip 2: Check Logs
Watch terminal logs to see which paths are being checked:
```
[Scanner] Home directory: /home/nikhil
[Scanner] Checking path 1: /home/nikhil/.vscode/extensions
[Scanner] Checking path 2: /home/nikhil/.vscode-server/extensions
[Scanner] Found extensions directory: /home/nikhil/.vscode-server/extensions
```

### Tip 3: Use Tab Completion in Terminal
When finding your path:
```bash
# Type and press Tab to auto-complete
ls ~/.vscode<TAB>
# Might show: .vscode/  .vscode-server/
```

### Tip 4: Save Your Path
Once you find the correct path, save it for reference:
```bash
# Create a note
echo "~/.vscode-server/extensions" > ~/my-vscode-path.txt
```

## Common Paths Cheatsheet

| Setup | Path |
|-------|------|
| Remote-WSL | `~/.vscode-server/extensions` |
| WSL Native | `~/.vscode/extensions` |
| Windows (from WSL) | `/mnt/c/Users/USERNAME/.vscode/extensions` |
| VSCodium | `~/.vscodium/extensions` |
| VS Code OSS | `~/.vscode-oss/extensions` |

## Need More Help?

1. Check `COMMON_PATHS.md` for detailed path information
2. Check `DEBUG.md` for troubleshooting steps
3. Look at terminal logs for error messages
4. Try the "Browse" button if manual entry doesn't work
