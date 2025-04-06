#!/bin/bash
set -euo pipefail

# Determine which directories to lint
# Always include the root directory (for main.go and other root-level Go files)
ROOT_DIR="."

if [ -d "./cmd" ] || [ -d "./internal" ]; then
  # Include both root dir and cmd/internal directories
  DIRS="$ROOT_DIR ./cmd/... ./internal/..."
else
  # If no cmd/internal dirs, just check the root directory
  DIRS="$ROOT_DIR"
fi

# Run go vet first
echo "Running go vet..."
go vet $DIRS >/dev/null

# Create a temporary YAML configuration file with .yml extension
TMP_CONFIG=$(mktemp -t golangci-XXXXXX.yml)

cat > "$TMP_CONFIG" << 'EOF'
run:
  timeout: 5m
  skip-dirs-use-default: true
  skip-dirs:
    - vendor
    - third_party
    - node_modules
    - .git
    - "pkg/mod"
  skip-files:
    - ".*_test.go"
  allow-parallel-runners: true
  modules-download-mode: readonly
  
linters:
  disable-all: true
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
  
linters-settings:
  goimports:
    local-prefixes: github.com/AccursedGalaxy/noidea
EOF

# Run linting with our custom config - use --path-prefix to only lint project code
echo "Linting project code..."
golangci-lint run \
  --config="$TMP_CONFIG" \
  --path-prefix="github.com/AccursedGalaxy/noidea" \
  $DIRS

# Clean up temp file
rm -f "$TMP_CONFIG"

echo "✅ Linting completed successfully" 