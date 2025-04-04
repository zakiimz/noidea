---
layout: default
title: Configuration
nav_order: 5
permalink: /docs/configuration
---

# Advanced Configuration
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Configuration File

noidea uses a TOML configuration file located at `~/.noidea/config.toml`. You can create or edit this file manually, or use the interactive configuration tool:

```bash
noidea config --init
```

Here's an example of a complete configuration file:

```toml
[llm]
enabled = true
provider = "xai"
api_key = "your_api_key_here"
model = "grok-2-1212"
temperature = 0.7

[moai]
use_lint = false
faces_mode = "random"
personality = "snarky_reviewer"
include_history = true
```

## LLM Configuration

The `[llm]` section controls AI-related settings:

| Option | Description | Default |
|--------|-------------|---------|
| `enabled` | Enable or disable AI features | `true` |
| `provider` | AI provider to use | `"xai"` |
| `api_key` | Your API key | `""` |
| `model` | Model to use | `"grok-2-1212"` |
| `temperature` | Creativity level (0.0-1.0) | `0.7` |

## Moai Configuration

The `[moai]` section controls the post-commit feedback:

| Option | Description | Default |
|--------|-------------|---------|
| `use_lint` | Enable linting feedback | `false` |
| `faces_mode` | How to select Moai faces (`"random"`, `"mood"`) | `"random"` |
| `personality` | Default personality to use | `"snarky_reviewer"` |
| `include_history` | Include commit history in analysis | `true` |

## Personalities

You can customize AI personalities by creating a `~/.noidea/personalities.toml` file. Here's an example:

```toml
[snarky_reviewer]
name = "Snarky Code Reviewer"
description = "A witty, slightly sarcastic reviewer who doesn't hold back"
prompt = """
You are a witty, slightly sarcastic code reviewer who provides feedback on Git commits.
Your tone is humorous but insightful. Focus on both technical accuracy and humor.
"""

[supportive_mentor]
name = "Supportive Mentor"
description = "A positive and encouraging mentor"
prompt = """
You are a supportive and encouraging mentor who provides positive feedback on Git commits.
Your tone is warm and constructive, always finding something to praise while gently suggesting improvements.
"""
```

## Git Configuration

You can configure Git-specific settings with:

```bash
# Enable automatic commit suggestions
git config noidea.suggest true

# Change the default Moai personality
git config noidea.personality supportive_mentor

# Set your preferred AI model
git config noidea.model grok-2-1212
``` 