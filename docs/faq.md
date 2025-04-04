---
layout: default
title: Frequently Asked Questions
---

# Frequently Asked Questions

## General Questions

### What is noidea?

noidea is a Git companion tool that enhances your Git workflow with AI-powered commit message suggestions and humorous feedback from a judgmental Moai after each commit.

### Why the name "noidea"?

The name represents the feeling many developers have when writing commit messages - "no idea" what to write. It's a playful nod to the struggle of writing meaningful commit messages.

### What is the Moai thing about?

The Moai (🗿) is a friendly but judgmental companion that provides sassy feedback after your commits. It's a way to make the Git experience more fun and less serious.

## Installation & Requirements

### What are the requirements to run noidea?

- Go 1.23+
- Git
- An API key from xAI, OpenAI, or DeepSeek (for AI features)

### Do I need an API key to use noidea?

For basic functionality, no. However, AI-powered features like commit message suggestions work best with an API key. Without one, noidea will fall back to simpler algorithms.

### How do I update noidea?

```bash
cd /path/to/noidea
git pull
./install.sh
```

## Features & Usage

### How do I get commit message suggestions?

```bash
# Manually request a suggestion
noidea suggest

# Enable automatic suggestions during commit
git config noidea.suggest true
```

### Can I use noidea without the Moai feedback?

Yes, you can disable the Moai feedback:

```bash
git config noidea.moai false
```

### How do I customize the personality of feedback?

You can change the personality using:

```bash
noidea moai --personality supportive_mentor
```

Or set a default in your Git config:

```bash
git config noidea.personality git_expert
```

## Troubleshooting

### noidea isn't suggesting commit messages

1. Check if you have staged changes with `git status`
2. Verify your API key is correctly set up
3. Ensure noidea is initialized in your repository with `noidea init`

### The Moai feedback doesn't appear after commits

1. Make sure the post-commit hook is installed: `noidea init`
2. Check if Moai is disabled: `git config noidea.moai`
3. Look for errors in the Git hook execution

### I'm getting API errors

1. Verify your API key is correctly set
2. Check your internet connection
3. Ensure the AI provider's services are operational

## Contributing

### How can I contribute to noidea?

We welcome contributions! Check out our [Contributing Guide](https://github.com/AccursedGalaxy/noidea/blob/main/CONTRIBUTING.md) for details on how to get started. 