# Changelog

All notable changes to PassForge are documented in this file.

## Unreleased

### Added
- Separate CI workflow (format, vet, tests, race, WASM build, privacy checks, govulncheck)
- Dependabot updates for GitHub Actions
- Property-based generator tests
- Estimated entropy bits in the UI
- PIN low-entropy warning
- Trust & Privacy and “Can I trust PassForge?” README sections
- docs/ architecture, security, threat model, development, and specification

### Changed
- Default: exclude ambiguous characters is **off** (better entropy for copy/paste use)
- Strength label clarified as **Estimated strength**
- Specs.md relocated to `docs/specification.md`
