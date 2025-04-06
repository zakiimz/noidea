# Testing Guide

This document provides guidelines and best practices for testing NoIdea components and plugins.

## Overview

Testing is a critical part of ensuring NoIdea's reliability and quality. The project uses a combination of:

- Unit tests
- Integration tests
- Simulation tests
- Linting and static analysis

## Test Directory Structure

The NoIdea testing infrastructure is organized as follows:

```
tests/
├── test_suites/          # Test scenarios and configuration
├── results/              # Test outputs and reports
├── test_repo/            # Test Git repository
├── simulation_tester.go  # Automated test runner
├── run_tests.sh          # Test executor script
└── setup_test_repo.sh    # Repository setup script
```

## Unit Testing

Unit tests should be added for all new functionality. Follow these guidelines for unit tests:

### Test File Naming

- Test files should be named `*_test.go`
- Place test files in the same package as the code being tested

### Test Function Naming

Use descriptive names for test functions following the convention:

```go
func TestFeature_Scenario(t *testing.T) {
    // Test implementation
}
```

For example:

```go
func TestGetRandomFace_ReturnsValidFace(t *testing.T) {
    // Test implementation
}
```

### Table-Driven Tests

Use table-driven tests for testing multiple inputs and expected outputs:

```go
func TestProcessFeedback(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "basic feedback",
            input:    "test message",
            expected: "processed test message",
        },
        {
            name:     "empty input",
            input:    "",
            expected: "",
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ProcessFeedback(tt.input)
            if result != tt.expected {
                t.Errorf("ProcessFeedback(%q) = %q, want %q", tt.input, result, tt.expected)
            }
        })
    }
}
```

### Mocking

Use mocks for external dependencies like:

- Git operations
- API calls
- File system operations

You can create mocks manually or use a mocking library:

```go
// MockGitClient implements the GitClient interface for testing
type MockGitClient struct {
    CommitMessages []string
    ShouldError    bool
}

func (m *MockGitClient) GetLastCommitMessage() (string, error) {
    if m.ShouldError {
        return "", errors.New("mock error")
    }
    if len(m.CommitMessages) > 0 {
        return m.CommitMessages[len(m.CommitMessages)-1], nil
    }
    return "", nil
}
```

## Integration Tests

Integration tests verify that different components work together correctly. For NoIdea, these typically involve:

1. Setting up a test Git repository
2. Running commands on that repository
3. Verifying the outcomes

The `tests/test_repo/` directory contains a Git repository specifically for integration testing.

### Running Integration Tests

Use the provided script to run integration tests:

```bash
./tests/run_tests.sh integration
```

## Simulation Tests

NoIdea uses simulation tests to verify end-to-end functionality. The `simulation_tester.go` file contains a test runner that:

1. Sets up a clean test environment
2. Executes a series of Git and NoIdea commands
3. Verifies the outputs match expected results

### Test Suites

Test suites are defined in JSON files in the `tests/test_suites/` directory:

```json
{
  "name": "Basic Commit Flow",
  "description": "Tests the basic commit suggestion and feedback flow",
  "steps": [
    {
      "command": "git init",
      "expected_output": "Initialized"
    },
    {
      "command": "touch test.txt",
      "expected_output": ""
    },
    {
      "command": "git add test.txt",
      "expected_output": ""
    },
    {
      "command": "noidea suggest",
      "expected_output_pattern": "Suggested commit message:"
    }
    // More steps...
  ]
}
```

### Running Simulation Tests

To run simulation tests:

```bash
./tests/run_tests.sh simulation
```

## Plugin Testing

When developing plugins for NoIdea, follow these additional guidelines:

### Plugin Unit Tests

Each plugin should have unit tests that verify:

1. The plugin can be loaded successfully
2. Core functionality works as expected
3. The plugin handles errors gracefully
4. The plugin cleans up resources properly

### Plugin Integration Tests

Create integration tests that verify your plugin works correctly with NoIdea:

```go
func TestMyPlugin_Integration(t *testing.T) {
    // Setup
    tmpDir := t.TempDir()
    gitInit := exec.Command("git", "init", tmpDir)
    gitInit.Run()
    
    // Run NoIdea with plugin
    cmd := exec.Command("noidea", "--plugin", "my-plugin", "command")
    cmd.Dir = tmpDir
    output, err := cmd.CombinedOutput()
    
    // Verify results
    if err != nil {
        t.Errorf("Command failed: %v", err)
    }
    if !strings.Contains(string(output), "Expected output") {
        t.Errorf("Output did not contain expected string")
    }
}
```

## Test Coverage

Aim for high test coverage, especially for critical components. You can check the current test coverage with:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Continuous Integration

NoIdea uses GitHub Actions for continuous integration. The CI workflow:

1. Runs all unit tests
2. Verifies code formatting
3. Runs linters
4. Builds the project for multiple platforms

See the `.github/workflows/` directory for the complete CI configuration.

## Debugging Tests

When debugging failing tests:

1. Use the `-v` flag for verbose output: `go test -v ./...`
2. Use `t.Logf()` to print debug information
3. For simulation tests, check the `tests/results/` directory for detailed logs

## Adding New Tests

When adding new features:

1. Add unit tests for all new functions
2. Update or add integration tests if the feature changes behavior
3. Consider adding a new simulation test suite for significant features

Following these guidelines will help ensure NoIdea remains stable and reliable as it evolves. 