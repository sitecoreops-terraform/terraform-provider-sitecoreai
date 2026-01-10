package apiclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCMClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/clients/v1/cm" {
			t.Errorf("Expected request to /api/clients/v1/cm, got %s", r.URL.Path)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"name": "test-cm-client",
			"description": "Test CM Client",
			"clientId": "test-client-id",
			"clientSecret": "test-client-secret"
		}`))
	}))
	defer server.Close()

	// Create a client
	client := &Client{
		BaseURL:    server.URL,
		AuthURL:    "https://auth.sitecorecloud.io/oauth/token",
		ClientID:   "test-client-id",
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	// Test the method by calling doRequest directly with a mock token
	// Create a mock JWT token that will pass validation
	client.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgwNTMxOTAsImlhdCI6MTc2ODA0OTU5MCwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	// Test the method
	response, err := client.CreateCMClient("test-project", "test-environment", "test-cm-client", "Test CM Client")
	if err != nil {
		t.Fatalf("CreateCMClient failed: %v", err)
	}

	// Verify the response
	if response.Name != "test-cm-client" {
		t.Errorf("Expected name 'test-cm-client', got '%s'", response.Name)
	}
	if response.Description != "Test CM Client" {
		t.Errorf("Expected description 'Test CM Client', got '%s'", response.Description)
	}
	if response.ClientID != "test-client-id" {
		t.Errorf("Expected client ID 'test-client-id', got '%s'", response.ClientID)
	}
	if response.ClientSecret != "test-client-secret" {
		t.Errorf("Expected client secret 'test-client-secret', got '%s'", response.ClientSecret)
	}
}

func TestCreateEdgeClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/clients/v1/edge" {
			t.Errorf("Expected request to /api/clients/v1/edge, got %s", r.URL.Path)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"name": "test-edge-client",
			"description": "Test Edge Client",
			"clientId": "test-client-id",
			"clientSecret": "test-client-secret"
		}`))
	}))
	defer server.Close()

	// Create a client
	client := &Client{
		BaseURL:    server.URL,
		AuthURL:    "https://auth.sitecorecloud.io/oauth/token",
		ClientID:   "test-client-id",
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	// Test the method by calling doRequest directly with a mock token
	// Create a mock JWT token that will pass validation
	client.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgwNTMxOTAsImlhdCI6MTc2ODA0OTU5MCwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	// Test the method
	response, err := client.CreateEdgeClient("test-project", "test-environment", "test-edge-client", "Test Edge Client")
	if err != nil {
		t.Fatalf("CreateEdgeClient failed: %v", err)
	}

	// Verify the response
	if response.Name != "test-edge-client" {
		t.Errorf("Expected name 'test-edge-client', got '%s'", response.Name)
	}
	if response.Description != "Test Edge Client" {
		t.Errorf("Expected description 'Test Edge Client', got '%s'", response.Description)
	}
	if response.ClientID != "test-client-id" {
		t.Errorf("Expected client ID 'test-client-id', got '%s'", response.ClientID)
	}
	if response.ClientSecret != "test-client-secret" {
		t.Errorf("Expected client secret 'test-client-secret', got '%s'", response.ClientSecret)
	}
}

func TestCreateDeployClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/clients/v1/deploy" {
			t.Errorf("Expected request to /api/clients/v1/deploy, got %s", r.URL.Path)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"name": "test-deploy-client",
			"description": "Test Deploy Client",
			"clientId": "test-client-id",
			"clientSecret": "test-client-secret"
		}`))
	}))
	defer server.Close()

	// Create a client
	client := &Client{
		BaseURL:    server.URL,
		AuthURL:    "https://auth.sitecorecloud.io/oauth/token",
		ClientID:   "test-client-id",
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	// Test the method by calling doRequest directly with a mock token
	// Create a mock JWT token that will pass validation
	client.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgwNTMxOTAsImlhdCI6MTc2ODA0OTU5MCwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	// Test the method
	response, err := client.CreateDeployClient("test-deploy-client", "Test Deploy Client")
	if err != nil {
		t.Fatalf("CreateDeployClient failed: %v", err)
	}

	// Verify the response
	if response.Name != "test-deploy-client" {
		t.Errorf("Expected name 'test-deploy-client', got '%s'", response.Name)
	}
	if response.Description != "Test Deploy Client" {
		t.Errorf("Expected description 'Test Deploy Client', got '%s'", response.Description)
	}
	if response.ClientID != "test-client-id" {
		t.Errorf("Expected client ID 'test-client-id', got '%s'", response.ClientID)
	}
	if response.ClientSecret != "test-client-secret" {
		t.Errorf("Expected client secret 'test-client-secret', got '%s'", response.ClientSecret)
	}
}

func TestCreateEditingHostBuildClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/clients/v1/ehbuild" {
			t.Errorf("Expected request to /api/clients/v1/ehbuild, got %s", r.URL.Path)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"name": "test-ehbuild-client",
			"description": "Test Editing Host Build Client",
			"clientId": "test-client-id",
			"clientSecret": "test-client-secret"
		}`))
	}))
	defer server.Close()

	// Create a client
	client := &Client{
		BaseURL:    server.URL,
		AuthURL:    "https://auth.sitecorecloud.io/oauth/token",
		ClientID:   "test-client-id",
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	// Test the method by calling doRequest directly with a mock token
	// Create a mock JWT token that will pass validation
	client.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgwNTMxOTAsImlhdCI6MTc2ODA0OTU5MCwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	// Test the method
	response, err := client.CreateEditingHostBuildClient("test-project", "test-environment", "test-ehbuild-client", "Test Editing Host Build Client")
	if err != nil {
		t.Fatalf("CreateEditingHostBuildClient failed: %v", err)
	}

	// Verify the response
	if response.Name != "test-ehbuild-client" {
		t.Errorf("Expected name 'test-ehbuild-client', got '%s'", response.Name)
	}
	if response.Description != "Test Editing Host Build Client" {
		t.Errorf("Expected description 'Test Editing Host Build Client', got '%s'", response.Description)
	}
	if response.ClientID != "test-client-id" {
		t.Errorf("Expected client ID 'test-client-id', got '%s'", response.ClientID)
	}
	if response.ClientSecret != "test-client-secret" {
		t.Errorf("Expected client secret 'test-client-secret', got '%s'", response.ClientSecret)
	}
}

func TestDeleteClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/clients/v1/test-client-id" {
			t.Errorf("Expected request to /api/clients/v1/test-client-id, got %s", r.URL.Path)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	// Create a client
	client := &Client{
		BaseURL:    server.URL,
		AuthURL:    "https://auth.sitecorecloud.io/oauth/token",
		ClientID:   "test-client-id",
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	// Test the method by calling doRequest directly with a mock token
	// Create a mock JWT token that will pass validation
	client.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgwNTMxOTAsImlhdCI6MTc2ODA0OTU5MCwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	// Test the method
	err := client.DeleteClient("test-client-id")
	if err != nil {
		t.Fatalf("DeleteClient failed: %v", err)
	}
}
