#!/bin/sh
#
# noidea - Post-commit hook
# This hook calls the 'noidea moai' command after each commit
# to show a Moai face with feedback about your commit.

# Define some terminal colors if supported
if [ -t 1 ]; then
    GREEN="\033[0;32m"
    YELLOW="\033[1;33m"
    CYAN="\033[0;36m"
    RED="\033[0;31m"
    RESET="\033[0m"
else
    GREEN=""
    YELLOW=""
    CYAN=""
    RED=""
    RESET=""
fi

# Print a divider
print_divider() {
    echo "${CYAN}───────────────────────────────────────────────────${RESET}"
}

# Find the noidea binary in various possible locations
find_noidea() {
    # Check if it's in PATH
    if command -v noidea >/dev/null 2>&1; then
        echo "noidea"
        return 0
    fi

    # Try to determine the git root
    GIT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
    if [ -n "$GIT_ROOT" ]; then
        # Check common binary locations
        for location in "$GIT_ROOT/noidea" "$GIT_ROOT/bin/noidea" "$GIT_ROOT/build/noidea" "$GIT_ROOT/dist/noidea"; do
            if [ -x "$location" ]; then
                echo "$location"
                return 0
            fi
        done
    fi
    
    # Not found
    return 1
}

# Find the noidea binary
NOIDEA_BIN=$(find_noidea)
if [ -z "$NOIDEA_BIN" ]; then
    echo "${YELLOW}⚠️  Warning: noidea binary not found${RESET}"
    echo "${YELLOW}   Please ensure noidea is in your PATH or at the repository root${RESET}"
    exit 0
fi

# Get the last commit message
COMMIT_MSG=$(git log -1 --pretty=%B)
if [ -z "$COMMIT_MSG" ]; then
    echo "${YELLOW}⚠️  Warning: Could not get the commit message${RESET}"
    COMMIT_MSG="unknown commit"
fi

# Print a divider before displaying the Moai
print_divider

# Call noidea with the commit message and history context
"$NOIDEA_BIN" moai --history "$COMMIT_MSG" || echo "${RED}Error running noidea moai command${RESET}"

# Print a final divider
print_divider

# Always exit with success so git continues normally
exit 0 