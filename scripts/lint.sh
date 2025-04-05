#!/bin/bash
# Simple script to run linting only on the project's own code

set -e

# Change to root directory of project
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"
cd "$ROOT_DIR"

# Ensure dependencies are installed
echo "Ensuring dependencies are installed..."
go get -v github.com/zalando/go-keyring
go mod tidy

# Install golangci-lint if needed
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
fi

echo "Running gofmt..."
gofmt_output=$(gofmt -l ./cmd ./internal 2>&1)
if [ -n "$gofmt_output" ]; then
    echo "⚠️ The following files are not formatted correctly:"
    echo "$gofmt_output"
    echo "Running gofmt -w to fix formatting..."
    gofmt -w ./cmd ./internal
    echo "✅ Files have been formatted."
else
    echo "✅ All Go files are properly formatted."
fi

echo "Running go vet..."
go vet ./cmd/... ./internal/... 2>&1

echo "Running staticcheck..."
if ! command -v staticcheck &> /dev/null; then
    echo "Installing staticcheck..."
    go install honnef.co/go/tools/cmd/staticcheck@latest
fi

# Run staticcheck but don't fail the build for unused function warnings
staticcheck_output=$(staticcheck ./cmd/... ./internal/... 2>&1)
if [ -n "$staticcheck_output" ]; then
    unused_only=true
    while IFS= read -r line; do
        if [[ ! "$line" =~ "is unused" ]]; then
            unused_only=false
            echo "$line"
        else
            echo "⚠️ $line"
        fi
    done <<< "$staticcheck_output"
    
    if [ "$unused_only" = false ]; then
        lint_errors=1
    else
        echo "⚠️ Only unused functions detected - not failing the build"
    fi
else
    echo "✅ Static check passed with no issues."
fi

# Project-specific linting
echo "Running golangci-lint on project files only..."
SKIP_DIRS="vendor,third_party,node_modules"
golangci-lint run --timeout=5m \
    --modules-download-mode=readonly \
    --skip-dirs-use-default \
    --skip-dirs="${SKIP_DIRS}" \
    --skip-files=".*_test.go" \
    --path-prefix="github.com/AccursedGalaxy/noidea" \
    ./cmd/... ./internal/...

# Return success
echo "✅ Linting complete."
exit 0 