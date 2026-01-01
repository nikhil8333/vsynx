# Vsynx Manager - Usage Guide

## What This App Does

Vsynx Manager is a security tool that validates VS Code extensions by comparing them against official sources:
- **Microsoft Marketplace**: The official VS Code extension registry
- **OpenVSX**: The open-source alternative registry

The app detects potentially malicious or tampered extensions by finding mismatches in metadata.

## Main Features

### 1. Extension Scanning
- Automatically detects your VS Code extensions directory
- Scans all installed extensions
- Shows extension details (name, version, publisher, path)

### 2. Validation
- Fetches metadata from Microsoft Marketplace
- Fetches metadata from OpenVSX
- Compares: publisher name, extension name, version, repository URL
- Assigns trust level based on findings

### 3. Trust Levels

#### 🟢 Legitimate
- Extension found in both sources
- All metadata matches perfectly
- **Action**: Safe to use

#### 🟡 Suspicious
- Minor metadata differences detected
- Version mismatch between sources
- Only found in one source
- **Action**: Review carefully before using

#### 🔴 Malicious
- Critical metadata mismatches
- Publisher name mismatch
- Repository URL mismatch
- **Action**: Do not use, consider removing

#### ⚪ Unknown
- Extension not found in either source
- API errors occurred
- **Action**: Manual verification required

### 4. Full Audit
- Validates ALL installed extensions at once
- Provides summary statistics
- Shows detailed report for each extension
- Helps identify security risks across your entire setup

### 5. Download Official Version
- Downloads the official VSIX package from Microsoft Marketplace
- Provides SHA256 hash for verification
- Allows you to replace suspicious extensions with official ones

## How to Use

### Initial Setup
1. Run `wails dev` to start the application
2. App automatically detects VS Code extensions directory
3. Extensions list loads automatically

### Change Extensions Path
1. Click "Change Path" button in toolbar
2. Select your VS Code extensions directory:
   - Windows: `%USERPROFILE%\.vscode\extensions`
   - Linux: `~/.vscode/extensions`
   - macOS: `~/.vscode/extensions`

### Validate Single Extension
1. Click on an extension in the list
2. Click "Validate" button
3. Wait for API calls to complete
4. Review trust level and recommendation
5. Check "Differences Found" section for details

### Run Full Audit
1. Click "Audit" button in header
2. Wait for validation of all extensions
3. Review summary cards:
   - Total extensions
   - Legitimate count
   - Suspicious count
   - Malicious count
4. Scroll through detailed results

### Download Official Extension
1. Select a suspicious extension
2. Click "Download Official" button
3. Note the download location (displayed in alert)
4. Install the official version manually in VS Code

## Understanding Results

### What to Look For

#### ✅ Good Signs
- Trust Level: Legitimate
- "Extension is verified - metadata matches across sources"
- No differences found

#### ⚠️ Warning Signs
- Version mismatch (could be due to update timing)
- Minor URL differences (e.g., http vs https)
- Extension only in one source

#### 🚨 Danger Signs
- Publisher name mismatch
- Repository URL points to different location
- Extension claims to be official but metadata doesn't match

### Example Scenarios

#### Scenario 1: Official Extension
```
Extension: ms-python.python
Trust Level: Legitimate
Marketplace: Microsoft (version 2024.0.1)
OpenVSX: Microsoft (version 2024.0.1)
Result: ✅ Safe to use
```

#### Scenario 2: Version Mismatch
```
Extension: publisher.extension
Trust Level: Suspicious
Differences:
- Version mismatch: 1.2.3 (marketplace) vs 1.2.2 (OpenVSX)
Result: ⚠️ Likely safe, just not synced yet
```

#### Scenario 3: Publisher Mismatch
```
Extension: fake-microsoft.python
Trust Level: Malicious
Differences:
- Publisher mismatch: Microsoft (marketplace) vs fake-microsoft (installed)
Result: 🚨 REMOVE IMMEDIATELY - Impersonation attempt
```

## Best Practices

### Security
1. **Run audits regularly**: Check monthly or after installing new extensions
2. **Validate before use**: Check new extensions before enabling
3. **Trust legitimate sources**: Prefer extensions from verified publishers
4. **Keep extensions updated**: Ensure you have latest versions

### Workflow
1. Install extension in VS Code
2. Run Vsynx audit
3. Review any suspicious findings
4. Download official version if needed
5. Remove suspicious extensions

### Interpreting Warnings
- **Version differences**: Usually safe if minor (e.g., 1.2.3 vs 1.2.2)
- **Missing from OpenVSX**: Many Microsoft extensions aren't in OpenVSX
- **Missing from Marketplace**: Could be community/custom extensions

## Limitations

### Known Limitations
1. **Rate Limiting**: APIs may throttle requests for large audits
2. **Network Required**: Requires internet connection for validation
3. **Manual Installation**: Downloaded VSIX must be installed manually
4. **Local Extensions**: Cannot validate extensions not in registries
5. **Timing Issues**: Recently updated extensions may show version mismatches

### What This App Cannot Do
- Does not scan extension source code
- Does not detect runtime malware
- Does not automatically remove extensions
- Does not guarantee 100% security
- Does not validate extension behavior

## FAQ

**Q: Why does my extension show as "Unknown"?**
A: Extension may be custom/private, or APIs may be temporarily unavailable.

**Q: Should I remove all "Suspicious" extensions?**
A: Not necessarily. Review the differences - minor version mismatches are often okay.

**Q: How often should I run audits?**
A: Monthly is good, or whenever you install new extensions.

**Q: Can I trust "Legitimate" extensions completely?**
A: This tool only validates metadata. You should still only install extensions from trusted publishers.

**Q: What if APIs are down?**
A: App will show errors. Try again later or check your internet connection.

**Q: Does this slow down VS Code?**
A: No, this is a separate application and doesn't affect VS Code performance.

## Getting Help

### Check Logs
- Backend logs: Terminal running `wails dev`
- Frontend logs: Browser console (F12)

### Common Issues
See `DEBUG.md` for detailed troubleshooting

### Report Issues
Include:
- Extension ID being validated
- Error messages from console
- Backend logs from terminal
- Network connectivity status
