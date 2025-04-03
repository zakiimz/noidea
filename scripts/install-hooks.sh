#!/bin/bash
#
# Installs Git hooks for the noidea project
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOOKS_DIR="$(git rev-parse --git-dir)/hooks"
GREEN='\033[0;32m'
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
    read -p "Do you want to enable interactive mode for suggestions? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git config noidea.suggest.interactive true
        echo -e "${GREEN}✓${NC} Enabled interactive mode"
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