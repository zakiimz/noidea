package github

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/git"
)

// InstallPostTagHook installs a Git hook that runs after tags are created
// to create GitHub releases automatically with enhanced release notes
func InstallPostTagHook() error {
	// Find the git directory
	gitDir, err := git.FindGitDir()
	if err != nil {
		return fmt.Errorf("failed to find git directory: %w", err)
	}

	// Path to hooks directory
	hooksDir := filepath.Join(gitDir, "hooks")

	// Ensure the hooks directory exists
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Path to post-tag hook
	hookPath := filepath.Join(hooksDir, "post-tag")

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

	// Create the hook content with enhanced release notes support
	hookContent := fmt.Sprintf(`#!/bin/sh
#
# noidea - Post-tag hook
# This hook provides information about GitHub release generation after a tag is created

# Get the tag name
TAG_NAME=$(git describe --tags --exact-match 2>/dev/null)
if [ -z "$TAG_NAME" ]; then
    echo "No tag found, skipping GitHub release info"
    exit 0
fi

# Information about GitHub release and recommended workflow
echo "============================================================="
echo "🏷️  Tag '$TAG_NAME' created successfully!"
echo "GitHub Actions workflows may have been triggered."
echo ""
echo "📝 For best results, wait for GitHub workflows to complete"
echo "   before generating enhanced release notes with:"
echo ""
echo "   %s github release notes --tag=\"$TAG_NAME\" %s"
echo "============================================================="

# Exit with success
exit 0
`, execPath, flags)

	// Write the hook file
	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("failed to write post-tag hook: %w", err)
	}

	fmt.Println("Installed post-tag hook at:", hookPath)
	return nil
}
