<div align="center">

# 🧠 noidea
### Smart Git Assistant with AI Commit Suggestions & Fun Feedback

<img src="docs/assets/images/banner.png" alt="noidea banner" width="600">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)](https://golang.org/doc/go1.18)
[![Made with](https://img.shields.io/badge/Made%20with-noidea-orange)](https://github.com/AccursedGalaxy/noidea)

> Intelligent commit suggestions with a side of sass. Be productive while being judged by a Moai.

[Features](#-features) • [Installation](#-getting-started) • [Usage](#-usage) • [Configuration](#-configuration) • [Documentation](https://accursedgalaxy.github.io/noidea/) • [Contributing](#-contributing)

</div>

**noidea** is a lightweight Git extension that enhances your workflow with AI-powered features:

✅ **Smart Commit Message Suggestions** - Get professional, conventional commit messages based on your staged changes  
✅ **Git Commit Analysis** - Receive insights about your commit patterns and code quality  
✅ **Fun Feedback** - Enjoy entertaining responses from a judgmental Moai after each commit  

Every `git commit` becomes both more productive and more entertaining. Whether you need help crafting the perfect commit message or just want to be roasted for your 3 AM code, noidea has you covered.

<div align="center">
<img src="docs/assets/images/demo.gif" alt="noidea in action" width="80%">
</div>

---

## 🗿 What It Does

After every `git commit`, you'll see something like:

```
───────────────────────────────────────────────────
🗿  (ಠ_ಠ) Your commit message was 'fix final final pls real'
"You've entered the 2AM hotfix arc. A legendary time."
───────────────────────────────────────────────────
```

Whether your code is clean or cursed, the Moai sees all.

And before committing, get AI-powered suggestions:

```
───────────────────────────────────────────────────
🧠 Analyzing staged changes and 10 recent commits
Generating professional commit message suggestion...
───────────────────────────────────────────────────
✨ Suggested commit message:
feat(user-auth): implement password reset functionality with email verification
───────────────────────────────────────────────────
```

<details>
<summary>👀 See noidea in action</summary>

<h3>Commit Message Suggestions</h3>
<img src="docs/assets/images/suggest.png" alt="Commit message suggestions" width="80%">

<h3>Post-Commit Moai Feedback</h3>
<img src="docs/assets/images/moai.png" alt="Moai feedback" width="80%">

<h3>Weekly Summary Analysis</h3>
<img src="docs/assets/images/summary.png" alt="Weekly summary" width="80%">
</details>

---

## 🚀 Getting Started

### Installation Options

Choose one of these methods to install noidea:

<details open>
<summary><b>1. One-Line Installation (Recommended)</b></summary>

```bash
# Install to /usr/local/bin (may require sudo)
curl -sSL https://raw.githubusercontent.com/AccursedGalaxy/noidea/main/quickinstall.sh | bash

# Or with sudo for system-wide installation
curl -sSL https://raw.githubusercontent.com/AccursedGalaxy/noidea/main/quickinstall.sh | sudo bash
```
</details>

<details>
<summary><b>2. Quick Installation Script</b></summary>

```bash
# Clone the repository
git clone https://github.com/AccursedGalaxy/noidea.git
cd noidea

# Run the installer (might need sudo)
./install.sh
# Or specify a custom location
./install.sh ~/bin
```
</details>

<details>
<summary><b>3. Using Make</b></summary>

```bash
# Clone the repository
git clone https://github.com/AccursedGalaxy/noidea.git
cd noidea

# Install to /usr/local/bin (default)
sudo make install
# Or specify a custom prefix
make install PREFIX=~/.local
```
</details>

<details>
<summary><b>4. Manual Installation</b></summary>

```bash
# Clone the repository
git clone https://github.com/AccursedGalaxy/noidea.git
cd noidea

# Build the binary
go build -o noidea

# Move it to a directory in your PATH
sudo cp noidea /usr/local/bin/
```
</details>

<details>
<summary><b>5. Pre-built Binaries (Coming Soon)</b></summary>

We'll soon provide pre-built binaries for various platforms on our releases page.
</details>

### Setting Up In Your Repository

Once noidea is installed, you can set it up in any Git repository:

```bash
# Navigate to your repository
cd /path/to/your/repo

# Initialize noidea (sets up Git hooks)
noidea init
```

This adds a post-commit hook to show the Moai judgments after each commit.

For commit message suggestions, enable them during initialization or run:

```bash
# Enable commit suggestions
git config noidea.suggest true
```

Now, when you commit, noidea will suggest a message for you!

## 📋 Features

<div align="center">
<img src="docs/assets/images/features.png" alt="noidea features overview" width="80%">
</div>

### Post-Commit Feedback

Get immediate feedback after each commit with the Moai:

```
───────────────────────────────────────────────────
🗿  (ಠ_ಠ)  This is definitely the final fix
"You've typed 'final fix' 17 times today. I'm not judging. (I am.)"
───────────────────────────────────────────────────
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

## 🧠 AI Integration

noidea supports AI-powered feedback using LLM providers that offer OpenAI-compatible APIs:

<div class="provider-grid" style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; margin-bottom: 20px;">
  <div style="text-align: center;">
    <img src="docs/assets/images/xai-logo.png" alt="xAI" width="100">
    <p>xAI (Grok)</p>
  </div>
  <div style="text-align: center;">
    <img src="docs/assets/images/openai-logo.png" alt="OpenAI" width="100">
    <p>OpenAI</p>
  </div>
  <div style="text-align: center;">
    <img src="docs/assets/images/deepseek-logo.png" alt="DeepSeek" width="100">
    <p>DeepSeek (coming soon)</p>
  </div>
</div>

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

<div class="personality-cards" style="display: flex; justify-content: space-between; margin-bottom: 20px;">
  <div style="border: 1px solid #ddd; padding: 10px; border-radius: 5px; width: 24%;">
    <h3>😈 Snarky Code Reviewer</h3>
    <p>A sarcastic, witty code reviewer that doesn't hold back</p>
  </div>
  <div style="border: 1px solid #ddd; padding: 10px; border-radius: 5px; width: 24%;">
    <h3>🤗 Supportive Mentor</h3>
    <p>Encouraging and positive feedback to keep you motivated</p>
  </div>
  <div style="border: 1px solid #ddd; padding: 10px; border-radius: 5px; width: 24%;">
    <h3>🧑‍💻 Git Expert</h3>
    <p>Technical feedback based on Git best practices</p>
  </div>
  <div style="border: 1px solid #ddd; padding: 10px; border-radius: 5px; width: 24%;">
    <h3>🚀 Motivational Speaker</h3>
    <p>Over-the-top enthusiasm for your commits!</p>
  </div>
</div>

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

## 📚 Documentation

For more detailed information, check out our [full documentation site](https://accursedgalaxy.github.io/noidea/).

<div align="center">
<a href="https://accursedgalaxy.github.io/noidea/" target="_blank">
  <img src="docs/assets/images/docs-preview.png" alt="Documentation Preview" width="80%">
</a>
</div>

We've created a comprehensive documentation site using GitHub Pages that includes:

- Detailed tutorials
- Command reference
- Configuration guide
- API documentation for integration
- Troubleshooting tips
- Advanced usage examples

## 💡 Feature Status

| Feature                   | Status          |
|---------------------------|-----------------|
| Moai face after commit    | ✅ Done         |
| AI-based commit feedback  | ✅ Done         |
| Config file support       | ✅ Done         |
| Weekly summaries          | ✅ Done         |
| On-demand commit analysis | ✅ Done         |
| Commit message suggestions| ✅ Done         |
| Enhanced terminal output  | ✅ Done         |
| POSIX-compatible hooks    | ✅ Done         |
| Lint feedback             | 🛠️ In progress  |
| Commit streak insights    | 🔜 Coming Soon  |
| Cross-platform releases   | 🔜 Coming Soon  |

## 🤯 Why tho?

Because Git is too serious. Coding is chaos. Let's embrace it.

## 🧪 Contributing

Got Moai faces? Snarky commit messages? Cursed feedback ideas?

PRs are welcome. Open an issue or just drop a meme.

Check out our test suite in the `tests/` directory to ensure your changes work as expected.

<div align="center">
<img src="docs/assets/images/contribute.png" alt="Contributing" width="50%">
</div>

## 🪦 Disclaimer

This tool will not improve your Git hygiene.
It will, however, make it more entertaining.

---

<div align="center">
Made with <code>noidea</code> and late-night energy.

<a href="https://github.com/AccursedGalaxy/noidea/stargazers"><img src="https://img.shields.io/github/stars/AccursedGalaxy/noidea?style=social" alt="GitHub stars"></a>
<a href="https://github.com/AccursedGalaxy/noidea/network/members"><img src="https://img.shields.io/github/forks/AccursedGalaxy/noidea?style=social" alt="GitHub forks"></a>
<a href="https://github.com/AccursedGalaxy/noidea/issues"><img src="https://img.shields.io/github/issues/AccursedGalaxy/noidea" alt="GitHub issues"></a>
</div>
