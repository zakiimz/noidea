package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/github"
	"github.com/AccursedGalaxy/noidea/internal/secure"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "GitHub integration commands",
	Long:  `Commands for interacting with GitHub repositories and services.`,
}

// githubAuthCmd represents the github auth command
var githubAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with GitHub",
	Long: `Authenticate with GitHub using a Personal Access Token (PAT).
This command will securely store your GitHub token for future use.

To create a new token, visit: https://github.com/settings/tokens
Required scopes: repo, read:user`,
	Run: func(cmd *cobra.Command, args []string) {
		runGitHubAuth()
	},
}

// githubStatusCmd represents the github status command
var githubStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check GitHub authentication status",
	Long:  `Check if you're authenticated with GitHub and display account information.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGitHubStatus()
	},
}

// githubLogoutCmd represents the github logout command
var githubLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored GitHub credentials",
	Long:  `Remove any stored GitHub Personal Access Tokens from your system.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGitHubLogout()
	},
}

// githubReleaseCmd represents the github release command
var githubReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "GitHub release commands",
	Long:  `Commands for managing GitHub releases.`,
}

// githubReleaseCreateCmd represents the github release create command
var githubReleaseCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a GitHub release",
	Long:  `Create a GitHub release for a specific tag.`,
	Run: func(cmd *cobra.Command, args []string) {
		tag, _ := cmd.Flags().GetString("tag")
		name, _ := cmd.Flags().GetString("name")
		draft, _ := cmd.Flags().GetBool("draft")
		prerelease, _ := cmd.Flags().GetBool("prerelease")
		runGitHubCreateRelease(tag, name, draft, prerelease)
	},
}

