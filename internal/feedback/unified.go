package feedback

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
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
	client          *openai.Client
	model           string
	provider        ProviderConfig
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
		client:          client,
		model:           model,
		provider:        providerConfig,
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

// GenerateSummaryFeedback provides insights for a weekly summary or on-demand analysis
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

	// Create a custom system prompt for summaries or on-demand feedback
	systemPrompt := personalityConfig.SystemPrompt
	if strings.Contains(systemPrompt, "one-liner") || strings.Contains(systemPrompt, "one sentence") {
		// Determine if this is a weekly summary or on-demand feedback
		isOnDemand := strings.Contains(ctx.Message, "On-Demand")

		if isOnDemand {
			// Specialized prompt for targeted code analysis with professional-with-sass tone
			systemPrompt = `You are an insightful Git expert who analyzes code practices and commit patterns.
Your task is to provide professional, actionable feedback on the specific set of commits being reviewed.
Include occasional subtle wit or clever observations (about 20% of your response) while maintaining 80% professionalism.
Focus on identifying patterns, potential issues, and specific suggestions for improvement.
Consider best practices related to commit message quality, code organization, and development workflow.
Your response should be 2-4 paragraphs with useful observations and actionable recommendations.
If diffs are provided, focus your analysis on the actual code changes too.`
		} else {
			// Weekly summary prompt with professional-with-sass tone
			systemPrompt = `You are an insightful Git expert who analyzes commit patterns with a professional tone and occasional wit.
Provide a thoughtful, detailed analysis of the commit history that's 80% professional insights and 20% clever observations.
Focus on patterns, trends, and actionable insights while occasionally adding a light touch of sass or humor.
Your response should be 3-5 paragraphs with useful observations and suggestions.
Maintain a helpful and informative tone throughout, with just enough personality to keep it engaging.`
		}
	}

	// Determine the appropriate user prompt based on context
	var userPrompt string
	isOnDemand := strings.Contains(ctx.Message, "On-Demand")

	if isOnDemand {
		// Specialized prompt for on-demand feedback
		userPrompt = fmt.Sprintf(`I'd like you to analyze this specific set of Git commits.

Commit messages:
%s

Commit statistics:
- Total commits: %v
- Files changed: %v
- Lines added: %v
- Lines deleted: %v

%s

Please provide targeted feedback about:
1. Code quality patterns visible in these commits
2. Commit message quality and clarity
3. Specific suggestions for improvement
4. Best practices that could be applied

Focus on giving actionable, specific feedback for these particular commits:`,
			formatCommitList(ctx.CommitHistory),
			ctx.CommitStats["total_commits"],
			ctx.CommitStats["total_files_changed"],
			ctx.CommitStats["total_insertions"],
			ctx.CommitStats["total_deletions"],
			diffContext(ctx.Diff))
	} else {
		// Original weekly summary prompt
		userPrompt = fmt.Sprintf(`I'd like you to analyze my Git commit history from the past week.

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
		Temperature: float32(personalityConfig.Temperature),
		MaxTokens:   800, // Increase token limit for more detailed analysis
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

// GenerateCommitSuggestion creates an AI-generated commit message based on staged changes
func (e *UnifiedFeedbackEngine) GenerateCommitSuggestion(ctx CommitContext) (string, error) {
	// Load personality configuration - only needed for fallback, not for styling
	personalities, err := personality.LoadPersonalities(e.personalityFile)
	if err != nil {
		// Fall back to default personalities if there's an error
		personalities = personality.DefaultPersonalities()
	}

	// Get the selected personality - we'll only use this for basic functionality
	_, err = personalities.GetPersonality(e.personalityName)
	if err != nil {
		// Just make sure we have a valid personality, but won't use its styling
		_, _ = personalities.GetPersonality("")
	}

	// Use a custom system prompt focused on commit message generation
	// This override ensures professional commit messages regardless of personality
	systemPrompt := `You are a professional Git expert who writes clear, precise, and effective commit messages.
Your task is to suggest a commit message that accurately describes the changes.
Follow these guidelines:
1. Use conventional commits format for the subject line: type(scope): description
2. First line must be under 50 characters
3. For SUBSTANTIAL changes (multiple files or significant code changes), ALWAYS add a blank line followed by 2-4 bullet points explaining key changes
4. Use present tense imperative mood (e.g., "fix bug" not "fixes bug")
5. The subject line should focus on the most significant aspect of the change
6. Include scope in parentheses when appropriate: type(scope): description
7. Common types: feat, fix, docs, style, refactor, test, chore
8. Make bullet points start with "- " and be concise but descriptive
9. If changes affect more than 3 files or have >100 line changes, DEFINITELY use a multi-line format
10. Respond with ONLY the commit message, no explanations

For small changes, a single line is sufficient.
For major changes (>100 lines or multiple files), ALWAYS use multi-line format with bullet points.`

	// Prepare the diff context - enhanced with file analysis
	diffContext := `
Here's the current diff of staged changes:

` + ctx.Diff + `

Analysis of changes:
`

	// Add simple analysis of the types of files changed and how many lines were modified
	var totalAdditions, totalDeletions int
	changedFiles := make(map[string]bool)
	
	// Track different types of files
	docFiles := make(map[string]bool)     // Documentation files (.md, .txt, etc)
	codeFiles := make(map[string]bool)    // Source code files (.go, .js, etc)
	configFiles := make(map[string]bool)  // Configuration files (.json, .yaml, etc)
	buildFiles := make(map[string]bool)   // Build files (Makefile, CMakeLists.txt, etc)
	testFiles := make(map[string]bool)    // Test files (*_test.go, etc)
	scriptFiles := make(map[string]bool)  // Scripts (.sh, .bat, etc)
	
	// Track file operations
	addedFiles := make(map[string]bool)
	modifiedFiles := make(map[string]bool)
	deletedFiles := make(map[string]bool)
	
	// File extension categorization maps for better maintainability
	extensionCategories := map[string]map[string]bool{
		"doc": docFiles,
		"code": codeFiles,
		"config": configFiles,
		"build": buildFiles,
		"script": scriptFiles,
	}
	
	// Map extensions to categories
	extensionMap := map[string]string{
		// Documentation files
		".md": "doc", ".txt": "doc", ".rst": "doc", ".adoc": "doc", 
		".markdown": "doc", ".wiki": "doc", ".org": "doc", 
		
		// Source code files
		".go": "code", ".js": "code", ".ts": "code", ".py": "code", 
		".java": "code", ".c": "code", ".cpp": "code", ".cc": "code", 
		".h": "code", ".hpp": "code", ".cs": "code", ".rb": "code", 
		".php": "code", ".swift": "code", ".kt": "code", ".rs": "code",
		
		// Configuration files
		".json": "config", ".yaml": "config", ".yml": "config", ".toml": "config", 
		".ini": "config", ".xml": "config", ".properties": "config", ".conf": "config",
		
		// Build files
		".bazel": "build", ".bzl": "build", ".mk": "build",
		
		// Script files
		".sh": "script", ".bash": "script", ".zsh": "script", 
		".bat": "script", ".cmd": "script", ".ps1": "script",
	}
	
	// Simple diff parser to count lines and identify files
	lines := strings.Split(ctx.Diff, "\n")
	currentFile := ""
	
	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				filePath := strings.TrimPrefix(parts[2], "a/")
				currentFile = filePath
				changedFiles[filePath] = true
				
				// Categorize by file type
				ext := filepath.Ext(filePath)
				baseName := filepath.Base(filePath)
				
				// Special file handling for common non-extension files
				if baseName == "Makefile" || baseName == "Dockerfile" || 
				   baseName == "CMakeLists.txt" || strings.HasPrefix(baseName, "Jenkinsfile") {
					buildFiles[filePath] = true
				} else if strings.Contains(filePath, "_test.") {
					// Test files get special handling
					testFiles[filePath] = true
				} else if category, found := extensionMap[ext]; found {
					// Use the extension map for categorization
					extensionCategories[category][filePath] = true
				}
			}
		} else if strings.HasPrefix(line, "new file mode") {
			addedFiles[currentFile] = true
		} else if strings.HasPrefix(line, "deleted file mode") {
			deletedFiles[currentFile] = true
		} else if !deletedFiles[currentFile] && !addedFiles[currentFile] {
			modifiedFiles[currentFile] = true
		}
		
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			totalAdditions++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			totalDeletions++
		}
	}

	// Remove files from modified if they're added or deleted
	for file := range addedFiles {
		delete(modifiedFiles, file)
	}
	for file := range deletedFiles {
		delete(modifiedFiles, file)
	}

	// Format the analysis
	diffContext += fmt.Sprintf("- Total files changed: %d (%d added, %d modified, %d deleted)\n",
		len(changedFiles), len(addedFiles), len(modifiedFiles), len(deletedFiles))
	diffContext += fmt.Sprintf("- Lines: +%d, -%d\n\n", totalAdditions, totalDeletions)

	// Add file categories analysis
	if len(docFiles) > 0 {
		fileList := make([]string, 0, len(docFiles))
		for file := range docFiles {
			fileList = append(fileList, file)
		}
		diffContext += fmt.Sprintf("Documentation files: %s\n", strings.Join(fileList, ", "))
	}

	if len(codeFiles) > 0 {
		fileList := make([]string, 0, len(codeFiles))
		for file := range codeFiles {
			fileList = append(fileList, file)
		}
		diffContext += fmt.Sprintf("Code files: %s\n", strings.Join(fileList, ", "))
	}

	if len(buildFiles) > 0 {
		fileList := make([]string, 0, len(buildFiles))
		for file := range buildFiles {
			fileList = append(fileList, file)
		}
		diffContext += fmt.Sprintf("Build files: %s\n", strings.Join(fileList, ", "))
	}

	if len(scriptFiles) > 0 {
		fileList := make([]string, 0, len(scriptFiles))
		for file := range scriptFiles {
			fileList = append(fileList, file)
		}
		diffContext += fmt.Sprintf("Script files: %s\n", strings.Join(fileList, ", "))
	}

	if len(configFiles) > 0 {
		fileList := make([]string, 0, len(configFiles))
		for file := range configFiles {
			fileList = append(fileList, file)
		}
		diffContext += fmt.Sprintf("Config files: %s\n", strings.Join(fileList, ", "))
	}

	if len(testFiles) > 0 {
		fileList := make([]string, 0, len(testFiles))
		for file := range testFiles {
			fileList = append(fileList, file)
		}
		diffContext += fmt.Sprintf("Test files: %s\n", strings.Join(fileList, ", "))
	}

	// Add operations analysis
	diffContext += "\nFile operations:\n"

	if len(addedFiles) > 0 {
		fileList := make([]string, 0, len(addedFiles))
		for file := range addedFiles {
			fileList = append(fileList, file)
		}
		diffContext += fmt.Sprintf("Added: %s\n", strings.Join(fileList, ", "))
	}

	if len(modifiedFiles) > 0 {
		fileList := make([]string, 0, len(modifiedFiles))
		for file := range modifiedFiles {
			fileList = append(fileList, file)
		}
		diffContext += fmt.Sprintf("Modified: %s\n", strings.Join(fileList, ", "))
	}

	if len(deletedFiles) > 0 {
		fileList := make([]string, 0, len(deletedFiles))
		for file := range deletedFiles {
			fileList = append(fileList, file)
		}
		diffContext += fmt.Sprintf("Deleted: %s\n", strings.Join(fileList, ", "))
	}

	// Create a user prompt focused on commit message generation with emphasis on changes
	isSubstantialChange := len(changedFiles) > 2 || totalAdditions+totalDeletions > 50
	
	var userPrompt string
	basePrompt := fmt.Sprintf(`I need a%s commit message for these staged changes.

%s

CODE CHANGES DETAIL:
%s

SEMANTIC ANALYSIS:
%s

CODE STRUCTURE ANALYSIS:
%s

Past commit messages for limited context (do not rely heavily on these patterns):
%s`,
		func() string {
			if isSubstantialChange {
				return " multi-line"
			}
			return ""
		}(),
		diffContext,
		formatCodeChanges(ctx.Diff),
		formatSemanticChanges(extractCodeSemantics(ctx.Diff)),
		formatCodeStructure(analyzeCodeStructure(ctx.Diff)),
		formatCommitList(ctx.CommitHistory))
		
	// Add instructions based on change size
	if isSubstantialChange {
		userPrompt = basePrompt + fmt.Sprintf(`

This is a SUBSTANTIAL change affecting %d files with %d insertions and %d deletions.
Therefore, please provide a multi-line commit message with:
1. A clear, concise subject line following conventional commit format (type(scope): description)
2. A blank line
3. 2-4 bullet points that summarize the key components or areas changed

Based primarily on the ACTUAL CODE CHANGES shown above, create a detailed commit message that accurately captures the scope and meaning of these changes:`, 
			len(changedFiles), totalAdditions, totalDeletions)
	} else {
		userPrompt = basePrompt + `

Based primarily on the ACTUAL CODE CHANGES shown above, suggest a BRIEF, CONCISE commit message that accurately describes the most important changes. Focus on being direct and to the point - every word must justify its inclusion:`
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
		Temperature: 0.3, // Slightly higher temperature for more nuanced messages
		MaxTokens:   250, // Increased token limit to accommodate multi-line messages
		N:           1,
	}

	// Send the request to the API
	response, err := e.client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return "", fmt.Errorf("%s API error: %w", e.provider.Name, err)
	}

	// Extract the response content
	if len(response.Choices) > 0 {
		// Get the raw response
		rawSuggestion := response.Choices[0].Message.Content

		// Clean up the response and extract only the actual commit message
		suggestion := extractCommitMessage(rawSuggestion)

		return suggestion, nil
	}

	return "", fmt.Errorf("no response from %s API", e.provider.Name)
}

// extractCommitMessage parses the LLM response to extract just the commit message
func extractCommitMessage(response string) string {
	// Trim whitespace
	response = strings.TrimSpace(response)

	// If wrapped in quotes, remove them
	if strings.HasPrefix(response, "\"") && strings.HasSuffix(response, "\"") {
		response = response[1 : len(response)-1]
	}

	// Check if the response contains a code block with ```
	if strings.Contains(response, "```") {
		// Extract content between code blocks
		parts := strings.Split(response, "```")
		if len(parts) >= 3 {
			// The code block content is in the even indices (1, 3, etc.)
			response = strings.TrimSpace(parts[1])
		}
	}

	// Split into lines to process
	lines := strings.Split(response, "\n")

	// Get first non-empty line
	var firstLine string
	for _, line := range lines {
		if trimmedLine := strings.TrimSpace(line); trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
			firstLine = trimmedLine
			break
		}
	}

	// Ensure first line is under 60 characters
	if len(firstLine) > 60 {
		firstLine = firstLine[:57] + "..."
	}

	// If we have a conventional commit format, ensure it's properly formatted
	if strings.Contains(firstLine, ":") {
		parts := strings.SplitN(firstLine, ":", 2)
		if len(parts) == 2 {
			prefix := parts[0]
			message := strings.TrimSpace(parts[1])
			
			// Known conventional commit types
			commitTypes := []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert"}
			
			// Check if prefix is a valid commit type or a type with scope
			isValidType := false
			for _, cType := range commitTypes {
				if prefix == cType || strings.HasPrefix(prefix, cType+"(") && strings.HasSuffix(prefix, ")") {
					isValidType = true
					break
				}
			}
			
			if isValidType {
				// Proper formatting with no space before colon and one space after
				firstLine = prefix + ": " + message
			}
		}
	}

	// Process body lines - preserve bullet points and maintain proper multi-line format
	var bodyLines []string
	var inBody = false
	
	for i := 1; i < len(lines); i++ {
		trimmedLine := strings.TrimSpace(lines[i])
		
		// Skip empty lines until we reach body content
		if !inBody && trimmedLine == "" {
			continue
		}
		
		// We've now reached body content
		inBody = true
		
		// Skip comment lines and empty lines after we've found body content
		if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
			// Ensure bullet points have proper format
			if strings.HasPrefix(trimmedLine, "* ") {
				trimmedLine = "- " + trimmedLine[2:]
			} else if strings.HasPrefix(trimmedLine, "•") {
				trimmedLine = "- " + trimmedLine[1:]
			} else if strings.HasPrefix(trimmedLine, "-") && !strings.HasPrefix(trimmedLine, "- ") {
				trimmedLine = "- " + trimmedLine[1:]
			}
			
			bodyLines = append(bodyLines, trimmedLine)
		}
	}
	
	// For significant changes, keep full body content with all bullet points
	if len(bodyLines) > 0 {
		// Ensure blank line after subject
		return firstLine + "\n\n" + strings.Join(bodyLines, "\n")
	}
	
	// If no body lines, return just the first line
	return firstLine
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

