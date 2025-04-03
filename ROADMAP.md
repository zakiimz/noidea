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

### 🧪 Phase 3: LLM Config & User Settings

**🔹 Goal:** Let the user configure LLM behavior.

#### ✅ Tasks:
- [x] Support multiple LLM providers (xAI, OpenAI, DeepSeek)
- [x] Support environment variables for configuration
- [ ] Create a config file:
  - Location: `~/.noidea/config.toml`
- [ ] Config structure:
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
- [ ] CLI command: `noidea config` (optional for manual key entry)
- [ ] Validate config + give feedback on missing keys
- [ ] Allow overriding config with env vars (e.g. `NOIDEA_API_KEY`)

---

### 📊 Phase 4: Commit Pattern Analyzer (Offline Logic)

**🔹 Goal:** Provide deeper insights based on Git behavior without needing AI.

#### ✅ Tasks:
- [ ] Track:
  - Time of day
  - Frequency of commits
  - Message patterns (e.g., "fix", "final", "pls work")
- [ ] Generate local-only messages based on trends:
  - "You haven't committed in 3 days"
  - "5 commits with the same message detected"
- [ ] Cache commit data in local SQLite or JSON

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

### 🛠️ Current Project Structure

```
noidea/
├── cmd/
│   ├── root.go            # Root command
│   ├── init.go            # Init command to install Git hook
│   └── moai.go            # Moai command for feedback generation
├── internal/
│   ├── config/
│   │   └── config.go      # Configuration loading and management
│   ├── feedback/
│   │   ├── engine.go      # FeedbackEngine interface
│   │   ├── unified.go     # Unified LLM engine for all providers
│   │   ├── local.go       # Local feedback engine (no API)
│   │   └── utils.go       # Shared utility functions
│   ├── git/
│   │   └── hooks.go       # Git hook installation logic
│   └── moai/
│       └── faces.go       # Moai faces and random feedback
├── scripts/
│   └── post-commit.sh     # Template Git hook
├── go.mod
├── main.go
└── README.md
```
