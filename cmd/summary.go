package cmd

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/feedback"
	"github.com/AccursedGalaxy/noidea/internal/history"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Summary command flags
	daysFlag              int
	exportFlag            string
	statsOnlyFlag         bool
	aiInsightFlag         bool
	personalityForSummary string
)

func init() {
	rootCmd.AddCommand(summaryCmd)

	// Add flags
	summaryCmd.Flags().IntVarP(&daysFlag, "days", "d", 7, "Number of days to include in summary (default: 7)")
	summaryCmd.Flags().StringVarP(&exportFlag, "export", "e", "", "Export format: text, markdown, or html")
	summaryCmd.Flags().BoolVarP(&statsOnlyFlag, "stats-only", "s", false, "Show only statistics without AI insights")
	summaryCmd.Flags().BoolVarP(&aiInsightFlag, "ai", "a", false, "Include AI insights (default: use config)")
	summaryCmd.Flags().StringVarP(&personalityForSummary, "personality", "p", "", "Personality to use for insights (default: from config)")
}

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Generate a summary of recent Git activity",
	Long:  `Analyze your Git history and provide statistics and insights about your recent commits.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg := config.LoadConfig()

		// Determine whether to use AI
		useAI := !statsOnlyFlag && (aiInsightFlag || cfg.LLM.Enabled)

		// Get personality name
		personalityName := cfg.Moai.Personality
		if personalityForSummary != "" {
			personalityName = personalityForSummary
		}

		// Get commit data for the specified period
		commits, err := history.GetCommitsFromLastNDays(daysFlag, useAI)
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to retrieve commit history:", err)
			return
		}

		// Check if we have any commits
		if len(commits) == 0 {
			fmt.Println(color.YellowString("No commits found in the last"), color.CyanString(strconv.Itoa(daysFlag)), color.YellowString("days."))
			return
		}

		// Generate statistics
		collector, err := history.NewHistoryCollector()
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to create history collector:", err)
			return
		}
		stats := collector.CalculateStats(commits)

		// Format statistics and get basic summary
		statsSummary := history.FormatStatsForDisplay(stats)

		// Get list of commits
		commitList := history.FormatCommitList(commits)

		var aiInsight string
		if useAI {
			aiInsight, err = generateAIInsights(commits, stats, personalityName, cfg)
			if err != nil {
				fmt.Println(color.YellowString("Note:"), "Unable to generate AI insights:", err)
			}
		}

		// Generate the complete summary
		summary := formatSummary(statsSummary, commitList, aiInsight, daysFlag)

		// Export if requested, otherwise print to console
		if exportFlag != "" {
			if err := exportSummary(summary, exportFlag); err != nil {
				fmt.Println(color.RedString("Error:"), "Failed to export summary:", err)
			} else {
				fmt.Println(color.GreenString("Summary exported successfully."))
			}
		} else {
			// Print to console
			fmt.Println(summary)
		}
	},
}

// generateAIInsights creates AI-powered insights for the commit history
func generateAIInsights(commits []history.CommitInfo, stats map[string]interface{}, personalityName string, cfg config.Config) (string, error) {
	// Build a condensed representation of commit messages
	var commitMessages []string
	for _, commit := range commits {
		commitMessages = append(commitMessages, commit.Message)
	}

	// Create summary context
	summaryContext := feedback.CommitContext{
		Message:       "Weekly Summary Analysis",
		Timestamp:     time.Now(),
		CommitHistory: commitMessages,
		CommitStats:   stats,
	}

	// Create feedback engine based on configuration
	engine := feedback.NewFeedbackEngine(
		cfg.LLM.Provider,
		cfg.LLM.Model,
		cfg.LLM.APIKey,
		personalityName,
		cfg.Moai.PersonalityFile,
	)

	// Generate AI insights
	return engine.GenerateSummaryFeedback(summaryContext)
}

// formatSummary combines all parts into a complete summary
func formatSummary(stats, commits, aiInsights string, days int) string {
	var result strings.Builder

	// Header
	result.WriteString(color.CyanString("📊 Git Activity Summary") + "\n")
	result.WriteString(color.CyanString(fmt.Sprintf("Last %d days - %s to %s\n\n",
		days,
		time.Now().AddDate(0, 0, -days).Format("2006-01-02"),
		time.Now().Format("2006-01-02"))))

	// Statistics
	result.WriteString(color.CyanString("## Statistics\n\n"))
	result.WriteString(stats)
	result.WriteString("\n")

	// AI Insights (if available)
	if aiInsights != "" {
		result.WriteString(color.CyanString("## AI Insights\n\n"))
		result.WriteString(aiInsights)
		result.WriteString("\n\n")
	}

	// Commit List
	result.WriteString(color.CyanString("## Commit History\n\n"))
	result.WriteString(commits)

	return result.String()
}

// exportSummary exports the summary to a file in the requested format
func exportSummary(summary, format string) error {
	// Determine output filename
	timestamp := time.Now().Format("2006-01-02")
	var filename string

	// Convert ANSI color codes to appropriate format
	plainSummary := stripANSIColors(summary)

	switch strings.ToLower(format) {
	case "text", "txt":
		filename = fmt.Sprintf("git-summary-%s.txt", timestamp)
		return ioutil.WriteFile(filename, []byte(plainSummary), 0644)

	case "markdown", "md":
		filename = fmt.Sprintf("git-summary-%s.md", timestamp)
		return ioutil.WriteFile(filename, []byte(convertToMarkdown(plainSummary)), 0644)

	case "html":
		filename = fmt.Sprintf("git-summary-%s.html", timestamp)
		return ioutil.WriteFile(filename, []byte(convertToHTML(plainSummary)), 0644)

	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// stripANSIColors removes ANSI color codes from a string
func stripANSIColors(s string) string {
	// Simple regex to remove ANSI color codes
	return color.New().SprintFunc()(s)
}

// convertToMarkdown converts the summary to Markdown format
func convertToMarkdown(summary string) string {
	lines := strings.Split(summary, "\n")
	var result strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "##") {
			// Convert "## Title" to "## Title"
			result.WriteString(line + "\n")
		} else if strings.HasPrefix(line, "📊 Git Activity Summary") {
			result.WriteString("# " + strings.TrimPrefix(line, "📊 ") + "\n")
		} else if strings.Contains(line, "Last") && strings.Contains(line, "days") {
			result.WriteString("*" + line + "*\n\n")
		} else {
			result.WriteString(line + "\n")
		}
	}

	return result.String()
}

// convertToHTML converts the summary to HTML format
func convertToHTML(summary string) string {
	markdown := convertToMarkdown(summary)

	// Simple HTML wrapper
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Git Activity Summary</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            color: #333;
        }
        h1 { color: #2c3e50; }
        h2 { color: #2980b9; }
        pre {
            background-color: #f5f5f5;
            padding: 10px;
            border-radius: 5px;
            overflow-x: auto;
        }
        .commit-list {
            list-style-type: none;
            padding-left: 0;
        }
        .commit-item {
            padding: 5px 0;
            border-bottom: 1px solid #eee;
        }
        .stats {
            display: flex;
            flex-wrap: wrap;
            gap: 20px;
            margin-bottom: 20px;
        }
        .stat-box {
            background-color: #f9f9f9;
            padding: 15px;
            border-radius: 5px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            flex: 1 1 200px;
        }
    </style>
</head>
<body>
    <div class="content">
        %s
    </div>
</body>
</html>
`

	// Replace Markdown with HTML tags
	htmlContent := strings.ReplaceAll(markdown, "# ", "<h1>")
	htmlContent = strings.ReplaceAll(htmlContent, "\n## ", "</h1>\n<h2>")
	htmlContent = strings.ReplaceAll(htmlContent, "\n", "<br>")
	htmlContent = strings.ReplaceAll(htmlContent, "</h1>", "</h1>")
	htmlContent = strings.ReplaceAll(htmlContent, "</h2>", "</h2>")

	return fmt.Sprintf(html, htmlContent)
}
