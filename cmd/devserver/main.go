package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	defaultPort := 8080
	if env := os.Getenv("PORT"); env != "" {
		if p, err := strconv.Atoi(env); err == nil {
			defaultPort = p
		}
	}

	dir := flag.String("dir", "dist", "directory to serve")
	port := flag.Int("port", defaultPort, "TCP port to listen on")
	flag.Parse()

	abs, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatal(err)
	}
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		log.Fatalf("serve directory not found: %s (run make build first)", abs)
	}

	mux := http.NewServeMux()
	mux.Handle("/", wasmFileServer(http.Dir(abs)))

	addr := fmt.Sprintf("localhost:%d", *port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			fmt.Fprintf(os.Stderr, "Error: port %d is already in use.\n", *port)
			fmt.Fprintf(os.Stderr, "Stop the other server with: lsof -ti :%d | xargs kill\n", *port)
			fmt.Fprintf(os.Stderr, "Or use another port: PORT=8081 make dev\n")
			os.Exit(1)
		}
		log.Fatal(err)
	}

	fmt.Printf("Serving %s at http://%s/\n", abs, addr)
	if err := http.Serve(ln, mux); err != nil {
		log.Fatal(err)
	}
}

// wasmFileServer serves static files and sets application/wasm for .wasm assets.
func wasmFileServer(root http.FileSystem) http.Handler {
	fs := http.FileServer(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if filepath.Ext(r.URL.Path) == ".wasm" {
			w.Header().Set("Content-Type", "application/wasm")
		}
		fs.ServeHTTP(w, r)
	})
}
