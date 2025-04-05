# Developer Guide

Welcome to the noidea developer documentation. This section provides information for contributors and developers who want to understand or extend noidea's functionality.

## Project Architecture

noidea is written in Go and organized into several key components:

- **cmd**: Command-line interface and entry points
- **internal**: Core functionality modules
- **scripts**: Git hooks and installation scripts

```
noidea/
├── cmd/               # Commands implementation
│   ├── root.go        # Base command and CLI setup
│   ├── suggest.go     # Commit suggestion command
│   ├── moai.go        # Feedback command
│   └── ...
├── internal/          # Internal packages
│   ├── config/        # Configuration handling
│   ├── feedback/      # Feedback generation
│   ├── git/           # Git operations
│   ├── history/       # Commit history analysis
│   ├── moai/          # Moai face and local feedback
│   ├── personality/   # AI personality system
│   └── secure/        # Secure API key storage
├── scripts/           # Installation and Git hooks
└── docs/              # Documentation
```

## Key Components

### Command Layer

The `cmd` package uses the [Cobra](https://github.com/spf13/cobra) library to implement the CLI commands. Each command is defined in its own file.

### Internal Packages

- **config**: Handles reading/writing configuration from files and environment
- **feedback**: Generates AI-powered feedback through different providers
- **git**: Manages Git operations like getting diffs and commit history
- **history**: Analyzes commit patterns and statistics
- **moai**: Manages Moai faces and local feedback generation
- **personality**: Handles the AI personality system and templates
- **secure**: Manages secure storage of API keys

## Getting Started with Development

1. **Clone the repository**:
   ```bash
   git clone https://github.com/AccursedGalaxy/noidea.git
   cd noidea
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Build the project**:
   ```bash
   go build
   ```

4. **Run tests**:
   ```bash
   go test ./...
   ```

## Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes and test them**:
   ```bash
   go build && ./noidea <command>
   ```

3. **Run linters**:
   ```bash
   # If you have golangci-lint installed
   golangci-lint run
   ```

4. **Submit a pull request**:
   - Ensure all tests pass
   - Update documentation if needed
   - Follow the code style of the project

## Documentation

Documentation is built using [MkDocs](https://www.mkdocs.org/) with the [Material](https://squidfunk.github.io/mkdocs-material/) theme.

To preview the documentation locally:

```bash
# Install mkdocs and the material theme
pip install mkdocs mkdocs-material

# Serve the documentation locally
mkdocs serve
```

## More Developer Resources

- [Architecture](architecture.md) - Detailed architecture documentation
- [Contributing](contributing.md) - How to contribute to noidea 