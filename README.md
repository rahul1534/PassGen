# PassForge

PassForge is a privacy-first, open-source password generator that runs entirely in your browser. Passwords are generated locally with cryptographically secure randomness and are never sent to a server.

## Features

- **Random passwords** with configurable length, character sets, minimums, exclusions, and ambiguous-character filtering
- **Passphrases** from a bundled word list (EFF Large Wordlist, 7,776 words)
- **PIN / numeric codes** with optional pattern avoidance
- **Strength indicator** based on estimated entropy
- **Copy to clipboard** with accessible feedback
- **Light, dark, and system themes**
- **Responsive, keyboard-accessible UI**
- **No backend** — static HTML, CSS, and Go WebAssembly only

## Why It's Secure

- Passwords are generated in the browser using `crypto.getRandomValues()`
- No analytics, no accounts, no database
- Generated passwords are not stored in localStorage or sent over the network
- Core generation logic is separated from UI and covered by unit tests

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
```

## Build

```bash
make build
```

Output is written to `dist/`:

```text
dist/
├── index.html
├── app.wasm
├── wasm_exec.js
├── styles.css
└── favicon.svg
```

## Deploy to GitHub Pages

1. Push to the `main` branch.
2. Enable GitHub Pages for the repository (GitHub Actions deployment).
3. The workflow in `.github/workflows/deploy.yml` builds WebAssembly and publishes `dist/`.

After deployment, the app is available at:

```text
https://<username>.github.io/<repository>/
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

## Word List License

Passphrases use the [EFF Large Wordlist](https://www.eff.org/deeplinks/2016/07/new-wordlists-random-passphrases), licensed under [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/us/).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

See [SECURITY.md](SECURITY.md) for our security policy and how to report vulnerabilities.

## License

MIT — see [LICENSE](LICENSE).
