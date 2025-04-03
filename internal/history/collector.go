package history

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// CommitInfo represents metadata about a single git commit
type CommitInfo struct {
	Hash        string      `json:"hash"`
	Author      string      `json:"author"`
	Email       string      `json:"email"`
	Timestamp   time.Time   `json:"timestamp"`
	Message     string      `json:"message"`
	Files       []string    `json:"files"`
	Stats       CommitStats `json:"stats"`
	DiffSummary string      `json:"diff_summary,omitempty"`
}

// CommitStats holds statistics about files changed in a commit
type CommitStats struct {
	FilesChanged int `json:"files_changed"`
	Insertions   int `json:"insertions"`
	Deletions    int `json:"deletions"`
}

// HistoryFilter defines parameters for filtering git history
type HistoryFilter struct {
	// Either use Since (time-based) or Count (number-based) filtering
	Since       time.Duration // e.g., 7*24*time.Hour for 7 days
	Count       int           // e.g., 10 for last 10 commits
	Author      string        // Filter by author, empty for all authors
	Branch      string        // Filter by branch, empty for current branch
	IncludeDiff bool          // Whether to include diff summaries
}

// HistoryCollector provides methods to fetch and analyze git commit history
type HistoryCollector struct {
	cacheDir  string
	cacheFile string
	cached    map[string]CommitInfo
}

// NewHistoryCollector creates a new collector with optional caching
func NewHistoryCollector() (*HistoryCollector, error) {
	// Setup cache directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	cacheDir := filepath.Join(home, ".noidea", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	cacheFile := filepath.Join(cacheDir, "history_cache.json")

	collector := &HistoryCollector{
		cacheDir:  cacheDir,
		cacheFile: cacheFile,
		cached:    make(map[string]CommitInfo),
	}

	// Load cache if exists
	collector.loadCache()

	return collector, nil
}

// loadCache attempts to load the commit cache from disk
func (h *HistoryCollector) loadCache() {
	data, err := os.ReadFile(h.cacheFile)
	if err != nil {
		// Cache doesn't exist yet, that's fine
		return
	}

	if err := json.Unmarshal(data, &h.cached); err != nil {
		// If cache is corrupted, start fresh
		h.cached = make(map[string]CommitInfo)
	}
}

// saveCache persists the commit cache to disk
func (h *HistoryCollector) saveCache() error {
	data, err := json.Marshal(h.cached)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	return os.WriteFile(h.cacheFile, data, 0644)
}

// GetCommitHistory retrieves commit history based on the provided filter
func (h *HistoryCollector) GetCommitHistory(filter HistoryFilter) ([]CommitInfo, error) {
	var args []string

	// Base command to get commit hashes
	args = append(args, "log", "--format=%H")

	// Apply filters
	if filter.Since != 0 {
		// Time-based filtering
		sinceStr := fmt.Sprintf("--since=%s", filter.Since.String())
		args = append(args, sinceStr)
	} else if filter.Count > 0 {
		// Count-based filtering
		args = append(args, fmt.Sprintf("-n%d", filter.Count))
	} else {
		// Default to last 10 commits if no filter specified
		args = append(args, "-n10")
	}

	// Author filter
	if filter.Author != "" {
		args = append(args, fmt.Sprintf("--author=%s", filter.Author))
	}

	// Branch filter
	if filter.Branch != "" {
		args = append(args, filter.Branch)
	}

	// Execute git command
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit hashes: %w", err)
	}

	// Parse commit hashes
	hashes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(hashes) == 0 || (len(hashes) == 1 && hashes[0] == "") {
		return nil, nil // No commits found
	}

	// Collect commit info for each hash
	commits := make([]CommitInfo, 0, len(hashes))
	for _, hash := range hashes {
		if hash == "" {
			continue
		}

		// Check cache first
		if commit, found := h.cached[hash]; found {
			// If we need diff but it's not in cache, we'll fetch it
			if filter.IncludeDiff && commit.DiffSummary == "" {
				diffSummary, err := h.getDiffSummary(hash)
				if err == nil {
					commit.DiffSummary = diffSummary
					h.cached[hash] = commit // Update cache
				}
			}

			commits = append(commits, commit)
			continue
		}

		// Fetch commit info for uncached commits
		commit, err := h.getCommitInfo(hash, filter.IncludeDiff)
		if err != nil {
			// Skip commits that can't be retrieved
			continue
		}

		// Add to cache
		h.cached[hash] = commit
		commits = append(commits, commit)
	}

	// Save updated cache
	h.saveCache()

	return commits, nil
}

