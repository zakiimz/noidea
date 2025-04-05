#!/bin/bash
# Setup git hooks for the project

set -e

# Get the root directory of the project
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"
HOOKS_DIR="$ROOT_DIR/.git/hooks"

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

echo "Setting up git hooks..."

# Create pre-commit hook
cat > "$HOOKS_DIR/pre-commit" << 'EOF'
#!/bin/bash
# Pre-commit hook for noidea

# Exit on error
set -e

echo "Running pre-commit checks..."

# Stash any changes not added to the index
git stash -q --keep-index || true

# Ensure we clean up on exit
function cleanup {
  # Restore the stashed changes regardless of outcome
  git stash pop -q 2>/dev/null || true
}
trap cleanup EXIT

# Run go mod tidy to manage dependencies
echo "Checking dependencies..."
go mod tidy

# Run goimports to format and organize imports
if command -v goimports &> /dev/null; then
  echo "Formatting and organizing imports..."
  find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w -local github.com/AccursedGalaxy/noidea
else
  echo "Installing goimports..."
  go install golang.org/x/tools/cmd/goimports@latest
  find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w -local github.com/AccursedGalaxy/noidea
fi

# Format the code
make format

# Run lint
make script-lint

# If there are any changes after formatting, add them
if git diff --name-only | grep -q "\.go$"; then
  echo "Adding automatically formatted files..."
  git add $(git diff --name-only | grep "\.go$")
fi

echo "✅ Pre-commit checks passed."
EOF

# Create pre-push hook
cat > "$HOOKS_DIR/pre-push" << 'EOF'
#!/bin/bash
# Pre-push hook for noidea

# Exit on error
set -e

echo "Running pre-push checks..."

# Run tests
echo "Running tests..."
make test

echo "✅ Pre-push checks passed."
EOF

# Make the hooks executable
chmod +x "$HOOKS_DIR/pre-commit"
chmod +x "$HOOKS_DIR/pre-push"

echo "✅ Git hooks setup complete."
echo "The following hooks are installed:"
echo "- pre-commit: manages dependencies, formats code, organizes imports, and runs linting checks"
echo "- pre-push: runs tests before pushing changes" 