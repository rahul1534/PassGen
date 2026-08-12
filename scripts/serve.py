#!/usr/bin/env python3
"""Static file server with correct WebAssembly MIME type."""

import os
import socket
import sys
from http.server import HTTPServer, SimpleHTTPRequestHandler
import mimetypes

mimetypes.add_type("application/wasm", ".wasm")

DEFAULT_PORT = 8080


class Handler(SimpleHTTPRequestHandler):
    extensions_map = {
        **getattr(SimpleHTTPRequestHandler, "extensions_map", {}),
        ".wasm": "application/wasm",
    }


class ReusableHTTPServer(HTTPServer):
    allow_reuse_address = True

    def server_bind(self):
        if hasattr(socket, "SO_REUSEADDR"):
            self.socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        super().server_bind()


def main():
    port = int(os.environ.get("PORT", DEFAULT_PORT))
    server = ReusableHTTPServer(("localhost", port), Handler)
    print(f"Serving HTTP on http://localhost:{port}/", flush=True)
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nStopped.", flush=True)
    finally:
        server.server_close()


if __name__ == "__main__":
    try:
        main()
    except OSError as err:
        port = os.environ.get("PORT", DEFAULT_PORT)
        if err.errno == 48:
            print(
                f"Error: port {port} is already in use.\n"
                f"Stop the other server with: lsof -ti :{port} | xargs kill\n"
                f"Or use another port: PORT=8081 make dev",
                file=sys.stderr,
            )
        else:
            print(f"Error: {err}", file=sys.stderr)
        sys.exit(1)
