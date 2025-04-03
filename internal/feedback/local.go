package feedback

import (
	"github.com/AccursedGalaxy/noidea/internal/moai"
)

// LocalFeedbackEngine generates feedback without using an LLM
type LocalFeedbackEngine struct{}

// NewLocalFeedbackEngine creates a new local feedback engine
func NewLocalFeedbackEngine() *LocalFeedbackEngine {
	return &LocalFeedbackEngine{}
}

// GenerateFeedback implements the FeedbackEngine interface
func (e *LocalFeedbackEngine) GenerateFeedback(context CommitContext) (string, error) {
	// Use our existing feedback generator
	return moai.GetRandomFeedback(context.Message), nil
} 