# Makefile for gwt (Git Worktree Manager)

# Variables
BINARY_NAME=gwt
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build directories
DIST_DIR=dist
BUILD_DIR=build

.PHONY: all build clean test test-coverage test-race test-bench deps run install uninstall fmt lint lint-fix security mod-verify mod-tidy vet quality ci pre-commit

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

test-race:
	$(GOTEST) -v -race ./...

test-bench:
	$(GOTEST) -v -bench=. -benchmem ./...

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR) $(DIST_DIR) coverage.out coverage.html

deps:
	$(GOMOD) download
	$(GOMOD) tidy

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installing shell integration..."
	@echo "Add 'source $(PWD)/scripts/gwt.sh' to your shell rc file"

uninstall:
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Cross compilation
build-all:
	mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe

# Development helpers
dev-init:
	@echo "Creating example configuration..."
	@mkdir -p .gitworktree/templates
	@cp examples/configs/basic.yaml .gitworktree/config.yaml
	@cp templates/* .gitworktree/templates/
	@echo "Configuration created at .gitworktree/"

fmt:
	go fmt ./...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

security:
	gosec ./...

mod-verify:
	$(GOMOD) verify

mod-tidy:
	$(GOMOD) tidy

vet:
	$(GOCMD) vet ./...

# Quality checks
quality: fmt lint vet test

# CI pipeline
ci: deps mod-verify quality test-coverage

# Pre-commit checks
pre-commit: fmt lint-fix vet test-race

# Homebrew formula generation (for later)
homebrew-formula:
	@echo "Generating Homebrew formula..."
	@cat homebrew/gwt.rb.template | \
		sed 's/{{VERSION}}/$(VERSION)/g' | \
		sed 's/{{SHA256}}/$(shell shasum -a 256 $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 | cut -d' ' -f1)/g' \
		> homebrew/gwt.rb