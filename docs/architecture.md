# Architecture

PassForge is a static, client-side password generator.

```text
Browser
   │
   ├── HTML / CSS
   ├── app.js (WASM bootstrap)
   └── Go WebAssembly
          │
          ├── UI state
          ├── password / passphrase / PIN generation
          ├── clipboard
          └── Web Crypto (crypto.getRandomValues)
```

## Package layout

| Path | Role |
|------|------|
| `cmd/passforge` | WASM entry point |
| `internal/generator` | Password, passphrase, PIN, strength (browser-independent) |
| `internal/random` | `Source` interface, WASM CryptoSource, test DeterministicSource |
| `internal/web` | DOM binding and UI logic |
| `web/` | Static HTML, CSS, bootstrap JS |
| `wordlist/` | Bundled EFF Large Wordlist |

## Design rules

- Generation never depends on a backend.
- Production randomness uses Web Crypto only; no `math/rand` fallback.
- Generated passwords are not logged, stored, or placed in URLs.
- Assets use relative paths so GitHub project Pages work under `/PassGen/`.

## Content Security Policy

See comments in `web/index.html`. The Go WASM runtime requires `'unsafe-eval'` / `'wasm-unsafe-eval'`. `connect-src 'self'` allows loading `app.wasm` and blocks third-party network calls.
