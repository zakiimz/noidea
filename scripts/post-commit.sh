#!/bin/sh
#
# noidea - Post-commit hook
# This hook calls the 'noidea moai' command after each commit
# to show a Moai face with feedback about your commit.

# Get the last commit message
COMMIT_MSG=$(git log -1 --pretty=%B)

# Call noidea with the commit message and history context
noidea moai --history "$COMMIT_MSG"

# Always exit with success so git continues normally
exit 0 