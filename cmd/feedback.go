package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/feedback"
	"github.com/AccursedGalaxy/noidea/internal/history"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Feedback command flags
	countFlag        int
	authorFlag       string
	branchFlag       string
	filesFlag        string
	feedbackPersonalityFlag string
	includeDiffFlag  bool
	exportFeedbackFlag string
)

func init() {
	rootCmd.AddCommand(feedbackCmd)

	// Add flags
	feedbackCmd.Flags().IntVarP(&countFlag, "count", "c", 5, "Number of recent commits to analyze (default: 5)")
	feedbackCmd.Flags().StringVarP(&authorFlag, "author", "a", "", "Filter commits by author")
	feedbackCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "Filter commits by branch")
	feedbackCmd.Flags().StringVarP(&filesFlag, "files", "f", "", "Filter commits by files (comma-separated)")
	feedbackCmd.Flags().StringVarP(&feedbackPersonalityFlag, "personality", "p", "", "Personality to use for feedback (default: from config)")
	feedbackCmd.Flags().BoolVarP(&includeDiffFlag, "diff", "d", false, "Include diff context in analysis")
	feedbackCmd.Flags().StringVarP(&exportFeedbackFlag, "export", "e", "", "Export format: text, markdown, or html")
}

var feedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Get targeted feedback on your recent commits",
	Long:  `Analyze your Git commit history and provide targeted code quality and practice insights.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg := config.LoadConfig()

		// Create a history filter based on command flags
		filter := history.HistoryFilter{
			Count:       countFlag,
			Author:      authorFlag,
			Branch:      branchFlag,
			IncludeDiff: includeDiffFlag,
		}

		// Get personality name from flag or config
		personalityName := cfg.Moai.Personality
		if feedbackPersonalityFlag != "" {
			personalityName = feedbackPersonalityFlag
		}

		// Create a history collector
		collector, err := history.NewHistoryCollector()
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to create history collector:", err)
			return
		}

		// Get commit history
		commits, err := collector.GetCommitHistory(filter)
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to retrieve commit history:", err)
			return
		}

		// Check if we have any commits
		if len(commits) == 0 {
			fmt.Println(color.YellowString("No commits found matching the criteria."))
			return
		}

		// Filter by files if needed
		if filesFlag != "" {
			fileFilters := strings.Split(filesFlag, ",")
			commits = filterCommitsByFiles(commits, fileFilters)
			
			if len(commits) == 0 {
				fmt.Println(color.YellowString("No commits found matching the file filters."))
				return
			}
		}

		// Calculate statistics for the filtered commits
		stats := collector.CalculateStats(commits)
		
		// Display basic info
		fmt.Printf("%s %s\n\n", 
			color.CyanString("🔍 Analyzing"),
			color.CyanString(fmt.Sprintf("%d commits", len(commits))))
		
		// Generate insights
		feedback, err := generateOnDemandFeedback(commits, stats, personalityName, cfg)
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to generate feedback:", err)
			return
		}
		
		// Create the complete feedback output
		output := formatFeedback(commits, stats, feedback)
		
		// Export if requested, otherwise print to console
		if exportFeedbackFlag != "" {
			if err := exportFeedback(output, exportFeedbackFlag); err != nil {
				fmt.Println(color.RedString("Error:"), "Failed to export feedback:", err)
			} else {
				fmt.Println(color.GreenString("Feedback exported successfully."))
			}
		} else {
			// Print to console
			fmt.Println(output)
		}
	},
}

// generateOnDemandFeedback creates specialized feedback based on filtered commits
func generateOnDemandFeedback(commits []history.CommitInfo, stats map[string]interface{}, personalityName string, cfg config.Config) (string, error) {
	// Create a list of commit messages
	var commitMessages []string
	var commitDiffs []string
	
	for _, commit := range commits {
		commitMessages = append(commitMessages, commit.Message)
		if commit.DiffSummary != "" {
			commitDiffs = append(commitDiffs, commit.DiffSummary)
		}
	}
	
	// Create a context tailored for on-demand feedback
	feedbackContext := feedback.CommitContext{
		Message:       "On-Demand Feedback Analysis",
		Timestamp:     time.Now(),
		CommitHistory: commitMessages,
		CommitStats:   stats,
	}
	
	// Add diff context if available
	if len(commitDiffs) > 0 {
		feedbackContext.Diff = strings.Join(commitDiffs[:min(3, len(commitDiffs))], "\n---\n")
	}
	
	// Create feedback engine based on configuration
	engine := feedback.NewFeedbackEngine(
		cfg.LLM.Provider,
		cfg.LLM.Model,
		cfg.LLM.APIKey,
		personalityName,
		cfg.Moai.PersonalityFile,
	)
	
	// Reuse the summary feedback method since it's designed for multi-commit analysis
	return engine.GenerateSummaryFeedback(feedbackContext)
}

// formatFeedback combines all parts into a complete feedback report
func formatFeedback(commits []history.CommitInfo, stats map[string]interface{}, insights string) string {
	var result strings.Builder
	
	// Header
	result.WriteString(color.CyanString("📋 Commit Analysis Report\n\n"))
	
	// Add filter info
	var filterInfo []string
	if authorFlag != "" {
		filterInfo = append(filterInfo, fmt.Sprintf("Author: %s", authorFlag))
	}
	if branchFlag != "" {
		filterInfo = append(filterInfo, fmt.Sprintf("Branch: %s", branchFlag))
	}
	if filesFlag != "" {
		filterInfo = append(filterInfo, fmt.Sprintf("Files: %s", filesFlag))
	}
	
	if len(filterInfo) > 0 {
		result.WriteString(color.CyanString("Filters: ") + strings.Join(filterInfo, ", ") + "\n\n")
	}
	
	// Statistics Summary
	result.WriteString(color.CyanString("## Statistics\n\n"))
	result.WriteString(fmt.Sprintf("Total Commits: %d\n", len(commits)))
	
	if val, ok := stats["unique_authors"].(int); ok && val > 0 {
		result.WriteString(fmt.Sprintf("Unique Authors: %d\n", val))
	}
	
	if val, ok := stats["total_files_changed"].(int); ok {
		result.WriteString(fmt.Sprintf("Files Changed: %d\n", val))
	}
	
	if ins, ok := stats["total_insertions"].(int); ok {
		if del, ok := stats["total_deletions"].(int); ok {
			result.WriteString(fmt.Sprintf("Lines Added: %d, Removed: %d (Net: %d)\n", 
				ins, del, ins-del))
		}
	}
	
	result.WriteString("\n")
	
	// AI Insights
	result.WriteString(color.CyanString("## Analysis\n\n"))
	result.WriteString(insights)
	result.WriteString("\n\n")
	
	// Commit List
	result.WriteString(color.CyanString("## Commits Analyzed\n\n"))
	result.WriteString(history.FormatCommitList(commits))
	
	return result.String()
}

// filterCommitsByFiles filters commits to only include those touching specified files
func filterCommitsByFiles(commits []history.CommitInfo, fileFilters []string) []history.CommitInfo {
	if len(fileFilters) == 0 {
		return commits
	}
	
	var filtered []history.CommitInfo
	
	for _, commit := range commits {
		for _, file := range commit.Files {
			for _, filter := range fileFilters {
				if strings.Contains(file, strings.TrimSpace(filter)) {
					filtered = append(filtered, commit)
					break // Break the innermost loop once matched
				}
			}
		}
	}
	
	return filtered
}

// exportFeedback exports the feedback to a file in the requested format
func exportFeedback(feedback, format string) error {
	// Use the same export functionality as the summary command
	return exportSummary(feedback, format)
}

// min returns the smaller of a and b
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 