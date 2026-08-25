#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST="$ROOT/dist"

if [[ ! -d "$DIST" ]]; then
  echo "FAIL: production directory does not exist: $DIST"
  exit 1
fi

for path in "$DIST" "$ROOT/web" "$ROOT/cmd" "$ROOT/internal" "$ROOT/wordlist"; do
  # Node.js and browser-test tooling must remain development-only dependencies.
  if matches=$(grep -RInI -E '\b(node(js)?|npm|playwright)\b' "$path"); then
    echo "FAIL: development tooling reference found in production path: $path"
    echo "$matches"
    exit 1
  fi
done

echo "OK: production paths contain no Node.js or Playwright references"