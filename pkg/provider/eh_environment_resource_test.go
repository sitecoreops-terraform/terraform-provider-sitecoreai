package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestEHEnvironmentResourceMetadata(t *testing.T) {
	r := ehEnvironmentResource{}

	req := resource.MetadataRequest{
		ProviderTypeName: "sitecore",
	}
	resp := resource.MetadataResponse{}

	r.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "sitecore_eh_environment" {
		t.Errorf("Expected TypeName to be 'sitecore_eh_environment', got '%s'", resp.TypeName)
	}
}

func TestEHEnvironmentResourceSchema(t *testing.T) {
	r := ehEnvironmentResource{}

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

	if _, ok := resp.Schema.Attributes["project_id"]; !ok {
		t.Error("Expected schema to have project_id attribute")
	}

	if _, ok := resp.Schema.Attributes["is_prod"]; !ok {
		t.Error("Expected schema to have is_prod attribute")
	}

	if _, ok := resp.Schema.Attributes["cm_environment_id"]; !ok {
		t.Error("Expected schema to have cm_environment_id attribute")
	}

	if _, ok := resp.Schema.Attributes["tenant_type"]; !ok {
		t.Error("Expected schema to have tenant_type attribute")
	}
}

func TestEHEnvironmentResourceConfigure(t *testing.T) {
	r := ehEnvironmentResource{}

	// Test with nil provider data
	req := resource.ConfigureRequest{}
	resp := resource.ConfigureResponse{}

	r.Configure(context.Background(), req, &resp)

	// Client should remain nil when no provider data is provided
	if r.client != nil {
		t.Error("Expected client to remain nil when no provider data is provided")
	}
}
