# Variables
BINARY_NAME=aphelion
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse HEAD)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILT_BY ?= $(shell whoami)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X github.com/Exmplr-AI/aphelion-cli/cmd.version=$(VERSION) \
                  -X github.com/Exmplr-AI/aphelion-cli/cmd.commit=$(COMMIT) \
                  -X github.com/Exmplr-AI/aphelion-cli/cmd.date=$(DATE) \
                  -X github.com/Exmplr-AI/aphelion-cli/cmd.builtBy=$(BUILT_BY)"

# Platforms
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: all build clean test deps help install

all: clean deps test build

## Build the binary for current platform
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) main.go

## Build for all platforms
build-all: clean deps
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_NAME=$(BINARY_NAME)-$$OS-$$ARCH; \
		if [ $$OS = "windows" ]; then OUTPUT_NAME=$$OUTPUT_NAME.exe; fi; \
		echo "Building $$OUTPUT_NAME..."; \
		GOOS=$$OS GOARCH=$$ARCH $(GOBUILD) $(LDFLAGS) -o dist/$$OUTPUT_NAME main.go; \
	done

## Install the binary to $GOPATH/bin
install: build
	cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

## Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf dist/

## Run tests
test:
	$(GOTEST) -v ./...

## Run tests with coverage
test-coverage:
	$(GOTEST) -v -cover ./...

## Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## Run the application (for development)
run: build
	./$(BINARY_NAME)

## Format code
fmt:
	$(GOCMD) fmt ./...

## Run linter
lint:
	golangci-lint run

## Initialize go modules
mod-init:
	$(GOMOD) init github.com/Exmplr-AI/aphelion-cli

## Show help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## Generate shell completions
completions: build
	@mkdir -p completions
	./$(BINARY_NAME) completion bash > completions/$(BINARY_NAME).bash
	./$(BINARY_NAME) completion zsh > completions/$(BINARY_NAME).zsh
	./$(BINARY_NAME) completion fish > completions/$(BINARY_NAME).fish
	./$(BINARY_NAME) completion powershell > completions/$(BINARY_NAME).ps1

## Create a new release (requires TAG variable)
release: clean deps test build-all
	@if [ -z "$(TAG)" ]; then echo "TAG is required. Usage: make release TAG=v1.0.0"; exit 1; fi
	@echo "Creating release $(TAG)..."
	git tag $(TAG)
	git push origin $(TAG)

.DEFAULT_GOAL := help