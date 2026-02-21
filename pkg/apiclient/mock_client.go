package apiclient

import (
	"fmt"
)

// ClientInterface defines the interface for the Sitecore API client
type ClientInterface interface {
	Authenticate() error
	GetProjects() ([]Project, error)
	GetProject(id string) (*Project, error)
	CreateProject(project Project) (*Project, error)
	UpdateProject(id string, project Project) error
	DeleteProject(id string) error
	GetProjectEnvironments(projectID string) ([]Environment, error)
	GetEnvironment(projectID string, environmentID string) (*Environment, error)
	CreateEnvironment(projectID string, name string, isProd bool, tenantType EnvironmentType, cmEnvironmentId string) (*Environment, error)
	DeleteEnvironment(projectID string, environmentID string) error
	CreateCMClient(projectID string, environmentID string, name string, description string) (*ClientCreateResponse, error)
	CreateEdgeClient(projectID string, environmentID string, name string, description string) (*ClientCreateResponse, error)
	CreateDeployClient(name string, description string) (*ClientCreateResponse, error)
	CreateEditingHostBuildClient(projectID string, environmentID string, name string, description string) (*ClientCreateResponse, error)
	DeleteClient(clientID string) error
	ObtainEditingSecret(environmentID string) (string, error)
	GetEnvironmentVariables(projectID string, environmentID string) (map[string]string, error)
	SetEnvironmentVariable(projectID string, environmentID string, variableName string, requestBody EnvironmentVariableUpsertRequestBodyDto) error
	DeleteEnvironmentVariable(projectID string, environmentID string, variableName string) error
}

// MockClient is a mock implementation of the API client for testing
type MockClient struct {
	Projects             []Project
	Environments         []Environment
	Clients              []ClientCreateResponse
	ShouldFail           bool
	FailOnMethod         string
	CreatedProjects      []Project
	UpdatedProjects      []Project
	DeletedProjects      []string
	CreatedEnvironments  []Environment
	DeletedEnvironments  []string
	CreatedClients       []ClientCreateResponse
	DeletedClients       []string
	EditingSecrets       map[string]string            // projectID+environmentID -> secret
	EnvironmentVariables map[string]map[string]string // projectID+environmentID -> variables
}

// NewMockClient creates a new mock client with test data
func NewMockClient() *MockClient {
	return &MockClient{
		Projects: []Project{
			{
				ID:          "project-1",
				Name:        "Test Project 1",
				Description: "Test Description 1",
			},
			{
				ID:          "project-2",
				Name:        "Test Project 2",
				Description: "Test Description 2",
			},
		},
		Environments: []Environment{
			{
				ID:                      "env-1",
				Name:                    "Development",
				ProjectID:               "project-1",
				ProjectName:             "Test Project 1",
				OrganizationID:          "org-1",
				OrganizationName:        "Test Organization",
				Zone:                    "westus",
				Host:                    "dev.example.com",
				SitecoreMajorVersion:    10,
				SitecoreMinorVersion:    1,
				PlatformTenantId:        "tenant-1",
				PlatformTenantName:      "Dev Tenant",
				RepositoryBranch:        "main",
				TenantType:              "cm",
				ProvisioningStatus:      2,
				DeployOnCommit:          true,
				IsDeleted:               false,
				PreviewContextId:        "preview-1",
				LiveContextId:           "live-1",
				HighAvailabilityEnabled: false,
				Type:                    "cm",
			},
		},
		Clients: []ClientCreateResponse{
			{
				Name:         "test-client",
				Description:  "Test Client",
				ClientID:     "client-1",
				ClientSecret: "secret-1",
			},
		},
		EditingSecrets: map[string]string{
			"project-1-env-1": "test-secret-value",
		},
		EnvironmentVariables: map[string]map[string]string{
			"project-1-env-1": {
				"TEST_VAR": "test-value",
			},
		},
	}
}

func (m *MockClient) Authenticate() error {
	if m.ShouldFail && m.FailOnMethod == "Authenticate" {
		return fmt.Errorf("mock authentication failed")
	}
	return nil
}

func (m *MockClient) GetProjects() ([]Project, error) {
	if m.ShouldFail && m.FailOnMethod == "GetProjects" {
		return nil, fmt.Errorf("mock get projects failed")
	}
	return m.Projects, nil
}

