# Installation Guide

There are several ways to install noidea on your system. Choose the option that works best for you.

## Quick Install (Recommended)

The fastest way to get started is with our one-line installer:

```bash
curl -sSL https://raw.githubusercontent.com/AccursedGalaxy/noidea/main/quickinstall.sh | bash
```

This script will:
- Download the latest noidea release
- Install it to `/usr/local/bin` (might require sudo)
- Make it executable

## Manual Installation

### Option 1: Clone and Install

```bash
# Clone the repository
git clone https://github.com/AccursedGalaxy/noidea.git
cd noidea

# Run the installer script
./install.sh
```

### Option 2: Download Binary

1. Go to the [Releases page](https://github.com/AccursedGalaxy/noidea/releases)
2. Download the appropriate binary for your platform
3. Make it executable: `chmod +x noidea`
4. Move it to your PATH: `sudo mv noidea /usr/local/bin/`

## Building from Source

If you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/AccursedGalaxy/noidea.git
cd noidea

# Build the binary
go build -o noidea

# Install to PATH
sudo mv noidea /usr/local/bin/
```

## Verifying Installation

To verify that noidea is installed correctly:

```bash
noidea --version
```

This should display the version information.

## Next Steps

After installation:

1. [Set up noidea in your repository](getting-started.md#quick-setup)
2. [Configure your API key](configuration.md#api-key-setup) for AI features
3. Explore the [available commands](commands/overview.md) 