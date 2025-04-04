---
layout: default
title: Usage Guide
---

# Usage Guide

## Getting Started

After installing noidea and setting it up in your repository, you're ready to use its features.

## Basic Workflow

### 1. Initialize in Your Repository

Before using noidea in a repository, initialize it:

```bash
cd /path/to/your/repo
noidea init
```

This sets up Git hooks that enable noidea's features.

### 2. Commit with Suggestions

When you're ready to commit, you can:

- **Get a suggestion before committing:**
  ```bash
  noidea suggest
  ```
  
- **Enable automatic suggestions during commit:**
  ```bash
  git config noidea.suggest true
  ```
  Then simply use `git commit` as usual.

### 3. Post-Commit Feedback

After each commit, the Moai will judge your work with sassy feedback.

## Command Examples

### Get Commit Message Suggestions

```bash
# View suggestions based on staged changes
noidea suggest

# Apply the suggested message directly
noidea suggest --apply

# Get a more detailed/verbose suggestion
noidea suggest --verbose
```

### View Commit History Summary

```bash
# Get a summary of recent activity
noidea summary

# Specify the number of days to analyze
noidea summary --days 30

# Focus on a specific author
noidea summary --author "your.name"
```

### Customize Moai Personality

```bash
# List available personalities
noidea moai --list-personalities

# Use a specific personality
noidea moai --personality supportive_mentor
```

## Tips & Tricks

- Use `noidea config --init` for interactive configuration
- Customize personalities in `~/.noidea/personalities.toml`
- Add your API key to `~/.noidea/.env` for AI-powered features 