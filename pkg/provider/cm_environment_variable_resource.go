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
	_ resource.Resource                = &cmEnvironmentVariableResource{}
	_ resource.ResourceWithConfigure   = &cmEnvironmentVariableResource{}
	_ resource.ResourceWithImportState = &cmEnvironmentVariableResource{}
)

// cmEnvironmentVariableResourceModel maps the resource schema data.
type cmEnvironmentVariableResourceModel struct {
	baseEnvironmentVariableResourceModel
}

// cmEnvironmentVariableResource delegates to the base resource with target="CM".
type cmEnvironmentVariableResource struct {
	base baseEnvironmentVariableResource
}

// NewCMEnvironmentVariableResource creates a new CM environment variable resource.
func NewCMEnvironmentVariableResource() resource.Resource {
	return &cmEnvironmentVariableResource{}
}

// Metadata returns the resource type name.
func (r *cmEnvironmentVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_environment_variable"
}

// Schema defines the schema for the resource.
func (r *cmEnvironmentVariableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Get the base schema
	r.base.Schema(ctx, req, resp)

	// Remove the 'target' attribute (hardcoded to "CM")
	delete(resp.Schema.Attributes, "target")
	resp.Schema.Description = "Environments ¤ Manages an environment variable for a SitecoreAI CM environment."
}

// Configure adds the provider-configured client to the resource.
func (r *cmEnvironmentVariableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.base.Configure(ctx, req, resp)
}

// Create creates the resource and sets the initial Terraform state.
func (r *cmEnvironmentVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Delegate to the base Create method with target="CM"
	r.base.Create(ctx, req, resp, "CM")
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state (ID is already set by base.Create)
	var plan cmEnvironmentVariableResourceModel
	diags := resp.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *cmEnvironmentVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Delegate to the base Update method with target="CM"
	r.base.Update(ctx, req, resp, "CM")
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	var plan cmEnvironmentVariableResourceModel
	diags := resp.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *cmEnvironmentVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Delegate to the base Read method
	r.base.Read(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *cmEnvironmentVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delegate to the base Delete method
	r.base.Delete(ctx, req, resp)
}

// ImportState imports an existing environment variable into Terraform state.
func (r *cmEnvironmentVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
	var state cmEnvironmentVariableResourceModel
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

	// Find our specific variable (must have target="CM")
	var foundVariable *apiclient.EnvironmentVariable
	for _, variable := range variables {
		if variable.Name == variableName && variable.Target == "CM" {
			foundVariable = &variable
			break
		}
	}

	if foundVariable == nil {
		resp.Diagnostics.AddError(
			"Environment variable not found",
			fmt.Sprintf("CM environment variable '%s' not found in environment '%s'", variableName, environmentID),
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
