package github

import (
	"fmt"
	"os"
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
func (m *ReleaseManager) UpdateReleaseNotes(tagName string, skipApproval bool) error {
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

	// Get diffs between tags for better context
	diffContent, err := getCodeDiffsBetweenTags(prevTagName, tagName)
	if err != nil {
		fmt.Printf("Warning: Could not get detailed code diffs: %s\n", err)
		// We can continue without diffs, it's not critical
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
			aiNotes, err := generator.GenerateReleaseNotes(tagName, commitMessages, prevTagName, diffContent)
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

	// Get the existing release body if available
	var existingBody string
	if err == nil {
		if body, ok := releases["body"].(string); ok {
			existingBody = body
		}
	}

	// If there's an existing body, extract the changelog
	changelog := extractChangelog(existingBody)

	// Combine our generated notes with the changelog if it exists
	if changelog != "" {
		releaseNotes = combineNotesWithChangelog(releaseNotes, changelog)
	}

	// Show the release notes to the user and ask for approval, unless skipped
	var approvedNotes string
	var approved bool

	if skipApproval {
		// Skip approval process
		approvedNotes = releaseNotes
		approved = true
		fmt.Println("Skipping approval process as requested.")
	} else {
		approvedNotes, approved = showAndApproveReleaseNotes(releaseNotes, tagName)
		if !approved {
			return fmt.Errorf("release notes update cancelled by user")
		}
	}

	// Use the approved notes (which might have been edited)
	releaseNotes = approvedNotes

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

	// Use a more detailed format for commit messages
	// %s = subject, %b = body, %h = abbreviated hash
	commitFormat := "%h %s"

	if prevTag == "" {
		// If there's no previous tag, get all commits up to the current tag
		// Limit to a reasonable number (e.g., 50) to avoid overwhelming output
		cmd = exec.Command("git", "log", "--pretty=format:"+commitFormat, "-n", "50", currentTag)
	} else {
		// Get commit messages between previous tag and current tag
		cmd = exec.Command("git", "log", "--pretty=format:"+commitFormat, prevTag+".."+currentTag)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split output into lines
	raw := string(output)
	if strings.TrimSpace(raw) == "" {
		// If no commits found, try a fallback approach
		return getRecentCommitsForTag(currentTag)
	}

	messages := strings.Split(strings.TrimSpace(raw), "\n")

	// Ensure we have at least some commit messages
	if len(messages) == 0 {
		return getRecentCommitsForTag(currentTag)
	}

	return messages, nil
}

// getRecentCommitsForTag gets recent commits up to a tag as a fallback
func getRecentCommitsForTag(tag string) ([]string, error) {
	// First, get the commit hash for the tag
	hashCmd := exec.Command("git", "rev-list", "-n", "1", tag)
	hashOutput, err := hashCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit hash for tag %s: %w", tag, err)
	}

	tagHash := strings.TrimSpace(string(hashOutput))
	if tagHash == "" {
		return nil, fmt.Errorf("couldn't determine commit hash for tag %s", tag)
	}

	// Get 10 commits leading up to and including the tag commit
	cmd := exec.Command("git", "log", "--pretty=format:%h %s", "-n", "10", tagHash)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent commits: %w", err)
	}

	messages := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(messages) == 0 || messages[0] == "" {
		// If still no commits, just return a placeholder message
		return []string{"Initial release or no commit history found"}, nil
	}

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

// getCodeDiffsBetweenTags returns a summary of code changes between two tags
func getCodeDiffsBetweenTags(prevTag, currentTag string) (string, error) {
	var cmd *exec.Cmd

	if prevTag == "" {
		// If no previous tag, just get stats for the current tag
		cmd = exec.Command("git", "show", "--stat", currentTag)
	} else {
		// Get shortened diff between tags
		cmd = exec.Command("git", "diff", "--stat", prevTag, currentTag)
	}

	statOutput, _ := cmd.Output()

	// Get a subset of actual diffs (limiting to avoid huge output)
	if prevTag == "" {
		cmd = exec.Command("git", "show", "--color=never", "--patch", "--unified=1", currentTag)
	} else {
		cmd = exec.Command("git", "diff", "--color=never", "--patch", "--unified=1",
			"--diff-filter=AM", // Only Added and Modified files
			"--no-prefix", prevTag, currentTag)
	}

	diffOutput, err := cmd.Output()
	if err != nil {
		return string(statOutput), nil // Return just stats if full diff fails
	}

	// Combine stats and limited diff output
	result := string(statOutput) + "\n" + limitDiffOutput(string(diffOutput))

	return result, nil
}

// limitDiffOutput truncates diff output to a reasonable size
func limitDiffOutput(diff string) string {
	lines := strings.Split(diff, "\n")

	// If the diff is small enough, return it entirely
	if len(lines) < 200 {
		return diff
	}

	// Otherwise, take the first ~150 lines with context
	var result []string
	result = append(result, lines[:150]...)
	result = append(result, "... [diff truncated for brevity] ...")

	return strings.Join(result, "\n")
}

// extractChangelog extracts the auto-generated GitHub changelog from release notes
func extractChangelog(notes string) string {
	if notes == "" {
		return ""
	}

	// Try to find the GitHub-generated changelog section
	changelogMarkers := []string{
		"## What's Changed",
		"**Full Changelog**",
	}

	for _, marker := range changelogMarkers {
		index := strings.Index(notes, marker)
		if index >= 0 {
			return notes[index:]
		}
	}

	return ""
}

// combineNotesWithChangelog combines AI-generated notes with GitHub changelog
func combineNotesWithChangelog(notes, changelog string) string {
	if changelog == "" {
		return notes
	}

	// Ensure there's a separator between our notes and the changelog
	return notes + "\n\n---\n\n" + changelog
}

// showAndApproveReleaseNotes shows the release notes to the user and asks for approval
func showAndApproveReleaseNotes(notes, tag string) (string, bool) {
	fmt.Println("\n==== Generated Release Notes for", tag, "====")
	fmt.Println(notes)
	fmt.Println("============================================")

	// Ask if user wants to approve, edit, or cancel
	fmt.Print("\nWould you like to: [a]pprove, [e]dit, or [c]ancel? ")
	var input string
	fmt.Scanln(&input)

	input = strings.ToLower(strings.TrimSpace(input))

	if input == "a" || input == "approve" {
		return notes, true
	} else if input == "c" || input == "cancel" {
		return "", false
	} else if input == "e" || input == "edit" {
		// Create a temp file with the notes
		tmpFile, err := os.CreateTemp("", "release-notes-*.md")
		if err != nil {
			fmt.Printf("Error creating temporary file: %s\n", err)
			fmt.Print("Do you still want to approve the unedited notes? [y/n] ")
			fmt.Scanln(&input)
			return notes, strings.ToLower(strings.TrimSpace(input)) == "y"
		}

		// Write notes to the temp file
		_, err = tmpFile.WriteString(notes)
		tmpFile.Close()
		if err != nil {
			fmt.Printf("Error writing to temporary file: %s\n", err)
			fmt.Print("Do you still want to approve the unedited notes? [y/n] ")
			fmt.Scanln(&input)
			return notes, strings.ToLower(strings.TrimSpace(input)) == "y"
		}

		// Open the editor
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano" // Fallback editor
		}

		cmd := exec.Command(editor, tmpFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error opening editor: %s\n", err)
			fmt.Print("Do you still want to approve the unedited notes? [y/n] ")
			fmt.Scanln(&input)
			return notes, strings.ToLower(strings.TrimSpace(input)) == "y"
		}

		// Read the edited content
		editedContent, err := os.ReadFile(tmpFile.Name())
		if err != nil {
			fmt.Printf("Error reading edited file: %s\n", err)
			fmt.Print("Do you still want to approve the unedited notes? [y/n] ")
			fmt.Scanln(&input)
			return notes, strings.ToLower(strings.TrimSpace(input)) == "y"
		}

		// Clean up
		os.Remove(tmpFile.Name())

		fmt.Println("Release notes edited successfully.")
		return string(editedContent), true
	}

	// Default case - ask again
	fmt.Println("Invalid choice. Please try again.")
	return showAndApproveReleaseNotes(notes, tag)
}
