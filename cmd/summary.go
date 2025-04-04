package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/feedback"
	"github.com/AccursedGalaxy/noidea/internal/history"
	"github.com/AccursedGalaxy/noidea/internal/personality"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	// Summary command flags
	daysFlag              int
	allHistoryFlag        bool
	exportFlag            string
	statsOnlyFlag         bool
	aiInsightFlag         bool
	personalityForSummary string
	showCommitHistoryFlag bool
)

func init() {
	rootCmd.AddCommand(summaryCmd)

	// Add flags
	summaryCmd.Flags().IntVarP(&daysFlag, "days", "d", 7, "Number of days to include in summary (default: 7, use 0 for all history)")
	summaryCmd.Flags().BoolVarP(&allHistoryFlag, "all", "A", false, "Show complete repository history regardless of --days value")
	summaryCmd.Flags().StringVarP(&exportFlag, "export", "e", "", "Export format: text, markdown, or html")
	summaryCmd.Flags().BoolVarP(&statsOnlyFlag, "stats-only", "s", false, "Show only statistics without AI insights")
	summaryCmd.Flags().BoolVarP(&aiInsightFlag, "ai", "a", false, "Include AI insights (default: use config)")
	summaryCmd.Flags().StringVarP(&personalityForSummary, "personality", "p", "", "Personality to use for insights (default: from config)")
	summaryCmd.Flags().BoolVarP(&showCommitHistoryFlag, "show-commits", "c", false, "Include detailed commit history in the output")
}

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Generate a summary of recent Git activity",
	Long:  `Analyze your Git history and provide statistics and insights about your recent commits.

By default, this command shows commits from the last 7 days. If no commits are found
in this period, it automatically shows all repository history.

Examples:
  noidea summary                # Show commits from the last 7 days
  noidea summary --days 30      # Show commits from the last 30 days
  noidea summary --all          # Show all repository history
  noidea summary --days 0       # Same as --all, shows all history
  noidea summary --show-commits # Include detailed commit history in output`,
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

		var commits []history.CommitInfo
		var err error

		// Check if user requested all history
		if allHistoryFlag || daysFlag == 0 {
			// Fetch all commits
			commits, err = history.GetLastNCommits(1000, useAI)
			if err != nil {
				fmt.Println(color.RedString("Error:"), "Failed to retrieve commit history:", err)
				return
			}
			// Set days to a large value to indicate complete history in the summary
			daysFlag = 365 * 10 // 10 years, arbitrary large number
		} else {
			// Get commit data for the specified period
			commits, err = history.GetCommitsFromLastNDays(daysFlag, useAI)
			if err != nil {
				fmt.Println(color.RedString("Error:"), "Failed to retrieve commit history:", err)
				return
			}

			// Only show the fallback message and fetch all history if we truly have zero commits
			if len(commits) == 0 {
				// No commits in the specified time period, automatically fetch all history
				fmt.Println(color.YellowString("No commits found in the last"), 
					color.CyanString(strconv.Itoa(daysFlag)), 
					color.YellowString("days. Showing complete history instead."))
				
				// Get all commits by using GetLastNCommits with a high number
				commits, err = history.GetLastNCommits(1000, useAI)
				if err != nil {
					fmt.Println(color.RedString("Error:"), "Failed to retrieve commit history:", err)
					return
				}
				
				// Set days to a large value to indicate complete history in the summary
				daysFlag = 365 * 10 // 10 years, arbitrary large number
			}
		}

		// Use a direct Git command to get commits as a test
		if len(commits) == 0 {
			// Execute a direct Git command to see if we can get commits
			cmd := exec.Command("git", "log", "--pretty=format:%s", "-n", "10")
			out, err := cmd.Output()
			if err == nil && len(out) > 0 {
				// We got direct git output but no commits from our history function
				fmt.Println(color.YellowString("Warning:"), "Git history is available but our history collector couldn't retrieve it.")
				fmt.Println(color.YellowString("Recent commits from Git:"))
				fmt.Println(string(out))
			}
		}

		// Check if we have any commits after all attempts
		if len(commits) == 0 {
			fmt.Println(color.YellowString("No commits found in this repository."))
			return
		}

		// If showing all history, update the days value to reflect the actual time span
		if daysFlag >= 365*10 && len(commits) > 0 {
			// Find the oldest commit timestamp
			oldestCommit := commits[len(commits)-1].Timestamp
			days := int(time.Since(oldestCommit).Hours() / 24) + 1
			
			// Only update if it's less than our arbitrary large number
			if days < 365*10 {
				daysFlag = days
			}
		}

		// Generate statistics
		// Directly get Git stats if collector doesn't provide data
		collector, err := history.NewHistoryCollector()
		if err != nil {
			fmt.Println(color.RedString("Error:"), "Failed to create history collector:", err)
			return
		}
		stats := collector.CalculateStats(commits)

		// Verify stats are not all zero
		allZeros := true
		if val, ok := stats["totalCommits"]; ok && val != nil {
			if v, ok := val.(int); ok && v > 0 {
				allZeros = false
			}
		}

		// If stats appear to be all zeros but we have commits, try to get stats directly
		if allZeros && len(commits) > 0 {
			// Directly calculate basic stats
			stats["totalCommits"] = len(commits)
			
			// Calculate unique authors
			authors := make(map[string]bool)
			for _, commit := range commits {
				authors[commit.Author] = true
			}
			stats["uniqueAuthors"] = len(authors)
			
			// Calculate timespan in hours
			if len(commits) >= 2 {
				newest := commits[0].Timestamp
				oldest := commits[len(commits)-1].Timestamp
				timeSpan := newest.Sub(oldest).Hours()
				stats["timeSpan"] = fmt.Sprintf("%.1f", timeSpan)
			} else {
				stats["timeSpan"] = "0.0"
			}
			
			// Calculate commits by day
			commitsByDay := make(map[string]int)
			for _, commit := range commits {
				day := commit.Timestamp.Weekday().String()
				commitsByDay[day]++
			}
			stats["commitsByDay"] = commitsByDay
			
			// Calculate commits by hour range
			commitsByHourRange := make(map[string]int)
			for _, commit := range commits {
				hour := commit.Timestamp.Hour()
				var hourRange string
				
				switch {
				case hour >= 4 && hour < 8:
					hourRange = "Morning (4-8)"
				case hour >= 8 && hour < 12:
					hourRange = "Work Hours (8-12)"
				case hour >= 12 && hour < 16:
					hourRange = "Afternoon (12-16)"
				case hour >= 16 && hour < 20:
					hourRange = "Evening (16-20)"
				case hour >= 20 && hour < 24:
					hourRange = "Late PM (20-24)"
				default:
					hourRange = "Night (0-4)"
				}
				
				commitsByHourRange[hourRange]++
			}
			stats["commitsByHourRange"] = commitsByHourRange
			
			// Try to get file stats using git command
			cmd := exec.Command("git", "diff", "--shortstat", commits[len(commits)-1].Hash, commits[0].Hash)
			out, err := cmd.Output()
			if err == nil {
				// Parse output like: " 10 files changed, 100 insertions(+), 50 deletions(-)"
				statStr := string(out)
				filesRe := regexp.MustCompile(`(\d+) files? changed`)
				addRe := regexp.MustCompile(`(\d+) insertions?\(\+\)`)
				delRe := regexp.MustCompile(`(\d+) deletions?\(-\)`)
				
				if matches := filesRe.FindStringSubmatch(statStr); len(matches) > 1 {
					if val, err := strconv.Atoi(matches[1]); err == nil {
						stats["filesChanged"] = val
					}
				}
				
				if matches := addRe.FindStringSubmatch(statStr); len(matches) > 1 {
					if val, err := strconv.Atoi(matches[1]); err == nil {
						stats["linesAdded"] = val
					}
				}
				
				if matches := delRe.FindStringSubmatch(statStr); len(matches) > 1 {
					if val, err := strconv.Atoi(matches[1]); err == nil {
						stats["linesRemoved"] = val
					}
				}
				
				// Calculate net change
				added := 0
				if val, ok := stats["linesAdded"].(int); ok {
					added = val
				}
				
				removed := 0
				if val, ok := stats["linesRemoved"].(int); ok {
					removed = val
				}
				
				stats["netChange"] = added - removed
			}
		}

		// Format statistics and get basic summary
		statsSummary := formatStatsForDisplay(stats)

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
		summary := formatSummary(statsSummary, commitList, aiInsight, daysFlag, showCommitHistoryFlag)

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
	// Check if we have any commits to analyze
	if len(commits) == 0 {
		// If no commits found, return a simple message
		return "No commits found in the specified time period to analyze.", nil
	}

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
	
	// Load personality configuration to modify
	personalities, err := personality.LoadPersonalities(cfg.Moai.PersonalityFile)
	if err != nil {
		// Fall back to default personalities if there's an error
		personalities = personality.DefaultPersonalities()
	}

	// Get the selected personality
	selectedPersonality, err := personalities.GetPersonality(personalityName)
	if err != nil {
		// Fall back to default personality
		selectedPersonality, _ = personalities.GetPersonality("")
	}

	// Get terminal width for dynamic token calculation
	width := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		width = w
	}
	
	// Calculate approximate tokens based on width
	// A rough estimate: 4 chars per token, with box borders taking about 4 chars
	// Also limit to max 60 chars per line for readability
	charsPerLine := minInt(width-4, 60)
	
	// For a box with N chars per line, we can fit approximately:
	// - About 10 lines of content maximum (conservative estimate)
	// This gives us about (10 * charsPerLine) / 4 tokens total
	// We'll be even more conservative with our estimate to ensure it fits
	maxTokens := (10 * charsPerLine) / 5 // Using 5 chars per token to be conservative
	
	// Constrain between reasonable min and max values
	maxTokens = minInt(maxInt(maxTokens, 80), 160)

	// Create a modified version with dynamic token limit
	modifiedPersonality := selectedPersonality
	modifiedPersonality.MaxTokens = maxTokens
	
	// Create a very targeted prompt to ensure the response fits
	modifiedPersonality.SystemPrompt = fmt.Sprintf(
		"You are a Git expert providing extremely concise insights. Respond with EXACTLY TWO short observations and ONE specific recommendation, formatted as bullet points. Keep each bullet to 1-2 sentences. Ensure your entire response is under %d words - this is critical. Be direct and clear.",
		maxTokens/2, // rough estimate of words based on tokens
	)

	// Create feedback engine with the modified personality
	engine := feedback.NewFeedbackEngineWithCustomPersonality(
		cfg.LLM.Provider,
		cfg.LLM.Model,
		cfg.LLM.APIKey,
		modifiedPersonality,
	)

	// Generate AI insights
	return engine.GenerateSummaryFeedback(summaryContext)
}

