BINARY := caching-proxy
PORT   ?= 3000
ORIGIN ?= http://dummyjson.com

.DEFAULT_GOAL := help

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'

## build: compile the binary
build:
	go build -o $(BINARY) .

## run: start the proxy (override with PORT=8080 ORIGIN=https://api.example.com)
run:
	go run . --port $(PORT) --origin $(ORIGIN)

## fmt: format the code
fmt:
	gofumpt -w .

## lint: vet + golangci-lint
lint:
	go vet ./...
	golangci-lint run ./...

## test: run tests
test:
	go test -race ./...

## check: verify formatting, then lint and test (run before committing)
check:
	@test -z "$$(gofumpt -l .)" || { echo "needs formatting:"; gofumpt -l .; exit 1; }
	go vet ./...
	golangci-lint run ./...
	go test -race ./...

## tidy: sync go.mod with imports
tidy:
	go mod tidy

## clean: remove build artifacts
clean:
	rm -f $(BINARY)
	go clean -testcache

.PHONY: help build run fmt lint test check tidy clean
