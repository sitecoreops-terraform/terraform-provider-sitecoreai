package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestEditingSecretDataSourceMetadata(t *testing.T) {
	d := editingSecretDataSource{}

	req := datasource.MetadataRequest{
		ProviderTypeName: "sitecore",
	}
	resp := datasource.MetadataResponse{}

	d.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "sitecore_editing_secret" {
		t.Errorf("Expected TypeName to be 'sitecore_editing_secret', got '%s'", resp.TypeName)
	}
}

func TestEditingSecretDataSourceSchema(t *testing.T) {
	d := editingSecretDataSource{}

	req := datasource.SchemaRequest{}
	resp := datasource.SchemaResponse{}

	d.Schema(context.Background(), req, &resp)

	if resp.Schema.Description != "Use this data source to get the editing secret for a Sitecore environment" {
		t.Errorf("Expected schema description to be 'Use this data source to get the editing secret for a Sitecore environment', got '%s'", resp.Schema.Description)
	}

	// Check that required attributes are present
	if _, ok := resp.Schema.Attributes["environment_id"]; !ok {
		t.Error("Expected schema to have environment_id attribute")
	}

	if _, ok := resp.Schema.Attributes["secret"]; !ok {
		t.Error("Expected schema to have secret attribute")
	}
}

func TestEditingSecretDataSourceConfigure(t *testing.T) {
	d := editingSecretDataSource{}

	// Test with nil provider data
	req := datasource.ConfigureRequest{}
	resp := datasource.ConfigureResponse{}

	d.Configure(context.Background(), req, &resp)

	// Client should remain nil when no provider data is provided
	if d.client != nil {
		t.Error("Expected client to remain nil when no provider data is provided")
	}
}
