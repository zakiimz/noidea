# noidea Simulation Testing Framework

This directory contains a testing framework for simulating and comparing outputs from noidea's AI-driven features. Since traditional unit testing doesn't work well for evaluating AI-based responses, this framework runs the same commands multiple times with different parameters to observe patterns and inconsistencies.

## Overview

The simulation tester can:

1. Run noidea commands with various configurations
2. Repeat each test multiple times to observe variations in AI responses
3. Save all outputs to organized directories for comparison
4. Generate commit simulations to test the Git hook functionality
5. Create comparison reports for manual review

## Getting Started

First, make sure the noidea binary is built:

```bash
# From the project root
go build -o noidea
```

### Using the Test Script

For convenience, a shell script is provided to run the tests:

```bash
# Make the script executable
chmod +x tests/run_tests.sh

# Create default test suites
./tests/run_tests.sh create

# Run all tests
./tests/run_tests.sh all

# Run a specific test suite
./tests/run_tests.sh run test_suites/moai_tests.json

# Run only the commit simulation
./tests/run_tests.sh commit

# Generate a comparison report
./tests/run_tests.sh compare
```

### Manual Usage

You can also run the simulation tester directly:

```bash
# Navigate to tests directory
cd tests

# Create test suites
go run simulation_tester.go create

# Run all test suites
go run simulation_tester.go run all

# Run a specific test suite
go run simulation_tester.go run test_suites/summary_tests.json

# Run commit simulation
go run simulation_tester.go commit

# Generate comparison report
go run simulation_tester.go compare
```

## Test Suites

Test suites are defined in JSON files in the `test_suites` directory:

- `summary_tests.json` - Tests for the `summary` command with different configurations
- `moai_tests.json` - Tests for the `moai` command with different personalities and flags

### Customizing Test Suites

You can modify the existing test suites or create new ones. The test suite format is:

```json
{
  "name": "Test Suite Name",
  "output_dir": "results/directory",
  "test_cases": [
    {
      "name": "test_case_name",
      "command": "./noidea",
      "args": ["summary", "--days", "14"],
      "description": "Description of the test case",
      "repetitions": 3
    },
    ...
  ]
}
```

## Understanding Test Results

After running tests, results are organized in the `results` directory:

```
results/
  ├── summary/
  │   ├── default_summary/
  │   │   ├── info.txt
  │   │   ├── run_1.txt
  │   │   ├── run_2.txt
  │   │   └── run_3.txt
  │   └── ...
  ├── moai/
  │   ├── default_moai/
  │   │   ├── info.txt
  │   │   ├── run_1.txt
  │   │   ...
  │   └── ...
  ├── commits/
  │   ├── commit_1.txt
  │   ├── commit_2.txt
  │   └── ...
  └── comparison_report.md
```

The comparison report (`results/comparison_report.md`) contains a consolidated view of all test results for easy review.

## What to Look For

When reviewing test results, pay attention to:

1. **Consistency** - Are responses for the same input reasonably consistent?
2. **Quality** - Does the AI provide useful and relevant feedback?
3. **Personality** - Do different personalities produce noticeably different styles?
4. **Error Handling** - Are errors handled gracefully?
5. **Execution Time** - Are commands responding within a reasonable time?

## Adding New Tests

To add new test types:

1. Modify `CreateDefaultTestSuites()` in `simulation_tester.go`
2. Add your new test cases with appropriate command arguments
3. Run `./run_tests.sh create` to generate the updated test suites 