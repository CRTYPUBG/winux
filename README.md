<p align="center">
  <img src="assets/logo.png" alt="WINUX Logo" width="300">
</p>

<h1 align="center">WINUX</h1>

<p align="center">
  <strong>Native Linux-like command line utilities for Windows.</strong><br>
  High performance · Zero dependencies · Enterprise ready
</p>

<p align="center">
  <a href="https://github.com/CRTYPUBG/winux/releases/latest">
    <img src="https://img.shields.io/github/v/release/CRTYPUBG/winux?style=for-the-badge&color=007ACC" alt="Latest Release">
  </a>
  <a href="https://github.com/CRTYPUBG/winux/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  </a>
  <img src="https://img.shields.io/badge/Platform-Windows%2010%2F11-0078D4?style=for-the-badge&logo=windows" alt="Platform">
</p>

---

## ⚡ Quick Installation

WINUX can be installed via several methods. Choose the one that fits your workflow.

### 1. Official WinGet (Coming Soon)
Once the manifest is merged, you can install with a single command:
```powershell
winget install CRTYPUBG.WINUX
```

### 2. Windows Installer (.exe)
Download the **v0.3.11-setup.exe** for a guided installation experience:
👉 **[Download Installer](https://github.com/CRTYPUBG/winux/releases/download/v0.3.11/winux-0.3.11-setup.exe)**

### 3. Portable Archives
Download and extract to your custom folder:
- [📦 ZIP Archive](https://github.com/CRTYPUBG/winux/releases/download/v0.3.11/winux-v0.3.11-windows-amd64.zip)
- [🗜️ 7-Zip Archive](https://github.com/CRTYPUBG/winux/releases/download/v0.3.11/winux-v0.3.11-windows-amd64.7z)

---

## ✨ Features

- 🚀 **Native performance** — no WSL, no emulation, no runtime overhead.
- 📦 **BusyBox-style** — a single binary that contains all commands.
- 🔄 **Auto-update system** — stay up to date with `update --check`.
- 🛡️ **Integrity verified** — all files are SHA256 hashed and verifiable.
- 🔗 **Pipe & Redirection** — full support for standard streams.

---

## 🛠️ Available Commands

| Command | Status | Description |
|:---:|:---:|---|
| `ls` | ✅ | List directory contents |
| `cat` | ✅ | Concatenate and print files |
| `grep` | ✅ | Search for patterns in files |
| `rm` | ✅ | Remove files or directories |
| `mkdir` | ✅ | Create directories |
| `touch` | ✅ | Create empty files or update timestamps |
| `pwd` | ✅ | Print working directory |
| `echo` | ✅ | Display text/variables |
| `whoami`| ✅ | Print effective username |
| `uptime`| ✅ | Display system uptime |
| `update`| ✅ | Self-updater utility |

---

## 🔄 Self-Update Utility

WINUX comes with a built-in update manager.

```powershell
# Check for latest version
update --check

# Apply latest update automatically
update --apply

# Check version in background on startup
update --startup
```

---

## 🗺️ Project Status

- [x] **v0.1** — Core logic & Dispatcher.
- [x] **v0.2** — Added basic file commands.
- [x] **v0.3** — Added Update system & Installer.
- [ ] **v0.4** — Recursive operations (`rm -rf`, `ls -R`).
- [ ] **v1.0** — Official WinGet release & Complete suite.

---

## 📜 License

Distributed under the **MIT License**. See `LICENSE` for more information.

---

<p align="center">
  <em>"Reclaiming the Windows CLI, one command at a time."</em>
</p>