// diffContext formats diff context for the prompt
func diffContext(diff string) string {
	if diff == "" {
		return ""
	}

	return fmt.Sprintf(`Code changes (diff context):
%s`, diff)
}

// formatCodeChanges formats code changes for the prompt
func formatCodeChanges(diff string) string {
	if diff == "" {
		return ""
	}

	// Split the diff into lines
	lines := strings.Split(diff, "\n")

	// Initialize a result string with a reasonable capacity to reduce allocations
	var result strings.Builder
	result.Grow(len(diff) / 2) // Estimate capacity at half the original diff size
	
	// Keep track of current file
	currentFile := ""
	
	// Capture context lines (unchanged lines around changes)
	const contextLines = 3
	lineBuffer := make([]string, 0, contextLines*2+1)
	inChangeBlock := false
	
	// Track if we're in a meaningful code section or just metadata
	inCodeSection := false
	
	// Iterate over each line
	for i, line := range lines {
		// Check if this is a new file
		if strings.HasPrefix(line, "diff --git") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				// Add a separator between files if we've already seen a file
				if currentFile != "" {
					result.WriteString("\n-----------------------------------\n\n")
				}
				
				// Extract the file name
				filePath := strings.TrimPrefix(parts[2], "a/")
				currentFile = filePath
				result.WriteString(fmt.Sprintf("==== CHANGES IN FILE: %s ====\n", filePath))
				inCodeSection = false
			}
			continue
		}
		
		// Skip git metadata lines
		if strings.HasPrefix(line, "index ") || 
		   strings.HasPrefix(line, "+++") || 
		   strings.HasPrefix(line, "---") {
			continue
		}
		
		// Check for chunk header
		if strings.HasPrefix(line, "@@") {
			inCodeSection = true
			continue
		}
		
		// Process only if we're in a code section
		if !inCodeSection {
			continue
		}
		
		// Handle code lines
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
			// If we're starting a new change block, add context before
			if !inChangeBlock {
				inChangeBlock = true
				
				// Add preceding context lines
				for j := i - contextLines; j < i; j++ {
					if j >= 0 && !strings.HasPrefix(lines[j], "diff --git") && 
					   !strings.HasPrefix(lines[j], "index ") && 
					   !strings.HasPrefix(lines[j], "+++") && 
					   !strings.HasPrefix(lines[j], "---") &&
					   !strings.HasPrefix(lines[j], "@@") {
						result.WriteString(fmt.Sprintf("  %s\n", lines[j]))
					}
				}
			}
			
			// Add the changed line with highlighting for better readability
			if strings.HasPrefix(line, "+") {
				result.WriteString(fmt.Sprintf("%s\n", line))
			} else {
				result.WriteString(fmt.Sprintf("%s\n", line))
			}
			
			// Clear the line buffer
			lineBuffer = lineBuffer[:0]
		} else {
			// Unchanged line
			if inChangeBlock {
				// Add the unchanged line
				lineBuffer = append(lineBuffer, line)
				
				// Check if we have enough context lines or reached the end
				if len(lineBuffer) >= contextLines || i == len(lines)-1 {
					// Write the context lines
					for _, bufLine := range lineBuffer {
						result.WriteString(fmt.Sprintf("  %s\n", bufLine))
					}
					
					// Reset the buffer and change block flag
					lineBuffer = lineBuffer[:0]
					inChangeBlock = false
					
					// Add a separator if not at the end
					if i < len(lines)-1 {
						result.WriteString("\n")
					}
				}
			}
		}
	}

	return result.String()
}

