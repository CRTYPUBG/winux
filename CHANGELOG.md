# Changelog

All notable changes to WINUX will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.3.0] - 2026-01-20

### Added
- **Startup Update Checker:** `update.exe --startup` for delayed background checking
- **Windows Popup Notification:** Native dialog with Yes/No/Cancel buttons
- **Changelog Parser:** Shows first 5 lines of release notes
- **"More Info" Button:** Opens GitHub release page in browser
- **Dynamic Version:** All binaries use build-time ldflags for version

### Changed
- `CurrentVersion` is now set via ldflags instead of hardcoded
- Update now auto-detects Program Files installation path
- Better error messages for admin permission requirements

### Fixed
- Update not applying to Program Files installation
- Version mismatch between binaries

---

## [0.2.0] - 2026-01-20

### Added
- `rm` — Remove files and directories
- `mkdir` — Create directories with `-p` flag
- `touch` — Create files or update timestamps
- `pwd` — Print working directory
- `echo` — Display text with escape sequences

---

## [0.1.0] - 2026-01-20

### Added
- **Core Commands:** `ls`, `cat`, `grep` with Linux-compatible flags
- **BusyBox-style dispatcher:** Single binary, multiple commands via argv[0]
- **Pipe support:** Full STDIN/STDOUT/STDERR handling
- **Linux exit codes:** 0, 1, 2, 127 compatibility
- **Self-update system:** `update.exe` with GitHub Releases integration
- **SHA256 verification:** Secure update downloads
- **Anti-debug protection:** Detects debuggers and analysis tools
- **Code signing:** All binaries signed with CRTYPUBG certificate
- **SLSA Level 3 pipeline:** GitHub Actions release workflow
- **Inno Setup installer:** Modern Windows 11 style UI
- **PATH integration:** Automatic system PATH update
- **Build script:** `build.ps1` for automated build & sign

### Security
- Anti-debug: IsDebuggerPresent, CheckRemoteDebuggerPresent, NtQueryInformationProcess
- Timing attack detection
- Process name scanning for known debuggers
- DigiCert timestamp on all signatures

### Documentation
- README.md with usage examples
- DISTRIBUTION.md for package management
- changelogs/0.1.0.md with full release notes

---

## [Unreleased]

### Planned for v0.4.0
- Recursive operations (`-r` flag)
- More POSIX flags
- `--help` for all commands

### Planned for v1.0.0
- Full coreutils suite
- Stable API
- Documentation site

---

[0.3.0]: https://github.com/CRTYPUBG/winux/releases/tag/v0.3.0
[0.2.0]: https://github.com/CRTYPUBG/winux/releases/tag/v0.2.0
[0.1.0]: https://github.com/CRTYPUBG/winux/releases/tag/v0.1.0
[Unreleased]: https://github.com/CRTYPUBG/winux/compare/v0.3.0...HEAD
