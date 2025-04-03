package feedback

import (
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/AccursedGalaxy/noidea/internal/moai"
)

// LocalFeedbackEngine generates feedback using pre-written local responses
type LocalFeedbackEngine struct{}

// NewLocalFeedbackEngine creates a new local feedback engine
func NewLocalFeedbackEngine() *LocalFeedbackEngine {
	return &LocalFeedbackEngine{}
}

// GenerateFeedback implements the FeedbackEngine interface
func (e *LocalFeedbackEngine) GenerateFeedback(ctx CommitContext) (string, error) {
	return moai.GetRandomFeedback(ctx.Message), nil
}

// GenerateSummaryFeedback provides basic insights for a weekly summary without using an LLM
func (e *LocalFeedbackEngine) GenerateSummaryFeedback(ctx CommitContext) (string, error) {
	summaries := []string{
		"Your commit history shows a consistent workflow. Keep up the good work!",
		"Looking at your commits, I notice you're making steady progress. Consider using more descriptive commit messages for better clarity.",
		"Your commit patterns suggest you're focused on quality. That's excellent!",
		"I see a mix of feature work and fixes in your commit history. Good balance!",
		"Your commit history indicates you're working through tasks methodically. Consider grouping related changes for cleaner history.",
		"Based on your commit times, you seem most productive in the middle of your work period. Interesting pattern!",
		"Your commits show attention to detail. Remember to take breaks too!",
		"I notice your commit messages are concise. For complex changes, a bit more detail might help future you.",
	}

	return summaries[rand.Intn(len(summaries))], nil
}

// GenerateCommitSuggestion creates a simple commit message suggestion based on diff stats
func (e *LocalFeedbackEngine) GenerateCommitSuggestion(ctx CommitContext) (string, error) {
	// Extract file paths from the diff
	lines := strings.Split(ctx.Diff, "\n")
	var filesChanged []string
	var fileExtensions = make(map[string]int)

	for _, line := range lines {
		if strings.HasPrefix(line, "+++ b/") {
			file := strings.TrimPrefix(line, "+++ b/")
			filesChanged = append(filesChanged, file)

			// Count file extensions
			ext := filepath.Ext(file)
			if ext != "" {
				fileExtensions[ext]++
			}
		}
	}

	// Determine the type prefix based on file types
	typePrefix := "feat"
	if len(fileExtensions) > 0 {
		// Find the most common extension
		maxCount := 0
		mostCommonExt := ""
		for ext, count := range fileExtensions {
			if count > maxCount {
				maxCount = count
				mostCommonExt = ext
			}
		}

		// Suggest type based on extension
		switch mostCommonExt {
		case ".md":
			typePrefix = "docs"
		case ".test.go", ".spec.js", "_test.go":
			typePrefix = "test"
		case ".css", ".scss", ".html", ".svg":
			typePrefix = "style"
		}
	}

	// Count number of modified lines (approximation)
	addedLines := 0
	removedLines := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			addedLines++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			removedLines++
		}
	}

	// Generate basic message based on files and line stats
	if len(filesChanged) == 0 {
		return typePrefix + ": update code", nil
	} else if len(filesChanged) == 1 {
		filename := filepath.Base(filesChanged[0])
		return typePrefix + ": update " + filename, nil
	} else if removedLines > addedLines*2 {
		return "refactor: simplify code in multiple files", nil
	} else if addedLines > removedLines*2 {
		return typePrefix + ": implement new functionality", nil
	} else {
		return typePrefix + ": update multiple files", nil
	}
}
