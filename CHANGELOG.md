# Changelog

All notable changes to taraDrop are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Planned
- Unit tests for HTTP handlers
- Security hardening (path traversal protection, input sanitization)
- CI/CD pipeline via GitHub Actions
- Automated binary releases

---

## [1.2.0] - 2026-02-23

### Added
- Port 80 fallback logic for Windows builds
- Linux binary distributed alongside Windows executable

### Changed
- Improved IP detection reliability across network interfaces

---

## [1.1.0] - 2026-02-18

### Added
- Multiple file upload support
- Chunked upload mechanism for large files (up to 16 GB)
- Progress tracking per file in the browser UI

---

## [1.0.0] - 2026-01-28

### Added
- Full rewrite from Python (Flask) to Go for zero-dependency distribution
- Fyne-based GUI control panel for Windows (no terminal required)
- Single-binary build for Linux and Windows
- Embedded HTML template via Go embed
- Smart local IP detection (UDP trick + interface fallback)
- Auto port selection (80, fallback to 5000)
- Cross-platform folder open support (xdg-open, explorer, open)

### Removed
- Python/Flask server (server.py, requirements.txt) replaced by Go implementation

---

## [0.1.0] - 2026-01-19

### Added
- Initial Python/Flask implementation (server.py)
- Glassmorphism HTML upload UI
- Basic file upload via browser
