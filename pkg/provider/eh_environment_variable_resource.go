package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreai/pkg/apiclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ehEnvironmentVariableResource{}
	_ resource.ResourceWithConfigure   = &ehEnvironmentVariableResource{}
	_ resource.ResourceWithImportState = &ehEnvironmentVariableResource{}
)

// ehEnvironmentVariableResourceModel maps the resource schema data.
type ehEnvironmentVariableResourceModel struct {
	baseEnvironmentVariableResourceModel
}

// ehEnvironmentVariableResource delegates to the base resource with target="EH".
type ehEnvironmentVariableResource struct {
	base baseEnvironmentVariableResource
}

// NewEHEnvironmentVariableResource creates a new EH environment variable resource.
func NewEHEnvironmentVariableResource() resource.Resource {
	return &ehEnvironmentVariableResource{}
}

// Metadata returns the resource type name.
func (r *ehEnvironmentVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eh_environment_variable"
}

// Schema defines the schema for the resource.
func (r *ehEnvironmentVariableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Get the base schema
	r.base.Schema(ctx, req, resp)

	// Remove the 'target' attribute (hardcoded to "EH")
	delete(resp.Schema.Attributes, "target")
	resp.Schema.Description = "Environments ¤ Manages an environment variable for a SitecoreAI editing host environment."
}

// Configure adds the provider-configured client to the resource.
func (r *ehEnvironmentVariableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.base.Configure(ctx, req, resp)
}

// Create creates the resource and sets the initial Terraform state.
func (r *ehEnvironmentVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Delegate to the base Create method with target="EH"
	r.base.Create(ctx, req, resp, "EH")
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state (ID is already set by base.Create)
	var plan ehEnvironmentVariableResourceModel
	diags := resp.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ehEnvironmentVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Delegate to the base Update method with target="EH"
	r.base.Update(ctx, req, resp, "EH")
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	var plan ehEnvironmentVariableResourceModel
	diags := resp.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ehEnvironmentVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Delegate to the base Read method
	r.base.Read(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ehEnvironmentVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delegate to the base Delete method
	r.base.Delete(ctx, req, resp)
}

// ImportState imports an existing environment variable into Terraform state.
func (r *ehEnvironmentVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: environment_id:name
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: environment_id:name",
		)
		return
	}

	environmentID := idParts[0]
	variableName := idParts[1]

	// Generate composite ID: environment_id:name
	compositeID := fmt.Sprintf("%s:%s", environmentID, variableName)

	// Set the composite ID and individual attributes
	var state ehEnvironmentVariableResourceModel
	state.ID = types.StringValue(compositeID)
	state.EnvironmentID = types.StringValue(environmentID)
	state.Name = types.StringValue(variableName)

	// Fetch the variable value from the API to ensure it exists
	variables, err := r.base.client.GetEnvironmentVariables(environmentID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading environment variables during import",
			"Could not read environment variables: "+err.Error(),
		)
		return
	}

	// Find our specific variable (must have target="EH")
	var foundVariable *apiclient.EnvironmentVariable
	for _, variable := range variables {
		if variable.Name == variableName && variable.Target == "EH" {
			foundVariable = &variable
			break
		}
	}

	if foundVariable == nil {
		resp.Diagnostics.AddError(
			"Environment variable not found",
			fmt.Sprintf("EH environment variable '%s' not found in environment '%s'", variableName, environmentID),
		)
		return
	}

	// Set the value based on whether it's a secret
	if foundVariable.Secret {
		state.SecretValue = types.StringValue(foundVariable.Value)
		state.Value = types.StringNull()
	} else {
		state.Value = types.StringValue(foundVariable.Value)
		state.SecretValue = types.StringNull()
	}

	// Set the state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
