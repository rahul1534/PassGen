#!/usr/bin/env bash
# Privacy regression checks for PassForge.
#
# Fails if storage/network/telemetry APIs or third-party resources appear in
# application source, or if the Content-Security-Policy drifts from the
# approved policy below.
#
# Design notes:
#  * Uses plain grep only (no ripgrep dependency).
#  * Fails CLOSED: a grep error (exit >= 2) is a failure, never a pass.
#  * All keyword searches are case-insensitive.
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

SCAN_DIRS=(web cmd internal)
INDEX_HTML="web/index.html"
fail=0

# search <ERE-pattern>: sets MATCHES; returns 0 if any match, 1 if none.
# Aborts the whole script if grep itself errors, so a broken check cannot pass.
search() {
  MATCHES=$(grep -rnIiE --exclude='*.md' --exclude-dir=node_modules -e "$1" "${SCAN_DIRS[@]}" 2>&1)
  local rc=$?
  if [[ $rc -gt 1 ]]; then
    echo "ERROR: grep failed (exit $rc) while checking: $1"
    echo "$MATCHES"
    exit 2
  fi
  return $((rc))
}

check_absent() {
  local pattern="$1" label="$2"
  if search "$pattern"; then
    echo "FAIL: found forbidden pattern ($label): $pattern"
    echo "$MATCHES"
    fail=1
  else
    echo "OK: $label"
  fi
}

echo "Running privacy regression checks..."

# --- Storage / network / telemetry APIs -------------------------------------
check_absent 'localStorage'               'localStorage'
check_absent 'sessionStorage'             'sessionStorage'
check_absent 'indexedDB'                  'indexedDB'
check_absent 'navigator\.sendBeacon'      'sendBeacon'
check_absent 'XMLHttpRequest'             'XMLHttpRequest'
check_absent 'WebSocket'                  'WebSocket'
check_absent 'analytics|telemetry'        'analytics/telemetry keywords'
check_absent 'gtag\(|google-analytics|googletagmanager' 'Google analytics/tag manager'
check_absent 'cloudflareinsights|data-cf-beacon|beacon\.min\.js' 'Cloudflare Web Analytics beacon'
check_absent 'plausible|mixpanel|hotjar|sentry|fullstory|amplitude|posthog|matomo|segment\.(com|io)' 'other analytics/monitoring vendors'

# --- Third-party resources ---------------------------------------------------
# No tag may load a script/style/image/frame/etc. from another origin.
check_absent "<(script|link|iframe|img|embed|object|source|audio|video)\\b[^>]*\\b(src|href|data)=[\"']?(https?:)?//" \
  'external resource references in HTML'
check_absent '@import|url\(["'"'"']?(https?:)?//' 'external references in CSS'

# --- fetch() is only allowed in the WASM bootstrap loader (web/app.js) --------
if search 'fetch\('; then
  other=$(printf '%s\n' "$MATCHES" | grep -v '^web/app\.js:' || true)
  if [[ -n "$other" ]]; then
    echo "FAIL: unexpected fetch() outside web/app.js"
    echo "$other"
    fail=1
  else
    echo "OK: fetch limited to WASM bootstrap"
  fi
else
  echo "OK: fetch limited to WASM bootstrap"
fi

# --- Password logging heuristics -------------------------------------------
if search 'log\.(Print|Fatal|Panic).*password|console\.(log|debug|info).*password'; then
  echo "FAIL: possible password logging"
  echo "$MATCHES"
  fail=1
else
  echo "OK: no password logging patterns"
fi

# --- Content-Security-Policy must match the approved policy exactly ---------
# Changing the CSP is a deliberate security decision: update this list, the
# threat model and SECURITY.md together, in the same change.
EXPECTED_CSP=(
  "default-src 'self'"
  "script-src 'self' 'wasm-unsafe-eval'"
  "style-src 'self' 'unsafe-inline'"
  "img-src 'self' data:"
  "connect-src 'self'"
  "font-src 'self'"
  "object-src 'none'"
  "base-uri 'none'"
  "form-action 'none'"
)

# Normalise a directive's value: split on whitespace, sort tokens, rejoin.
norm() { printf '%s\n' "$1" | tr -s ' \t' '\n' | sed '/^$/d' | sort | tr '\n' ' ' | sed 's/ $//'; }

csp_meta=$(grep -oiE '<meta[^>]*http-equiv="Content-Security-Policy"[^>]*>' "$INDEX_HTML" || true)
csp_count=$(printf '%s' "$csp_meta" | grep -c . || true)
if [[ "$csp_count" -ne 1 ]]; then
  echo "FAIL: expected exactly one CSP <meta> in $INDEX_HTML, found $csp_count"
  fail=1
else
  csp=$(printf '%s' "$csp_meta" | sed -E 's/.*content="([^"]*)".*/\1/')
  csp_ok=1

  # Every expected directive must be present with exactly the approved sources.
  for entry in "${EXPECTED_CSP[@]}"; do
    name="${entry%% *}"; want="${entry#* }"
    got=$(printf '%s' "$csp" | tr ';' '\n' | sed -E 's/^[[:space:]]+//; s/[[:space:]]+$//' | awk -v n="$name" '$1==n { $1=""; sub(/^ /,""); print; exit }')
    if [[ "$(norm "$got")" != "$(norm "$want")" ]]; then
      echo "FAIL: CSP directive '$name' is \"$got\", approved value is \"$want\""
      csp_ok=0
    fi
  done

  # No directives beyond the approved set.
  want_names=$(printf '%s\n' "${EXPECTED_CSP[@]}" | awk '{print $1}' | sort)
  got_names=$(printf '%s' "$csp" | tr ';' '\n' | awk 'NF{print $1}' | sort)
  if [[ "$want_names" != "$got_names" ]]; then
    echo "FAIL: CSP directive set differs from the approved set"
    echo "  approved: $(echo $want_names)"
    echo "  found:    $(echo $got_names)"
    csp_ok=0
  fi

  # Belt and braces: never allow remote hosts or 'unsafe-eval' anywhere.
  if printf '%s' "$csp" | grep -qiE 'https?:|(^|[[:space:]])\*|(^|[[:space:]])'"'"'unsafe-eval'"'"; then
    echo "FAIL: CSP contains a remote host, wildcard, or 'unsafe-eval'"
    csp_ok=0
  fi

  if [[ "$csp_ok" -eq 1 ]]; then
    echo "OK: CSP matches the approved policy"
  else
    fail=1
  fi
fi

if [[ "$fail" -ne 0 ]]; then
  echo "Privacy checks failed."
  exit 1
fi

echo "All privacy checks passed."
