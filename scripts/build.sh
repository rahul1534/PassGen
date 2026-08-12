#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST="$ROOT/dist"
GOROOT="$(go env GOROOT)"

mkdir -p "$DIST/assets"

echo "Building WebAssembly..."
GOOS=js GOARCH=wasm go build -o "$DIST/app.wasm" "$ROOT/cmd/passforge"

echo "Copying static assets..."
WASM_EXEC="$GOROOT/lib/wasm/wasm_exec.js"
if [[ ! -f "$WASM_EXEC" ]]; then
  WASM_EXEC="$GOROOT/misc/wasm/wasm_exec.js"
fi
cp "$WASM_EXEC" "$DIST/wasm_exec.js"
cp "$ROOT/web/index.html" "$DIST/index.html"
cp "$ROOT/web/app.js" "$DIST/app.js"
cp "$ROOT/web/styles.css" "$DIST/styles.css"
cp "$ROOT/web/favicon.svg" "$DIST/favicon.svg"

echo "Build complete: $DIST"
