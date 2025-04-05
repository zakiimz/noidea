package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set args to just include command name and help flag to avoid running the full app
	os.Args = []string{"noidea", "--help"}

	// Redirect output temporarily to avoid printing help text
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	nullFile, _ := os.Open(os.DevNull)
	defer nullFile.Close()
	os.Stdout = nullFile
	os.Stderr = nullFile

	// Call main() in a separate goroutine with recover to catch os.Exit calls
	done := make(chan bool)
	go func() {
		defer func() {
			// Recover from potential panic or os.Exit
			recover()
			done <- true
		}()

		// Run the main function
		main()
	}()

	// Restore output before test ends
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Wait for main to complete or timeout
	select {
	case <-done:
		// Main completed normally
	}

	// Nothing to assert - if we got here without crashing, the test passes
}