// getCommitInfo fetches detailed info for a specific commit
func (h *HistoryCollector) getCommitInfo(hash string, includeDiff bool) (CommitInfo, error) {
	var commit CommitInfo
	commit.Hash = hash

	// Get commit metadata
	cmd := exec.Command("git", "show", "--format=%an%n%ae%n%at%n%B", "--name-only", hash)
	output, err := cmd.Output()
	if err != nil {
		return commit, fmt.Errorf("failed to get commit metadata: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 4 {
		return commit, fmt.Errorf("invalid commit data format")
	}

	commit.Author = lines[0]
	commit.Email = lines[1]

	// Parse timestamp
	timestamp, err := strconv.ParseInt(lines[2], 10, 64)
	if err != nil {
		return commit, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	commit.Timestamp = time.Unix(timestamp, 0)

	// Parse message (might span multiple lines)
	var messageBuilder strings.Builder
	lineIndex := 3
	for ; lineIndex < len(lines); lineIndex++ {
		if lines[lineIndex] == "" {
			break
		}
		if messageBuilder.Len() > 0 {
			messageBuilder.WriteString("\n")
		}
		messageBuilder.WriteString(lines[lineIndex])
	}
	commit.Message = messageBuilder.String()

	// Skip any blank lines
	for ; lineIndex < len(lines) && lines[lineIndex] == ""; lineIndex++ {
	}

	// Collect changed files
	for ; lineIndex < len(lines); lineIndex++ {
		if lines[lineIndex] != "" {
			commit.Files = append(commit.Files, lines[lineIndex])
		}
	}

	// Get commit stats
	commit.Stats = h.getCommitStats(hash)

	// Get diff summary if requested
	if includeDiff {
		diffSummary, err := h.getDiffSummary(hash)
		if err == nil {
			commit.DiffSummary = diffSummary
		}
	}

	return commit, nil
}

// getCommitStats retrieves stats about files changed in the commit
func (h *HistoryCollector) getCommitStats(hash string) CommitStats {
	var stats CommitStats

	// Run git show with stat option
	cmd := exec.Command("git", "show", "--stat", hash)
	output, err := cmd.Output()
	if err != nil {
		return stats
	}

	// Extract stats from the last line which looks like:
	// " 3 files changed, 24 insertions(+), 4 deletions(-)"
	lines := strings.Split(string(output), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.Contains(line, "file") && strings.Contains(line, "changed") {
			// Parse stats
			parts := strings.Split(line, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasSuffix(part, "files changed") || strings.HasSuffix(part, "file changed") {
					numStr := strings.TrimSuffix(strings.TrimSuffix(part, "files changed"), "file changed")
					stats.FilesChanged, _ = strconv.Atoi(strings.TrimSpace(numStr))
				} else if strings.HasSuffix(part, "insertions(+)") || strings.HasSuffix(part, "insertion(+)") {
					numStr := strings.TrimSuffix(strings.TrimSuffix(part, "insertions(+)"), "insertion(+)")
					stats.Insertions, _ = strconv.Atoi(strings.TrimSpace(numStr))
				} else if strings.HasSuffix(part, "deletions(-)") || strings.HasSuffix(part, "deletion(-)") {
					numStr := strings.TrimSuffix(strings.TrimSuffix(part, "deletions(-)"), "deletion(-)")
					stats.Deletions, _ = strconv.Atoi(strings.TrimSpace(numStr))
				}
			}
			break
		}
	}

	return stats
}

// getDiffSummary generates a summarized version of the diff for LLM consumption
func (h *HistoryCollector) getDiffSummary(hash string) (string, error) {
	// Get the diff with context
	cmd := exec.Command("git", "show", hash)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}

	// For LLM consumption, we want to keep the diff relatively concise
	// This is a simple approach - for a more sophisticated summary,
	// you might want to process the diff to extract only key changes
	diffText := string(output)

	// Truncate very large diffs to avoid token explosion when sent to LLMs
	const maxDiffLength = 5000
	if len(diffText) > maxDiffLength {
		return diffText[:maxDiffLength] + "... [diff truncated]", nil
	}

	return diffText, nil
}

// GetCommitRange retrieves commits between two dates
func (h *HistoryCollector) GetCommitRange(startTime, endTime time.Time) ([]CommitInfo, error) {
	args := []string{
		"log",
		"--format=%H",
		fmt.Sprintf("--since=%s", startTime.Format(time.RFC3339)),
		fmt.Sprintf("--until=%s", endTime.Format(time.RFC3339)),
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit range: %w", err)
	}

	hashes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(hashes) == 0 || (len(hashes) == 1 && hashes[0] == "") {
		return nil, nil // No commits in range
	}

	commits := make([]CommitInfo, 0, len(hashes))
	for _, hash := range hashes {
		if hash == "" {
			continue
		}

		// Check cache first
		if commit, found := h.cached[hash]; found {
			commits = append(commits, commit)
			continue
		}

		// Fetch commit info
		commit, err := h.getCommitInfo(hash, false)
		if err != nil {
			continue
		}

		// Add to cache
		h.cached[hash] = commit
		commits = append(commits, commit)
	}

	h.saveCache()

	return commits, nil
}

// CalculateStats generates aggregated statistics for a set of commits
func (h *HistoryCollector) CalculateStats(commits []CommitInfo) map[string]interface{} {
	stats := make(map[string]interface{})

	if len(commits) == 0 {
		return stats
	}

	// Basic counts
	stats["total_commits"] = len(commits)

	// Time range
	earliest := commits[len(commits)-1].Timestamp
	latest := commits[0].Timestamp
	stats["time_span_hours"] = latest.Sub(earliest).Hours()

	// Author stats
	authors := make(map[string]int)
	for _, c := range commits {
		authors[c.Author]++
	}
	stats["unique_authors"] = len(authors)
	stats["author_distribution"] = authors

	// File stats
	totalFiles := 0
	totalInsertions := 0
	totalDeletions := 0
	for _, c := range commits {
		totalFiles += len(c.Files)
		totalInsertions += c.Stats.Insertions
		totalDeletions += c.Stats.Deletions
	}
	stats["total_files_changed"] = totalFiles
	stats["total_insertions"] = totalInsertions
	stats["total_deletions"] = totalDeletions

	// Commits by day of week
	dayOfWeek := make(map[string]int)
	for _, c := range commits {
		day := c.Timestamp.Weekday().String()
		dayOfWeek[day]++
	}
	stats["commits_by_day"] = dayOfWeek

	// Commits by hour
	hourOfDay := make(map[int]int)
	for _, c := range commits {
		hour := c.Timestamp.Hour()
		hourOfDay[hour]++
	}
	stats["commits_by_hour"] = hourOfDay

	return stats
}

// ClearCache removes the cache file
func (h *HistoryCollector) ClearCache() error {
	h.cached = make(map[string]CommitInfo)
	if _, err := os.Stat(h.cacheFile); err == nil {
		return os.Remove(h.cacheFile)
	}
	return nil
}
