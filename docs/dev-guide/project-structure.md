# Project Structure

This document outlines the organization and structure of the NoIdea codebase, helping developers understand where to find specific functionality and how components relate to each other.

## Directory Structure

```
noidea/
├── .github/               # GitHub-related files
│   ├── ISSUE_TEMPLATE/    # Issue templates
│   └── workflows/         # GitHub Actions workflows
├── .golangci-lint/        # Linting configuration
├── assets/                # Images and static assets
├── cmd/                   # CLI commands implementation
├── docs/                  # Documentation
│   ├── assets/            # Documentation assets
│   ├── blog/              # Blog posts
│   ├── dev-guide/         # Developer documentation
│   ├── stylesheets/       # Documentation CSS
│   └── user-guide/        # User documentation
├── internal/              # Internal packages
│   ├── config/            # Configuration handling
│   ├── feedback/          # Feedback generation
│   ├── git/               # Git operations
│   ├── github/            # GitHub integration
│   ├── history/           # Commit history analysis
│   ├── moai/              # Moai face and feedback
│   ├── personality/       # AI personality system
│   ├── plugin/            # Plugin system (future)
│   ├── releaseai/         # Release note generation
│   └── secure/            # Secure storage
├── scripts/               # Helper scripts and Git hooks
└── tests/                 # Test infrastructure
    ├── results/           # Test results
    ├── test_repo/         # Test Git repository
    └── test_suites/       # Test scenarios
```

## Key Components

### Entry Point

- `main.go`: The application entry point that initializes the CLI

### Command Layer

The `cmd/` directory contains all CLI commands defined using the Cobra library:

- `root.go`: Base command and shared functionality
- `suggest.go`: Commit message suggestion command
- `moai.go`: Post-commit feedback command
- `summary.go`: Git history summarization
- `config.go`: Configuration management
- `github.go`: GitHub integration
- `init.go`: Repository initialization
- `update.go`: Self-update functionality

### Internal Packages

The `internal/` directory contains packages with core functionality:

#### Configuration (`internal/config/`)

Handles loading, saving, and validating configuration from various sources:
- Environment variables
- Configuration files
- Command-line flags

#### Feedback (`internal/feedback/`)

Generates AI-powered feedback:
- Local feedback engine (no AI required)
- Unified feedback engine (works with multiple LLM providers)
- Different feedback types (commit suggestions, commit feedback, summaries)

#### Git Operations (`internal/git/`)

Abstracts Git operations:
- Getting diffs
- Retrieving commit history
- Working with branches

#### GitHub Integration (`internal/github/`)

Handles GitHub API interactions:
- Authentication
- Release management
- Release note generation
- Workflow status checks

#### History Analysis (`internal/history/`)

Analyzes Git commit patterns:
- Commit statistics
- Time-based patterns
- Author patterns
- Commit message analysis

#### Moai (`internal/moai/`)

Manages the Moai face system:
- ASCII Moai faces
- Local feedback templates
- Face selection based on context

#### Personality (`internal/personality/`)

Implements the AI personality system:
- Loading personality definitions
- Managing context and prompts
- Personality selection

#### Security (`internal/secure/`)

Provides secure storage for sensitive information:
- API key management
- Secure credential storage
- Environment variable handling

### Scripts

The `scripts/` directory contains helper scripts:
- Git hooks
- Installation scripts
- Release management
- Document generation

## Code Flow

### Commit Suggestion Flow

1. User runs `noidea suggest`
2. Command parsed in `cmd/suggest.go`
3. Git changes retrieved via `internal/git`
4. Configuration loaded from `internal/config`
5. Feedback engine created in `internal/feedback`
6. AI provider selected based on config
7. Suggestion generated and displayed to user

### Moai Feedback Flow

1. User makes a Git commit
2. Post-commit hook runs `noidea moai`
3. Command handled in `cmd/moai.go`
4. Last commit retrieved via `internal/git`
5. Moai face selected from `internal/moai`
6. Feedback generated via `internal/feedback`
7. Feedback displayed to user

## Plugin Architecture

The future plugin system will be centered around the `internal/plugin` package:

- Plugin registry for managing plugins
- Plugin loader for loading different plugin types
- Core interfaces for plugin developers
- Event system for plugin hooks

Plugins will integrate with NoIdea through several hook points:
- Command hooks
- Pre/post commit hooks
- Feedback hooks
- UI hooks
- Data hooks

## Build System

NoIdea uses a `Makefile` for build automation:
- `make build`: Builds the binary
- `make install`: Installs NoIdea
- `make test`: Runs tests
- `make lint`: Runs linters
- `make release`: Creates release builds

## Documentation

Documentation is built using MkDocs with the Material theme:
- User guide for end users
- Developer guide for contributors
- Blog posts for announcements and tips

## Future Structure

As NoIdea evolves, the following structural changes are planned:

1. Add the `internal/plugin` directory for the plugin system
2. Add a dedicated `api` package for stable API interfaces
3. Implement a modular LLM provider system
4. Enhance the GitHub integration capabilities 