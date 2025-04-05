// Package releaseai provides AI-powered release notes generation
package releaseai

import (
	"fmt"
	"strings"

	"github.com/AccursedGalaxy/noidea/internal/config"
)

// ReleaseNotesGenerator handles creating AI-enhanced release notes
type ReleaseNotesGenerator struct {
	// Using direct client instead of feedback engine to avoid pattern confusion
	directClient *DirectLLMClient
	config       config.Config
}

// NewReleaseNotesGenerator creates a new release notes generator
func NewReleaseNotesGenerator(cfg config.Config) (*ReleaseNotesGenerator, error) {
	// Check if LLM is enabled
	if !cfg.LLM.Enabled {
		return nil, fmt.Errorf("LLM is not enabled in config, cannot generate AI release notes")
	}

	// Create a direct client specifically for release notes
	directClient := NewDirectLLMClient(
		cfg.LLM.Provider,
		cfg.LLM.Model,
		cfg.LLM.APIKey,
		cfg.LLM.Temperature,
	)

	// Set a custom system prompt specifically for release notes
	directClient.SetSystemPrompt(`You are a professional software release notes writer.
Your task is to describe changes, features, and fixes in a clear, organized manner.
IMPORTANT RULES:
1. NEVER analyze commit message patterns or quality
2. Focus ONLY on actual software changes and features
3. Begin with a clear overview of key changes
4. Organize changes into relevant categories
5. Keep the tone professional and factual
6. Remove any sections if there are no relevant changes
7. Do not introduce yourself or explain your role`)

	return &ReleaseNotesGenerator{
		directClient: directClient,
		config:       cfg,
	}, nil
}

// GenerateReleaseNotes creates AI-enhanced release notes from commit messages
func (g *ReleaseNotesGenerator) GenerateReleaseNotes(
	version string,
	commitMessages []string,
	previousVersion string,
	diffContent string,
) (string, error) {
	// Check if we have enough data
	if len(commitMessages) == 0 {
		return generateBasicReleaseNotes(version, []string{"Version update"}), nil
	}

	// Build specialized prompt for release notes
	prompt := buildReleaseNotesPrompt(version, commitMessages, previousVersion, diffContent)

	// Use direct LLM client for generation (separate from feedback system)
	notes, err := g.directClient.GenerateReleaseNotes(prompt, 3) // Try up to 3 times
	if err != nil {
		fmt.Printf("Warning: Direct release notes generation failed: %s\n", err)
		return generateBasicReleaseNotes(version, commitMessages), nil
	}

	// Clean up the response and check if it's usable
	notes = cleanReleaseNotes(notes)

	// Fallback to basic notes if we got nothing useful
	if strings.TrimSpace(notes) == "" {
		return generateBasicReleaseNotes(version, commitMessages), nil
	}

	// Make sure we have a proper release title
	if !strings.Contains(notes, "# Release") && !strings.Contains(notes, "#Release") {
		notes = "# Release " + version + "\n\n" + notes
	}

	return notes, nil
}

// cleanReleaseNotes removes any self-introduction or AI mentions from notes
func cleanReleaseNotes(notes string) string {
	if notes == "" {
		return ""
	}

	// Check for common meta-analysis patterns and remove entire sections
	metaPatterns := []string{
		"Commit Message Patterns",
		"Commit Quality",
		"Analysis of commit",
		"Analysis:",
		"Message analysis",
		"Analyzing your commit",
		"1. Commit Message",
		"2. Commit Message",
		"# 1. Commit",
		"# Analysis",
		"# Overview of Commit",
		"Looking at your commit",
		"Based on the commit patterns",
	}

	for _, pattern := range metaPatterns {
		if idx := strings.Index(strings.ToLower(notes), strings.ToLower(pattern)); idx >= 0 {
			// Look for the next section header to truncate
			afterIdx := strings.Index(notes[idx:], "# ")
			if afterIdx > 0 {
				// There's another section after the meta section, remove just the meta section
				notes = notes[:idx] + notes[idx+afterIdx:]
			} else {
				// This meta section might be the start of the notes - check if anything comes before it
				if idx == 0 || !strings.Contains(notes[:idx], "# ") {
					// Find the first real header after this point
					nextHeaderIdx := strings.Index(notes, "# ")
					if nextHeaderIdx > 0 {
						notes = notes[nextHeaderIdx:]
					} else {
						// No section headers, delete everything (will fall back to basic notes)
						return ""
					}
				}
			}
		}
	}

	// Check for common self-introduction patterns and remove them
	introPatterns := []string{
		"Hello!", "I'm an AI", "As an AI", "I'd be happy to",
		"Hello, I am", "Hi there", "Here's the",
		"Below is", "I'll create", "I've generated", "I'll generate",
	}

	for _, pattern := range introPatterns {
		if strings.Contains(notes, pattern) {
			// Find the first markdown header after the intro text
			headerIndex := strings.Index(notes, "# ")
			if headerIndex > 0 {
				notes = notes[headerIndex:]
				break
			} else {
				// Try to find the first line break that might signify the end of an intro
				parts := strings.SplitN(notes, "\n\n", 2)
				if len(parts) > 1 {
					notes = parts[1]
					break
				}
			}
		}
	}

	// Make sure the release notes don't have placeholder text
	if containsPlaceholderText(notes) {
		// The model returned our template with placeholders - generate basic notes instead
		return ""
	}

	return notes
}

