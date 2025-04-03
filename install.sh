#!/bin/bash
#
# noidea - Simple installation script
#
# This script:
# 1. Builds the noidea binary
# 2. Installs it to /usr/local/bin (or custom location)
# 3. Creates the ~/.noidea directory for configuration
#

set -e

# Text formatting
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Preserve the real user and their home directory
REAL_USER=${SUDO_USER:-$USER}
REAL_HOME=$(eval echo ~$REAL_USER)

# Default installation directory
INSTALL_DIR="/usr/local/bin"
if [ -n "$1" ]; then
    INSTALL_DIR="$1"
fi

echo -e "${BLUE}🧠 noidea installer${NC}"
echo "============================="

# Find Go - try different approaches
find_go() {
    # First check if go is directly available
    if command -v go &> /dev/null; then
        echo "go"
        return 0
    fi

    # Check common Go installation paths
    for path in "/usr/local/go/bin/go" "/usr/lib/go/bin/go" "/usr/bin/go" "$HOME/go/bin/go" "$HOME/.local/bin/go"; do
        if [ -x "$path" ]; then
            echo "$path"
            return 0
        fi
    done

    # If running as sudo, try the real user's PATH
    if [ -n "$SUDO_USER" ]; then
        # Get the real user's PATH
        USER_PATH=$(sudo -u "$SUDO_USER" bash -c 'echo $PATH')
        
        # Check each directory in the user's PATH
        IFS=':' read -r -a path_dirs <<< "$USER_PATH"
        for dir in "${path_dirs[@]}"; do
            if [ -x "$dir/go" ]; then
                echo "$dir/go"
                return 0
            fi
        done
    fi

    return 1
}

# Try to find Go
GO_CMD=$(find_go)

if [ -z "$GO_CMD" ]; then
    echo -e "${RED}Error: Go is not installed or not found in common locations${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo "Using Go at: $GO_CMD"

# Check if directory exists and is writable
if [ ! -d "$INSTALL_DIR" ]; then
    echo -e "${YELLOW}Directory $INSTALL_DIR does not exist. Creating it...${NC}"
    if ! mkdir -p "$INSTALL_DIR" 2>/dev/null; then
        echo -e "${RED}Error: Failed to create $INSTALL_DIR${NC}"
        if [ "$(id -u)" -ne 0 ]; then
            echo "Try running with sudo: sudo ./install.sh"
        fi
        exit 1
    fi
fi

if [ ! -w "$INSTALL_DIR" ]; then
    echo -e "${RED}Error: Cannot write to $INSTALL_DIR${NC}"
    if [ "$(id -u)" -ne 0 ]; then
        echo "Try running with sudo: sudo ./install.sh"
    fi
    exit 1
fi

# Create config directory for the real user
CONFIG_DIR="$REAL_HOME/.noidea"
echo "Creating config directory at $CONFIG_DIR"
mkdir -p "$CONFIG_DIR"
# If running as root, ensure the directory is owned by the real user
if [ "$(id -u)" -eq 0 ] && [ -n "$SUDO_USER" ]; then
    chown -R "$SUDO_USER" "$CONFIG_DIR"
fi

# Build the binary
echo "Building noidea..."
"$GO_CMD" build -o noidea

# Install the binary
echo "Installing noidea to $INSTALL_DIR"
cp noidea "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/noidea"

# Output success message
echo -e "${GREEN}✅ noidea has been successfully installed!${NC}"
echo ""
echo "To use noidea in a Git repository:"
echo "  1. cd /path/to/your/repo"
echo "  2. noidea init                 # Set up Git hooks"
echo ""
echo "For commit message suggestions:"
echo "  noidea suggest                 # Run directly"
echo "  git config noidea.suggest true # Enable automatic suggestions"
echo ""
echo "For help and available commands:"
echo "  noidea --help"
echo ""

# Check if the installation directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}⚠️  Warning: $INSTALL_DIR is not in your PATH${NC}"
    echo "You might need to add it to your shell configuration:"
    echo ""
    echo "  For bash/zsh:"
    echo "    echo 'export PATH=\$PATH:$INSTALL_DIR' >> ~/.bashrc"
    echo "    # or"
    echo "    echo 'export PATH=\$PATH:$INSTALL_DIR' >> ~/.zshrc"
    echo ""
    echo "  For fish:"
    echo "    echo 'set -gx PATH \$PATH $INSTALL_DIR' >> ~/.config/fish/config.fish"
    echo ""
fi

echo -e "${BLUE}Thank you for installing noidea!${NC}" 