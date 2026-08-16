# Open-Source Password Generator — Product & Technical Specification

## 1. Overview

Build a privacy-first, open-source web application for generating strong passwords and passphrases.

The application must:

- Be written entirely in **Go**
- Run entirely in the user's browser
- Require **no backend server**
- Never transmit generated passwords anywhere
- Be deployable to **GitHub Pages**
- Require minimal configuration to deploy
- Have sensible, security-focused defaults
- Allow users to customize password requirements
- Support multiple password generation modes
- Be responsive and usable on desktop and mobile
- Be accessible
- Be suitable for an open-source GitHub repository

The application should feel like a polished standalone utility rather than a developer demo.

Suggested name:

**PassForge**

Alternative names can be considered by the implementation agent.

---

# 2. Core Product Principles

### Privacy first

Passwords must be generated locally in the browser.

The application must never:

- Send generated passwords to a server
- Store generated passwords remotely
- Send password configuration to analytics services
- Log generated passwords
- Include passwords in URLs
- Persist generated passwords to localStorage

The application should work completely offline after the initial page load.

### Secure randomness

Password generation must use a **cryptographically secure random number generator**.

The preferred source is:

```text
Web Crypto API → crypto.getRandomValues()
```

Because the application is written in Go and compiled to WebAssembly, implement a Go abstraction around the browser's cryptographically secure RNG.

Do NOT use:

- math/rand
- timestamps
- counters
- predictable seeds
- UUIDs as a password RNG
- pseudo-random browser APIs

If Web Crypto cannot be accessed, the application should fail safely rather than silently fall back to an insecure RNG.

---

# 3. Technology Requirements

## Language

Primary language:

**Go**

Target Go version:

```text
Go 1.26+
```

The implementation should avoid unnecessary third-party dependencies.

## Frontend

The UI should be generated and controlled by Go/Wasm.

Preferred architecture:

```text
Browser
   │
   ├── HTML
   ├── CSS
   └── Go WebAssembly
          │
          ├── UI state
          ├── password generation
          ├── clipboard integration
          └── Web Crypto integration
```

Do not introduce React, Vue, Angular, Node.js, or another frontend framework.

A small amount of JavaScript is acceptable only where required to bootstrap WebAssembly or access browser APIs that are impractical to access directly from Go.

## Build output

The build must produce static files such as:

```text
dist/
├── index.html
├── app.wasm
├── wasm_exec.js
├── styles.css
├── favicon.svg
└── assets/
```

These files must be directly hostable by GitHub Pages.

---

# 4. Deployment Requirements

The project must support:

```text
GitHub repository
        ↓
GitHub Actions
        ↓
Go build
        ↓
WebAssembly
        ↓
GitHub Pages
```

The repository must include a GitHub Actions workflow:

```text
.github/workflows/deploy.yml
```

The workflow should:

1. Checkout repository
2. Install/setup Go
3. Build the WebAssembly application
4. Generate the static deployment directory
5. Deploy to GitHub Pages

No external deployment service should be required.

The README must contain simple deployment instructions.

Target experience:

```text
git clone ...
cd ...
go build / make build
```

for local development and:

```text
git push
```

for GitHub Pages deployment.

---

# 5. Application Modes

The application must support at least four password generation modes.

## Mode 1 — Random Password

Traditional random password.

Example:

```text
vT8#qL2!xP7@rK
```

Configurable:

- Length
- Uppercase letters
- Lowercase letters
- Numbers
- Symbols
- Excluded characters
- Minimum number of each character type

---

## Mode 2 — Strong Password

A simplified mode designed for users who don't want to understand password configuration.

The UI should expose:

```text
Length
[ 20 ]

Strength
[ Very Strong ]

[ Generate ]
```

Recommended default:

```text
Length: 20
Uppercase: enabled
Lowercase: enabled
Numbers: enabled
Symbols: enabled
```

The user should still be able to expand advanced settings.

---

## Mode 3 — Passphrase

Generate human-readable random passphrases.

Example:

