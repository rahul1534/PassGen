.PHONY: build dev test clean privacy

build:
	@chmod +x scripts/build.sh
	@./scripts/build.sh

dev: build
	@go run ./cmd/devserver -dir dist -port $${PORT:-8080}

test:
	@go test ./...

privacy:
	@chmod +x scripts/privacy-check.sh
	@./scripts/privacy-check.sh

clean:
	@rm -rf dist
