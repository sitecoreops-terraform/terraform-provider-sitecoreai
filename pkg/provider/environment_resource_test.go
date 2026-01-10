package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnvironmentResourceMetadata(t *testing.T) {
	r := environmentResource{}

	req := resource.MetadataRequest{
		ProviderTypeName: "sitecore",
	}
	resp := resource.MetadataResponse{}

	r.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "sitecore_environment" {
		t.Errorf("Expected TypeName to be 'sitecore_environment', got '%s'", resp.TypeName)
	}
}

func TestEnvironmentResourceSchema(t *testing.T) {
	r := environmentResource{}

	req := resource.SchemaRequest{}
	resp := resource.SchemaResponse{}

	r.Schema(context.Background(), req, &resp)

	if resp.Schema.Description != "Manages a Sitecore environment" {
		t.Errorf("Expected schema description to be 'Manages a Sitecore environment', got '%s'", resp.Schema.Description)
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

	if _, ok := resp.Schema.Attributes["is_prod"]; !ok {
		t.Error("Expected schema to have is_prod attribute")
	}

	if _, ok := resp.Schema.Attributes["tenant_type"]; !ok {
		t.Error("Expected schema to have tenant_type attribute")
	}
}

func TestEnvironmentResourceConfigure(t *testing.T) {
	r := environmentResource{}

	// Test with nil provider data
	req := resource.ConfigureRequest{}
	resp := resource.ConfigureResponse{}

	r.Configure(context.Background(), req, &resp)

	// Client should remain nil when no provider data is provided
	if r.client != nil {
		t.Error("Expected client to remain nil when no provider data is provided")
	}
}

// Note: ImportState test is skipped for now as it requires complex state setup
// func TestEnvironmentResourceImportState(t *testing.T) {
//  	r := environmentResource{}
//  	req := resource.ImportStateRequest{ID: "test-project-id/test-environment-id"}
//  	resp := resource.ImportStateResponse{}
//  	r.ImportState(context.Background(), req, &resp)
// }
