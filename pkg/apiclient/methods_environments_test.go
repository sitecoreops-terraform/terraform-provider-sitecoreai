package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateEnvironment_EditingHostEnvironmentDetails(t *testing.T) {
	t.Run("Without cmEnvironmentId - EditingHostEnvironmentDetails should be omitted", func(t *testing.T) {
		// Create a mock HTTP server
		requestBody := ""
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read body", http.StatusInternalServerError)
				return
			}
			requestBody = string(body)

			// Return a mock response
			w.WriteHeader(http.StatusOK)
			mockResponse := `{
				"id": "test-environment-id",
				"name": "test-environment",
				"projectId": "test-project-id",
				"type": "Development"
			}`
			_, _ = fmt.Fprint(w, mockResponse)
		}))
		defer server.Close()

		// Create a client that uses the mock server
		client := &Client{
			BaseURL:    server.URL,
			HTTPClient: server.Client(),
			Token:      "test-token",
		}

		// Call CreateEnvironment without cmEnvironmentId
		createdEnv, err := client.CreateEnvironment("test-project-id", "test-environment", false, EnvironmentTypeCombined, "")
		if err != nil {
			t.Fatalf("CreateEnvironment failed: %v", err)
		}

		// Verify the environment was "created" (mock response)
		if createdEnv == nil {
			t.Fatal("Expected created environment, got nil")
		}

		// Parse the captured request body using json.NewDecoder
		var requestData map[string]interface{}
		err = json.NewDecoder(strings.NewReader(requestBody)).Decode(&requestData)
		if err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		// Verify EditingHostEnvironmentDetails is not present
		if _, exists := requestData["editingHostEnvironmentDetails"]; exists {
			t.Errorf("EditingHostEnvironmentDetails should not be present when cmEnvironmentId is empty. Request data: %+v", requestData)
		}

		// Verify other expected fields are present
		if requestData["name"] != "test-environment" {
			t.Errorf("Expected name 'test-environment', got '%v'", requestData["name"])
		}
		// For EnvironmentTypeCombined, type should not be present or should be empty
		if requestData["type"] != nil && requestData["type"] != "" {
			t.Errorf("Expected type to be nil or empty for EnvironmentTypeCombined, got '%v'", requestData["type"])
		}
	})

	t.Run("With cmEnvironmentId - EditingHostEnvironmentDetails should be included", func(t *testing.T) {
		// Create a mock HTTP server
		requestBody := ""
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read body", http.StatusInternalServerError)
				return
			}
			requestBody = string(body)

			// Return a mock response
			w.WriteHeader(http.StatusOK)
			mockResponse := `{
				"id": "test-environment-id",
				"name": "test-environment",
				"projectId": "test-project-id",
				"type": "Development"
			}`
			_, _ = fmt.Fprint(w, mockResponse)
		}))
		defer server.Close()

		// Create a client that uses the mock server
		client := &Client{
			BaseURL:    server.URL,
			HTTPClient: server.Client(),
			Token:      "test-token",
		}

		// Call CreateEnvironment with cmEnvironmentId
		createdEnv, err := client.CreateEnvironment("test-project-id", "test-environment", false, EnvironmentTypeCombined, "test-cm-env-id")
		if err != nil {
			t.Fatalf("CreateEnvironment failed: %v", err)
		}

		// Verify the environment was "created" (mock response)
		if createdEnv == nil {
			t.Fatal("Expected created environment, got nil")
		}

		// Parse the captured request body using json.NewDecoder
		var requestData map[string]interface{}
		err = json.NewDecoder(strings.NewReader(requestBody)).Decode(&requestData)
		if err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		// Verify EditingHostEnvironmentDetails is present
		if _, exists := requestData["editingHostEnvironmentDetails"]; !exists {
			t.Error("EditingHostEnvironmentDetails should be present when cmEnvironmentId is provided")
		}

		// Verify the cmEnvironmentId is correct
		editingDetails := requestData["editingHostEnvironmentDetails"].(map[string]interface{})
		if editingDetails["cmEnvironmentId"] != "test-cm-env-id" {
			t.Errorf("Expected cmEnvironmentId 'test-cm-env-id', got '%v'", editingDetails["cmEnvironmentId"])
		}

		// Verify other expected fields are present
		if requestData["name"] != "test-environment" {
			t.Errorf("Expected name 'test-environment', got '%v'", requestData["name"])
		}
		// For EnvironmentTypeCombined, type should not be present or should be empty
		if requestData["type"] != nil && requestData["type"] != "" {
			t.Errorf("Expected type to be nil or empty for EnvironmentTypeCombined, got '%v'", requestData["type"])
		}
	})
}
