package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AccursedGalaxy/noidea/internal/git"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize noidea in your Git repository",
	Long:  `Install the post-commit hook in your Git repository to show Moai after each commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we're in a Git repository
		gitDir, err := git.FindGitDir()
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Not in a Git repository.")
			os.Exit(1)
		}

		// Install the post-commit hook
		hooksDir := filepath.Join(gitDir, "hooks")
		err = git.InstallPostCommitHook(hooksDir)
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to install post-commit hook:", err)
			os.Exit(1)
		}

		fmt.Println(color.GreenString("Success!"), "Moai is now watching your commits. 🗿")
	},
} 