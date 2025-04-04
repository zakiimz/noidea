package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Config operation flags
	showConfig   bool
	initConfig   bool
	validateFlag bool

	// Config file path
	configPath string
)

func init() {
	rootCmd.AddCommand(configCmd)

	// Add flags
	configCmd.Flags().BoolVarP(&showConfig, "show", "s", false, "Show current configuration")
	configCmd.Flags().BoolVarP(&initConfig, "init", "i", false, "Initialize a new config file")
	configCmd.Flags().BoolVarP(&validateFlag, "validate", "v", false, "Validate the current configuration")
	configCmd.Flags().StringVarP(&configPath, "path", "p", "", "Path to config file (default: ~/.noidea/config.toml)")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage noidea configuration",
	Long:  `View, edit, and validate noidea configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Determine config path if not specified
		if configPath == "" {
			home, err := os.UserHomeDir()
			if err == nil {
				configPath = filepath.Join(home, ".noidea", "config.toml")
			} else {
				configPath = ".noidea.toml"
			}
		}

		// Load configuration
		cfg := config.LoadConfig()

		// Show configuration
		if showConfig {
			printConfig(cfg)
			return
		}

		// Validate configuration
		if validateFlag {
			issues := config.ValidateConfig(cfg)
			if len(issues) == 0 {
				fmt.Println(color.GreenString("✓ Configuration is valid"))
			} else {
				fmt.Println(color.RedString("✗ Configuration has issues:"))
				for _, issue := range issues {
					fmt.Println(color.YellowString("  - " + issue))
				}
			}
			return
		}

		// Initialize new config
		if initConfig {
			createConfigInteractive(configPath)
			return
		}

		// If no flag specified, show help
		cmd.Help()
	},
}

// printConfig displays the current configuration
func printConfig(cfg config.Config) {
	fmt.Println(color.CyanString("🧠 noidea configuration:"))

	fmt.Println(color.CyanString("\n[LLM]"))
	fmt.Printf("Enabled: %v\n", cfg.LLM.Enabled)
	fmt.Printf("Provider: %s\n", cfg.LLM.Provider)

	// Don't show the full API key for security
	apiKey := cfg.LLM.APIKey
	if apiKey != "" {
		if len(apiKey) > 8 {
			apiKey = apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
		} else {
			apiKey = "***"
		}
	}

	fmt.Printf("API Key: %s\n", apiKey)
	fmt.Printf("Model: %s\n", cfg.LLM.Model)
	fmt.Printf("Temperature: %.1f\n", cfg.LLM.Temperature)

	fmt.Println(color.CyanString("\n[Moai]"))
	fmt.Printf("Use Lint: %v\n", cfg.Moai.UseLint)
	fmt.Printf("Faces Mode: %s\n", cfg.Moai.FacesMode)
}

// createConfigInteractive creates a new config file with user input
func createConfigInteractive(path string) {
	// Start with default config
	cfg := config.DefaultConfig()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(color.CyanString("Creating a new configuration file at:"), path)

	// LLM Settings
	fmt.Println(color.CyanString("\nLLM Settings:"))

	// Enable LLM
	fmt.Print("Enable LLM? (y/n): ")
	response, _ := reader.ReadString('\n')
	cfg.LLM.Enabled = strings.TrimSpace(response) == "y"

	if cfg.LLM.Enabled {
		// Provider
		fmt.Print("Provider (xai, openai, deepseek): ")
		response, _ = reader.ReadString('\n')
		provider := strings.TrimSpace(response)
		if provider != "" {
			cfg.LLM.Provider = provider
		}

		// API Key
		fmt.Print("API Key: ")
		response, _ = reader.ReadString('\n')
		apiKey := strings.TrimSpace(response)
		if apiKey != "" {
			cfg.LLM.APIKey = apiKey
		}

		// Model
		fmt.Print("Model: ")
		response, _ = reader.ReadString('\n')
		model := strings.TrimSpace(response)
		if model != "" {
			cfg.LLM.Model = model
		}

		// Temperature
		fmt.Print("Temperature (0.0-1.0): ")
		response, _ = reader.ReadString('\n')
		temp := strings.TrimSpace(response)
		if temp != "" {
			if t, err := strconv.ParseFloat(temp, 64); err == nil {
				cfg.LLM.Temperature = t
			}
		}
	}

	// Moai Settings
	fmt.Println(color.CyanString("\nMoai Settings:"))

	// Use Lint
	fmt.Print("Enable linting feedback? (y/n): ")
	response, _ = reader.ReadString('\n')
	cfg.Moai.UseLint = strings.TrimSpace(response) == "y"

	// Faces Mode
	fmt.Print("Faces mode (random, sequential, mood): ")
	response, _ = reader.ReadString('\n')
	facesMode := strings.TrimSpace(response)
	if facesMode != "" {
		cfg.Moai.FacesMode = facesMode
	}
	
	// Personality
	fmt.Println(color.CyanString("\nPersonality Settings:"))
	fmt.Println("1. Professional with Sass (professional with a touch of wit)")
	fmt.Println("2. Snarky Code Reviewer (witty and sarcastic)")
	fmt.Println("3. Supportive Mentor (encouraging and positive)")
	fmt.Println("4. Git Expert (technical and professional)")
	fmt.Println("5. Motivational Speaker (enthusiastic and energetic)")
	
	fmt.Print("Choose personality (1-5, default: 1): ")
	response, _ = reader.ReadString('\n')
	personalityChoice := strings.TrimSpace(response)
	
	// Map choices to personality names
	switch personalityChoice {
	case "2":
		cfg.Moai.Personality = "snarky_reviewer"
	case "3":
		cfg.Moai.Personality = "supportive_mentor"
	case "4":
		cfg.Moai.Personality = "git_expert"
	case "5":
		cfg.Moai.Personality = "motivational_speaker"
	default:
		// Default or "1" option
		cfg.Moai.Personality = "professional_sass"
	}

	// Validate the config
	issues := config.ValidateConfig(cfg)
	if len(issues) > 0 {
		fmt.Println(color.YellowString("\nWarning: Configuration has issues:"))
		for _, issue := range issues {
			fmt.Println(color.YellowString("  - " + issue))
		}

		fmt.Print("Save anyway? (y/n): ")
		response, _ = reader.ReadString('\n')
		if strings.TrimSpace(response) != "y" {
			fmt.Println("Configuration not saved.")
			return
		}
	}

	// Save the config
	err := config.SaveConfig(cfg)
	if err != nil {
		fmt.Println(color.RedString("Error saving configuration:"), err)
		return
	}

	fmt.Println(color.GreenString("Configuration saved to:"), path)
}
