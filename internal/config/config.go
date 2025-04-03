package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

// LLMConfig holds LLM-specific configuration
type LLMConfig struct {
	Enabled     bool    `toml:"enabled"`
	Provider    string  `toml:"provider"`
	APIKey      string  `toml:"api_key"`
	Model       string  `toml:"model"`
	Temperature float64 `toml:"temperature"`
}

// MoaiConfig holds Moai-specific configuration
type MoaiConfig struct {
	UseLint         bool   `toml:"use_lint"`
	FacesMode       string `toml:"faces_mode"`
	Personality     string `toml:"personality"`
	PersonalityFile string `toml:"personality_file"`
	IncludeHistory  bool   `toml:"include_history"`
}

// Config holds the application configuration
type Config struct {
	LLM  LLMConfig  `toml:"llm"`
	Moai MoaiConfig `toml:"moai"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		LLM: LLMConfig{
			Enabled:     false,
			Provider:    "xai",
			APIKey:      "",
			Model:       "grok-2-1212",
			Temperature: 0.7,
		},
		Moai: MoaiConfig{
			UseLint:         false,
			FacesMode:       "random",
			Personality:     "snarky_reviewer",
			PersonalityFile: "",
			IncludeHistory:  false,
		},
	}
}

// ConfigPaths returns the possible configuration file paths
func ConfigPaths() []string {
	// Check for configuration in current directory
	paths := []string{".noidea.toml"}

	// Check in home directory
	home, err := os.UserHomeDir()
	if err == nil {
		paths = append(paths, filepath.Join(home, ".noidea", "config.toml"))
		paths = append(paths, filepath.Join(home, ".config", "noidea", "config.toml"))
	}

	return paths
}

// LoadConfig loads configuration from files and environment variables
func LoadConfig() Config {
	// Start with default config
	config := DefaultConfig()

	// Load from config files
	for _, path := range ConfigPaths() {
		if _, err := os.Stat(path); err == nil {
			_, err := toml.DecodeFile(path, &config)
			if err == nil {
				// Successfully loaded a config file, break
				break
			}
		}
	}

	// Load from .env files
	_ = godotenv.Load()

	// Try to load from home directory .env file
	home, err := os.UserHomeDir()
	if err == nil {
		_ = godotenv.Load(filepath.Join(home, ".noidea", ".env"))
	}

	// Override with environment variables (higher priority)
	if apiKey := os.Getenv("XAI_API_KEY"); apiKey != "" {
		config.LLM.APIKey = apiKey
		config.LLM.Provider = "xai"
		config.LLM.Enabled = true
	} else if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		config.LLM.APIKey = apiKey
		config.LLM.Provider = "openai"
		config.LLM.Enabled = true
	} else if apiKey := os.Getenv("DEEPSEEK_API_KEY"); apiKey != "" {
		config.LLM.APIKey = apiKey
		config.LLM.Provider = "deepseek"
		config.LLM.Enabled = true
	}

	// Other environment variable overrides
	if model := os.Getenv("NOIDEA_MODEL"); model != "" {
		config.LLM.Model = model
	}

	if temp := os.Getenv("NOIDEA_TEMPERATURE"); temp != "" {
		// Ignore errors - if we can't parse, we keep the default
		if t, err := ParseFloat(temp); err == nil {
			config.LLM.Temperature = t
		}
	}

	if enabled := os.Getenv("NOIDEA_LLM_ENABLED"); enabled != "" {
		config.LLM.Enabled = enabled == "true" || enabled == "1" || enabled == "yes"
	}

	if facesMode := os.Getenv("NOIDEA_FACES_MODE"); facesMode != "" {
		config.Moai.FacesMode = facesMode
	}

	if useLint := os.Getenv("NOIDEA_USE_LINT"); useLint != "" {
		config.Moai.UseLint = useLint == "true" || useLint == "1" || useLint == "yes"
	}
	
	if personality := os.Getenv("NOIDEA_PERSONALITY"); personality != "" {
		config.Moai.Personality = personality
	}
	
	if personalityFile := os.Getenv("NOIDEA_PERSONALITY_FILE"); personalityFile != "" {
		config.Moai.PersonalityFile = personalityFile
	}

	if includeHistory := os.Getenv("NOIDEA_INCLUDE_HISTORY"); includeHistory != "" {
		config.Moai.IncludeHistory = includeHistory == "true" || includeHistory == "1" || includeHistory == "yes"
	}

	return config
}

// SaveConfig saves the configuration to a file
func SaveConfig(config Config, path string) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create or truncate the file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	// Encode the config to TOML
	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// ValidateConfig checks if the configuration is valid
func ValidateConfig(config Config) []string {
	var issues []string

	// Check LLM configuration
	if config.LLM.Enabled {
		if config.LLM.APIKey == "" {
			issues = append(issues, "LLM is enabled but no API key is provided")
		}

		if config.LLM.Model == "" {
			issues = append(issues, "LLM is enabled but no model is specified")
		}

		if config.LLM.Temperature < 0 || config.LLM.Temperature > 1 {
			issues = append(issues, "Temperature must be between 0 and 1")
		}

		// Check provider
		switch config.LLM.Provider {
		case "xai", "openai", "deepseek":
			// Valid providers
		default:
			issues = append(issues, fmt.Sprintf("Unknown provider: %s", config.LLM.Provider))
		}
	}

	// Check Moai configuration
	if config.Moai.FacesMode != "random" && config.Moai.FacesMode != "sequential" && config.Moai.FacesMode != "mood" {
		issues = append(issues, fmt.Sprintf("Unknown faces mode: %s", config.Moai.FacesMode))
	}

	return issues
}

// Helper function to parse float
func ParseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
} 