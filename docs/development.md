# Development

## Prerequisites

- Go 1.25+
- `ripgrep` (`rg`) for privacy checks (optional locally; installed in CI)

## Commands

```bash
make test     # go test ./...
make build    # produce dist/
make dev      # build + serve dist/ with Go (PORT, default 8080)
make clean    # remove dist/
make privacy  # privacy regression checks
```

## Testing notes

- Generator tests use `DeterministicSource` and property suites; they do not need a browser.
- WASM UI code is behind `//go:build js && wasm`.
- Avoid `math/rand` in `internal/generator` and production random sources.

## Release checklist

1. `make test` and `make build` succeed locally.
2. CI is green on the PR.
3. Tag a release if publishing a versioned changelog (`v0.x.y`).
4. Confirm GitHub Pages deploy from `main`.

## Specs

Historical product/technical requirements live in [specification.md](specification.md). Prefer this guide and the README for day-to-day contribution.
