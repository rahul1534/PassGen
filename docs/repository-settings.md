# Repository settings (manual)

These items from the review require GitHub UI (or `gh`) and are not encoded in git:

## Branch protection for `main`

Settings → Branches → Add rule / ruleset:

- Require a pull request before merging (optional for solo; recommended)
- Require status checks: **CI** / `test` and `security`
- Do not allow force pushes
- Do not allow deletions

## About section

Settings → General → About:

- **Description:** Privacy-first, open-source password, passphrase and PIN generator that runs entirely in your browser.
- **Website:** https://rahul1534.github.io/PassGen/
- **Topics:** `password-generator`, `password`, `security`, `privacy`, `passphrase`, `cryptography`, `golang`, `go`, `webassembly`, `wasm`, `github-pages`

## Visibility

Keep the repository **public**. For this project, public source is part of the trust model. A private repo with a public Pages site requires GitHub Pro (personal) and still serves a public site by default.
