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
