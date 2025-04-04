// Package secure provides secure storage capabilities for sensitive data like API keys.
// It uses the system's native keyring/keychain when available, with fallbacks for different platforms.
package secure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	// ServiceName is the name used to identify this application in the keyring
	ServiceName = "noidea-git-tool"
	
	// FallbackDir is the directory used for fallback storage if keyring is unavailable
	FallbackDir = ".noidea/secure"
	
	// FallbackFile is the filename used for fallback storage
	FallbackFile = "keyring.enc"
)

// ErrKeyNotFound indicates that a key was not found in the secure storage
var ErrKeyNotFound = errors.New("key not found in secure storage")

// StoreAPIKey securely stores an API key for a given provider
func StoreAPIKey(provider, apiKey string) error {
	// Standardize the provider name for consistency
	provider = normalizeProviderName(provider)
	
	err := keyring.Set(ServiceName, provider, apiKey)
	if err != nil {
		// If keyring failed, try to use fallback storage
		return storeInFallbackStorage(provider, apiKey)
	}
	
	return nil
}

// GetAPIKey retrieves an API key for a given provider from secure storage
func GetAPIKey(provider string) (string, error) {
	// Standardize the provider name for consistency
	provider = normalizeProviderName(provider)
	
	// Try to get from keyring first
	apiKey, err := keyring.Get(ServiceName, provider)
	if err == nil && apiKey != "" {
		return apiKey, nil
	}
	
	// If keyring failed, try fallback storage
	return getFromFallbackStorage(provider)
}

// DeleteAPIKey removes an API key from secure storage
func DeleteAPIKey(provider string) error {
	// Standardize the provider name for consistency
	provider = normalizeProviderName(provider)
	
	// Try to delete from keyring
	err := keyring.Delete(ServiceName, provider)
	
	// Also delete from fallback if it exists (regardless of keyring result)
	fallbackErr := deleteFromFallbackStorage(provider)
	
	// If keyring succeeded or fallback succeeded, return nil
	if err == nil || fallbackErr == nil {
		return nil
	}
	
	// Both failed
	return fmt.Errorf("failed to delete key: %v", err)
}

// normalizeProviderName standardizes provider names
func normalizeProviderName(provider string) string {
	provider = strings.ToLower(provider)
	
	// Map common variations to standard names
	switch provider {
	case "openai", "open-ai", "gpt":
		return "openai"
	case "xai", "x-ai", "grok":
		return "xai" 
	case "deepseek", "deep-seek":
		return "deepseek"
	default:
		return provider
	}
}

// storeInFallbackStorage stores API keys in an encrypted file as fallback
func storeInFallbackStorage(provider, apiKey string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	secureDir := filepath.Join(homeDir, FallbackDir)
	if err := os.MkdirAll(secureDir, 0700); err != nil {
		return fmt.Errorf("failed to create secure directory: %w", err)
	}
	
	// In a real implementation, this would encrypt the data
	// For now, we'll create a simple obfuscation
	filePath := filepath.Join(secureDir, FallbackFile)
	
	// Read existing data first
	existingData := make(map[string]string)
	if fileData, err := os.ReadFile(filePath); err == nil {
		lines := strings.Split(string(fileData), "\n")
		for _, line := range lines {
			if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
				existingData[parts[0]] = parts[1]
			}
		}
	}

	// Update or add the new key
	existingData[provider] = obfuscate(apiKey)
	
	// Write all data back
	var sb strings.Builder
	for k, v := range existingData {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
		sb.WriteString("\n")
	}
	
	return os.WriteFile(filePath, []byte(sb.String()), 0600)
}

// getFromFallbackStorage retrieves API keys from fallback storage
func getFromFallbackStorage(provider string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	filePath := filepath.Join(homeDir, FallbackDir, FallbackFile)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", ErrKeyNotFound
	}
	
	lines := strings.Split(string(fileData), "\n")
	for _, line := range lines {
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			if parts[0] == provider {
				return deobfuscate(parts[1]), nil
			}
		}
	}
	
	return "", ErrKeyNotFound
}

// deleteFromFallbackStorage removes API keys from fallback storage
func deleteFromFallbackStorage(provider string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	filePath := filepath.Join(homeDir, FallbackDir, FallbackFile)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		// If file doesn't exist, consider it a success
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	
	// Read existing data
	existingData := make(map[string]string)
	lines := strings.Split(string(fileData), "\n")
	for _, line := range lines {
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			if parts[0] != provider {
				existingData[parts[0]] = parts[1]
			}
		}
	}
	
	// Write remaining data back
	var sb strings.Builder
	for k, v := range existingData {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
		sb.WriteString("\n")
	}
	
	return os.WriteFile(filePath, []byte(sb.String()), 0600)
}

// obfuscate provides a simple obfuscation for the fallback storage
// Note: This is NOT secure encryption, just simple obfuscation to prevent casual viewing
func obfuscate(text string) string {
	// Simple XOR with a fixed key - this is NOT secure encryption
	// In a real implementation, use proper encryption with a secure key
	key := []byte("noiDeA-SEcUrE-ObfUsCaTiOn-KeY")
	result := make([]byte, len(text))
	
	for i := 0; i < len(text); i++ {
		result[i] = text[i] ^ key[i%len(key)]
	}
	
	// Return as hex string for storage
	var sb strings.Builder
	for _, b := range result {
		sb.WriteString(fmt.Sprintf("%02x", b))
	}
	
	return sb.String()
}

// deobfuscate reverses the obfuscation
func deobfuscate(hexText string) string {
	// Convert hex to bytes
	if len(hexText) == 0 || len(hexText)%2 != 0 {
		return ""
	}
	
	result := make([]byte, len(hexText)/2)
	for i := 0; i < len(hexText); i += 2 {
		var b byte
		fmt.Sscanf(hexText[i:i+2], "%02x", &b)
		result[i/2] = b
	}
	
	// Apply XOR with the same key
	key := []byte("noiDeA-SEcUrE-ObfUsCaTiOn-KeY")
	for i := 0; i < len(result); i++ {
		result[i] = result[i] ^ key[i%len(key)]
	}
	
	return string(result)
}

// GetSecureStorageStatus returns information about the secure storage status
func GetSecureStorageStatus() map[string]string {
	status := make(map[string]string)
	
	// Check if keyring is available
	testKey := "noidea-test-key"
	testValue := "noidea-test-value"
	
	err := keyring.Set(ServiceName, testKey, testValue)
	if err == nil {
		// Successfully stored, now try to retrieve
		value, err := keyring.Get(ServiceName, testKey)
		if err == nil && value == testValue {
			status["keyring"] = "available"
			// Clean up test key
			keyring.Delete(ServiceName, testKey)
		} else {
			status["keyring"] = "retrieval-failed"
		}
	} else {
		status["keyring"] = "unavailable"
	}
	
	// Check fallback storage
	homeDir, err := os.UserHomeDir()
	if err == nil {
		secureDir := filepath.Join(homeDir, FallbackDir)
		if _, err := os.Stat(secureDir); err == nil {
			status["fallback"] = "directory-exists"
		} else if os.IsNotExist(err) {
			status["fallback"] = "directory-not-exists"
		} else {
			status["fallback"] = "directory-error"
		}
	} else {
		status["fallback"] = "homedir-error"
	}
	
	// Add platform information
	status["platform"] = runtime.GOOS
	
	return status
} 