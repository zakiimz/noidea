package git

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestRepo creates a temporary Git repository for testing
func setupTestRepo(t *testing.T) string {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "noidea-test-git")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to init git repository: %v", err)
	}

	// Set up git config for the test repo
	configCmd := exec.Command("git", "config", "user.name", "NoIdea Test")
	configCmd.Dir = tempDir
	if err := configCmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to set git config user.name: %v", err)
	}

	configCmd = exec.Command("git", "config", "user.email", "test@noidea.test")
	configCmd.Dir = tempDir
	if err := configCmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to set git config user.email: %v", err)
	}

	return tempDir
}

// cleanupTestRepo removes the temporary Git repository
func cleanupTestRepo(path string) {
	os.RemoveAll(path)
}

// TestFindGitDir tests the FindGitDir function
func TestFindGitDir(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git executable not available, skipping test")
	}

	// Create a test repository
	repoPath := setupTestRepo(t)
	defer cleanupTestRepo(repoPath)

	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(currentDir)

	// Change to test repository directory
	if err := os.Chdir(repoPath); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test FindGitDir function
	gitDir, err := FindGitDir()
	if err != nil {
		t.Fatalf("FindGitDir failed: %v", err)
	}

	// Verify that the git directory exists and is within the repo path
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf("Git directory does not exist: %s", gitDir)
	}

	// Test outside a git repository
	outsideDir, err := ioutil.TempDir("", "noidea-test-outside")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(outsideDir)

	if err := os.Chdir(outsideDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// FindGitDir should fail outside a git repository
	_, err = FindGitDir()
	if err == nil {
		t.Error("FindGitDir should fail outside a git repository")
	}
}

// TestInstallPostCommitHook tests the installation of the post-commit hook
func TestInstallPostCommitHook(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git executable not available, skipping test")
	}

	// Create a test repository
	repoPath := setupTestRepo(t)
	defer cleanupTestRepo(repoPath)

	// Create hooks directory
	hooksDir := filepath.Join(repoPath, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	// Install the hook
	err := InstallPostCommitHook(hooksDir)
	if err != nil {
		t.Fatalf("InstallPostCommitHook failed: %v", err)
	}

	// Verify that the hook file exists
	hookPath := filepath.Join(hooksDir, "post-commit")
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Errorf("Post-commit hook was not created at %s", hookPath)
	}

	// Verify that the hook file is executable
	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("Failed to get hook file info: %v", err)
	}
	if info.Mode()&0111 == 0 {
		t.Error("Post-commit hook is not executable")
	}

	// Verify that the hook file contains the expected content
	content, err := ioutil.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read hook file: %v", err)
	}
	if !strings.Contains(string(content), "noidea moai") {
		t.Error("Post-commit hook does not contain expected content")
	}
}

// TestInstallPrepareCommitMsgHook tests the installation of the prepare-commit-msg hook
func TestInstallPrepareCommitMsgHook(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git executable not available, skipping test")
	}

	// Create a test repository
	repoPath := setupTestRepo(t)
	defer cleanupTestRepo(repoPath)

	// Create hooks directory
	hooksDir := filepath.Join(repoPath, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	// Install the hook
	err := InstallPrepareCommitMsgHook(hooksDir)
	if err != nil {
		t.Fatalf("InstallPrepareCommitMsgHook failed: %v", err)
	}

	// Verify that the hook file exists
	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Errorf("Prepare-commit-msg hook was not created at %s", hookPath)
	}

	// Verify that the hook file is executable
	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("Failed to get hook file info: %v", err)
	}
	if info.Mode()&0111 == 0 {
		t.Error("Prepare-commit-msg hook is not executable")
	}

	// Verify that the hook file contains the expected content
	content, err := ioutil.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read hook file: %v", err)
	}
	if !strings.Contains(string(content), "noidea suggest") {
		t.Error("Prepare-commit-msg hook does not contain expected content")
	}
}
