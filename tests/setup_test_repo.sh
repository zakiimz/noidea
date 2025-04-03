#!/bin/bash
# setup_test_repo.sh
# This script creates a test git repository with an initial commit
# for use in testing the commit suggestion feature

# Get the absolute path to the noidea binary
NOIDEA_BIN=$(realpath "../noidea")
echo "Using noidea binary: $NOIDEA_BIN"

# Create test repo directory if it doesn't exist
mkdir -p test_repo

# Navigate to the test repo
cd test_repo || exit 1

# Initialize git repository if it doesn't exist
if [ ! -d .git ]; then
    echo "Initializing git repository..."
    git init
    
    # Configure git user (needed for commits)
    git config user.name "Test User"
    git config user.email "test@example.com"
else
    echo "Git repository already exists, cleaning..."
    git clean -fd
    # Try to reset, but don't fail if there's no commit yet
    git reset --hard HEAD 2>/dev/null || true
fi

# Create a simple README file
echo "# Test Repository

This is a test repository for noidea commit suggestion tests.

## Purpose

This repository contains various staged changes to test the commit 
suggestion functionality of the noidea tool.
" > README.md

# Make an initial commit if needed
if ! git rev-parse --verify HEAD >/dev/null 2>&1; then
    echo "Creating initial commit..."
    git add README.md
    git commit -m "Initial commit"
fi

# Initialize noidea in the test repository
echo "Initializing noidea in test repository..."
"$NOIDEA_BIN" init

# Enable suggest feature explicitly
echo "Enabling noidea suggest feature..."
git config noidea.suggest true
git config noidea.suggest.full-diff true

# Set API key environment variable if it exists in parent environment
API_KEY=""
if [ -n "$XAI_API_KEY" ]; then
    API_KEY="$XAI_API_KEY"
elif [ -n "$OPENAI_API_KEY" ]; then
    API_KEY="$OPENAI_API_KEY"
elif [ -n "$DEEPSEEK_API_KEY" ]; then
    API_KEY="$DEEPSEEK_API_KEY"
fi

if [ -n "$API_KEY" ]; then
    echo "Setting API key for testing..."
    # Truncate for security in logs
    MASKED_KEY="${API_KEY:0:4}...${API_KEY: -4}"
    echo "Using API key: $MASKED_KEY"
    
    # Pass API key via git config
    if [[ "$API_KEY" == *"xai"* ]]; then
        git config noidea.xai.api_key "$API_KEY"
    elif [[ "$API_KEY" == *"openai"* ]]; then
        git config noidea.openai.api_key "$API_KEY"
    else
        git config noidea.deepseek.api_key "$API_KEY"
    fi
else
    echo "No API key found, using mock key for testing"
    git config noidea.llm.enabled true
    git config noidea.xai.api_key "mock-api-key-for-testing"
fi

echo "Test repository setup complete."

# Go back to original directory
cd - > /dev/null 