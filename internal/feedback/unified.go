package feedback

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// ProviderConfig contains configuration for different LLM providers
type ProviderConfig struct {
	BaseURL    string
	DefaultModel string
	Name       string
}

// Known providers
var (
	ProviderXAI = ProviderConfig{
		BaseURL:      "https://api.x.ai/v1",
		DefaultModel: "grok-2-1212",
		Name:         "xAI",
	}
	
	ProviderOpenAI = ProviderConfig{
		BaseURL:      "", // Default OpenAI URL
		DefaultModel: "gpt-3.5-turbo",
		Name:         "OpenAI",
	}
	
	ProviderDeepSeek = ProviderConfig{
		BaseURL:      "https://api.deepseek.com/v1", // This is a placeholder, replace with actual URL
		DefaultModel: "deepseek-chat",
		Name:         "DeepSeek",
	}
)

// UnifiedFeedbackEngine generates feedback using any OpenAI-compatible API
type UnifiedFeedbackEngine struct {
	client *openai.Client
	model  string
	provider ProviderConfig
}

// NewUnifiedFeedbackEngine creates a new unified feedback engine
func NewUnifiedFeedbackEngine(provider string, model string, apiKey string) *UnifiedFeedbackEngine {
	var providerConfig ProviderConfig
	
	// Select provider configuration
	switch provider {
	case "xai":
		providerConfig = ProviderXAI
	case "openai":
		providerConfig = ProviderOpenAI
	case "deepseek":
		providerConfig = ProviderDeepSeek
	default:
		// Default to OpenAI if unknown provider
		providerConfig = ProviderOpenAI
	}
	
	// Use provider's default model if none specified
	if model == "" {
		model = providerConfig.DefaultModel
	}
	
	// Configure the client
	config := openai.DefaultConfig(apiKey)
	if providerConfig.BaseURL != "" {
		config.BaseURL = providerConfig.BaseURL
	}
	
	client := openai.NewClientWithConfig(config)
	return &UnifiedFeedbackEngine{
		client:   client,
		model:    model,
		provider: providerConfig,
	}
}

// GenerateFeedback implements the FeedbackEngine interface
func (e *UnifiedFeedbackEngine) GenerateFeedback(ctx CommitContext) (string, error) {
	// Create the system prompt
	systemPrompt := `You are a snarky but insightful Git expert named Moai. 
Given a commit message and time of day, give a short and funny, but helpful comment.
Your responses must be ONE sentence only and should be witty, memorable, and concise.
Responses should be between 50-120 characters.`

	// Format the user prompt with commit information
	timeOfDay := GetTimeOfDay(ctx.Timestamp)
	userPrompt := fmt.Sprintf("Commit message: \"%s\"\nTime of day: %s", 
		ctx.Message, timeOfDay)
	
	// Add diff information if available
	if ctx.Diff != "" {
		userPrompt += fmt.Sprintf("\n\nCommit changes:\n%s", ctx.Diff)
	}

	// Create the chat completion request
	request := openai.ChatCompletionRequest{
		Model: e.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   150,
		N:           1,
	}

	// Send the request to the API
	response, err := e.client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return "", fmt.Errorf("%s API error: %w", e.provider.Name, err)
	}

	// Extract the response content
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from %s API", e.provider.Name)
} 