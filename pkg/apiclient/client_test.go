package apiclient

import (
	"os"
	"testing"
)

func TestGetProjects(t *testing.T) {
	// Get client credentials from environment variables
	clientID := os.Getenv("SITECOREAI_CLIENT_ID")
	clientSecret := os.Getenv("SITECOREAI_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skip("SITECOREAI_CLIENT_ID and SITECOREAI_CLIENT_SECRET environment variables must be set to run this test")
	}

	// Create new client
	client, err := NewClientFromEnv()
	if err != nil {
		t.Errorf("Client instatiation failed: %v", err)
	}

	// Authenticate
	err = client.Authenticate()
	if err != nil {
		t.Fatalf("Authentication failed: %v", err)
	}

	// Test GetProjects method
	projects, err := client.GetProjects()
	if err != nil {
		t.Errorf("GetProjects failed: %v", err)
	}

	// Verify we got some projects
	if len(projects) == 0 {
		t.Error("No projects returned")
	}

	t.Logf("GetProjects test passed successfully. Found %d projects", len(projects))
}

func TestGetProjectAndEnvironments(t *testing.T) {
	// Get client credentials from environment variables
	clientID := os.Getenv("SITECOREAI_CLIENT_ID")
	clientSecret := os.Getenv("SITECOREAI_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skip("SITECOREAI_CLIENT_ID and SITECOREAI_CLIENT_SECRET environment variables must be set to run this test")
	}

	// Create new client
	client, err := NewClientFromEnv()
	if err != nil {
		t.Errorf("Client instatiation failed: %v", err)
	}

	// Authenticate
	err = client.Authenticate()
	if err != nil {
		t.Fatalf("Authentication failed: %v", err)
	}

	// Test GetProjects method to get a project
	projects, err := client.GetProjects()
	if err != nil {
		t.Fatalf("GetProjects failed: %v", err)
	}

	// Verify we got some projects
	if len(projects) == 0 {
		t.Skip("No projects available to test environments")
	}

	// Use the first project for testing
	project := projects[0]
	t.Logf("Testing with project: %s (ID: %s)", project.Name, project.ID)

	// Test GetProject method
	projectDetails, err := client.GetProject(project.ID)
	if err != nil {
		t.Errorf("GetProject failed: %v", err)
	}

	// Verify project details
	if projectDetails == nil {
		t.Error("Project details are nil")
	} else {
		t.Logf("Project details retrieved successfully: %s (ID: %s)", projectDetails.Name, projectDetails.ID)
	}

	// Test GetProjectEnvironments method
	environments, err := client.GetProjectEnvironments(project.ID)
	if err != nil {
		t.Errorf("GetProjectEnvironments failed: %v", err)
	}

	// Verify we got some environments
	if len(environments) == 0 {
		t.Log("No environments found for the project")
	} else {
		t.Logf("Found %d environments for project %s", len(environments), project.Name)
		for _, env := range environments {
			t.Logf("  - Environment: %s (ID: %s)", env.Name, env.ID)
		}
	}

	t.Log("GetProjectAndEnvironments test passed successfully")
}
