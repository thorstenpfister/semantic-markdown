.PHONY: help test test-verbose test-short bench bench-verbose coverage coverage-html lint fmt vet build build-cli install clean parity golden all pre-release release-build release-archives release-checksums release-validate release-notes homebrew-formula homebrew-update release github-release clean-release

# Version and release configuration
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X 'github.com/thorstenpfister/semantic-markdown/cmd/semantic-md/cmd.Version=$(VERSION)' -X 'github.com/thorstenpfister/semantic-markdown/cmd/semantic-md/cmd.GitCommit=$(GIT_COMMIT)' -X 'github.com/thorstenpfister/semantic-markdown/cmd/semantic-md/cmd.BuildDate=$(BUILD_DATE)'

# Release directories
DIST_DIR := dist
RELEASE_DIR := $(DIST_DIR)/release
TARBALL_DIR := $(DIST_DIR)/tarballs

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Testing
test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./test

test-verbose: ## Run tests with verbose output
	@go test -v ./test -count=1

test-short: ## Run tests in short mode
	@go test -short ./test

test-all: ## Run all tests including in all packages
	@go test -v ./...

# Benchmarking
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem -run=^$$ ./test

bench-verbose: ## Run benchmarks with verbose output
	@go test -bench=. -benchmem -benchtime=3s -run=^$$ ./test

bench-cpu: ## Run CPU benchmarks
	@go test -bench=. -benchmem -cpuprofile=cpu.prof -run=^$$ ./test

bench-mem: ## Run memory benchmarks
	@go test -bench=. -benchmem -memprofile=mem.prof -run=^$$ ./test

# Coverage
coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -cover ./test

coverage-html: ## Generate HTML coverage report
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./test
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

coverage-func: ## Show coverage per function
	@go test -coverprofile=coverage.out ./test
	@go tool cover -func=coverage.out

# Parity and Golden Tests
parity: ## Run parity tests only
	@go test -v ./test -run TestParity

golden: ## Run golden file tests only
	@go test -v ./test -run TestGoldenFiles

# Code Quality
lint: ## Run linters (requires golangci-lint)
	@echo "Running linters..."
	@golangci-lint run --timeout=5m

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Building
build: ## Build library
	@echo "Building library..."
	@go build .

build-cli: ## Build CLI binary
	@echo "Building CLI..."
	@mkdir -p bin
	@go build -o bin/semantic-md ./cmd/semantic-md
	@echo "Binary created: bin/semantic-md"

build-all: build build-cli ## Build library and CLI

install: ## Install CLI binary
	@echo "Installing CLI..."
	@go install ./cmd/semantic-md

# Pre-release checks
pre-release: clean deps-verify lint test coverage ## Run all pre-release checks
	@echo "Pre-release checks passed!"
	@echo "Current version: $(VERSION)"
	@if [ "$(VERSION)" = "dev" ] || echo "$(VERSION)" | grep -q dirty; then \
		echo "Warning: Not on a clean tagged commit"; \
		echo "Create a tag first: git tag -a v1.0.0 -m 'Release v1.0.0'"; \
		exit 1; \
	fi

# Build release binaries with version info
release-build: ## Build all release binaries with version information
	@echo "Building release binaries for version $(VERSION)..."
	@mkdir -p $(RELEASE_DIR)

	@echo "Building Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(RELEASE_DIR)/semantic-md-linux-amd64 ./cmd/semantic-md

	@echo "Building Linux ARM64..."
	@GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(RELEASE_DIR)/semantic-md-linux-arm64 ./cmd/semantic-md

	@echo "Building macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(RELEASE_DIR)/semantic-md-darwin-amd64 ./cmd/semantic-md

	@echo "Building macOS ARM64..."
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(RELEASE_DIR)/semantic-md-darwin-arm64 ./cmd/semantic-md

	@echo "Building Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(RELEASE_DIR)/semantic-md-windows-amd64.exe ./cmd/semantic-md

	@echo "Release binaries built in $(RELEASE_DIR)/"