// extractCodeSemantics analyzes the diff to identify key semantic changes
// for better commit message suggestions
func extractCodeSemantics(diff string) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Track function additions/modifications/removals
	functionChanges := make(map[string]string)
	
	// Track package/import changes
	importChanges := make([]string, 0)
	
	// Track meaningful variable declarations
	variableChanges := make(map[string]string)
	
	// Split the diff into lines
	lines := strings.Split(diff, "\n")
	
	// State tracking
	currentFile := ""
	inImportBlock := false
	
	// Regex patterns for semantic analysis
	functionPattern := regexp.MustCompile(`^[+-](func\s+\w+)`)
	methodPattern := regexp.MustCompile(`^[+-](func\s+\([^)]+\)\s+\w+)`)
	importPattern := regexp.MustCompile(`^[+-]\s*import\s+(?:\w+\s+)?"([^"]+)"`)
	variablePattern := regexp.MustCompile(`^[+-]\s*(\w+)\s*:?=\s*(.+)$`)
	
	for _, line := range lines {
		// Track current file
		if strings.HasPrefix(line, "diff --git") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				filePath := strings.TrimPrefix(parts[2], "a/")
				currentFile = filePath
			}
			inImportBlock = false
			continue
		}
		
		// Skip metadata lines
		if strings.HasPrefix(line, "index ") || 
		   strings.HasPrefix(line, "+++") || 
		   strings.HasPrefix(line, "---") ||
		   strings.HasPrefix(line, "@@") {
			continue
		}
		
		// Detect import block
		if strings.Contains(line, "import (") {
			inImportBlock = true
		} else if inImportBlock && strings.Contains(line, ")") {
			inImportBlock = false
		}
		
		// Check for function changes
		if matches := functionPattern.FindStringSubmatch(line); len(matches) > 1 {
			op := string(line[0])
			funcDecl := matches[1]
			functionChanges[funcDecl] = op
		}
		
		// Check for method changes
		if matches := methodPattern.FindStringSubmatch(line); len(matches) > 1 {
			op := string(line[0])
			methodDecl := matches[1]
			functionChanges[methodDecl] = op
		}
		
		// Check for import changes
		if matches := importPattern.FindStringSubmatch(line); len(matches) > 1 {
			importChanges = append(importChanges, matches[1])
		} else if inImportBlock && (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-")) {
			// For multi-line import blocks
			importLine := strings.TrimLeft(strings.TrimSpace(line[1:]), "\"")
			importLine = strings.TrimRight(importLine, "\"")
			if importLine != "" && !strings.HasPrefix(importLine, "//") {
				importChanges = append(importChanges, importLine)
			}
		}
		
		// Check for variable changes
		if matches := variablePattern.FindStringSubmatch(line); len(matches) > 1 {
			varName := matches[1]
			varValue := matches[2]
			variableChanges[varName] = varValue
		}
	}
	
	// Store the collected changes
	result["files"] = []string{currentFile}
	result["functions"] = functionChanges
	result["imports"] = importChanges
	result["variables"] = variableChanges
	
	return result
}

