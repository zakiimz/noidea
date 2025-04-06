# Core Components

This document provides a detailed overview of the core components that make up NoIdea, explaining their purpose, architecture, and interactions.

## Feedback Engine

The feedback engine is the heart of NoIdea's AI capabilities, responsible for generating commit suggestions, feedback, and summaries.

### Architecture

```
┌─────────────────────────────────────────────┐
│                Feedback Engine              │
├─────────────┬───────────────┬───────────────┤
│    Local    │    Unified    │    Provider   │
│   Engine    │    Engine     │   Adapters    │
└─────────────┴───────────────┴───────────────┘
       ▲               ▲               ▲
       │               │               │
       ├───────────────┼───────────────┤
       │               │               │
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ Personality │ │ Commit Data │ │    LLM      │
│   System    │ │  Context    │ │  Providers  │
└─────────────┘ └─────────────┘ └─────────────┘
```

### Key Components

#### Local Engine

The local engine provides feedback without requiring API access. It uses predefined templates and simple pattern matching.

**Key Files:**
- `internal/feedback/local.go`: Implementation of the local feedback engine
- `internal/moai/faces.go`: Moai face definitions and templates

#### Unified Engine

The unified engine is responsible for interacting with LLM providers like xAI, OpenAI, and DeepSeek.

**Key Files:**
- `internal/feedback/unified.go`: Unified API for different LLM providers
- `internal/feedback/engine.go`: Common engine interface definitions

#### Personality System

The personality system manages different AI personas, allowing customized feedback styles.

**Key Files:**
- `internal/personality/personalities.go`: Personality definition and loading
- `personalities.toml.example`: Example personality configuration

## Configuration System

The configuration system manages user preferences, API credentials, and behavior settings.

### Architecture

```
┌─────────────────────────────────────────────┐
│             Configuration System            │
├─────────────┬───────────────┬───────────────┤
│   Config    │    Secure     │  Environment  │
│   Files     │    Storage    │   Variables   │
└─────────────┴───────────────┴───────────────┘
       ▲               ▲               ▲
       │               │               │
       └───────────────┼───────────────┘
                       │
               ┌───────────────┐
               │ Configuration │
               │    Manager    │
               └───────────────┘
                       ▲
                       │
               ┌───────────────┐
               │  Application  │
               │    Components │
               └───────────────┘
```

### Key Components

#### Configuration Files

Configuration files store user preferences in TOML format.

**Key Files:**
- `internal/config/config.go`: Configuration loading and parsing
- `internal/config/default.go`: Default configuration values

#### Secure Storage

Secure storage manages sensitive information like API keys.

**Key Files:**
- `internal/secure/keyring.go`: Secure credential storage
- `internal/secure/apikey.go`: API key validation and management

## Command System

The command system provides the CLI interface for NoIdea, built on the Cobra library.

### Architecture

```
┌─────────────────────────────────────────────┐
│                Command System               │
├─────────────┬───────────────┬───────────────┤
│    Root     │   Feature     │   Utility     │
│   Command   │   Commands    │   Commands    │
└─────────────┴───────────────┴───────────────┘
       ▲               ▲               ▲
       │               │               │
       └───────────────┼───────────────┘
                       │
               ┌───────────────┐
               │    Cobra      │
               │   Framework   │
               └───────────────┘
```

### Key Components

#### Root Command

The root command is the entry point for all CLI operations.

**Key Files:**
- `cmd/root.go`: Root command definition and global flags

#### Feature Commands

Feature commands implement core NoIdea functionality.

**Key Files:**
- `cmd/suggest.go`: Commit suggestion command
- `cmd/moai.go`: Feedback command
- `cmd/summary.go`: Summary generation command

#### Utility Commands

Utility commands provide supporting functionality.

**Key Files:**
- `cmd/config.go`: Configuration management
- `cmd/init.go`: Repository initialization
- `cmd/update.go`: Self-update functionality

## Git Integration

The Git integration system interacts with Git repositories to analyze changes and history.

### Architecture

```
┌─────────────────────────────────────────────┐
│               Git Integration               │
├─────────────┬───────────────┬───────────────┤
│   Command   │    History    │     Diff      │
│  Execution  │    Analysis   │    Parsing    │
└─────────────┴───────────────┴───────────────┘
       ▲               ▲               ▲
       │               │               │
       └───────────────┼───────────────┘
                       │
               ┌───────────────┐
               │  Application  │
               │    Logic      │
               └───────────────┘
```

### Key Components

#### Command Execution

Executes Git commands and processes results.

