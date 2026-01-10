package apiclient

import (
	"os"
	"testing"
)

func TestCreateEnvironment(t *testing.T) {
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

	// Create environment
	env, err := client.CreateEnvironment(project.ID, "inttestenv", true, EnvironmentTypeCmOnly)
	if err != nil {
		t.Errorf("CreateEnvironment failed: %v", err)
	}

	t.Logf("CreateEnvironment test passed successfully. result %s", env.ID)
}

func TestDeleteEnvironment(t *testing.T) {
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

	// Find the environment with the matching name
	var foundEnvironment *Environment
	for i := range environments {
		if environments[i].Name == "inittestenv" {
			foundEnvironment = &environments[i]
			break
		}
	}

	if foundEnvironment == nil {
		t.Skip("No environments named 'inittestenv' found to test delete-environment")
	}

	// Delete environment
	err = client.DeleteEnvironment(foundEnvironment.ID)
	if err != nil {
		t.Errorf("DeleteEnvironment failed: %v", err)
	}

	t.Logf("DeleteEnvironment test passed successfully. deleted environment %s", foundEnvironment.ID)
}
