# Security overview

PassForge treats the browser as the trust boundary.

## Guarantees

- Passwords are generated locally with `crypto.getRandomValues()`.
- There is no password-generation server.
- Generated passwords are not written to `localStorage`, `sessionStorage`, IndexedDB, cookies, or URL parameters.
- No analytics or telemetry in V1.
- Production generation paths must not use `math/rand`.

## What GitHub Pages does

GitHub Pages only serves static files (HTML, CSS, JS, WASM). It does not generate passwords. Users still need to verify that the served assets match the open-source repository.

## Verification

- Unit and property tests in `internal/generator`
- CI: `go test`, `go vet`, race tests, WASM build
- Privacy regression script: `scripts/privacy-check.sh`
- Vulnerability scanning via `govulncheck` in CI

## Reporting

See [SECURITY.md](../SECURITY.md). Do not include real passwords in reports.
