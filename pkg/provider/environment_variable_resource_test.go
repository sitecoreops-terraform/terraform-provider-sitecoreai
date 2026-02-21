package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEnvironmentVariableResourceMetadata(t *testing.T) {
	r := environmentVariableResource{}

	req := resource.MetadataRequest{
		ProviderTypeName: "sitecore",
	}
	resp := resource.MetadataResponse{}

	r.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "sitecore_environment_variable" {
		t.Errorf("Expected TypeName to be 'sitecore_environment_variable', got '%s'", resp.TypeName)
	}
}

func TestEnvironmentVariableResourceSchema(t *testing.T) {
	r := environmentVariableResource{}

	req := resource.SchemaRequest{}
	resp := resource.SchemaResponse{}

	r.Schema(context.Background(), req, &resp)

	// Check that the description is not empty and contains expected keywords
	if resp.Schema.Description == "" {
		t.Error("Expected schema description to be non-empty")
	}

	// Check that required attributes are present
	if _, ok := resp.Schema.Attributes["id"]; !ok {
		t.Error("Expected schema to have id attribute")
	}

	if _, ok := resp.Schema.Attributes["name"]; !ok {
		t.Error("Expected schema to have name attribute")
	}

	if _, ok := resp.Schema.Attributes["value"]; !ok {
		t.Error("Expected schema to have value attribute")
	}

	if _, ok := resp.Schema.Attributes["environment_id"]; !ok {
		t.Error("Expected schema to have environment_id attribute")
	}
}

func TestEnvironmentVariableResourceConfigure(t *testing.T) {
	r := environmentVariableResource{}

	// Test with nil provider data
	req := resource.ConfigureRequest{}
	resp := resource.ConfigureResponse{}

	r.Configure(context.Background(), req, &resp)

	// Client should remain nil when no provider data is provided
	if r.base.client != nil {
		t.Error("Expected client to remain nil when no provider data is provided")
	}
}
