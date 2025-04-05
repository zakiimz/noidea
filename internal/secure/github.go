// Package secure provides secure storage capabilities for sensitive data like API keys.
package secure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// GitHubTokenKey is the key used to store the GitHub token in the secure storage
	GitHubTokenKey = "github-token"

	// GitHubAPIURL is the base URL for GitHub API
	GitHubAPIURL = "https://api.github.com"
)

// StoreGitHubToken securely stores a GitHub Personal Access Token
func StoreGitHubToken(token string) error {
	return StoreAPIKey(GitHubTokenKey, token)
}

// GetGitHubToken retrieves the GitHub Personal Access Token from secure storage
func GetGitHubToken() (string, error) {
	return GetAPIKey(GitHubTokenKey)
}

// DeleteGitHubToken removes the GitHub Personal Access Token from secure storage
func DeleteGitHubToken() error {
	return DeleteAPIKey(GitHubTokenKey)
}

// ValidateGitHubToken checks if the GitHub token is valid by making a request to the GitHub API
func ValidateGitHubToken(token string) (bool, map[string]interface{}, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", GitHubAPIURL+"/user", nil)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return false, nil, fmt.Errorf("invalid token or API error, status code: %d", resp.StatusCode)
	}

	// Parse user information
	var userData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return true, nil, fmt.Errorf("failed to parse user data: %w", err)
	}

	return true, userData, nil
}
