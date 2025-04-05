package github

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/AccursedGalaxy/noidea/internal/config"
	"github.com/AccursedGalaxy/noidea/internal/releaseai"
)

// ReleaseManager handles GitHub release operations
type ReleaseManager struct {
	client *Client
	config config.Config
}

// NewReleaseManager creates a new release manager
func NewReleaseManager(config config.Config) (*ReleaseManager, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	return &ReleaseManager{
		client: client,
		config: config,
	}, nil
}

// UpdateReleaseNotes creates or updates GitHub release notes with AI-generated content
func (m *ReleaseManager) UpdateReleaseNotes(tagName string) error {
	// Extract owner and repo from git remote
	owner, repo, err := ExtractRepoInfo("")
	if err != nil {
		return fmt.Errorf("failed to determine repository info: %w", err)
	}

	// Get the previous tag name
	prevTagName, err := getPreviousTag(tagName)
	if err != nil {
		// Not a critical error, we can proceed without previous tag
		prevTagName = ""
	}

	// Get commit messages between tags
	commitMessages, err := getCommitMessagesBetweenTags(prevTagName, tagName)
	if err != nil {
		return fmt.Errorf("failed to get commit messages: %w", err)
	}

	// Generate AI release notes if LLM is enabled
	var releaseNotes string
	if m.config.LLM.Enabled {
		generator, err := releaseai.NewReleaseNotesGenerator(m.config)
		if err != nil {
			// Fallback to basic notes if AI generation fails
			releaseNotes = generateBasicReleaseNotes(tagName, commitMessages)
			fmt.Printf("Warning: Could not initialize AI release notes generator: %s\n", err)
			fmt.Println("Falling back to basic release notes.")
		} else {
			aiNotes, err := generator.GenerateReleaseNotes(tagName, commitMessages, prevTagName)
			if err != nil {
				// Fallback to basic notes if AI generation fails
				releaseNotes = generateBasicReleaseNotes(tagName, commitMessages)
				fmt.Printf("Warning: AI release notes generation failed: %s\n", err)
				fmt.Println("Falling back to basic release notes.")
			} else {
				releaseNotes = aiNotes
			}
		}
	} else {
		// Generate basic release notes if LLM is not enabled
		releaseNotes = generateBasicReleaseNotes(tagName, commitMessages)
	}

	// Check if a release for this tag already exists
	releases, err := m.client.get(fmt.Sprintf("/repos/%s/%s/releases/tags/%s", owner, repo, tagName))
	if err == nil {
		// Release exists, update it
		releaseID, ok := releases["id"].(float64)
		if !ok {
			return fmt.Errorf("failed to extract release ID")
		}

		// Prepare update payload
		payload := map[string]interface{}{
			"body": releaseNotes,
		}

		// Update the release
		_, err = m.client.patch(fmt.Sprintf("/repos/%s/%s/releases/%d", owner, repo, int(releaseID)), payload)
		if err != nil {
			return fmt.Errorf("failed to update release notes: %w", err)
		}

		fmt.Printf("✅ Updated release notes for %s\n", tagName)
		return nil
	}

	// Release doesn't exist, create a new one
	payload := map[string]interface{}{
		"tag_name": tagName,
		"name":     fmt.Sprintf("Release %s", tagName),
		"body":     releaseNotes,
		"draft":    false,
	}

	_, err = m.client.post(fmt.Sprintf("/repos/%s/%s/releases", owner, repo), payload)
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	fmt.Printf("✅ Created release for %s with enhanced notes\n", tagName)
	return nil
}

// getPreviousTag returns the tag before the specified tag
func getPreviousTag(tag string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0", tag+"^")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getCommitMessagesBetweenTags returns commit messages between two tags
func getCommitMessagesBetweenTags(prevTag, currentTag string) ([]string, error) {
	var cmd *exec.Cmd
	if prevTag == "" {
		// If there's no previous tag, get all commits up to the current tag
		cmd = exec.Command("git", "log", "--pretty=format:%s", currentTag)
	} else {
		// Get commit messages between previous tag and current tag
		cmd = exec.Command("git", "log", "--pretty=format:%s", prevTag+".."+currentTag)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split output into lines
	messages := strings.Split(strings.TrimSpace(string(output)), "\n")
	return messages, nil
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
