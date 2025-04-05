package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/secure"
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

	// Add key management commands
	configCmd.AddCommand(configAPIKeyCmd)
	configCmd.AddCommand(configAPIKeyRemoveCmd)
	configCmd.AddCommand(configAPIKeyStatusCmd)
	configCmd.AddCommand(configAPIKeyCleanEnvCmd)

	// Add flags to API key commands
	configAPIKeyCmd.Flags().Bool("skip-validation", false, "Skip API key validation")
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
		fmt.Println("\nAPI Key (will be stored securely):")
		apiKey := readPassword("Enter your API key: ")

		if apiKey != "" {
			// Save the API key securely
			err := config.SaveAPIKey(cfg.LLM.Provider, apiKey)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to save API key securely: %v\n", err)
				fmt.Fprintf(os.Stderr, "Your API key will need to be provided via environment variables.\n")
			}

			// For UX purposes, show we've set the API key in memory
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

// readPassword prompts for a password securely without echoing input
func readPassword(prompt string) string {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // Add a newline after the password input
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return ""
	}
	return strings.TrimSpace(string(password))
}

// configAPIKeyCmd handles API key configuration
var configAPIKeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Configure API key securely",
	Long:  `Securely store an API key for use with LLM features.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if skip validation flag is set
		skipValidation, _ := cmd.Flags().GetBool("skip-validation")

		// Get config
		cfg := config.LoadConfig()

		fmt.Printf("Current provider: %s\n", cfg.LLM.Provider)

		// Allow changing provider
		fmt.Printf("Would you like to change the provider? (current: %s) [y/N]: ", cfg.LLM.Provider)
		var changeProvider string
		fmt.Scanln(&changeProvider)

		provider := cfg.LLM.Provider

		if strings.ToLower(changeProvider) == "y" || strings.ToLower(changeProvider) == "yes" {
			fmt.Println("\nAvailable providers:")
			fmt.Println("1. xAI (Grok)")
			fmt.Println("2. OpenAI")
			fmt.Println("3. DeepSeek (experimental)")

			var choice string
			fmt.Print("Select provider (1-3): ")
			fmt.Scanln(&choice)

			switch choice {
			case "1":
				provider = "xai"
			case "2":
				provider = "openai"
			case "3":
				provider = "deepseek"
			default:
				// Keep current provider
			}

			if provider != cfg.LLM.Provider {
				fmt.Printf("Changing provider to: %s\n", provider)
			}
		}

		// Prompt for API key securely
		apiKey := readPassword("\nEnter API key (input will be hidden): ")

		if apiKey == "" {
			fmt.Println("No API key provided, operation cancelled.")
			return
		}

		// Validate the key before saving (unless skipped)
		if !skipValidation {
			fmt.Print("Validating API key... ")
			isValid, err := secure.ValidateAPIKey(provider, apiKey)
			if err != nil {
				fmt.Println(color.YellowString("Warning: Error during validation"))
				fmt.Printf("Error details: %v\n", err)

				// Ask if user wants to continue anyway
				fmt.Print("Continue saving this key anyway? [y/N]: ")
				var continueAnyway string
				fmt.Scanln(&continueAnyway)

				if strings.ToLower(continueAnyway) != "y" && strings.ToLower(continueAnyway) != "yes" {
					fmt.Println("Operation cancelled.")
					return
				}
			} else if !isValid {
				fmt.Println(color.RedString("Invalid"))
				fmt.Println("The API key was not accepted by the provider. It may be expired or incorrect.")

				// Ask if user wants to continue anyway
				fmt.Print("Continue saving this key anyway? [y/N]: ")
				var continueAnyway string
				fmt.Scanln(&continueAnyway)

				if strings.ToLower(continueAnyway) != "y" && strings.ToLower(continueAnyway) != "yes" {
					fmt.Println("Operation cancelled.")
					return
				}
			} else {
				fmt.Println(color.GreenString("Valid"))
			}
		} else {
			fmt.Println("Skipping API key validation as requested.")
		}

		// Save API key securely
		if err := config.SaveAPIKey(provider, apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save API key: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nAPI key saved securely.")

		// Inform the user about environment variables that should be removed
		if os.Getenv("XAI_API_KEY") != "" ||
			os.Getenv("OPENAI_API_KEY") != "" ||
			os.Getenv("DEEPSEEK_API_KEY") != "" ||
			os.Getenv("NOIDEA_API_KEY") != "" {
			fmt.Println(color.YellowString("\nWARNING: API key environment variables detected!"))
			fmt.Println("For secure storage to work properly, you should remove these environment variables:")

			if os.Getenv("XAI_API_KEY") != "" {
				fmt.Println("  - XAI_API_KEY")
			}
			if os.Getenv("OPENAI_API_KEY") != "" {
				fmt.Println("  - OPENAI_API_KEY")
			}
			if os.Getenv("DEEPSEEK_API_KEY") != "" {
				fmt.Println("  - DEEPSEEK_API_KEY")
			}
			if os.Getenv("NOIDEA_API_KEY") != "" {
				fmt.Println("  - NOIDEA_API_KEY")
			}

			fmt.Println("\nIn your shell, run:")
			fmt.Println("  unset XAI_API_KEY OPENAI_API_KEY DEEPSEEK_API_KEY NOIDEA_API_KEY")
			fmt.Println("\nOr remove these from your .bashrc/.zshrc file if they are set there.")
		}

		// Check if LLM is enabled
		if !cfg.LLM.Enabled {
			fmt.Println("\nNote: LLM features are currently disabled.")
			fmt.Print("Would you like to enable LLM features now? [y/N]: ")
			var enableLLM string
			fmt.Scanln(&enableLLM)

			if strings.ToLower(enableLLM) == "y" || strings.ToLower(enableLLM) == "yes" {
				cfg.LLM.Enabled = true
				cfg.LLM.Provider = provider

				if err := config.SaveConfig(cfg); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to update config: %v\n", err)
				} else {
					fmt.Println("LLM features enabled successfully.")
				}
			}
		}
	},
}

// configAPIKeyRemoveCmd handles API key removal
var configAPIKeyRemoveCmd = &cobra.Command{
	Use:   "apikey-remove",
	Short: "Remove stored API key",
	Long:  `Remove an API key from secure storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config
		cfg := config.LoadConfig()

		// Ask for confirmation
		fmt.Printf("Remove API key for %s? This cannot be undone. [y/N]: ", cfg.LLM.Provider)
		var confirm string
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
			fmt.Println("Operation cancelled.")
			return
		}

		// Delete the API key
		if err := config.DeleteAPIKey(cfg.LLM.Provider); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to remove API key: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("API key removed successfully.")
	},
}

