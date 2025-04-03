package feedback

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/AccursedGalaxy/noidea/internal/personality"
	openai "github.com/sashabaranov/go-openai"
)

// ProviderConfig contains configuration for different LLM providers
type ProviderConfig struct {
	BaseURL      string
	DefaultModel string
	Name         string
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
	client        *openai.Client
	model         string
	provider      ProviderConfig
	personalityName string
	personalityFile string
}

// NewUnifiedFeedbackEngine creates a new unified feedback engine
func NewUnifiedFeedbackEngine(provider string, model string, apiKey string, personalityName string, personalityFile string) *UnifiedFeedbackEngine {
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
		client:        client,
		model:         model,
		provider:      providerConfig,
		personalityName: personalityName,
		personalityFile: personalityFile,
	}
}

// GenerateFeedback implements the FeedbackEngine interface
func (e *UnifiedFeedbackEngine) GenerateFeedback(ctx CommitContext) (string, error) {
	// Load personality configuration
	personalities, err := personality.LoadPersonalities(e.personalityFile)
	if err != nil {
		// Fall back to default personalities if there's an error
		personalities = personality.DefaultPersonalities()
	}

	// Get the selected personality
	personalityConfig, err := personalities.GetPersonality(e.personalityName)
	if err != nil {
		// Fall back to default personality
		personalityConfig, _ = personalities.GetPersonality("")
	}

	// Create personality context for template rendering
	personalityCtx := personality.Context{
		Message:       ctx.Message,
		TimeOfDay:     GetTimeOfDay(ctx.Timestamp),
		Diff:          ctx.Diff,
		Username:      getUserName(),
		RepoName:      getRepoName(),
		CommitHistory: ctx.CommitHistory,
		CommitStats:   ctx.CommitStats,
	}

	// Generate the prompt using the personality template
	userPrompt, err := personalityConfig.GeneratePrompt(personalityCtx)
	if err != nil {
		return "", fmt.Errorf("failed to generate prompt: %w", err)
	}

	// Create the chat completion request
	request := openai.ChatCompletionRequest{
		Model: e.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: personalityConfig.SystemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Temperature: float32(personalityConfig.Temperature),
		MaxTokens:   personalityConfig.MaxTokens,
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

// GenerateSummaryFeedback provides insights for a weekly summary
func (e *UnifiedFeedbackEngine) GenerateSummaryFeedback(ctx CommitContext) (string, error) {
	// Load personality configuration
	personalities, err := personality.LoadPersonalities(e.personalityFile)
	if err != nil {
		// Fall back to default personalities if there's an error
		personalities = personality.DefaultPersonalities()
	}

	// Get the selected personality
	personalityConfig, err := personalities.GetPersonality(e.personalityName)
	if err != nil {
		// Fall back to default personality
		personalityConfig, _ = personalities.GetPersonality("")
	}

	// Create a custom system prompt for summaries
	systemPrompt := personalityConfig.SystemPrompt
	if strings.Contains(systemPrompt, "one-liner") || strings.Contains(systemPrompt, "one sentence") {
		// Replace one-liner instruction with more detailed analysis for summaries
		systemPrompt = `You are an insightful Git expert who analyzes commit patterns.
Provide a thoughtful, detailed analysis of the commit history.
Focus on patterns, trends, and actionable insights.
Your response should be 3-5 paragraphs with useful observations and suggestions.`
	}

	// Create user prompt with commit history
	userPrompt := fmt.Sprintf(`I'd like you to analyze my Git commit history from the past week.

Commit messages:
%s

Commit statistics:
- Total commits: %v
- Unique authors: %v
- Files changed: %v
- Lines added: %v
- Lines deleted: %v

Please provide insights about:
1. Commit message patterns and quality
2. Work focus areas (based on commit messages)
3. Time distribution patterns
4. Suggestions for improving workflow or commit habits

Respond with thoughtful analysis and actionable suggestions:`,
		formatCommitList(ctx.CommitHistory),
		ctx.CommitStats["total_commits"],
		ctx.CommitStats["unique_authors"],
		ctx.CommitStats["total_files_changed"],
		ctx.CommitStats["total_insertions"],
		ctx.CommitStats["total_deletions"])

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
		Temperature: 0.7, // Slightly higher temperature for creative insights
		MaxTokens:   800, // Longer response for summary
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

// formatCommitList creates a formatted string of commit messages
func formatCommitList(commits []string) string {
	var result strings.Builder
	
	for i, commit := range commits {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, commit))
	}
	
	return result.String()
}

// getUserName attempts to get the Git user name
func getUserName() string {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "User"
	}
	return strings.TrimSpace(string(output))
}

// getRepoName attempts to get the Git repository name
func getRepoName() string {
	// Try to get the remote origin URL
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "repository"
	}

	// Extract repo name from URL
	repoURL := strings.TrimSpace(string(output))
	parts := strings.Split(repoURL, "/")
	if len(parts) > 0 {
		repoName := parts[len(parts)-1]
		// Remove .git suffix if present
		repoName = strings.TrimSuffix(repoName, ".git")
		return repoName
	}

	return "repository"
} 