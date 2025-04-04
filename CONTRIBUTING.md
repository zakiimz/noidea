# Contributing to noidea

Thank you for your interest in contributing to noidea! This document provides guidelines and instructions for contributing to this project.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Coding Standards](#coding-standards)
- [Submitting Contributions](#submitting-contributions)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)
- [Testing](#testing)

## Code of Conduct

This project adheres to the Contributor Covenant [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork to your local machine
3. Create a new branch for your changes
4. Make your changes and commit them with clear, descriptive messages
5. Push your changes to your fork
6. Submit a pull request

## Development Environment

### Prerequisites

- Go 1.23+
- Git
- An API key from xAI, OpenAI, or DeepSeek (for testing AI features)

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/AccursedGalaxy/noidea.git
   cd noidea
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run tests:
   ```bash
   make test
   ```

### Environment Variables

For testing AI features, set up your API key:
```bash
export XAI_API_KEY=your_api_key_here
```

Or create a `.env` file in the project root:
```
XAI_API_KEY=your_api_key_here
```

## Coding Standards

### Go Style Guide

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Write meaningful comments and documentation
- Keep functions focused and small
- Use descriptive variable and function names

### Commit Messages

Follow conventional commit format:
```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Changes that don't affect the meaning of the code
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding or modifying tests
- `chore`: Changes to the build process or auxiliary tools

Example:
```
feat(suggest): add support for multi-line commit messages

This change improves the commit message suggestions to support multi-line
formats according to conventional commit standards.

Closes #42
```

## Submitting Contributions

### For Bug Fixes

1. Ensure the bug is documented in the GitHub issues
2. If the issue doesn't exist, create it
3. Reference the issue in your pull request

### For Features

1. Discuss the feature in an issue before implementing
2. Ensure the feature aligns with the project roadmap
3. Include tests and documentation

## Pull Request Process

1. Update the README.md if needed
2. Update documentation if you're changing or adding functionality
3. Add tests for new features
4. Ensure your code passes all tests
5. Make sure your code doesn't introduce linting errors
6. Reference any relevant issues in your PR description
7. Wait for review from maintainers

## Issue Reporting

When reporting issues:

1. Use a clear, descriptive title
2. Describe the expected behavior and the actual behavior
3. Include steps to reproduce the issue
4. Include version information (OS, Go version, noidea version)
5. Add screenshots if applicable

## Testing

- Write unit tests for all new features
- Run tests before submitting a PR: `make test`
- Test your changes with different API providers if modifying AI functionality

Thank you for contributing to noidea! 