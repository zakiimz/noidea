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
  - [x] Implement "Motivational Speaker" personality
- [x] Documentation:
  - [x] Add personality configuration guide
  - [x] Document prompt template syntax
  - [x] Include personality customization examples

---

### 📈 Phase 5: AI-Driven Commit Insights ✅

**🔹 Goal:** Provide actionable and interesting insights into commit history patterns using LLM analysis.

#### ✅ Tasks:
- [x] Git History Collection:
  - [x] Create `HistoryCollector` module in `internal/history/`
  - [x] Implement flexible date range filtering (past 7 days, past X commits)
  - [x] Extract commit metadata (messages, authors, timestamps, files changed)
  - [x] Add diff summarization for contextual understanding
  - [x] Implement caching to reduce repeat Git operations

- [x] Weekly Summary Feature:
  - [x] Add `noidea summary` command
  - [x] Design comprehensive prompt for weekly work analysis
  - [x] Create statistical overview (commits/day, files touched, etc.)
  - [x] Generate human-readable summary with AI-powered insights
  - [x] Support markdown/HTML export

- [x] On-Demand Feedback:
  - [x] Implement `noidea feedback [--count X]` command
  - [x] Add filtering options (author, branch, files)
  - [x] Create specialized prompts based on filter context
  - [x] Generate targeted code quality and practice insights

- [x] Commit Message Suggestions:
  - [x] Implement `noidea suggest` command
  - [x] Create specialized system prompt for professional commit messages
  - [x] Add support for Git hook integration (`prepare-commit-msg`)
  - [x] Implement interactive mode for approval/editing
  - [x] Add message extraction for clarity and compatibility

---

### 🧼 Phase 6: Polish & Ship 🚧

#### Tasks:
- [x] Installation & Distribution:
  - [x] Add install instructions (`go install`, brew tap)
  - [x] Create cross-platform release binaries (macOS, Linux, Windows)
  - [x] Add one-line install script for quick setup
- [x] User Experience Improvements:
  - [x] Add `--verbose` and `--silent` flags for control over output
  - [x] Create interactive `noidea setup` wizard for first-time users
  - [x] Add colorful dividers and formatting for improved terminal output
  - [x] Enhance visual feedback with emoji and color-coded statuses
  - [x] Ensure shell compatibility across platforms (POSIX-compliant)
  - [ ] Add usage GIFs and examples in `README`
- [x] Quality & Testing:
  - [x] Implement test suite with simulation framework
  - [x] Create test repository fixtures for consistent testing
  - [x] Add shell script and Go-based test harnesses

---