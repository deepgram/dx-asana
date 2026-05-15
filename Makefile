.PHONY: help build run test smoke lint fmt clean install deps release

BINARY := dx-asana
MODULE := github.com/deepgram/dx-asana

VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -ldflags "\
	-X $(MODULE)/cmd.Version=$(VERSION) \
	-X $(MODULE)/cmd.Commit=$(COMMIT) \
	-X $(MODULE)/cmd.Date=$(BUILD_DATE)"

help:
	@echo "$(BINARY) - Deepgram DX Asana CLI with TUI and sync daemon"
	@echo ""
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  install     - Install the binary"
	@echo "  run         - Run the CLI"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linter"
	@echo "  fmt         - Format code"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Download dependencies"
	@echo "  release     - Build release binaries for all platforms"

build:
	go build $(LDFLAGS) -o $(BINARY) main.go

install:
	go install $(LDFLAGS)

run:
	go run $(LDFLAGS) main.go

test:
	go test -v ./...

smoke: build
	./scripts/smoke.sh ./$(BINARY)

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

clean:
	rm -f $(BINARY)
	go clean

deps:
	go mod download
	go mod tidy

dev-setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Build release binaries for all platforms
release: clean
	@echo "Building release binaries for version $(VERSION)"
	@mkdir -p dist
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 main.go
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 main.go
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 main.go
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe main.go
	@echo "Binaries built in dist/"
	@ls -lh dist/
