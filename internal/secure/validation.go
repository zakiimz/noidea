// Package secure provides secure storage and validation for API keys
package secure

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// ValidateAPIKey checks if the API key works with the provider
func ValidateAPIKey(provider, apiKey string) (bool, error) {
	// For all providers, try to validate against their API
	var baseURL string
	
	switch provider {
	case "xai":
		// Use the correct xAI endpoint from docs.x.ai
		baseURL = "https://api.x.ai/v1/models"
	case "openai":
		baseURL = "https://api.openai.com/v1/models"
	case "deepseek":
		baseURL = "https://api.deepseek.com/v1/models"
	default:
		// Default to OpenAI for unknown providers
		baseURL = "https://api.openai.com/v1/models"
	}
	
	// Debug info to help troubleshoot
	fmt.Printf("Debug: Validating %s API key with API endpoint: %s\n", provider, baseURL)
	
	return validateAPIKeyWithEndpoint(apiKey, baseURL)
}

// validateAPIKeyWithEndpoint checks if an API key is valid for any API endpoint
func validateAPIKeyWithEndpoint(apiKey, baseURL string) (bool, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return false, err
	}
	
	req.Header.Add("Authorization", "Bearer "+apiKey)
	
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()
	
	// Debug response
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Debug: API response status: %d\n", resp.StatusCode)
	if len(body) < 1000 {
		fmt.Printf("Debug: API response body: %s\n", string(body))
	} else {
		fmt.Printf("Debug: API response body length: %d bytes\n", len(body))
	}
	
	// For our purposes, consider any response (even error) as valid
	// Since many providers will return errors for invalid models, etc.
	// but a 401/403 specifically indicates an auth problem
	return resp.StatusCode != 401 && resp.StatusCode != 403, nil
} 