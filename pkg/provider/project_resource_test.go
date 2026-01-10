package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestProjectResourceMetadata(t *testing.T) {
	r := projectResource{}

	req := resource.MetadataRequest{
		ProviderTypeName: "sitecore",
	}
	resp := resource.MetadataResponse{}

	r.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "sitecore_project" {
		t.Errorf("Expected TypeName to be 'sitecore_project', got '%s'", resp.TypeName)
	}
}

func TestProjectResourceSchema(t *testing.T) {
	r := projectResource{}

	req := resource.SchemaRequest{}
	resp := resource.SchemaResponse{}

	r.Schema(context.Background(), req, &resp)

	if resp.Schema.Description != "Manages a Sitecore project" {
		t.Errorf("Expected schema description to be 'Manages a Sitecore project', got '%s'", resp.Schema.Description)
	}

	// Check that required attributes are present
	if _, ok := resp.Schema.Attributes["id"]; !ok {
		t.Error("Expected schema to have id attribute")
	}

	if _, ok := resp.Schema.Attributes["name"]; !ok {
		t.Error("Expected schema to have name attribute")
	}

	if _, ok := resp.Schema.Attributes["description"]; !ok {
		t.Error("Expected schema to have description attribute")
	}

	// Check attribute properties - simplified for now
	// Note: A more comprehensive test would check the exact attribute properties
	// but this requires more complex attribute type checking
}

func TestProjectResourceConfigure(t *testing.T) {
	r := projectResource{}

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
// func TestProjectResourceImportState(t *testing.T) {
//  	r := projectResource{}
//  	req := resource.ImportStateRequest{ID: "test-project-id"}
//  	resp := resource.ImportStateResponse{}
//  	r.ImportState(context.Background(), req, &resp)
// }
