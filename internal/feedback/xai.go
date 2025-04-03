package feedback

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// XAIFeedbackEngine generates feedback using the X AI API
type XAIFeedbackEngine struct {
	client *openai.Client
	model  string
}

// NewXAIFeedbackEngine creates a new xAI feedback engine
func NewXAIFeedbackEngine(apiKey string, model string) *XAIFeedbackEngine {
	// Set default model if not provided
	if model == "" {
		model = "grok-2-1212"
	}

	// Configure the client with xAI API endpoint
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.x.ai/v1"

	client := openai.NewClientWithConfig(config)
	return &XAIFeedbackEngine{
		client: client,
		model:  model,
	}
}

// GenerateFeedback implements the FeedbackEngine interface
func (e *XAIFeedbackEngine) GenerateFeedback(ctx CommitContext) (string, error) {
	// Create the system prompt
	systemPrompt := `You are a snarky but insightful Git expert named Moai. 
Given a commit message and time of day, give a short and funny, but helpful comment.
Your responses must be ONE sentence only and should be witty, memorable, and concise.
Responses should be between 50-120 characters.`

	// Format the user prompt with commit information
	timeOfDay := getTimeOfDay(ctx.Timestamp)
	userPrompt := fmt.Sprintf("Commit message: \"%s\"\nTime of day: %s", 
		ctx.Message, timeOfDay)

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
		return "", fmt.Errorf("xAI API error: %w", err)
	}

	// Extract the response content
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from xAI API")
}

// getTimeOfDay returns a string representation of the time of day
func getTimeOfDay(t time.Time) string {
	hour := t.Hour()
	
	switch {
	case hour >= 5 && hour < 12:
		return "morning"
	case hour >= 12 && hour < 17:
		return "afternoon"
	case hour >= 17 && hour < 21:
		return "evening"
	default:
		return "night"
	}
} 