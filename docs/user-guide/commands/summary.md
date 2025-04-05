# Summary Command

The `summary` command analyzes your Git history and provides statistics and insights about your recent commits.

## Usage

```bash
noidea summary [flags]
```

## Description

The summary command provides a detailed overview of your Git activity, including:

- Total number of commits
- Lines added and removed
- Files changed
- Commit patterns by day and time
- Contribution trends
- AI-powered insights (when enabled)

By default, the command shows commits from the last 7 days. If no commits are found in this period, it automatically shows your entire repository history.

## Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--days` | `-d` | `7` | Number of days to include in summary (use 0 for all history) |
| `--all` | `-A` | `false` | Show complete repository history regardless of --days value |
| `--export` | `-e` | | Export format: text, markdown, or html |
| `--stats-only` | `-s` | `false` | Show only statistics without AI insights |
| `--ai` | `-a` | `false` | Include AI insights (default: use config setting) |
| `--personality` | `-p` | | Personality to use for insights (default: from config) |
| `--show-commits` | `-c` | `false` | Include detailed commit history in the output |

## Examples

### Basic Usage

```bash
# Show commits from the last 7 days
noidea summary

# Show commits from the last 30 days
noidea summary --days 30

# Show all repository history
noidea summary --all
# Or equivalently:
noidea summary --days 0
```

### Output Options

```bash
# Show only statistics without AI insights
noidea summary --stats-only

# Include detailed commit history in the output
noidea summary --show-commits

# Use a specific AI personality for insights
noidea summary --ai --personality supportive_mentor
```

### Exporting Results

```bash
# Export as plain text
noidea summary --export text

# Export as markdown
noidea summary --export markdown

# Export as HTML
noidea summary --export html
```

## AI Insights

When AI integration is enabled (either by default in your config or using the `--ai` flag), the summary includes AI-powered analysis of your coding patterns and provides personalized insights.

The AI insights will use the personality specified in your configuration or via the `--personality` flag.

## Related Commands

- [`moai`](moai.md) - Display feedback for your commits
- [`suggest`](suggest.md) - Generate commit message suggestions
- [`config`](config.md) - Configure noidea settings 