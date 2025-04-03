#!/bin/sh
#
# noidea - Post-commit hook
# This hook calls the 'noidea moai' command after each commit
# to show a Moai face with feedback about your commit.

# Define some terminal colors if supported
if [ -t 1 ]; then
    GREEN="\033[0;32m"
    CYAN="\033[0;36m"
    RESET="\033[0m"
else
    GREEN=""
    CYAN=""
    RESET=""
fi

# Print a divider
print_divider() {
    echo "${CYAN}───────────────────────────────────────────────────${RESET}"
}

# Get the last commit message
COMMIT_MSG=$(git log -1 --pretty=%B)

# Print a divider before displaying the Moai
print_divider

# Call noidea with the commit message and history context
# Capture the output to check for errors
noidea moai --history "$COMMIT_MSG" 

# Print a final divider
print_divider

# Always exit with success so git continues normally
exit 0 