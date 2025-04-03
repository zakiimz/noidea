package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/feedback"
	"github.com/AccursedGalaxy/noidea/internal/moai"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Flag to enable/disable AI feedback
	useAI bool
	// Flag to get the diff of the last commit
	includeDiff bool
)

func init() {
	rootCmd.AddCommand(moaiCmd)

	// Add flags
	moaiCmd.Flags().BoolVarP(&useAI, "ai", "a", false, "Use AI to generate feedback")
	moaiCmd.Flags().BoolVarP(&includeDiff, "diff", "d", false, "Include the diff in AI context")
}

var moaiCmd = &cobra.Command{
	Use:   "moai [commit message]",
	Short: "Display a Moai with feedback on your commit",
	Long:  `Show a Moai face and random feedback about your most recent commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		var commitMsg string
		var commitDiff string
		
		// If commit message was provided as args, use it
		if len(args) > 0 {
			commitMsg = strings.Join(args, " ")
		} else {
			// Otherwise, try to get the latest commit message
			gitCmd := exec.Command("git", "log", "-1", "--pretty=%B")
			output, err := gitCmd.Output()
			if err != nil {
				commitMsg = "unknown commit"
			} else {
				commitMsg = strings.TrimSpace(string(output))
			}
		}

		// If diff flag is set, get the diff too
		if includeDiff {
			gitCmd := exec.Command("git", "show", "--stat", "HEAD")
			output, err := gitCmd.Output()
			if err == nil {
				commitDiff = string(output)
			}
		}

		// Get the Moai face
		face := moai.GetRandomFace()

		// Load configuration
		cfg := config.LoadConfig()

		// Override AI flag from config if set
		if !useAI && cfg.LLMEnabled {
			useAI = true
		}

		// Display the commit message
		fmt.Printf("%s  %s\n", face, commitMsg)

		// Generate feedback based on AI flag
		if useAI {
			// Create commit context
			commitContext := feedback.CommitContext{
				Message:   commitMsg,
				Timestamp: time.Now(),
				Diff:      commitDiff,
			}

			// Create feedback engine based on configuration
			engine := feedback.NewFeedbackEngine(cfg.LLMProvider, cfg.LLMModel, cfg.APIKey)

			// Generate AI feedback
			aiResponse, err := engine.GenerateFeedback(commitContext)
			if err != nil {
				// On error, fallback to local feedback
				fmt.Println(color.YellowString(moai.GetRandomFeedback(commitMsg)))
				fmt.Println(color.RedString("AI Error:"), err)
			} else {
				// Display AI-generated feedback
				fmt.Println(color.CyanString(aiResponse))
			}
		} else {
			// Use local feedback
			fmt.Println(color.YellowString(moai.GetRandomFeedback(commitMsg)))
		}
	},
} 