<div align="center">

# 🧠 noidea

<p align="center">
  <b>Git assistant with AI commit messages and sassy feedback</b>
</p>

<p align="center">
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://golang.org/doc/go1.18"><img src="https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go" alt="Go Version"></a>
</p>

</div>

## 🗿 What is noidea?

**noidea** is a Git companion that makes your commits better and funnier:

- **Get smart commit messages** based on your changes
- **Receive sassy feedback** from a judgmental Moai after each commit
- **Analyze your Git history** for insights and patterns

<div align="center">
<img src="docs/assets/images/demo.gif" alt="noidea in action" width="80%">
</div>

## ⚡ Quick Start

```bash
# Install noidea
git clone https://github.com/AccursedGalaxy/noidea
cd noidea

./install.sh (might require sudo)

# Set up in your Git repo
cd /path/to/your/repo
noidea init

# Enable auto commit suggestions (optional)
git config noidea.suggest true

# Add your API key for AI features
echo "XAI_API_KEY=your_api_key_here" > ~/.noidea/.env
```

Now make commits as usual and enjoy both helpful suggestions and sassy feedback!

## 🧠 Features in Action

### 1. Commit Message Suggestions

When you're ready to commit, run:

```bash
noidea suggest
```

And get professional commit messages:

```
───────────────────────────────────────────────────
🧠 Analyzing staged changes...
───────────────────────────────────────────────────
✨ Suggested commit message:
feat(user-auth): implement password reset functionality with email verification
───────────────────────────────────────────────────
```

### 2. Post-Commit Feedback

After each commit, the Moai will judge your work:

```
───────────────────────────────────────────────────
🗿  (ಠ_ಠ) Your commit message was 'fix final final pls real'
"You've entered the 2AM hotfix arc. A legendary time."
───────────────────────────────────────────────────
```

### 3. Weekly Summaries

Get insights about your work:

```bash
noidea summary
```

## 🚀 Setup Options

### Installation

Choose one of these methods:

```bash
# One-line quick install
curl -sSL https://raw.githubusercontent.com/AccursedGalaxy/noidea/main/quickinstall.sh | bash

# OR Clone and install
git clone https://github.com/AccursedGalaxy/noidea.git
cd noidea
./install.sh
```

### AI Configuration

For AI-powered features, add your API key:

1. Add to environment:
   ```bash
   export XAI_API_KEY=your_api_key_here
   ```

2. Or create `~/.noidea/.env`:
   ```
   XAI_API_KEY=your_api_key_here
   ```

3. Configure interactively:
   ```bash
   noidea config --init
   ```

## 🔧 Commands

```bash
# Set up noidea in your repo
noidea init

# Get commit message suggestions
noidea suggest

# View weekly summary
noidea summary [--days 30]

# Get feedback on recent commits
noidea feedback [--count 5]

# Configure noidea
noidea config --init
```

## 🤖 AI Personalities

noidea has several AI personalities for feedback:

- **Snarky Code Reviewer** - Witty, sarcastic feedback
- **Supportive Mentor** - Encouraging, positive feedback
- **Git Expert** - Technical, professional feedback
- **Motivational Speaker** - Energetic enthusiasm

```bash
# List all personalities
noidea moai --list-personalities

# Use a specific personality
noidea moai --personality supportive_mentor
```

## 📘 More Information

<details>
<summary>Advanced Configuration</summary>

Create a `~/.noidea/config.toml` file:

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
</details>

<details>
<summary>Available Commands</summary>

- `noidea init` - Set up Git hooks in your repository
- `noidea suggest` - Get commit message suggestions
- `noidea moai` - Display Moai feedback for the last commit
- `noidea summary` - Generate summary of recent Git activity
- `noidea feedback` - Analyze specific commits
- `noidea config` - Configure noidea

Run `noidea --help` for more information.
</details>

<details>
<summary>Feature Status</summary>

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
</details>

## 🤯 Why tho?

Because Git is too serious. Coding is chaos. Let's embrace it.

This tool won't improve your Git hygiene, but it will make it more entertaining.

---

<div align="center">
Made with <code>noidea</code> and late-night energy.

<a href="https://github.com/AccursedGalaxy/noidea/stargazers"><img src="https://img.shields.io/github/stars/AccursedGalaxy/noidea?style=social" alt="GitHub stars"></a>
<a href="https://github.com/AccursedGalaxy/noidea/issues"><img src="https://img.shields.io/github/issues/AccursedGalaxy/noidea" alt="GitHub issues"></a>
</div>
