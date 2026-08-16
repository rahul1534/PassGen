# Threat model

## Assets

- Generated passwords / passphrases / PINs (high sensitivity)
- User configuration (mode, length, character options) — non-sensitive

## Trust boundaries

1. **User browser** — generates and displays secrets.
2. **GitHub Pages CDN** — delivers static assets over HTTPS.
3. **Source repository** — public code reviewers can inspect generators and CSP.

## Threats and mitigations

| Threat | Mitigation |
|--------|------------|
| Weak / predictable RNG | Web Crypto only; fail closed if unavailable; rejection sampling for uniform ints |
| Password exfiltration | No backend; CSP `connect-src 'self'`; privacy CI checks |
| Accidental persistence | No password history; no localStorage for secrets |
| Compromised CDN / hosting | Open source + reproducible `make build`; HTTPS from Pages |
| XSS / injection | Strict CSP, no third-party scripts, minimal JS |
| Misleading strength claims | UI labels estimates as theoretical entropy, not site-specific safety |

## Out of scope

- Compromised user machine / malicious browser extensions
- Shoulder surfing of displayed passwords
- Server-side password policies on third-party sites
- Social engineering

## Residual risk

Anyone who can change the deployed static assets can change the generator. Public source + CI + user-runnable local builds reduce, but do not eliminate, that risk.
