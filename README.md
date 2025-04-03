# 🧠 noidea — The Git Extension You Never Knew You Needed

> Commit with confidence. Or shame. Probably shame.

**noidea** is a lightweight, plug-and-play Git extension that adds ✨fun and occasionally useful ✨feedback into your normal Git workflow.

Every time you commit, a mysterious Moai appears to judge your code.

---

## 🗿 What It Does

After every `git commit`, you'll see something like:

```
🗿  (ಠ_ಠ) Your commit message was 'fix final final pls real'
"You've entered the 2AM hotfix arc. A legendary time."
```

Whether your code is clean or cursed, the Moai sees all.

---

## 🚀 Getting Started

1. **Install the binary**

(coming soon — cross-platform release)

For now, build from source:

```
git clone https://github.com/AccursedGalaxy/noidea.git
cd noidea
go build -o noidea
```

2. **Run `noidea init`**

```
./noidea init
```

This installs a Git `post-commit` hook in your repo.

3. **Commit your sins**

```
git commit -m "fix maybe this time"
```

And witness the Moai's judgment.

---

## 📋 Features

### Post-Commit Feedback

Get immediate feedback after each commit with the Moai:

```
🗿  (ಠ_ಠ)  This is definitely the final fix
"You've typed 'final fix' 17 times today. I'm not judging. (I am.)"
```

Options:
- `--ai` - Use AI to generate feedback (default: use config setting)
- `--diff` - Include the diff in AI context for better analysis
- `--personality <n>` - Personality to use for feedback
- `--history` - Include recent commit history for context
- `--list-personalities` - List available personalities

Example:
```bash
# Get AI-powered feedback with recent history context
noidea moai --ai --history
```

### Weekly Summaries

Generate insightful summaries of your Git activity:

```
noidea summary
```

Options:
- `--days <N>` - Analyze the last N days (default: 7)
- `--personality <n>` - Use a specific personality for insights
- `--export <format>` - Export to text, markdown, or HTML
- `--stats-only` - Show only statistics without AI insights
- `--ai` - Include AI insights (default: use config)

Example:
```bash
# Generate a 30-day summary with AI insights and export as markdown
noidea summary --days 30 --ai --export markdown
```

### On-Demand Feedback

Get targeted analysis of specific commits with powerful filtering:

```
noidea feedback
```

Options:
- `--count <N>` - Analyze last N commits (default: 5)
- `--author <n>` - Filter by commit author
- `--branch <n>` - Filter by specific branch
- `--files <list>` - Filter by files touched (comma-separated)
- `--diff` - Include diff context for deeper analysis
- `--personality <n>` - Use a specific personality
- `--export <format>` - Export to text, markdown, or HTML

Examples:

```bash
# Basic feedback on last 3 commits
noidea feedback --count 3

# Analyze commits affecting specific files
noidea feedback --files "src/main.go,pkg/utils"

# Analyze commits from a specific author
noidea feedback --author "Your Name"

# Use a supportive personality with diff context
noidea feedback --personality supportive_mentor --diff

# Export your feedback to share with the team
noidea feedback --export markdown
```

### Commit Message Suggestions

Get AI-powered commit message suggestions based on your staged changes:

```
noidea suggest
```

Options:
- `--history <N>` - Number of recent commits to analyze for context (default: 10)
- `--full-diff` - Include full diff instead of summary
- `--interactive` - Interactive mode to approve/reject suggestions
- `--file <path>` - Path to commit message file (for hooks)

> **Note:** Commit suggestions always use a professional conventional commit format, regardless of any personality settings used elsewhere in the tool.

Git Hook Integration:
- Automatically suggests commit messages during the commit process
- Easily enable with the included script:
  ```
  ./scripts/install-hooks.sh
  ```
  This installs the `prepare-commit-msg` hook and sets up your Git config with interactive prompts

Examples:

```bash
# Get a suggestion for your staged changes
noidea suggest

# Interactive mode to approve, regenerate, or edit suggestions
noidea suggest --interactive

# Consider more context from previous commits
noidea suggest --history 20

# Include the full diff for more detailed analysis
noidea suggest --full-diff
```

---

## 🧠 AI Integration


noidea supports AI-powered feedback using LLM providers that offer OpenAI-compatible APIs:

