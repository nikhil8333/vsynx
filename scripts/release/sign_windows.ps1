# Windows Code Signing Script
# Requires: WINDOWS_CERTIFICATE (base64 encoded .pfx) and WINDOWS_CERTIFICATE_PASSWORD secrets

param(
    [Parameter(Mandatory=$true)]
    [string]$FilePath,
    
    [Parameter(Mandatory=$false)]
    [string]$CertificateBase64 = $env:WINDOWS_CERTIFICATE,
    
    [Parameter(Mandatory=$false)]
    [string]$CertificatePassword = $env:WINDOWS_CERTIFICATE_PASSWORD,
    
    [Parameter(Mandatory=$false)]
    [string]$TimestampServer = "http://timestamp.digicert.com"
)

$ErrorActionPreference = "Stop"

# Check if signing is configured
if ([string]::IsNullOrEmpty($CertificateBase64)) {
    Write-Host "WINDOWS_CERTIFICATE not set - skipping code signing"
    Write-Host "To enable signing, add WINDOWS_CERTIFICATE (base64 .pfx) and WINDOWS_CERTIFICATE_PASSWORD secrets"
    exit 0
}

Write-Host "=== Windows Code Signing ==="
Write-Host "Signing: $FilePath"

# Create temp directory for certificate
$tempDir = Join-Path $env:RUNNER_TEMP "codesign"
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null
$certPath = Join-Path $tempDir "certificate.pfx"

try {
    # Decode certificate from base64
    Write-Host "Decoding certificate..."
    [System.IO.File]::WriteAllBytes($certPath, [System.Convert]::FromBase64String($CertificateBase64))
    
    # Import certificate to store
    Write-Host "Importing certificate..."
    $securePassword = ConvertTo-SecureString -String $CertificatePassword -Force -AsPlainText
    $cert = Import-PfxCertificate -FilePath $certPath -CertStoreLocation Cert:\CurrentUser\My -Password $securePassword
    
    Write-Host "Certificate thumbprint: $($cert.Thumbprint)"
    Write-Host "Certificate subject: $($cert.Subject)"
    
    # Find signtool
    $signtoolPaths = @(
        "${env:ProgramFiles(x86)}\Windows Kits\10\bin\10.0.22621.0\x64\signtool.exe",
        "${env:ProgramFiles(x86)}\Windows Kits\10\bin\10.0.19041.0\x64\signtool.exe",
        "${env:ProgramFiles(x86)}\Windows Kits\10\bin\x64\signtool.exe"
    )
    
    $signtool = $null
    foreach ($path in $signtoolPaths) {
        if (Test-Path $path) {
            $signtool = $path
            break
        }
    }
    
    # Fallback: search for signtool
    if (-not $signtool) {
        $signtool = Get-ChildItem -Path "${env:ProgramFiles(x86)}\Windows Kits" -Recurse -Filter "signtool.exe" -ErrorAction SilentlyContinue | 
                    Where-Object { $_.FullName -match "x64" } | 
                    Select-Object -First 1 -ExpandProperty FullName
    }
    
    if (-not $signtool) {
        throw "signtool.exe not found. Install Windows SDK."
    }
    
    Write-Host "Using signtool: $signtool"
    
    # Sign the file
    Write-Host "Signing file..."
    & $signtool sign /fd SHA256 /sha1 $cert.Thumbprint /tr $TimestampServer /td SHA256 /v $FilePath
    
    if ($LASTEXITCODE -ne 0) {
        throw "Signing failed with exit code: $LASTEXITCODE"
    }
    
    # Verify signature
    Write-Host "Verifying signature..."
    & $signtool verify /pa /v $FilePath
    
    Write-Host "=== Signing completed successfully ==="
    
} finally {
    # Cleanup
    if (Test-Path $certPath) {
        Remove-Item -Path $certPath -Force
    }
    
    # Remove certificate from store
    if ($cert) {
        Remove-Item -Path "Cert:\CurrentUser\My\$($cert.Thumbprint)" -Force -ErrorAction SilentlyContinue
    }
}
