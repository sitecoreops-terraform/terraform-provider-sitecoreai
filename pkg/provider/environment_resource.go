// Environment resource implementation
package provider

import (
	"context"

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
	_ resource.Resource                = &environmentResource{}
	_ resource.ResourceWithConfigure   = &environmentResource{}
	_ resource.ResourceWithImportState = &environmentResource{}
)

// NewEnvironmentResource is a helper function to simplify the provider implementation
func NewEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

// environmentResource is the resource implementation
type environmentResource struct {
	client *apiclient.Client
}

// environmentResourceModel maps the resource schema data
type environmentResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	ProjectID               types.String `tfsdk:"project_id"`
	IsProd                  types.Bool   `tfsdk:"is_prod"`
	TenantType              types.String `tfsdk:"tenant_type"`
	Host                    types.String `tfsdk:"host"`
	PlatformTenantId        types.String `tfsdk:"platform_tenant_id"`
	PlatformTenantName      types.String `tfsdk:"platform_tenant_name"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	LastUpdatedBy           types.String `tfsdk:"last_updated_by"`
	LastUpdatedAt           types.String `tfsdk:"last_updated_at"`
	IsDeleted               types.Bool   `tfsdk:"is_deleted"`
	PreviewContextId        types.String `tfsdk:"preview_context_id"`
	LiveContextId           types.String `tfsdk:"live_context_id"`
	HighAvailabilityEnabled types.Bool   `tfsdk:"high_availability_enabled"`
}

// Metadata returns the resource type name
func (r *environmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// Schema defines the schema for the resource
func (r *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Environments ¤ Manages a traditional SitecoreAI combined environment with both authoring and editing.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the environment",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the environment",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The ID of the project to which the environment belongs",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_prod": schema.BoolAttribute{
				Description: "Whether this is a production environment",
				Optional:    true,
			},
			"tenant_type": schema.StringAttribute{
				Description: "Indicates if it is production or not, can have the values 'prod' or 'nonprod'",
				Computed:    true,
			},
			"host": schema.StringAttribute{
				Description: "The host of the environment",
				Computed:    true,
			},
			"platform_tenant_id": schema.StringAttribute{
				Description: "The platform tenant ID",
				Computed:    true,
			},
			"platform_tenant_name": schema.StringAttribute{
				Description: "The platform tenant name",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When the environment was created",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "Who created the environment",
				Computed:    true,
			},
			"last_updated_by": schema.StringAttribute{
				Description: "Who last updated the environment",
				Computed:    true,
			},
			"last_updated_at": schema.StringAttribute{
				Description: "When the environment was last updated",
				Computed:    true,
			},
			"is_deleted": schema.BoolAttribute{
				Description: "Whether the environment is deleted",
				Computed:    true,
			},
			"preview_context_id": schema.StringAttribute{
				Description: "The preview context ID",
				Computed:    true,
			},
			"live_context_id": schema.StringAttribute{
				Description: "The live context ID",
				Computed:    true,
			},
			"high_availability_enabled": schema.BoolAttribute{
				Description: "Whether high availability is enabled",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *environmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*apiclient.Client)
}

// Create creates the resource and sets the initial Terraform state
func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan environmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse is_prod parameter (default: false)
	isProd := false
	if !plan.IsProd.IsNull() && !plan.IsProd.IsUnknown() {
		isProd = plan.IsProd.ValueBool()
	}

	// Call API with Combined environment type
	createdEnvironment, err := r.client.CreateEnvironment(
		plan.ProjectID.ValueString(),
		plan.Name.ValueString(),
		isProd,
		apiclient.EnvironmentTypeCombined,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating environment",
			"Could not create environment, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(createdEnvironment.ID)
	plan.Name = types.StringValue(createdEnvironment.Name)
	plan.ProjectID = types.StringValue(createdEnvironment.ProjectID)
	plan.Host = types.StringValue(createdEnvironment.Host)
	plan.PlatformTenantId = types.StringValue(createdEnvironment.PlatformTenantId)
	plan.PlatformTenantName = types.StringValue(createdEnvironment.PlatformTenantName)
	plan.TenantType = types.StringValue(createdEnvironment.TenantType)
	plan.CreatedAt = types.StringValue(createdEnvironment.CreatedAt)
	plan.CreatedBy = types.StringValue(createdEnvironment.CreatedBy)
	plan.LastUpdatedBy = types.StringValue(createdEnvironment.LastUpdatedBy)
	plan.LastUpdatedAt = types.StringValue(createdEnvironment.LastUpdatedAt)
	plan.IsDeleted = types.BoolValue(createdEnvironment.IsDeleted)
	plan.PreviewContextId = types.StringValue(createdEnvironment.PreviewContextId)
	plan.LiveContextId = types.StringValue(createdEnvironment.LiveContextId)
	plan.HighAvailabilityEnabled = types.BoolValue(createdEnvironment.HighAvailabilityEnabled)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data
func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state environmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get environment from API
	environment, err := r.client.GetEnvironment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading environment",
			"Could not read environment ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Handle environment not found (soft delete scenario)
	if environment == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(environment.Name)
	state.Host = types.StringValue(environment.Host)
	state.PlatformTenantId = types.StringValue(environment.PlatformTenantId)
	state.PlatformTenantName = types.StringValue(environment.PlatformTenantName)
	state.TenantType = types.StringValue(environment.TenantType)
	state.CreatedAt = types.StringValue(environment.CreatedAt)
	state.CreatedBy = types.StringValue(environment.CreatedBy)
	state.LastUpdatedBy = types.StringValue(environment.LastUpdatedBy)
	state.LastUpdatedAt = types.StringValue(environment.LastUpdatedAt)
	state.IsDeleted = types.BoolValue(environment.IsDeleted)
	state.PreviewContextId = types.StringValue(environment.PreviewContextId)
	state.LiveContextId = types.StringValue(environment.LiveContextId)
	state.HighAvailabilityEnabled = types.BoolValue(environment.HighAvailabilityEnabled)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success
func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan environmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state environmentResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the environment update
	environment := apiclient.Environment{
		ID:        plan.ID.ValueString(),
		Name:      plan.Name.ValueString(),
		ProjectID: plan.ProjectID.ValueString(),
	}

	// Set tenant type if provided
	if !plan.TenantType.IsNull() && !plan.TenantType.IsUnknown() {
		environment.TenantType = plan.TenantType.ValueString()
	}

	err := r.client.UpdateEnvironment(plan.ProjectID.ValueString(), plan.ID.ValueString(), environment)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating environment",
			"Could not update environment, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated environment from API
	updatedEnvironment, err := r.client.GetEnvironment(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated environment",
			"Could not read updated environment ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update state with refreshed values
	plan.Name = types.StringValue(updatedEnvironment.Name)
	plan.Host = types.StringValue(updatedEnvironment.Host)
	plan.PlatformTenantId = types.StringValue(updatedEnvironment.PlatformTenantId)
	plan.PlatformTenantName = types.StringValue(updatedEnvironment.PlatformTenantName)
	plan.TenantType = types.StringValue(updatedEnvironment.TenantType)
	plan.CreatedAt = types.StringValue(updatedEnvironment.CreatedAt)
	plan.CreatedBy = types.StringValue(updatedEnvironment.CreatedBy)
	plan.LastUpdatedBy = types.StringValue(updatedEnvironment.LastUpdatedBy)
	plan.LastUpdatedAt = types.StringValue(updatedEnvironment.LastUpdatedAt)
	plan.IsDeleted = types.BoolValue(updatedEnvironment.IsDeleted)
	plan.PreviewContextId = types.StringValue(updatedEnvironment.PreviewContextId)
	plan.LiveContextId = types.StringValue(updatedEnvironment.LiveContextId)
	plan.HighAvailabilityEnabled = types.BoolValue(updatedEnvironment.HighAvailabilityEnabled)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success
func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state environmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the environment
	err := r.client.DeleteEnvironment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting environment",
			"Could not delete environment, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing environment into Terraform state
func (r *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// Expected format: project_id,environment_id
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