// Helper function for min/max
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// formatSummary combines all parts into a complete summary
func formatSummary(stats, commits, aiInsights string, days int, showHistory bool) string {
	var result strings.Builder

	// Get terminal width for better formatting
	width := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		width = w
	}

	// Create styled boxes
	boxStylePrimary := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#2980b9")).
		Padding(0, 1).
		Width(width - 4)

	boxStyleSecondary := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#27ae60")).
		Padding(0, 1).
		Width(width - 4)

	subHeaderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2980b9")).
		Bold(true)

	// Statistics section with combined date range and header
	var statsHeader string
	if days >= 365*10 || days == 0 {
		statsHeader = subHeaderStyle.Render("Git Statistics: Complete repository history")
	} else {
		statsHeader = subHeaderStyle.Render(fmt.Sprintf("Git Statistics: Last %d days (%s to %s)",
			days,
			time.Now().AddDate(0, 0, -days).Format("2006-01-02"),
			time.Now().Format("2006-01-02")))
	}
	result.WriteString(statsHeader + "\n")
	result.WriteString(boxStylePrimary.Render(stats))
	result.WriteString("\n\n")

	// AI Insights (if available)
	if aiInsights != "" {
		insightsHeader := subHeaderStyle.Render("AI Insights")
		result.WriteString(insightsHeader + "\n")
		result.WriteString(boxStyleSecondary.Render(aiInsights))
		result.WriteString("\n\n")
	}

	// Only include commit history if explicitly requested
	if showHistory {
		historyHeader := subHeaderStyle.Render("Commit History")
		result.WriteString(historyHeader + "\n")
		result.WriteString(boxStylePrimary.Render(commits))
	}

	return result.String()
}

