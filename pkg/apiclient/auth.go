package apiclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// authenticate requests a JWT token from SitecoreAI API
// using client ID and client secret
func (c *Client) Authenticate() error {
	// Create request payload
	payload := url.Values{}
	payload.Set("audience", "https://api.sitecorecloud.io")
	payload.Set("grant_type", "client_credentials")
	payload.Set("client_id", c.ClientID)
	payload.Set("client_secret", c.ClientSecret)

	// Create HTTP request
	req, err := http.NewRequest(
		"POST",
		c.AuthURL,
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
