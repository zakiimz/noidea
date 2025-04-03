package cmd

// suggest.go - Commit message suggestion functionality
//
// This file implements the 'suggest' command for generating conventional commit messages.
// Conventional commit format follows the pattern: <type>[(scope)]: <description>
//
// Common types include:
// - feat: A new feature
// - fix: A bug fix
// - docs: Documentation changes
// - style: Code style changes (formatting, etc)
// - refactor: Code changes that neither fix bugs nor add features
// - test: Adding or fixing tests
// - chore: Maintenance tasks, dependencies, etc
//
// Example: feat(auth): implement password reset functionality

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/feedback"
	"github.com/AccursedGalaxy/noidea/internal/history"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Suggest command flags
	historyCountFlag int
	fullDiffFlag     bool
	interactiveFlag  bool
	commitMsgFileFlag string
)

func init() {
	rootCmd.AddCommand(suggestCmd)

	// Add flags
	suggestCmd.Flags().IntVarP(&historyCountFlag, "history", "n", 10, "Number of recent commits to analyze for context")
	suggestCmd.Flags().BoolVarP(&fullDiffFlag, "full-diff", "f", false, "Include full diff instead of summary")
	suggestCmd.Flags().BoolVarP(&interactiveFlag, "interactive", "i", false, "Interactive mode to approve/reject suggestions")
	suggestCmd.Flags().StringVarP(&commitMsgFileFlag, "file", "F", "", "Path to commit message file (for prepare-commit-msg hook)")
}

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Suggest a commit message based on staged changes",
	Long:  `Analyze staged changes and commit history to suggest a descriptive commit message.
	
Commit message suggestions always follow professional conventional commit format,
regardless of the personality settings used elsewhere in noidea.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg := config.LoadConfig()

		// Get staged changes
		diff, err := getStagedDiff()
		if err != nil {
			fmt.Println(color.RedString("❌ Error:"), "Failed to get staged changes:", err)
			return
		}

		// Check if there are staged changes
		if strings.TrimSpace(diff) == "" {
			fmt.Println(color.YellowString("⚠️ No staged changes found. Stage files with 'git add' first."))
			return
		}

		// Get recent commit history for context
		commits, err := history.GetLastNCommits(historyCountFlag, false)
		if err != nil {
			fmt.Println(color.YellowString("⚠️ Warning:"), "Failed to get commit history. Continuing with staged changes only.")
		}

		// Extract commit messages and stats
		var commitMessages []string
		for _, commit := range commits {
			commitMessages = append(commitMessages, commit.Message)
		}

		// Create a history collector to calculate stats
		collector, _ := history.NewHistoryCollector()
		stats := collector.CalculateStats(commits)

		// Print a divider
		divider := "───────────────────────────────────────────────────"
		fmt.Println(color.HiBlackString(divider))
		
		// Print analysis info
		fmt.Printf("%s %s\n", 
			color.CyanString("🧠 Analyzing staged changes and"),
			color.CyanString(fmt.Sprintf("%d recent commits", len(commitMessages))))
		
		fmt.Printf("%s\n",
			color.CyanString("Generating professional commit message suggestion..."))

		// Create feedback engine based on config
		engineProvider := cfg.LLM.Provider
		engineModel := cfg.LLM.Model
		apiKey := cfg.LLM.APIKey
		personality := cfg.Moai.Personality
		personalityFile := cfg.Moai.PersonalityFile
		
		engine := feedback.NewFeedbackEngine(engineProvider, engineModel, apiKey, personality, personalityFile)
		
		// Create commit context for the suggestion
		ctx := feedback.CommitContext{
			Diff:          diff,
			CommitHistory: commitMessages,
			CommitStats:   stats,
			Timestamp:     time.Now(),
		}
		
		// Generate suggested commit message
		suggestion, err := engine.GenerateCommitSuggestion(ctx)
		if err != nil {
			fmt.Println(color.RedString("❌ Error:"), "Failed to generate suggestion:", err)
			return
		}

		// Print another divider
		fmt.Println(color.HiBlackString(divider))

		if interactiveFlag {
			// Handle interactive mode
			finalMessage := handleInteractiveMode(suggestion)
			
			// If we have a commit message file, write to it
			if commitMsgFileFlag != "" {
				err := writeToCommitMsgFile(finalMessage, commitMsgFileFlag)
				if err != nil {
					fmt.Println(color.RedString("❌ Error:"), "Failed to write commit message:", err)
					return
				}
				fmt.Println(color.GreenString("✅ Commit message suggestion applied"))
			} else {
				// Print another divider
				fmt.Println(color.HiBlackString(divider))
				fmt.Println(color.GreenString("✅ Final commit message:"))
				fmt.Println(color.HiWhiteString(finalMessage))
				fmt.Println(color.HiBlackString(divider))
			}
		} else {
			// Just print the suggestion
			fmt.Println(color.GreenString("✨ Suggested commit message:"))
			fmt.Println(color.HiWhiteString(suggestion))
			fmt.Println(color.HiBlackString(divider))
			
			// If we have a commit message file, write to it
			if commitMsgFileFlag != "" {
				err := writeToCommitMsgFile(suggestion, commitMsgFileFlag)
				if err != nil {
					fmt.Println(color.RedString("❌ Error:"), "Failed to write commit message:", err)
					return
				}
				fmt.Println(color.GreenString("✅ Commit message suggestion applied"))
			}
		}
	},
}

// getStagedDiff gets the diff of staged changes
func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}
	return string(output), nil
}

// handleInteractiveMode presents the suggestion to the user and allows interaction
func handleInteractiveMode(suggestion string) string {
	fmt.Println(color.GreenString("✨ Suggested commit message:"))
	fmt.Println(color.HiWhiteString(suggestion))
	
	// Check if we're in a terminal/interactive environment
	isTTY := isRunningInTerminal()
	if !isTTY {
		fmt.Println(color.YellowString("⚠️ Not running in an interactive terminal. Accepting suggestion automatically."))
		return suggestion
	}
	
	for {
		fmt.Println()
		fmt.Print(color.CyanString("Accept (a), Regenerate (r), Edit (e), or Cancel (c)? "))
		
		var choice string
		fmt.Scanln(&choice)
		
		switch strings.ToLower(choice) {
		case "a", "accept", "y", "yes":
			return suggestion
		case "r", "regenerate":
			fmt.Println(color.YellowString("🔄 Regenerating suggestion..."))
			// In a real implementation, we would regenerate here
			// For now, we'll just return the original
			return suggestion
		case "e", "edit":
			fmt.Println(color.CyanString("✏️ Enter your edited message (type 'done' on a new line when finished):"))
			
			var lines []string
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "done" {
					break
				}
				lines = append(lines, line)
			}
			
			return strings.Join(lines, "\n")
		case "c", "cancel":
			fmt.Println(color.YellowString("❌ Cancelled. Using default commit message."))
			return ""
		default:
			fmt.Println(color.RedString("❓ Invalid choice. Please try again."))
		}
	}
}

// isRunningInTerminal checks if the program is running in an interactive terminal
func isRunningInTerminal() bool {
	// Check if stdin is a terminal
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	
	// Check terminal mode bits
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// writeToCommitMsgFile writes the message to the Git commit message file
func writeToCommitMsgFile(message string, filePath string) error {
	return ioutil.WriteFile(filePath, []byte(message), 0644)
} 