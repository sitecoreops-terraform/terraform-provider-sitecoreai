package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreai/pkg/apiclient"
)

// baseEnvironmentVariableResourceModel maps the resource schema data.
type baseEnvironmentVariableResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Value         types.String `tfsdk:"value"`
	SecretValue   types.String `tfsdk:"secret_value"`
	EnvironmentID types.String `tfsdk:"environment_id"`
}

// baseEnvironmentVariableResource contains shared logic for environment variable resources.
type baseEnvironmentVariableResource struct {
	client *apiclient.Client
}

// Configure adds the provider-configured client to the resource.
func (r *baseEnvironmentVariableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*apiclient.Client)
}

// Schema defines the schema for the base resource.
func (r *baseEnvironmentVariableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Base environment variable resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the environment variable (composite of environment_id and name)",
				Computed:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "The ID of the environment to which the variable belongs",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the environment variable",
				Required:    true,
			},
			"value": schema.StringAttribute{
				Description: "The non-sensitive value of the environment variable",
				Optional:    true,
				Sensitive:   false,
			},
			"secret_value": schema.StringAttribute{
				Description: "The sensitive value of the environment variable",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *baseEnvironmentVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, target string) {
	// Retrieve values from plan
	var plan baseEnvironmentVariableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate mutual exclusivity of value and secret_value
	if !plan.Value.IsNull() && !plan.SecretValue.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Either 'value' or 'secret_value' must be set, but not both.",
		)
		return
	}
	if plan.Value.IsNull() && plan.SecretValue.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'value' or 'secret_value' must be set.",
		)
		return
	}

	// Prepare the request body
	requestBody := apiclient.EnvironmentVariableUpsertRequestBodyDto{
		Value:  plan.Value.ValueString(),
		Secret: !plan.SecretValue.IsNull(),
		Target: &target,
	}

	// Set the environment variable using the API
	err := r.client.SetEnvironmentVariable(
		plan.EnvironmentID.ValueString(),
		plan.Name.ValueString(),
		requestBody,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating environment variable",
			"Could not create environment variable, unexpected error: "+err.Error(),
		)
		return
	}

	// Generate composite ID: environment_id:name
	compositeID := fmt.Sprintf("%s:%s", plan.EnvironmentID.ValueString(), plan.Name.ValueString())
	plan.ID = types.StringValue(compositeID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *baseEnvironmentVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, target string) {
	// Retrieve values from plan
	var plan baseEnvironmentVariableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate mutual exclusivity of value and secret_value
	if !plan.Value.IsNull() && !plan.SecretValue.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Either 'value' or 'secret_value' must be set, but not both.",
		)
		return
	}
	if plan.Value.IsNull() && plan.SecretValue.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'value' or 'secret_value' must be set.",
		)
		return
	}

	// Prepare the request body
	requestBody := apiclient.EnvironmentVariableUpsertRequestBodyDto{
		Value:  plan.Value.ValueString(),
		Secret: !plan.SecretValue.IsNull(),
		Target: &target,
	}

	// Update the environment variable using the API
	err := r.client.SetEnvironmentVariable(
		plan.EnvironmentID.ValueString(),
		plan.Name.ValueString(),
		requestBody,
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
}

// Read refreshes the Terraform state with the latest data.
func (r *baseEnvironmentVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state baseEnvironmentVariableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all environment variables from API
	variables, err := r.client.GetEnvironmentVariables(state.EnvironmentID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading environment variables",
			"Could not read environment variables: "+err.Error(),
		)
		return
	}

	// Find our specific variable
	var foundVariable *apiclient.EnvironmentVariable
	for _, variable := range variables {
		if variable.Name == state.Name.ValueString() {
			foundVariable = &variable
			break
		}
	}

	if foundVariable == nil {
		// Variable was deleted outside of Terraform
		resp.State.RemoveResource(ctx)
		return
	}

	// Update the state based on whether the variable is a secret
	if foundVariable.Secret {
		state.SecretValue = types.StringValue(foundVariable.Value)
		state.Value = types.StringNull()
	} else {
		state.Value = types.StringValue(foundVariable.Value)
		state.SecretValue = types.StringNull()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *baseEnvironmentVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state baseEnvironmentVariableResourceModel
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
