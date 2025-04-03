package feedback

import (
	"log"
	"strings"
	"time"
)

// CommitContext contains information about a commit
type CommitContext struct {
	Message       string
	Timestamp     time.Time
	Diff          string                 // Optional
	CommitHistory []string               // Recent commit messages
	CommitStats   map[string]interface{} // Stats about recent commits
}

// FeedbackEngine defines the interface for generating commit feedback
type FeedbackEngine interface {
	// Generate feedback based on commit context
	GenerateFeedback(context CommitContext) (string, error)

	// Generate insights for a weekly summary
	GenerateSummaryFeedback(context CommitContext) (string, error)

	// Generate commit message suggestions based on staged changes and history
	GenerateCommitSuggestion(context CommitContext) (string, error)
}

// EngineName returns a string identifier for an engine type
type EngineName string

const (
	// Local feedback engine (no LLM)
	EngineLocal EngineName = "local"
	// xAI feedback engine
	EngineXAI EngineName = "xai"
	// OpenAI feedback engine
	EngineOpenAI EngineName = "openai"
	// DeepSeek feedback engine
	EngineDeepSeek EngineName = "deepseek"
)

// NewFeedbackEngine creates a new feedback engine based on the provided configuration
func NewFeedbackEngine(provider string, model string, apiKey string, personalityName string, personalityFile string) FeedbackEngine {
	// Normalize provider name to lowercase for case-insensitive comparison
	provider = strings.ToLower(strings.TrimSpace(provider))
	
	// Validate provider if API key is provided
	if apiKey != "" {
		validProviders := map[string]bool{
			"xai":      true,
			"openai":   true,
			"deepseek": true,
		}
		
		if !validProviders[provider] {
			// Log warning and default to a known provider
			log.Printf("Warning: Unknown provider '%s', defaulting to 'xai'", provider)
			provider = "xai"
		}
		
		// Ensure we have a valid model name
		if model == "" {
			// Set default model based on provider
			switch provider {
			case "xai":
				model = "grok-2-1212"
			case "openai":
				model = "gpt-3.5-turbo"
			case "deepseek":
				model = "deepseek-chat"
			}
		}
		
		return NewUnifiedFeedbackEngine(provider, model, apiKey, personalityName, personalityFile)
	}

	// Fallback to local feedback engine if no API key is provided
	return NewLocalFeedbackEngine()
}