```text
orbit-cactus-lantern-river-velvet
```

Configurable:

- Number of words
- Separator
- Capitalization
- Add number
- Add symbol
- Word-list selection

Recommended default:

```text
Words: 5
Separator: -
Capitalize: No
Add number: Yes
Add symbol: No
```

Example:

```text
orbit-cactus-lantern-river-velvet7
```

The word list must be bundled with the application.

The application must NOT download a word list from a remote server during generation.

---

## Mode 4 — PIN / Numeric

Generate numeric codes.

Configurable:

- Length
- Allow repeated digits
- Avoid ambiguous patterns

Examples:

```text
482917
```

Recommended default:

```text
Length: 6
```

Allow lengths from:

```text
4–32
```

---

# 6. Optional Future Modes

Design the architecture so additional modes can be added easily.

Potential future modes:

- Memorable password
- Alphanumeric password
- Wi-Fi password
- API token
- Hex token
- Base64 token
- Diceware passphrase
- Pronounceable password

These should not need a major rewrite of the application.

---

# 7. User Interface

## Main layout

Desktop:

```text
┌──────────────────────────────────────────────────────────┐
│ PassForge                                  GitHub / About │
├──────────────────────────────────────────────────────────┤
│                                                          │
│                 Password Generator                      │
│                                                          │
│  [ Random Password ] [ Passphrase ] [ PIN ]              │
│                                                          │
│  Generated Password                                      │
│  ┌────────────────────────────────────────────────────┐  │
│  │  vT8#qL2!xP7@rK                                    │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ Copy ]                    [ Generate ]               │
│                                                          │
│  Strength: ████████████████████ Very Strong             │
│                                                          │
│  ──────────────────────────────────────────────────────  │
│                                                          │
│  Password Requirements                                   │
│                                                          │
│  Length              [ 20 ]                              │
│                                                          │
│  ☑ Uppercase         ☑ Lowercase                        │
│  ☑ Numbers           ☑ Symbols                          │
│                                                          │
│  Advanced Options ▼                                      │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

Mobile layout should stack controls vertically.

---

# 8. Quality Defaults

The application must open with a useful configuration.

Default mode:

```text
Random Password
```

Default settings:

```text
Length: 20

Uppercase: enabled
Lowercase: enabled
Numbers: enabled
Symbols: enabled

Minimum uppercase: 1
Minimum lowercase: 1
Minimum numbers: 1
Minimum symbols: 1

Exclude ambiguous characters: enabled

No repeated characters: disabled
```

The first generated password should already be considered strong.

Do not force the user to configure anything before generating a password.

---

# 9. Password Requirements UI

## Basic settings

Display the most important options immediately:

```text
Password length
[ 20 ]

☑ Uppercase letters
☑ Lowercase letters
☑ Numbers
☑ Symbols
```

## Advanced settings

Expandable section:

```text
Advanced Options

Minimum uppercase characters [ 1 ]
Minimum lowercase characters [ 1 ]
Minimum numbers             [ 1 ]
Minimum symbols             [ 1 ]

☑ Exclude ambiguous characters

Characters to exclude
[ ]

☐ Prevent repeated characters
```

---

# 10. Validation

The application must validate the configuration before generating a password.

Examples:

### Invalid character configuration

If the user disables all character groups:

```text
Please select at least one character type.
```

### Impossible minimums

If:

```text
Length = 8
Minimum uppercase = 5
Minimum lowercase = 5
```

Display:

```text
Minimum character requirements exceed password length.
```

### Exclusions

If exclusions remove every character from a selected character group:

```text
Your excluded characters remove all available symbols.
```

The Generate button should either be disabled or generation should show a clear validation error.

---

# 11. Character Sets

Default character sets:

### Lowercase

```text
abcdefghijklmnopqrstuvwxyz
```

### Uppercase

```text
ABCDEFGHIJKLMNOPQRSTUVWXYZ
```

### Numbers

```text
0123456789
```

### Symbols

Use a conservative set that works well across websites:

```text
!@#$%^&*()-_=+[]{};:,.?
```

Avoid including characters that are likely to cause problems when copied into URLs, shells, HTML, SQL, or configuration files unless explicitly requested.

---

# 12. Ambiguous Characters

When enabled, exclude commonly confused characters.

Suggested exclusions:

```text
0
O
o
1
I
l
```

The implementation should apply this consistently across the selected character sets.

Display a tooltip:

```text
Removes characters that can be difficult to distinguish,
such as 0/O and 1/l/I.
```

---

# 13. Generation Algorithm

The algorithm must guarantee minimum requirements.

Example:

```text
Length = 20

