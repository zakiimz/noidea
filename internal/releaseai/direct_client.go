// Package releaseai provides AI-powered release notes generation
package releaseai

import (
	"context"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// DirectLLMClient provides direct access to LLM APIs for release notes generation
// Completely separate from the feedback system to avoid pattern interference
type DirectLLMClient struct {
	client       *openai.Client
	model        string
	apiKey       string
	provider     string
	maxTokens    int
	temperature  float32
	systemPrompt string
}

// NewDirectLLMClient creates a new LLM client with direct API access
func NewDirectLLMClient(provider, model, apiKey string, temperature float64) *DirectLLMClient {
	// Configure client based on provider
	baseURL := ""
	switch provider {
	case "xai":
		baseURL = "https://api.x.ai/v1"
	case "openai":
		// Default OpenAI URL is used automatically
	case "deepseek":
		baseURL = "https://api.deepseek.com/v1"
	}

	// Create configuration
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	// Set max tokens based on model
	maxTokens := 2000
	if strings.Contains(model, "gpt-4") || strings.Contains(model, "grok-1") {
		maxTokens = 4000
	}

	return &DirectLLMClient{
		client:      openai.NewClientWithConfig(config),
		model:       model,
		apiKey:      apiKey,
		provider:    provider,
		maxTokens:   maxTokens,
		temperature: float32(temperature),
	}
}

// SetSystemPrompt sets a custom system prompt
func (c *DirectLLMClient) SetSystemPrompt(prompt string) {
	c.systemPrompt = prompt
}

// GenerateReleaseNotes generates release notes directly using the LLM API
func (c *DirectLLMClient) GenerateReleaseNotes(
	prompt string,
	maxAttempts int,
) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var generationErr error
	var response string

	// Try multiple attempts if needed, with increasing temperatures
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Adjust temperature slightly for each retry
		temperature := c.temperature
		if attempt > 0 {
			temperature += float32(attempt) * 0.1
			if temperature > 1.0 {
				temperature = 1.0
			}
		}

		// Use custom system prompt if set, otherwise use default
		systemContent := "You are a professional release notes writer. Generate detailed, accurate release notes for software updates. Your task is to describe changes, features, and fixes. IMPORTANT: Do not analyze commit message patterns or formatting - focus only on the actual software changes. Always begin with a clear overview and organize changes into relevant sections."
		if c.systemPrompt != "" {
			systemContent = c.systemPrompt
		}

		messages := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemContent,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		}

		// Create the request
		req := openai.ChatCompletionRequest{
			Model:       c.model,
			Messages:    messages,
			Temperature: temperature,
			MaxTokens:   c.maxTokens,
		}

		// Send the request
		resp, err := c.client.CreateChatCompletion(ctx, req)
		if err != nil {
			generationErr = err
			continue // Try again if there's an error
		}

		// Check if we got a valid response
		if len(resp.Choices) > 0 {
			response = resp.Choices[0].Message.Content
			if strings.TrimSpace(response) != "" {
				return response, nil
			}
		}
	}

	if generationErr != nil {
		return "", fmt.Errorf("failed to generate release notes after %d attempts: %w", maxAttempts, generationErr)
	}

	return "", fmt.Errorf("failed to generate meaningful release notes after %d attempts", maxAttempts)
}
