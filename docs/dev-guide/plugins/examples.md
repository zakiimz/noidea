# Plugin Examples

This document provides practical examples for developing plugins for NoIdea. These examples demonstrate how to implement common plugin types and use the plugin interfaces effectively.

## Basic Plugin Structure

Every NoIdea plugin follows this basic structure:

```go
package main

import (
    "github.com/AccursedGalaxy/noidea/internal/plugin"
)

// MyPlugin is a basic NoIdea plugin
type MyPlugin struct {
    ctx    plugin.PluginContext
    logger plugin.Logger
    config plugin.PluginConfig
}

// Info returns plugin metadata
func (p *MyPlugin) Info() plugin.PluginInfo {
    return plugin.PluginInfo{
        Name:            "my-plugin",
        Version:         "1.0.0",
        Description:     "My first NoIdea plugin",
        Author:          "Your Name",
        Website:         "https://example.com/my-plugin",
        MinNoideaVersion: "v0.4.0",
    }
}

// Initialize sets up the plugin
func (p *MyPlugin) Initialize(ctx plugin.PluginContext) error {
    p.ctx = ctx
    p.logger = ctx.Logger()
    p.config = ctx.Config()
    
    // Register hooks
    hooks := plugin.Hooks{
        // Add your hooks here
    }
    
    return ctx.RegisterHooks(hooks)
}

// Shutdown performs cleanup
func (p *MyPlugin) Shutdown() error {
    p.logger.Info("Plugin shutting down")
    return nil
}

// Plugin entry point
func CreatePlugin() plugin.Plugin {
    return &MyPlugin{}
}
```

## Command Hook Example

This example demonstrates how to add a new command to NoIdea:

```go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/AccursedGalaxy/noidea/internal/plugin"
)

// CommandPlugin adds custom commands
type CommandPlugin struct {
    ctx plugin.PluginContext
}

// Info returns plugin metadata
func (p *CommandPlugin) Info() plugin.PluginInfo {
    return plugin.PluginInfo{
        Name:            "custom-command",
        Version:         "1.0.0",
        Description:     "Adds custom commands to NoIdea",
        Author:          "Your Name",
        MinNoideaVersion: "v0.4.0",
    }
}

// Initialize sets up the plugin
func (p *CommandPlugin) Initialize(ctx plugin.PluginContext) error {
    p.ctx = ctx
    
    // Register command hooks
    hooks := plugin.Hooks{
        Command: &CommandHooks{},
    }
    
    return ctx.RegisterHooks(hooks)
}

// Shutdown performs cleanup
func (p *CommandPlugin) Shutdown() error {
    return nil
}

// CommandHooks implements the CommandHooks interface
type CommandHooks struct{}

// AddCommands returns custom commands
func (h *CommandHooks) AddCommands() []plugin.Command {
    statsCmd := &cobra.Command{
        Use:   "custom-stats",
        Short: "Display custom Git statistics",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Custom Git statistics:")
            fmt.Println("- Lines added: 1,024")
            fmt.Println("- Lines removed: 512")
            fmt.Println("- Commits: 48")
        },
    }
    
    return []plugin.Command{statsCmd}
}

// ExtendCommand modifies existing commands
func (h *CommandHooks) ExtendCommand(name string, extender plugin.CommandExtender) error {
    return nil // Not implementing command extension in this example
}

// Plugin entry point
func CreatePlugin() plugin.Plugin {
    return &CommandPlugin{}
}
```

## Feedback Hook Example

This example shows how to modify the Moai feedback system:

