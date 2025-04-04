---
layout: default
title: Installation Guide
---

# Installation Guide

## Prerequisites

- Go 1.23+
- Git
- An API key from xAI, OpenAI, or DeepSeek (for AI features)

## Quick Install

The easiest way to install noidea is with our one-line installer:

```bash
curl -sSL https://raw.githubusercontent.com/AccursedGalaxy/noidea/main/quickinstall.sh | bash
```

## Manual Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/AccursedGalaxy/noidea.git
   cd noidea
   ```

2. Run the installation script:
   ```bash
   ./install.sh
   ```
   
   This script will:
   - Build the noidea binary
   - Install it to `/usr/local/bin` (or custom location)
   - Create the `~/.noidea` directory for configuration
   - Install the necessary Git hooks and scripts

3. Set up in your repository:
   ```bash
   cd /path/to/your/repo
   noidea init
   ```

## Configuration

For AI-powered features, you'll need to configure your API key:

### Add to Environment

```bash
export XAI_API_KEY=your_api_key_here
```

### Create Configuration File

Create `~/.noidea/.env`:

```
XAI_API_KEY=your_api_key_here
```

### Interactive Configuration

Use the built-in configuration tool:

```bash
noidea config --init
```

## Verifying Installation

Check that noidea is installed correctly:

```bash
noidea --version
```

You should see version information and a list of available commands. 