package apiclient

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClientAuthentication(t *testing.T) {
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

	// Test authentication
	err = client.Authenticate()
	if err != nil {
		t.Errorf("Authentication failed: %v", err)
	}

	// Verify token is not empty
	if client.Token == "" {
		t.Error("Authentication token is empty")
	}

	t.Logf("Authentication test passed successfully: %s", client.Token)
}

func TestFindCLIUserConfig(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create .sitecore directory
	sitecoreDir := filepath.Join(tmpDir, ".sitecore")
	if err := os.Mkdir(sitecoreDir, 0755); err != nil {
		t.Fatalf("Failed to create .sitecore directory: %v", err)
	}

	// Create a test user.json file
	userJSON := `{
		"endpoints": {
			"xmCloud": {
				"host": "https://test-api.sitecorecloud.io/",
				"authority": "https://test-auth.sitecorecloud.io/",
				"accessToken": "test-access-token",
				"refreshToken": "test-refresh-token"
			}
		}
	}`

	userJSONPath := filepath.Join(sitecoreDir, "user.json")
	if err := os.WriteFile(userJSONPath, []byte(userJSON), 0644); err != nil {
		t.Fatalf("Failed to write user.json: %v", err)
	}

	// Change to the temp directory
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	err = os.Chdir(oldWD)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test finding the config
	config, err := findCLIUserConfig()
	if err != nil {
		t.Fatalf("Failed to find CLI user config: %v", err)
	}

	if config == nil {
		t.Fatal("Expected to find CLI user config, got nil")
	}

	// Verify the token
	if config.Endpoints.XMCloud.AccessToken != "test-access-token" {
		t.Errorf("Expected access token 'test-access-token', got '%s'", config.Endpoints.XMCloud.AccessToken)
	}

	if config.Endpoints.XMCloud.Host != "https://test-api.sitecorecloud.io/" {
		t.Errorf("Expected host 'https://test-api.sitecorecloud.io/', got '%s'", config.Endpoints.XMCloud.Host)
	}
}

func TestFindCLIUserConfigNotFound(t *testing.T) {
	// Create a temporary directory without user.json
	tmpDir := t.TempDir()

	// Change to the temp directory
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	err = os.Chdir(oldWD)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test that config is not found
	config, err := findCLIUserConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if config != nil {
		t.Error("Expected nil config when user.json not found")
	}
}

func TestCLIAuthenticationUsingRealConfig(t *testing.T) {
	client, err := NewClientFromCLI("")
	if err != nil {
		t.Fatalf("Failed during instantiation: %v", err)
	}

	err = client.Authenticate()
	if err != nil {
		t.Fatalf("Failed during authenticate: %v", err)
	}
}
