package feedback

import (
	"math/rand"

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