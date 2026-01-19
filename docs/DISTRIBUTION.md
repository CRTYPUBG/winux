# WINUX – Distribution, Update & Package Management Specification

This document defines the **official distribution**, **update system**, and
**package manager integrations** for the WINUX project.

WINUX is a **native Windows CLI tool written in Go**, designed to provide
Linux-like command behavior on Windows **without runtime dependencies**.

---

## 1. Distribution Philosophy

WINUX follows these core principles:

- Single native binary (`.exe`)
- No runtime dependencies (no Node.js, no .NET, no Java)
- Predictable, secure updates
- Enterprise-ready distribution model
- Fully open-source

Because of this philosophy:

- ❌ npm, NuGet, Maven, RubyGems are **not used**
- ✅ GitHub Releases is the **single source of truth**
- ✅ Windows-native package managers are preferred

---

## 2. Official Distribution Channels

### 2.1 GitHub Releases (Primary)

Repository:
```
https://github.com/CRTYPUBG/winux/releases
```

Each release **must contain**:
```
winux-vX.Y.Z-windows-amd64.exe
winux-vX.Y.Z-windows-amd64.exe.sha256
release-notes.md
```

The SHA256 file contains:
```
<sha256_hash>  winux-vX.Y.Z-windows-amd64.exe
```

---

### 2.2 Windows Package Manager (winget)

#### 2.2.1 Why winget

- Native Windows solution
- Trusted by Microsoft
- No additional tools required
- Ideal for enterprise and power users

#### 2.2.2 Winget Architecture

WINUX **does not host binaries in winget**.

Instead:
- winget pulls from **GitHub Releases**
- Integrity is verified via SHA256

#### 2.2.3 Winget Manifest Structure

Repository:
```
microsoft/winget-pkgs
```

Manifest example:
```yaml
PackageIdentifier: CRTYPUBG.WINUX
PackageVersion: 1.0.0
PackageName: WINUX
Publisher: CRTYPUBG
License: MIT
ShortDescription: Linux-like CLI commands for Windows
Installers:
  - Architecture: x64
    InstallerType: portable
    InstallerUrl: https://github.com/CRTYPUBG/winux/releases/download/v1.0.0/winux-v1.0.0-windows-amd64.exe
    InstallerSha256: <SHA256_HASH>
```

Installation:
```powershell
winget install CRTYPUBG.WINUX
```

Upgrade:
```powershell
winget upgrade CRTYPUBG.WINUX
```

---

### 2.3 Scoop (Developer Distribution)

#### 2.3.1 Why Scoop

- Popular among developers
- Portable-first philosophy
- Easy version pinning

#### 2.3.2 Scoop Bucket Structure

Custom bucket repository:
```
winux-scoop-bucket/
 └── bucket/
     └── winux.json
```

`winux.json` example:
```json
{
  "version": "1.0.0",
  "description": "Linux-like CLI commands for Windows",
  "homepage": "https://github.com/CRTYPUBG/winux",
  "license": "MIT",
  "architecture": {
    "64bit": {
      "url": "https://github.com/CRTYPUBG/winux/releases/download/v1.0.0/winux-v1.0.0-windows-amd64.exe",
      "hash": "<SHA256_HASH>"
    }
  },
  "bin": "winux.exe"
}
```

Installation:
```powershell
scoop bucket add winux https://github.com/CRTYPUBG/winux-scoop-bucket
scoop install winux
```

---

## 3. Update System Architecture (Built-in)

### 3.1 Overview

WINUX includes a self-update system implemented via a dedicated binary:
```
update.exe
```

Design goals:
- Atomic updates
- Secure verification
- No partial installs
- No locked binary issues

### 3.2 Update Flow

```
winux.exe
   |
   |---> GitHub Releases API
   |
   |---> Compare versions
   |
   |---> Download update via aria2
   |
   |---> Verify SHA256
   |
   |---> Replace binary safely
   |
   |---> Restart WINUX
```

### 3.3 GitHub API Usage

API endpoint:
```
https://api.github.com/repos/CRTYPUBG/winux/releases/latest
```

WINUX extracts:
- `tag_name`
- asset download URL
- `.sha256` asset

No authentication required.

### 3.4 Secure Download with aria2

Repository:
```
https://github.com/aria2/aria2
```

Why aria2:
- Multi-connection downloads
- Resume support
- Proven reliability
- Scriptable CLI

aria2 is bundled or dynamically fetched depending on build mode.

Example command:
```powershell
aria2c -x 8 -s 8 -o winux.new.exe <DOWNLOAD_URL>
```

### 3.5 SHA256 Verification

After download:
```
Downloaded hash == Published hash
```

If mismatch:
- Update is aborted
- Binary is deleted
- Error is logged

No update is applied without successful verification.

### 3.6 Safe Binary Replacement Strategy

Windows does not allow replacing a running executable.

Solution:
- `update.exe` performs replacement
- `winux.exe` exits gracefully

Process:
1. `winux.exe` launches `update.exe`
2. `winux.exe` exits
3. `update.exe`:
   - renames old binary
   - moves new binary into place
4. `update.exe` relaunches `winux.exe`
5. `update.exe` exits

This ensures:
- No file locks
- No corruption
- Atomic upgrade

---

## 4. Versioning Policy

WINUX follows **Semantic Versioning**:
```
MAJOR.MINOR.PATCH
```

- **MAJOR**: breaking behavior changes
- **MINOR**: new commands / features
- **PATCH**: bug fixes / performance

---

## 5. Security Model

- All downloads are verified via SHA256
- Single trusted source (GitHub Releases)
- No mirrors
- No auto-executed scripts
- No elevation required

---

## 6. Why Not npm / NuGet / GitHub Packages

WINUX is **not a library**, it is a **native CLI tool**.

Publishing to language-specific package managers would:
- Introduce runtime dependencies
- Break the single-binary principle
- Increase attack surface

This decision is intentional and final.

---

## 7. Future Extensions (Optional)

- Docker image for CI environments
- WINUX SDKs (Go / JS bindings)
- Plugin system (out of scope for v1)

---

## 8. Summary

WINUX uses:
- **GitHub Releases** → source of truth
- **winget** → Windows users
- **Scoop** → developers
- **aria2** → fast, reliable updates
- **SHA256** → integrity
- **update.exe** → atomic upgrades

This architecture is:
- Secure
- Maintainable
- Enterprise-grade
- Fully open-source

---

© CRTYPUBG – WINUX Project
