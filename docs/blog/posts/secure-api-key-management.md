---
date: 2024-04-05
authors:
  - accursedgalaxy
categories:
  - Security
  - Development
tags:
  - api-keys
  - security
  - cross-platform
  - golang
readtime: 8
---

# Building a Cross-Platform Secure API Key Manager in Go

While developing NoIdea, my Git assistant with AI-powered features, I faced a significant challenge: how could I securely store API keys across different operating systems without compromising security or frustrating users? In this article, I'll walk you through my solution and the lessons I learned along the way.

<!-- more -->

## The Security Challenge Explained

First, let me explain the problem. Modern AI tools like NoIdea need to communicate with AI services (like OpenAI or xAI), which require API keys for authentication. These keys are essentially long strings that look something like `xai_sk_1a2b3c4d5e6f7g8h9i0j`.

The tricky part? API keys are like digital passwords to paid services. If someone gets your key, they could:
- Use your AI services and run up charges on your account
- Potentially access your data
- Impersonate you to the service

Common approaches for storing these keys have serious drawbacks:

- ❌ **Environment variables**: These are variables stored in your computer's memory while you're using it. The problem? They're often visible in process listings, meaning other programs can potentially see them.

- ❌ **Plaintext `.env` files**: Many developers store keys in simple text files called `.env` files. If you accidentally commit these to GitHub or other version control systems, your keys are exposed to the world (there are even bots that scan GitHub for leaked API keys).

- ❌ **Hardcoded credentials**: Some developers embed keys directly in their code. This is perhaps the worst approach since your keys are permanently recorded in your code history.

I needed a solution that would work across Windows, macOS, and Linux while keeping keys secure and still providing a smooth experience for users.

## My Tiered Security Solution

After researching the problem, I designed a system with multiple layers of protection:

```go
// Example from my codebase
func GetAPIKey(provider string) (string, error) {
    // Standardize the provider name for consistency
    provider = normalizeProviderName(provider)
    
    // Try system keyring first (macOS Keychain, Windows Credential Manager, etc.)
    apiKey, err := keyring.Get(ServiceName, provider)
    if err == nil && apiKey != "" {
        return apiKey, nil
    }
    
    // Fall back to encrypted storage if keyring unavailable
    return getFromFallbackStorage(provider)
}
```

Let me break down what this does:

1. **First Security Tier**: The function tries to use the operating system's built-in credential storage:
   - On a Mac, this is the Keychain (the same secure storage that remembers your WiFi passwords)
   - On Windows, it uses the Windows Credential Manager
   - On Linux, it uses the Secret Service API (programs like GNOME Keyring or KDE Wallet)

   These are secure, encrypted storage systems managed by your operating system with proper permissions.

2. **Second Security Tier**: If the system keyring isn't available (which happens on some Linux setups or in certain environments), I fall back to a custom encrypted storage in the user's home directory.

This "defense in depth" approach ensures that even if one security measure fails, there's a backup. It's like having both a lock on your door and a safe inside your house.

## How I Implemented It (With Real-World Examples)

### 1. Solving the Provider Name Problem

One of the first challenges I encountered was that users might refer to the same AI provider in different ways. For example:
- "openai", "open-ai", "gpt" all refer to OpenAI
- "xai", "x-ai", "grok" all refer to xAI

This might seem like a minor issue, but it can lead to confusing situations. Imagine storing your key under "openai" but then your program looking for it under "open-ai" - it would fail to find it!

My first implementation was simple but limited:

```go
// My initial implementation - works but isn't flexible
func normalizeProviderName(provider string) string {
    provider = strings.ToLower(provider)
    
    // Hardcoded mappings - not ideal for extensibility
    switch provider {
    case "openai", "open-ai", "gpt":
        return "openai"
    case "xai", "x-ai", "grok":
        return "xai" 
    case "deepseek", "deep-seek":
        return "deepseek"
    default:
        return provider
    }
}
```

This worked initially, but I quickly realized it had limitations:
- What if a new AI provider emerged? I'd need to update my code.
- What if users wanted to use their own custom naming? They'd be out of luck.
- What if a company rebranded? Again, code updates would be needed.

So I refactored it to a more flexible system:

```go
// Default provider alias mapping - maps standard provider names to their known aliases
var defaultProviderAliases = map[string][]string{
    "openai":    {"open-ai", "gpt", "chatgpt", "davinci"},
    "xai":       {"x-ai", "grok", "x.ai"},
    "deepseek":  {"deep-seek", "deepseek-ai"},
    "anthropic": {"claude", "anthropic-ai"},
    "mistral":   {"mistral-ai", "mistralai"},
}

// Reverse lookup map built at init time
var aliasToProvider map[string]string

func init() {
    // Load provider aliases (default + user-defined)
    providerAliases := loadProviderAliases()
    
    // Build reverse lookup map
    aliasToProvider = make(map[string]string)
    for provider, aliases := range providerAliases {
        aliasToProvider[provider] = provider // Map standard name to itself
        for _, alias := range aliases {
            aliasToProvider[alias] = provider
        }
    }
}

func normalizeProviderName(provider string) string {
    provider = strings.ToLower(provider)
    
    // Look up in our alias map
    if standardName, exists := aliasToProvider[provider]; exists {
        return standardName
    }
    
    // If no match, return as-is
    return provider
}
```

