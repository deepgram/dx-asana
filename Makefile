.PHONY: help build run test lint fmt clean install deps release

VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -ldflags "\
	-X github.com/TheCoolRobot/asana-cli/cmd.Version=$(VERSION) \
	-X github.com/TheCoolRobot/asana-cli/cmd.Commit=$(COMMIT) \
	-X github.com/TheCoolRobot/asana-cli/cmd.Date=$(BUILD_DATE)"

help:
	@echo "asana-cli - Asana CLI with TUI and sync daemon"
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
	go build $(LDFLAGS) -o asana-cli main.go

install:
	go install $(LDFLAGS)

run:
	go run $(LDFLAGS) main.go

test:
	go test -v ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

clean:
	rm -f asana-cli
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
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/asana-cli-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/asana-cli-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/asana-cli-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/asana-cli-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/asana-cli-windows-amd64.exe main.go
	@echo "Binaries built in dist/"
	@ls -lh dist/