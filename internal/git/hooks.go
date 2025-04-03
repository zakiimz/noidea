package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AccursedGalaxy/noidea/internal/config"
)

// FindGitDir finds the .git directory from the current path
func FindGitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	// Trim newline from the end of the output
	gitDir := string(output)
	if len(gitDir) > 0 && gitDir[len(gitDir)-1] == '\n' {
		gitDir = gitDir[:len(gitDir)-1]
	}
	
	// Convert to absolute path if it's relative
	if !filepath.IsAbs(gitDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		gitDir = filepath.Join(cwd, gitDir)
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

// InstallPostCommitHook installs the post-commit hook into the Git repository
func InstallPostCommitHook(hooksDir string) error {
	postCommitPath := filepath.Join(hooksDir, "post-commit")
	
	// Create hooks directory if it doesn't exist
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return err
	}
	
	// Get the absolute path to the noidea executable
	execPath, err := os.Executable()
	if err != nil {
		return err
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
		return err
	}
	
	fmt.Println("Installed post-commit hook at:", postCommitPath)
	return nil
} 