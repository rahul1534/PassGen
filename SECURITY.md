# Security Policy

## Our Commitment

PassForge is designed so that passwords are generated entirely in your browser and never transmitted to a server.

## What PassForge Does Not Do

- Does not send generated passwords over the network
- Does not store generated passwords in localStorage, sessionStorage, IndexedDB, cookies, or URLs
- Does not include analytics or telemetry in V1
- Does not log generated passwords

## Randomness

Production password generation uses the browser Web Crypto API (`crypto.getRandomValues()`). The application fails safely if secure randomness is unavailable rather than falling back to insecure sources.

## Content Security Policy

The app ships a strict CSP. Go’s WebAssembly runtime requires `'unsafe-eval'` / `'wasm-unsafe-eval'`. `connect-src 'self'` allows loading local `app.wasm` and blocks third-party network requests. See comments in `web/index.html` and [docs/architecture.md](docs/architecture.md).

## Automated checks

CI runs unit/property tests, `go vet`, a privacy regression script (`scripts/privacy-check.sh`), and `govulncheck`.

## Reporting a Vulnerability

If you discover a security issue, please report it responsibly:

1. Open a private security advisory on GitHub, or
2. Open an issue labeled `security` with a description of the problem and steps to reproduce.

Please do **not** include real passwords in your report.

We aim to acknowledge reports within a few days and will work on a fix as quickly as possible.

## Scope

Security reports are welcome for:

- Weak or predictable password generation
- Password leakage through network, storage, or logging
- Clipboard handling issues that expose secrets unexpectedly
- Content Security Policy or deployment misconfigurations

Out of scope:

- Social engineering
- Issues in third-party browsers unrelated to PassForge code
- Strength of passwords chosen by users outside the generator defaults
