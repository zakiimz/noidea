# Plugin Interface Specifications

This document defines the core interfaces and types that plugins must implement to integrate with NoIdea.

## Core Plugin Interface

Every NoIdea plugin must implement the `Plugin` interface:

```go
// Plugin is the core interface that all plugins must implement
type Plugin interface {
    // Info returns metadata about the plugin
    Info() PluginInfo
    
    // Initialize prepares the plugin for use
    Initialize(ctx PluginContext) error
    
    // Shutdown performs cleanup when the plugin is disabled
    Shutdown() error
}

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
    // Name is the unique identifier for the plugin
    Name string
    
    // Version is the semantic version of the plugin
    Version string
    
    // Description explains what the plugin does
    Description string
    
    // Author identifies the plugin developer
    Author string
    
    // Website is the URL to the plugin's homepage
    Website string
    
    // MinNoideaVersion is the minimum compatible NoIdea version
    MinNoideaVersion string
}

// PluginContext provides access to NoIdea's API and services
type PluginContext interface {
    // Logger returns a logger instance for the plugin
    Logger() Logger
    
    // Config returns the plugin's configuration
    Config() PluginConfig
    
    // RegisterHooks allows the plugin to register its hooks
    RegisterHooks(hooks Hooks) error
    
    // GetService returns a shared NoIdea service by name
    GetService(name string) (interface{}, error)
}
```

## Hook Interfaces

Plugins can implement various hook interfaces to extend NoIdea functionality:

```go
// Hooks is a collection of all available hook types
type Hooks struct {
    Command     CommandHooks
    Commit      CommitHooks
    Feedback    FeedbackHooks
    UI          UIHooks
    Data        DataHooks
    Persistence PersistenceHooks
}

// CommandHooks adds or modifies CLI commands
type CommandHooks interface {
    // AddCommands returns new commands to add to the CLI
    AddCommands() []Command
    
    // ExtendCommand modifies an existing command
    ExtendCommand(name string, extender CommandExtender) error
}

// CommitHooks integrates with the Git commit process
type CommitHooks interface {
    // BeforeCommit runs before a commit is created
    BeforeCommit(ctx CommitContext) error
    
    // AfterCommit runs after a commit is created
    AfterCommit(ctx CommitContext) error
    
    // ModifySuggestion can alter a suggested commit message
    ModifySuggestion(message string, diff string) (string, error)
}

// FeedbackHooks integrates with the Moai feedback system
type FeedbackHooks interface {
    // ProcessFeedback can modify the feedback message
    ProcessFeedback(feedback string, commit string) (string, error)
    
    // AddFeedbackType registers a new feedback template type
    AddFeedbackType(name string, templates []string) error
}

// UIHooks adds custom UI components and styling
type UIHooks interface {
    // BeforeOutput runs before CLI output is displayed
    BeforeOutput(output string) (string, error)
    
    // AfterOutput runs after CLI output is displayed
    AfterOutput(output string) error
    
    // CustomUI creates a custom UI element
    CustomUI(ctx UIContext) error
}

// DataHooks accesses and modifies Git analytics
type DataHooks interface {
    // OnCollectStats is called when Git history is analyzed
    OnCollectStats(stats map[string]interface{}) error
    
    // ProvideAnalytics returns custom analytics data
    ProvideAnalytics() (map[string]interface{}, error)
}
```

## Plugin Configuration

Plugins can access and modify their configuration:

```go
// PluginConfig manages plugin settings
type PluginConfig interface {
    // Get retrieves a configuration value
    Get(key string) (interface{}, error)
    
    // Set stores a configuration value
    Set(key string, value interface{}) error
    
    // Has checks if a configuration key exists
    Has(key string) bool
    
    // Remove deletes a configuration value
    Remove(key string) error
}
```

## Plugin Manifest

Each plugin should include a `plugin.json` manifest file:

```json
{
    "name": "my-awesome-plugin",
    "version": "1.0.0",
    "description": "Adds awesome functionality to NoIdea",
    "author": "Your Name",
    "website": "https://example.com/my-plugin",
    "minNoideaVersion": "v0.4.0",
    "main": "plugin.go",
    "hooks": ["command", "feedback"],
    "permissions": ["git", "config"]
}
```

## Plugin Loading

Plugins can be loaded through several methods:

1. **Local Go Plugin**: Embedded Go plugins compiled as shared objects
2. **Script Plugin**: External script files using a scripting language
3. **HTTP Plugin**: Remotely loaded plugins over HTTPS
4. **Container Plugin**: Isolated plugins running in containers

## Error Handling

Plugins should implement proper error handling:

```go
// PluginError provides detailed error information
type PluginError struct {
    // Code is a unique error identifier
    Code string
    
    // Message is a human-readable error message
    Message string
    
    // Details contains additional error information
    Details map[string]interface{}
    
    // Cause is the underlying error
    Cause error
}
```

## Best Practices

1. **Error Handling**: Always return meaningful errors
2. **Performance**: Minimize impact on NoIdea's performance
3. **Isolation**: Don't interfere with other plugins
4. **Resources**: Clean up resources during Shutdown
5. **Configuration**: Use the provided PluginConfig interface
6. **Logging**: Use the provided Logger interface 