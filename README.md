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
git clone https://github.com/yourname/noidea.git
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

---

## 💡 Upcoming Features

| Feature                   | Status          |
|---------------------------|-----------------|
| Moai face after commit    | ✅ Done          |
| AI-based commit feedback  | ✅ Done          |
| Lint feedback             | 🛠️ In progress   |
| Commit streak insights    | 🔜 Coming Soon   |
| Config file support       | 🔜 Coming Soon   |

---

## 🔧 Configuration (coming soon)

You'll be able to configure:
- Whether linting is checked
- Types of Moai reactions
- AI mode on/off
- Local vs. cloud model

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
