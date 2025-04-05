# GitHub Integration and Enhanced Release Notes

This document explains how to use the GitHub integration features of NoIdea, including the AI-enhanced release notes generator.

## Setting Up GitHub Integration

The easiest way to set up GitHub integration is to run the setup script:

```bash
./scripts/setup-github.sh
```

This script will guide you through:
1. GitHub authentication using a Personal Access Token (PAT)
2. Installing GitHub hooks for automated release creation

### Manual Setup

If you prefer to set up GitHub integration manually, follow these steps:

1. Authenticate with GitHub:
   ```bash
   noidea github auth
   ```

2. Install GitHub hooks:
   ```bash
   noidea github hook-install
   ```

## Enhanced Release Notes

NoIdea can automatically generate well-structured, user-friendly release notes when you create a new tag or release.

### How it Works

When you create a new Git tag (either manually or using `./scripts/version.sh`), NoIdea will:

1. Gather all commit messages between the previous tag and the current tag
2. If AI is enabled, process these messages using a language model to create user-friendly notes
3. Generate a structured, organized release document grouped by change type
4. Create or update a GitHub release with these enhanced notes

### Using Enhanced Release Notes

The enhanced release notes feature integrates seamlessly with your existing workflow:

1. When running `./scripts/version.sh` to bump a version, you are prompted whether to generate release notes immediately or wait for GitHub workflows
2. If you create tags manually, the post-tag hook will suggest the best command to run for release notes generation
3. You can manually generate or update release notes for any tag:
   ```bash
   noidea github release notes --tag=v1.2.3
   ```

### Workflow-Aware Release Notes

One of the most powerful features is the ability to wait for GitHub Actions workflows to complete before generating release notes:

```bash
noidea github release notes --wait-for-workflows
```

This command:
1. Checks the status of GitHub workflow runs for the specified tag
2. Displays a live progress indicator with elapsed time
3. Once all workflows complete, generates release notes that preserve GitHub's auto-generated content
4. Enhances only the overview section with AI-generated content while keeping all commit links intact

You can adjust the maximum wait time (default: 5 minutes):

```bash
noidea github release notes --wait-for-workflows --max-wait=600  # Wait up to 10 minutes
```

### AI-Powered Release Notes

If you have LLM features enabled in your NoIdea configuration, release notes will be generated using AI. This provides:

- Better organization of changes into logical sections
- More user-friendly language explaining technical changes
- Consistent formatting and style
- Focus on user impact rather than raw commit messages

To force AI-generation even if LLM is disabled in your config:

```bash
noidea github release notes --tag=v1.2.3 --ai
```

### Integration with GitHub's Release Notes

When using the `--wait-for-workflows` flag, NoIdea intelligently preserves GitHub's auto-generated content:

1. It waits for GitHub Actions workflows to complete, including any that generate release notes
2. It preserves GitHub's "What's Changed" section with commit links and PR references
3. It enhances only the "Overview" section with AI-generated content for better readability

This gives you the best of both worlds: GitHub's detailed changelog with commit links, and an AI-enhanced overview that's more user-friendly.

### Automatic Mode

For complete automation, you can use the `--auto` flag, which enables both AI generation and skips approval:

```bash
noidea github release notes --auto
```

### Examples

Standard release notes (without AI):
```markdown
# Release v1.2.3

## Changes

- Add GitHub integration
- Fix bug in config loading
- Update dependencies
```

AI-enhanced release notes:
```markdown
# Release v1.2.3

## Overview
This release adds GitHub integration capabilities, fixes several configuration bugs, and updates dependencies for improved security.

## 🚀 New Features
- **GitHub Integration**: Added complete GitHub API integration with secure token storage
- **Release Note Generation**: Automated creation of release notes from commit history

## 🛠️ Bug Fixes
- Fixed configuration loading issues when user directory contains spaces
- Resolved error handling in API key validation

## 🔧 Maintenance
- Updated all dependencies to latest versions
- Improved documentation for setup process
```

## Command Reference

| Command | Description |
|---------|-------------|
| `noidea github auth` | Authenticate with GitHub using a Personal Access Token |
| `noidea github status` | Check GitHub authentication status |
| `noidea github logout` | Remove stored GitHub credentials |
| `noidea github release create --tag=TAG` | Manually create a GitHub release |
| `noidea github release notes --tag=TAG` | Generate enhanced release notes |
| `noidea github release notes --wait-for-workflows` | Wait for GitHub Actions to complete before generating notes |
| `noidea github release notes --auto` | Automatically generate and update notes without interaction |
| `noidea github hook-install` | Install GitHub hooks for automation | 