func (m *MockClient) GetProject(id string) (*Project, error) {
	if m.ShouldFail && m.FailOnMethod == "GetProject" {
		return nil, fmt.Errorf("mock get project failed")
	}
	for _, project := range m.Projects {
		if project.ID == id {
			return &project, nil
		}
	}
	return nil, fmt.Errorf("project not found")
}

func (m *MockClient) CreateProject(project Project) (*Project, error) {
	if m.ShouldFail && m.FailOnMethod == "CreateProject" {
		return nil, fmt.Errorf("mock create project failed")
	}
	project.ID = fmt.Sprintf("project-%d", len(m.CreatedProjects)+1)
	m.CreatedProjects = append(m.CreatedProjects, project)
	m.Projects = append(m.Projects, project)
	return &project, nil
}

func (m *MockClient) UpdateProject(id string, project Project) error {
	if m.ShouldFail && m.FailOnMethod == "UpdateProject" {
		return fmt.Errorf("mock update project failed")
	}
	for i, p := range m.Projects {
		if p.ID == id {
			m.Projects[i] = project
			m.UpdatedProjects = append(m.UpdatedProjects, project)
			return nil
		}
	}
	return fmt.Errorf("project not found")
}

func (m *MockClient) DeleteProject(id string) error {
	if m.ShouldFail && m.FailOnMethod == "DeleteProject" {
		return fmt.Errorf("mock delete project failed")
	}
	for i, project := range m.Projects {
		if project.ID == id {
			m.Projects = append(m.Projects[:i], m.Projects[i+1:]...)
			m.DeletedProjects = append(m.DeletedProjects, id)
			return nil
		}
	}
	return fmt.Errorf("project not found")
}

func (m *MockClient) GetProjectEnvironments(projectID string) ([]Environment, error) {
	if m.ShouldFail && m.FailOnMethod == "GetProjectEnvironments" {
		return nil, fmt.Errorf("mock get project environments failed")
	}
	var environments []Environment
	for _, env := range m.Environments {
		if env.ProjectID == projectID {
			environments = append(environments, env)
		}
	}
	return environments, nil
}

func (m *MockClient) GetEnvironment(projectID string, environmentID string) (*Environment, error) {
	if m.ShouldFail && m.FailOnMethod == "GetEnvironment" {
		return nil, fmt.Errorf("mock get environment failed")
	}
	for _, env := range m.Environments {
		if env.ProjectID == projectID && env.ID == environmentID {
			return &env, nil
		}
	}
	return nil, fmt.Errorf("environment not found")
}

func (m *MockClient) CreateEnvironment(projectID string, name string, isProd bool, tenantType EnvironmentType, cmEnvironmentId string) (*Environment, error) {
	if m.ShouldFail && m.FailOnMethod == "CreateEnvironment" {
		return nil, fmt.Errorf("mock create environment failed")
	}
	env := Environment{
		ID:         fmt.Sprintf("env-%d", len(m.CreatedEnvironments)+1),
		Name:       name,
		ProjectID:  projectID,
		Host:       fmt.Sprintf("%s.example.com", name),
		TenantType: fmt.Sprintf("%d", tenantType),
		IsDeleted:  false,
	}
	m.CreatedEnvironments = append(m.CreatedEnvironments, env)
	m.Environments = append(m.Environments, env)
	return &env, nil
}

func (m *MockClient) DeleteEnvironment(projectID string, environmentID string) error {
	if m.ShouldFail && m.FailOnMethod == "DeleteEnvironment" {
		return fmt.Errorf("mock delete environment failed")
	}
	for i, env := range m.Environments {
		if env.ProjectID == projectID && env.ID == environmentID {
			m.Environments = append(m.Environments[:i], m.Environments[i+1:]...)
			m.DeletedEnvironments = append(m.DeletedEnvironments, environmentID)
			return nil
		}
	}
	return fmt.Errorf("environment not found")
}

func (m *MockClient) CreateCMClient(projectID string, environmentID string, name string, description string) (*ClientCreateResponse, error) {
	if m.ShouldFail && m.FailOnMethod == "CreateCMClient" {
		return nil, fmt.Errorf("mock create CM client failed")
	}
	client := ClientCreateResponse{
		Name:         name,
		Description:  description,
		ClientID:     fmt.Sprintf("cm-client-%d", len(m.CreatedClients)+1),
		ClientSecret: fmt.Sprintf("cm-secret-%d", len(m.CreatedClients)+1),
	}
	m.CreatedClients = append(m.CreatedClients, client)
	m.Clients = append(m.Clients, client)
	return &client, nil
}

