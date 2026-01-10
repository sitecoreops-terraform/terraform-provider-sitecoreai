package apiclient

import (
	"os"
	"testing"
)

func TestObtainEditingSecret(t *testing.T) {
	// Get client credentials from environment variables
	clientID := os.Getenv("SITECORE_CLIENT_ID")
	clientSecret := os.Getenv("SITECORE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skip("SITECORE_CLIENT_ID and SITECORE_CLIENT_SECRET environment variables must be set to run this test")
	}

	// Create new client
	client := NewClientFromEnv()

	// Authenticate
	err := client.Authenticate()
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

	// Use the first project for testing
	project := projects[0]
	t.Logf("Testing with project: %s (ID: %s)", project.Name, project.ID)

	// Get environments for the project
	environments, err := client.GetProjectEnvironments(project.ID)
	if err != nil {
		t.Errorf("GetProjectEnvironments failed: %v", err)
	}

	// Verify we got some environments
	if len(environments) == 0 {
		t.Skip("No environments available to test obtain-editing-secret")
	}

	// Use the first environment for testing
	environment := environments[0]
	t.Logf("Testing with environment: %s (ID: %s)", environment.Name, environment.ID)

	// Test ObtainEditingSecret method
	secret, err := client.ObtainEditingSecret(environment.ID)
	if err != nil {
		t.Errorf("ObtainEditingSecret failed: %v", err)
	}

	// Verify we got a secret
	if secret == "" {
		t.Error("Obtained editing secret is empty")
	}

	t.Logf("ObtainEditingSecret test passed successfully. Secret : %s", secret)
}