The key improvement is that now users can define their own aliases through a JSON configuration file:

```go
// loadProviderAliases combines default aliases with user-defined ones
func loadProviderAliases() map[string][]string {
    // Start with default aliases
    combined := make(map[string][]string)
    for provider, aliases := range defaultProviderAliases {
        combined[provider] = aliases
    }
    
    // Try to load user-defined aliases from ~/.noidea/secure/provider_aliases.json
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return combined // Fall back to defaults
    }
    
    // Load and parse user aliases from JSON file
    aliasPath := filepath.Join(homeDir, FallbackDir, AliasFile)
    data, err := os.ReadFile(aliasPath)
    if err != nil {
        // Create template file for the user
        createDefaultAliasFile(homeDir)
        return combined
    }
    
    // Merge user aliases with defaults
    var userAliases map[string][]string
    if err := json.Unmarshal(data, &userAliases); err == nil {
        for provider, aliases := range userAliases {
            if existing, ok := combined[provider]; ok {
                // Add user aliases, avoiding duplicates
                for _, alias := range aliases {
                    if !contains(existing, alias) {
                        combined[provider] = append(combined[provider], alias)
                    }
                }
            } else {
                // Add new provider
                combined[provider] = aliases
            }
        }
    }
    
    return combined
}
```

Think of this like a translation dictionary. When a user types "gpt", my code looks it up and says "Oh, you mean OpenAI!" This way, the actual provider name is standardized for storage, but users can use whatever terminology they're familiar with.

#### Why This Matters

This might seem like excessive effort for a simple name mapping, but it has real usability benefits:

1. **User Freedom**: Different communities use different terminology. AI researchers might say "OpenAI," while developers might say "GPT," and a team might have their own internal name like "chat-provider-1".

2. **Future-Proofing**: New AI providers emerge regularly. With this system, users can add support for new providers without waiting for a software update.

3. **Team Standardization**: In a team environment, you can establish standardized names that match your documentation or internal systems.

When users first use the system, I automatically create a template configuration file they can edit:

```json
// Example of the JSON template created at ~/.noidea/secure/provider_aliases.json
{
  "example-provider": ["alias1", "alias2"],
  "openai": ["gpt4", "oai"]
}
```

### 2. Securing Keys on Different Operating Systems

The next challenge was creating a secure storage system that works across platforms. Different operating systems have different security models:

```go
func storeInFallbackStorage(provider, apiKey string) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }
    
    secureDir := filepath.Join(homeDir, FallbackDir)
    if err := os.MkdirAll(secureDir, 0700); err != nil {
        return fmt.Errorf("failed to create secure directory: %w", err)
    }
    
    // Update or add the new key with obfuscation
    existingData[provider] = obfuscate(apiKey)
    
    // Write with strict permissions - only the user can read/write
    return os.WriteFile(filePath, []byte(sb.String()), 0600)
}
```

This fallback system has several security features:
- A dedicated directory with restricted permissions (0700 means only the owner can access it)
- File permissions that limit access (0600 again ensures only the owner can read or write)
- Simple obfuscation to prevent casual viewing

While not as secure as the OS keyring, it's significantly better than plaintext storage in most cases.

### 3. Validating Keys Before Storage

An often-overlooked security feature is validation. I've seen too many systems that happily store invalid API keys, leading to mysterious failures later. 

My system checks if a key is valid before storing it:

```go
isValid, err := secure.ValidateAPIKey(provider, apiKey)
if err != nil {
    fmt.Println(color.YellowString("Error during validation"))
    fmt.Printf("Error details: %v\n", err)
} else if isValid {
    fmt.Println(color.GreenString("Valid"))
} else {
    fmt.Println(color.RedString("Invalid"))
}
```

Here's how the validation works:

```go
// ValidateAPIKey checks if the API key works with the provider
func ValidateAPIKey(provider, apiKey string) (bool, error) {
    // For all providers, try to validate against their API
    var baseURL string
    
    switch provider {
    case "xai":
        baseURL = "https://api.x.ai/v1/models"
    case "openai":
        baseURL = "https://api.openai.com/v1/models"
    case "deepseek":
        baseURL = "https://api.deepseek.com/v1/models"
    default:
        // Default to OpenAI for unknown providers
        baseURL = "https://api.openai.com/v1/models"
    }
    
    // Make a test API call that doesn't cost money
    client := &http.Client{Timeout: 5 * time.Second}
    req, err := http.NewRequest("GET", baseURL, nil)
    if err != nil {
        return false, err
    }
    
    req.Header.Add("Authorization", "Bearer "+apiKey)
    
    resp, err := client.Do(req)
    if err != nil {
        return false, fmt.Errorf("connection error: %w", err)
    }
    defer resp.Body.Close()
    
    // 401/403 means unauthorized - invalid key
    return resp.StatusCode != 401 && resp.StatusCode != 403, nil
}
```

