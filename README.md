# PassForge

[![CI](https://github.com/rahul1534/PassGen/actions/workflows/ci.yml/badge.svg)](https://github.com/rahul1534/PassGen/actions/workflows/ci.yml)
[![Deploy](https://github.com/rahul1534/PassGen/actions/workflows/deploy.yml/badge.svg)](https://github.com/rahul1534/PassGen/actions/workflows/deploy.yml)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

PassForge is a privacy-first, open-source password generator that runs entirely in your browser. Passwords are generated locally with cryptographically secure randomness and are never sent to a server.

**Live demo:** https://rahul1534.github.io/PassGen/

## Privacy

PassForge is entirely client-side.

- Passwords generated locally
- Web Crypto API (`crypto.getRandomValues`)
- No backend
- No analytics
- No accounts
- No cookies for secrets
- No password history
- No `localStorage` for generated passwords
- No password transmitted over the network
- Works offline after assets are loaded

See [SECURITY.md](SECURITY.md) and [docs/security.md](docs/security.md).

## Can I trust PassForge?

PassForge is open source. The password generation code is available for inspection, and passwords are generated locally using the browser’s Web Crypto API.

There is no password-generation server. You can also [run the application locally](#quick-start) yourself.

## Features

- **Random passwords** with configurable length, character sets, minimums, exclusions, and optional ambiguous-character filtering
- **Passphrases** from a bundled word list (EFF Large Wordlist, 7,776 words)
- **PIN / numeric codes** with optional pattern avoidance (and a clear low-entropy warning)
- **Estimated strength** based on theoretical entropy (honestly labeled)
- **Copy to clipboard** with accessible feedback
- **Light, dark, and system themes**
- **Responsive, keyboard-accessible UI**
- **No backend** — static HTML, CSS, and Go WebAssembly only

## Why It's Secure

- Passwords are generated in the browser using `crypto.getRandomValues()`
- No analytics, no accounts, no database
- Generated passwords are not stored in localStorage or sent over the network
- Core generation logic is separated from UI and covered by unit + property tests
- CI runs tests, `go vet`, privacy regression checks, and vulnerability scanning

## Architecture

```text
┌───────────────────────────────┐
│          Browser              │
│                               │
│  ┌─────────┐   ┌──────────┐  │
│  │   UI    │──▶│ Go Wasm  │  │
│  └─────────┘   └────┬─────┘  │
│                     │        │
│               Web Crypto     │
│                     │        │
│               Password       │
│               Generator      │
└───────────────────────────────┘
```

More detail: [docs/architecture.md](docs/architecture.md) · [docs/threat-model.md](docs/threat-model.md)

## Quick Start

### Prerequisites

- Go 1.25+

### Run locally

```bash
git clone https://github.com/rahul1534/PassGen.git
cd PassGen
make build
make dev
```

Open http://localhost:8080

### Run tests

```bash
make test
make privacy
```

## Build

```bash
make build
```

Output is written to `dist/`:

```text
dist/
├── index.html
├── app.js
├── app.wasm
├── wasm_exec.js
├── styles.css
└── favicon.svg
```

## Deploy to GitHub Pages

### One-time setup

1. Open **Settings → Pages** for this repository:
   https://github.com/rahul1534/PassGen/settings/pages
2. Under **Build and deployment**, set **Source** to **GitHub Actions** (not “Deploy from a branch”).
3. Save. GitHub creates the `github-pages` deployment environment automatically.

### Deploy

1. Push to the `main` branch.
2. The workflow in `.github/workflows/deploy.yml` builds WebAssembly and publishes `dist/`.
3. If a deploy failed with `404` before Pages was enabled, re-run the workflow from the **Actions** tab after completing setup above.

After deployment, the app is available at:

```text
https://rahul1534.github.io/PassGen/
```

All asset paths are relative, so the app works on project pages.

## Configuration

Repository-specific values live in `internal/config/config.go`:

```go
const (
    AppName   = "PassForge"
    GitHubURL = "https://github.com/rahul1534/PassGen"
)
```

## Documentation

- [Architecture](docs/architecture.md)
- [Security overview](docs/security.md)
- [Threat model](docs/threat-model.md)
- [Development](docs/development.md)
- [Full specification](docs/specification.md) (historical product/tech requirements)

## Word List License

Passphrases use the [EFF Large Wordlist](https://www.eff.org/deeplinks/2016/07/new-wordlists-random-passphrases), licensed under [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/us/).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

See [SECURITY.md](SECURITY.md) for our security policy and how to report vulnerabilities.

## License

MIT — see [LICENSE](LICENSE).
