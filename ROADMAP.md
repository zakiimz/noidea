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

### 🧼 Phase 5: Polish & Ship

#### ✅ Tasks:
- [ ] Add install instructions (`go install`, releases)
- [ ] Add `--verbose` and `--silent` flags
- [ ] Add `noidea feedback` command (manual insight trigger)
- [ ] Write tests for core components
- [ ] Prepare cross-platform release binaries
- [ ] Add usage GIF in `README`

---