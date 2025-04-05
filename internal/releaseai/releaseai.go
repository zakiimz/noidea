// Package releaseai provides AI-powered release notes generation
package releaseai

import (
	"fmt"
	"strings"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/feedback"
)

// ReleaseNotesGenerator handles creating AI-enhanced release notes
type ReleaseNotesGenerator struct {
	engine feedback.FeedbackEngine
	config config.Config
}

// NewReleaseNotesGenerator creates a new release notes generator
func NewReleaseNotesGenerator(cfg config.Config) (*ReleaseNotesGenerator, error) {
	// Check if LLM is enabled
	if !cfg.LLM.Enabled {
		return nil, fmt.Errorf("LLM is not enabled in config, cannot generate AI release notes")
	}

	// Initialize feedback engine using the existing unified feedback engine
	engine := feedback.NewFeedbackEngine(
		cfg.LLM.Provider,
		cfg.LLM.Model,
		cfg.LLM.APIKey,
		"professional", // Default to professional personality for release notes
		cfg.Moai.PersonalityFile,
	)

	return &ReleaseNotesGenerator{
		engine: engine,
		config: cfg,
	}, nil
}

// GenerateReleaseNotes creates AI-enhanced release notes from commit messages
func (g *ReleaseNotesGenerator) GenerateReleaseNotes(
	version string,
	commitMessages []string,
	previousVersion string,
) (string, error) {
	// Create context for the LLM
	context := feedback.CommitContext{
		Message:       fmt.Sprintf("Generate release notes for version %s", version),
		CommitHistory: commitMessages,
		CommitStats: map[string]interface{}{
			"version":         version,
			"previousVersion": previousVersion,
		},
	}

	// Use the feedback engine to generate enhanced notes
	prompt := buildReleaseNotesPrompt(version, commitMessages, previousVersion)
	context.Message = prompt

	notes, err := g.engine.GenerateSummaryFeedback(context)
	if err != nil {
		return "", fmt.Errorf("failed to generate release notes: %w", err)
	}

	return notes, nil
}

// buildReleaseNotesPrompt creates a prompt for the LLM
func buildReleaseNotesPrompt(version string, commitMessages []string, previousVersion string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Generate comprehensive release notes for version %s", version))

	if previousVersion != "" {
		sb.WriteString(fmt.Sprintf(" (previous version was %s)", previousVersion))
	}

	sb.WriteString(".\n\n")
	sb.WriteString("Here are the raw commit messages since the last release:\n\n")

	for _, msg := range commitMessages {
		sb.WriteString("- ")
		sb.WriteString(msg)
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString("Please create organized, user-friendly release notes with the following characteristics:\n")
	sb.WriteString("1. Start with a brief overview summary of the release\n")
	sb.WriteString("2. Group similar changes into sections (Features, Bug Fixes, Improvements, etc.)\n")
	sb.WriteString("3. Use clear, concise language that end-users can understand\n")
	sb.WriteString("4. Format with Markdown for GitHub\n")
	sb.WriteString("5. Focus on the user impact and benefits of each change\n")

	return sb.String()
}
