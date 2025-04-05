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

	// Track function and method additions/changes
	var functionChanges []string

	// Parse the diff to extract semantic information

	for i, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			// Just detect the file boundary, no need to track the current file
		} else if strings.HasPrefix(line, "+++ b/") {
			file := strings.TrimPrefix(line, "+++ b/")
			filesChanged = append(filesChanged, file)

			// Count file extensions
			ext := filepath.Ext(file)
			if ext != "" {
				fileExtensions[ext]++
			}
		} else if strings.HasPrefix(line, "+func ") && !strings.HasPrefix(line, "+++") {
			// Capture added function declarations
			funcDecl := strings.TrimPrefix(line, "+func ")
			funcName := strings.Split(funcDecl, "(")[0]
			funcName = strings.TrimSpace(funcName)
			if funcName != "" {
				functionChanges = append(functionChanges, funcName)
			}
		} else if i > 0 && strings.HasPrefix(line, "+") && strings.Contains(line, " struct {") {
			// Capture added struct declarations
			prevLine := lines[i-1]
			if strings.HasPrefix(prevLine, "+type ") {
				typeLine := strings.TrimPrefix(prevLine, "+type ")
				structName := strings.Split(typeLine, " ")[0]
				if structName != "" {
					functionChanges = append(functionChanges, "struct "+structName)
				}
			}
		} else if i > 0 && strings.HasPrefix(line, "+") && strings.Contains(line, " interface {") {
			// Capture added interface declarations
			prevLine := lines[i-1]
			if strings.HasPrefix(prevLine, "+type ") {
				typeLine := strings.TrimPrefix(prevLine, "+type ")
				interfaceName := strings.Split(typeLine, " ")[0]
				if interfaceName != "" {
					functionChanges = append(functionChanges, "interface "+interfaceName)
				}
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

	// Generate more specific message based on the detected changes
	if len(functionChanges) > 0 {
		if len(functionChanges) == 1 {
			// Single function change
			return typePrefix + ": add " + functionChanges[0] + " function", nil
		} else if len(functionChanges) <= 3 {
			// Multiple function changes, but not too many to list
			return typePrefix + ": add " + strings.Join(functionChanges[:len(functionChanges)-1], ", ") +
				" and " + functionChanges[len(functionChanges)-1] + " functions", nil
		} else {
			// Too many function changes to list individually
			return typePrefix + ": add multiple new functions including " + functionChanges[0], nil
		}
	} else if len(filesChanged) == 0 {
		return typePrefix + ": update code structure", nil
	} else if len(filesChanged) == 1 {
		filename := filepath.Base(filesChanged[0])
		if addedLines > 50 && removedLines < 10 {
			return typePrefix + ": add new functionality to " + filename, nil
		} else if removedLines > 50 && addedLines < 10 {
			return typePrefix + ": remove unused code from " + filename, nil
		} else if addedLines > 0 && removedLines > 0 {
			return typePrefix + ": refactor code in " + filename, nil
		} else {
			return typePrefix + ": make changes to " + filename, nil
		}
	} else if removedLines > addedLines*2 {
		return "refactor: simplify code across multiple files", nil
	} else if addedLines > removedLines*2 {
		if addedLines > 100 {
			return typePrefix + ": implement major new functionality", nil
		} else {
			return typePrefix + ": add new features", nil
		}
	} else {
		return typePrefix + ": update implementation in multiple files", nil
	}
}
