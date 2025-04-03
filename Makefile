.PHONY: build install uninstall clean

# Binary name
BINARY=noidea
# Version from git tag or default to 'dev'
VERSION=$(shell git describe --tags 2>/dev/null || echo "dev")
# Build with version info
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Installation paths
PREFIX?=/usr/local
BINDIR?=$(PREFIX)/bin

# Cross-compilation targets
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Determine the real user (in case of sudo)
REAL_USER?=$(shell echo $${SUDO_USER:-$$USER})
REAL_HOME?=$(shell eval echo ~$(REAL_USER))

# Find Go
GO?=$(shell command -v go 2>/dev/null || \
       command -v /usr/local/go/bin/go 2>/dev/null || \
	   command -v /usr/bin/go 2>/dev/null || \
	   command -v $(REAL_HOME)/go/bin/go 2>/dev/null || \
	   echo "go")

# Default: build the binary
build:
	@echo "Building $(BINARY) version $(VERSION)..."
	@if ! $(GO) version >/dev/null 2>&1; then \
		echo "Error: Go not found. Please install Go or specify GO=/path/to/go"; \
		exit 1; \
	fi
	$(GO) build $(LDFLAGS) -o $(BINARY)
	@echo "✅ Build complete: $(BINARY)"

# Install binary and set up global config directory
install: build
	@echo "Installing $(BINARY) to $(BINDIR)..."
	mkdir -p $(BINDIR)
	cp $(BINARY) $(BINDIR)/
	mkdir -p $(REAL_HOME)/.noidea
	@# Fix ownership if running as root
	@if [ $$(id -u) -eq 0 ] && [ -n "$(SUDO_USER)" ]; then \
		chown -R $(SUDO_USER) $(REAL_HOME)/.noidea; \
	fi
	@echo "✅ Installation complete."
	@echo "Run 'noidea init' in any repository to set up git hooks."

# Create release binaries for multiple platforms
release:
	@echo "Building release binaries for:"
	@mkdir -p dist
	@rm -f dist/*
	@$(foreach platform,$(PLATFORMS),\
		echo "  - $(platform)"; \
		export GOOS=$$(echo $(platform) | cut -d/ -f1); \
		export GOARCH=$$(echo $(platform) | cut -d/ -f2); \
		export OUTPUT=dist/$(BINARY)_$${GOOS}_$${GOARCH}; \
		if [ "$${GOOS}" = "windows" ]; then export OUTPUT=$${OUTPUT}.exe; fi; \
		$(GO) build $(LDFLAGS) -o $${OUTPUT}; \
	)
	@echo "✅ Release builds complete. See dist/ directory."

# Install dependencies for development
deps:
	$(GO) get -v ./...

# Run tests
test:
	$(GO) test -v ./...

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY)..."
	rm -f $(BINDIR)/$(BINARY)
	@echo "✅ Uninstallation complete."

# Clean built binaries and artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
	rm -rf dist/
	@echo "✅ Clean complete."

# Show help
help:
	@echo "noidea Makefile"
	@echo ""
	@echo "Commands:"
	@echo "  make build      - Build noidea binary"
	@echo "  make install    - Install noidea to $(BINDIR)"
	@echo "  make uninstall  - Remove noidea from $(BINDIR)"
	@echo "  make release    - Build binaries for all platforms"
	@echo "  make clean      - Remove built binaries"
	@echo "  make deps       - Install dependencies"
	@echo "  make test       - Run tests"
	@echo ""
	@echo "Options:"
	@echo "  PREFIX=<path>   - Set installation prefix (default: /usr/local)"
	@echo "  GO=<path>       - Path to Go executable (default: auto-detected)" 