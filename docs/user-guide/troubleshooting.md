# Troubleshooting

This guide addresses common issues you might encounter while using noidea.

## Installation Issues

### noidea Command Not Found

**Problem**: After installation, running `noidea` results in "command not found".

**Solutions**:

1. Check if the binary is in your PATH:
   ```bash
   which noidea
   ```

2. If not found, verify the installation location:
   ```bash
   # Installed to /usr/local/bin
   ls -la /usr/local/bin/noidea
   
   # Or in your Go bin directory
   ls -la ~/go/bin/noidea
   ```

3. Add the installation directory to your PATH if needed:
   ```bash
   # Add to your .bashrc or .zshrc
   export PATH=$PATH:/path/to/directory
   ```

### Installation Fails with Permission Error

**Problem**: `./install.sh` fails with permission errors.

**Solution**: Run with sudo or specify a user-writable location:
   ```bash
   sudo ./install.sh
   # Or
   ./install.sh --prefix ~/bin
   ```

## API Key Issues

### API Key Validation Fails

**Problem**: "API key validation error" or "Invalid API key" messages.

**Solutions**:

1. Verify your API key:
   ```bash
   noidea config apikey-status
   ```

2. Re-enter your API key:
   ```bash
   noidea config apikey
   ```

3. Check your internet connection and provider status.

### Can't Access AI Features

**Problem**: AI-powered features like `suggest --ai` or `moai --ai` don't work.

**Solutions**:

1. Ensure AI features are enabled:
   ```bash
   # Check configuration
   noidea config --show | grep enabled
   
   # Enable if needed
   noidea config set llm.enabled true
   ```

2. Set up your API key if you haven't already:
   ```bash
   noidea config apikey
   ```

## Git Integration Issues

### Git Hooks Not Working

**Problem**: noidea's Git hooks (commit suggestions, Moai feedback) aren't running.

**Solutions**:

1. Verify hooks are installed in your repository:
   ```bash
   ls -la .git/hooks/prepare-commit-msg .git/hooks/post-commit
   ```

2. If missing, run:
   ```bash
   noidea init
   ```

3. Ensure hooks are executable:
   ```bash
   chmod +x .git/hooks/prepare-commit-msg .git/hooks/post-commit
   ```

4. Check if hooks are bypassed:
   ```bash
   # Make sure you're not using --no-verify
   git config --get noidea.suggest
   ```

### Commit Suggestions Not Appearing

**Problem**: The commit message suggestion feature isn't working.

**Solutions**:

1. Make sure it's enabled:
   ```bash
   git config --get noidea.suggest
   
   # Enable if needed
   git config noidea.suggest true
   ```

2. Verify you have staged changes:
   ```bash
   git status
   ```

3. Try running the suggest command directly:
   ```bash
   noidea suggest
   ```

## Performance Issues

### Slow Commit Suggestions

**Problem**: Generating commit suggestions takes too long.

**Solutions**:

1. Disable full-diff mode:
   ```bash
   git config noidea.suggest.full-diff false
   ```

2. Reduce history context:
   ```bash
   git config noidea.suggest.history 5
   ```

## Configuration Issues

### Configuration Changes Not Applied

**Problem**: Changes to configuration don't seem to take effect.

**Solutions**:

1. Check current configuration:
   ```bash
   noidea config --show
   ```

2. Verify configuration file location:
   ```bash
   ls -la ~/.noidea/config.json
   ```

3. Ensure you're using the correct method (environment variables might override file settings).

## Getting More Help

If your issue isn't addressed here:

1. Check the logs for more information (if enabled):
   ```bash
   cat ~/.noidea/logs/noidea.log
   ```

2. [Open an issue](https://github.com/AccursedGalaxy/noidea/issues) on GitHub with:
   - noidea version (`noidea --version`)
   - OS details
   - Error messages
   - Steps to reproduce 