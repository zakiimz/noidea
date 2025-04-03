package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AccursedGalaxy/noidea/internal/config"
)

// FindGitDir returns the path to the .git directory for the current repository.
// If not in a git repository, returns an error.
func FindGitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}
	
	gitDir := strings.TrimSpace(string(output))
	if gitDir == "" {
		return "", fmt.Errorf("unable to determine git directory")
	}
	
	// If the git dir is relative (usually .git), make it absolute
	if !filepath.IsAbs(gitDir) {
		workDir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
		gitDir = filepath.Join(workDir, gitDir)
	}
	
	return gitDir, nil
}

// GetScriptPath returns the absolute path to the scripts directory
func GetScriptPath() (string, error) {
	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Get the directory of the executable
	execDir := filepath.Dir(execPath)

	// The scripts directory should be in the same directory as the executable
	scriptsDir := filepath.Join(execDir, "..", "scripts")

	return scriptsDir, nil
}

// InstallPostCommitHook installs the post-commit hook script in the specified
// hooks directory. The hook will call 'noidea moai' after each commit to show
// feedback about the commit message.
func InstallPostCommitHook(hooksDir string) error {
	postCommitPath := filepath.Join(hooksDir, "post-commit")
	
	// Create hooks directory if it doesn't exist
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}
	
	// Get the absolute path to the noidea executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	
	// Load configuration
	cfg := config.LoadConfig()
	
	// Build command flags
	flags := ""
	
	// Add AI flag if enabled
	if cfg.LLM.Enabled {
		flags += "--ai "
	}
	
	// Add personality flag if set
	if cfg.Moai.Personality != "" {
		flags += fmt.Sprintf("--personality=%s ", cfg.Moai.Personality)
	}
	
	// Create the post-commit hook content
	hookContent := fmt.Sprintf(`#!/bin/sh
#
# noidea - Post-commit hook
# This hook calls the 'noidea moai' command after each commit
# to show a Moai face with feedback about your commit.

# Get the last commit message
COMMIT_MSG=$(git log -1 --pretty=%%B)

# Call noidea with the commit message (using absolute path)
%s moai %s"$COMMIT_MSG"

# Always exit with success so git continues normally
exit 0
`, execPath, flags)
	
	// Write the hook file
	if err := os.WriteFile(postCommitPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("failed to write post-commit hook: %w", err)
	}
	
	fmt.Println("Installed post-commit hook at:", postCommitPath)
	return nil
}

// InstallPrepareCommitMsgHook installs the prepare-commit-msg hook for commit message suggestions.
// This hook runs before Git creates a commit and offers AI-generated commit message suggestions
// based on the staged changes.
func InstallPrepareCommitMsgHook(hooksDir string) error {
	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")
	
	// Create hooks directory if it doesn't exist
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}
	
	// Get the absolute path to the noidea executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	
	// Create the hook content
	hookContent := fmt.Sprintf(`#!/bin/sh
#
# noidea - prepare-commit-msg hook
# This hook calls 'noidea suggest' to generate commit message suggestions
# To disable, run: git config noidea.suggest false

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

# Get commit message file
COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2

# Check if noidea's suggestion feature is enabled
if [ "$(git config --get noidea.suggest)" != "true" ]; then
    exit 0
fi

# Skip if it's a merge, rebase, or cherry-pick
if [ "$COMMIT_SOURCE" = "merge" ] || [ "$COMMIT_SOURCE" = "squash" ] || [ -n "$COMMIT_SOURCE" ]; then
    exit 0
fi

# Check if the commit message already has content
if [ -s "$COMMIT_MSG_FILE" ]; then
    # Has content already - user may have specified a message with -m
    # Skip if the file already has content beyond comments
    if grep -v "^#" "$COMMIT_MSG_FILE" | grep -q "[^[:space:]]"; then
        exit 0
    fi
fi

# Always use non-interactive mode for hooks to prevent stdin issues
INTERACTIVE_FLAG=""

# Get history setting from config
HISTORY_FLAG="--history 10"

# Get full diff setting from config
FULL_DIFF=$(git config --get noidea.suggest.full-diff)
if [ "$FULL_DIFF" = "true" ]; then
    DIFF_FLAG="--full-diff"
else
    DIFF_FLAG=""
fi

# Generate a suggested commit message
echo "${CYAN}🧠 Generating commit message suggestion...${RESET}"
%s suggest $INTERACTIVE_FLAG $HISTORY_FLAG $DIFF_FLAG --file "$COMMIT_MSG_FILE"

exit 0
`, execPath)
	
	// Write the hook file
	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("failed to write prepare-commit-msg hook: %w", err)
	}
	
	fmt.Println("Installed prepare-commit-msg hook at:", hookPath)
	return nil
}
