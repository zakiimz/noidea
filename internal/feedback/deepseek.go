package feedback

import (
	"fmt"
	"time"
)

// DeepSeekFeedbackEngine generates feedback using the DeepSeek API
type DeepSeekFeedbackEngine struct {
	apiKey string
	model  string
}

// NewDeepSeekFeedbackEngine creates a new DeepSeek feedback engine
func NewDeepSeekFeedbackEngine(apiKey string, model string) *DeepSeekFeedbackEngine {
	// Set default model if not provided
	if model == "" {
		model = "deepseek-chat"
	}

	return &DeepSeekFeedbackEngine{
		apiKey: apiKey,
		model:  model,
	}
}

// GenerateFeedback implements the FeedbackEngine interface
func (e *DeepSeekFeedbackEngine) GenerateFeedback(ctx CommitContext) (string, error) {
	// NOTE: This is a placeholder implementation
	// TODO: Implement actual DeepSeek API integration when/if needed
	
	// For now, return a message explaining the DeepSeek API is not yet implemented
	return fmt.Sprintf(
		"DeepSeek LLM integration not yet implemented (commit: '%s', time: %s)", 
		ctx.Message, 
		ctx.Timestamp.Format(time.RFC3339),
	), nil
} 