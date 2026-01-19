# ============================================================================
# WINUX Build & Sign Script
# ============================================================================
# 
# Usage: .\build.ps1 [-Version "0.1.0"]
#
# This script:
# 1. Builds winux.exe and update.exe
# 2. Signs both binaries
# 3. Compiles the installer
# 4. Signs the installer
#
# ============================================================================

param(
    [string]$Version = "0.1.0"
)

# Configuration
$GoPath = "C:\Program Files\Go\bin\go.exe"
$SignToolPath = "C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x64\signtool.exe"
$InnoPath = "C:\Program Files (x86)\Inno Setup 6\ISCC.exe"
$CertPath = "C:\Users\LenovoPC\cert.pfx"
$CertPass = "ueo586_crty555"
$TimestampURL = "http://timestamp.digicert.com"

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  WINUX Build & Sign Script v$Version" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# Function to sign a file
function Sign-File {
    param([string]$FilePath)
    
    Write-Host "  Signing: $FilePath" -ForegroundColor Yellow
    & $SignToolPath sign /f $CertPath /p $CertPass /fd SHA256 /tr $TimestampURL /td SHA256 /q $FilePath
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  ERROR: Signing failed!" -ForegroundColor Red
        exit 1
    }
    Write-Host "  ✓ Signed successfully" -ForegroundColor Green
}

# Step 1: Build winux.exe
Write-Host "[1/6] Building winux.exe..." -ForegroundColor Cyan
$BuildTime = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
& $GoPath build -trimpath -ldflags="-s -w -X main.Version=v$Version -X main.BuildTime=$BuildTime" -o winux.exe ./cmd/winux
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  ✓ winux.exe built" -ForegroundColor Green

# Step 2: Build update.exe
Write-Host "[2/6] Building update.exe..." -ForegroundColor Cyan
& $GoPath build -trimpath -ldflags="-s -w" -o update.exe ./cmd/update
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  ✓ update.exe built" -ForegroundColor Green

# Step 3: Sign winux.exe
Write-Host "[3/6] Signing winux.exe..." -ForegroundColor Cyan
Sign-File "winux.exe"

# Step 4: Sign update.exe
Write-Host "[4/6] Signing update.exe..." -ForegroundColor Cyan
Sign-File "update.exe"

# Step 5: Compile installer
Write-Host "[5/6] Compiling installer..." -ForegroundColor Cyan
& $InnoPath "winux.iss"
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Installer compilation failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  ✓ Installer compiled" -ForegroundColor Green

# Step 6: Sign installer
Write-Host "[6/6] Signing installer..." -ForegroundColor Cyan
Sign-File "installer\winux-$Version-setup.exe"

# Generate SHA256
Write-Host ""
Write-Host "Generating SHA256 checksums..." -ForegroundColor Cyan
$files = @("winux.exe", "update.exe", "installer\winux-$Version-setup.exe")
foreach ($file in $files) {
    $hash = (Get-FileHash $file -Algorithm SHA256).Hash
    "$hash  $file" | Out-File -Append -FilePath "installer\checksums.sha256" -Encoding UTF8
    Write-Host "  $file : $hash" -ForegroundColor Gray
}

# Summary
Write-Host ""
Write-Host "============================================" -ForegroundColor Green
Write-Host "  BUILD COMPLETE!" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Green
Write-Host ""
Write-Host "Output files:" -ForegroundColor White
Write-Host "  • winux.exe                    ($('{0:N2}' -f ((Get-Item winux.exe).Length / 1MB)) MB)" -ForegroundColor Gray
Write-Host "  • update.exe                   ($('{0:N2}' -f ((Get-Item update.exe).Length / 1MB)) MB)" -ForegroundColor Gray
Write-Host "  • installer\winux-$Version-setup.exe ($('{0:N2}' -f ((Get-Item "installer\winux-$Version-setup.exe").Length / 1MB)) MB)" -ForegroundColor Gray
Write-Host ""
Write-Host "All files are signed with CRTYPUBG certificate." -ForegroundColor Green
Write-Host ""
