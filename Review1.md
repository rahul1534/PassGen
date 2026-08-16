Executive assessment

I would keep this repository public. In fact, for this particular project, public is arguably the better configuration.

The project is already positioned as a privacy-first, open-source password generator, and its security story actually benefits from transparency: users can inspect the generator, verify that it uses browser cryptographic randomness, and verify that there is no backend. Your README explicitly makes those claims.

My current assessment:

Area	Assessment
Overall architecture	🟢 Very good
Cryptographic RNG design	🟢 Good
Privacy model	🟢 Excellent
GitHub Pages deployment	🟢 Good
CSP	🟢 Good
Password generation algorithm	🟢 Good
Passphrase design	🟢 Good
Accessibility	🟢 Good foundation
Testing	🟡 Needs more depth
CI/security automation	🟡 In progress
Repository presentation	🟡 Needs polish
Production/open-source readiness	🟡 Close, but not quite
Keep repo public?	🟢 Yes

The biggest thing I'd change is not repository visibility. I'd spend the effort making the repository look and behave like a mature security utility.

1. What you have right

The repository currently has a surprisingly solid structure for a project with only two commits on main. GitHub currently shows the project as public, with cmd/passforge, internal, scripts, web, wordlist, CI configuration, documentation, MIT license, and security/contribution policies.

The architecture is essentially:

Browser
   │
   ├── HTML/CSS
   │
   └── Go WebAssembly
          │
          ├── Password generation
          ├── Passphrase generation
          ├── PIN generation
          └── Web Crypto API

That is a very reasonable architecture for this application.

There is no server, database, API, analytics platform, authentication, or cloud backend. The README explicitly states that generated passwords aren't sent over the network or persisted.

That's exactly the architecture I'd want for a password generator.

2. The most important part: randomness

I specifically looked at your random-number implementation.

Your production WASM implementation accesses:

crypto.getRandomValues()

and fails if the browser doesn't provide it.

That's excellent.

More importantly, you aren't simply doing:

random number → modulo character-set-size

without considering modulo bias.

Your implementation uses rejection sampling before taking the modulo. That's the right general approach for uniform cryptographic selection.

Your generic random abstraction is also nicely separated:

generator
    ↓
random.Source
    ↓
CryptoSource

and the shuffle uses Fisher-Yates with the secure random source.

Recommendation

Keep this design.

I'd consider the RNG implementation one of the strongest parts of the project.

3. Password generation algorithm is also well designed

Your password generator doesn't just repeatedly generate a password until it happens to satisfy requirements.

Instead it:

selects the required minimum characters from each category
fills the remaining positions
securely shuffles the result

That's exactly what your specification calls for.

For example:

minimum uppercase
minimum lowercase
minimum numbers
minimum symbols
        ↓
fill remaining positions
        ↓
CSPRNG shuffle
        ↓
password

This is much better than the common implementation:

generate()
if invalid:
    generate again

because the latter can introduce unnecessary inefficiencies and makes reasoning about distributions harder.

One thing I'd add

Add property-based tests / fuzz tests around this.

For example, for thousands of generated passwords, automatically verify:

requested length
minimum uppercase
minimum lowercase
minimum numbers
minimum symbols
exclusions
no-repeat requirement
valid Unicode/rune handling
impossible configurations return errors
secure RNG errors propagate

You already have generator tests, which is a good start, but this is the area I'd expand substantially. The repository currently has only a small test footprint relative to the security claims you're making.

4. There is one important security/documentation distinction I'd make

Your README says:

"Passwords are generated locally in the browser using crypto.getRandomValues()."

That's good.

But I'd make an even stronger distinction:

The application itself doesn't need to be secret.

This is important because it answers your original question.

A password generator should ideally be something users can inspect.

If the repo were:

PRIVATE SOURCE
       ↓
mysterious password generator
       ↓
public website

a security-conscious user could reasonably ask:

"How do I know this website isn't sending my passwords somewhere?"

With your current architecture:

PUBLIC SOURCE
       ↓
users can inspect it
       ↓
Go/WASM
       ↓
Web Crypto
       ↓
password stays in browser

that's much more trustworthy.

So public source code is actually an advantage for PassForge.

5. Your CSP is a strong point

I found this in web/index.html:

default-src 'self'
script-src 'self' 'unsafe-eval' 'wasm-unsafe-eval'
connect-src 'self'
object-src 'none'
base-uri 'none'
form-action 'none'

That's a very good start.

In particular:

connect-src 'self'

This is useful because it limits outbound connections from the page.

object-src 'none'

Good hardening.

base-uri 'none'

Good protection against certain injection/base-tag attacks.