// formatSemanticChanges formats semantic changes for the prompt
func formatSemanticChanges(semantics map[string]interface{}) string {
	var result strings.Builder
	
	// Format files
	if files, ok := semantics["files"].([]string); ok && len(files) > 0 {
		result.WriteString("Modified files:\n")
		for _, file := range files {
			result.WriteString(fmt.Sprintf("- %s\n", file))
		}
		result.WriteString("\n")
	}
	
	// Format function changes
	if functions, ok := semantics["functions"].(map[string]string); ok && len(functions) > 0 {
		result.WriteString("Function changes:\n")
		for funcName, op := range functions {
			if op == "+" {
				result.WriteString(fmt.Sprintf("- Added: %s\n", funcName))
			} else if op == "-" {
				result.WriteString(fmt.Sprintf("- Removed: %s\n", funcName))
			} else {
				result.WriteString(fmt.Sprintf("- Modified: %s\n", funcName))
			}
		}
		result.WriteString("\n")
	}
	
	// Format import changes
	if imports, ok := semantics["imports"].([]string); ok && len(imports) > 0 {
		result.WriteString("Import changes:\n")
		for _, imp := range imports {
			result.WriteString(fmt.Sprintf("- %s\n", imp))
		}
		result.WriteString("\n")
	}
	
	// Format variable changes
	if variables, ok := semantics["variables"].(map[string]string); ok && len(variables) > 0 {
		result.WriteString("Variable changes:\n")
		for varName, varValue := range variables {
			result.WriteString(fmt.Sprintf("- %s = %s\n", varName, varValue))
		}
	}
	
	return result.String()
}

