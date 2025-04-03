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
echo -e "${BLUE}You can now delete the temporary files with:${NC}"
echo "rm -rf $TEMP_DIR" 