// githubHookInstallCmd represents the github hook install command
var githubHookInstallCmd = &cobra.Command{
	Use:   "hook-install",
	Short: "Install GitHub-related Git hooks",
	Long:  `Install Git hooks for GitHub integration, such as automatic release creation.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGitHubHookInstall()
	},
}

// githubReleaseNotesCmd represents the release notes command
var githubReleaseNotesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Generate and update release notes",
	Long: `Generate AI-enhanced release notes from commit messages and update GitHub release.
This command uses LLM (if enabled) to create comprehensive, user-friendly release notes.`,
	Run: func(cmd *cobra.Command, args []string) {
		tag, _ := cmd.Flags().GetString("tag")
		useAI, _ := cmd.Flags().GetBool("ai")
		skipApproval, _ := cmd.Flags().GetBool("skip-approval")
		runGitHubReleaseNotes(tag, useAI, skipApproval)
	},
}

func init() {
	rootCmd.AddCommand(githubCmd)
	githubCmd.AddCommand(githubAuthCmd)
	githubCmd.AddCommand(githubStatusCmd)
	githubCmd.AddCommand(githubLogoutCmd)
	githubCmd.AddCommand(githubReleaseCmd)
	githubCmd.AddCommand(githubHookInstallCmd)

	// Release command
	githubReleaseCmd.AddCommand(githubReleaseCreateCmd)

	// Release notes command
	githubReleaseCmd.AddCommand(githubReleaseNotesCmd)

	// Flags for release create command
	githubReleaseCreateCmd.Flags().String("tag", "", "Tag name for the release (required)")
	githubReleaseCreateCmd.Flags().String("name", "", "Release name (defaults to tag name)")
	githubReleaseCreateCmd.Flags().Bool("draft", false, "Mark as a draft release")
	githubReleaseCreateCmd.Flags().Bool("prerelease", false, "Mark as a prerelease")
	githubReleaseCreateCmd.MarkFlagRequired("tag")

	// Flags for release notes command
	githubReleaseNotesCmd.Flags().String("tag", "", "Tag name to generate notes for (defaults to latest tag)")
	githubReleaseNotesCmd.Flags().Bool("ai", false, "Force AI-generated notes even if LLM is disabled in config")
	githubReleaseNotesCmd.Flags().Bool("skip-approval", false, "Skip approval before updating release notes")
}

// runGitHubAuth handles the GitHub authentication flow
func runGitHubAuth() {
	fmt.Println("GitHub Authentication")
	fmt.Println("---------------------")
	fmt.Println("This will store a GitHub Personal Access Token (PAT) for noidea to use.")
	fmt.Println("To create a new token, visit: https://github.com/settings/tokens")
	fmt.Println("Required scopes: repo, read:user")
	fmt.Println()

	// Ask if the user wants to proceed
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Would you like to proceed? (y/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Authentication cancelled.")
		return
	}

	// Prompt for token
	fmt.Print("Enter your GitHub Personal Access Token (input will be hidden): ")
	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after hidden input
	if err != nil {
		fmt.Printf("Error reading token: %s\n", err)
		return
	}

	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		fmt.Println("Token cannot be empty. Authentication cancelled.")
		return
	}

	// Validate the token
	fmt.Println("Validating token...")
	valid, userData, err := secure.ValidateGitHubToken(token)
	if err != nil || !valid {
		if err != nil {
			fmt.Printf("Error validating token: %s\n", err)
		} else {
			fmt.Println("Invalid token. Please check your token and try again.")
		}
		return
	}

	// Store the token
	err = secure.StoreGitHubToken(token)
	if err != nil {
		fmt.Printf("Error storing token: %s\n", err)
		return
	}

	// Show success message with user info
	username := "Unknown"
	if userData != nil {
		if login, ok := userData["login"].(string); ok {
			username = login
		}
	}

	fmt.Printf("Successfully authenticated as: %s\n", username)
	fmt.Println("Your GitHub token has been securely stored.")
}

// runGitHubStatus checks and displays GitHub authentication status
func runGitHubStatus() {
	token, err := secure.GetGitHubToken()
	if err != nil {
		fmt.Println("Not authenticated with GitHub.")
		fmt.Println("Run 'noidea github auth' to authenticate.")
		return
	}

	// Token exists, validate it
	fmt.Println("Checking GitHub authentication status...")
	valid, userData, err := secure.ValidateGitHubToken(token)
	if err != nil || !valid {
		fmt.Println("Your GitHub token is invalid or expired.")
		fmt.Println("Run 'noidea github auth' to re-authenticate.")
		return
	}

	// Display user information
	fmt.Println("GitHub Authentication: ✅ Active")
	if userData != nil {
		if login, ok := userData["login"].(string); ok {
			fmt.Printf("Username: %s\n", login)
		}
		if name, ok := userData["name"].(string); ok && name != "" {
			fmt.Printf("Name: %s\n", name)
		}
	}
}

// runGitHubLogout removes stored GitHub credentials
func runGitHubLogout() {
	// Check if we have a token first
	_, err := secure.GetGitHubToken()
	if err != nil {
		fmt.Println("No GitHub credentials found.")
		return
	}

	// Confirm with the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure you want to remove your GitHub credentials? (y/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Operation cancelled.")
		return
	}

	// Delete the token
	err = secure.DeleteGitHubToken()
	if err != nil {
		fmt.Printf("Error removing credentials: %s\n", err)
		return
	}

	fmt.Println("GitHub credentials successfully removed.")
}

// runGitHubCreateRelease handles creating a GitHub release
func runGitHubCreateRelease(tag, name string, draft, prerelease bool) {
	// Initialize GitHub client
	client, err := github.NewClient()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	// If name is empty, use tag name
	if name == "" {
		name = tag
	}

	// Get repository owner and name
	owner, repo, err := github.ExtractRepoInfo("")
	if err != nil {
		fmt.Printf("Error: Failed to determine repository info: %s\n", err)
		fmt.Println("Make sure you're in a GitHub repository with a valid remote.")
		return
	}

	// Generate release notes from tag (get commit messages since last tag)
	body, err := generateReleaseNotes(tag)
	if err != nil {
		fmt.Printf("Warning: Failed to generate release notes: %s\n", err)
		body = "Release " + tag
	}

	fmt.Printf("Creating GitHub release for tag '%s' in %s/%s\n", tag, owner, repo)

	// Create the release
	release, err := client.CreateRelease(owner, repo, tag, name, body, draft, prerelease)
	if err != nil {
		fmt.Printf("Error creating release: %s\n", err)
		return
	}

	// Display success message
	fmt.Println("✅ Release created successfully!")
	if url, ok := release["html_url"].(string); ok {
		fmt.Printf("URL: %s\n", url)
	}
}

// generateReleaseNotes creates release notes from Git commit messages
func generateReleaseNotes(tag string) (string, error) {
	// Get the previous tag
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0", tag+"^")
	prevTag, err := cmd.Output()
	if err != nil {
		// If there's no previous tag, get all commits up to this tag
		cmd = exec.Command("git", "log", "--pretty=format:- %s", tag)
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("## Release %s\n\n%s", tag, string(output)), nil
	}

	// Get commit messages between previous tag and this tag
	cmd = exec.Command("git", "log", "--pretty=format:- %s", strings.TrimSpace(string(prevTag))+".."+tag)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("## Release %s\n\n%s", tag, string(output)), nil
}

// runGitHubHookInstall installs GitHub hooks
func runGitHubHookInstall() {
	// Check GitHub authentication
	_, err := secure.GetGitHubToken()
	if err != nil {
		fmt.Println("GitHub authentication required to install hooks.")
		fmt.Println("Run 'noidea github auth' to authenticate first.")
		return
	}

	// Install post-tag hook
	err = github.InstallPostTagHook()
	if err != nil {
		fmt.Printf("Error installing post-tag hook: %s\n", err)
		return
	}

	fmt.Println("GitHub hooks installed successfully!")
	fmt.Println("Now when you create a Git tag, a GitHub release will be created automatically.")
}

// runGitHubReleaseNotes handles generating and updating release notes
func runGitHubReleaseNotes(tag string, forceAI bool, skipApproval bool) {
	// Check if we're authenticated with GitHub
	_, err := secure.GetGitHubToken()
	if err != nil {
		fmt.Println("GitHub authentication required.")
		fmt.Println("Run 'noidea github auth' to authenticate.")
		return
	}

	// If no tag specified, try to get the latest tag
	if tag == "" {
		var err error
		tag, err = getLatestTag()
		if err != nil {
			fmt.Printf("Error: Please specify a tag with --tag or ensure you're in a Git repository with tags: %s\n", err)
			return
		}
		fmt.Printf("No tag specified, using latest tag: %s\n", tag)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Override LLM enabled setting if --ai flag is provided
	if forceAI {
		cfg.LLM.Enabled = true
	}

	// Create release manager
	manager, err := github.NewReleaseManager(cfg)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Generating %s release notes for tag %s...\n",
		getGenerationTypeString(cfg.LLM.Enabled), tag)

	// Update the release notes
	err = manager.UpdateReleaseNotes(tag, skipApproval)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}

// getLatestTag returns the latest tag in the Git repository
func getLatestTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get latest tag: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// getGenerationTypeString returns a string describing the type of generation
func getGenerationTypeString(llmEnabled bool) string {
	if llmEnabled {
		return "AI-enhanced"
	}
	return "standard"
}
