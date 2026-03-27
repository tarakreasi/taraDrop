# Contributing to taraDrop

Thank you for your interest in contributing. taraDrop is a simple, focused tool and contributions should align with that philosophy: keep it lean, portable, and dependency-free.

---

## Code of Conduct

Be respectful. Keep discussions focused on the project.

---

## How to Contribute

### Reporting Bugs

Open an issue and include:

- Your OS and version
- Steps to reproduce
- Expected vs actual behavior
- Any relevant logs or screenshots

### Suggesting Features

Open an issue with the `enhancement` label. Describe:

- The problem you are trying to solve
- Your proposed solution
- Any alternatives you considered

### Submitting a Pull Request

1. Fork the repository
2. Create a branch from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```
3. Make your changes
4. Run the build to verify nothing is broken:
   ```bash
   make linux
   ```
5. Commit using [Conventional Commits](https://www.conventionalcommits.org/):
   ```
   feat: add clipboard paste upload support
   fix: correct IP detection on dual-stack interfaces
   docs: update README with macOS instructions
   ```
6. Push and open a Pull Request against `main`

---

## Development Setup

**Requirements:**

- Go 1.21 or later
- For Windows cross-compilation: `x86_64-w64-mingw32-gcc` (MinGW on Linux)

**Build:**

```bash
# Linux binary
make linux

# Windows binary (cross-compiled from Linux)
make windows

# Both
make all
```

**Note:** Do NOT commit built binaries (`taraDrop`, `taraDrop.exe`) to the repository. Binaries are distributed via GitHub Releases only.

---

## Project Principles

- Single binary, zero runtime dependencies
- Works without configuration
- Accessible to non-technical users on Windows (GUI, no terminal)
- Security-conscious: no arbitrary file access, validated paths

---

## Contact

Maintainer: Tri Wantoro
GitHub: [tarakreasi](https://github.com/tarakreasi)
Email: ajarsinau@gmail.com
