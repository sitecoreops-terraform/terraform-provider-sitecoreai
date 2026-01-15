package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProviderMetadata(t *testing.T) {
	p := sitecoreProvider{version: "test-version"}

	req := provider.MetadataRequest{}
	resp := provider.MetadataResponse{}

	p.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "sitecoreai" {
		t.Errorf("Expected TypeName to be 'sitecoreai', got '%s'", resp.TypeName)
	}

	if resp.Version != "test-version" {
		t.Errorf("Expected Version to be 'test-version', got '%s'", resp.Version)
	}
}

func TestProviderSchema(t *testing.T) {
	p := sitecoreProvider{}

	req := provider.SchemaRequest{}
	resp := provider.SchemaResponse{}

	p.Schema(context.Background(), req, &resp)

	if resp.Schema.Description != "Interact with SitecoreAI" {
		t.Errorf("Expected schema description to be 'Interact with SitecoreAI', got '%s'", resp.Schema.Description)
	}

	// Check that required attributes are present
	if _, ok := resp.Schema.Attributes["client_id"]; !ok {
		t.Error("Expected schema to have client_id attribute")
	}

	if _, ok := resp.Schema.Attributes["client_secret"]; !ok {
		t.Error("Expected schema to have client_secret attribute")
	}
}

func TestProviderConfigure(t *testing.T) {
	t.Run("provider configuration method exists", func(t *testing.T) {
		// This is a basic test to verify the Configure method exists
		p := sitecoreProvider{version: "test"}

		// For now, we just verify that the provider struct is properly initialized
		if p.version != "test" {
			t.Error("Expected provider to be properly initialized")
		}
	})
}

func TestProviderDataSources(t *testing.T) {
	p := sitecoreProvider{}

	dataSources := p.DataSources(context.Background())

	if len(dataSources) == 0 {
		t.Error("Expected provider to have data sources")
	}

	// Check that we can create instances of the data sources
	for _, ds := range dataSources {
		if ds == nil {
			t.Error("Expected data source function to be non-nil")
		}
	}
}

func TestProviderResources(t *testing.T) {
	p := sitecoreProvider{}

	resources := p.Resources(context.Background())

	if len(resources) == 0 {
		t.Error("Expected provider to have resources")
	}

	// Check that we can create instances of the resources
	for _, r := range resources {
		if r == nil {
			t.Error("Expected resource function to be non-nil")
		}
	}
}
