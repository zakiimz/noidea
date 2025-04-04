# API Key Management

This document provides best practices for managing API keys with noidea.

## Secure Storage

As of version v0.3.0, noidea securely stores API keys using your system's native keyring/keychain when available:

- **macOS**: Uses the Keychain
- **Windows**: Uses the Windows Credential Manager
- **Linux**: Uses the Secret Service API (requires libsecret)

If the system keyring is unavailable, a fallback encrypted storage is used in `~/.noidea/secure/`.

## Setting Up Your API Key

You can set up your API key in several ways:

### 1. Using the CLI (Recommended)

```bash
# Set up API key securely
noidea config apikey

# Check API key storage status
noidea config apikey-status

# Remove a stored API key
noidea config apikey-remove

# Generate commands to clean environment variables
noidea config clean-env
```

When setting up a key with `noidea config apikey`, the system will:
1. Prompt for the API key (input is hidden for security)
2. Validate the key with the provider to ensure it works
3. Store the key securely in your system's keyring or fallback storage

The `apikey-status` command will:
1. Show which storage system is being used
2. Check if your API key is valid with a test request
3. Display whether the key is working correctly

### 2. Using Environment Variables (Alternative)

You can still use environment variables if you prefer:

```bash
# For xAI
export XAI_API_KEY=your_api_key_here

# For OpenAI
export OPENAI_API_KEY=your_api_key_here

# For DeepSeek
export DEEPSEEK_API_KEY=your_api_key_here

# Generic (will use whatever provider is configured)
export NOIDEA_API_KEY=your_api_key_here
```

**Important Note**: Environment variables will take precedence over secure storage. If you want to use secure storage, make sure these environment variables are not set.

### 3. Using .env Files (Not Recommended)

While still supported for backward compatibility, we recommend transitioning away from .env files:

```bash
# Create or edit ~/.noidea/.env
XAI_API_KEY=your_api_key_here
```

## API Key Priority Order

The system uses the following order of precedence when looking for API keys:

1. Environment variables (highest priority)
2. Secure storage (keyring/keychain or fallback encrypted file)
3. Config file (lowest priority - not recommended for API keys)

If you've set up a key using secure storage but it's not being used, check if any environment variables are overriding it with:

```bash
noidea config apikey-status
```

To clean environment variables and use secure storage instead:

```bash
noidea config clean-env
```

This will generate commands you can run to remove API key environment variables.

## API Key Security Best Practices

1. **Never commit API keys to version control**
   - Ensure `.env` files are in your `.gitignore`
   - Use secure storage or environment variables instead

2. **Rotate keys periodically**
   - Change your API keys regularly
   - Use `noidea config apikey` to update your stored key

3. **Use the least privileged key possible**
   - Only use keys with the permissions your application needs

4. **Monitor key usage**
   - Check the provider's dashboard for unusual activity
   - Set up usage alerts if available

5. **Environment separation**
   - Use different keys for development and production

## Troubleshooting

If you encounter issues with secure storage:

1. **Check storage status and key validity**
   ```bash
   noidea config apikey-status
   ```
   This command will verify:
   - If your keyring is available
   - If your API key is properly stored
   - If your API key is valid and working

2. **Ensure dependencies are installed**
   - On Linux, install libsecret: `sudo apt-get install libsecret-1-dev`

3. **Remove environment variables**
   - If secure storage is working but not being used, environment variables may be taking precedence
   ```bash
   noidea config clean-env
   ```

4. **If validation fails**
   - Check if your API key is correct
   - Ensure you have an active subscription with the provider
   - Check if your network can reach the provider's API servers

## Migration from Previous Versions

If you're upgrading from a version before v0.3.0:

1. Use the migration script to move your API key to secure storage:
   ```bash
   ./scripts/migrate_to_secure.sh
   ```

2. Or run `noidea config apikey` to set up your API key securely

3. Remove any API keys from `.env` files or your config file

4. Remove API key environment variables with:
   ```bash
   noidea config clean-env
   ```

For more information, please see the [Configuration Guide](./configuration.md). 