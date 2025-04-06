#!/bin/bash
#
# setup-github.sh
# Sets up GitHub integration for NoIdea
#

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if noidea executable is in path
if ! command -v noidea >/dev/null 2>&1; then
    echo -e "${RED}Error: 'noidea' command not found.${NC}"
    echo "Please make sure noidea is installed and in your PATH."
    exit 1
fi

# Welcome message
echo -e "${CYAN}NoIdea - GitHub Integration Setup${NC}"
echo ""
echo -e "This script will set up GitHub integration for NoIdea."
echo ""

# Ask to proceed
read -p "Do you want to proceed? (y/n): " PROCEED
if [[ ! $PROCEED =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Setup cancelled.${NC}"
    exit 0
fi

echo ""
echo -e "${CYAN}Step 1: GitHub Authentication${NC}"
echo -e "You need a GitHub Personal Access Token (PAT)."
echo -e "Create one at: ${YELLOW}https://github.com/settings/tokens${NC}"
echo -e "Required scopes: ${YELLOW}repo, read:user${NC}"
echo ""

read -p "Authenticate with GitHub now? (y/n): " AUTH
if [[ $AUTH =~ ^[Yy]$ ]]; then
    noidea github auth
    if [ $? -ne 0 ]; then
        echo -e "${RED}GitHub authentication failed.${NC}"
        echo "You can try again later with: noidea github auth"
        exit 1
    fi
else
    echo -e "${YELLOW}Skipping authentication.${NC}"
fi

echo ""
echo -e "${CYAN}Step 2: GitHub Hook Installation${NC}"
echo -e "NoIdea can automatically create GitHub releases when you create Git tags."
echo ""

read -p "Install GitHub hooks now? (y/n): " HOOKS
if [[ $HOOKS =~ ^[Yy]$ ]]; then
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo -e "${RED}Error: Not in a git repository.${NC}"
        echo "Please run this script from within a git repository."
        exit 1
    fi
    
    noidea github hook-install
    if [ $? -ne 0 ]; then
        echo -e "${RED}Hook installation failed.${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}Skipping hook installation.${NC}"
fi

echo ""
echo -e "${GREEN}GitHub integration setup complete!${NC}"
echo ""
echo -e "Available commands:"
echo -e "  ${CYAN}noidea github status${NC}"
echo -e "  ${CYAN}noidea github release create${NC}" 