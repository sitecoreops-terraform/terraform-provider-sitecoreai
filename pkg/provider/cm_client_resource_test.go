package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCMClientResourceMetadata(t *testing.T) {
	r := cmClientResource{}

	req := resource.MetadataRequest{
		ProviderTypeName: "sitecore",
	}
	resp := resource.MetadataResponse{}

	r.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "sitecore_cm_client" {
		t.Errorf("Expected TypeName to be 'sitecore_cm_client', got '%s'", resp.TypeName)
	}
}

func TestCMClientResourceSchema(t *testing.T) {
	r := cmClientResource{}

	req := resource.SchemaRequest{}
	resp := resource.SchemaResponse{}

	r.Schema(context.Background(), req, &resp)

	if resp.Schema.Description != "Automation Clients ¤ Manages a Sitecore CM automation client" {
		t.Errorf("Expected schema description to be 'Automation Clients ¤ Manages a Sitecore CM automation client', got '%s'", resp.Schema.Description)
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

	if _, ok := resp.Schema.Attributes["environment_id"]; !ok {
		t.Error("Expected schema to have environment_id attribute")
	}
}

func TestCMClientResourceConfigure(t *testing.T) {
	r := cmClientResource{}

	// Test with nil provider data
	req := resource.ConfigureRequest{}
	resp := resource.ConfigureResponse{}

	r.Configure(context.Background(), req, &resp)

	// Client should remain nil when no provider data is provided
	if r.client != nil {
		t.Error("Expected client to remain nil when no provider data is provided")
	}
}
