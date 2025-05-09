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

	// Check if a release for this tag already exists
	releases, err := m.client.get(fmt.Sprintf("/repos/%s/%s/releases/tags/%s", owner, repo, tagName))

	// Get the existing release body if available
	var existingBody string
	var releaseID float64 = 0
	if err == nil {
		if body, ok := releases["body"].(string); ok {
			existingBody = body
		}
		// Store release ID for later
		if id, ok := releases["id"].(float64); ok {
			releaseID = id
		}
	}

	// Parse the existing GitHub-generated release notes
	var hasGitHubContent bool
	var overviewSection, whatsChangedSection, summarySection string
	if existingBody != "" {
		overviewSection, whatsChangedSection, summarySection, hasGitHubContent = parseGitHubReleaseNotes(existingBody)
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

	// If there's existing GitHub-generated content and it has our expected structure,
	// we'll only replace the Overview section, preserving the rest
	var releaseNotes string
	if hasGitHubContent {
		// Generate AI content for the overview section only
		overviewContent, err := generateAIOverview(m.config, tagName, commitMessages)
		if err != nil {
			fmt.Printf("Warning: Failed to generate AI overview: %s\n", err)
			// Keep the existing overview if AI generation fails
			if overviewSection == "" {
				overviewSection = "## Overview\n\nThis release includes several improvements and updates."
			}
		} else if overviewContent != "" {
			overviewSection = "## Overview\n\n" + overviewContent
		}

		// Combine the AI-enhanced overview with GitHub's What's Changed and Summary sections
		releaseNotes = overviewSection
		if whatsChangedSection != "" {
			releaseNotes += "\n\n" + whatsChangedSection
		}
		if summarySection != "" {
			releaseNotes += "\n\n" + summarySection
		}

		fmt.Println("Enhanced the Overview section of GitHub-generated release notes.")
	} else {
		// No existing GitHub content with expected structure, generate complete notes
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

		// If there's existing content, but not in our expected format,
		// still try to preserve any GitHub-generated changelog
		if existingBody != "" {
			githubChangelog := extractChangelog(existingBody)
			if githubChangelog != "" {
				releaseNotes = combineNotesWithChangelog(releaseNotes, githubChangelog)
			}
		}
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

	// Check for breaking changes to mark as prerelease if needed
	isBreaking := detectBreakingChanges(commitMessages)

	if releaseID > 0 {
		// Release exists, update it
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
		"tag_name":   tagName,
		"name":       formatReleaseTitle(tagName),
		"body":       releaseNotes,
		"draft":      false,
		"prerelease": isBreaking, // Mark as prerelease if it contains breaking changes
	}

	_, err = m.client.post(fmt.Sprintf("/repos/%s/%s/releases", owner, repo), payload)
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	fmt.Printf("✅ Created release for %s with enhanced notes\n", tagName)
	return nil
}

// parseGitHubReleaseNotes parses GitHub-generated release notes to extract sections
func parseGitHubReleaseNotes(notes string) (overview, whatsChanged, summary string, structured bool) {
	// Check for the expected sections in GitHub's workflow-generated notes
	lines := strings.Split(notes, "\n")

	var currentSection string
	var overviewLines, whatsChangedLines, summaryLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Detect section headers
		if strings.HasPrefix(trimmedLine, "## Overview") {
			currentSection = "overview"
			overviewLines = append(overviewLines, line)
			continue
		} else if strings.HasPrefix(trimmedLine, "## What's Changed") {
			currentSection = "whatsChanged"
			whatsChangedLines = append(whatsChangedLines, line)
			continue
		} else if strings.HasPrefix(trimmedLine, "## Summary") {
			currentSection = "summary"
			summaryLines = append(summaryLines, line)
			continue
		}

		// Add lines to their respective sections
		switch currentSection {
		case "overview":
			overviewLines = append(overviewLines, line)
		case "whatsChanged":
			whatsChangedLines = append(whatsChangedLines, line)
		case "summary":
			summaryLines = append(summaryLines, line)
		}
	}

	// Convert sections back to strings
	if len(overviewLines) > 0 {
		overview = strings.Join(overviewLines, "\n")
	}
	if len(whatsChangedLines) > 0 {
		whatsChanged = strings.Join(whatsChangedLines, "\n")
	}
	if len(summaryLines) > 0 {
		summary = strings.Join(summaryLines, "\n")
	}

	// Determine if we found the expected GitHub-generated structure
	structured = overview != "" && whatsChanged != ""

	return overview, whatsChanged, summary, structured
}

