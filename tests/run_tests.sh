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
        go run simulation_tester.go run all
        
        echo "3. Running commit simulation..."
        go run simulation_tester.go commit
        
        echo "4. Generating comparison report..."
        go run simulation_tester.go compare
        
        echo "=== Testing complete ==="
        echo "Results saved in: results/"
        echo "Comparison report: results/comparison_report.md"
        ;;
    
    *)
        echo "Unknown action: $ACTION"
        echo "Usage: $0 [create|run|commit|compare|all]"
        echo "  create    - Create default test suites"
        echo "  run all   - Run all test suites"
        echo "  run file  - Run a specific test suite JSON file"
        echo "  commit    - Run a commit simulation"
        echo "  compare   - Generate a comparison report"
        echo "  all       - Run the complete test cycle"
        exit 1
        ;;
esac

exit 0 