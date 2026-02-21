// Environment variable resource implementation
package provider

import (
	"context"
	"fmt"
	"strings"

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
	SecretValue   types.String `tfsdk:"secret_value"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	Target        types.String `tfsdk:"target"`
}

// Metadata returns the resource type name
func (r *environmentVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_variable"
}

// Schema defines the schema for the resource
func (r *environmentVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Environments ¤ Manages an environment variable for a Sitecore environment.`,
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
				Description: "The non-sensitive value of the environment variable",
				Optional:    true,
				Sensitive:   false,
			},
			"secret_value": schema.StringAttribute{
				Description: "The sensitive value of the environment variable",
				Optional:    true,
				Sensitive:   true,
			},
			"target": schema.StringAttribute{
				Description: "The target for the environment variable (CM, EH, or custom editing host name). Leave empty for all targets.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{}, // Add validation for mutual exclusivity of value/secret_value
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
	}

	// Set target if provided
	target := plan.Target.ValueString()
	if target != "" {
		requestBody.Target = &target
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

	// Generate composite ID: environment_id:target:name
	var compositeID string
	if target == "" {
		compositeID = fmt.Sprintf("%s::%s", plan.EnvironmentID.ValueString(), plan.Name.ValueString())
	} else {
		compositeID = fmt.Sprintf("%s:%s:%s", plan.EnvironmentID.ValueString(), target, plan.Name.ValueString())
	}
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
	}

	// Set target if provided
	target := plan.Target.ValueString()
	if target != "" {
		requestBody.Target = &target
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
	// Expected format: environment_id:target:name (consistent with Create method)
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
	var compositeID string
	if target == "" {
		compositeID = fmt.Sprintf("%s::%s", environmentID, variableName)
	} else {
		compositeID = fmt.Sprintf("%s:%s:%s", environmentID, target, variableName)
	}

	// Set the composite ID and individual attributes
	var state environmentVariableResourceModel
	state.ID = types.StringValue(compositeID)
	state.EnvironmentID = types.StringValue(environmentID)
	state.Name = types.StringValue(variableName)
	state.Target = types.StringValue(target)

	// Fetch the variable value from the API to ensure it exists
	variables, err := r.client.GetEnvironmentVariables(environmentID)
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
		if variable.Name == variableName {
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
