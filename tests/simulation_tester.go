package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TestCase defines a simulation test case
type TestCase struct {
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	Description string   `json:"description"`
	Repetitions int      `json:"repetitions"`
}

// TestSuite defines a set of test cases
type TestSuite struct {
	Name      string     `json:"name"`
	OutputDir string     `json:"output_dir"`
	TestCases []TestCase `json:"test_cases"`
}

// LoadTestSuite loads a test suite from a JSON file
func LoadTestSuite(filename string) (TestSuite, error) {
	var suite TestSuite
	data, err := os.ReadFile(filename)
	if err != nil {
		return suite, err
	}

	err = json.Unmarshal(data, &suite)
	return suite, err
}

// SaveTestSuite saves a test suite to a JSON file
func SaveTestSuite(suite TestSuite, filename string) error {
	data, err := json.MarshalIndent(suite, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// RunTestSuite runs all test cases in a test suite and captures output
func RunTestSuite(suite TestSuite) error {
	// Create output directory
	outputDir := suite.OutputDir
	if outputDir == "" {
		outputDir = "test_results"
	}

	os.MkdirAll(outputDir, 0755)

	// Save summary of test suite
	summaryFile := filepath.Join(outputDir, "summary.txt")
	summary, err := os.Create(summaryFile)
	if err != nil {
		return err
	}
	defer summary.Close()

	fmt.Fprintf(summary, "# Test Suite: %s\n", suite.Name)
	fmt.Fprintf(summary, "Run at: %s\n\n", time.Now().Format(time.RFC3339))

	// Try to load .env file from project root
	LoadEnvFile("../.env")

	// Check if we have an API key, if not use a mock one
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
	}

	// Log API key status (truncated for security)
	if apiKey != "" {
		maskedKey := "****"
		if len(apiKey) > 8 {
			maskedKey = apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
		}
		fmt.Printf("Using API key: %s\n", maskedKey)
	} else {
		fmt.Println("No API key found, using mock key for testing")
		os.Setenv("XAI_API_KEY", "mock-api-key-for-testing")
	}

	// Always enable LLM for testing
	os.Setenv("NOIDEA_LLM_ENABLED", "true")

	// Run each test case
	for _, testCase := range suite.TestCases {
		fmt.Printf("Running test case: %s\n", testCase.Name)

		// Add test case to summary
		fmt.Fprintf(summary, "## %s\n", testCase.Name)
		fmt.Fprintf(summary, "Description: %s\n", testCase.Description)
		fmt.Fprintf(summary, "Command: `%s %s`\n", testCase.Command, strings.Join(testCase.Args, " "))
		fmt.Fprintf(summary, "Repetitions: %d\n\n", testCase.Repetitions)

		// Create directory for test case results
		testCaseDir := filepath.Join(outputDir, testCase.Name)
		os.MkdirAll(testCaseDir, 0755)

		// Save info file
		infoFile := filepath.Join(testCaseDir, "info.txt")
		info, err := os.Create(infoFile)
		if err != nil {
			return err
		}

		fmt.Fprintf(info, "Test Case: %s\n", testCase.Name)
		fmt.Fprintf(info, "Description: %s\n", testCase.Description)
		fmt.Fprintf(info, "Command: %s %s\n", testCase.Command, strings.Join(testCase.Args, " "))
		fmt.Fprintf(info, "Repetitions: %d\n", testCase.Repetitions)
		info.Close()

		// Run the command multiple times
		for i := 1; i <= testCase.Repetitions; i++ {
			fmt.Printf("  Run %d/%d...\n", i, testCase.Repetitions)

			// Run the command
			cmd := exec.Command(testCase.Command, testCase.Args...)

			// Capture output
			output, err := cmd.CombinedOutput()

			// Save output to file
			outputFile := filepath.Join(testCaseDir, fmt.Sprintf("run_%d.txt", i))
			err = os.WriteFile(outputFile, output, 0644)
			if err != nil {
				fmt.Printf("  Error saving output: %v\n", err)
			}
		}
	}

	fmt.Printf("Test suite completed. Results saved to: %s\n", outputDir)
	return nil
}

// CreateDefaultTestSuites creates default test suites for each main feature
func CreateDefaultTestSuites() error {
	// Create test suites directory
	suiteDir := "test_suites"
	err := os.MkdirAll(suiteDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create test suites directory: %w", err)
	}

	// Summary command test suite
	summarySuite := TestSuite{
		Name:      "Summary Command Tests",
		OutputDir: "results/summary",
		TestCases: []TestCase{
			{
				Name:        "default_summary",
				Command:     "../noidea",
				Args:        []string{"summary"},
				Description: "Default summary command with default settings",
				Repetitions: 3,
			},
			{
				Name:        "summary_14_days",
				Command:     "../noidea",
				Args:        []string{"summary", "--days", "14"},
				Description: "Summary for the last 14 days",
				Repetitions: 3,
			},
			{
				Name:        "summary_snarky_reviewer",
				Command:     "../noidea",
				Args:        []string{"summary", "--personality", "snarky_reviewer"},
				Description: "Summary with snarky reviewer personality",
				Repetitions: 3,
			},
			{
				Name:        "summary_supportive_mentor",
				Command:     "../noidea",
				Args:        []string{"summary", "--personality", "supportive_mentor"},
				Description: "Summary with supportive mentor personality",
				Repetitions: 3,
			},
			{
				Name:        "summary_git_expert",
				Command:     "../noidea",
				Args:        []string{"summary", "--personality", "git_expert"},
				Description: "Summary with git expert personality",
				Repetitions: 3,
			},
			{
				Name:        "summary_markdown_export",
				Command:     "../noidea",
				Args:        []string{"summary", "--export", "markdown"},
				Description: "Summary with markdown export",
				Repetitions: 2,
			},
			{
				Name:        "summary_html_export",
				Command:     "../noidea",
				Args:        []string{"summary", "--export", "html"},
				Description: "Summary with HTML export",
				Repetitions: 2,
			},
			{
				Name:        "summary_stats_only",
				Command:     "../noidea",
				Args:        []string{"summary", "--stats-only"},
				Description: "Summary with statistics only (no AI insights)",
				Repetitions: 1,
			},
		},
	}

	// Moai command test suite
	moaiSuite := TestSuite{
		Name:      "Moai Command Tests",
		OutputDir: "results/moai",
		TestCases: []TestCase{
			{
				Name:        "default_moai",
				Command:     "../noidea",
				Args:        []string{"moai", "Test commit message"},
				Description: "Default moai command with static message",
				Repetitions: 5,
			},
			{
				Name:        "moai_with_ai",
				Command:     "../noidea",
				Args:        []string{"moai", "--ai", "Feature implementation complete"},
				Description: "Moai with AI feedback enabled",
				Repetitions: 5,
			},
			{
				Name:        "moai_with_diff",
				Command:     "../noidea",
				Args:        []string{"moai", "--ai", "--diff", "Bug fix for critical issue"},
				Description: "Moai with diff context",
				Repetitions: 3,
			},
			{
				Name:        "moai_snarky_reviewer",
				Command:     "../noidea",
				Args:        []string{"moai", "--ai", "--personality", "snarky_reviewer", "Code cleanup"},
				Description: "Moai with snarky reviewer personality",
				Repetitions: 5,
			},
			{
				Name:        "moai_supportive_mentor",
				Command:     "../noidea",
				Args:        []string{"moai", "--ai", "--personality", "supportive_mentor", "First attempt at new feature"},
				Description: "Moai with supportive mentor personality",
				Repetitions: 5,
			},
			{
				Name:        "moai_git_expert",
				Command:     "../noidea",
				Args:        []string{"moai", "--ai", "--personality", "git_expert", "Merge branch 'feature/x' into main"},
				Description: "Moai with git expert personality",
				Repetitions: 5,
			},
			{
				Name:        "moai_with_history",
				Command:     "../noidea",
				Args:        []string{"moai", "--ai", "--history", "Follow-up fix"},
				Description: "Moai with commit history context",
				Repetitions: 3,
			},
		},
	}

	// Save the test suites
	err = SaveTestSuite(summarySuite, filepath.Join(suiteDir, "summary_tests.json"))
	if err != nil {
		return fmt.Errorf("failed to save summary test suite: %w", err)
	}

	err = SaveTestSuite(moaiSuite, filepath.Join(suiteDir, "moai_tests.json"))
	if err != nil {
		return fmt.Errorf("failed to save moai test suite: %w", err)
	}

	fmt.Println("Default test suites created in:", suiteDir)
	return nil
}

// Helper to generate and run a commit simulation with different messages
func RunCommitSimulation() error {
	// Create a temporary test repo
	testRepoDir := "test_repo"
	resultsDir := "results/commits"

	// Create dirs
	os.MkdirAll(testRepoDir, 0755)
	os.MkdirAll(resultsDir, 0755)

	// Try to load .env file from project root
	LoadEnvFile("../.env")

	// Check if we have an API key, if not use a mock one
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
	}

	// Log API key status (truncated for security)
	if apiKey != "" {
		maskedKey := "****"
		if len(apiKey) > 8 {
			maskedKey = apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
		}
		fmt.Printf("Using API key: %s\n", maskedKey)
	} else {
		fmt.Println("No API key found, using mock key for testing")
		os.Setenv("XAI_API_KEY", "mock-api-key-for-testing")
	}

	// Always enable LLM for testing
	os.Setenv("NOIDEA_LLM_ENABLED", "true")

	// Navigate to test repo dir
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	testRepoAbsPath := filepath.Join(currentDir, testRepoDir)
	fmt.Printf("Test repo path: %s\n", testRepoAbsPath)

	// Set up Git repo if it doesn't exist
	if _, err := os.Stat(filepath.Join(testRepoDir, ".git")); os.IsNotExist(err) {
		fmt.Println("Initializing new git repository...")

		// Initialize Git repo
		cmd := exec.Command("git", "init")
		cmd.Dir = testRepoDir
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to init git repo: %w", err)
		}

		// Set up Git config
		cmd = exec.Command("git", "config", "user.name", "Test User")
		cmd.Dir = testRepoDir
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("git", "config", "user.email", "test@example.com")
		cmd.Dir = testRepoDir
		err = cmd.Run()
		if err != nil {
			return err
		}

		// Create README
		err = os.WriteFile(filepath.Join(testRepoDir, "README.md"),
			[]byte("# Test Repository\n\nThis is a test repository for noidea commit simulations.\n"), 0644)
		if err != nil {
			return err
		}

		// Initial commit
		cmd = exec.Command("git", "add", "README.md")
		cmd.Dir = testRepoDir
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("git", "commit", "-m", "Initial commit")
		cmd.Dir = testRepoDir
		err = cmd.Run()
		if err != nil {
			return err
		}

		// Install noidea
		absPath, err := filepath.Abs("../noidea")
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		fmt.Printf("Using noidea binary (absolute path): %s\n", absPath)

		cmd = exec.Command(absPath, "init")
		cmd.Dir = testRepoDir
		initOutput, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install noidea: %w\nOutput: %s", err, string(initOutput))
		}
		fmt.Printf("noidea initialized in test repo: %s\n", string(initOutput))
	}

	// Test commit messages
	commitMessages := []string{
		"Add new feature X",
		"Fix bug in login flow",
		"Refactor authentication module",
		"Update documentation for API",
		"Performance optimization for database queries",
		"Remove deprecated code",
		"Merge branch 'feature/y' into main",
		"CSS styling updates",
		"Add unit tests",
		"Version bump to 1.2.0",
	}

	// Perform commits
	for i, message := range commitMessages {
		// Create a test file
		filename := fmt.Sprintf("test_file_%d.txt", i+1)
		filePath := filepath.Join(testRepoDir, filename)
		err = os.WriteFile(filePath, []byte(fmt.Sprintf("Test content for %s\n", filename)), 0644)
		if err != nil {
			return err
		}

		// Add the file
		addCmd := exec.Command("git", "add", filename)
		addCmd.Dir = testRepoDir
		err = addCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to add file %s: %w", filename, err)
		}

		// Commit with output capture
		outputFile := filepath.Join(resultsDir, fmt.Sprintf("commit_%d.txt", i+1))
		commitCmd := exec.Command("git", "commit", "-m", message)
		commitCmd.Dir = testRepoDir

		// Capture stdout and stderr
		commitOutput, err := commitCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed commit %d: %w", i+1, err)
		}

		// Save the output
		err = os.WriteFile(outputFile, commitOutput, 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("Commit %d/%d completed\n", i+1, len(commitMessages))

		// Sleep briefly to avoid rate limits on API calls
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("Commit simulation completed. Results saved to: %s\n", resultsDir)
	return nil
}

