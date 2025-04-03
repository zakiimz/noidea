#!/bin/bash
#
# Installs Git hooks for the noidea project
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOOKS_DIR="$(git rev-parse --git-dir)/hooks"
GIT_ROOT="$(git rev-parse --show-toplevel)"
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

install_hook() {
    local hook_name="$1"
    local source_path="$SCRIPT_DIR/$hook_name"
    local target_path="$HOOKS_DIR/$hook_name"
    
    # Check if source exists and is readable
    if [ ! -r "$source_path" ]; then
        echo "Error: Hook $hook_name not found at $source_path"
        return 1
    fi
    
    # Copy the hook
    cp "$source_path" "$target_path"
    chmod +x "$target_path"
    echo -e "${GREEN}✓${NC} Installed $hook_name hook"
}

# Install the prepare-commit-msg hook
install_hook "prepare-commit-msg"

# Ask if the user wants to enable the commit message suggestion feature
read -p "Do you want to enable commit message suggestions? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git config noidea.suggest true
    echo -e "${GREEN}✓${NC} Enabled commit message suggestions"
    
    # Ask about interactive mode
    echo -e "${BLUE}Note:${NC} Interactive mode only applies when running 'noidea suggest' directly."
    echo "      Git hooks always use non-interactive mode to avoid input issues."
    echo "      You can still edit the message in your editor after suggestion."
    read -p "Do you want to enable interactive mode for direct command usage? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git config noidea.suggest.interactive true
        echo -e "${GREEN}✓${NC} Enabled interactive mode for direct command usage"
    else
        git config noidea.suggest.interactive false
        echo -e "${GREEN}✓${NC} Disabled interactive mode"
    fi
    
    # Ask about full diff mode
    read -p "Do you want to include full diffs in analysis? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git config noidea.suggest.full-diff true
        echo -e "${GREEN}✓${NC} Enabled full diff analysis"
    else
        git config noidea.suggest.full-diff false
        echo -e "${GREEN}✓${NC} Disabled full diff analysis"
    fi
    
    # Check if noidea binary is in PATH or in the repository
    if command -v noidea >/dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} Found noidea in PATH"
    elif [ -x "$GIT_ROOT/noidea" ]; then
        echo -e "${GREEN}✓${NC} Found noidea in repository root"
    else
        echo -e "${YELLOW}!${NC} The noidea binary was not found in PATH or repository root"
        echo "   For the hook to work properly, either:"
        echo "   1. Add noidea to your PATH"
        echo "   2. Build noidea in the repository root (./noidea)"
        echo "   3. Place the noidea binary in a common location like ./bin/ or ./dist/"
    fi
else
    git config noidea.suggest false
    echo -e "${GREEN}✓${NC} Disabled commit message suggestions"
fi

echo -e "\n${GREEN}✓${NC} Git hooks installation complete"
echo "To uninstall, run: git config noidea.suggest false"
echo "To change settings, use: git config noidea.suggest.interactive [true|false]"
echo "                          git config noidea.suggest.full-diff [true|false]"
echo -e "\n${GREEN}Note:${NC} Commit message suggestions always use a professional format"
echo "      regardless of any personality settings used elsewhere in noidea." 