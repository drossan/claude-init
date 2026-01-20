# Makefile para claude-init CLI
#
# Este Makefile proporciona comandos convenientes para testing, building, linting, y más.

# Variables
BINARY_NAME=claude-init
MAIN_PATH=./main.go
BUILD_DIR=bin
DIST_DIR=dist
COVER_FILE=coverage.out
COVER_HTML=coverage.html

# Variables de versión
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")

# Variables de Go
GO=go
GOFLAGS=-v
LDFLAGS=-s -w -X "github.com/danielrossellosanchez/claude-init/cmd/version.Version=$(VERSION)" -X "github.com/danielrossellosanchez/claude-init/cmd/version.Commit=$(COMMIT)" -X "github.com/danielrossellosanchez/claude-init/cmd/version.BuildDate=$(BUILD_DATE)"

# Plataformas para build multi-plataforma
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Variables de testing
TEST_TIMEOUT=10m
TEST_PKGS=./...
TEST_FLAGS=-timeout=$(TEST_TIMEOUT) -v

# Variables de coverage
COVER_FLAGS=-coverprofile=$(COVER_FILE) -covermode=atomic

# Colores para output
BLUE=\033[0;34m
GREEN=\033[0;32m
RED=\033[0;31m
YELLOW=\033[0;33m
NC=\033[0m # No Color

# Objetivo por defecto
.PHONY: all
all: test lint build
	@echo "$(GREEN)✓ All checks passed!$(NC)"

## ============================================================================
## Comandos de Build
## ============================================================================

.PHONY: build
build: clean
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

.PHONY: build-local
build-local:
	@echo "$(BLUE)Building $(BINARY_NAME) in current directory...$(NC)"
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete: ./$(BINARY_NAME)$(NC)"

.PHONY: build-all
build-all:
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@echo "$(BLUE)Version: $(VERSION)$(NC)"
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		goos=$$(echo $$platform | cut -d'/' -f1); \
		goarch=$$(echo $$platform | cut -d'/' -f2); \
		binary_name=$(BINARY_NAME)-$$goos-$$goarch; \
		if [ "$$goos" = "windows" ]; then \
			binary_name=$$binary_name.exe; \
		fi; \
		binary_path=$(BUILD_DIR)/$$binary_name; \
		echo "  Building $$goos/$$goarch..."; \
		GOOS=$$goos GOARCH=$$goarch $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $$binary_path $(MAIN_PATH) || exit 1; \
		echo "    ✓ $$binary_path"; \
	done
	@echo "$(GREEN)✓ Multi-platform build complete$(NC)"
	@echo "$(BLUE)Binaries:$(NC)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)-*

.PHONY: clean
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVER_FILE) $(COVER_HTML)
	@rm -f $(BINARY_NAME)
	@echo "$(GREEN)✓ Clean complete$(NC)"

.PHONY: clean-all
clean-all: clean
	@echo "$(BLUE)Cleaning distribution artifacts...$(NC)"
	@rm -rf $(DIST_DIR)
	@echo "$(GREEN)✓ All artifacts cleaned$(NC)"

## ============================================================================
## Comandos de Testing
## ============================================================================

.PHONY: test
test:
	@echo "$(BLUE)Running tests...$(NC)"
	$(GO) test $(TEST_FLAGS) $(TEST_PKGS)
	@echo "$(GREEN)✓ Tests passed$(NC)"

.PHONY: test-short
test-short:
	@echo "$(BLUE)Running short tests...$(NC)"
	$(GO) test -short $(TEST_FLAGS) $(TEST_PKGS)
	@echo "$(GREEN)✓ Short tests passed$(NC)"

.PHONY: test-race
test-race:
	@echo "$(BLUE)Running tests with race detector...$(NC)"
	$(GO) test -race $(TEST_FLAGS) $(TEST_PKGS)
	@echo "$(GREEN)✓ Race detector tests passed$(NC)"

.PHONY: test-cover
test-cover:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GO) test $(COVER_FLAGS) $(TEST_FLAGS) $(TEST_PKGS)
	@echo "$(BLUE)Generating coverage report...$(NC)"
	$(GO) tool cover -html=$(COVER_FILE) -o $(COVER_HTML)
	@echo "$(GREEN)✓ Coverage report generated: $(COVER_HTML)$(NC)"
	@echo "$(BLUE)Total coverage:$(NC)"
	@$(GO) tool cover -func=$(COVER_FILE) | grep total

.PHONY: test-integration
test-integration: build-local
	@echo "$(BLUE)Running integration tests...$(NC)"
	$(GO) test -tags=integration $(TEST_FLAGS) ./tests/integration/...
	@echo "$(GREEN)✓ Integration tests passed$(NC)"

.PHONY: test-all
test-all: test test-race test-integration
	@echo "$(GREEN)✓ All test suites passed$(NC)"

## ============================================================================
## Comandos de Linting
## ============================================================================

.PHONY: lint
lint:
	@echo "$(BLUE)Running linters...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Linting passed$(NC)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not found. Installing...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Linting passed$(NC)"; \
	fi

.PHONY: lint-fix
lint-fix:
	@echo "$(BLUE)Running linters with auto-fix...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix ./...; \
		echo "$(GREEN)✓ Linting with fixes passed$(NC)"; \
	else \
		echo "$(RED)✗ golangci-lint not found$(NC)"; \
		exit 1; \
	fi

.PHONY: fmt
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "$(YELLOW)⚠ goimports not found. Installing...$(NC)"; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		goimports -w .; \
	fi
	@echo "$(GREEN)✓ Code formatted$(NC)"

.PHONY: vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GO) vet ./...
	@echo "$(GREEN)✓ Vet passed$(NC)"

## ============================================================================
## Comandos de Benchmarks
## ============================================================================

.PHONY: benchmark
benchmark:
	@echo "$(BLUE)Running benchmarks...$(NC)"
	$(GO) test -bench=. -benchmem -run=^$$ $(TEST_PKGS)
	@echo "$(GREEN)✓ Benchmarks complete$(NC)"

.PHONY: benchmark-cpu
benchmark-cpu:
	@echo "$(BLUE)Running CPU profiling benchmarks...$(NC)"
	$(GO) test -bench=. -cpuprofile=cpu.prof -run=^$$ $(TEST_PKGS)
	@echo "$(GREEN)✓ CPU profiling complete: cpu.prof$(NC)"

.PHONY:benchmark-mem
benchmark-mem:
	@echo "$(BLUE)Running memory profiling benchmarks...$(NC)"
	$(GO) test -bench=. -memprofile=mem.prof -run=^$$ $(TEST_PKGS)
	@echo "$(GREEN)✓ Memory profiling complete: mem.prof$(NC)"

## ============================================================================
## Comandos de Dependencias
## ============================================================================

.PHONY: deps
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

.PHONY: deps-update
deps-update:
	@echo "$(BLUE)Updating dependencies...$(NC)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

.PHONY: deps-verify
deps-verify:
	@echo "$(BLUE)Verifying dependencies...$(NC)"
	$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies verified$(NC)"

## ============================================================================
## Comandos de Release
## ============================================================================

.PHONY: release
release: clean-all build-all
	@echo "$(BLUE)Creating release $(VERSION)...$(NC)"
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		goos=$$(echo $$platform | cut -d'/' -f1); \
		goarch=$$(echo $$platform | cut -d'/' -f2); \
		binary_name=$(BINARY_NAME)-$$goos-$$goarch; \
		if [ "$$goos" = "windows" ]; then \
			binary_name=$$binary_name.exe; \
		fi; \
		binary_path=$(BUILD_DIR)/$$binary_name; \
		archive_name=$(BINARY_NAME)-$(VERSION)-$$goos-$$goarch.tar.gz; \
		archive_path=$(DIST_DIR)/$$archive_name; \
		echo "  Creating $$archive_name..."; \
		tar -czf $$archive_path -C $(BUILD_DIR) $$binary_name; \
		echo "    ✓ $$archive_path"; \
	done
	@echo "$(BLUE)Generating checksums...$(NC)"
	@cd $(DIST_DIR) && \
		for file in $(BINARY_NAME)-*.tar.gz; do \
			if [ -f "$$file" ]; then \
				openssl sha256 -r "$$file" | awk '{print $$1}' > "$$file.sha256"; \
				echo "  ✓ $$file.sha256"; \
			fi; \
		done
	@echo "$(GREEN)✓ Release $(VERSION) ready in $(DIST_DIR)$(NC)"
	@echo "$(BLUE)Release artifacts:$(NC)"
	@ls -lh $(DIST_DIR)/

.PHONY: release-checksums
release-checksums:
	@echo "$(BLUE)Generating checksums for existing release...$(NC)"
	@if [ ! -d "$(DIST_DIR)" ]; then \
		echo "$(RED)✗ Distribution directory not found. Run 'make release' first.$(NC)"; \
		exit 1; \
	fi
	@cd $(DIST_DIR) && \
		for file in $(BINARY_NAME)-*.tar.gz; do \
			if [ -f "$$file" ]; then \
				openssl sha256 -r "$$file" | awk '{print $$1}' > "$$file.sha256"; \
				echo "✓ $$file.sha256"; \
			fi; \
		done
	@echo "$(GREEN)✓ Checksums generated$(NC)"

.PHONY: verify-checksums
verify-checksums:
	@echo "$(BLUE)Verifying checksums...$(NC)"
	@if [ ! -d "$(DIST_DIR)" ]; then \
		echo "$(RED)✗ Distribution directory not found.$(NC)"; \
		exit 1; \
	fi
	@cd $(DIST_DIR) && \
		for sum in *.sha256; do \
			if [ -f "$$sum" ]; then \
				openssl dgst -sha256 -verify "$$sum" "$${sum%.sha256}" || exit 1; \
			fi; \
		done
	@echo "$(GREEN)✓ All checksums verified$(NC)"

## ============================================================================
## Comandos de Instalación
## ============================================================================

.PHONY: install
install: build
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	$(GO) install $(GOFLAGS) -ldflags "$(LDFLAGS)"
	@echo "$(GREEN)✓ Installation complete$(NC)"

.PHONY: install-tools
install-tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	@echo "$(BLUE)Installing golangci-lint...$(NC)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(BLUE)Installing goimports...$(NC)"
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)✓ Development tools installed$(NC)"

## ============================================================================
## Comandos de Utilidad
## ============================================================================

.PHONY: run
run: build-local
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	./$(BINARY_NAME)

.PHONY: check
check: fmt vet lint test
	@echo "$(GREEN)✓ All checks passed!$(NC)"

.PHONY: ci
ci: deps fmt vet lint test-race test-cover
	@echo "$(GREEN)✓ CI checks passed!$(NC)"

.PHONY: help
help:
	@echo "$(BLUE)claude-init CLI - Available commands:$(NC)"
	@echo ""
	@echo "$(GREEN)Build:$(NC)"
	@echo "  make build          - Build the binary for current platform"
	@echo "  make build-local    - Build in current directory"
	@echo "  make build-all      - Build for all platforms"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make clean-all      - Clean all artifacts (build + dist)"
	@echo ""
	@echo "$(GREEN)Release:$(NC)"
	@echo "  make release        - Create release (build-all + archives + checksums)"
	@echo "  make release-checksums - Generate checksums for existing release"
	@echo "  make verify-checksums  - Verify all checksums"
	@echo ""
	@echo "$(GREEN)Testing:$(NC)"
	@echo "  make test           - Run all tests"
	@echo "  make test-short     - Run short tests"
	@echo "  make test-race      - Run with race detector"
	@echo "  make test-cover     - Run with coverage"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-all       - Run all test suites"
	@echo ""
	@echo "$(GREEN)Linting:$(NC)"
	@echo "  make lint           - Run linters"
	@echo "  make lint-fix       - Run linters with auto-fix"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo ""
	@echo "$(GREEN)Benchmarks:$(NC)"
	@echo "  make benchmark      - Run benchmarks"
	@echo "  make benchmark-cpu  - Run CPU profiling"
	@echo "  make benchmark-mem  - Run memory profiling"
	@echo ""
	@echo "$(GREEN)Dependencies:$(NC)"
	@echo "  make deps           - Install dependencies"
	@echo "  make deps-update    - Update dependencies"
	@echo "  make deps-verify    - Verify dependencies"
	@echo ""
	@echo "$(GREEN)Installation:$(NC)"
	@echo "  make install        - Install to GOBIN"
	@echo "  make install-tools  - Install dev tools"
	@echo ""
	@echo "$(GREEN)Utility:$(NC)"
	@echo "  make run            - Build and run"
	@echo "  make check          - Run all checks"
	@echo "  make ci             - Run CI checks"
	@echo "  make help           - Show this help"
	@echo ""
	@echo "$(YELLOW)Version info:$(NC)"
	@echo "  VERSION=$(VERSION)"
	@echo "  COMMIT=$(COMMIT)"
	@echo "  BUILD_DATE=$(BUILD_DATE)"
	@echo ""

# Phony targets para evitar conflictos con archivos
.PHONY: all build build-local build-all clean clean-all test test-short test-race test-cover test-integration test-all lint lint-fix fmt vet benchmark benchmark-cpu benchmark-mem deps deps-update deps-verify install install-tools run check ci help release release-checksums verify-checksums