Uppercase minimum = 2
Lowercase minimum = 2
Numbers minimum = 2
Symbols minimum = 2
```

Algorithm:

1. Build the available character sets.
2. Validate configuration.
3. Select the required minimum number from each required character set using CSPRNG.
4. Fill remaining positions from the combined allowed character pool.
5. Securely shuffle the final password using CSPRNG.
6. Return the password.

Do NOT simply generate a random string and check whether requirements happen to be satisfied.

The algorithm must intentionally satisfy the requirements.

---

# 14. Secure Random Selection

Create an abstraction such as:

```go
type RandomSource interface {
    RandomInt(max int) (int, error)
    Shuffle[T any](items []T) error
}
```

Production implementation:

```text
CryptoRandomSource
```

Testing implementation:

```text
DeterministicRandomSource
```

The deterministic implementation should only be used in tests.

Production code must use browser cryptographic randomness.

---

# 15. Entropy / Strength Indicator

Display a password strength indicator.

Example:

```text
Strength: Very Strong
```

Possible levels:

```text
Very Weak
Weak
Fair
Strong
Very Strong
```

Strength should be based primarily on:

- Length
- Character pool size
- Character diversity
- Estimated entropy

For a random password, approximate entropy:

```text
entropy = length × log2(characterPoolSize)
```

Do not claim that this represents the actual security of every website.

For passphrases, calculate entropy based on:

```text
number of words × log2(word-list-size)
```

The UI should use simple language.

Example tooltip:

```text
Estimated entropy assumes the password was generated randomly
from the selected character set.
```

---

# 16. Generated Password Display

The generated password must:

- Be visually prominent
- Use a monospace font
- Be selectable
- Not be truncated without a way to view the full value

Provide:

```text
Copy
Regenerate
```

Optional:

```text
Show / Hide
```

for extremely long passwords.

The password should never be inserted into:

- URL query parameters
- browser history
- console logs
- analytics
- error reporting
- localStorage

---

# 17. Clipboard

Provide a Copy button.

Behavior:

```text
Copy
↓
Copied!
```

After approximately 2 seconds:

```text
Copied!
↓
Copy
```

Use the browser Clipboard API where available.

If clipboard access is unavailable:

```text
Unable to access clipboard. Select and copy the password manually.
```

Do not expose the password in error messages.

---

# 18. Regeneration

Clicking Generate should immediately create a new password.

Keyboard support:

```text
Enter
```

should generate when focus is inside relevant controls where appropriate.

Optional shortcut:

```text
Cmd/Ctrl + Enter
```

Generate password.

---

# 19. Password History

Do NOT implement password history by default.

Reason:

Password history creates unnecessary sensitive data retention.

If history is eventually implemented, it should be:

- Opt-in
- Memory-only
- Cleared on page refresh
- Never persisted

For V1:

```text
No password history.
```

---

# 20. Local Storage

Do not store generated passwords.

It is acceptable to optionally store **non-sensitive UI preferences**, such as:

```text
Selected mode
Password length
Theme
```

However, V1 should preferably avoid persistence unless it provides clear value.

A simple reset button should restore defaults.

---

# 21. Theme

Support:

```text
Light
Dark
System
```

Default:

```text
System
```

The UI should respect:

```text
prefers-color-scheme
```

Avoid excessive visual complexity.

The application should feel similar to a modern developer utility.

---

# 22. Accessibility

Target WCAG 2.2 AA where practical.

Requirements:

- Keyboard navigable
- Visible focus states
- Proper labels
- Semantic HTML
- Screen-reader-friendly controls
- Sufficient contrast
- No color-only indicators
- Buttons must have accessible names
- Form validation errors must be accessible

Do not rely solely on icons.

For example:

Bad:

```text
📋
```

Better:

```text
Copy
```

An icon can accompany the text.

---

# 23. Responsive Design

Support:

- Desktop
- Tablet
- Mobile

Minimum target width:

```text
320px
```

No horizontal scrolling for normal usage.

Recommended content width:

```text
720–900px
```

---

# 24. About / Privacy Section

Include a small footer:

```text
PassForge is open source.

