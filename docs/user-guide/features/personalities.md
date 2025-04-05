# AI Personalities

noidea's personality system allows you to customize the tone and style of AI-generated feedback, making your Git experience more enjoyable and personalized.

## Built-in Personalities

noidea comes with several built-in personalities:

| Personality | Description | Best For |
|-------------|-------------|----------|
| `professional_sass` | Professional with subtle hints of sass | Day-to-day work |
| `snarky_reviewer` | Witty, sarcastic code reviewer | Fun environments |
| `supportive_mentor` | Encouraging and helpful | Learning/teaching |
| `git_expert` | Technical Git best practices | Professional teams |
| `motivational_speaker` | Energetic and enthusiastic | Motivation boosts |

## Using Personalities

You can select a personality in several ways:

### Command Line

```bash
# Use with moai command
noidea moai --ai --personality supportive_mentor

# List all available personalities
noidea moai --list-personalities
```

### Configuration File

In your `~/.noidea/config.json`:

```json
{
  "moai": {
    "personality": "snarky_reviewer"
  }
}
```

### Environment Variable

```bash
export NOIDEA_PERSONALITY="git_expert"
```

## Creating Custom Personalities

You can create your own personalities by creating a `personalities.toml` file in your `~/.noidea/` directory:

```toml
[personalities.my_custom_personality]
name = "My Custom Personality"
description = "A personality that matches my team's culture"
system_prompt = """
You are a Git expert with my company's specific style.
Your responses should be concise and follow our team conventions.
Highlight performance improvements and security considerations.
Keep your responses under 100 characters.
"""
user_prompt_format = """
Commit message: "{{.Message}}"
Time of day: {{.TimeOfDay}}
{{if .Diff}}Commit diff summary: {{.Diff}}{{end}}

Provide feedback that matches our team culture:
"""
max_tokens = 150
temperature = 0.5
```

### Personality Template Variables

Your custom personalities can use these template variables:

| Variable | Description |
|----------|-------------|
| `{{.Message}}` | The commit message |
| `{{.TimeOfDay}}` | Current time (e.g., "morning", "afternoon") |
| `{{.Diff}}` | Commit diff (if included) |
| `{{.Username}}` | Git username |
| `{{.RepoName}}` | Repository name |
| `{{.CommitHistory}}` | Recent commit messages |

### Personality Settings

| Setting | Description | Default |
|---------|-------------|---------|
| `name` | Display name for the personality | Required |
| `description` | Brief description | Required |
| `system_prompt` | Instructions for the AI | Required |
| `user_prompt_format` | Template for the user prompt | Required |
| `max_tokens` | Maximum response length | 150 |
| `temperature` | Randomness (0.0-1.0) | 0.7 |

## Setting a Default Personality

To set the default personality in your configuration:

```json
{
  "moai": {
    "personality": "my_custom_personality",
    "personality_file": "~/.noidea/personalities.toml"
  }
}
```

## Tips for Creating Personalities

- Keep system prompts concise and specific
- For consistent results, use lower temperature values (0.2-0.5)
- For more creative results, use higher values (0.7-0.9)
- Test your personality with different types of commits
- Include specific guidelines about response formatting and length 