package personality

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
)

// Context represents the commit context for template rendering
type Context struct {
	Message       string
	TimeOfDay     string
	Diff          string
	Username      string
	RepoName      string
	CommitHistory []string               // Recent commit messages
	CommitStats   map[string]interface{} // Stats about recent commits
}

// Personality defines a configurable AI personality
type Personality struct {
	Name             string  `toml:"name"`
	Description      string  `toml:"description"`
	SystemPrompt     string  `toml:"system_prompt"`
	UserPromptFormat string  `toml:"user_prompt_format"`
	MaxTokens        int     `toml:"max_tokens"`
	Temperature      float64 `toml:"temperature"`
}

// PersonalityConfig holds multiple personality configurations
type PersonalityConfig struct {
	Default       string                 `toml:"default"`
	Personalities map[string]Personality `toml:"personalities"`
}

// DefaultPersonalities returns the built-in personality configurations
func DefaultPersonalities() PersonalityConfig {
	return PersonalityConfig{
		Default: "professional_sass",
		Personalities: map[string]Personality{
			"professional_sass": {
				Name:        "Professional with Sass",
				Description: "A professional Git expert with a subtle hint of sass",
				SystemPrompt: `You are a professional Git expert named Moai with subtle hints of wit and sass.
Your responses should be primarily informative and useful while occasionally delivering a clever observation.
Focus on providing actionable insights about the commit with 80% professionalism and 20% subtle humor.
Keep your responses concise (one sentence) and to the point.`,
				UserPromptFormat: `Commit message: "{{.Message}}"
Time of day: {{.TimeOfDay}}
{{if .Diff}}Commit diff summary: {{.Diff}}{{end}}
{{if .CommitHistory}}
Recent commit messages:
{{range .CommitHistory}}- "{{.}}"
{{end}}{{end}}

Provide professional feedback with a subtle touch of wit about this commit:`,
				MaxTokens:   150,
				Temperature: 0.6,
			},
			"snarky_reviewer": {
				Name:        "Snarky Code Reviewer",
				Description: "A code reviewer with a sarcastic and witty attitude",
				SystemPrompt: `You are a snarky but insightful Git expert named Moai. 
Your responses should be witty, memorable, and concise.
Always aim to be funny while also providing insight about the commit.
Keep your responses between 50-120 characters and as a single sentence.`,
				UserPromptFormat: `Commit message: "{{.Message}}"
Time of day: {{.TimeOfDay}}
{{if .Diff}}Commit diff summary: {{.Diff}}{{end}}
{{if .CommitHistory}}
Recent commit messages:
{{range .CommitHistory}}- "{{.}}"
{{end}}{{end}}

Provide a snarky, funny one-liner about this commit:`,
				MaxTokens:   150,
				Temperature: 0.7,
			},
			"supportive_mentor": {
				Name:        "Supportive Mentor",
				Description: "A supportive and encouraging mentor",
				SystemPrompt: `You are a supportive and encouraging Git mentor.
Your responses should be positive, helpful, and motivating.
You want to help the developer feel good about their progress while subtly suggesting improvements.
Keep your responses concise (one sentence) and encouraging.`,
				UserPromptFormat: `Commit message: "{{.Message}}"
Time of day: {{.TimeOfDay}}
{{if .Diff}}Commit diff summary: {{.Diff}}{{end}}
{{if .CommitHistory}}
Recent commit messages:
{{range .CommitHistory}}- "{{.}}"
{{end}}{{end}}

Provide a supportive, encouraging comment about this commit:`,
				MaxTokens:   150,
				Temperature: 0.6,
			},
			"git_expert": {
				Name:        "Git Expert",
				Description: "A professional Git expert providing technical feedback",
				SystemPrompt: `You are a professional Git expert with deep knowledge of best practices.
Your responses should be technical, insightful, and focused on Git best practices.
Provide specific technical advice to improve the commit or commend good practices you notice.
Keep your responses concise (one sentence) and informative.`,
				UserPromptFormat: `Commit message: "{{.Message}}"
Time of day: {{.TimeOfDay}}
{{if .Diff}}Commit diff summary: {{.Diff}}{{end}}
{{if .CommitHistory}}
Recent commit messages:
{{range .CommitHistory}}- "{{.}}"
{{end}}{{end}}
{{if .CommitStats}}
Commit patterns:
- Recent commits: {{index .CommitStats "total_commits"}}
- Common commit times: {{if index .CommitStats "commits_by_hour"}}{{index (index .CommitStats "commits_by_hour") (printf "%d" (time "15:04" .TimeOfDay | hour))}} commits at this hour{{end}}
{{end}}

Provide concise, technical Git feedback about this commit:`,
				MaxTokens:   150,
				Temperature: 0.4,
			},
		},
	}
}

