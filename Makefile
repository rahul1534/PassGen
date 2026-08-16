.PHONY: build dev test clean privacy

build:
	@chmod +x scripts/build.sh
	@./scripts/build.sh

dev: build
	@echo "Serving dist/ at http://localhost:$${PORT:-8080}"
	@cd dist && python3 ../scripts/serve.py

test:
	@go test ./...

privacy:
	@chmod +x scripts/privacy-check.sh
	@./scripts/privacy-check.sh

clean:
	@rm -rf dist