Your passwords are generated locally in your browser
and are never sent to a server.

View source
```

Link:

```text
GitHub
```

The GitHub URL should be configurable through a single build/config constant rather than scattered throughout the code.

---

# 25. Offline Support

The application should function without a network connection after loading.

Do not depend on:

- CDN-hosted JavaScript
- CDN-hosted CSS
- External fonts
- Remote word lists
- External APIs
- Analytics services

All runtime assets must be bundled.

Optional future enhancement:

```text
Service Worker / PWA
```

But this is not required for V1.

---

# 26. Project Structure

Recommended structure:

```text
passforge/
│
├── cmd/
│   └── passforge/
│       └── main.go
│
├── internal/
│   ├── generator/
│   │   ├── generator.go
│   │   ├── password.go
│   │   ├── passphrase.go
│   │   ├── pin.go
│   │   ├── charset.go
│   │   └── strength.go
│   │
│   ├── random/
│   │   ├── random.go
│   │   ├── crypto.go
│   │   └── deterministic.go
│   │
│   └── web/
│       ├── app.go
│       ├── clipboard.go
│       └── browser.go
│
├── web/
│   ├── index.html
│   ├── styles.css
│   └── favicon.svg
│
├── wordlist/
│   └── words.txt
│
├── tests/
│
├── scripts/
│   └── build.sh
│
├── .github/
│   └── workflows/
│       └── deploy.yml
│
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── LICENSE
└── .gitignore
```

The implementation agent may modify the structure if it results in a cleaner architecture.

---

# 27. Go Architecture

Separate password generation from UI.

The core generator must be completely independent of WebAssembly/browser APIs.

Example:

```go
type PasswordOptions struct {
    Length                  int
    Uppercase               bool
    Lowercase               bool
    Numbers                 bool
    Symbols                 bool
    MinUppercase            int
    MinLowercase            int
    MinNumbers              int
    MinSymbols              int
    ExcludeAmbiguous        bool
    ExcludedCharacters      string
    PreventRepeated         bool
}
```

Generator:

```go
type PasswordGenerator interface {
    Generate(options PasswordOptions) (string, error)
}
```

Passphrase:

```go
type PassphraseOptions struct {
    Words             int
    Separator         string
    Capitalize        bool
    AddNumber         bool
    AddSymbol         bool
}
```

PIN:

```go
type PINOptions struct {
    Length              int
    AllowRepeatedDigits bool
}
```

---

# 28. Error Handling

Use typed errors where useful.

Example:

```go
var (
    ErrInvalidLength       = errors.New("invalid password length")
    ErrNoCharacterSet      = errors.New("no character set selected")
    ErrImpossibleMinimums  = errors.New("minimum requirements exceed length")
    ErrNoAvailableChars    = errors.New("no available characters")
    ErrRandomSourceFailure = errors.New("secure random source unavailable")
)
```

UI should convert technical errors into human-readable messages.

Never display stack traces to users.

---

# 29. Testing

The project must have comprehensive unit tests.

## Password generator tests

Test:

- Minimum length
- Maximum length
- Character inclusion
- Minimum character requirements
- Excluded characters
- Ambiguous character exclusion
- Impossible configurations
- Empty character sets
- Repeated character option
- Different lengths
- Different modes

Example property:

For every generated password:

```text
len(password) == requested length
```

And when requirements are specified:

```text
uppercaseCount >= minUppercase
lowercaseCount >= minLowercase
numberCount >= minNumbers
symbolCount >= minSymbols
```

---

# 30. Statistical Tests

Do not rely on statistical randomness tests as the primary security test.

However, basic sanity tests may verify that:

- Results are not identical repeatedly
- All characters in the configured pool can eventually appear
- Shuffle does not systematically bias obvious positions

Avoid flaky tests.

---

# 31. Browser Tests

If practical, add browser-level tests for:

- Generate button
- Copy button
- Mode selection
- Configuration changes
- Validation
- Theme
- Responsive UI

Do not make browser tests a requirement for every local build.

---

# 32. Security Requirements

The implementation must follow these rules.

### Never log passwords

Forbidden:

```go
log.Printf("Generated password: %s", password)
```

### Never include passwords in telemetry

No analytics in V1.

### Never persist passwords

No:

```text
localStorage
sessionStorage
IndexedDB
cookies
URL parameters
```

for generated passwords.

### Never send password data over the network

The only network request after initial loading should be unnecessary for normal operation.

Ideally:

```text
Network requests during generation = 0
```

---

# 33. Content Security Policy

Add a strict Content Security Policy where compatible.

Prefer something similar to:

```text
default-src 'self';
script-src 'self';
style-src 'self';
img-src 'self' data:;
connect-src 'none';
font-src 'self';
object-src 'none';
base-uri 'none';
form-action 'none';
```

The exact CSP must be compatible with the WebAssembly runtime and GitHub Pages.

Do not blindly copy the above policy if it prevents the application from functioning.

---

# 34. SEO

This is a utility application, so SEO should remain simple.

Include:

```html
<title>PassForge — Secure Password Generator</title>