```go
package main

import (
    "fmt"
    "strings"
    "github.com/AccursedGalaxy/noidea/internal/plugin"
)

// FeedbackPlugin customizes the feedback system
type FeedbackPlugin struct {
    ctx plugin.PluginContext
}

// Info returns plugin metadata
func (p *FeedbackPlugin) Info() plugin.PluginInfo {
    return plugin.PluginInfo{
        Name:            "pirate-feedback",
        Version:         "1.0.0",
        Description:     "Adds pirate-themed feedback",
        Author:          "Your Name",
        MinNoideaVersion: "v0.4.0",
    }
}

// Initialize sets up the plugin
func (p *FeedbackPlugin) Initialize(ctx plugin.PluginContext) error {
    p.ctx = ctx
    
    // Register feedback hooks
    hooks := plugin.Hooks{
        Feedback: &PirateFeedbackHooks{},
    }
    
    return ctx.RegisterHooks(hooks)
}

// Shutdown performs cleanup
func (p *FeedbackPlugin) Shutdown() error {
    return nil
}

// PirateFeedbackHooks implements the FeedbackHooks interface
type PirateFeedbackHooks struct{}

// ProcessFeedback modifies feedback messages
func (h *PirateFeedbackHooks) ProcessFeedback(feedback string, commit string) (string, error) {
    // Convert feedback to pirate speak
    pirateFeedback := feedback
    
    // Replace common words with pirate equivalents
    replacements := map[string]string{
        "Hello": "Ahoy",
        "yes": "aye",
        "you": "ye",
        "your": "yer",
        "is": "be",
        "are": "be",
        "the": "th'",
        "great": "mighty fine",
        "good": "shipshape",
        "bad": "scurvy",
    }
    
    for word, replacement := range replacements {
        pirateFeedback = strings.ReplaceAll(pirateFeedback, word, replacement)
    }
    
    // Add pirate suffix
    pirateFeedback += " Arrr! 🏴‍☠️"
    
    return pirateFeedback, nil
}

// AddFeedbackType registers pirate feedback templates
func (h *PirateFeedbackHooks) AddFeedbackType(name string, templates []string) error {
    if name == "pirate" {
        pirateTemplates := []string{
            "Shiver me timbers! What kind of code be this?",
            "Blimey! This commit be worthy of me treasure chest!",
            "Ye code like a landlubber! Walk the plank!",
            "This commit be a fine addition to the ship's log!",
            "Avast! This code be needin' more rum!",
        }
        
        templates = pirateTemplates
        return nil
    }
    
    return fmt.Errorf("unknown feedback type: %s", name)
}

// Plugin entry point
func CreatePlugin() plugin.Plugin {
    return &FeedbackPlugin{}
}
```

## Data Hook Example

This example demonstrates how to access and modify Git analytics:

