package apiclient

import (
	"encoding/json"
	"fmt"
)

// GetEnvironmentVariables retrieves variables for an environment
func (c *Client) GetEnvironmentVariables(projectID string, environmentID string) (map[string]string, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/projects/v1/%s/environments/%s/variables", projectID, environmentID),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variables: %v", err)
	}

	defer resp.Body.Close()

	// Parse the response
	var variables map[string]string
	err = json.NewDecoder(resp.Body).Decode(&variables)
	if err != nil {
		return nil, fmt.Errorf("failed to decode environment variables: %v", err)
	}

	return variables, nil
}

// SetEnvironmentVariable sets a variable for an environment
func (c *Client) SetEnvironmentVariable(projectID string, environmentID string, variableName string, variableValue string) error {
	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/api/projects/v1/%s/environments/%s/variables/%s", projectID, environmentID, variableName),
		Body:   map[string]string{"value": variableValue},
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return fmt.Errorf("failed to set environment variable: %v", err)
	}

	defer resp.Body.Close()

	return nil
}

// DeleteEnvironmentVariable deletes a variable from an environment
func (c *Client) DeleteEnvironmentVariable(projectID string, environmentID string, variableName string) error {
	// Create request options
	opts := RequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("/api/projects/v1/%s/environments/%s/variables/%s", projectID, environmentID, variableName),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return fmt.Errorf("failed to delete environment variable: %v", err)
	}

	defer resp.Body.Close()

	return nil
}