func (m *MockClient) CreateEdgeClient(projectID string, environmentID string, name string, description string) (*ClientCreateResponse, error) {
	if m.ShouldFail && m.FailOnMethod == "CreateEdgeClient" {
		return nil, fmt.Errorf("mock create edge client failed")
	}
	client := ClientCreateResponse{
		Name:         name,
		Description:  description,
		ClientID:     fmt.Sprintf("edge-client-%d", len(m.CreatedClients)+1),
		ClientSecret: fmt.Sprintf("edge-secret-%d", len(m.CreatedClients)+1),
	}
	m.CreatedClients = append(m.CreatedClients, client)
	m.Clients = append(m.Clients, client)
	return &client, nil
}

func (m *MockClient) CreateDeployClient(name string, description string) (*ClientCreateResponse, error) {
	if m.ShouldFail && m.FailOnMethod == "CreateDeployClient" {
		return nil, fmt.Errorf("mock create deploy client failed")
	}
	client := ClientCreateResponse{
		Name:         name,
		Description:  description,
		ClientID:     fmt.Sprintf("deploy-client-%d", len(m.CreatedClients)+1),
		ClientSecret: fmt.Sprintf("deploy-secret-%d", len(m.CreatedClients)+1),
	}
	m.CreatedClients = append(m.CreatedClients, client)
	m.Clients = append(m.Clients, client)
	return &client, nil
}

func (m *MockClient) CreateEditingHostBuildClient(projectID string, environmentID string, name string, description string) (*ClientCreateResponse, error) {
	if m.ShouldFail && m.FailOnMethod == "CreateEditingHostBuildClient" {
		return nil, fmt.Errorf("mock create editing host build client failed")
	}
	client := ClientCreateResponse{
		Name:         name,
		Description:  description,
		ClientID:     fmt.Sprintf("ehbuild-client-%d", len(m.CreatedClients)+1),
		ClientSecret: fmt.Sprintf("ehbuild-secret-%d", len(m.CreatedClients)+1),
	}
	m.CreatedClients = append(m.CreatedClients, client)
	m.Clients = append(m.Clients, client)
	return &client, nil
}

func (m *MockClient) DeleteClient(clientID string) error {
	if m.ShouldFail && m.FailOnMethod == "DeleteClient" {
		return fmt.Errorf("mock delete client failed")
	}
	for i, client := range m.Clients {
		if client.ClientID == clientID {
			m.Clients = append(m.Clients[:i], m.Clients[i+1:]...)
			m.DeletedClients = append(m.DeletedClients, clientID)
			return nil
		}
	}
	return fmt.Errorf("client not found")
}

func (m *MockClient) ObtainEditingSecret(environmentID string) (string, error) {
	if m.ShouldFail && m.FailOnMethod == "ObtainEditingSecret" {
		return "", fmt.Errorf("mock get editing secret failed")
	}
	// Simple mock - return a secret based on environment ID
	if environmentID == "env-1" {
		return "test-secret-value", nil
	}
	return "", nil
}

func (m *MockClient) GetEnvironmentVariables(projectID string, environmentID string) (map[string]string, error) {
	if m.ShouldFail && m.FailOnMethod == "GetEnvironmentVariables" {
		return nil, fmt.Errorf("mock get environment variables failed")
	}
	key := projectID + "-" + environmentID
	if variables, ok := m.EnvironmentVariables[key]; ok {
		return variables, nil
	}
	return map[string]string{}, nil
}

func (m *MockClient) SetEnvironmentVariable(projectID string, environmentID string, variableName string, requestBody EnvironmentVariableUpsertRequestBodyDto) error {
	if m.ShouldFail && m.FailOnMethod == "SetEnvironmentVariable" {
		return fmt.Errorf("mock set environment variable failed")
	}
	key := projectID + "-" + environmentID
	if _, ok := m.EnvironmentVariables[key]; !ok {
		m.EnvironmentVariables[key] = map[string]string{}
	}
	m.EnvironmentVariables[key][variableName] = requestBody.Value
	return nil
}

func (m *MockClient) DeleteEnvironmentVariable(projectID string, environmentID string, variableName string) error {
	if m.ShouldFail && m.FailOnMethod == "DeleteEnvironmentVariable" {
		return fmt.Errorf("mock delete environment variable failed")
	}
	key := projectID + "-" + environmentID
	if variables, ok := m.EnvironmentVariables[key]; ok {
		delete(variables, variableName)
	}
	return nil
}