```go
package main

import (
    "time"
    "github.com/AccursedGalaxy/noidea/internal/plugin"
)

// AnalyticsPlugin adds custom analytics
type AnalyticsPlugin struct {
    ctx plugin.PluginContext
}

// Info returns plugin metadata
func (p *AnalyticsPlugin) Info() plugin.PluginInfo {
    return plugin.PluginInfo{
        Name:            "time-analytics",
        Version:         "1.0.0",
        Description:     "Adds time-based commit analytics",
        Author:          "Your Name",
        MinNoideaVersion: "v0.4.0",
    }
}

// Initialize sets up the plugin
func (p *AnalyticsPlugin) Initialize(ctx plugin.PluginContext) error {
    p.ctx = ctx
    
    // Register data hooks
    hooks := plugin.Hooks{
        Data: &TimeAnalyticsHooks{},
    }
    
    return ctx.RegisterHooks(hooks)
}

// Shutdown performs cleanup
func (p *AnalyticsPlugin) Shutdown() error {
    return nil
}

// TimeAnalyticsHooks implements the DataHooks interface
type TimeAnalyticsHooks struct{
    morningCommits   int
    afternoonCommits int
    eveningCommits   int
    nightCommits     int
    weekdayCommits   int
    weekendCommits   int
}

// OnCollectStats analyzes commit times
func (h *TimeAnalyticsHooks) OnCollectStats(stats map[string]interface{}) error {
    // Extract timestamps from commits
    if commits, ok := stats["commits"].([]interface{}); ok {
        for _, c := range commits {
            if commit, ok := c.(map[string]interface{}); ok {
                if timestamp, ok := commit["timestamp"].(time.Time); ok {
                    // Analyze commit time
                    hour := timestamp.Hour()
                    weekday := timestamp.Weekday()
                    
                    // Time of day analysis
                    switch {
                    case hour >= 5 && hour < 12:
                        h.morningCommits++
                    case hour >= 12 && hour < 17:
                        h.afternoonCommits++
                    case hour >= 17 && hour < 22:
                        h.eveningCommits++
                    default:
                        h.nightCommits++
                    }
                    
                    // Day of week analysis
                    if weekday == time.Saturday || weekday == time.Sunday {
                        h.weekendCommits++
                    } else {
                        h.weekdayCommits++
                    }
                }
            }
        }
    }
    
    return nil
}

// ProvideAnalytics returns time-based analytics
func (h *TimeAnalyticsHooks) ProvideAnalytics() (map[string]interface{}, error) {
    return map[string]interface{}{
        "time_analytics": map[string]interface{}{
            "morning_commits":   h.morningCommits,
            "afternoon_commits": h.afternoonCommits,
            "evening_commits":   h.eveningCommits,
            "night_commits":     h.nightCommits,
            "weekday_commits":   h.weekdayCommits,
            "weekend_commits":   h.weekendCommits,
        },
    }, nil
}

// Plugin entry point
func CreatePlugin() plugin.Plugin {
    return &AnalyticsPlugin{}
}
```

## UI Hook Example

This example shows how to customize NoIdea's UI:

```go
package main

import (
    "strings"
    "github.com/fatih/color"
    "github.com/AccursedGalaxy/noidea/internal/plugin"
)

// CustomUIPlugin enhances NoIdea's UI
type CustomUIPlugin struct {
    ctx plugin.PluginContext
}

// Info returns plugin metadata
func (p *CustomUIPlugin) Info() plugin.PluginInfo {
    return plugin.PluginInfo{
        Name:            "rainbow-ui",
        Version:         "1.0.0",
        Description:     "Adds colorful UI elements",
        Author:          "Your Name",
        MinNoideaVersion: "v0.4.0",
    }
}

// Initialize sets up the plugin
func (p *CustomUIPlugin) Initialize(ctx plugin.PluginContext) error {
    p.ctx = ctx
    
    // Register UI hooks
    hooks := plugin.Hooks{
        UI: &RainbowUIHooks{},
    }
    
    return ctx.RegisterHooks(hooks)
}

// Shutdown performs cleanup
func (p *CustomUIPlugin) Shutdown() error {
    return nil
}

// RainbowUIHooks implements the UIHooks interface
type RainbowUIHooks struct{}

// BeforeOutput modifies CLI output
func (h *RainbowUIHooks) BeforeOutput(output string) (string, error) {
    // Add rainbow dividers to the output
    if strings.Contains(output, "-----") {
        divider := color.New(color.FgRed).Sprint("❤️ ") +
                  color.New(color.FgYellow).Sprint("💛 ") +
                  color.New(color.FgGreen).Sprint("💚 ") +
                  color.New(color.FgBlue).Sprint("💙 ") +
                  color.New(color.FgMagenta).Sprint("💜 ")
        
        // Repeat the pattern to match divider length
        fullDivider := strings.Repeat(divider, 3)
        
        // Replace all dividers with rainbow dividers
        output = strings.ReplaceAll(output, "------------------------------------------------------", fullDivider)
    }
    
    return output, nil
}

// AfterOutput runs after CLI output is displayed
func (h *RainbowUIHooks) AfterOutput(output string) error {
    // No post-processing needed
    return nil
}

// CustomUI creates a custom UI element
func (h *RainbowUIHooks) CustomUI(ctx plugin.UIContext) error {
    // Not implementing custom UI elements in this example
    return nil
}

// Plugin entry point
func CreatePlugin() plugin.Plugin {
    return &CustomUIPlugin{}
}
```

