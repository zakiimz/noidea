package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/git"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	enableSuggestions bool
	enableInteractive bool
	enableFullDiff    bool
)

func init() {
	initCmd.Flags().BoolVarP(&enableSuggestions, "suggest", "s", true, "Enable commit message suggestions")
	initCmd.Flags().BoolVarP(&enableInteractive, "interactive", "i", false, "Enable interactive mode for direct command usage")
	initCmd.Flags().BoolVarP(&enableFullDiff, "full-diff", "f", false, "Include full diffs in commit message analysis")

	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize noidea in your Git repository",
	Long:  `Install the Git hooks for noidea in your repository, including Moai feedback and commit message suggestions.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we're in a Git repository
		gitDir, err := git.FindGitDir()
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Not in a Git repository.")
			os.Exit(1)
		}

		// Create hooks directory if it doesn't exist
		hooksDir := filepath.Join(gitDir, "hooks")
		if err := os.MkdirAll(hooksDir, 0755); err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to create hooks directory:", err)
			os.Exit(1)
		}

		// Install the post-commit hook for Moai feedback
		err = git.InstallPostCommitHook(hooksDir)
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to install post-commit hook:", err)
			os.Exit(1)
		}
		fmt.Println(color.GreenString("✓"), "Installed post-commit hook for Moai feedback")

		// Install the prepare-commit-msg hook for commit suggestions
		err = git.InstallPrepareCommitMsgHook(hooksDir)
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to install prepare-commit-msg hook:", err)
			os.Exit(1)
		}
		fmt.Println(color.GreenString("✓"), "Installed prepare-commit-msg hook for commit suggestions")

		// Configure git settings based on flags
		gitConfigRunner := func(key, value string) {
			cmd := exec.Command("git", "config", key, value)
			if err := cmd.Run(); err != nil {
				fmt.Println(color.YellowString("Warning:"), "Failed to set git config", key, ":", err)
			}
		}

		// Set suggestion configuration
		gitConfigRunner("noidea.suggest", fmt.Sprintf("%t", enableSuggestions))
		status := "enabled"
		if !enableSuggestions {
			status = "disabled"
		}
		fmt.Println(color.GreenString("✓"), "Commit message suggestions", status)

		// Only configure these if suggestions are enabled
		if enableSuggestions {
			gitConfigRunner("noidea.suggest.interactive", fmt.Sprintf("%t", enableInteractive))
			if enableInteractive {
				fmt.Println(color.GreenString("✓"), "Interactive mode enabled for direct command usage")
				fmt.Println(color.BlueString("Note:"), "Interactive mode only applies when running 'noidea suggest' directly.")
				fmt.Println("      Git hooks always use non-interactive mode to avoid input issues.")
			}

			gitConfigRunner("noidea.suggest.full-diff", fmt.Sprintf("%t", enableFullDiff))
			if enableFullDiff {
				fmt.Println(color.GreenString("✓"), "Full diff analysis enabled")
			}
		}

		// Check if noidea is properly available
		execPath, _ := os.Executable()
		fmt.Println(color.GreenString("Success!"), "noidea hooks installed - executable at:", execPath)
		fmt.Println(color.BlueString("Note:"), "To change settings, run 'git config noidea.suggest [true|false]'")
		
		// Load config and check if API key is set
		cfg := config.LoadConfig()
		if cfg.LLM.Enabled && cfg.LLM.APIKey == "" {
			fmt.Println()
			fmt.Println(color.YellowString("⚠️  Warning:"), "LLM is enabled but no API key is configured.")
			fmt.Println("     For better commit message suggestions, configure your API key:")
			fmt.Println("     Run 'noidea config --init' or edit ~/.noidea/config.json")
			fmt.Println()
			fmt.Println("     Without an API key, commit suggestions will use a simple local algorithm")
			fmt.Println("     that's less detailed than the AI-powered suggestions.")
		} else if !cfg.LLM.Enabled && enableSuggestions {
			fmt.Println()
			fmt.Println(color.BlueString("💡 Tip:"), "For better commit message suggestions, enable AI integration:")
			fmt.Println("     Run 'noidea config --init' to configure AI settings.")
		}
	},
}