// generateAIOverview generates AI-enhanced content for the overview section only
func generateAIOverview(cfg config.Config, tagName string, commitMessages []string) (string, error) {
	if !cfg.LLM.Enabled {
		return "", fmt.Errorf("LLM not enabled")
	}

	// Create a new generator
	generator, err := releaseai.NewReleaseNotesGenerator(cfg)
	if err != nil {
		return "", err
	}

	// Create a prompt specifically for generating just the overview
	prompt := fmt.Sprintf(`Generate a concise, user-friendly overview paragraph (3-5 sentences) for release %s.
Focus exclusively on explaining the key changes, improvements, or fixes in non-technical, straightforward language.
DO NOT include section headers, bullet points, emojis, or detailed descriptions of individual changes.
Just provide the paragraph text that would go under an "Overview" section.

Based on these commits:
%s`, tagName, strings.Join(commitMessages, "\n"))

	// Generate the overview content
	overview, err := generator.GenerateCustomContent(prompt)
	if err != nil {
		return "", err
	}

	// Clean up the result - remove any accidentally included headers or formatting
	overview = cleanGeneratedOverview(overview)

	return overview, nil
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
// This ensures GitHub's auto-generated changelog content (with links to commits and PRs)
// is preserved when we update the release with AI-generated content
func extractChangelog(notes string) string {
	if notes == "" {
		return ""
	}

	// Try to find the GitHub-generated changelog section
	// GitHub uses several different markers depending on the repo settings
	changelogMarkers := []string{
		"## What's Changed",
		"**Full Changelog**",
		"## Changelog",
		"### Changelog",
		"<!-- Release notes generated using",
		"<details><summary>Changelog</summary>",
	}

	for _, marker := range changelogMarkers {
		index := strings.Index(notes, marker)
		if index >= 0 {
			return notes[index:]
		}
	}

	// If we can't find specific markers, look for common GitHub auto-generated patterns
	// GitHub often includes links to compare views or PR lists
	commonPatterns := []string{
		"compare/",
		"commits/",
		"pull/",
		"Full Changelog:",
	}

	for _, pattern := range commonPatterns {
		if strings.Contains(notes, pattern) {
			// The entire content might be GitHub-generated
			// Try to find a reasonable split point
			lines := strings.Split(notes, "\n")
			for i, line := range lines {
				if strings.Contains(line, pattern) {
					// Return from this line to the end
					return strings.Join(lines[i:], "\n")
				}
			}
		}
	}

	return ""
}

// combineNotesWithChangelog combines AI-generated notes with GitHub changelog
func combineNotesWithChangelog(notes, changelog string) string {
	if changelog == "" {
		return notes
	}

	// Check if the AI-generated notes already contain the changelog
	if strings.Contains(notes, changelog) {
		return notes
	}

	// Ensure there's a clear separator between our notes and the GitHub changelog
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

// formatReleaseTitle formats a release title nicely
func formatReleaseTitle(tagName string) string {
	// If tag name is already properly formatted (e.g., "v1.2.3"), add a prefix
	if strings.HasPrefix(tagName, "v") && len(tagName) > 1 {
		// Check if it looks like a semver tag (vX.Y.Z)
		parts := strings.Split(tagName[1:], ".")
		if len(parts) >= 2 {
			// Return a nicely formatted title
			return fmt.Sprintf("Release %s", tagName)
		}
	}

	// For other formats, just return the tag name
	return tagName
}

// detectBreakingChanges looks for indications of breaking changes in commits
func detectBreakingChanges(commitMessages []string) bool {
	breakingPatterns := []string{
		"BREAKING CHANGE",
		"breaking change",
		"BREAKING-CHANGE",
		"BREAKING_CHANGE",
		"!:",     // Conventional commit breaking change marker
		"feat!:", // Exclamation mark indicates breaking change
		"fix!:",
		"refactor!:",
	}

	for _, msg := range commitMessages {
		for _, pattern := range breakingPatterns {
			if strings.Contains(msg, pattern) {
				return true
			}
		}
	}

	return false
}

// cleanGeneratedOverview removes unwanted formatting from AI-generated overview text
func cleanGeneratedOverview(text string) string {
	// Remove common header prefixes AI might add
	headerPrefixes := []string{
		"# Overview",
		"## Overview",
		"### Overview",
		"Overview:",
		"Overview",
	}

	lines := strings.Split(text, "\n")
	var cleaned []string

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines
		if trimmedLine == "" {
			continue
		}

		// Skip header lines
		isHeader := false
		for _, prefix := range headerPrefixes {
			if strings.HasPrefix(trimmedLine, prefix) {
				isHeader = true
				break
			}
		}

		if isHeader {
			continue
		}

		// Skip bullet points or numbering
		if strings.HasPrefix(trimmedLine, "- ") ||
			strings.HasPrefix(trimmedLine, "* ") ||
			strings.HasPrefix(trimmedLine, "+ ") ||
			(len(trimmedLine) > 2 && trimmedLine[0] >= '0' && trimmedLine[0] <= '9' && trimmedLine[1] == '.') {
			// Extract just the content after the bullet/number
			parts := strings.SplitN(trimmedLine, " ", 2)
			if len(parts) > 1 {
				trimmedLine = parts[1]
			}
		}

		// Add line if it's not a header or empty
		if i == 0 && strings.Contains(trimmedLine, "This release") {
			// Keep the first line as is if it starts a proper description
			cleaned = append(cleaned, trimmedLine)
		} else if i > 0 || len(trimmedLine) > 0 {
			cleaned = append(cleaned, trimmedLine)
		}
	}

	result := strings.Join(cleaned, " ")

	// Replace multiple spaces with single space
	result = strings.Join(strings.Fields(result), " ")

	return result
}

// UpdateReleaseNotesWithWorkflowCheck creates or updates GitHub release notes after checking workflow status
func (m *ReleaseManager) UpdateReleaseNotesWithWorkflowCheck(tagName string, skipApproval bool, waitForWorkflows bool, maxWaitSeconds int) error {
	// Extract owner and repo from git remote
	owner, repo, err := ExtractRepoInfo("")
	if err != nil {
		return fmt.Errorf("failed to determine repository info: %w", err)
	}

	// Wait for workflows to complete if requested
	if waitForWorkflows {
		// Wait for GitHub workflows to finish
		err := m.client.WaitForWorkflowsToComplete(owner, repo, tagName, maxWaitSeconds)
		if err != nil {
			fmt.Printf("Warning: %s\n", err)
			fmt.Println("Proceeding anyway...")
		}
	}

	// Call the regular update method
	return m.UpdateReleaseNotes(tagName, skipApproval)
}
