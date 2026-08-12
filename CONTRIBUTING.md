# Contributing to PassForge

Thank you for your interest in contributing!

## Development Setup

1. Install Go 1.25 or later.
2. Clone the repository.
3. Run tests:

```bash
make test
```

4. Build the WebAssembly app:

```bash
make build
```

5. Preview locally:

```bash
make dev
```

Then open http://localhost:8080

## Project Structure

- `cmd/passforge/` — WebAssembly entry point
- `internal/generator/` — password, passphrase, and PIN generation (browser-independent)
- `internal/random/` — secure and deterministic random sources
- `internal/web/` — DOM integration for the browser
- `web/` — static HTML, CSS, and favicon
- `wordlist/` — bundled EFF word list for passphrases

## Guidelines

- Keep password generation logic independent of browser APIs so it remains unit-testable.
- Never log or persist generated passwords.
- Use `crypto.getRandomValues()` for production randomness; use `DeterministicSource` only in tests.
- Match existing code style and keep changes focused.
- Add or update tests for generator behavior changes.

## Pull Requests

1. Fork the repository and create a feature branch.
2. Ensure `make test` passes.
3. Ensure `make build` succeeds.
4. Describe what changed and why in the pull request.

## Word List

The passphrase word list is derived from the [EFF Large Wordlist](https://www.eff.org/deeplinks/2016/07/new-wordlists-random-passphrases) (CC BY 3.0). If you replace or extend the list, document the source and license.

## Security

See [SECURITY.md](SECURITY.md) for reporting security issues.