// generateBasicReleaseNotes creates a simple release notes from commit messages
func generateBasicReleaseNotes(version string, commitMessages []string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Release %s\n\n", version))
	sb.WriteString("## Changes\n\n")

	for _, msg := range commitMessages {
		sb.WriteString("- ")
		sb.WriteString(msg)
		sb.WriteString("\n")
	}

	return sb.String()
}

// buildReleaseNotesPrompt creates a prompt for the LLM
func buildReleaseNotesPrompt(version string, commitMessages []string, previousVersion string, diffContent string) string {
	var sb strings.Builder

	// Very explicit instructions with template format - emphasizing NOT to analyze commits
	sb.WriteString("IMPORTANT: DO NOT ANALYZE COMMIT MESSAGES OR PATTERNS! This is critical.\n\n")
	sb.WriteString("You are writing release notes for version ")
	sb.WriteString(version)

	if previousVersion != "" {
		sb.WriteString(fmt.Sprintf(" (compared to previous version %s)", previousVersion))
	}

	sb.WriteString(".\n\n")
	sb.WriteString("Based on the following commits and changes, create user-focused release notes:\n\n")

	for _, msg := range commitMessages {
		sb.WriteString("- ")
		sb.WriteString(msg)
		sb.WriteString("\n")
	}

	// Add diff content if available
	if diffContent != "" {
		sb.WriteString("\n\nCode changes summary:\n\n```diff\n")
		// Limit diff content to avoid overwhelming the model
		if len(diffContent) > 2000 {
			sb.WriteString(diffContent[:2000])
			sb.WriteString("\n... [additional changes truncated] ...\n")
		} else {
			sb.WriteString(diffContent)
		}
		sb.WriteString("\n```\n")
	}

	sb.WriteString("\n\nCRITICAL INSTRUCTIONS:\n")
	sb.WriteString("1. DO NOT ANALYZE COMMIT QUALITY OR PATTERNS - THIS IS VERY IMPORTANT!\n")
	sb.WriteString("2. DO NOT introduce yourself or explain what you're doing\n")
	sb.WriteString("3. FOCUS ONLY on the actual software changes and features\n")
	sb.WriteString("4. START DIRECTLY with the release notes\n")
	sb.WriteString("5. DO NOT number sections (like '1. Features')\n")

	sb.WriteString("\n\nOUTPUT FORMAT:\n")
	sb.WriteString("# Release " + version + "\n\n")
	sb.WriteString("## Overview\n")
	sb.WriteString("[Brief summary of key changes]\n\n")
	sb.WriteString("## 🚀 New Features\n")
	sb.WriteString("[New capabilities]\n\n")
	sb.WriteString("## 🔧 Improvements\n")
	sb.WriteString("[Enhancements]\n\n")
	sb.WriteString("## 🐛 Bug Fixes\n")
	sb.WriteString("[Fixed issues]\n\n")

	sb.WriteString("\nRemove any section that has no relevant changes. NEVER analyze commit formats or patterns. Replace placeholder text with actual changes.\n")

	return sb.String()
}

// containsPlaceholderText checks if the generated notes still contain placeholder text
func containsPlaceholderText(notes string) bool {
	placeholders := []string{
		"[Brief summary",
		"[List of",
		"[Add other",
		"[New capabilities",
		"[Enhancements",
		"[Fixed issues",
		"[Brief overview",
		"[description",
		"[placeholder",
	}

	for _, placeholder := range placeholders {
		if strings.Contains(notes, placeholder) {
			return true
		}
	}

	// Check for sections with no content
	emptyPatterns := []string{
		"## Overview\n\n##",
		"## 🚀 New Features\n\n##",
		"## 🔧 Improvements\n\n##",
		"## 🐛 Bug Fixes\n\n##",
	}

	for _, pattern := range emptyPatterns {
		if strings.Contains(notes, pattern) {
			return true
		}
	}

	return false
}
