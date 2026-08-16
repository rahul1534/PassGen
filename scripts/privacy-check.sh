#!/usr/bin/env bash
# Privacy regression checks for PassForge.
# Fails if storage/network/telemetry APIs appear in application source.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

fail=0

check_absent() {
  local pattern="$1"
  local label="$2"
  local matches
  matches=$(rg -n --glob '!**/Review1.md' --glob '!**/Specs.md' --glob '!**/docs/**' --glob '!**/*.md' \
    -g '!scripts/privacy-check.sh' \
    -g '!scripts/smoke-test.mjs' \
    "$pattern" web/ cmd/ internal/ 2>/dev/null || true)
  if [[ -n "$matches" ]]; then
    echo "FAIL: found forbidden pattern ($label): $pattern"
    echo "$matches"
    fail=1
  else
    echo "OK: $label"
  fi
}

echo "Running privacy regression checks..."

check_absent 'localStorage' 'localStorage'
check_absent 'sessionStorage' 'sessionStorage'
check_absent 'indexedDB' 'indexedDB'
check_absent 'navigator\.sendBeacon' 'sendBeacon'
check_absent 'XMLHttpRequest' 'XMLHttpRequest'
check_absent 'WebSocket' 'WebSocket'
check_absent 'analytics' 'analytics'
check_absent 'telemetry' 'telemetry'
check_absent 'gtag\(' 'gtag'
check_absent 'google-analytics' 'google-analytics'

# fetch is only allowed in the WASM bootstrap loader (app.js).
fetch_hits=$(rg -n 'fetch\(' web/ cmd/ internal/ --glob '!scripts/**' 2>/dev/null || true)
allowed=$(echo "$fetch_hits" | rg 'web/app\.js:' || true)
other=$(echo "$fetch_hits" | rg -v 'web/app\.js:' || true)
if [[ -n "$other" ]]; then
  echo "FAIL: unexpected fetch() outside web/app.js"
  echo "$other"
  fail=1
else
  echo "OK: fetch limited to WASM bootstrap"
fi

# Password logging heuristics
log_hits=$(rg -n 'log\.(Print|Fatal|Panic).*[Pp]assword|console\.(log|debug|info).*password' \
  web/ cmd/ internal/ 2>/dev/null || true)
if [[ -n "$log_hits" ]]; then
  echo "FAIL: possible password logging"
  echo "$log_hits"
  fail=1
else
  echo "OK: no password logging patterns"
fi

if [[ "$fail" -ne 0 ]]; then
  echo "Privacy checks failed."
  exit 1
fi

echo "All privacy checks passed."