// analyzeCodeStructure performs deeper analysis of code structure in the diff
// to identify structural changes like interface implementations, struct modifications, etc.
// This function scans the diff for type definitions, interfaces, structs, and constants
// to provide more semantic understanding of the code changes.
func analyzeCodeStructure(diff string) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Track structural elements
	interfaces := make(map[string][]string)       // Interface -> list of methods
	structs := make(map[string][]string)          // Struct -> list of fields
	typeDeclarations := make(map[string]string)   // Type name -> definition
	constants := make(map[string]string)          // Const name -> value
	
	// Track current file and module
	currentFile := ""
	currentPackage := ""
	
	// State tracking
	inTypeBlock := false
	inConstBlock := false
	currentType := ""
	
	// Regex patterns for semantic analysis
	typeDefPattern := regexp.MustCompile(`^[+-]type\s+(\w+)\s+(struct|interface)`)
	structFieldPattern := regexp.MustCompile(`^[+-]\s*(\w+)(\s+\w+\s*(?:\`+"`"+`[^`+"`"+`]*\`+"`"+`)?)`)
	interfaceMethodPattern := regexp.MustCompile(`^[+-]\s*(\w+\([^)]*\))`)
	constPattern := regexp.MustCompile(`^[+-]const\s+(\w+)\s+=\s+(.*)`)
	packagePattern := regexp.MustCompile(`^[+-]package\s+(\w+)`)
	
	// Split the diff into lines
	lines := strings.Split(diff, "\n")
	
	for _, line := range lines {
		// Track current file
		if strings.HasPrefix(line, "diff --git") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				filePath := strings.TrimPrefix(parts[2], "a/")
				currentFile = filePath
				// Reset state for new file
				inTypeBlock = false
				inConstBlock = false
				currentType = ""
			}
			continue
		}
		
		// Skip metadata lines
		if strings.HasPrefix(line, "index ") || 
		   strings.HasPrefix(line, "+++") || 
		   strings.HasPrefix(line, "---") ||
		   strings.HasPrefix(line, "@@") {
			continue
		}
		
		// Only process added or removed lines
		if !strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "-") {
			continue
		}
		
		// Extract the actual code line without the +/-
		codeLine := line[1:]
		
		// Check for package declaration
		if matches := packagePattern.FindStringSubmatch(line); len(matches) > 1 {
			currentPackage = matches[1]
			continue
		}
		
		// Check for type declarations
		if matches := typeDefPattern.FindStringSubmatch(line); len(matches) > 1 {
			typeName := matches[1]
			typeKind := matches[2]
			currentType = typeName
			
			if typeKind == "struct" {
				inTypeBlock = true
				inConstBlock = false
				structs[typeName] = []string{}
			} else if typeKind == "interface" {
				inTypeBlock = true
				inConstBlock = false
				interfaces[typeName] = []string{}
			}
			
			continue
		}
		
		// Check for const blocks
		if strings.HasPrefix(codeLine, "const (") {
			inConstBlock = true
			inTypeBlock = false
			continue
		} else if strings.HasPrefix(codeLine, ")") {
			inConstBlock = false
			inTypeBlock = false
			continue
		}
		
		// Process content inside blocks
		if inTypeBlock {
			if currentType != "" {
				if _, ok := structs[currentType]; ok {
					// This is a struct field
					if matches := structFieldPattern.FindStringSubmatch(line); len(matches) > 1 {
						fieldName := matches[1]
						structs[currentType] = append(structs[currentType], fieldName)
					}
				} else if _, ok := interfaces[currentType]; ok {
					// This is an interface method
					if matches := interfaceMethodPattern.FindStringSubmatch(line); len(matches) > 1 {
						methodSig := matches[1]
						interfaces[currentType] = append(interfaces[currentType], methodSig)
					}
				}
			}
		} else if inConstBlock {
			// Check for constants
			if matches := constPattern.FindStringSubmatch(line); len(matches) > 1 {
				constName := matches[1]
				constValue := matches[2]
				constants[constName] = constValue
			}
		} else {
			// Check for standalone constants
			if matches := constPattern.FindStringSubmatch(line); len(matches) > 1 {
				constName := matches[1]
				constValue := matches[2]
				constants[constName] = constValue
			}
		}
	}
	
	// Store the collected changes
	result["file"] = currentFile
	result["package"] = currentPackage
	result["interfaces"] = interfaces
	result["structs"] = structs
	result["types"] = typeDeclarations
	result["constants"] = constants
	
	return result
}

