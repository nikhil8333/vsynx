# New Features - Marketplace Search & Publisher Verification

## 🎉 What's New

### 1. Marketplace Extension Search
**Search ANY extension** by ID without needing it installed locally!
- Direct marketplace lookup
- Compare with OpenVSX registry
- Get instant validation results

### 2. Publisher Verification
**Know who you trust!**
- ✅ **Verified Publisher Badge** - Shows if publisher is verified by Microsoft
- 🏢 **Publisher Domain** - See the verified domain
- 🛡️ **Trust Indicators** - Visual badges for verified publishers

### 3. SHA256 Digest Comparison
**Binary-level security!**
- Compare SHA256 hashes between marketplaceand OpenVSX
- Detect if binaries are different
- **Critical Alert** if SHA mismatch detected (potential supply chain attack)

## 🚀 How to Use

### Marketplace Search Tab

1. **Click "Marketplace" tab** in the header
2. **Enter Extension ID** (e.g., `ms-python.python`)
3. **Click Search** or press Enter
4. **View Results**:
   - Trust level and recommendation
   - Publisher verification status
   - SHA256 comparison
   - Side-by-side comparison of both registries

### Quick Search Examples

Click these example buttons in the search interface:
- `ms-python.python` - Official Python extension
- `esbenp.prettier-vscode` - Prettier formatter
- `dbaeumer.vscode-eslint` - ESLint integration
- `GitHub.copilot` - GitHub Copilot (marketplace only)

## 📊 New Validation Features

### Publisher Verification

**Verified Publishers Show**:
```
✓ Verified Publisher: Microsoft
  Publisher Domain: microsoft.com
```

**Non-Verified Publishers**:
- No verification badge
- Recommendation includes caution

### SHA256 Comparison

**When Hashes Match** ✅:
```
✓ SHA256 hashes match
Trust Level: Legitimate
```

**When Hashes Don't Match** ⚠️:
```
⚠ SHA256 hash mismatch - binaries are different!
Trust Level: Malicious
Recommendation: DANGER: SHA256 mismatch detected - binaries are DIFFERENT. Potential supply chain attack!
```

## 🔍 What Gets Validated

### For Each Extension:

1. **Metadata Comparison**
   - Publisher name
   - Extension name
   - Version number
   - Repository URL

2. **Publisher Trust**
   - Domain verification
   - Verified publisher flag
   - Publisher reputation

3. **Binary Integrity** (when available)
   - SHA256 hash comparison
   - File size validation
   - Binary authenticity

## 📱 Updated UI

### New "Marketplace" Tab
- Search any extension without installing
- Beautiful results display
- Quick example buttons
- Real-time validation

### Enhanced Results Display
- **Trust Level Cards** with color coding
- **Verified Publisher Badges** (🛡️)
- **SHA Mismatch Alerts** (⚠️)
- **Side-by-side Comparison** of both registries

### Visual Indicators

| Icon | Meaning |
|------|---------|
| 🛡️ | Verified Publisher |
| ✓ | Verified/Matched |
| ⚠️ | Warning/Mismatch |
| 📦 | Microsoft Marketplace |
| 🔓 | OpenVSX Registry |

## 🔒 Security Improvements

### SHA256 Verification
**The most critical new feature!**

If an extension shows different SHA256 hashes:
- **Trust Level: Malicious**
- **Red alert displayed**
- **Clear warning message**
- **DO NOT USE the extension**

This detects:
- Supply chain attacks
- Tampered binaries
- Malicious modifications
- Compromised packages

### Publisher Trust
**Verified publishers are authenticated by Microsoft**

Benefits:
- Known organizations
- Verified domain ownership
- Reduced impersonation risk
- Better trust assurance

## 💡 Use Cases

### 1. Pre-Installation Check
Before installing an extension:
1. Search for it in Marketplace tab
2. Check trust level
3. Verify publisher
4. Confirm SHA matches (if available)
5. **Then** install if legitimate

### 2. Supply Chain Audit
Check popular extensions:
1. Search each one
2. Verify SHA digests match
3. Confirm publisher verification
4. Document any suspicious findings

### 3. Security Investigation
Investigate suspicious extensions:
1. Compare marketplace vs OpenVSX
2. Check for metadata mismatches
3. Verify binary integrity
4. Check publisher authenticity

