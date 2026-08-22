.PHONY: build dev test format setup-hooks clean privacy

build:
	@chmod +x scripts/build.sh
	@./scripts/build.sh

dev: build
	@go run ./cmd/devserver -dir dist -port $${PORT:-8080}

test:
	@go test ./...

format:
	@gofmt -w $$(git ls-files '*.go')

setup-hooks:
	@git config core.hooksPath .githooks
	@echo "Git hooks enabled from .githooks"

privacy:
	@chmod +x scripts/privacy-check.sh
	@./scripts/privacy-check.sh

clean:
	@rm -rf dist