// configAPIKeyStatusCmd shows the status of secure storage
var configAPIKeyStatusCmd = &cobra.Command{
	Use:   "apikey-status",
	Short: "Check API key storage status",
	Long:  `Check the status of secure API key storage on your system.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get secure storage status
		status := secure.GetSecureStorageStatus()

		fmt.Println(color.CyanString("Secure Storage Status:"))
		fmt.Printf("Platform: %s\n", status["platform"])

		// Check if we have keyring support
		if status["keyring"] == "available" {
			fmt.Printf("System keyring: %s\n", color.GreenString("Available"))
		} else {
			fmt.Printf("System keyring: %s (%s)\n",
				color.YellowString("Unavailable"),
				status["keyring"])
			fmt.Println("Using fallback encrypted storage.")
		}

		// Check fallback status
		fmt.Printf("Fallback storage: %s\n", status["fallback"])

		// Check if API key is set in environment
		envApiKey := ""
		envSource := ""

		if os.Getenv("XAI_API_KEY") != "" {
			envApiKey = os.Getenv("XAI_API_KEY")
			envSource = "XAI_API_KEY"
		} else if os.Getenv("OPENAI_API_KEY") != "" {
			envApiKey = os.Getenv("OPENAI_API_KEY")
			envSource = "OPENAI_API_KEY"
		} else if os.Getenv("DEEPSEEK_API_KEY") != "" {
			envApiKey = os.Getenv("DEEPSEEK_API_KEY")
			envSource = "DEEPSEEK_API_KEY"
		} else if os.Getenv("NOIDEA_API_KEY") != "" {
			envApiKey = os.Getenv("NOIDEA_API_KEY")
			envSource = "NOIDEA_API_KEY"
		}

		// Check if API key is in secure storage
		cfg := config.LoadConfig()
		secureApiKey, secureErr := secure.GetAPIKey(cfg.LLM.Provider)

		// Show information about both keys
		fmt.Println(color.CyanString("\nAPI Key Status:"))
		fmt.Printf("Provider: %s\n", cfg.LLM.Provider)

		// Environment variable key
		if envApiKey != "" {
			fmt.Printf("Environment: %s (from %s)\n", color.YellowString("Set"), envSource)
		} else {
			fmt.Printf("Environment: %s\n", color.HiBlackString("Not set"))
		}

		// Secure storage key
		if secureErr == nil && secureApiKey != "" {
			fmt.Printf("Secure storage: %s\n", color.GreenString("Set"))
		} else {
			fmt.Printf("Secure storage: %s\n", color.RedString("Not set"))
			if secureErr != nil && secureErr != secure.ErrKeyNotFound {
				fmt.Printf("  Error: %v\n", secureErr)
			}
		}

		// Show which key is actually being used
		fmt.Println(color.CyanString("\nActive Key:"))
		if envApiKey != "" {
			fmt.Printf("Using: %s (environment takes precedence)\n", color.YellowString("Environment key"))

			// Validate the environment key
			fmt.Print("Validating environment key... ")
			isValid, err := secure.ValidateAPIKey(cfg.LLM.Provider, envApiKey)
			if err != nil {
				fmt.Println(color.YellowString("Error during validation"))
				fmt.Printf("Error details: %v\n", err)
			} else if isValid {
				fmt.Println(color.GreenString("Valid"))
			} else {
				fmt.Println(color.RedString("Invalid"))
				fmt.Println("The environment API key was not accepted by the provider.")
			}

			// If we also have a secure key, validate that too
			if secureApiKey != "" && secureApiKey != envApiKey {
				fmt.Print("Validating secure storage key... ")
				isValid, err := secure.ValidateAPIKey(cfg.LLM.Provider, secureApiKey)
				if err != nil {
					fmt.Println(color.YellowString("Error during validation"))
					fmt.Printf("Error details: %v\n", err)
				} else if isValid {
					fmt.Println(color.GreenString("Valid"))
					fmt.Println(color.YellowString("\nRecommendation:"))
					fmt.Println("Your secure key is valid, but the environment key is being used.")
					fmt.Println("To use your secure key, unset the environment variables:")
					fmt.Println("  unset XAI_API_KEY OPENAI_API_KEY DEEPSEEK_API_KEY NOIDEA_API_KEY")
				} else {
					fmt.Println(color.RedString("Invalid"))
				}
			}
		} else if secureApiKey != "" {
			fmt.Printf("Using: %s\n", color.GreenString("Secure storage key"))

			// Validate the secure key
			fmt.Print("Validating secure key... ")
			isValid, err := secure.ValidateAPIKey(cfg.LLM.Provider, secureApiKey)
			if err != nil {
				fmt.Println(color.YellowString("Error during validation"))
				fmt.Printf("Error details: %v\n", err)
			} else if isValid {
				fmt.Println(color.GreenString("Valid"))
			} else {
				fmt.Println(color.RedString("Invalid"))
				fmt.Println("The API key was not accepted by the provider. It may be expired or incorrect.")
				fmt.Println("Use 'noidea config apikey' to update your key.")
			}
		} else {
			fmt.Printf("Using: %s\n", color.RedString("No key available"))
			fmt.Println("No API key found in environment or secure storage.")
			fmt.Println("Use 'noidea config apikey' to set up a key.")
		}
	},
}

// configAPIKeyCleanEnvCmd helps users clean up environment variables
var configAPIKeyCleanEnvCmd = &cobra.Command{
	Use:   "clean-env",
	Short: "Generate commands to clean environment variables",
	Long:  `Generate commands to remove API key environment variables that might override secure storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check for API keys in environment
		hasEnvVars := false
		envVars := []string{}

		if os.Getenv("XAI_API_KEY") != "" {
			hasEnvVars = true
			envVars = append(envVars, "XAI_API_KEY")
		}
		if os.Getenv("OPENAI_API_KEY") != "" {
			hasEnvVars = true
			envVars = append(envVars, "OPENAI_API_KEY")
		}
		if os.Getenv("DEEPSEEK_API_KEY") != "" {
			hasEnvVars = true
			envVars = append(envVars, "DEEPSEEK_API_KEY")
		}
		if os.Getenv("NOIDEA_API_KEY") != "" {
			hasEnvVars = true
			envVars = append(envVars, "NOIDEA_API_KEY")
		}

		if !hasEnvVars {
			fmt.Println("No API key environment variables detected!")
			return
		}

		// Check if secure storage has a key
		cfg := config.LoadConfig()
		secureApiKey, secureErr := secure.GetAPIKey(cfg.LLM.Provider)

		if secureErr != nil || secureApiKey == "" {
			fmt.Println(color.YellowString("Warning: No API key found in secure storage!"))
			fmt.Println("You should set up a secure key before removing environment variables:")
			fmt.Println("  noidea config apikey")
			fmt.Println("")
		} else {
			// Validate the secure key
			fmt.Print("Validating secure key... ")
			isValid, err := secure.ValidateAPIKey(cfg.LLM.Provider, secureApiKey)
			if err != nil {
				fmt.Println(color.YellowString("Warning: Error during validation"))
				fmt.Printf("Error details: %v\n", err)
			} else if !isValid {
				fmt.Println(color.RedString("Invalid"))
				fmt.Println(color.YellowString("Warning: Your secure API key is invalid!"))
				fmt.Println("You should update it before removing environment variables:")
				fmt.Println("  noidea config apikey")
				fmt.Println("")
			} else {
				fmt.Println(color.GreenString("Valid"))
			}
		}

		// Show detected environment variables
		fmt.Println(color.CyanString("Detected API key environment variables:"))
		for _, envVar := range envVars {
			fmt.Printf("  - %s\n", envVar)
		}

		// Generate commands to clean environment
		fmt.Println(color.CyanString("\nTo clean your environment, run:"))

		// For Bash/Zsh
		fmt.Println(color.GreenString("\n# For Bash/Zsh:"))
		fmt.Printf("unset %s\n", strings.Join(envVars, " "))

		// For Fish
		fmt.Println(color.GreenString("\n# For Fish shell:"))
		for _, envVar := range envVars {
			fmt.Printf("set -e %s\n", envVar)
		}

		// RC files
		fmt.Println(color.CyanString("\nIf these variables are set in your shell startup files:"))
		fmt.Println("Check these files and remove or comment out any lines setting these variables:")
		fmt.Println("  - ~/.bashrc")
		fmt.Println("  - ~/.zshrc")
		fmt.Println("  - ~/.config/fish/config.fish")
		fmt.Println("  - ~/.profile")
		fmt.Println("  - ~/.env (or other .env files)")
	},
}
