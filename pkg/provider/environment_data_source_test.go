package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEnvironmentDataSourceMetadata(t *testing.T) {
	d := environmentDataSource{}

	req := datasource.MetadataRequest{
		ProviderTypeName: "sitecore",
	}
	resp := datasource.MetadataResponse{}

	d.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "sitecore_environment" {
		t.Errorf("Expected TypeName to be 'sitecore_environment', got '%s'", resp.TypeName)
	}
}

func TestEnvironmentDataSourceSchema(t *testing.T) {
	d := environmentDataSource{}

	req := datasource.SchemaRequest{}
	resp := datasource.SchemaResponse{}

	d.Schema(context.Background(), req, &resp)

	if resp.Schema.Description != "Use this data source to get information about a Sitecore environment by project ID and name" {
		t.Errorf("Expected schema description to be 'Use this data source to get information about a Sitecore environment by project ID and name', got '%s'", resp.Schema.Description)
	}

	// Check that required attributes are present
	if _, ok := resp.Schema.Attributes["id"]; !ok {
		t.Error("Expected schema to have id attribute")
	}

	if _, ok := resp.Schema.Attributes["name"]; !ok {
		t.Error("Expected schema to have name attribute")
	}

	if _, ok := resp.Schema.Attributes["project_id"]; !ok {
		t.Error("Expected schema to have project_id attribute")
	}
}

func TestEnvironmentDataSourceConfigure(t *testing.T) {
	d := environmentDataSource{}

	// Test with nil provider data
	req := datasource.ConfigureRequest{}
	resp := datasource.ConfigureResponse{}

	d.Configure(context.Background(), req, &resp)

	// Client should remain nil when no provider data is provided
	if d.client != nil {
		t.Error("Expected client to remain nil when no provider data is provided")
	}
}
