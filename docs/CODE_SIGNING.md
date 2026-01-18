# Code Signing Setup

This document explains how to set up code signing for Windows and macOS to avoid "Unknown Publisher" warnings.

## Overview

Without code signing, users will see security warnings when running the application:
- **Windows**: "Windows protected your PC" / "Unknown publisher"
- **macOS**: "App can't be opened because it is from an unidentified developer"

## Required GitHub Secrets

Add these secrets to your repository under **Settings → Secrets and variables → Actions**.

### Windows Code Signing

| Secret | Description |
|--------|-------------|
| `WINDOWS_CERTIFICATE` | Base64-encoded `.pfx` code signing certificate |
| `WINDOWS_CERTIFICATE_PASSWORD` | Password for the `.pfx` file |

#### Getting a Windows Code Signing Certificate

1. **Purchase a certificate** from a Certificate Authority:
   - [DigiCert](https://www.digicert.com/signing/code-signing-certificates)
   - [Sectigo](https://sectigo.com/ssl-certificates-tls/code-signing)
   - [GlobalSign](https://www.globalsign.com/en/code-signing-certificate)
   - For best SmartScreen reputation, get an **EV (Extended Validation)** certificate

2. **Export as .pfx**:
   ```powershell
   # If you have separate .cer and .key files:
   openssl pkcs12 -export -out certificate.pfx -inkey private.key -in certificate.cer
   ```

3. **Encode to base64**:
   ```powershell
   [Convert]::ToBase64String([IO.File]::ReadAllBytes("certificate.pfx")) | Set-Clipboard
   ```
   Or on Linux/macOS:
   ```bash
   base64 -i certificate.pfx | pbcopy  # macOS
   base64 certificate.pfx | xclip      # Linux
   ```

4. **Add to GitHub Secrets**:
   - `WINDOWS_CERTIFICATE`: Paste the base64 string
   - `WINDOWS_CERTIFICATE_PASSWORD`: The password you set when exporting

### macOS Code Signing & Notarization

| Secret | Description |
|--------|-------------|
| `APPLE_CERTIFICATE` | Base64-encoded `.p12` Developer ID Application certificate |
| `APPLE_CERTIFICATE_PASSWORD` | Password for the `.p12` file |
| `APPLE_ID` | Your Apple ID email address |
| `APPLE_ID_PASSWORD` | App-specific password (not your Apple ID password) |
| `APPLE_TEAM_ID` | Your Apple Developer Team ID |

#### Prerequisites

1. **Apple Developer Program membership** ($99/year): https://developer.apple.com/programs/

#### Getting a Developer ID Certificate

1. **Create a Certificate Signing Request (CSR)**:
   - Open **Keychain Access** on macOS
   - Go to **Keychain Access → Certificate Assistant → Request a Certificate from a Certificate Authority**
   - Enter your email, select "Saved to disk"

2. **Create the certificate** at https://developer.apple.com/account/resources/certificates:
   - Click "+" to create new certificate
   - Select **"Developer ID Application"**
   - Upload your CSR
   - Download the certificate and double-click to install

3. **Export as .p12**:
   - Open **Keychain Access**
   - Find your "Developer ID Application" certificate
   - Right-click → Export
   - Save as `.p12` with a strong password

4. **Encode to base64**:
   ```bash
   base64 -i Certificates.p12 | pbcopy
   ```

5. **Create an app-specific password**:
   - Go to https://appleid.apple.com/account/manage
   - Under "Sign-In and Security", select "App-Specific Passwords"
   - Generate a new password for "GitHub Actions"

6. **Find your Team ID**:
   - Go to https://developer.apple.com/account
   - Your Team ID is shown in the top right or under Membership

7. **Add to GitHub Secrets**:
   - `APPLE_CERTIFICATE`: Base64 string of your .p12
   - `APPLE_CERTIFICATE_PASSWORD`: The .p12 export password
   - `APPLE_ID`: Your Apple ID email
   - `APPLE_ID_PASSWORD`: The app-specific password you generated
   - `APPLE_TEAM_ID`: Your 10-character Team ID

## Testing

After adding secrets, trigger a release:

```bash
# Test on a branch first
gh workflow run release.yml --ref your-branch -f branch=your-branch

# Or create a tag for a full release
git tag v1.0.3
git push origin v1.0.3
```

## Verification

### Windows
```powershell
# Check signature
signtool verify /pa /v "vsynx-manager-windows-amd64.exe"

# Or right-click the .exe → Properties → Digital Signatures
```

### macOS
```bash
# Check signature
codesign --verify --verbose=4 "Vsynx Manager.app"

# Check notarization
spctl --assess --verbose=4 --type execute "Vsynx Manager.app"

# Check stapled ticket
xcrun stapler validate "Vsynx Manager.app"
```

## Troubleshooting

### Windows: "Certificate not found"
- Ensure the .pfx is correctly base64-encoded (no line breaks)
- Verify the password is correct

### macOS: "Developer ID Application certificate not found"
- Make sure you exported the correct certificate type (Developer ID Application, not Developer ID Installer)
- The certificate must include the private key

### macOS: Notarization fails
- Ensure hardened runtime is enabled (our scripts do this automatically)
- Check that your app-specific password is valid
- Review notarization logs: `xcrun notarytool log <submission-id> --apple-id ... --password ... --team-id ...`

## Cost Summary

| Platform | Certificate Type | Approximate Cost | Validity |
|----------|-----------------|------------------|----------|
| Windows | Standard Code Signing | $200-400/year | 1-3 years |
| Windows | EV Code Signing | $400-700/year | 1-3 years |
| macOS | Apple Developer Program | $99/year | 1 year |

**Note**: EV certificates provide immediate SmartScreen reputation on Windows. Standard certificates may still show warnings until enough users have downloaded your software.
