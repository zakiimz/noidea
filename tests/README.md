# Noidea Testing Suite

This directory contains automated tests for the noidea tool to ensure it functions correctly across different environments and use cases.

## Test Structure

The testing suite is organized as follows:

```
tests/
  ├── README.md             # This documentation file
  ├── run_tests.sh          # Main script to execute all tests
  ├── simulation_tester.go  # Go program that simulates user interactions
  ├── test_suites/          # Collection of individual test scenarios
  │   ├── basic/            # Basic functionality tests
  │   ├── advanced/         # Advanced feature tests
  │   └── edge_cases/       # Tests for unusual scenarios
  ├── test_repo/            # Git repository used for testing
  └── results/              # Test output and logs
```

## Running Tests

To run the full test suite:

```bash
# From the project root directory:
./tests/run_tests.sh

# To run a specific test suite:
./tests/run_tests.sh basic

# To run with verbose output:
./tests/run_tests.sh --verbose
```

## Test Types

### Basic Tests

- **Installation Tests**: Verify that noidea can be installed correctly
- **Hook Setup Tests**: Ensure Git hooks are properly installed and configured
- **Command Tests**: Test basic command functionality (help, version, etc.)

### Feature Tests

- **Commit Message Suggestion**: Tests for commit message generation
  - Single-line suggestions for simple changes
  - Multi-line commit messages with body for complex changes
  - Different command flags (history context, full diff, etc.)
  - Interactive mode and file output options
  - Conventional commit format verification
  - Breaking change detection
- **Moai Feedback**: Verify that post-commit feedback works correctly
- **Configuration**: Test configuration loading and saving
- **Summary Reports**: Test the generation of summary reports

### Edge Cases

- **No Git Repo**: Behavior when not in a Git repository
- **Empty Commits**: Handling of commits with no changes
- **Large Diffs**: Performance with large file changes
- **Unicode Support**: Handling of non-ASCII characters in commit messages

## Writing New Tests

To add a new test:

1. Create a new directory in `test_suites/` with a descriptive name
2. Add a `setup.sh` script that prepares the test environment
3. Add a `test.sh` script that runs your test steps
4. Add a `cleanup.sh` script to restore the environment
5. Add a `README.md` describing the test purpose and expected results

## Test Output

Test results are stored in the `results` directory with the following structure:

- `YYYY-MM-DD_HH-MM-SS/` - Timestamp directory for the test run
  - `summary.log` - Overview of all test results
  - `test_name/` - Directory for each test
    - `output.log` - Standard output and error logs
    - `results.json` - Machine-readable results

## Common Issues

- **Permission Errors**: Ensure `run_tests.sh` is executable (`chmod +x run_tests.sh`)
- **Git Config**: Tests may modify Git configuration; use a test repository!
- **Dependencies**: Some tests require specific Go packages or system tools 