<p align="center">
  <img src="assets/logo.png" alt="WINUX Logo" width="400">
</p>

<h1 align="center">WINUX</h1>

<p align="center">
  <strong>Native Linux-like coreutils for Windows</strong><br>
  Single binary · No WSL · No aliases
</p>

<p align="center">
  <img src="https://img.shields.io/badge/platform-Windows%2010%2F11-blue?style=flat-square" alt="Platform">
  <img src="https://img.shields.io/badge/language-Go-00ADD8?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/license-MIT-green?style=flat-square" alt="License">
  <img src="https://img.shields.io/github/v/release/CRTYPUBG/winux?style=flat-square" alt="Release">
</p>

---

## ✨ Features

- ✅ **Single static binary** — one executable, no dependencies
- ✅ **Written in Go** — fast compilation, cross-platform potential
- ✅ **Native Windows executable** — no emulation layer
- ✅ **No WSL required** — works on any Windows 10/11
- ✅ **No aliases or shell wrappers** — real executables
- ✅ **Real STDIN / STDOUT / STDERR** — proper stream handling
- ✅ **Pipe and redirection support** — `type file.txt | winux grep error`
- ✅ **Linux-compatible exit codes** — scripts work as expected
- ✅ **BusyBox-style dispatch** — `argv[0]` command resolution

---

## 🚀 Quick Start

### Download

Download the latest release from [Releases](https://github.com/CRTYPUBG/winux/releases).

### Usage

```powershell
# Basic commands
winux ls
winux ls -la
winux cat file.txt
winux grep error log.txt

# Pipe support
type log.txt | winux grep -i error
winux cat file.txt | winux grep pattern
```

### BusyBox-style (symlink)

```powershell
# Rename or symlink winux.exe to command name
copy winux.exe ls.exe
.\ls.exe -la
```

---

## 📦 Available Commands

| Command | Description | Flags |
|---------|-------------|-------|
| `ls` | List directory contents | `-a`, `-l`, `-h` |
| `cat` | Concatenate and print files | `-n`, `-b` |
| `grep` | Search for patterns | `-i`, `-v`, `-n`, `-c`, `-l`, `-E` |

*More commands coming in future releases.*

---

## 🏗️ Build from Source

### Requirements

- Go 1.21+

### Build

```powershell
go build -ldflags="-s -w" -o winux.exe ./cmd/winux
```

---

## 📁 Project Structure

```
winux/
├── cmd/winux/main.go          # Entry point & dispatcher
├── internal/
│   ├── commands/              # Command implementations
│   │   ├── cat.go
│   │   ├── grep.go
│   │   └── ls.go
│   ├── core/dispatcher.go     # BusyBox-style command dispatch
│   ├── io/stdin.go            # Pipe detection
│   └── utils/exitcodes.go     # Linux exit codes
├── assets/                    # Branding assets
└── go.mod
```

---

## 🔌 Exit Codes

| Condition | Exit Code |
|-----------|-----------|
| Success | `0` |
| No matches / failure | `1` |
| Invalid usage / error | `2` |
| Command not found | `127` |

---

## 🗺️ Roadmap

- [x] v0.1 — Core commands (`ls`, `cat`, `grep`)
- [ ] v0.2 — More commands (`rm`, `mkdir`, `touch`, `pwd`, `echo`)
- [ ] v0.3 — POSIX-style flags, recursive operations
- [ ] v1.0 — Full coreutils suite, installer, PATH integration

---

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Follow Go best practices
4. Keep behavior Linux-compatible
5. Submit a pull request

---

## 📜 License

MIT License © 2026 CRTYPUBG

---

<p align="center">
  <em>"Linux tools should feel native on Windows, not emulated."</em>
</p>