<meta
  name="description"
  content="Generate secure passwords and passphrases locally in your browser."
>
```

Add Open Graph metadata where practical.

Do not include generated passwords in metadata.

---

# 35. GitHub Repository

The repository should contain:

```text
README.md
LICENSE
CONTRIBUTING.md
SECURITY.md
```

README should explain:

1. What the project does
2. Why it is secure
3. How randomness works
4. How to run locally
5. How to build
6. How to deploy to GitHub Pages
7. How to contribute
8. Security reporting

Recommended license:

```text
MIT
```

unless otherwise specified by the repository owner.

---

# 36. Security Documentation

Create:

```text
SECURITY.md
```

Explain:

- Passwords never leave the browser
- No analytics
- CSPRNG is used
- No password persistence
- How security issues should be reported

Do not encourage users to submit actual passwords when reporting issues.

---

# 37. GitHub Actions

Create:

```text
.github/workflows/deploy.yml
```

The workflow should:

```text
on:
  push:
    branches:
      - main
```

Build:

```text
GOOS=js
GOARCH=wasm
go build
```

Copy the Go WebAssembly runtime support file.

Generate:

```text
dist/index.html
dist/app.wasm
dist/wasm_exec.js
dist/styles.css
```

Then deploy `dist/` to GitHub Pages.

Use official GitHub Pages deployment actions where possible.

---

# 38. GitHub Pages Compatibility

The generated application must work when hosted at:

```text
https://<username>.github.io/<repository>/
```

Do NOT assume the application is hosted at:

```text
/
```

All assets should use relative paths or correctly derive the base path.

For example:

```html
<script src="./wasm_exec.js"></script>
```

rather than:

```html
<script src="/wasm_exec.js"></script>
```

This is important for GitHub project pages.

---

# 39. Local Development

Provide:

```text
make build
make dev
make test
make clean
```

Example:

```text
make build
```

should produce:

```text
dist/
```

Example:

```text
make test
```

should execute all Go tests.

For development, a simple static HTTP server should be sufficient.

Do not require Node.js.

---

# 40. Configuration

Avoid unnecessary configuration.

Repository-level configuration may contain:

```go
const (
    AppName       = "PassForge"
    GitHubURL     = "..."
    DefaultLength = 20
)
```

The implementation agent should make repository-specific values easy to modify.

---

# 41. UX Details

When the page loads:

1. Display Random Password mode.
2. Load quality defaults.
3. Generate a password immediately.
4. Focus should NOT automatically be placed in a text field.
5. Display strength.
6. Allow immediate Copy.

When the user changes a requirement:

```text
Do not automatically generate on every keystroke.
```

Instead:

```text
[ Generate ]
```

should explicitly generate the new password.

Exception:

Mode changes may optionally generate a new password automatically.

---

# 42. Reset

Provide:

```text
Reset to defaults
```

This restores:

```text
Random Password
20 characters
Uppercase enabled
Lowercase enabled
Numbers enabled
Symbols enabled
Ambiguous characters excluded
Minimum one of each selected type
```

The current generated password should be replaced with a newly generated password after reset.

---

# 43. Copy Security Consideration

When copying:

```text
navigator.clipboard.writeText(password)
```

Do not expose the password in the DOM more than necessary.

The displayed password itself is necessarily visible, but avoid creating unnecessary duplicate copies.

Do not attempt aggressive clipboard clearing because this can interfere with normal user workflows.

---

# 44. Passphrase Word List

Use a reasonably large English word list.

Target:

```text
≥ 2,000 words
```

Prefer:

```text
4,000+ words
```

Words should be:

- Common
- Easy to spell
- Free of offensive terms
- Free of proper nouns where possible
- Suitable for random passphrases

The word list should be distributed with the application.

Document its source and license.

Do not use a word list whose license is incompatible with the project's open-source license.

---

# 45. Passphrase Security

Do not claim that:

```text
5 random words
```

is automatically secure in every context.

The strength calculation should account for the actual word list size.

For example, if the word list contains 4,096 words:

```text
12 bits/word
```

Five independent words:

```text
~60 bits
```

before considering optional additions.

The UI can explain this simply:

```text
More randomly selected words generally provide more security.
```

---

# 46. Visual Design

Use a clean, modern utility-style design.

Characteristics:

- Minimal
- Professional
- Fast
- High contrast
- Rounded controls
- Subtle borders
- Limited animations
- No unnecessary gradients
- No distracting illustrations

The password itself should be the visual focal point.

Use a monospace font stack such as:

```css
font-family:
  ui-monospace,
  SFMono-Regular,
  Menlo,
  Monaco,
  Consolas,
  "Liberation Mono",
  monospace;
