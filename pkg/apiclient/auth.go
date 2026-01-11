package apiclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AuthResponse represents the JWT authentication response
// from SitecoreAI API
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// CLIUserConfig represents the structure of .sitecore/user.json file
type CLIUserConfig struct {
	Endpoints struct {
		XMCloud struct {
			Host         string `json:"host"`
			Authority    string `json:"authority"`
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		} `json:"xmCloud"`
	} `json:"endpoints"`
}

// authenticate requests a JWT token from SitecoreAI API
// using client ID and client secret
func (c *Client) Authenticate() error {

	// If we already have a token, no need to authenticate
	if c.Token != "" {
		return nil
	}

	// If we have token in cli config, no need to authenticate
	if c.CliConfig != nil && len(c.CliConfig.Endpoints.XMCloud.AccessToken) > 0 {
		c.Token = c.CliConfig.Endpoints.XMCloud.AccessToken
		return nil
	}

	// Create request payload
	payload := url.Values{}
	payload.Set("audience", "https://api.sitecorecloud.io")
	payload.Set("grant_type", "client_credentials")
	payload.Set("client_id", c.ClientID)
	payload.Set("client_secret", c.ClientSecret)

	// Create HTTP request
	req, err := http.NewRequest(
		"POST",
		c.AuthURL+"/oauth/token",
		strings.NewReader(payload.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResponse AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Set token
	c.Token = authResponse.AccessToken

	return nil
}

// EnsureTokenValid checks if the current token is valid and
// refreshes it if needed
func (c *Client) EnsureTokenValid() error {
	if c.Token == "" {
		return c.Authenticate()
	}

	// Parse token to check expiration
	// This is a simplified check - in production you would properly parse the JWT
	parts := strings.Split(c.Token, ".")
	if len(parts) != 3 {
		return c.Authenticate()
	}

	// Decode payload
	payload, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return c.Authenticate()
	}

	var data map[string]interface{}
	err = json.Unmarshal(payload, &data)
	if err != nil {
		return c.Authenticate()
	}

	exp, ok := data["exp"].(float64)
	if !ok {
		return c.Authenticate()
	}

	// Check if token is expired or about to expire (within 5 minutes)
	if time.Now().Unix() > int64(exp) || time.Now().Unix() > int64(exp)-300 {
		return c.Authenticate()
	}

	return nil
}

func findCLIUserConfig() (*CLIUserConfig, error) {

	configPath, err := findCLIUserConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %v", err)
	}

	return readCLIUserConfig(configPath)
}

func readCLIUserConfig(configPath string) (*CLIUserConfig, error) {
	if configPath == "" {
		return nil, fmt.Errorf("no configPath specified")
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read user.json: %v", err)
	}

	var config CLIUserConfig
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user.json: %v", err)
	}

	return &config, nil
}

// findCLIUserConfig searches for .sitecore/user.json in current and parent directories
func findCLIUserConfigPath() (string, error) {
	// Start from current directory and move up the directory tree
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	for {
		configPath := filepath.Join(currentDir, ".sitecore", "user.json")

		// Check if file exists
		if _, err := os.Stat(configPath); err == nil {
			// File exists, try to read it
			_, err := os.ReadFile(configPath)
			if err != nil {
				return "", fmt.Errorf("failed to read user.json: %v", err)
			}

			return configPath, nil
		}

		// Move up to parent directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached root directory
			break
		}
		currentDir = parentDir
	}

	return "", nil
}
