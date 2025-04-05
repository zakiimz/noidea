---
hide:
  - navigation
---

<div class="hero">
  <div class="hero-content">
    <h1><span class="animated-emoji">🧠</span> NoIdea</h1>
    <p>Git Assistant with AI Powers & Sassy Feedback</p>
    <div class="hero-buttons">
      <a href="#installation" class="md-button md-button--primary">Get Started</a>
      <a href="https://github.com/AccursedGalaxy/noidea" class="md-button">GitHub</a>
    </div>
  </div>
</div>

<div class="content-section">
  <h2 id="what-is-noidea">🗿 What is NoIdea?</h2>

  <p>NoIdea is a Git companion that makes your commits better and funnier:</p>

  <ul>
    <li><span class="highlight-text">Get smart commit messages</span> based on your changes</li>
    <li><span class="highlight-text">Receive sassy feedback</span> from a judgmental Moai after each commit</li>
    <li><span class="highlight-text">Analyze your Git history</span> for insights and patterns</li>
  </ul>
</div>

<div class="grid cards" markdown>

- :fontawesome-brands-git-alt: __Smart Commit Messages__

    Get AI-powered commit message suggestions that accurately describe your changes, making your commit history more professional and useful

    [:octicons-arrow-right-24: Learn more](api-key-management.md)

- :material-message-text: __Sassy Feedback__

    Receive witty, personalized feedback from our judgmental Moai after each commit, with multiple AI personalities to choose from

    [:octicons-arrow-right-24: Features](#features)

- :material-chart-timeline: __Insights & Analysis__

    Track coding patterns, analyze commit quality, and get weekly summaries of your Git activity

    [:octicons-arrow-right-24: Commands](#commands)

</div>

<div class="content-section">
  <h2 id="installation">Installation</h2>

  <div class="terminal">
  ```bash
  # Clone and install
  git clone https://github.com/AccursedGalaxy/noidea.git
  cd noidea
  ./install.sh
  ```
  </div>

  <h3>Repository Setup</h3>

  <p>After installation, you need to set up NoIdea in your Git repository:</p>

  <div class="terminal">
  ```bash
  # Navigate to your repository
  cd /path/to/your/repo

  # Initialize NoIdea
  noidea init

  # Enable auto commit suggestions (optional)
  git config noidea.suggest true
  ```
  </div>
</div>

<div class="content-section">
  <h2 id="features">Features</h2>

  <h3>AI-powered Commit Suggestions</h3>

  <p>When you're ready to commit, NoIdea analyzes your changes and suggests professional commit messages:</p>

  <div class="terminal">
  ```
  ───────────────────────────────────────────────────
  🧠 Analyzing staged changes...
  ───────────────────────────────────────────────────
  ✨ Suggested commit message:
  feat(user-auth): implement password reset functionality with email verification
  ───────────────────────────────────────────────────
  ```
  </div>

  <h3>Post-commit Feedback</h3>

  <p>After each commit, the Moai will judge your work with witty commentary:</p>

  <div class="terminal">
  ```
  ───────────────────────────────────────────────────
  🗿  (ಠ_ಠ) Your commit message was 'fix final final pls real'
  "You've entered the 2AM hotfix arc. A legendary time."
  ───────────────────────────────────────────────────
  ```
  </div>

  <h3>AI Personalities</h3>

  <p>NoIdea offers several AI personalities for feedback:</p>

  <ul>
    <li><strong>Snarky Code Reviewer</strong> - Witty, sarcastic feedback</li>
    <li><strong>Supportive Mentor</strong> - Encouraging, positive feedback</li>
    <li><strong>Git Expert</strong> - Technical, professional feedback</li>
    <li><strong>Motivational Speaker</strong> - Energetic enthusiasm</li>
  </ul>

  <h3>Weekly Summaries</h3>

  <p>Track your coding patterns and get insights about your work habits with detailed summaries.</p>
</div>

<div class="content-section">
  <h2 id="commands">Commands</h2>

  <table>
    <thead>
      <tr>
        <th>Command</th>
        <th>Description</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>noidea init</code></td>
        <td>Set up Git hooks in your repository</td>
      </tr>
      <tr>
        <td><code>noidea suggest</code></td>
        <td>Get commit message suggestions</td>
      </tr>
      <tr>
        <td><code>noidea moai</code></td>
        <td>Display Moai feedback for the last commit</td>
      </tr>
      <tr>
        <td><code>noidea moai --personality supportive_mentor</code></td>
        <td>Use a specific AI personality</td>
      </tr>
      <tr>
        <td><code>noidea moai --list-personalities</code></td>
        <td>List all available personalities</td>
      </tr>
      <tr>
        <td><code>noidea summary [--days 30]</code></td>
        <td>Generate summary of recent Git activity</td>
      </tr>
      <tr>
        <td><code>noidea feedback [--count 5]</code></td>
        <td>Analyze specific commits</td>
      </tr>
      <tr>
        <td><code>noidea config</code></td>
        <td>Configure noidea</td>
      </tr>
    </tbody>
  </table>
</div>

<div class="content-section">
  <h2 id="configuration">Configuration</h2>

  <p>Configure NoIdea to match your workflow:</p>

  <div class="terminal">
  ```bash
  # Configure interactively
  noidea config --init

  # Add your API key (for AI-powered features)
  noidea config apikey
  ```
  </div>

  <h3>Advanced Configuration</h3>

  <p>For more advanced setup, you can create a <code>~/.noidea/config.toml</code> file:</p>

  <div class="terminal">
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
  </div>

  <div class="admonition tip">
    <p class="admonition-title">Pro Tip</p>
    <p>Enable auto-commit suggestions with <code>git config noidea.suggest true</code> to get suggestions every time you commit!</p>
  </div>
</div>

<div class="content-section">
  <h2 id="why-noidea">Why NoIdea?</h2>

  <p>Because Git is too serious. Coding is chaos. Let's embrace it.</p>

  <p>This tool won't improve your Git hygiene, but it will make it more entertaining.</p>

  <div class="feature-status">
  <h3>Project Status</h3>

  <table>
    <thead>
      <tr>
        <th>Feature</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>Moai face after commit</td>
        <td>✅ Done</td>
      </tr>
      <tr>
        <td>AI-based commit feedback</td>
        <td>✅ Done</td>
      </tr>
      <tr>
        <td>Config file support</td>
        <td>✅ Done</td>
      </tr>
      <tr>
        <td>Weekly summaries</td>
        <td>✅ Done</td>
      </tr>
      <tr>
        <td>On-demand commit analysis</td>
        <td>✅ Done</td>
      </tr>
      <tr>
        <td>Commit message suggestions</td>
        <td>✅ Done</td>
      </tr>
      <tr>
        <td>Enhanced terminal output</td>
        <td>✅ Done</td>
      </tr>
      <tr>
        <td>POSIX-compatible hooks</td>
        <td>✅ Done</td>
      </tr>
      <tr>
        <td>Lint feedback</td>
        <td>🛠️ In progress</td>
      </tr>
      <tr>
        <td>Commit streak insights</td>
        <td>🔜 Coming Soon</td>
      </tr>
    </tbody>
  </table>
  </div>
</div>

<div class="github-card">
  <a href="https://github.com/AccursedGalaxy/noidea" title="Star AccursedGalaxy/noidea on GitHub">
    <img src="https://img.shields.io/github/stars/AccursedGalaxy/noidea?style=social" alt="GitHub stars">
  </a>
  <a href="https://github.com/AccursedGalaxy/noidea/issues">
    <img src="https://img.shields.io/github/issues/AccursedGalaxy/noidea" alt="GitHub issues">
  </a>
</div>