## Commit Hook Example

This example shows how to integrate with the Git commit process:

```go
package main

import (
    "regexp"
    "strings"
    "github.com/AccursedGalaxy/noidea/internal/plugin"
)

// CommitRulesPlugin enforces commit message rules
type CommitRulesPlugin struct {
    ctx plugin.PluginContext
}

// Info returns plugin metadata
func (p *CommitRulesPlugin) Info() plugin.PluginInfo {
    return plugin.PluginInfo{
        Name:            "commit-rules",
        Version:         "1.0.0",
        Description:     "Enforces commit message rules",
        Author:          "Your Name",
        MinNoideaVersion: "v0.4.0",
    }
}

// Initialize sets up the plugin
func (p *CommitRulesPlugin) Initialize(ctx plugin.PluginContext) error {
    p.ctx = ctx
    
    // Register commit hooks
    hooks := plugin.Hooks{
        Commit: &CommitRulesHooks{logger: ctx.Logger()},
    }
    
    return ctx.RegisterHooks(hooks)
}

// Shutdown performs cleanup
func (p *CommitRulesPlugin) Shutdown() error {
    return nil
}

// CommitRulesHooks implements the CommitHooks interface
type CommitRulesHooks struct{
    logger plugin.Logger
}

// BeforeCommit validates commit messages
func (h *CommitRulesHooks) BeforeCommit(ctx plugin.CommitContext) error {
    // Check if commit message follows conventional format
    message := ctx.Message
    
    // Conventional commit regex pattern
    pattern := `^(feat|fix|docs|style|refactor|test|chore)(\([a-z0-9-]+\))?: .+`
    match, err := regexp.MatchString(pattern, message)
    if err != nil {
        return err
    }
    
    if !match {
        h.logger.Warn("Commit message does not follow conventional format")
        // We don't block the commit, just warn
    }
    
    return nil
}

// AfterCommit runs after a commit is created
func (h *CommitRulesHooks) AfterCommit(ctx plugin.CommitContext) error {
    // No post-commit processing needed
    return nil
}

// ModifySuggestion enhances the commit message
func (h *CommitRulesHooks) ModifySuggestion(message string, diff string) (string, error) {
    // Add issue reference if found in branch name
    // Example: If branch is "feature/PROJ-123-new-feature", add "Refs: PROJ-123"
    
    // Get current branch name
    // Note: In a real plugin, you'd get this from git
    branchName := "feature/PROJ-123-new-feature" // Example value
    
    // Extract issue ID with regex
    issuePattern := `(([A-Z]+)-\d+)`
    re := regexp.MustCompile(issuePattern)
    matches := re.FindStringSubmatch(branchName)
    
    if len(matches) > 1 {
        issueID := matches[1]
        
        // Check if issue ID is already in message
        if !strings.Contains(message, issueID) {
            // Add reference to the end of the first line
            lines := strings.SplitN(message, "\n", 2)
            
            if len(lines) == 1 {
                message = lines[0] + " (Refs: " + issueID + ")"
            } else {
                message = lines[0] + " (Refs: " + issueID + ")\n" + lines[1]
            }
        }
    }
    
    return message, nil
}

// Plugin entry point
func CreatePlugin() plugin.Plugin {
    return &CommitRulesPlugin{}
}
```

## Best Practices and Tips

1. **Error Handling**: Always check errors and provide meaningful error messages
2. **Configuration**: Make your plugin configurable through the Config interface
3. **Logging**: Use the provided Logger rather than fmt.Println
4. **Testing**: Write tests for your plugin functionality
5. **Documentation**: Add usage instructions and examples in your plugin's README
6. **Dependencies**: Minimize external dependencies to avoid conflicts
7. **Performance**: Be mindful of performance impact, especially for commit hooks
8. **Versioning**: Use semantic versioning for your plugin 