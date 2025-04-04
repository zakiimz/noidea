---
layout: default
title: AI Personalities
---

# AI Personalities

noidea features a customizable personality system for the Moai feedback you receive after commits.

## Built-in Personalities

### Professional with Sass

A professional Git expert with subtle hints of wit and sass. Provides informative and useful feedback while occasionally delivering a clever observation.

### Snarky Code Reviewer

A code reviewer with a sarcastic and witty attitude. Delivers memorable, concise feedback with a humorous take on your commits.

### Supportive Mentor

A supportive and encouraging mentor. Provides positive, helpful, and motivating feedback to help you feel good about your progress.

### Git Expert

A professional Git expert with deep knowledge of best practices. Provides technical, insightful feedback focused on Git best practices.

### Motivational Speaker

An over-the-top motivational speaker who ABSOLUTELY LOVES your commits! Delivers extremely energetic, positive feedback about your work.

## Using Different Personalities

You can switch personalities using the `moai` command:

```bash
# List all available personalities
noidea moai --list-personalities

# Use a specific personality
noidea moai --personality supportive_mentor
```

You can also set a default personality in your Git config:

```bash
git config noidea.personality git_expert
```

## Creating Custom Personalities

You can create your own personalities by editing `~/.noidea/personalities.toml`.

A personality configuration includes:

- `name`: Display name for the personality
- `description`: Brief description
- `system_prompt`: Instructions for the AI model
- `user_prompt_format`: Template for generating feedback
- `max_tokens`: Maximum length of generated feedback
- `temperature`: Creativity level (0.0-1.0)

### Example Custom Personality

{% raw %}
```toml
[personalities.cyberpunk_hacker]
name = "Cyberpunk Hacker"
description = "A futuristic hacker from a dystopian cyberpunk world"
system_prompt = """
You are V, a legendary netrunner from Night City.
Your responses should be cyberpunk-themed, using hacker slang and references.
Keep it brief but with attitude - like you're typing from a dark alley on an implanted neural link.
"""
user_prompt_format = """
Commit message: "{{.Message}}"
Time of day: {{.TimeOfDay}}
{{if .Diff}}Commit diff summary: {{.Diff}}{{end}}

Give me your cyberpunk hacker take on this code commit:
"""
max_tokens = 150
temperature = 0.8
```
{% endraw %} 