```

Avoid external fonts.

---

# 47. Performance

The application should load quickly.

Targets:

- Minimal JavaScript
- Small Wasm binary where practical
- No external dependencies at runtime
- No network requests for generation
- Generation should feel instantaneous

Do not optimize prematurely, but avoid unnecessary Go dependencies.

---

# 48. Browser Compatibility

Target current versions of:

- Chrome
- Edge
- Firefox
- Safari
- iOS Safari

If Web Crypto or Clipboard APIs are unavailable, provide graceful fallback messaging.

Do not compromise security by falling back to insecure randomness.

---

# 49. Definition of Done

The implementation is complete when:

### Functionality

- [ ] Random password mode works
- [ ] Passphrase mode works
- [ ] PIN mode works
- [ ] Password requirements are configurable
- [ ] Quality defaults exist
- [ ] Minimum requirements are guaranteed
- [ ] Character exclusions work
- [ ] Ambiguous character exclusion works
- [ ] Strength indicator works
- [ ] Copy works
- [ ] Reset works
- [ ] Dark/light/system themes work

### Security

- [ ] CSPRNG is used
- [ ] No math/rand in production password generation
- [ ] Passwords are never sent to a server
- [ ] Passwords are not stored
- [ ] No analytics
- [ ] No password logging
- [ ] No password in URLs
- [ ] Security documentation exists

### Engineering

- [ ] Go tests pass
- [ ] Build works from a clean checkout
- [ ] No Node.js dependency
- [ ] GitHub Actions workflow works
- [ ] GitHub Pages deployment works
- [ ] Relative asset paths work under `/repository-name/`
- [ ] README contains setup/deployment instructions
- [ ] MIT license included

### UX

- [ ] Mobile responsive
- [ ] Keyboard accessible
- [ ] Screen-reader labels exist
- [ ] Validation errors are understandable
- [ ] Password is immediately visible
- [ ] User can generate without configuring anything

---

# 50. AI Coding Agent Instructions

The implementation agent should follow this process:

## Phase 1 — Architecture

Before writing code:

1. Inspect the specification.
2. Define the Go package structure.
3. Define interfaces for:
  - Password generation
  - Randomness
  - Passphrase generation
  - Strength calculation
4. Define browser/Wasm integration boundaries.

Do not over-engineer the project.

---

## Phase 2 — Core Generator

Implement and test:

```text
Random password
Passphrase
PIN
Character sets
Validation
CSPRNG
Shuffle
Strength estimation
```

The generator must be testable without a browser.

---

## Phase 3 — Browser Integration

Implement:

```text
WebAssembly bootstrap
DOM interaction
Clipboard
Web Crypto
Theme
UI state
```

Keep browser-specific code separate from the core generator.

---

## Phase 4 — UI

Implement the UI according to this specification.

Prioritize:

1. Generated password
2. Generate button
3. Copy button
4. Basic settings
5. Advanced settings
6. Mode selection
7. Strength indicator
8. Privacy information

---

## Phase 5 — Testing

Run:

```text
go test ./...
```

Add tests for all generator edge cases.

Verify that production generation never uses a deterministic or insecure RNG.

---

## Phase 6 — GitHub Pages

Implement:

```text
.github/workflows/deploy.yml
```

Verify the generated application works from:

```text
https://username.github.io/repository/
```

Pay particular attention to relative asset paths.

---

## Phase 7 — Documentation

Create/update:

```text
README.md
CONTRIBUTING.md
SECURITY.md
LICENSE
```

README should include an architecture diagram similar to:

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

---

# 51. Important Implementation Constraints

The agent MUST NOT:

- Introduce a backend server
- Introduce Node.js unless absolutely required by tooling
- Use React/Vue/etc.
- Use `math/rand` for production generation
- Send generated passwords to external services
- Add analytics
- Store passwords
- Use remote CDNs for runtime dependencies
- Download word lists dynamically
- Add unnecessary authentication
- Add a database

The final application should be able to run entirely as:

```text
Static HTML
+
CSS
+
Go WebAssembly
```

---

# 52. Recommended V1 Scope

Keep V1 focused.

### Must have

```text
✓ Random password
✓ Passphrase
✓ PIN
✓ Configurable requirements
✓ Secure randomness
✓ Strength indicator
✓ Copy
✓ Reset
✓ Dark/light/system theme
✓ Responsive UI
✓ Accessibility
✓ GitHub Pages
✓ Offline after load
✓ Open-source documentation
```

### V2

```text
○ API token generator
○ Diceware
○ PWA/offline installation
○ Import/export preferences
○ More passphrase dictionaries
○ Browser extension
○ QR code generation
```

Do not implement V2 features unless the core V1 implementation is complete and tested.

---

# 53. Final Product Goal

The finished application should answer the user's question:

> "I need a secure password. What should I use?"

with a frictionless experience:

```text
Open website
     ↓
Strong password already generated
     ↓
See strength
     ↓
Click Copy
     ↓
Done
```

At the same time, an advanced user should be able to say:

```text
I need a 32-character password,
at least 4 uppercase characters,
at least 4 numbers,
no ambiguous characters,
and no repeated characters.
```

and configure exactly that.

The application should combine:

**security + customization + simplicity + privacy + zero infrastructure.**

The ideal deployment experience is:

```text
GitHub repository
       ↓
git push
       ↓
GitHub Actions
       ↓
GitHub Pages
       ↓
Public password generator
```

No server. No database. No API keys. No accounts. No tracking.