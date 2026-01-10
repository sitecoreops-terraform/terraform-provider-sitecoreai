package apiclient

import (
	"fmt"
	"io"
	"strings"
)

// ObtainEditingSecret calls the obtain-editing-secret endpoint for an environment
func (c *Client) ObtainEditingSecret(environmentID string) (string, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/environments/v1/%s/obtain-editing-secret", environmentID),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return "", fmt.Errorf("failed to obtain editing secret: %v", err)
	}

	defer resp.Body.Close()

	// If JSON decoding fails, try to read as plain text (new API format)
	// Read the response body as plain text
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read editing secret response: %v", err)
	}

	// Trim whitespace and return as ApiKey
	secret := strings.TrimSpace(string(bodyBytes))

	// Return response with just the API key (EdgeUrl will be empty)
	return secret, nil
}
