package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Version information
var (
	Version   = "v0.1.1" // Will be overridden during build
	BuildDate = "dev"   // Will be overridden during build
	Commit    = "none"  // Will be overridden during build
)

// Flag variables
var versionFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "noidea",
	Short: "noidea - The Git Extension You Never Knew You Needed",
	Long: `🧠 noidea - A lightweight, plug-and-play Git extension that adds 
✨fun and occasionally useful ✨feedback into your normal Git workflow.

Every time you commit, a mysterious Moai appears to judge your code.

Main commands:
  suggest     Generate commit message suggestions based on staged changes
  moai        Show feedback about your most recent commit
  summary     Generate a summary of your recent Git activity
  feedback    Get detailed feedback on your recent commits
  init        Set up noidea in your Git repository
  config      Manage noidea configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		// If version flag is set, print version and exit
		if versionFlag {
			printVersion()
			return
		}
		
		// If no subcommand is provided, print help
		cmd.Help()
	},
}

func init() {
	// Add version flag
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print version information and exit")
}

// Execute executes the root command and handles any errors
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// printVersion prints detailed version information
func printVersion() {
	fmt.Printf("noidea version %s\n", Version)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Git commit: %s\n", Commit)
}
