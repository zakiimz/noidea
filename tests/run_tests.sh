#!/bin/bash

# Simple script to run noidea simulation tests

cd "$(dirname "$0")"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go to run the tests."
    exit 1
fi

# Make sure the noidea binary is built
if [ ! -f "../noidea" ]; then
    echo "Building noidea binary..."
    (cd .. && go build -o noidea)
    if [ $? -ne 0 ]; then
        echo "Error: Failed to build noidea binary."
        exit 1
    fi
    echo "noidea binary built successfully."
fi

# Create test repository directory if it doesn't exist
if [ ! -d "test_repo" ]; then
    echo "Creating test repository directory..."
    mkdir -p test_repo
fi

# Function to clean up test repo
cleanup_test_repo() {
    echo "Cleaning up test repository..."
    rm -rf test_repo
    mkdir -p test_repo
}

# Function to setup environment for suggest tests
setup_suggest_env() {
    echo "Setting up environment for suggest tests..."
    
    # Make sure LLM is enabled for testing
    export NOIDEA_LLM_ENABLED=true
    
    # Set a mock API key if none is provided
    if [ -z "$XAI_API_KEY" ] && [ -z "$OPENAI_API_KEY" ] && [ -z "$DEEPSEEK_API_KEY" ]; then
        echo "No API key found, using mock key for testing"
        export XAI_API_KEY="mock-api-key-for-testing"
    fi
    
    # Ensure Git has a user name and email set
    if [ -z "$(git config --global user.name)" ]; then
        git config --global user.name "Test User"
    fi
    
    if [ -z "$(git config --global user.email)" ]; then
        git config --global user.email "test@example.com"
    fi
    
    # Get absolute path to noidea
    NOIDEA_BIN=$(realpath "../noidea")
    echo "Using noidea binary: $NOIDEA_BIN"
    
    # Make the setup script executable
    chmod +x ./setup_test_repo.sh
    
    echo "Environment setup complete for suggest tests."
}

ACTION=$1

case "$ACTION" in
    create)
        echo "Creating test suites..."
        go run simulation_tester.go create
        ;;
    
    run)
        TEST_SUITE=$2
        if [ -z "$TEST_SUITE" ]; then
            echo "Error: Missing test suite argument."
            echo "Usage: $0 run [all|file.json]"
            exit 1
        fi
        
        # Clean up test repo before running tests
        cleanup_test_repo
        
        # Setup environment for suggest tests specifically
        if [[ "$TEST_SUITE" == *"suggest"* ]]; then
            setup_suggest_env
        fi
        
        echo "Running test suite: $TEST_SUITE..."
        go run simulation_tester.go run "$TEST_SUITE"
        ;;
    
    commit)
        echo "Running commit simulation..."
        go run simulation_tester.go commit
        ;;
    
    compare)
        echo "Generating comparison report..."
        go run simulation_tester.go compare
        ;;
    
    all)
        echo "=== Running complete test cycle ==="
        
        echo "1. Creating test suites..."
        go run simulation_tester.go create
        
        echo "2. Running all test suites..."
        # Clean up test repo before running all tests
        cleanup_test_repo
        
        # Setup environment for all tests including suggest
        setup_suggest_env
        
        go run simulation_tester.go run all
        
        echo "3. Running commit simulation..."
        go run simulation_tester.go commit
        
        echo "4. Generating comparison report..."
        go run simulation_tester.go compare
        
        echo "=== Testing complete ==="
        echo "Results saved in: results/"
        echo "Comparison report: results/comparison_report.md"
        ;;
    
    clean)
        echo "Cleaning up test artifacts..."
        cleanup_test_repo
        echo "Test repository cleaned."
        ;;
    
    *)
        echo "Unknown action: $ACTION"
        echo "Usage: $0 [create|run|commit|compare|all|clean]"
        echo "  create    - Create default test suites"
        echo "  run all   - Run all test suites"
        echo "  run file  - Run a specific test suite JSON file"
        echo "  commit    - Run a commit simulation"
        echo "  compare   - Generate a comparison report"
        echo "  all       - Run the complete test cycle"
        echo "  clean     - Clean up test repository"
        exit 1
        ;;
esac

exit 0 