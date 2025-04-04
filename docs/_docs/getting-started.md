---
layout: default
title: Getting Started
nav_order: 2
permalink: /docs/getting-started
---

# Getting Started
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Installation

Choose one of these methods to install noidea:

### One-line Quick Install

```bash
curl -sSL https://raw.githubusercontent.com/AccursedGalaxy/noidea/main/quickinstall.sh | bash
```

### Manual Install

```bash
# Clone the repository
git clone https://github.com/AccursedGalaxy/noidea.git
cd noidea

# Run the install script
./install.sh # (might require sudo)
```

## Setting Up

After installation, you need to initialize noidea in your Git repository:

```bash
# Navigate to your Git repository
cd /path/to/your/repo

# Initialize noidea
noidea init

# Enable auto commit suggestions (optional)
git config noidea.suggest true
```

## AI Configuration

For AI-powered features, you'll need to add your API key:

### Option 1: Environment Variable

```bash
export XAI_API_KEY=your_api_key_here
```

### Option 2: Config File

Create or edit `~/.noidea/.env`:

```
XAI_API_KEY=your_api_key_here
```

### Option 3: Interactive Setup

```bash
noidea config --init
```

## Basic Usage

Once set up, your new Git workflow with noidea will be:

```bash
# Add your changes
git add .

# When you commit, noidea will suggest a commit message
git commit
```

This will open your default editor with the suggested commit message. Saving and closing will approve and commit.

You can also explicitly request a commit message suggestion:

```bash
noidea suggest
``` 