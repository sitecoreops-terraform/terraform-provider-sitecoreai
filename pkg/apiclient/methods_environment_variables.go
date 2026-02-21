package apiclient

import (
	"encoding/json"
	"fmt"
)

// EnvironmentVariable represents an environment variable from the API
type EnvironmentVariable struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
	Target string `json:"target,omitempty"`
}

// EnvironmentVariableUpsertRequestBodyDto represents the request body for setting environment variables
type EnvironmentVariableUpsertRequestBodyDto struct {
	Secret bool    `json:"secret"`
	Value  string  `json:"value"`
	Target *string `json:"target,omitempty"`
}

// GetEnvironmentVariables retrieves variables for an environment
func (c *Client) GetEnvironmentVariables(environmentID string) ([]EnvironmentVariable, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/environments/v1/%s/variables", environmentID),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variables: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response as an array of EnvironmentVariable
	var variables []EnvironmentVariable
	err = json.NewDecoder(resp.Body).Decode(&variables)
	if err != nil {
		return nil, fmt.Errorf("failed to decode environment variables: %v", err)
	}

	return variables, nil
}

// SetEnvironmentVariable sets a variable for an environment
func (c *Client) SetEnvironmentVariable(environmentID string, variableName string, requestBody EnvironmentVariableUpsertRequestBodyDto) error {
	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/api/environments/v1/%s/variables/%s", environmentID, variableName),
		Body:   requestBody,
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return fmt.Errorf("failed to set environment variable: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// DeleteEnvironmentVariable deletes a variable from an environment
func (c *Client) DeleteEnvironmentVariable(environmentID string, variableName string) error {
	// Create request options
	opts := RequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("/api/environments/v1/%s/variables/%s", environmentID, variableName),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return fmt.Errorf("failed to delete environment variable: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}
