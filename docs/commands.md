---
layout: default
title: Command Reference
---

# Command Reference

This page provides a detailed reference for all available noidea commands.

## Core Commands

### init

```bash
noidea init [--force]
```

Sets up noidea in the current Git repository by installing hooks and creating necessary configurations.

**Options:**
- `--force`: Override existing hooks

### suggest

```bash
noidea suggest [--apply] [--verbose] [--no-ai]
```

Analyzes staged changes and suggests a commit message.

**Options:**
- `--apply`: Apply the suggested message directly
- `--verbose`: Generate a more detailed suggestion
- `--no-ai`: Use simple algorithm instead of AI (when API key unavailable)

### moai

```bash
noidea moai [--personality NAME] [--list-personalities]
```

Displays Moai feedback for the last commit or lists available personalities.

**Options:**
- `--personality NAME`: Use the specified personality
- `--list-personalities`: Display all available personalities

### summary

```bash
noidea summary [--days N] [--author NAME]
```

Generates a summary of recent Git activity.

**Options:**
- `--days N`: Number of days to analyze (default: 7)
- `--author NAME`: Focus on a specific author

## Configuration

### config

```bash
noidea config [--init] [--get KEY] [--set KEY VALUE]
```

View or modify noidea configuration.

**Options:**
- `--init`: Interactive configuration setup
- `--get KEY`: Get a specific configuration value
- `--set KEY VALUE`: Set a configuration value

## Utility Commands

### help

```bash
noidea help [COMMAND]
```

Shows help for all commands or a specific command.

### version

```bash
noidea version
```

Displays version information.

## Advanced Options

Most commands support the following global options:

- `--config PATH`: Use an alternative config file
- `--verbose`: Enable verbose output
- `--no-color`: Disable colorized output 