// Helper to compare results and generate a comparison report
func GenerateComparisonReport(resultsDir string, outputFile string) error {
	// Parse all test results and generate a comparison report
	report, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer report.Close()

	fmt.Fprintf(report, "# noidea Test Results Comparison\n\n")
	fmt.Fprintf(report, "Generated: %s\n\n", time.Now().Format(time.RFC3339))

	// Walk through the results directory
	err = filepath.Walk(resultsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .txt files
		if !strings.HasSuffix(info.Name(), ".txt") {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(resultsDir, path)
		if err != nil {
			return err
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Add to report
		fmt.Fprintf(report, "## %s\n\n", relPath)
		fmt.Fprintf(report, "```\n%s\n```\n\n", content)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk results directory: %w", err)
	}

	fmt.Printf("Comparison report generated: %s\n", outputFile)
	return nil
}

// LoadEnvFile loads environment variables from a .env file
func LoadEnvFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Warning: Could not read .env file: %v\n", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Skip if variable is already set in environment
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}

func main() {
	// Make sure results directory exists
	err := os.MkdirAll("results", 0755)
	if err != nil {
		fmt.Printf("Error creating results directory: %v\n", err)
		os.Exit(1)
	}

	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run simulation_tester.go [create|run|commit|compare]")
		fmt.Println("  create    - Create default test suites")
		fmt.Println("  run file  - Run a specific test suite from a JSON file")
		fmt.Println("  run all   - Run all test suites in test_suites directory")
		fmt.Println("  commit    - Run a commit simulation with various commit messages")
		fmt.Println("  compare   - Generate a comparison report of all results")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		err = CreateDefaultTestSuites()
		if err != nil {
			fmt.Printf("Error creating test suites: %v\n", err)
			os.Exit(1)
		}

	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing test suite file argument")
			os.Exit(1)
		}

		file := os.Args[2]

		if file == "all" {
			// Run all test suites
			files, err := filepath.Glob("test_suites/*.json")
			if err != nil {
				fmt.Printf("Error finding test suites: %v\n", err)
				os.Exit(1)
			}

			for _, f := range files {
				fmt.Printf("Running test suite: %s\n", f)
				suite, err := LoadTestSuite(f)
				if err != nil {
					fmt.Printf("Error loading test suite %s: %v\n", f, err)
					continue
				}

				err = RunTestSuite(suite)
				if err != nil {
					fmt.Printf("Error running test suite %s: %v\n", f, err)
				}
			}
		} else {
			// Run a specific test suite
			suite, err := LoadTestSuite(file)
			if err != nil {
				fmt.Printf("Error loading test suite: %v\n", err)
				os.Exit(1)
			}

			err = RunTestSuite(suite)
			if err != nil {
				fmt.Printf("Error running test suite: %v\n", err)
				os.Exit(1)
			}
		}

	case "commit":
		err = RunCommitSimulation()
		if err != nil {
			fmt.Printf("Error running commit simulation: %v\n", err)
			os.Exit(1)
		}

	case "compare":
		err = GenerateComparisonReport("results", "results/comparison_report.md")
		if err != nil {
			fmt.Printf("Error generating comparison report: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
