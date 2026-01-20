# WINUX Wiki

> Native Linux-like command-line utilities for Windows

---

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Commands Reference](#commands-reference)
4. [Usage Examples](#usage-examples)
5. [Piping & Redirection](#piping--redirection)
6. [Self-Update](#self-update)
7. [Building from Source](#building-from-source)
8. [Architecture](#architecture)
9. [FAQ](#faq)

---

## Introduction

WINUX provides native Linux-style command-line utilities on Windows without:
- ❌ WSL (Windows Subsystem for Linux)
- ❌ Emulation layers
- ❌ PowerShell aliases
- ❌ Shell wrappers

### Key Features

- **Single binary** — One EXE, all commands
- **Native performance** — Pure Go, no runtime dependencies
- **Linux-compatible** — Same flags, same behavior
- **Pipe support** — Full STDIN/STDOUT/STDERR
- **Exit codes** — Linux-compatible (0, 1, 2, 127)

---

## Installation

### Option 1: Download Binary

1. Download `winux.exe` from [Releases](https://github.com/CRTYPUBG/winux/releases)
2. Place in `C:\Windows\System32\` or add to PATH

### Option 2: Installer

1. Download `winux-x.x.x-setup.exe`
2. Run installer
3. Select "Add to PATH"

### Option 3: Portable

Just download `winux.exe` and run from any folder.

### Verify Installation

```powershell
winux --version
winux --help
```

---

## Commands Reference

### ls — List Directory Contents

```
Usage: ls [OPTION]... [FILE]...

Options:
  -a, --all    Show hidden files (starting with .)
  -l           Long listing format
  -h           Human-readable sizes (with -l)
  --help       Display help
```

**Examples:**
```powershell
winux ls
winux ls -la
winux ls -lah C:\Users
```

---

### cat — Concatenate Files

```
Usage: cat [OPTION]... [FILE]...

Options:
  -n, --number             Number all lines
  -b, --number-nonblank    Number non-empty lines
  --help                   Display help
```

**Examples:**
```powershell
winux cat file.txt
winux cat -n file.txt
winux cat file1.txt file2.txt
```

---

### grep — Search Patterns

```
Usage: grep [OPTION]... PATTERN [FILE]...

Options:
  -i, --ignore-case         Case insensitive
  -v, --invert-match        Select non-matching lines
  -n, --line-number         Show line numbers
  -c, --count               Count matches only
  -l, --files-with-matches  Show only filenames
  -E, --extended-regexp     Extended regex
  -e PATTERN                Pattern to match
  --help                    Display help
```

**Examples:**
```powershell
winux grep error log.txt
winux grep -i ERROR log.txt
winux grep -n "pattern" *.txt
winux grep -c TODO *.go
```

---

### rm — Remove Files

```
Usage: rm [OPTION]... FILE...

Options:
  -r, -R, --recursive    Remove directories recursively
  -f, --force            Ignore errors, never prompt
  -v, --verbose          Explain what is being done
  --help                 Display help
```

**Examples:**
```powershell
winux rm file.txt
winux rm -f file.txt
winux rm -rf directory/
winux rm -v file1.txt file2.txt
```

⚠️ **Warning:** `rm -rf` is permanent. Use with caution.

---

### mkdir — Create Directories

```
Usage: mkdir [OPTION]... DIRECTORY...

Options:
  -p, --parents    Create parent directories as needed
  -v, --verbose    Print message for each directory
  --help           Display help
```

**Examples:**
```powershell
winux mkdir newdir
winux mkdir -p path/to/deep/dir
winux mkdir -v dir1 dir2 dir3
```

---

### touch — Create/Update Files

```
Usage: touch [OPTION]... FILE...

Options:
  -c, --no-create    Do not create new files
  --help             Display help
```

**Examples:**
```powershell
winux touch newfile.txt
winux touch -c existing.txt
winux touch file1.txt file2.txt
```

---

### pwd — Print Working Directory

```
Usage: pwd

Options:
  --help    Display help
```

**Examples:**
```powershell
winux pwd
```

---

### echo — Display Text

```
Usage: echo [OPTION]... [STRING]...

Options:
  -n    No trailing newline
  -e    Interpret escape sequences
  -E    Disable escape interpretation (default)
```

**Escape sequences (with -e):**
- `\n` — newline
- `\t` — tab
- `\r` — carriage return
- `\\` — backslash

**Examples:**
```powershell
winux echo Hello World
winux echo -n "no newline"
winux echo -e "line1\nline2"
```

---

## Usage Examples

### Basic File Operations

```powershell
# List files
winux ls -la

# View file
winux cat readme.txt

# Search in file
winux grep TODO main.go

# Create directory
winux mkdir -p src/components

# Create file
winux touch src/index.js

# Remove file
winux rm -f temp.txt

# Remove directory
winux rm -rf build/
```

### Combined with Windows Commands

```powershell
# List and filter
winux ls | findstr ".go"

# Count files
winux ls | find /c /v ""

# Save output
winux cat log.txt > backup.txt
```

---

## Piping & Redirection

WINUX fully supports Windows piping:

### Pipe Input

```powershell
# Using type (Windows cat)
type file.txt | winux grep error

# Using echo
echo test | winux cat -n
```

### Pipe Output

```powershell
# Filter output
winux ls | findstr ".exe"

# Chain commands
winux cat file.txt | winux grep pattern
```

### Redirection

```powershell
# Output to file
winux ls > files.txt

# Append to file
winux ls >> files.txt

# Errors to file
winux cat missing.txt 2> errors.txt
```

---

## Self-Update

WINUX includes a built-in update utility.

### Check for Updates

```powershell
update.exe --check
```

### Apply Update

```powershell
update.exe --apply
```

### Force Reinstall

```powershell
update.exe --force
```

### Update Process

1. Checks GitHub Releases API
2. Compares versions
3. Downloads new binary
4. Verifies SHA256 checksum
5. Backs up current version
6. Installs new version
7. Removes backup

---

## Building from Source

### Requirements

- Go 1.22+
- Windows 10/11

### Clone & Build

```powershell
git clone https://github.com/CRTYPUBG/winux.git
cd winux

# Build main binary
go build -trimpath -ldflags="-s -w" -o winux.exe ./cmd/winux

# Build updater
go build -trimpath -ldflags="-s -w" -o update.exe ./cmd/update
```

### Using Build Script

```powershell
.\build.ps1 -Version "0.2.0"
```

This script:
1. Builds winux.exe
2. Builds update.exe
3. Signs both binaries
4. Compiles installer
5. Signs installer
6. Generates SHA256 checksums

---

## Architecture

### BusyBox-Style Dispatch

WINUX uses a single-binary dispatcher model:

```
winux.exe
    ├── argv[0] = "winux" → parse first arg as command
    ├── argv[0] = "ls"    → execute ls directly
    ├── argv[0] = "grep"  → execute grep directly
    └── ...
```

### Directory Structure

```
winux/
├── cmd/
│   ├── winux/main.go       # Entry point
│   └── update/main.go      # Updater
├── internal/
│   ├── commands/           # Command implementations
│   │   ├── cat.go
│   │   ├── echo.go
│   │   ├── grep.go
│   │   ├── ls.go
│   │   ├── mkdir.go
│   │   ├── pwd.go
│   │   ├── rm.go
│   │   └── touch.go
│   ├── core/
│   │   └── dispatcher.go   # Command router
│   ├── io/
│   │   └── stdin.go        # Pipe detection
│   ├── protection/
│   │   └── antidebug.go    # Security
│   └── utils/
│       └── exitcodes.go    # Exit codes
└── go.mod
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General failure (no matches) |
| 2 | Usage error / permission denied |
| 127 | Command not found |

---

## FAQ

### Q: Why not just use WSL?

WSL requires Hyper-V, uses more resources, and has filesystem overhead. WINUX is a native Windows binary with zero dependencies.

### Q: Is this a full Linux replacement?

No. WINUX provides the most common file utilities. For full Linux compatibility, use WSL.

### Q: Are there more commands coming?

Yes! Planned commands:
- `cp` — Copy files
- `mv` — Move files
- `head` / `tail` — View file parts
- `wc` — Word count
- `find` — Search files
- `xargs` — Build commands

### Q: Can I use symlinks for BusyBox-style?

Yes! Create symlinks:
```powershell
# PowerShell (Admin)
New-Item -ItemType SymbolicLink -Path "C:\bin\ls.exe" -Target "C:\bin\winux.exe"
New-Item -ItemType SymbolicLink -Path "C:\bin\grep.exe" -Target "C:\bin\winux.exe"
```

### Q: How do I report bugs?

Open an issue: https://github.com/CRTYPUBG/winux/issues

---

## Links

- **Repository:** https://github.com/CRTYPUBG/winux
- **Releases:** https://github.com/CRTYPUBG/winux/releases
- **Issues:** https://github.com/CRTYPUBG/winux/issues

---

*© 2026 CRTYPUBG — MIT License*
