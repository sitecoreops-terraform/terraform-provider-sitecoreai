// Environment variable resource implementation
package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreai/pkg/apiclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &environmentVariableResource{}
	_ resource.ResourceWithConfigure   = &environmentVariableResource{}
	_ resource.ResourceWithImportState = &environmentVariableResource{}
)

// environmentVariableResourceModel maps the resource schema data.
type environmentVariableResourceModel struct {
	baseEnvironmentVariableResourceModel
	Target types.String `tfsdk:"target"`
}

// environmentVariableResource extends the base resource to support the target attribute.
type environmentVariableResource struct {
	base baseEnvironmentVariableResource
}

// NewEnvironmentVariableResource creates a new environment variable resource.
func NewEnvironmentVariableResource() resource.Resource {
	return &environmentVariableResource{}
}

// Metadata returns the resource type name.
func (r *environmentVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_variable"
}

// Schema defines the schema for the resource.
func (r *environmentVariableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Get the base schema
	r.base.Schema(ctx, req, resp)

	// Add the 'target' attribute
	resp.Schema.Attributes["target"] = schema.StringAttribute{
		Description: "The target for the environment variable (CM, EH, or custom editing host name). Leave empty for all targets.",
		Optional:    true,
	}
	resp.Schema.Description = "Environments ¤ Manages an environment variable for SitecoreAI environment."
}

// Configure adds the provider-configured client to the resource.
func (r *environmentVariableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.base.Configure(ctx, req, resp)
}

// Create creates the resource and sets the initial Terraform state.
func (r *environmentVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan environmentVariableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delegate to the base Create method with the target field
	target := plan.Target.ValueString()
	r.base.Create(ctx, resource.CreateRequest{Plan: req.Plan}, resp, target)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate composite ID: environment_id:target:name
	var compositeID string
	if target == "" {
		compositeID = fmt.Sprintf("%s::%s", plan.EnvironmentID.ValueString(), plan.Name.ValueString())
	} else {
		compositeID = fmt.Sprintf("%s:%s:%s", plan.EnvironmentID.ValueString(), target, plan.Name.ValueString())
	}
	plan.ID = types.StringValue(compositeID)

	// Set the state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *environmentVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan environmentVariableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delegate to the base Update method with the target field
	target := plan.Target.ValueString()
	r.base.Update(ctx, resource.UpdateRequest{Plan: req.Plan}, resp, target)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *environmentVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Delegate to the base Read method
	r.base.Read(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *environmentVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delegate to the base Delete method
	r.base.Delete(ctx, req, resp)
}

// ImportState imports an existing environment variable into Terraform state.
func (r *environmentVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: environment_id:target:name (target is optional)
	idParts := strings.Split(req.ID, ":")
	if len(idParts) < 2 || idParts[0] == "" || idParts[len(idParts)-1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: environment_id:target:name (target is optional)",
		)
		return
	}

	environmentID := idParts[0]
	variableName := idParts[len(idParts)-1]
	target := ""
	if len(idParts) == 3 {
		target = idParts[1]
	}

	// Generate composite ID: environment_id:target:name
	compositeID := fmt.Sprintf("%s:%s:%s", environmentID, target, variableName)

	// Set the composite ID and individual attributes
	var state environmentVariableResourceModel
	state.ID = types.StringValue(compositeID)
	state.EnvironmentID = types.StringValue(environmentID)
	state.Name = types.StringValue(variableName)
	state.Target = types.StringValue(target)

	// Fetch the variable value from the API to ensure it exists
	variables, err := r.base.client.GetEnvironmentVariables(environmentID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading environment variables during import",
			"Could not read environment variables: "+err.Error(),
		)
		return
	}

	// Find our specific variable
	var foundVariable *apiclient.EnvironmentVariable
	for _, variable := range variables {
		if variable.Name == variableName && (target == "" || variable.Target == target) {
			foundVariable = &variable
			break
		}
	}

	if foundVariable == nil {
		resp.Diagnostics.AddError(
			"Environment variable not found",
			fmt.Sprintf("Environment variable '%s' not found in environment '%s'", variableName, environmentID),
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