# Create distribution archives
release-archives: release-build ## Create tar.gz archives for distribution
	@echo "Creating distribution archives..."
	@mkdir -p $(TARBALL_DIR)

	@echo "Creating Linux AMD64 archive..."
	@tar -czf $(TARBALL_DIR)/semantic-md-$(VERSION)-linux-amd64.tar.gz -C $(RELEASE_DIR) semantic-md-linux-amd64

	@echo "Creating Linux ARM64 archive..."
	@tar -czf $(TARBALL_DIR)/semantic-md-$(VERSION)-linux-arm64.tar.gz -C $(RELEASE_DIR) semantic-md-linux-arm64

	@echo "Creating macOS AMD64 archive..."
	@tar -czf $(TARBALL_DIR)/semantic-md-$(VERSION)-darwin-amd64.tar.gz -C $(RELEASE_DIR) semantic-md-darwin-amd64

	@echo "Creating macOS ARM64 archive..."
	@tar -czf $(TARBALL_DIR)/semantic-md-$(VERSION)-darwin-arm64.tar.gz -C $(RELEASE_DIR) semantic-md-darwin-arm64

	@echo "Creating Windows AMD64 archive..."
	@cd $(RELEASE_DIR) && zip ../tarballs/semantic-md-$(VERSION)-windows-amd64.zip semantic-md-windows-amd64.exe

	@echo "Archives created in $(TARBALL_DIR)/"

# Generate checksums
release-checksums: release-archives ## Generate SHA256 checksums for all archives
	@echo "Generating checksums..."
	@cd $(TARBALL_DIR) && sha256sum *.tar.gz *.zip > checksums.txt
	@echo "Checksums generated:"
	@cat $(TARBALL_DIR)/checksums.txt

# Validate release artifacts
release-validate: release-checksums ## Validate release artifacts exist and checksums are correct
	@echo "Validating release artifacts..."
	@cd $(TARBALL_DIR) && sha256sum -c checksums.txt
	@echo "All release artifacts validated successfully!"

# Generate release notes from CHANGELOG
release-notes: ## Extract release notes for current version from CHANGELOG
	@echo "Extracting release notes for version $(VERSION)..."
	@if [ ! -f CHANGELOG.md ]; then \
		echo "Error: CHANGELOG.md not found"; \
		exit 1; \
	fi
	@mkdir -p $(DIST_DIR)
	@sed -n "/## \[$(VERSION:v%=%)\]/,/## \[/p" CHANGELOG.md | sed '$$d' | tail -n +2 > $(DIST_DIR)/release-notes.md
	@echo "Release notes extracted to $(DIST_DIR)/release-notes.md"

