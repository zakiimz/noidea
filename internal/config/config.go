// Package config provides configuration management for the noidea tool.
// It handles loading, saving, and validating user configuration settings.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config represents the application configuration
type Config struct {
	// LLM contains settings for the AI language model integration
	LLM struct {
		Enabled     bool    `json:"enabled"`
		Provider    string  `json:"provider"`    // "xai", "openai", "deepseek"
		APIKey      string  `json:"api_key"`     // API key for the language model provider
		Model       string  `json:"model"`       // Model name to use
		Temperature float64 `json:"temperature"` // Temperature for AI responses (0.0-1.0)
	} `json:"llm"`

	// Moai contains settings for the Moai feedback system
	Moai struct {
		UseLint         bool   `json:"use_lint"`          // Include linting feedback
		FacesMode       string `json:"faces_mode"`        // "random", "sequential", "mood"
		Personality     string `json:"personality"`       // Selected personality
		PersonalityFile string `json:"personality_file"`  // Custom personality definitions
	} `json:"moai"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	var cfg Config
	
	// LLM settings
	cfg.LLM.Enabled = false
	cfg.LLM.Provider = "xai"
	cfg.LLM.Model = "grok-2-1212"
	cfg.LLM.Temperature = 0.7
	
	// Moai settings
	cfg.Moai.UseLint = false
	cfg.Moai.FacesMode = "random"
	cfg.Moai.Personality = "professional_sass"
	
	// Get home directory for default personality file path
	homeDir, err := os.UserHomeDir()
	if err == nil {
		cfg.Moai.PersonalityFile = filepath.Join(homeDir, ".noidea", "personalities.json")
	} else {
		// Fallback to current directory if we can't get home dir
		cfg.Moai.PersonalityFile = "personalities.json"
	}
	
	return cfg
}

// LoadConfig loads the configuration from the default location or environment variables
// If the config file doesn't exist, it returns the default config
func LoadConfig() Config {
	// Start with default config
	cfg := DefaultConfig()
	
	// Try to get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not determine user home directory: %v\n", err)
		// Continue with defaults
		return applyEnvironmentOverrides(cfg)
	}
	
	// Config directory path
	configDir := filepath.Join(homeDir, ".noidea")
	
	// Config file path
	configFile := filepath.Join(configDir, "config.json")
	
	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Info: No config file found at %s, using defaults\n", configFile)
		// Check also for .toml format for backward compatibility
		tomlConfigFile := filepath.Join(configDir, "config.toml")
		if _, err := os.Stat(tomlConfigFile); os.IsNotExist(err) {
			return applyEnvironmentOverrides(cfg)
		}
		configFile = tomlConfigFile
	}
	
	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not read config file %s: %v\n", configFile, err)
		return applyEnvironmentOverrides(cfg)
	}
	
	// Parse config based on file extension
	if filepath.Ext(configFile) == ".toml" {
		// Handle TOML format if needed
		fmt.Fprintf(os.Stderr, "Warning: TOML format not fully supported yet\n")
		// TODO: Implement TOML parsing
	} else {
		// Parse JSON
		if err := json.Unmarshal(data, &cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not parse config file %s: %v\n", configFile, err)
			// Continue with defaults
			return applyEnvironmentOverrides(DefaultConfig())
		}
	}
	
	// Ensure all fields are set properly
	ensureDefaults(&cfg)
	
	// Apply environment variable overrides
	return applyEnvironmentOverrides(cfg)
}

// applyEnvironmentOverrides applies environment variable settings to override config file values
func applyEnvironmentOverrides(cfg Config) Config {
	// LLM settings
	if val := os.Getenv("NOIDEA_LLM_ENABLED"); val != "" {
		cfg.LLM.Enabled = val == "true" || val == "1" || val == "yes"
	}
	
	if val := os.Getenv("NOIDEA_LLM_PROVIDER"); val != "" {
		cfg.LLM.Provider = val
	}
	
	// API keys from multiple possible environment variables
	if val := os.Getenv("NOIDEA_API_KEY"); val != "" {
		cfg.LLM.APIKey = strings.TrimSpace(val)
	}
	
	// Provider-specific API keys take precedence
	switch cfg.LLM.Provider {
	case "xai":
		if val := os.Getenv("XAI_API_KEY"); val != "" {
			// Ensure the key is properly formatted and trimmed
			cfg.LLM.APIKey = strings.TrimSpace(val)
			// Log a warning if key doesn't have expected prefix
			if !strings.HasPrefix(cfg.LLM.APIKey, "xai-") {
				fmt.Fprintf(os.Stderr, "Warning: XAI API key doesn't start with 'xai-' prefix\n")
			}
		}
	case "openai":
		if val := os.Getenv("OPENAI_API_KEY"); val != "" {
			cfg.LLM.APIKey = strings.TrimSpace(val)
		}
	case "deepseek":
		if val := os.Getenv("DEEPSEEK_API_KEY"); val != "" {
			cfg.LLM.APIKey = strings.TrimSpace(val)
		}
	}
	
	if val := os.Getenv("NOIDEA_MODEL"); val != "" {
		cfg.LLM.Model = val
	}
	
	if val := os.Getenv("NOIDEA_TEMPERATURE"); val != "" {
		if temp, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.LLM.Temperature = temp
		}
	}
	
	// Moai settings
	if val := os.Getenv("NOIDEA_USE_LINT"); val != "" {
		cfg.Moai.UseLint = val == "true" || val == "1" || val == "yes"
	}
	
	if val := os.Getenv("NOIDEA_FACES_MODE"); val != "" {
		cfg.Moai.FacesMode = val
	}
	
	if val := os.Getenv("NOIDEA_PERSONALITY"); val != "" {
		cfg.Moai.Personality = val
	}
	
	if val := os.Getenv("NOIDEA_PERSONALITY_FILE"); val != "" {
		cfg.Moai.PersonalityFile = val
	}
	
	return cfg
}

// ensureDefaults ensures that all config fields have valid values
// by applying defaults to any missing or invalid values
func ensureDefaults(cfg *Config) {
	defaultCfg := DefaultConfig()
	
	// Ensure LLM defaults
	if cfg.LLM.Provider == "" {
		cfg.LLM.Provider = defaultCfg.LLM.Provider
	}
	
	if cfg.LLM.Model == "" {
		cfg.LLM.Model = defaultCfg.LLM.Model
	}
	
	if cfg.LLM.Temperature <= 0 || cfg.LLM.Temperature > 1.0 {
		cfg.LLM.Temperature = defaultCfg.LLM.Temperature
	}
	
	// Ensure Moai defaults
	if cfg.Moai.FacesMode == "" {
		cfg.Moai.FacesMode = defaultCfg.Moai.FacesMode
	}
	
	if cfg.Moai.Personality == "" {
		cfg.Moai.Personality = defaultCfg.Moai.Personality
	}
	
	if cfg.Moai.PersonalityFile == "" {
		cfg.Moai.PersonalityFile = defaultCfg.Moai.PersonalityFile
	}
}

// SaveConfig saves the configuration to the default location
func SaveConfig(cfg Config) error {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	// Create config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".noidea")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Config file path
	configFile := filepath.Join(configDir, "config.json")
	
	// Marshal config to JSON
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	
	// Write config file
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// ValidateConfig checks the configuration for errors or inconsistencies
// Returns a list of issues or an empty slice if the config is valid
func ValidateConfig(config Config) []string {
	var issues []string
	
	// Validate LLM settings
	if config.LLM.Enabled {
		// Check that provider is valid
		validProviders := map[string]bool{
			"xai":      true,
			"openai":   true,
			"deepseek": true,
		}
		
		if !validProviders[config.LLM.Provider] {
			issues = append(issues, fmt.Sprintf("Unknown provider: %s", config.LLM.Provider))
		}
		
		// Check that API key is set
		if config.LLM.APIKey == "" {
			issues = append(issues, "API key is required when LLM is enabled")
		}
		
		// Check temperature range
		if config.LLM.Temperature < 0 || config.LLM.Temperature > 1.0 {
			issues = append(issues, fmt.Sprintf("Temperature value must be between 0.0 and 1.0 (got %.1f)", 
				config.LLM.Temperature))
		}
	}
	
	// Validate Moai settings
	validFacesModes := map[string]bool{
		"random":     true,
		"sequential": true,
		"mood":       true,
	}
	
	if !validFacesModes[config.Moai.FacesMode] {
		issues = append(issues, fmt.Sprintf("Unknown faces mode: %s", config.Moai.FacesMode))
	}
	
	// Check that personality file exists if a custom personality is set
	if config.Moai.Personality != "default" && 
	   config.Moai.Personality != "friendly" && 
	   config.Moai.Personality != "professional" && 
	   config.Moai.Personality != "sarcastic" {
		
		// Check if the file exists
		if _, err := os.Stat(config.Moai.PersonalityFile); os.IsNotExist(err) {
			issues = append(issues, "Custom personality file not found: " + config.Moai.PersonalityFile)
		}
	}
	
	return issues
}

// ParseFloat parses a string to a float64 with a default value if parsing fails
func ParseFloat(s string, defaultVal float64) float64 {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		return defaultVal
	}
	return f
}
