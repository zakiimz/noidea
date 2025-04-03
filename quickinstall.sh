#!/bin/bash
#
# noidea - One-line quick installer
#
# This script:
# 1. Clones the noidea repository into a temporary directory
# 2. Runs the installer
# 3. Cleans up the temporary files
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/AccursedGalaxy/noidea/main/quickinstall.sh | bash
#   # Or with sudo:
#   curl -sSL https://raw.githubusercontent.com/AccursedGalaxy/noidea/main/quickinstall.sh | sudo bash
#

set -e

# Text formatting
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧠 noidea quick installer${NC}"
echo "============================="

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
echo "Creating temporary directory: $TEMP_DIR"

# Clean up on exit
cleanup() {
    echo "Cleaning up temporary files..."
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Clone the repository
echo "Cloning noidea repository..."
git clone https://github.com/AccursedGalaxy/noidea.git "$TEMP_DIR" --depth=1

# Navigate to the directory
cd "$TEMP_DIR"

# Run the installer
echo "Running installer..."
./install.sh

echo -e "${GREEN}✅ noidea has been successfully installed!${NC}"

# Detect if we're in a Git repository
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo -e "${YELLOW}Detected that you're inside a Git repository.${NC}"
    echo -e "To initialize noidea in this repository, run:"
    echo -e "${GREEN}noidea init${NC}"
    echo ""
    echo -e "Would you like to run 'noidea init' now? (y/n)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        noidea init
    fi
else
    echo -e "${BLUE}To use noidea in a Git repository:${NC}"
    echo "  1. cd /path/to/your/repo"
    echo "  2. noidea init                 # Set up Git hooks"
fi 