.PHONY: run build test clean help

help:
	@echo "Available commands:"
	@echo "  make run     - Run the application"
	@echo "  make build   - Build the application"
	@echo "  make test    - Run tests"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make fmt     - Format code"
	@echo "  make lint    - Run linter"

run:
	go run cmd/api/main.go

build:
	go build -o bin/products-api cmd/api/main.go

test:
	go test -v -cover ./...

clean:
	rm -rf bin/
	go clean

fmt:
	go fmt ./...
	goimports -w .

lint:
	golangci-lint run

dev:
	air
