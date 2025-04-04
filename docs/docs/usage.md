---
layout: default
title: Usage
nav_order: 3
permalink: /docs/usage
---

# Usage
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Commands Overview

noidea provides several commands to make your Git experience more enjoyable:

| Command | Description |
|---------|-------------|
| `noidea init` | Set up Git hooks in your repository |
| `noidea suggest` | Get commit message suggestions |
| `noidea moai` | Display Moai feedback for the last commit |
| `noidea summary` | Generate summary of recent Git activity |
| `noidea feedback` | Analyze specific commits |
| `noidea config` | Configure noidea |

Run `noidea --help` for more detailed information.

## Commit Message Suggestions

When you're ready to commit, you can use:

```bash
noidea suggest
```

This will analyze your staged changes and suggest a professional commit message:

```
───────────────────────────────────────────────────
🧠 Analyzing staged changes...
───────────────────────────────────────────────────
✨ Suggested commit message:
feat(user-auth): implement password reset functionality with email verification
───────────────────────────────────────────────────
```

## Post-Commit Feedback

After each commit, the Moai will automatically provide feedback on your work:

```
───────────────────────────────────────────────────
🗿  (ಠ_ಠ) Your commit message was 'fix final final pls real'
"You've entered the 2AM hotfix arc. A legendary time."
───────────────────────────────────────────────────
```

## Weekly Summaries

Get insights about your recent Git activity:

```bash
noidea summary
```

You can specify the number of days to include in the summary:

```bash
noidea summary --days 30
```

## AI Personalities

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