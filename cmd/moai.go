package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/AccursedGalaxy/noidea/internal/moai"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(moaiCmd)
}

var moaiCmd = &cobra.Command{
	Use:   "moai [commit message]",
	Short: "Display a Moai with feedback on your commit",
	Long:  `Show a Moai face and random feedback about your most recent commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		var commitMsg string
		
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

		// Display the Moai face and feedback
		face := moai.GetRandomFace()
		feedback := moai.GetRandomFeedback(commitMsg)
		
		fmt.Printf("%s  %s\n", face, commitMsg)
		fmt.Println(feedback)
	},
} 