// LoadPersonalities loads personality configurations from the given path
func LoadPersonalities(path string) (PersonalityConfig, error) {
	// Start with default personalities
	config := DefaultPersonalities()

	// If no path provided, return defaults
	if path == "" {
		return config, nil
	}

	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		return config, fmt.Errorf("personality file not found: %s", path)
	}

	// Load and parse TOML file
	var fileConfig PersonalityConfig
	_, err := toml.DecodeFile(path, &fileConfig)
	if err != nil {
		return config, fmt.Errorf("failed to decode personality file: %w", err)
	}

	// Merge with defaults - any custom personalities override defaults
	for name, personality := range fileConfig.Personalities {
		config.Personalities[name] = personality
	}

	// Override default if specified
	if fileConfig.Default != "" {
		// Check if the specified default exists
		if _, exists := config.Personalities[fileConfig.Default]; exists {
			config.Default = fileConfig.Default
		}
	}

	return config, nil
}

// GetPersonality returns a personality by name, falling back to default if not found
func (pc PersonalityConfig) GetPersonality(name string) (Personality, error) {
	// If name is empty, use default
	if name == "" {
		name = pc.Default
	}

	personality, exists := pc.Personalities[name]
	if !exists {
		return Personality{}, fmt.Errorf("personality not found: %s", name)
	}

	return personality, nil
}

// FindPersonalityFile returns the path to the personality configuration file
func FindPersonalityFile() string {
	// Check in current directory
	if _, err := os.Stat(".noidea-personalities.toml"); err == nil {
		return ".noidea-personalities.toml"
	}

	// Check in home directory
	home, err := os.UserHomeDir()
	if err == nil {
		homeConfig := filepath.Join(home, ".noidea", "personalities.toml")
		if _, err := os.Stat(homeConfig); err == nil {
			return homeConfig
		}
	}

	// No file found
	return ""
}

// GeneratePrompt generates a formatted prompt from a personality and context
func (p Personality) GeneratePrompt(ctx Context) (string, error) {
	// Define template functions
	funcMap := template.FuncMap{
		"hour": func(timeStr string) int {
			t, err := time.Parse("15:04", timeStr)
			if err != nil {
				return 0
			}
			return t.Hour()
		},
		"time": func(format string, timeStr string) int {
			t, err := time.Parse(format, timeStr)
			if err != nil {
				return 0
			}
			return t.Hour()
		},
	}

	// Parse the template
	tmpl, err := template.New("userPrompt").Funcs(funcMap).Parse(p.UserPromptFormat)
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	// Apply the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return buf.String(), nil
}

// ValidatePersonality checks if a personality is valid
func ValidatePersonality(p Personality) error {
	if p.Name == "" {
		return errors.New("personality name cannot be empty")
	}

	if p.SystemPrompt == "" {
		return errors.New("system prompt cannot be empty")
	}

	if p.UserPromptFormat == "" {
		return errors.New("user prompt format cannot be empty")
	}

	// Try to parse the template to validate it
	_, err := template.New("validation").Parse(p.UserPromptFormat)
	if err != nil {
		return fmt.Errorf("invalid user prompt template: %w", err)
	}

	if p.Temperature < 0 || p.Temperature > 1 {
		return errors.New("temperature must be between 0 and 1")
	}

	return nil
}
