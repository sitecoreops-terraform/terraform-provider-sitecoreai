package apiclient

import (
	"encoding/json"
	"fmt"
)

// EnvironmentVariable represents an environment variable from the API
type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetEnvironmentVariables retrieves variables for an environment
func (c *Client) GetEnvironmentVariables(environmentID string) (map[string]string, error) {
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

	// Convert the array to a map
	variablesMap := make(map[string]string)
	for _, variable := range variables {
		variablesMap[variable.Name] = variable.Value
	}

	return variablesMap, nil
}

// SetEnvironmentVariable sets a variable for an environment
func (c *Client) SetEnvironmentVariable(environmentID string, variableName string, variableValue string) error {
	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/api/environments/v1/%s/variables/%s", environmentID, variableName),
		Body:   map[string]string{"value": variableValue},
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