form-action 'none'

Good for an application that doesn't need form submissions.

Why unsafe-eval?

Go/WASM commonly needs special CSP allowances, so I wouldn't blindly remove it.

However, I'd document why those CSP exceptions exist.

That will make your security posture much easier for reviewers to understand.

6. Your JavaScript is pleasantly tiny

The actual web/app.js is only about 27 lines and essentially:

checks Web Crypto availability
loads WASM
starts Go
displays an error if startup fails.

That's excellent.

You haven't created a huge JavaScript application sitting around the security-critical code.

I'd keep this architecture.

7. Privacy model is excellent—but make it verifiable

Your SECURITY.md makes strong claims:

no network transmission of passwords
no localStorage
no sessionStorage
no IndexedDB
no cookies
no URLs
no analytics
no telemetry
no password logging

This is exactly the kind of security policy I'd expect from a privacy-focused password generator.

But there's an important principle:

The stronger the security claim, the stronger the automated verification should be.

I'd add automated checks that make those claims harder to accidentally break.

For example:

CI security checks

Search source for things such as:

localStorage
sessionStorage
indexedDB
navigator.sendBeacon
fetch(
XMLHttpRequest
WebSocket
analytics
telemetry

You obviously don't want to ban every occurrence forever, but you can make a lightweight privacy regression check.

For example:

PR changes
   ↓
privacy test
   ↓
"Did network/storage capabilities change?"

That would be particularly valuable for this project.

8. Your security workflow is currently the biggest process gap

I noticed something interesting.

main currently has only the original two commits, while you have an open PR #1 called security.yml with six commits.

The PR is specifically intended to add a security scanning workflow.

That's exactly the right direction.

However, I would clean up that PR before merging it.

The six commits include several iterative fixes:

"Fix workflow: add top-level keys and structure"
"reformat"
"Fix workflow: robust install step + logging"
"Fix workflow: quote debug step names..."
etc.

That's perfectly normal while developing, but before merging an open-source project, I'd squash this into something like:

ci: add security scanning workflow

Then merge it.

After that I'd have CI run:
go test ./...
go vet ./...
go test -race ./...
go build
security scanning

and ideally:

go test ./... -cover
9. I would add a separate CI workflow

Right now your deployment workflow does:

push main
   ↓
setup Go
   ↓
make build
   ↓
deploy GitHub Pages

That's fine for deployment.

But deployment should not be your primary test gate.

I'd separate:

ci.yml
PR / push
    ↓
format
    ↓
test
    ↓
vet
    ↓
build
    ↓
security checks

and:

deploy.yml
main
 ↓
build
 ↓
deploy Pages

This gives you a much cleaner open-source workflow.

10. Protect main

This is one of the most important GitHub repository recommendations I'd make.

Once the project is public, enable branch protection / rules for main.

I'd configure approximately:

main
 ├── PR required
 ├── CI must pass
 ├── conversation resolution
 ├── no force pushes
 └── no direct pushes

For a solo project, you don't necessarily need mandatory code review from another person.

But requiring your CI checks before merging is worthwhile.

That becomes particularly important because you're making security-related claims.

11. Add Dependabot

Your project currently has very few dependencies, which is excellent.

go.mod is only:

module github.com/rahul1534/PassGen
go 1.25

with no third-party dependencies.

That's a huge positive from a supply-chain perspective.

Still, I'd enable:

Dependabot
GitHub Actions dependency updates

because your workflow itself depends on:

actions/checkout
actions/setup-go
actions/upload-pages-artifact
actions/deploy-pages

Your deployment workflow currently uses those Actions.

12. Pin GitHub Actions more tightly if you want a security-focused project

Currently you're using versions such as:

actions/checkout@v4
actions/setup-go@v5
actions/upload-pages-artifact@v3
actions/deploy-pages@v4

That's reasonable.

For a project explicitly marketed around security, though, I'd eventually consider pinning Actions to immutable commit SHAs.

For example conceptually:

actions/checkout@<commit-sha>

rather than:

actions/checkout@v4

It's more cumbersome, but gives stronger supply-chain guarantees.

This isn't something I'd call a blocker for your project.

13. Passphrase implementation looks good

The passphrase implementation uses a bundled word list, rather than downloading words dynamically.

That's exactly what you want.

Your default is:

5 words
-
number: yes
symbol: no

and the repository documents the EFF Large Wordlist and CC BY 3.0 licensing.

That's good open-source hygiene.

One improvement I'd make

Make the UI explicitly explain:

"Security depends primarily on the number of randomly selected words, not on adding arbitrary symbols."

This prevents users from thinking:

word-word-word-word!7

is automatically stronger than simply using more words.

14. Your strength meter needs careful wording

Your entropy calculation is:

password length × log2(character pool size)

which is a legitimate theoretical entropy estimate under the assumption of uniform random generation.

Your specification actually acknowledges that limitation.

That's good.

But I'd make the UI wording even clearer.

Instead of:

Strength: Very Strong

I'd display:

Estimated strength: Very Strong

and perhaps:

~118 bits estimated entropy

with an information tooltip:

"This estimate assumes the password was generated randomly from the selected character set. It does not account for site-specific password policies, breached-password lists, or human choices."

That would make the product more technically honest.

15. One thing I would reconsider: "Exclude ambiguous characters"

Your default is:

Exclude ambiguous characters = ON

and your ambiguous list includes:

0 O o 1 I l

That's reasonable for passwords that humans may need to transcribe.

But it reduces the available entropy pool.

For a password generator intended primarily for copy/paste, I'd probably make the default:

Exclude ambiguous characters: OFF

and perhaps have:

☐ Exclude ambiguous characters

with an explanation:

Useful when passwords need to be read or typed manually.

This is a product decision rather than a security vulnerability.

16. Your PIN mode deserves a stronger warning

Your PIN defaults to:

6 digits
repeated digits allowed
avoid patterns disabled

That's technically reasonable.

But I'd add a prominent warning:

PINs provide much less entropy than passwords. Use them only where the system requires a numeric code.

A 6-digit PIN has only:

10^6 = 1,000,000

possible combinations.

That's fundamentally different from your 20-character password mode.

The UI should make this distinction clear.

17. Your repository name vs application name

This is currently:

Repository: PassGen
Application: PassForge

The README explicitly calls the application PassForge while the GitHub repository remains PassGen.

I would actually rename the repository to:

PassForge

if the name is available and you intend to keep the project long-term.

Then the URL becomes:

rahul1534.github.io/PassForge/

instead of:

rahul1534.github.io/PassGen/

But don't do this casually because it changes your public URL and existing links.

If you're still early—and you are, with essentially two commits on main—now is the easiest time to do it.

18. More important: add repository metadata

GitHub currently shows:

"No description, website, or topics provided."

even though your README is good.

That's a missed opportunity.

I'd set:

Description

Privacy-first, open-source password, passphrase and PIN generator that runs entirely in your browser.

Website
https://rahul1534.github.io/PassGen/
Topics

I'd use roughly:

password-generator
password
security
privacy
passphrase
cryptography
golang
go
webassembly
wasm
github-pages

This will make the project much more discoverable.

19. Add badges to README

I'd add something like:

CI
Security
Go version
License
GitHub Pages

For example conceptually:

[CI] [Security] [Go] [License: MIT]

This is particularly useful because you're selling the project on security and privacy.

20. Add a "Trust & Privacy" section to the README

This is probably the single biggest README improvement I'd make.

Something like:

## Privacy


PassForge is entirely client-side.


✓ Passwords generated locally
✓ Web Crypto API
✓ No backend
✓ No analytics
✓ No accounts
✓ No cookies
✓ No password history
✓ No localStorage
✓ No password transmitted over the network
✓ Works offline after assets are loaded

Then link to SECURITY.md.

Your current README already says most of this, but bringing it into a visually obvious section would make the product much easier to trust.

21. Add a "How can I trust this?" section

This would differentiate PassForge from many random password-generator websites.

Something like:

## Can I trust PassForge?


PassForge is open source.


The password generation code is available for inspection,
and passwords are generated locally using the browser's
Web Crypto API.


There is no password-generation server.


You can also run the application locally yourself.

That's a very strong selling point.

22. Add an offline/PWA mode

This would be a very good future feature for this particular application.

Your specification already says the app should work completely offline after initial load.

You could go one step further:

PassForge
    ↓
Service Worker
    ↓
Cache HTML
Cache CSS
Cache JS
Cache WASM
Cache word list
    ↓
Offline password generation

Then users could literally disconnect from the internet and still use it.

That's a compelling privacy feature.

But:

I wouldn't make this your next change.

First get the security/CI foundation right.

23. Consider Subresource Integrity carefully

For a conventional JavaScript application, I'd normally recommend SRI for external scripts.

But your application has intentionally minimized external dependencies and uses local assets. That's already better.

Don't introduce CDN dependencies merely to make the project look modern.

Your current "everything local" philosophy is the right one.

24. GitHub Pages itself is not the security boundary

This is worth emphasizing.

Your security model is:

GitHub Pages
      ↓
serves static files
      ↓
Browser
      ↓
Web Crypto
      ↓
password

GitHub isn't generating the password.

Therefore the key question isn't:

"Can I trust GitHub with my password?"

It's:

"Can I verify that the JavaScript/WASM running in my browser is the code I expect?"

That's why public source + reproducible builds + strong CI are valuable.

GitHub Pages supports HTTPS automatically for github.io sites, and GitHub documents that Pages sites are publicly accessible even in configurations where the repository itself can be private.

25. If you eventually want a custom domain

I would consider:

passforge.dev

or something similar if you decide to turn this into a serious portfolio/open-source project.

GitHub Pages supports custom domains and HTTPS.

But don't buy a domain yet unless you're committed to the name.

I'd first settle:

PassGen vs PassForge

and get the repository polished.

26. Your current project is actually a good portfolio project

Given your goal of having this as an open-source GitHub project, I'd lean into that.

The repository demonstrates:

Go
WebAssembly
cryptographic randomness
security engineering
browser APIs
CSP
GitHub Actions
GitHub Pages
testing
accessibility
privacy-by-design
static deployment
clean architecture

That's considerably more interesting than a basic password generator.

The architecture is already good enough that I'd polish it rather than rewrite it.

27. What I would NOT do

I would not:

❌ Make the repository private

There's little security benefit and it hurts transparency.

❌ Add a backend

Absolutely unnecessary.

❌ Add analytics

It conflicts with your privacy positioning.

❌ Add Firebase/Supabase/etc.

No reason.

❌ Add React

Your current Go/WASM architecture is a differentiator.

❌ Add a third-party password-strength API

That would undermine your privacy model.

❌ Store password history

Your current decision not to do this is correct.

❌ Add unnecessary dependencies

Your go.mod being dependency-free is excellent.

28. My recommended roadmap

I'd do this in roughly this order.

Phase 1 — Security foundation

Do now

 Web Crypto
 no backend
 CSP
 secure shuffle
 minimum-character algorithm
 SECURITY.md
 Merge security workflow PR
 Add CI workflow
 Add go vet
 Add fuzz/property tests
 Add privacy regression tests
 Enable Dependabot
 Protect main
Phase 2 — Repository polish

Next

 Repository description
 Website URL
 Topics
 README badges
 Trust & Privacy section
 Screenshots
 Better installation instructions
 Better contribution instructions
 Add version/release process
 Decide PassGen vs PassForge
Phase 3 — Product polish

After that

 Better entropy explanation
 PIN security warning
 Improve ambiguous-character UX
 Better mobile UX
 More accessibility testing
 Offline/PWA support
 Copy-to-clipboard security UX
Phase 4 — Open-source growth

Only if you want the project to grow

 GitHub Discussions
 Issue templates
 Feature request template
 Security disclosure process
 Releases
 Changelog
 Demo screenshot/GIF
 Custom domain
29. One thing I would do differently from your current approach

Your repository currently has a very large Specs.md—about 1,900 lines / 31.7 KB.

That's useful as an AI-agent implementation specification, but it's not ideal as the main open-source project's product documentation.

I'd separate:

docs/
   architecture.md
   security.md
   threat-model.md
   development.md

and keep Specs.md either:

docs/specification.md

or remove it once implementation has stabilized.

The public repository should primarily tell a human contributor:

What is this?
Why should I trust it?
How does it work?
How do I run it?
How do I contribute?

rather than exposing the entire AI implementation specification.

30. My final recommendation
Keep the repository public.

For this specific project, I think public is actually the correct choice.

Your ideal architecture is:

                 PUBLIC
                   │
        ┌──────────▼──────────┐
        │ GitHub Repository   │
        │                     │
        │ Go source           │
        │ WASM                │
        │ Tests               │
        │ Security policy     │
        └──────────┬──────────┘
                   │
              GitHub Actions
                   │
              build + test
                   │
                   ▼
            GitHub Pages
                   │
                   ▼
              User Browser
                   │
            crypto.getRandomValues()
                   │
                   ▼
              PASSWORD
                   │
              stays local

That is a much better security story than hiding the source.

The biggest gap isn't secrecy—it is verification and project maturity.

I'd prioritize:

CI + security scanning → stronger tests → branch protection → repository polish → product polish.

And I would merge your current security workflow PR after cleaning/squashing it, rather than keeping it open indefinitely. GitHub currently shows that PR as the only open PR, with six commits and no review yet.

One final observation: the repository was only created/updated around August 12, 2026, has two commits on main, zero stars/forks, and is therefore at a perfect stage to make these structural changes without worrying about breaking an established community.

Overall: I would keep the project public and invest in making it a demonstrably secure, well-tested open-source utility rather than trying to hide the implementation.