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

# Provider API keys
export XAI_API_KEY=your_api_key_here
export OPENAI_API_KEY=your_api_key_here
export DEEPSEEK_API_KEY=your_api_key_here
```

---

## 💡 Upcoming Features

| Feature                   | Status          |
|---------------------------|-----------------|
| Moai face after commit    | ✅ Done          |
| AI-based commit feedback  | ✅ Done          |
| Config file support       | ✅ Done          |
| Lint feedback             | 🛠️ In progress   |
| Commit streak insights    | 🔜 Coming Soon   |

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