**Key Files:**
- `internal/git/git.go`: Git command execution
- `internal/git/repo.go`: Repository interaction

#### History Analysis

Analyzes Git history for patterns and insights.

**Key Files:**
- `internal/history/collector.go`: Gathers commit history data
- `internal/history/analysis.go`: Analyzes commit patterns

## GitHub Integration

The GitHub integration system interacts with the GitHub API for managing releases and workflows.

### Architecture

```
┌─────────────────────────────────────────────┐
│              GitHub Integration             │
├─────────────┬───────────────┬───────────────┤
│     Auth     │    Release    │   Workflow   │
│    Manager   │   Management  │    Status    │
└─────────────┴───────────────┴───────────────┘
       ▲               ▲               ▲
       │               │               │
       └───────────────┼───────────────┘
                       │
               ┌───────────────┐
               │  GitHub API   │
               │    Client     │
               └───────────────┘
```

### Key Components

#### Auth Manager

Manages GitHub authentication and credentials.

**Key Files:**
- `internal/github/auth.go`: Authentication handling
- `cmd/github.go`: GitHub command implementation

#### Release Management

Manages GitHub releases and release notes.

**Key Files:**
- `internal/github/release.go`: Release creation and management
- `internal/releaseai/generator.go`: AI-enhanced release notes

## Plugin System (Future)

The plugin system will allow extending NoIdea with custom functionality.

### Architecture

```
┌─────────────────────────────────────────────┐
│                Plugin System                │
├─────────────┬───────────────┬───────────────┤
│   Plugin    │    Plugin     │    Hook       │
│  Registry   │    Loader     │    Points     │
└─────────────┴───────────────┴───────────────┘
       ▲               ▲               ▲
       │               │               │
       └───────────────┼───────────────┘
                       │
               ┌───────────────┐
               │    Plugin     │
               │  Interfaces   │
               └───────────────┘
                       ▲
                       │
               ┌───────────────┐
               │   External    │
               │    Plugins    │
               └───────────────┘
```

### Key Components

#### Plugin Registry

Manages plugin registration and discovery.

**Future Files:**
- `internal/plugin/registry.go`: Plugin registration and management

#### Plugin Loader

Loads plugins from various sources.

**Future Files:**
- `internal/plugin/loader.go`: Plugin loading mechanisms

#### Hook Points

Defines extension points for plugins.

**Future Files:**
- `internal/plugin/hooks.go`: Hook point definitions
- `internal/plugin/events.go`: Event system for plugins

## Interaction Between Components

### Commit Suggestion Flow

```
┌──────────┐     ┌───────────┐     ┌─────────────┐     ┌───────────┐
│  Suggest │────▶│   Git     │────▶│   Feedback  │────▶│   Output  │
│  Command │     │  System   │     │    Engine   │     │  Renderer │
└──────────┘     └───────────┘     └─────────────┘     └───────────┘
      │                │                 │                   │
      │                │                 │                   │
      ▼                ▼                 ▼                   ▼
┌──────────┐     ┌───────────┐     ┌─────────────┐     ┌───────────┐
│  Config  │     │  History  │     │ Personality │     │  Terminal │
│  System  │     │  Analysis │     │   System    │     │   Output  │
└──────────┘     └───────────┘     └─────────────┘     └───────────┘
```

### GitHub Release Flow

```
┌──────────┐     ┌───────────┐     ┌─────────────┐     ┌───────────┐
│  GitHub  │────▶│   Auth    │────▶│   Release   │────▶│ ReleaseAI │
│  Command │     │  Manager  │     │   Manager   │     │ Generator │
└──────────┘     └───────────┘     └─────────────┘     └───────────┘
      │                │                 │                   │
      │                │                 │                   │
      ▼                ▼                 ▼                   ▼
┌──────────┐     ┌───────────┐     ┌─────────────┐     ┌───────────┐
│  Config  │     │  Secure   │     │   GitHub    │     │   Output  │
│  System  │     │  Storage  │     │     API     │     │  Renderer │
└──────────┘     └───────────┘     └─────────────┘     └───────────┘
```

## Future Development

Future development of core components will focus on:

1. **Modularization**: Further separating components for improved maintainability
2. **Plugin Support**: Adding robust plugin infrastructure
3. **API Stability**: Creating stable APIs for plugin developers
4. **Performance**: Optimizing performance for large repositories
5. **Additional Providers**: Supporting more LLM providers

These core components work together to create a cohesive system that provides intelligent Git assistance while maintaining a modular, extensible architecture. 