.PHONY: build dev test clean

build:
	@chmod +x scripts/build.sh
	@./scripts/build.sh

dev: build
	@echo "Serving dist/ at http://localhost:$${PORT:-8080}"
	@cd dist && python3 ../scripts/serve.py

test:
	@go test ./...

clean:
	@rm -rf dist
