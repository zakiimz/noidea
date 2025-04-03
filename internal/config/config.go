package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	// LLM provider settings
	LLMEnabled bool
	LLMProvider string
	LLMModel string
	APIKey string
	
	// Future config options can be added here
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		LLMEnabled: false,
		LLMProvider: "xai",
		LLMModel: "grok-2-1212",
		APIKey: "",
	}
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() Config {
	// Load default config
	config := DefaultConfig()
	
	// Try to load .env file from current directory
	_ = godotenv.Load()
	
	// Try to load from home directory
	home, err := os.UserHomeDir()
	if err == nil {
		_ = godotenv.Load(filepath.Join(home, ".noidea", ".env"))
	}
	
	// Check environment variables (these have priority over .env file)
	if apiKey := os.Getenv("XAI_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
		config.LLMEnabled = true
		config.LLMProvider = "xai"
	} else if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
		config.LLMEnabled = true
		config.LLMProvider = "openai"
	} else if apiKey := os.Getenv("DEEPSEEK_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
		config.LLMEnabled = true
		config.LLMProvider = "deepseek"
	}
	
	// Allow overriding model
	if model := os.Getenv("NOIDEA_MODEL"); model != "" {
		config.LLMModel = model
	}
	
	// Allow enabling/disabling LLM
	if enabled := os.Getenv("NOIDEA_LLM_ENABLED"); enabled == "false" {
		config.LLMEnabled = false
	} else if enabled == "true" {
		config.LLMEnabled = true
	}
	
	return config
} 