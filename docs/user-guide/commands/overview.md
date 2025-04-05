# Command Reference

noidea offers several commands to enhance your Git workflow. This page provides an overview of all available commands.

## Available Commands

| Command | Description |
|---------|-------------|
| `init` | Set up noidea in your Git repository |
| `suggest` | Generate commit message suggestions based on staged changes |
| `moai` | Display feedback about your most recent commit |
| `summary` | Generate a summary of your recent Git activity |
| `config` | Manage noidea configuration |

## Getting Help

For any command, you can use the `--help` flag to see available options:

```bash
noidea <command> --help
```

## Common Options

These options are available for most commands:

| Option | Description |
|--------|-------------|
| `--version`, `-v` | Show version information |
| `--help`, `-h` | Show help for a command |

## Detailed Command Documentation

Each command has its own detailed documentation page:

- [`init`](init.md) - Setup noidea in your repository
- [`suggest`](suggest.md) - Generate commit message suggestions
- [`moai`](moai.md) - Get feedback on your commits
- [`summary`](summary.md) - Analyze your Git history
- [`config`](config.md) - Configure noidea

## Examples

Here are some common usage examples:

### Setting up noidea

```bash
# Initialize noidea in your repository
noidea init
```

### Getting commit suggestions

```bash
# Get a commit message suggestion
noidea suggest

# Pipe the suggestion directly to commit
noidea suggest | git commit -F-
```

### Getting Moai feedback

```bash
# Get feedback on your last commit
noidea moai

# Use AI-powered feedback with a specific personality
noidea moai --ai --personality supportive_mentor
```

### Generating summaries

```bash
# Get a summary of the last 30 days
noidea summary --days 30

# Get a weekly summary
noidea summary --weeks 1
``` 