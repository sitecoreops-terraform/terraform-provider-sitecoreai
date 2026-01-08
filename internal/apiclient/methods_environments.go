package apiclient

import (
	"encoding/json"
	"fmt"
)

// Environment represents a Sitecore environment
type Environment struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	ProjectID   string            `json:"project_id"`
	Description string            `json:"description,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
	// Add other environment fields as needed based on API specification
}

// CreateEnvironment creates a new environment for a project
func (c *Client) CreateEnvironment(projectID string, environment Environment) (*Environment, error) {
	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/api/projects/v1/%s/environments", projectID),
		Body:   environment,
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment: %v", err)
	}

	defer resp.Body.Close()

	// Parse the response
	var createdEnvironment Environment
	err = json.NewDecoder(resp.Body).Decode(&createdEnvironment)
	if err != nil {
		return nil, fmt.Errorf("failed to decode created environment: %v", err)
	}

	return &createdEnvironment, nil
}

// GetEnvironment retrieves details of an existing environment
func (c *Client) GetEnvironment(projectID string, environmentID string) (*Environment, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/projects/v1/%s/environments/%s", projectID, environmentID),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %v", err)
	}

	defer resp.Body.Close()

	// Parse the response
	var environment Environment
	err = json.NewDecoder(resp.Body).Decode(&environment)
	if err != nil {
		return nil, fmt.Errorf("failed to decode environment: %v", err)
	}

	return &environment, nil
}

// UpdateEnvironment updates an existing environment
func (c *Client) UpdateEnvironment(projectID string, environmentID string, environment Environment) error {
	// Create request options
	opts := RequestOptions{
		Method: "PUT",
		Path:   fmt.Sprintf("/api/projects/v1/%s/environments/%s", projectID, environmentID),
		Body:   environment,
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return fmt.Errorf("failed to update environment: %v", err)
	}

	defer resp.Body.Close()

	return nil
}

// DeleteEnvironment deletes an existing environment
func (c *Client) DeleteEnvironment(projectID string, environmentID string) error {
	// Create request options
	opts := RequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("/api/projects/v1/%s/environments/%s", projectID, environmentID),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return fmt.Errorf("failed to delete environment: %v", err)
	}

	defer resp.Body.Close()

	return nil
}

// GetProjectEnvironments lists environments for a specific project
func (c *Client) GetProjectEnvironments(projectID string) ([]Environment, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/projects/v1/%s/environments", projectID),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get project environments: %v", err)
	}

	defer resp.Body.Close()

	// Parse the response
	var environments []Environment
	err = json.NewDecoder(resp.Body).Decode(&environments)
	if err != nil {
		return nil, fmt.Errorf("failed to decode project environments: %v", err)
	}

	return environments, nil
}