- xAI (Grok models)
- OpenAI
- DeepSeek (coming soon)

To enable AI integration:

1. Set your API key in an environment variable:
   ```
   # For xAI
   export XAI_API_KEY=your_api_key_here
   
   # For OpenAI
   export OPENAI_API_KEY=your_api_key_here
   ```

2. Or create a `.env` file in your project directory or in `~/.noidea/.env`
   ```
   XAI_API_KEY=your_api_key_here
   ```

3. Run with the `--ai` flag or enable it permanently:
   ```
   # Run with the flag
   noidea moai --ai
   
   # Enable permanently 
   export NOIDEA_LLM_ENABLED=true
   ```

4. Configure your model (optional):
   ```
   export NOIDEA_MODEL=grok-2-1212
   ```

## 🤖 AI Personalities

noidea supports multiple AI personalities to provide different types of feedback:

- **Snarky Code Reviewer** - A sarcastic, witty code reviewer (default)
- **Supportive Mentor** - Encouraging and positive feedback 
- **Git Expert** - Technical feedback based on Git best practices
- **Motivational Speaker** - Over-the-top enthusiasm for your commits!

> **Note:** Personalities affect post-commit feedback and analysis, but commit message suggestions (via `noidea suggest`) always use a professional format regardless of the selected personality.

### Using personalities

```bash
# List available personalities
noidea moai --list-personalities

# Use a specific personality
noidea moai --ai --personality=supportive_mentor

# Set default personality in config
export NOIDEA_PERSONALITY=git_expert
```

### Creating custom personalities

Create a `personalities.toml` file in your project or in `~/.noidea/` directory:

```toml
# Default personality to use
default = "my_custom_personality"

[personalities.my_custom_personality]
name = "My Custom Personality"
description = "A custom personality that fits my workflow"
system_prompt = """
Instructions for the AI on how to respond.
Keep it concise and clear.
"""
user_prompt_format = """
Commit message: "{{.Message}}"
Time of day: {{.TimeOfDay}}
{{if .Diff}}Commit diff summary: {{.Diff}}{{end}}

Provide feedback about this commit:
"""
max_tokens = 150
temperature = 0.7
```

Reference the example file at `personalities.toml.example` for more details.

---

## 🔧 Configuration

noidea can be configured using environment variables, a `.env` file, or a TOML configuration file.

### Using the config command

noidea provides a config command to help you manage your configuration:

```
# Show current configuration
noidea config --show

# Create a new configuration file interactively
noidea config --init

# Validate your configuration
noidea config --validate

# Specify a custom config path
noidea config --path /custom/path/config.toml --show
```

### Configuration file

Default location: `~/.noidea/config.toml`

Example configuration:
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

### Environment variables

You can override any configuration setting using environment variables:

```
# LLM settings
export NOIDEA_LLM_ENABLED=true
export NOIDEA_MODEL=grok-2-1212
export NOIDEA_TEMPERATURE=0.5

# Moai settings
export NOIDEA_FACES_MODE=random
export NOIDEA_USE_LINT=false
export NOIDEA_PERSONALITY=snarky_reviewer
export NOIDEA_INCLUDE_HISTORY=true

# Provider API keys
export XAI_API_KEY=your_api_key_here
export OPENAI_API_KEY=your_api_key_here
export DEEPSEEK_API_KEY=your_api_key_here
```

---

## 💡 Feature Status

| Feature                   | Status          |
|---------------------------|-----------------|
| Moai face after commit    | ✅ Done         |
| AI-based commit feedback  | ✅ Done         |
| Config file support       | ✅ Done         |
| Weekly summaries          | ✅ Done         |
| On-demand commit analysis | ✅ Done         |
| Commit message suggestions| ✅ Done         |
| Lint feedback             | 🛠️ In progress  |
| Commit streak insights    | 🔜 Coming Soon  |

---

## 🤯 Why tho?

Because Git is too serious. Coding is chaos. Let's embrace it.

---

## 🧪 Contributing

Got Moai faces? Snarky commit messages? Cursed feedback ideas?

PRs are welcome. Open an issue or just drop a meme.

---

## 🪦 Disclaimer

This tool will not improve your Git hygiene.
It will, however, make it more entertaining.

---

Made with `noidea` and late-night energy.
