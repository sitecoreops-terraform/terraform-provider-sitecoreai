package apiclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestObtainEditingSecret_MockedResponse(t *testing.T) {
	t.Run("Successful response should return secret", func(t *testing.T) {
		// Create a mock HTTP server that returns success
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, "test-secret-key-12345")
		}))
		defer server.Close()

		// Create a client that uses the mock server
		client := &Client{
			BaseURL:    server.URL,
			HTTPClient: server.Client(),
			Token:      "test-token",
		}

		// Call ObtainEditingSecret
		secret, err := client.ObtainEditingSecret("test-env-id")
		if err != nil {
			t.Fatalf("ObtainEditingSecret failed: %v", err)
		}

		// Verify we got the secret
		if secret != "test-secret-key-12345" {
			t.Errorf("Expected 'test-secret-key-12345', got '%s'", secret)
		}
	})

	t.Run("Not Found errour should return empty string as it is awaiting deployment", func(t *testing.T) {
		// Create a mock HTTP server that returns 404 with the expected JSON format
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
			errorResponse := `{
				"type":"https://tools.ietf.org/html/rfc9110#section-15.5.5",
				"title":"Not Found",
				"status":404,
				"traceId":"00-d472dcd68fbe5895cd474e70312e2110-3155286964e7ddd8-01"
			}`
			_, _ = fmt.Fprint(w, errorResponse)
		}))
		defer server.Close()

		// Create a client that uses the mock server
		client := &Client{
			BaseURL:    server.URL,
			HTTPClient: server.Client(),
			Token:      "test-token",
		}

		// Call ObtainEditingSecret
		secret, err := client.ObtainEditingSecret("test-env-id")
		if err != nil {
			t.Fatalf("ObtainEditingSecret failed: %v", err)
		}

		// Verify we got an empty string for 404
		if secret != "" {
			t.Errorf("Expected empty string for 404, got '%s'", secret)
		}
	})

	t.Run("Other errors should return error", func(t *testing.T) {
		// Create a mock HTTP server that returns 500 error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			errorResponse := `{
				"type":"https://tools.ietf.org/html/rfc9110#section-15.5.1",
				"title":"Internal Server Error",
				"status":500,
				"traceId":"00-d472dcd68fbe5895cd474e70312e2110-3155286964e7ddd8-01"
			}`
			_, _ = fmt.Fprint(w, errorResponse)
		}))
		defer server.Close()

		// Create a client that uses the mock server
		client := &Client{
			BaseURL:    server.URL,
			HTTPClient: server.Client(),
			Token:      "test-token",
		}

		// Call ObtainEditingSecret
		_, err := client.ObtainEditingSecret("test-env-id")

		// Should return an error for 500
		if err == nil {
			t.Error("Expected error for 500 status, got nil")
		}
	})
}
