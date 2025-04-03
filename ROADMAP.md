## 🧠 `noidea` — Development Roadmap
**Goal:** A CLI Git extension that provides fun, insightful, AI-powered feedback after each commit — with no extra work for the user.

---

### 🏗️ Phase 1: Core CLI + Moai Hook ✅

**🔹 Goal:** Fully working post-commit Moai display.

#### ✅ Tasks:
- [x] Set up Go module + project structure (`main.go`, `cmd/`, etc.)
- [x] Use [`spf13/cobra`](https://github.com/spf13/cobra) for CLI:
  - `noidea init` → installs Git `post-commit` hook
  - `noidea moai` → renders commit-based feedback
- [x] Create `scripts/githook-post-commit.sh` template
- [x] Generate the Git hook dynamically and place it via `init`
- [x] Create a few Moai faces and feedback messages (ASCII + random text)
- [x] Implement colored terminal output using [`fatih/color`](https://github.com/fatih/color)

---

### ⚙️ Phase 2: LLM Integration via OpenAI-compatible Framework ✅

**🔹 Goal:** Use AI to give fun, context-aware Git commit feedback.

#### ✅ Tasks:
- [x] Define `FeedbackEngine` interface (abstraction for LLM agent)
- [x] Implement OpenAI-compatible backend using `openai-go` or `ollama`, `openfga`, etc.
  - [x] Create unified engine for all OpenAI-compatible APIs (xAI, OpenAI, DeepSeek)
  - [x] Define provider configurations for easy switching
- [x] Capture basic commit context:
  - Latest commit message
  - Timestamp
  - Commit diff (optional)
- [x] Craft base prompt:
  _"You are a snarky but insightful Git expert. Given the following commit message and time of day, give a short and funny, but helpful comment."_
- [x] Add `--ai` flag to `noidea moai` or auto-toggle via config
- [x] Fallback to local jokes if no API key is provided

---

### 🧪 Phase 3: LLM Config & User Settings ✅

**🔹 Goal:** Let the user configure LLM behavior.

#### ✅ Tasks:
- [x] Support multiple LLM providers (xAI, OpenAI, DeepSeek)
- [x] Support environment variables for configuration
- [x] Create a config file:
  - Location: `~/.noidea/config.toml`
- [x] Config structure:
  ```toml
  [llm]
  enabled = true
  provider = "openai"
  api_key = "sk-..."
  model = "gpt-3.5-turbo"
  temperature = 0.7

  [moai]
  use_lint = true
  faces_mode = "random"
  ```
- [x] CLI command: `noidea config` (optional for manual key entry)
- [x] Validate config + give feedback on missing keys
- [x] Allow overriding config with env vars (e.g. `NOIDEA_API_KEY`)

---

### 📊 Phase 4: Personality & Customization ✅

**🔹 Goal:** Enhance the AI feedback system with customizable personalities and improved prompts.

#### ✅ Tasks:
- [x] Enhance LLM Prompt System:
  - [x] Design modular prompt template system
  - [x] Add context-aware prompt generation
  - [x] Implement prompt validation and testing
- [x] Personality Configuration:
  - [x] Create TOML-based personality configuration schema
  - [x] Add support for custom prompt templates
  - [x] Implement personality hot-reloading
- [x] Default Personalities:
  - [x] Implement "Snarky Code Reviewer" personality
  - [x] Implement "Supportive Mentor" personality  
  - [x] Implement "Git Expert" personality
- [x] Documentation:
  - [x] Add personality configuration guide
  - [x] Document prompt template syntax
  - [x] Include personality customization examples

---

### 📈 Phase 5: AI-Driven Commit Insights

**🔹 Goal:** Provide actionable and interesting insights into commit history patterns using LLM analysis.

#### ✅ Tasks:
- [ ] Git History Collection:
  - [ ] Create `HistoryCollector` module in `internal/history/`
  - [ ] Implement flexible date range filtering (past 7 days, past X commits)
  - [ ] Extract commit metadata (messages, authors, timestamps, files changed)
  - [ ] Add diff summarization for contextual understanding
  - [ ] Implement caching to reduce repeat Git operations
- [ ] Weekly Summary Feature:
  - [ ] Add `noidea summary` command with time range options
  - [ ] Design comprehensive prompt for weekly work analysis
  - [ ] Create statistical overview (commit frequency, time patterns, etc.)
  - [ ] Generate human-readable summary with actionable insights
  - [ ] Support markdown/HTML export for team sharing
- [ ] On-Demand Feedback:
  - [ ] Implement `noidea feedback [--count X]` command
  - [ ] Add filtering options (author, branch, files)
  - [ ] Create specialized prompts based on filter context
  - [ ] Generate targeted code quality and practice insights
- [ ] Configuration Extensions:
  - [ ] Add summary settings to config file structure:
    ```toml
    [summary]
    default_timespan = "7d"  # 7d, 14d, 30d
    include_stats = true
    export_format = "text"   # text, markdown, html
    ```
  - [ ] Support custom prompt templates for summaries
  - [ ] Add personality-specific summary templates
- [ ] Reporting & Visualization:
  - [ ] Design simple ASCII/Unicode charts for terminal
  - [ ] Add commit heatmap visualization
  - [ ] Implement pattern detection algorithms
  - [ ] Support configurable highlighting of notable trends

---

### 🧼 Phase 6: Polish & Ship

#### ✅ Tasks:
- [ ] Installation & Distribution:
  - [ ] Add install instructions (`go install`, brew tap)
  - [ ] Create cross-platform release binaries (macOS, Linux, Windows)
  - [ ] Add one-line install script for quick setup
- [ ] User Experience Improvements:
  - [ ] Add `--verbose` and `--silent` flags for control over output
  - [ ] Create interactive `noidea setup` wizard for first-time users
  - [ ] Add colorful progress bars for long-running operations
  - [ ] Add usage GIFs and examples in `README`
- [ ] Developer Insights:
  - [ ] Create achievement system for commit streaks and milestones
  - [ ] Add gamification elements (levels, badges)
  - [ ] Support for personalized improvement suggestions
- [ ] Team & Collaboration:
  - [ ] Add shareable personality configurations via gists
  - [ ] Create team leaderboard for most commits/streaks (optional)
  - [ ] Support for team-specific feedback rules
- [ ] Extensibility:
  - [ ] Create plugin system for community extensions
  - [ ] Add integration hooks for CI/CD systems
  - [ ] Support for custom feedback triggers beyond commits
- [ ] Quality & Testing:
  - [ ] Implement comprehensive test suite
  - [ ] Add telemetry opt-in for improving personalities
  - [ ] Create benchmark suite for performance testing

---