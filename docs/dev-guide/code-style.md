# Code Style Guide

This document outlines the coding standards and style guidelines for the NoIdea project. Following these guidelines ensures consistency across the codebase and makes collaboration easier.

## General Principles

- **Readability**: Code should be easy to read and understand
- **Simplicity**: Prefer simple solutions over complex ones
- **Consistency**: Follow established patterns in the codebase
- **Documentation**: Document code thoroughly, especially public APIs

## Go Style Guide

NoIdea follows the official Go style guidelines with some project-specific additions.

### Formatting

- Use `gofmt` or `goimports` to format your code
- Line length should be kept reasonable (aim for under 100 characters)
- Group imports in the standard Go way:
  1. Standard library imports
  2. Third-party imports
  3. Local project imports

```go
import (
    "fmt"
    "strings"
    
    "github.com/fatih/color"
    "github.com/spf13/cobra"
    
    "github.com/AccursedGalaxy/noidea/internal/config"
)
```

### Naming Conventions

- Use meaningful, descriptive names
- Follow Go conventions:
  - `CamelCase` for exported names
  - `mixedCase` for non-exported names
  - `acronyms` should be all caps: `HTTPClient` not `HttpClient`
- Package names should be short, concise, and lowercase

### Error Handling

- Always check errors and handle them appropriately
- Use descriptive error messages
- Consider using custom error types for domain-specific errors
- Wrap errors with context when propagating them up the call stack

```go
if err != nil {
    return fmt.Errorf("failed to load configuration: %w", err)
}
```

### Comments

- Write comments for all exported functions, types, and constants
- Follow the Go convention for documentation comments:

```go
// GetRandomFace returns a random Moai face string.
// It selects from a predefined list of ASCII Moai faces.
func GetRandomFace() string {
    // Implementation
}
```

### Project-Specific Conventions

#### Command Structure

All CLI commands should:

1. Have a descriptive `Use` field
2. Include a short and long description
3. Handle errors consistently
4. Follow the established command pattern

```go
var exampleCmd = &cobra.Command{
    Use:   "example [options]",
    Short: "Short description of command",
    Long:  `A longer description that explains the command in detail.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Command implementation
        if err := doSomething(); err != nil {
            fmt.Println(color.RedString("Error:"), err)
            return
        }
    },
}
```

#### Configuration Handling

- Configuration should be accessed through the `config` package
- Don't use hardcoded configuration values
- Respect user-defined settings

#### Feedback and UI

- Use the `color` package consistently for terminal output
- Follow the established color scheme:
  - Red for errors
  - Yellow for warnings
  - Green for success
  - Cyan for information
  - White for normal output

## Code Organization

### File Structure

- Place related functionality in the same package
- Break large files into smaller, focused ones
- Keep the `main` package minimal, delegating to other packages

### Package Organization

NoIdea uses the following package organization:

- `cmd/`: CLI commands
- `internal/`: Internal packages (not exported)
  - `config/`: Configuration handling
  - `feedback/`: Feedback generation
  - `git/`: Git operations
  - `moai/`: Moai face and local feedback
  - `plugin/`: Plugin system
  - ...
- `scripts/`: Helper scripts and Git hooks

## Testing

- Write tests for all new functionality
- Use table-driven tests when appropriate
- Mock external dependencies
- Aim for high test coverage, especially for critical components

```go
func TestGetRandomFace(t *testing.T) {
    face := GetRandomFace()
    if face == "" {
        t.Error("Expected non-empty face, got empty string")
    }
    
    // Check if face is in the valid faces list
    validFace := false
    for _, f := range moaiFaces {
        if f == face {
            validFace = true
            break
        }
    }
    
    if !validFace {
        t.Errorf("Got unexpected face: %s", face)
    }
}
```

## Linting and Static Analysis

NoIdea uses the following tools for code quality:

- `golangci-lint` with custom configuration in `.golangci.yml`
- Git hooks that run linters automatically

Make sure to run linters before submitting code:

```bash
make lint
```

## Plugin Development

When developing plugins, follow these additional guidelines:

- Use interfaces defined in the `plugin` package
- Follow the plugin lifecycle correctly
- Provide clear error messages
- Include thorough documentation
- Respect resource limitations

See the [Plugin System](plugins/index.md) documentation for more details. 