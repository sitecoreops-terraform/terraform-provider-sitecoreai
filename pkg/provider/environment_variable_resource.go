// Environment variable resource implementation
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreai/pkg/apiclient"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &environmentVariableResource{}
	_ resource.ResourceWithConfigure   = &environmentVariableResource{}
	_ resource.ResourceWithImportState = &environmentVariableResource{}
)

// NewEnvironmentVariableResource is a helper function to simplify the provider implementation
func NewEnvironmentVariableResource() resource.Resource {
	return &environmentVariableResource{}
}

// environmentVariableResource is the resource implementation
type environmentVariableResource struct {
	client *apiclient.Client
}

// environmentVariableResourceModel maps the resource schema data
type environmentVariableResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Value         types.String `tfsdk:"value"`
	EnvironmentID types.String `tfsdk:"environment_id"`
}

// Metadata returns the resource type name
func (r *environmentVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_variable"
}

// Schema defines the schema for the resource
func (r *environmentVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manages an environment variable for a Sitecore environment.
		
		Environment variables are key-value pairs that can be used to configure
		environment-specific settings for Sitecore environments.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the environment variable (composite of environment_id and name)",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Description: "The ID of the environment to which the variable belongs",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the environment variable",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "The value of the environment variable",
				Required:    true,
				Sensitive:   false,
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *environmentVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*apiclient.Client)
}

// Create creates the resource and sets the initial Terraform state
func (r *environmentVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan environmentVariableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the environment variable using the API
	err := r.client.SetEnvironmentVariable(
		plan.EnvironmentID.ValueString(),
		plan.Name.ValueString(),
		plan.Value.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating environment variable",
			"Could not create environment variable, unexpected error: "+err.Error(),
		)
		return
	}

	// Generate composite ID: project_id:environment_id:name
	compositeID := fmt.Sprintf("%s:%s",
		plan.EnvironmentID.ValueString(),
		plan.Name.ValueString(),
	)

	// Set the composite ID and other attributes
	plan.ID = types.StringValue(compositeID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data
func (r *environmentVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state environmentVariableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all environment variables from API
	variables, err := r.client.GetEnvironmentVariables(
		state.EnvironmentID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading environment variables",
			"Could not read environment variables: "+err.Error(),
		)
		return
	}

	// Check if our specific variable exists
	variableValue, exists := variables[state.Name.ValueString()]
	if !exists {
		// Variable was deleted outside of Terraform
		resp.State.RemoveResource(ctx)
		return
	}

	// Update the value from the API response
	state.Value = types.StringValue(variableValue)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success
func (r *environmentVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan environmentVariableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the environment variable using the API
	err := r.client.SetEnvironmentVariable(
		plan.EnvironmentID.ValueString(),
		plan.Name.ValueString(),
		plan.Value.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating environment variable",
			"Could not update environment variable, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success
func (r *environmentVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state environmentVariableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the environment variable
	err := r.client.DeleteEnvironmentVariable(
		state.EnvironmentID.ValueString(),
		state.Name.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting environment variable",
			"Could not delete environment variable, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing environment variable into Terraform state
func (r *environmentVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: environment_id,name
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