## 🎯 Example Workflows

### Workflow 1: Validating New Extension

```
1. User wants to install "some-extension"
2. Click "Marketplace" tab
3. Enter: "publisher.some-extension"
4. Click "Search"
5. Results show:
   ✓ Verified Publisher: TrustedPublisher
   ✓ SHA256 hashes match
   Trust Level: Legitimate
6. Safe to install!
```

### Workflow 2: Detecting Malicious Extension

```
1. User searches suspicious extension
2. Results show:
   ⚠ SHA256 hash mismatch - binaries are different!
   Publisher mismatch: TrustedPublisher (marketplace) vs FakePublisher (OpenVSX)
   Trust Level: Malicious
   Recommendation: DANGER - do not use!
3. User avoids security threat!
```

### Workflow 3: Verifying Microsoft Extensions

```
1. Search: "GitHub.copilot"
2. Results:
   🛡️ Verified Publisher: GitHub
   Publisher Domain: github.com
   Extension not found in OpenVSX (expected)
   Trust Level: Legitimate
3. Confirmed authentic Microsoft extension
```

## 🔧 Technical Details

### New Backend Methods

**app.go**:
- `SearchMarketplaceExtension(extensionID string)` - Search without local install

### Updated Models

**ExtensionMetadata**:
- `IsVerifiedPublisher` - Publisher verification status
- `PublisherDomain` - Verified domain
- `SHA256Hash` - Binary hash
- `FileSize` - Package size

**ValidationResult**:
- `SHAMatch` - Boolean for hash comparison
- `SHAMismatchDetails` - Details if hashes differ

### Enhanced Validation Logic

**Priority-based Trust Levels**:
1. **SHA Mismatch** → Malicious (highest priority)
2. **Critical Metadata Mismatch** → Malicious
3. **Minor Differences** → Suspicious
4. **All Match + Verified Publisher** → Legitimate
5. **All Match** → Legitimate

## 📝 Notes

### SHA256 Limitations
- Not all APIs provide SHA256 hashes
- When unavailable, validation uses metadata only
- Message displayed: "SHA256 hash not available from one source"

### Verified Publisher
- Only available from Microsoft Marketplace
- OpenVSX doesn't have verification system
- Microsoft-only extensions won't be in OpenVSX

### Rate Limiting
- Both APIs may rate-limit requests
- Space out searches if needed
- Bulk audits may take time

## 🚀 Getting Started

### Step 1: Restart Wails Dev
```bash
# Stop current server (Ctrl+C)
wails dev
```

This regenerates bindings with the new `SearchMarketplaceExtension` method.

### Step 2: Open Marketplace Tab
Click the **"Marketplace"** button in the header (with 🔍 icon).

### Step 3: Try a Search
Enter an extension ID and click **Search**.

### Step 4: Review Results
- Check trust level
- Look for verified publisher badge
- Review SHA comparison
- Read recommendation

## 🎨 UI Screenshots (Description)

### Marketplace Search View
- Large search input with placeholder
- Quick example buttons
- Search button with loading state
- Results card with trust level color
- Verified publisher badge (if applicable)
- SHA comparison details
- Side-by-side registry comparison

### Results Display
- **Green** - Legitimate extensions
- **Yellow** - Suspicious extensions
- **Red** - Malicious extensions with warnings
- **Shields** - Verified publishers
- **Checkmarks** - Validated items
- **Warning signs** - Mismatches

## 🔐 Best Practices

1. **Always search before installing** new extensions
2. **Check for verified publisher** badge
3. **Verify SHA256 matches** when available
4. **Read all differences** carefully
5. **Never install** extensions with SHA mismatches
6. **Be cautious** with unverified publishers
7. **Report suspicious** extensions to marketplace

## 📚 Additional Documentation

- See `DEBUG.md` for troubleshooting
- See `USAGE.md` for general usage
- See `CHANGES.md` for all recent changes

## 🎯 Summary

You can now:
- ✅ Search marketplace extensions without installing
- ✅ Verify publisher authenticity
- ✅ Compare SHA256 digests
- ✅ Detect supply chain attacks
- ✅ Make informed installation decisions
- ✅ Audit extension security proactively

**Stay secure! 🛡️**
