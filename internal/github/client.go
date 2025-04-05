// Package github provides functionality for interacting with the GitHub API
package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/AccursedGalaxy/noidea/internal/secure"
)

// Client represents a GitHub API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// NewClient creates a new GitHub API client
func NewClient() (*Client, error) {
	token, err := secure.GetGitHubToken()
	if err != nil {
		return nil, fmt.Errorf("GitHub authentication required. Run 'noidea github auth' to authenticate: %w", err)
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: secure.GitHubAPIURL,
		token:   token,
	}, nil
}

// GetUser retrieves the authenticated user's information
func (c *Client) GetUser() (map[string]interface{}, error) {
	return c.get("/user")
}

// GetRepository retrieves a repository by owner and repo name
func (c *Client) GetRepository(owner, repo string) (map[string]interface{}, error) {
	return c.get(fmt.Sprintf("/repos/%s/%s", owner, repo))
}

// CreateRelease creates a new release in the specified repository
func (c *Client) CreateRelease(owner, repo, tagName, name, body string, draft, prerelease bool) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"tag_name":   tagName,
		"name":       name,
		"body":       body,
		"draft":      draft,
		"prerelease": prerelease,
	}

	return c.post(fmt.Sprintf("/repos/%s/%s/releases", owner, repo), payload)
}

// get performs a GET request to the GitHub API
func (c *Client) get(path string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// post performs a POST request to the GitHub API
func (c *Client) post(path string, payload interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// patch performs a PATCH request to the GitHub API
func (c *Client) patch(path string, payload interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", c.baseURL+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// doRequest executes the HTTP request and processes the response
func (c *Client) doRequest(req *http.Request) (map[string]interface{}, error) {
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GitHub API error: %s (status code: %d)", string(body), resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// IsAuthenticated checks if the client has a valid GitHub token
func (c *Client) IsAuthenticated() (bool, error) {
	_, err := c.GetUser()
	return err == nil, err
}

// ExtractRepoInfo extracts owner and repo name from a Git remote URL or the current repository
func ExtractRepoInfo(remoteURL string) (string, string, error) {
	// If no remote URL provided, try to get it from the current git repository
	if remoteURL == "" {
		var err error
		remoteURL, err = getOriginRemoteURL()
		if err != nil {
			return "", "", err
		}
	}

	// Handle SSH URLs (git@github.com:owner/repo.git)
	sshPattern := regexp.MustCompile(`git@github\.com:([^/]+)/([^.]+)\.git`)
	if matches := sshPattern.FindStringSubmatch(remoteURL); len(matches) == 3 {
		return matches[1], matches[2], nil
	}

	// Handle HTTPS URLs (https://github.com/owner/repo.git)
	httpsPattern := regexp.MustCompile(`https://github\.com/([^/]+)/([^./]+)(?:\.git)?`)
	if matches := httpsPattern.FindStringSubmatch(remoteURL); len(matches) == 3 {
		return matches[1], matches[2], nil
	}

	return "", "", fmt.Errorf("could not parse GitHub repository URL: %s", remoteURL)
}

// getOriginRemoteURL gets the origin remote URL from the current git repository
func getOriginRemoteURL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git remote: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetWorkflowRunsForRef gets the workflow runs triggered by a specific git ref (tag/branch)
func (c *Client) GetWorkflowRunsForRef(owner, repo, ref string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("/repos/%s/%s/actions/runs?event=push&branch=%s", owner, repo, ref)
	response, err := c.get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow runs: %w", err)
	}

	// Extract workflow runs from response
	workflowRunsObj, ok := response["workflow_runs"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format: workflow_runs not found")
	}

	workflowRuns, ok := workflowRunsObj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format: workflow_runs is not an array")
	}

	// Convert to maps
	result := make([]map[string]interface{}, 0, len(workflowRuns))
	for _, run := range workflowRuns {
		if runMap, ok := run.(map[string]interface{}); ok {
			result = append(result, runMap)
		}
	}

	return result, nil
}

// AreAllWorkflowsComplete checks if all workflows for a ref have completed (success or failure)
func (c *Client) AreAllWorkflowsComplete(owner, repo, ref string) (bool, error) {
	runs, err := c.GetWorkflowRunsForRef(owner, repo, ref)
	if err != nil {
		return false, err
	}

	// If no runs found, assume workflows are complete
	if len(runs) == 0 {
		return true, nil
	}

	// Check if all runs have a conclusive status
	for _, run := range runs {
		status, ok := run["status"].(string)
		if !ok {
			continue
		}

		// If any workflow is still in progress, return false
		if status == "queued" || status == "in_progress" || status == "waiting" {
			return false, nil
		}
	}

	// All workflows have completed
	return true, nil
}

// WaitForWorkflowsToComplete waits for all workflows to complete with a max wait time
func (c *Client) WaitForWorkflowsToComplete(owner, repo, ref string, maxWaitSeconds int) error {
	fmt.Printf("Checking GitHub workflow status for %s...\n", ref)

	// Start a timeout context
	timeoutChan := time.After(time.Duration(maxWaitSeconds) * time.Second)
	ticker := time.NewTicker(2 * time.Second) // Check every 2 seconds
	defer ticker.Stop()

	// Spinner animation chars
	spinChars := []string{"⋮", "⋰", "⋮", "⋱"}
	spinIdx := 0
	count := 0

	for {
		select {
		case <-timeoutChan:
			// Clear the current line before error message
			fmt.Print("\r\033[K")
			return fmt.Errorf("timed out waiting for workflows to complete after %d seconds", maxWaitSeconds)
		case <-ticker.C:
			complete, err := c.AreAllWorkflowsComplete(owner, repo, ref)
			if err != nil {
				// Clear the current line
				fmt.Print("\r\033[K")
				fmt.Printf("Warning: Failed to check workflow status: %s\n", err)
				// Continue checking despite errors
			} else if complete {
				// Clear the current line
				fmt.Print("\r\033[K")
				fmt.Println("✅ All GitHub workflows completed successfully!")
				return nil
			} else {
				// Increment the spinner index
				spinIdx = (spinIdx + 1) % len(spinChars)
				count++

				// Calculate elapsed time
				elapsedSecs := count * 2
				timeStr := ""
				if elapsedSecs >= 60 {
					timeStr = fmt.Sprintf(" (%dm%02ds)", elapsedSecs/60, elapsedSecs%60)
				} else {
					timeStr = fmt.Sprintf(" (%ds)", elapsedSecs)
				}

				// Clear line and show spinning animation with elapsed time
				fmt.Printf("\r\033[K⏳ Workflows still running... %s %s", spinChars[spinIdx], timeStr)
			}
		}
	}
}
