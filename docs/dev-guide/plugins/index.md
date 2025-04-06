# NoIdea Plugin System

Welcome to the NoIdea plugin system documentation. This section provides comprehensive information about the plugin architecture, interface specifications, and examples for developers who want to create plugins for NoIdea.

## Overview

The NoIdea plugin system allows developers to extend the functionality of the CLI tool with custom commands, feedback mechanisms, UI elements, and more. Plugins are modular, maintainable, and designed to enhance the user experience.

## Documentation Sections

- [Plugin Architecture](architecture.md) - Detailed architecture and design principles
- [Interface Specifications](interfaces.md) - Core interfaces and types for plugin development
- [Plugin Examples](examples.md) - Practical examples for developing various types of plugins

## Quick Start

To create a basic NoIdea plugin:

1. Create a new Go module:
   ```bash
   mkdir my-noidea-plugin
   cd my-noidea-plugin
   go mod init github.com/yourusername/my-noidea-plugin
   ```

2. Add NoIdea as a dependency:
   ```bash
   go get github.com/AccursedGalaxy/noidea
   ```

3. Create a basic plugin structure:
   ```go
   package main

   import (
       "github.com/AccursedGalaxy/noidea/internal/plugin"
   )

   // MyPlugin is a basic NoIdea plugin
   type MyPlugin struct {
       ctx plugin.PluginContext
   }

   // Info returns plugin metadata
   func (p *MyPlugin) Info() plugin.PluginInfo {
       return plugin.PluginInfo{
           Name:            "my-plugin",
           Version:         "1.0.0",
           Description:     "My first NoIdea plugin",
           Author:          "Your Name",
           MinNoideaVersion: "v0.4.0",
       }
   }

   // Initialize sets up the plugin
   func (p *MyPlugin) Initialize(ctx plugin.PluginContext) error {
       p.ctx = ctx
       return nil
   }

   // Shutdown performs cleanup
   func (p *MyPlugin) Shutdown() error {
       return nil
   }

   // Plugin entry point
   func CreatePlugin() plugin.Plugin {
       return &MyPlugin{}
   }
   ```

4. Build your plugin:
   ```bash
   go build -buildmode=plugin -o my-plugin.so
   ```

5. Install your plugin:
   ```bash
   mkdir -p ~/.noidea/plugins
   cp my-plugin.so ~/.noidea/plugins/
   ```

## Plugin Distribution

Plugins can be distributed as:

1. **Shared Object Files (.so)**: For direct installation
2. **Source Code**: For users to build themselves
3. **Plugin Packages**: (Future) For installation via the NoIdea plugin manager

## Community Plugins

We encourage the community to develop and share plugins for NoIdea. When your plugin is ready, consider submitting it to our upcoming plugin directory.

## Getting Help

If you need assistance with plugin development, you can:

- Check the [examples](examples.md) for practical guidance
- Review the [interface specifications](interfaces.md) for technical details
- Join our community discussions on GitHub 