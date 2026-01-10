package apiclient

import (
	"encoding/json"
	"fmt"
	"time"
)

// Environment represents a Sitecore environment
type Environment struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	ProjectID               string `json:"projectId"`
	Host                    string `json:"host,omitempty"`
	PlatformTenantId        string `json:"platformTenantId,omitempty"`
	PlatformTenantName      string `json:"platformTenantName,omitempty"`
	TenantType              string `json:"tenantType,omitempty"`
	CreatedAt               string `json:"createdAt,omitempty"`
	CreatedBy               string `json:"createdBy,omitempty"`
	LastUpdatedBy           string `json:"lastUpdatedBy,omitempty"`
	LastUpdatedAt           string `json:"lastUpdatedAt,omitempty"`
	IsDeleted               bool   `json:"isDeleted,omitempty"`
	PreviewContextId        string `json:"previewContextId,omitempty"`
	LiveContextId           string `json:"liveContextId,omitempty"`
	HighAvailabilityEnabled bool   `json:"highAvailabilityEnabled,omitempty"`
}

type CreateEnvironmentRequest struct {
	Name       string `json:"name"`
	TenantType int    `json:"tenantType,omitempty"`
	Type       string `json:"type"`
}

type EnvironmentType int

const (
	EnvironmentTypeCombined EnvironmentType = 0
	EnvironmentTypeCmOnly   EnvironmentType = 1
)

// CreateEnvironment creates a new environment for a project
func (c *Client) CreateEnvironment(projectID string, name string, isProd bool, environmentType EnvironmentType) (*Environment, error) {

	tenantType := 0
	if isProd {
		tenantType = 1
	}

	body := CreateEnvironmentRequest{
		Name:       name,
		TenantType: tenantType,
	}

	if environmentType == EnvironmentTypeCmOnly {
		body.Type = "cm"
	}

	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/api/projects/v2/%s/environments", projectID),
		Body:   body,
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

// WaitForEnvironmentReady polls the environment until it's ready and has context IDs
func (c *Client) WaitForEnvironmentReady(projectID string, environmentID string, timeoutMinutes int) (*Environment, error) {
	// Set timeout
	timeout := time.Duration(timeoutMinutes) * time.Minute
	startTime := time.Now()

	// Polling interval
	pollInterval := 1 * time.Second

	for {
		// Check if we've exceeded the timeout
		if time.Since(startTime) > timeout {
			return nil, fmt.Errorf("timed out waiting for environment to be ready after %d minutes", timeoutMinutes)
		}

		// Get the current environment status
		environment, err := c.GetEnvironment(projectID, environmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get environment status: %v", err)
		}

		// Check if environment has the required context IDs
		if environment.PreviewContextId != "" && environment.LiveContextId != "" {
			return environment, nil
		}

		// Wait before polling again
		time.Sleep(pollInterval)
	}
}

// GetEnvironment retrieves details of an existing environment
func (c *Client) GetEnvironment(projectID string, environmentID string) (*Environment, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/environments/v2/%s", environmentID),
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
		Path:   fmt.Sprintf("/api/environments/v1/%s", environmentID),
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
		Path:   fmt.Sprintf("/api/environments/v1/%s", environmentID),
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