// formatCodeStructure formats code structure analysis for the prompt
func formatCodeStructure(structure map[string]interface{}) string {
	var result strings.Builder
	
	// Format interfaces
	if interfaces, ok := structure["interfaces"].(map[string][]string); ok && len(interfaces) > 0 {
		result.WriteString("Modified interfaces:\n")
		for interfaceName, methods := range interfaces {
			result.WriteString(fmt.Sprintf("- %s\n", interfaceName))
			for _, method := range methods {
				result.WriteString(fmt.Sprintf("  - %s\n", method))
			}
		}
		result.WriteString("\n")
	}
	
	// Format structs
	if structs, ok := structure["structs"].(map[string][]string); ok && len(structs) > 0 {
		result.WriteString("Modified structs:\n")
		for structName, fields := range structs {
			result.WriteString(fmt.Sprintf("- %s\n", structName))
			for _, field := range fields {
				result.WriteString(fmt.Sprintf("  - %s\n", field))
			}
		}
		result.WriteString("\n")
	}
	
	// Format type declarations
	if types, ok := structure["types"].(map[string]string); ok && len(types) > 0 {
		result.WriteString("Modified type declarations:\n")
		for typeName, typeDef := range types {
			result.WriteString(fmt.Sprintf("- %s\n", typeName))
			result.WriteString(fmt.Sprintf("  - Definition: %s\n", typeDef))
		}
		result.WriteString("\n")
	}
	
	// Format constants
	if constants, ok := structure["constants"].(map[string]string); ok && len(constants) > 0 {
		result.WriteString("Modified constants:\n")
		for constName, constValue := range constants {
			result.WriteString(fmt.Sprintf("- %s = %s\n", constName, constValue))
		}
		result.WriteString("\n")
	}
	
	return result.String()
}