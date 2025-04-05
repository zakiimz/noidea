# suggest

The `suggest` command generates AI-powered commit message suggestions based on your staged changes, helping you create clear, consistent, and professional commit messages.

## Usage

```bash
noidea suggest [options]
```

## Description

This command analyzes your staged Git changes and generates a conventional commit message suggestion that follows best practices. When used with Git hooks, it can automatically pre-fill your commit message template.

## Options

| Option | Description |
|--------|-------------|
| `--history`, `-n` | Number of recent commits to analyze for context (default: 10) |
| `--full-diff`, `-f` | Include the full diff instead of a summary for better (but slower) suggestions |
| `--interactive`, `-i` | Enable interactive mode to approve/reject suggestions |
| `--file`, `-F` | Path to commit message file (for Git hooks) |
| `--quiet`, `-q` | Output only the message without UI elements (for scripts) |

## Examples

### Basic Usage

```bash
# Generate a commit message suggestion
noidea suggest
```

### With More Context

```bash
# Include more commit history for better context
noidea suggest --history 20
```

### Detailed Analysis

```bash
# Use full diff for more detailed analysis
noidea suggest --full-diff
```

### Git Integration

```bash
# Pipe directly to Git commit
noidea suggest | git commit -F-

# Use with Git hooks
# This happens automatically if you've run 'noidea init'
git config noidea.suggest true
```

## How It Works

1. **Analysis**: The command extracts your staged changes and recent commit history
2. **Context Building**: It builds context about your repository's commit style
3. **AI Processing**: The staged diff is analyzed by an AI model
4. **Suggestion**: A conventional commit message is suggested, typically following the format:
   ```
   type(scope): short description
   
   Longer description if needed
   ```

## Common Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code changes that neither fix bugs nor add features
- `test`: Adding or fixing tests
- `chore`: Maintenance tasks, dependencies, etc.

## Tips

- Stage only related changes in a single commit for better suggestions
- Use `--full-diff` when you need more detailed analysis
- For complex changes, review and edit the suggestion as needed 