# Generate Homebrew formula
homebrew-formula: release-checksums ## Generate Homebrew formula file
	@echo "Generating Homebrew formula..."
	@mkdir -p $(DIST_DIR)/homebrew
	@AMD64_SHA256=$$(grep "darwin-amd64.tar.gz" $(TARBALL_DIR)/checksums.txt | awk '{print $$1}'); \
	ARM64_SHA256=$$(grep "darwin-arm64.tar.gz" $(TARBALL_DIR)/checksums.txt | awk '{print $$1}'); \
	LINUX_AMD64_SHA256=$$(grep "linux-amd64.tar.gz" $(TARBALL_DIR)/checksums.txt | awk '{print $$1}'); \
	LINUX_ARM64_SHA256=$$(grep "linux-arm64.tar.gz" $(TARBALL_DIR)/checksums.txt | awk '{print $$1}'); \
	echo "# typed: false" > $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "# frozen_string_literal: true" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "# This file was generated by semantic-markdown's Makefile" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "class SemanticMd < Formula" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  desc \"Convert HTML to clean, semantic Markdown optimized for LLMs\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  homepage \"https://github.com/thorstenpfister/semantic-markdown\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  version \"$(VERSION:v%=%)\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  license \"MIT\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  on_macos do" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    if Hardware::CPU.arm?" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "      url \"https://github.com/thorstenpfister/semantic-markdown/releases/download/$(VERSION)/semantic-md-$(VERSION)-darwin-arm64.tar.gz\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "      sha256 \"$$ARM64_SHA256\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    end" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    if Hardware::CPU.intel?" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "      url \"https://github.com/thorstenpfister/semantic-markdown/releases/download/$(VERSION)/semantic-md-$(VERSION)-darwin-amd64.tar.gz\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "      sha256 \"$$AMD64_SHA256\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    end" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  end" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  on_linux do" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "      url \"https://github.com/thorstenpfister/semantic-markdown/releases/download/$(VERSION)/semantic-md-$(VERSION)-linux-arm64.tar.gz\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "      sha256 \"$$LINUX_ARM64_SHA256\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    end" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    if Hardware::CPU.intel?" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "      url \"https://github.com/thorstenpfister/semantic-markdown/releases/download/$(VERSION)/semantic-md-$(VERSION)-linux-amd64.tar.gz\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "      sha256 \"$$LINUX_AMD64_SHA256\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    end" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  end" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  def install" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    bin.install \"semantic-md-darwin-arm64\" => \"semantic-md\" if OS.mac? && Hardware::CPU.arm?" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    bin.install \"semantic-md-darwin-amd64\" => \"semantic-md\" if OS.mac? && Hardware::CPU.intel?" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    bin.install \"semantic-md-linux-arm64\" => \"semantic-md\" if OS.linux? && Hardware::CPU.arm?" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    bin.install \"semantic-md-linux-amd64\" => \"semantic-md\" if OS.linux? && Hardware::CPU.intel?" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  end" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  test do" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    assert_match \"semantic-md version\", shell_output(\"#{bin}/semantic-md version\")" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    html = \"<h1>Test</h1><p>Content</p>\"" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    output = pipe_output(\"#{bin}/semantic-md convert\", html)" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "    assert_match \"# Test\", output" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "  end" >> $(DIST_DIR)/homebrew/semantic-md.rb; \
	echo "end" >> $(DIST_DIR)/homebrew/semantic-md.rb
	@echo "Homebrew formula generated at $(DIST_DIR)/homebrew/semantic-md.rb"

# Update Homebrew tap (requires tap repository to be cloned locally)
homebrew-update: homebrew-formula ## Update local Homebrew tap repository
	@echo "Updating Homebrew tap..."
	@if [ -z "$(TAP_DIR)" ]; then \
		echo "Error: TAP_DIR not set. Usage: make homebrew-update TAP_DIR=/path/to/homebrew-tap"; \
		exit 1; \
	fi
	@if [ ! -d "$(TAP_DIR)" ]; then \
		echo "Error: TAP_DIR '$(TAP_DIR)' does not exist"; \
		exit 1; \
	fi
	@mkdir -p $(TAP_DIR)/Formula
	@cp $(DIST_DIR)/homebrew/semantic-md.rb $(TAP_DIR)/Formula/
	@echo "Formula copied to $(TAP_DIR)/Formula/semantic-md.rb"
	@echo "Don't forget to commit and push the tap repository!"

# Complete release workflow
release: pre-release release-validate release-notes homebrew-formula ## Run complete release workflow
	@echo ""
	@echo "=========================================="
	@echo "Release $(VERSION) built successfully!"
	@echo "=========================================="
	@echo ""
	@echo "Release artifacts:"
	@ls -lh $(TARBALL_DIR)/
	@echo ""
	@echo "Next steps:"
	@echo "1. Review release notes in $(DIST_DIR)/release-notes.md"
	@echo "2. Create GitHub release: gh release create $(VERSION) $(TARBALL_DIR)/* --notes-file $(DIST_DIR)/release-notes.md"
	@echo "3. Update Homebrew tap: make homebrew-update TAP_DIR=/path/to/homebrew-tap"
	@echo "4. Push Homebrew tap changes"
	@echo ""

# GitHub release creation (requires gh CLI)
github-release: release ## Create GitHub release with artifacts
	@echo "Creating GitHub release..."
	@if ! command -v gh &> /dev/null; then \
		echo "Error: GitHub CLI (gh) not installed. Install with: brew install gh"; \
		exit 1; \
	fi
	@if [ ! -f $(DIST_DIR)/release-notes.md ]; then \
		echo "Error: Release notes not found. Run 'make release' first."; \
		exit 1; \
	fi
	@gh release create $(VERSION) $(TARBALL_DIR)/* \
		--title "Release $(VERSION)" \
		--notes-file $(DIST_DIR)/release-notes.md
	@echo "GitHub release created successfully!"
	@echo "View at: https://github.com/thorstenpfister/semantic-markdown/releases/tag/$(VERSION)"

# Legacy build-release for backward compatibility
build-release: release-build ## Legacy target - use release-build instead

# Cleaning
clean-release: ## Clean release build artifacts
	@echo "Cleaning release artifacts..."
	@rm -rf $(DIST_DIR)
	@echo "Release artifacts cleaned"

clean: clean-release ## Clean build artifacts and test cache
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html cpu.prof mem.prof
	@go clean -testcache
	@echo "Clean complete"

clean-all: clean ## Clean everything including go module cache
	@go clean -modcache

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

deps-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

deps-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

# Development
dev: fmt vet test ## Run formatter, vet, and tests

ci: deps-verify lint test coverage ## Run CI checks (lint, test, coverage)

# Quick checks
quick: test-short bench ## Run short tests and benchmarks quickly

# Complete validation
all: clean deps-verify fmt vet lint test coverage bench build-cli ## Run all checks and build
	@echo "All checks passed!"

# Examples
run-example-basic: ## Run basic example
	@go run examples/basic/main.go

run-example-metadata: ## Run metadata example
	@go run examples/metadata/main.go

# Documentation
doc: ## Open package documentation
	@echo "Opening documentation..."
	@open http://localhost:6060/pkg/github.com/thorstenpfister/semantic-markdown/
	@godoc -http=:6060

# Git helpers
tag: ## Create a new git tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make tag VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "Tag created. Push with: git push origin $(VERSION)"