This validation step offers several benefits:
- Catches typos or incorrect keys immediately
- Avoids storing keys that won't work
- Provides immediate feedback to users
- Helps distinguish between network errors and invalid keys

### 4. Making Migration Easy

Recognizing that many users would already have API keys stored in environment variables or `.env` files, I created a migration script:

```bash
# Automatically detect and migrate API keys from various sources
for location in "${LOCATIONS[@]}"; do
    if [ -f "$location" ]; then
        echo -e "${YELLOW}Found .env file at $location${NC}"
        
        # Check for API keys
        if grep -q "XAI_API_KEY" "$location"; then
            API_KEY=$(grep "XAI_API_KEY" "$location" | cut -d= -f2)
            PROVIDER="xai"
            echo "Found xAI API key"
        elif grep -q "OPENAI_API_KEY" "$location"; then
            # Handle other providers...
        fi
        
        # ...migration code...
    fi
done
```

This script:
1. Scans common locations for `.env` files
2. Checks for known API key environment variables
3. Extracts keys and their associated providers
4. Validates the keys
5. Stores them in the secure storage system
6. Optionally removes them from the insecure locations

## A User-Friendly Interface for Security

I've found that security mechanisms are only effective if people actually use them. To encourage adoption, I designed a clean CLI interface:

```
$ noidea config apikey-status

Secure Storage Status:
Platform: linux
System keyring: Available
Fallback storage: Available

API Key Status:
Provider: xai
Environment: Not set
Secure storage: Set

Active Key:
Using: Secure storage key
Validating secure key... Valid
```

For users who want to customize provider aliases, they can edit a simple JSON file:

```bash
# Edit the provider aliases file
nano ~/.noidea/secure/provider_aliases.json

# Example custom configuration
{
  "openai": ["my-openai", "chatgpt-custom"],
  "my-company-api": ["internal-llm", "company-ai"]
}
```

## Lessons I Learned (That You Can Apply)

Building this system taught me several valuable lessons that you can apply to your own projects:

1. **Security doesn't have to hurt UX**: By providing fallbacks and clear guidance, you can maintain strong security without frustrating users. Always aim for secure defaults with clear paths to customize.

2. **Trust the platform when possible**: Native OS security features are typically well-audited and robust. I use the system keyring when available because it's maintained by security professionals and integrated with the OS.

3. **Defense in depth works**: Multiple security layers provide protection even when one layer fails. My tiered approach means that even if one security mechanism is compromised, there are additional protections.

4. **Validation improves reliability**: Checking keys before storage prevents frustrating invalid credential issues later. This seems obvious but is often overlooked.

5. **Prefer configuration over code**: For elements that might change or expand (like provider aliases), use configuration-based approaches rather than hardcoding values. This makes your software more adaptable without needing code changes.

## How You Can Apply This in Your Projects

If you're building an application that needs to store sensitive credentials, here are some practical steps:

1. **Look for existing secure storage APIs**:
   - On Windows: [Windows Credential Manager API](https://learn.microsoft.com/en-us/windows/win32/secauthn/credential-manager)
   - On macOS: [Keychain Services API](https://developer.apple.com/documentation/security/keychain_services)
   - On Linux: [Secret Service API](https://specifications.freedesktop.org/secret-service/latest/)
   - Cross-platform: Consider libraries like [go-keyring](https://github.com/zalando/go-keyring) or [keyring](https://pypi.org/project/keyring/) for Python

2. **Implement a fallback mechanism**:
   - Store in user's home directory with restricted permissions
   - Use basic encryption/obfuscation at minimum
   - Clearly document the security model

3. **Validate credentials at storage time**:
   - Make a test request that doesn't consume resources
   - Provide clear feedback about validity
   - Allow force-storing with warnings if needed

4. **Design for user customization**:
   - Allow configuration without code changes
   - Provide sensible defaults
   - Create templates for common scenarios

5. **Help users migrate**:
   - Provide tools to detect and import credentials from common locations
   - Guide users through the transition
   - Don't break existing workflows during migration

## Future Improvements

While my current system works well, I'm considering several improvements:

1. **More robust encryption** for the fallback storage mechanism - perhaps using industry-standard encryption libraries rather than simple obfuscation
2. **Additional validation methods** for different API providers to handle their specific authentication patterns
3. **Enhanced key rotation policies** and workflows to encourage regular key updates
4. **Web interface** for managing provider aliases and credentials for less technical users

## Conclusion

Securing API keys is a critical aspect of modern application development, especially for tools that interact with AI services. By implementing this tiered approach, I've built a system that:

- ✅ Provides strong security guarantees
- ✅ Works seamlessly across platforms
- ✅ Offers a smooth user experience
- ✅ Handles migration gracefully
- ✅ Remains flexible and extensible through user-configurable aliases

The most important lesson I've learned is that good security doesn't have to come at the expense of usability. With thoughtful design and attention to user needs, you can create security systems that people actually want to use.

For more information about NoIdea's approach to API key management, check out the [API Key Management documentation](../api-key-management.md). 