// Format the stats sections in a more visually appealing way
func formatStatsForDisplay(stats map[string]interface{}) string {
	var result strings.Builder

	// Basic stats with highlighted numbers - with nil checks
	totalCommits := safeGetValue(stats, "totalCommits", "0")
	timeSpan := safeGetValue(stats, "timeSpan", "0")
	uniqueAuthors := safeGetValue(stats, "uniqueAuthors", "0")
	
	result.WriteString(fmt.Sprintf("Total Commits: %s\n", color.New(color.FgHiGreen, color.Bold).Sprint(totalCommits)))
	result.WriteString(fmt.Sprintf("Time Span: %s hours\n", color.New(color.FgHiGreen, color.Bold).Sprint(timeSpan)))
	result.WriteString(fmt.Sprintf("Unique Authors: %s\n\n", color.New(color.FgHiGreen, color.Bold).Sprint(uniqueAuthors)))

	// File changes with highlighted numbers - with nil checks
	filesChanged := safeGetValue(stats, "filesChanged", "0")
	linesAdded := safeGetValue(stats, "linesAdded", "0")
	linesRemoved := safeGetValue(stats, "linesRemoved", "0")
	
	result.WriteString(fmt.Sprintf("Files Changed: %s\n", color.New(color.FgHiYellow, color.Bold).Sprint(filesChanged)))
	result.WriteString(fmt.Sprintf("Lines Added: %s\n", color.New(color.FgGreen, color.Bold).Sprint(linesAdded)))
	result.WriteString(fmt.Sprintf("Lines Removed: %s\n", color.New(color.FgRed, color.Bold).Sprint(linesRemoved)))
	
	netChange := 0
	if val, ok := stats["netChange"]; ok && val != nil {
		if intVal, ok := val.(int); ok {
			netChange = intVal
		}
	}
	
	netChangeColor := color.New(color.Bold)
	if netChange > 0 {
		netChangeColor = color.New(color.FgGreen, color.Bold)
	} else if netChange < 0 {
		netChangeColor = color.New(color.FgRed, color.Bold)
	}
	
	result.WriteString(fmt.Sprintf("Net Change: %s\n\n", netChangeColor.Sprint(netChange)))

	// Commits by day section
	result.WriteString(color.New(color.FgHiMagenta, color.Bold).Sprint("📅 Commits by Day:\n"))
	
	if commitsByDay, ok := stats["commitsByDay"].(map[string]int); ok && commitsByDay != nil {
		maxDay := 0
		for _, count := range commitsByDay {
			if count > maxDay {
				maxDay = count
			}
		}
		
		// Days of week in order
		daysOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		
		for _, day := range daysOrder {
			if count, exists := commitsByDay[day]; exists && count > 0 {
				barLength := int(float64(count) / float64(maxDay) * 50)
				if maxDay == 0 {
					barLength = 0
				}
				bar := strings.Repeat("█", barLength)
				dayFormatted := fmt.Sprintf("%-10s", day)
				result.WriteString(fmt.Sprintf("%s : %s %s\n", 
					color.New(color.FgHiWhite).Sprint(dayFormatted),
					color.New(color.FgBlue).Sprint(bar),
					color.New(color.FgHiBlue).Sprintf("(%d)", count)))
			}
		}
		result.WriteString("\n")
	}

	// Commits by hour with emoji
	result.WriteString(color.New(color.FgHiCyan, color.Bold).Sprint("🕒 Commits by Hour:\n"))
	
	if commitsByHour, ok := stats["commitsByHourRange"].(map[string]int); ok && commitsByHour != nil {
		maxHour := 0
		for _, count := range commitsByHour {
			if count > maxHour {
				maxHour = count
			}
		}
		
		// Hour ranges in chronological order
		hourRanges := []string{"Morning (4-8)", "Work Hours (8-12)", "Afternoon (12-16)", "Evening (16-20)", "Late PM (20-24)", "Night (0-4)"}
		
		for _, hourRange := range hourRanges {
			if count, exists := commitsByHour[hourRange]; exists && count > 0 {
				barLength := int(float64(count) / float64(maxHour) * 50)
				if maxHour == 0 {
					barLength = 0
				}
				bar := strings.Repeat("█", barLength)
				rangeFormatted := fmt.Sprintf("%-16s", hourRange)
				result.WriteString(fmt.Sprintf("%s : %s %s\n", 
					color.New(color.FgHiWhite).Sprint(rangeFormatted),
					color.New(color.FgCyan).Sprint(bar),
					color.New(color.FgHiCyan).Sprintf("(%d)", count)))
			}
		}
	}

	return result.String()
}

// safeGetValue safely extracts a value from a map, returning defaultValue if nil or not found
func safeGetValue(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key]; ok && val != nil {
		return fmt.Sprintf("%v", val)
	}
	return defaultValue
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
		return os.WriteFile(filename, []byte(plainSummary), 0644)

	case "markdown", "md":
		filename = fmt.Sprintf("git-summary-%s.md", timestamp)
		return os.WriteFile(filename, []byte(convertToMarkdown(plainSummary)), 0644)

	case "html":
		filename = fmt.Sprintf("git-summary-%s.html", timestamp)
		return os.WriteFile(filename, []byte(convertToHTML(plainSummary)), 0644)

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
