package apiclient

import (
	"encoding/json"
	"fmt"
)

// ClientCreateResponse represents the response when creating a client
// This struct should be updated based on the actual API response structure
type ClientCreateResponse struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

type ClientType int

const (
	ClientTypeCM          ClientType = 1
	ClientTypeEdge        ClientType = 2
	ClientTypeDeploy      ClientType = 3
	ClientTypeEditingHost ClientType = 4
)

// ClientDto represents a client in the response
type ClientDto struct {
	ID              string     `json:"id,omitempty"`
	Name            string     `json:"name,omitempty"`
	Description     string     `json:"description,omitempty"`
	ClientID        string     `json:"clientId,omitempty"`
	CreatedAt       string     `json:"createdAt,omitempty"`
	ClientType      ClientType `json:"clientType,omitempty"`
	ProjectName     string     `json:"projectName,omitempty"`
	EnvironmentName string     `json:"environmentName,omitempty"`
}

// OrganizationClientDto represents an organization client in the response
type OrganizationClientDto struct {
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	ClientID    string     `json:"clientId,omitempty"`
	CreatedAt   string     `json:"createdAt,omitempty"`
	ClientType  ClientType `json:"clientType,omitempty"`
}

// ClientsListResponse represents the response when getting environment clients
type ClientsListResponse struct {
	Items []ClientDto `json:"items,omitempty"`
}

// OrganizationClientsListResponse represents the response when getting organization clients
type OrganizationClientsListResponse struct {
	Items []OrganizationClientDto `json:"items,omitempty"`
}

// CMClientCreateRequest represents the request for creating a CM client
type CMClientCreateRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	ProjectID     string `json:"projectId"`
	EnvironmentID string `json:"environmentId"`
}

// EdgeClientCreateRequest represents the request for creating an Edge client
type EdgeClientCreateRequest struct {
	ProjectID     string `json:"projectId"`
	EnvironmentID string `json:"environmentId"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
}

// DeployClientRequest represents the request for creating a Deploy client
type DeployClientRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// EditingHostBuildClientRequest represents the request for creating an Editing Host Build client
type EditingHostBuildClientRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	ProjectID     string `json:"projectId"`
	EnvironmentID string `json:"environmentId"`
}

// CreateCMClient creates a new CM automation client
func (c *Client) CreateCMClient(projectID string, environmentID string, name string, description string) (*ClientCreateResponse, error) {
	body := CMClientCreateRequest{
		Name:          name,
		Description:   description,
		ProjectID:     projectID,
		EnvironmentID: environmentID,
	}

	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   "/api/clients/v1/cm",
		Body:   body,
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create CM client: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response
	var clientResponse ClientCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&clientResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CM client response: %v", err)
	}

	return &clientResponse, nil
}

// CreateEdgeClient creates a new Edge automation client
func (c *Client) CreateEdgeClient(projectID string, environmentID string, name string, description string) (*ClientCreateResponse, error) {
	body := EdgeClientCreateRequest{
		ProjectID:     projectID,
		EnvironmentID: environmentID,
		Name:          name,
		Description:   description,
	}

	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   "/api/clients/v1/edge",
		Body:   body,
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create Edge client: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response
	var clientResponse ClientCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&clientResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Edge client response: %v", err)
	}

	return &clientResponse, nil
}

// CreateDeployClient creates a new Deploy automation client
func (c *Client) CreateDeployClient(name string, description string) (*ClientCreateResponse, error) {
	body := DeployClientRequest{
		Name:        name,
		Description: description,
	}

	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   "/api/clients/v1/deploy",
		Body:   body,
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create Deploy client: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response
	var clientResponse ClientCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&clientResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Deploy client response: %v", err)
	}

	return &clientResponse, nil
}

// CreateEditingHostBuildClient creates a new Editing Host Build automation client
func (c *Client) CreateEditingHostBuildClient(projectID string, environmentID string, name string, description string) (*ClientCreateResponse, error) {
	body := EditingHostBuildClientRequest{
		Name:          name,
		Description:   description,
		ProjectID:     projectID,
		EnvironmentID: environmentID,
	}

	// Create request options
	opts := RequestOptions{
		Method: "POST",
		Path:   "/api/clients/v1/ehbuild",
		Body:   body,
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create Editing Host Build client: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response
	var clientResponse ClientCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&clientResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Editing Host Build client response: %v", err)
	}

	return &clientResponse, nil
}

// DeleteClient deletes an automation client by ID
func (c *Client) DeleteClient(id string) error {
	// Create request options
	opts := RequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("/api/clients/v1/%s", id),
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return fmt.Errorf("failed to delete client: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

// GetClientsForOrganization retrieves a list of automation clients for the organization
func (c *Client) GetClientsForOrganization() (*OrganizationClientsListResponse, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   "/api/clients/v1/organization",
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization clients: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response
	var response OrganizationClientsListResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode organization clients response: %v", err)
	}

	return &response, nil
}

// GetClientsForEnvironment retrieves a list of automation clients for environments
func (c *Client) GetClientsForEnvironment() (*ClientsListResponse, error) {
	// Create request options
	opts := RequestOptions{
		Method: "GET",
		Path:   "/api/clients/v1/environment",
	}

	// Make the request
	resp, err := c.doRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment clients: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Parse the response
	var response ClientsListResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode environment clients response: %v", err)
	}

	return &response, nil
}
