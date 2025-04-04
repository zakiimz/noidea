---
layout: default
title: Configuration Guide
---

# Configuration Guide

noidea provides several ways to configure its behavior to suit your workflow.

## Configuration File

The main configuration is stored in `~/.noidea/config.json`:

```json
{
  "llm": {
    "enabled": true,
    "provider": "xai",
    "api_key": "your_api_key_here",
    "model": "grok-2-1212",
    "temperature": 0.7
  },
  "moai": {
    "use_lint": false,
    "faces_mode": "random",
    "personality": "snarky_reviewer",
    "personality_file": "/home/user/.noidea/personalities.toml"
  }
}
```

## Interactive Configuration

The easiest way to configure noidea is with the interactive configuration tool:

```bash
noidea config --init
```

This will guide you through setting up your configuration.

## AI Configuration

### API Providers

noidea supports multiple AI providers:

- **xAI (Default)**: Uses Grok models
- **OpenAI**: Uses GPT models (3.5-turbo, 4, etc.)
- **DeepSeek**: Uses DeepSeek models (experimental)

### API Keys

Set your API key using one of these methods:

1. In your `~/.noidea/config.json` file
2. In a `~/.noidea/.env` file: `XAI_API_KEY=your_key_here`
3. As an environment variable: `export XAI_API_KEY=your_key_here`

## Personality Configuration

Personality configurations are stored in `~/.noidea/personalities.toml`:

```toml
# Default personality to use
default = "snarky_reviewer"

# Personality definitions
[personalities]

[personalities.snarky_reviewer]
name = "Snarky Code Reviewer"
description = "A code reviewer with a sarcastic and witty attitude"
system_prompt = """
You are a snarky but insightful Git expert named Moai. 
Your responses should be witty, memorable, and concise.
Always aim to be funny while also providing insight about the commit.
Keep your responses between 50-120 characters and as a single sentence.
"""
# ... more configuration
```

You can create your own personalities by adding them to this file.

## Git Repository Configuration

noidea can be configured per-repository using Git config:

```bash
# Enable automatic commit suggestions
git config noidea.suggest true

# Set the default personality for this repository
git config noidea.personality supportive_mentor

# Disable Moai feedback for this repository
git config noidea.moai false
``` 