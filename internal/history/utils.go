package history

import (
	"fmt"
	"time"
)

// GetCommitsFromLastNDays retrieves commits from the past N days
func GetCommitsFromLastNDays(days int, includeDiff bool) ([]CommitInfo, error) {
	collector, err := NewHistoryCollector()
	if err != nil {
		return nil, fmt.Errorf("failed to create history collector: %w", err)
	}

	filter := HistoryFilter{
		Since:       time.Duration(days) * 24 * time.Hour,
		IncludeDiff: includeDiff,
	}

	return collector.GetCommitHistory(filter)
}

// GetLastNCommits retrieves the last N commits
func GetLastNCommits(count int, includeDiff bool) ([]CommitInfo, error) {
	collector, err := NewHistoryCollector()
	if err != nil {
		return nil, fmt.Errorf("failed to create history collector: %w", err)
	}

	filter := HistoryFilter{
		Count:       count,
		IncludeDiff: includeDiff,
	}

	return collector.GetCommitHistory(filter)
}

// FormatCommitSummary creates a human-readable summary of a commit
func FormatCommitSummary(commit CommitInfo) string {
	timeStr := commit.Timestamp.Format("2006-01-02 15:04:05")

	summary := fmt.Sprintf("Commit: %s\n", commit.Hash[:8])
	summary += fmt.Sprintf("Author: %s <%s>\n", commit.Author, commit.Email)
	summary += fmt.Sprintf("Date: %s\n\n", timeStr)
	summary += fmt.Sprintf("%s\n\n", commit.Message)

	summary += fmt.Sprintf("Files changed: %d\n", commit.Stats.FilesChanged)
	summary += fmt.Sprintf("Insertions: %d\n", commit.Stats.Insertions)
	summary += fmt.Sprintf("Deletions: %d\n", commit.Stats.Deletions)

	return summary
}

// FormatCommitList creates a concise summary list of commits
func FormatCommitList(commits []CommitInfo) string {
	if len(commits) == 0 {
		return "No commits found."
	}

	var summary string
	for i, commit := range commits {
		shortHash := commit.Hash
		if len(shortHash) > 8 {
			shortHash = shortHash[:8]
		}

		date := commit.Timestamp.Format("2006-01-02")
		time := commit.Timestamp.Format("15:04:05")

		// Truncate message if too long
		message := commit.Message
		if len(message) > 50 {
			message = message[:47] + "..."
		}

		summary += fmt.Sprintf("%d. [%s] %s %s - %s\n",
			i+1, shortHash, date, time, message)
	}

	return summary
}

// GetWeeklyStats gets stats for the last week with diffs if requested
func GetWeeklyStats(includeDiff bool) ([]CommitInfo, map[string]interface{}, error) {
	commits, err := GetCommitsFromLastNDays(7, includeDiff)
	if err != nil {
		return nil, nil, err
	}

	collector, _ := NewHistoryCollector()
	stats := collector.CalculateStats(commits)

	return commits, stats, nil
}

// FormatStatsForDisplay formats statistics into a human-readable string
func FormatStatsForDisplay(stats map[string]interface{}) string {
	result := "📊 Commit Statistics:\n\n"

	// Basic stats
	if total, ok := stats["total_commits"].(int); ok {
		result += fmt.Sprintf("Total Commits: %d\n", total)
	}

	if timeSpan, ok := stats["time_span_hours"].(float64); ok {
		days := timeSpan / 24
		if days < 1 {
			result += fmt.Sprintf("Time Span: %.1f hours\n", timeSpan)
		} else {
			result += fmt.Sprintf("Time Span: %.1f days\n", days)
		}
	}

	if authors, ok := stats["unique_authors"].(int); ok {
		result += fmt.Sprintf("Unique Authors: %d\n", authors)
	}

	// File stats
	if files, ok := stats["total_files_changed"].(int); ok {
		result += fmt.Sprintf("\nFiles Changed: %d\n", files)
	}

	if ins, ok := stats["total_insertions"].(int); ok {
		if del, ok := stats["total_deletions"].(int); ok {
			result += fmt.Sprintf("Lines Added: %d\n", ins)
			result += fmt.Sprintf("Lines Removed: %d\n", del)
			result += fmt.Sprintf("Net Change: %d\n", ins-del)
		}
	}

	// Day of week distribution
	if daysMap, ok := stats["commits_by_day"].(map[string]int); ok && len(daysMap) > 0 {
		result += "\n📅 Commits by Day:\n"
		// Days in order
		days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		for _, day := range days {
			if count, exists := daysMap[day]; exists {
				// Basic ASCII bar
				bar := ""
				for i := 0; i < count; i++ {
					bar += "█"
				}
				result += fmt.Sprintf("%-10s: %s (%d)\n", day, bar, count)
			}
		}
	}

	// Hour distribution (simplified)
	if hoursMap, ok := stats["commits_by_hour"].(map[int]int); ok && len(hoursMap) > 0 {
		result += "\n🕒 Commits by Hour:\n"

		// Group in 4-hour blocks for simplicity
		timeBlocks := map[string]int{
			"Night (0-4)":       0,
			"Early AM (4-8)":    0,
			"Morning (8-12)":    0,
			"Afternoon (12-16)": 0,
			"Evening (16-20)":   0,
			"Late PM (20-24)":   0,
		}

		for hour, count := range hoursMap {
			switch {
			case hour < 4:
				timeBlocks["Night (0-4)"] += count
			case hour < 8:
				timeBlocks["Early AM (4-8)"] += count
			case hour < 12:
				timeBlocks["Morning (8-12)"] += count
			case hour < 16:
				timeBlocks["Afternoon (12-16)"] += count
			case hour < 20:
				timeBlocks["Evening (16-20)"] += count
			default:
				timeBlocks["Late PM (20-24)"] += count
			}
		}

		// Print time blocks
		for name, count := range timeBlocks {
			if count > 0 {
				bar := ""
				for i := 0; i < count; i++ {
					bar += "█"
				}
				result += fmt.Sprintf("%-18s: %s (%d)\n", name, bar, count)
			}
		}
	}

	return result
}
