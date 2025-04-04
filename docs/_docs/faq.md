---
layout: default
title: FAQ
nav_order: 4
permalink: /docs/faq
---

# Frequently Asked Questions
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## General Questions

### What is noidea?

noidea is a Git companion that makes your commits better and funnier by providing AI-powered commit message suggestions and sassy feedback after each commit.

### Why should I use noidea?

Because Git is too serious, and coding is chaos. noidea won't necessarily improve your Git hygiene, but it will make your Git experience more entertaining and potentially more informative.

### Is noidea free to use?

Yes, noidea is open-source and free to use. However, some AI features require an API key from supported providers.

## Setup and Configuration

### How do I install noidea?

You can install noidea using the one-line install command:
```bash
curl -sSL https://raw.githubusercontent.com/AccursedGalaxy/noidea/main/quickinstall.sh | bash
```

Or by cloning the repository and running the install script:
```bash
git clone https://github.com/AccursedGalaxy/noidea.git
cd noidea
./install.sh
```

### Do I need an API key?

Yes, for AI-powered features like commit message suggestions and AI feedback, you'll need an API key from a supported provider.

### How do I configure noidea?

You can configure noidea interactively:
```bash
noidea config --init
```

Or manually by creating a `~/.noidea/config.toml` file. See the [advanced configuration section](/noidea/docs/configuration) for more details.

## Features and Usage

### How do I get commit message suggestions?

You can get commit message suggestions by running:
```bash
noidea suggest
```

Or by enabling auto-suggestions:
```bash
git config noidea.suggest true
```

### Can I customize the Moai personalities?

Yes, you can choose from several built-in personalities or create your own by editing the personalities configuration file.

### Does noidea work with all Git workflows?

noidea should work with most standard Git workflows. However, it might not integrate perfectly with all Git GUIs or specialized workflows.

### Can I disable specific features?

Yes, you can disable specific features by modifying your configuration. For example, to disable the Moai feedback:
```toml
[moai]
enabled = false
``` 