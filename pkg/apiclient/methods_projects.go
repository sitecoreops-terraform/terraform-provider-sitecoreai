package apiclient

import (
	"encoding/json"
	"fmt"
)

// Project represents a Sitecore project
// This struct should be updated based on the actual API response structure
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Add other project fields as needed based on API specification
}

// GetProjects retrieves all projects
func (c *Client) GetProjects() ([]Project, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   "/api/projects/v1",
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response
	var projects []Project
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		return nil, fmt.Errorf("failed to decode projects: %v", err)
	}

	return projects, nil
}

// GetProject retrieves a specific project by ID
func (c *Client) GetProject(projectID string) (*Project, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/projects/v1/%s", projectID),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response
	var project Project
	err = json.NewDecoder(resp.Body).Decode(&project)
	if err != nil {
		return nil, fmt.Errorf("failed to decode project: %v", err)
	}

	return &project, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(project Project) (*Project, error) {
	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   "/api/projects/v1",
		Body:   project,
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response
	var createdProject Project
	err = json.NewDecoder(resp.Body).Decode(&createdProject)
	if err != nil {
		return nil, fmt.Errorf("failed to decode created project: %v", err)
	}

	return &createdProject, nil
}

// UpdateProject updates an existing project
func (c *Client) UpdateProject(projectID string, project Project) error {
	// Create request options
	opts := RequestOptions{
		Method: "PUT",
		Path:   fmt.Sprintf("/api/projects/v1/%s", projectID),
		Body:   project,
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return fmt.Errorf("failed to update project: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectID string) error {
	// Create request options
	opts := RequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("/api/projects/v1/%s", projectID),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return fmt.Errorf("failed to delete project: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}
