# Changelog

All notable changes to PassForge are documented in this file.

## Unreleased

### Changed
- Redesigned interface: light-first layout with a segmented mode selector, a large colour-coded readout (digits blue, symbols orange), a five-segment strength meter, character-type chips that double as a colour legend, and sentence-case copy throughout
- Length, word-count and PIN-length controls now have sliders synced with their number fields
- Options and sliders regenerate the result immediately (Generate and Ctrl/⌘+Enter still create a new one)
- The strength meter is empty (not "Very Weak") when there is no result; Copy is disabled until there is one
- Typography uses the system font stacks only: no font downloads, no new binary assets

### Accessibility
- Mode tabs and character chips now show a visible keyboard focus ring (previously the visually hidden radio inputs swallowed it)
- The strength meter exposes `aria-valuenow`/`aria-valuetext`; the readout has an accessible name; copying is announced to screen readers
- Fixed low contrast on the copied-state and hover styles (hover no longer fades the label with opacity)

### Security
- Removed the Cloudflare Web Analytics beacon and the CSP allowances added for it; `script-src` and `connect-src` are same-origin again, matching the "no analytics / nothing sent to a server" guarantees
- Dropped `'unsafe-eval'` from the CSP (only `'wasm-unsafe-eval'` is required)
- Password options: negative minimum counts are now rejected (they could offset oversized minimums and produce output longer than requested)
- `scripts/privacy-check.sh` is now case-insensitive, fails closed on errors, bans external resources and analytics vendors, and asserts the CSP matches the approved policy exactly
- Browser tests assert no cross-origin requests and no CSP violations at runtime
- Deploy workflow now runs the privacy check, `go vet` and tests before publishing, and only the deploy job has Pages/OIDC permissions
- Security workflow: `govulncheck` is blocking, `gosec`/`staticcheck` are explicit advisory steps, least-privilege token permissions; removed the never-run gitleaks placeholder

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
- Local `make dev` server is Go (`cmd/devserver`) instead of Python
