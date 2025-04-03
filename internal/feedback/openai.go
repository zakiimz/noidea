package feedback

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIFeedbackEngine generates feedback using the OpenAI API
type OpenAIFeedbackEngine struct {
	client *openai.Client
	model  string
}

// NewOpenAIFeedbackEngine creates a new OpenAI feedback engine
func NewOpenAIFeedbackEngine(apiKey string, model string) *OpenAIFeedbackEngine {
	// Set default model if not provided
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	client := openai.NewClient(apiKey)
	return &OpenAIFeedbackEngine{
		client: client,
		model:  model,
	}
}

// GenerateFeedback implements the FeedbackEngine interface
func (e *OpenAIFeedbackEngine) GenerateFeedback(ctx CommitContext) (string, error) {
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
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	// Extract the response content
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from OpenAI API")
} 