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

# Only stash if there are unstaged changes
if git diff --quiet; then
  NEED_STASH=0
else
  NEED_STASH=1
  echo "Stashing unstaged changes..."
  git stash push -q --keep-index --include-untracked --message "pre-commit-hook"
fi

# Ensure we clean up on exit
function cleanup {
  # Only pop the stash if we needed to stash
  if [ $NEED_STASH -eq 1 ]; then
    echo "Restoring unstaged changes..."
    git stash pop -q
  fi
}

# Use ERR trap instead of EXIT to keep stashed changes if there's an error
trap cleanup EXIT

# Run go mod tidy to manage dependencies
echo "Checking dependencies..."
go mod tidy -quiet

# Run goimports to format and organize imports
if command -v goimports &> /dev/null; then
  echo "Formatting code..."
  find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w -local github.com/AccursedGalaxy/noidea >/dev/null 2>&1
else
  echo "Installing goimports..."
  go install golang.org/x/tools/cmd/goimports@latest >/dev/null 2>&1
  find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w -local github.com/AccursedGalaxy/noidea >/dev/null 2>&1
fi

# Format the code
make format -s

# Run lint
make script-lint -s

# If there are any changes after formatting, add them
if ! git diff --quiet; then
  echo "Adding formatted files..."
  git add -u
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
make test -s

echo "✅ Pre-push checks passed."
EOF

# Make the hooks executable
chmod +x "$HOOKS_DIR/pre-commit"
chmod +x "$HOOKS_DIR/pre-push"

echo "✅ Git hooks setup complete."
echo "Installed hooks: pre-commit, pre-push" 