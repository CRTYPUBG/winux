# ============================================================================
# WINUX Build & Sign Script
# ============================================================================
# Usage: .\build.ps1 [-Version "0.3.11"]

param(
    [string]$Version = "0.3.11"
)

# Configuration
$GoPath = "go"
$InnoPath = "C:\Program Files (x86)\Inno Setup 6\ISCC.exe"
$ArtifactsDir = "installer"

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  WINUX Build Script v$Version" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

# Step 1: Build winux.exe
Write-Host "[1/3] Building winux.exe..." -ForegroundColor Cyan
$LdFlags = "-s -w -X main.Version=$Version -X github.com/CRTYPUBG/winux/internal/core.Version=$Version"
& $GoPath build -trimpath -ldflags $LdFlags -o winux.exe ./cmd/winux
if ($LASTEXITCODE -ne 0) { throw "Winux build failed" }

# Step 2: Build update.exe
Write-Host "[2/3] Building update.exe..." -ForegroundColor Cyan
$UpdateLdFlags = "-s -w -X main.CurrentVersion=$Version -X github.com/CRTYPUBG/winux/internal/updater.CurrentVersion=$Version"
& $GoPath build -trimpath -ldflags $UpdateLdFlags -o update.exe ./cmd/update
if ($LASTEXITCODE -ne 0) { throw "Update build failed" }

# Step 3: Compile installer
Write-Host "[3/3] Compiling installer..." -ForegroundColor Cyan
& $InnoPath "/DMyAppVersion=$Version" winux.iss
if ($LASTEXITCODE -ne 0) { throw "Installer compilation failed" }

# Generate SHA256
Write-Host ""
Write-Host "Generating checksums..." -ForegroundColor Cyan
$Files = @("winux.exe", "update.exe", "$ArtifactsDir\winux-$Version-setup.exe")
foreach ($f in $Files) {
    if (Test-Path $f) {
        $hash = (Get-FileHash $f -Algorithm SHA256).Hash.ToLower()
        Write-Host "  $f : $hash"
        "$hash  $(Split-Path $f -Leaf)" | Out-File -Append -FilePath "$ArtifactsDir\checksums.sha256" -Encoding ascii
    }
}

Write-Host ""
Write-Host "BUILD COMPLETE!" -ForegroundColor Green
