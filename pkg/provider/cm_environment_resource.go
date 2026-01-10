// CM Environment resource implementation
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops/terraform-provider-sitecoreai/pkg/apiclient"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &cmEnvironmentResource{}
	_ resource.ResourceWithConfigure   = &cmEnvironmentResource{}
	_ resource.ResourceWithImportState = &cmEnvironmentResource{}
)

// NewCMEnvironmentResource is a helper function to simplify the provider implementation
func NewCMEnvironmentResource() resource.Resource {
	return &cmEnvironmentResource{}
}

// cmEnvironmentResource is the resource implementation
type cmEnvironmentResource struct {
	client *apiclient.Client
}

// cmEnvironmentResourceModel maps the resource schema data
type cmEnvironmentResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	ProjectID               types.String `tfsdk:"project_id"`
	IsProd                  types.Bool   `tfsdk:"is_prod"`
	Host                    types.String `tfsdk:"host"`
	PlatformTenantId        types.String `tfsdk:"platform_tenant_id"`
	PlatformTenantName      types.String `tfsdk:"platform_tenant_name"`
	TenantType              types.String `tfsdk:"tenant_type"`
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
func (r *cmEnvironmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_environment"
}

// Schema defines the schema for the resource
func (r *cmEnvironmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Sitecore CM-only environment",
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
			"tenant_type": schema.StringAttribute{
				Description: "The tenant type for the environment",
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
func (r *cmEnvironmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*apiclient.Client)
}

// Create creates the resource and sets the initial Terraform state
func (r *cmEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cmEnvironmentResourceModel
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

	// Call API with CM environment type
	createdEnvironment, err := r.client.CreateEnvironment(
		plan.ProjectID.ValueString(),
		plan.Name.ValueString(),
		isProd,
		apiclient.EnvironmentTypeCmOnly,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating CM environment",
			"Could not create CM environment, unexpected error: "+err.Error(),
		)
		return
	}

	// Wait for environment to be ready with context IDs (timeout after 30 minutes)
	readyEnvironment, err := r.client.WaitForEnvironmentReady(
		createdEnvironment.ID,
		10, // 10 minutes timeout
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for CM environment to be ready",
			"Could not wait for CM environment to be ready: "+err.Error(),
		)
		return
	}

	// Use the ready environment instead of the initially created one
	createdEnvironment = readyEnvironment

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
func (r *cmEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cmEnvironmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get environment from API
	environment, err := r.client.GetEnvironment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading CM environment",
			"Could not read CM environment ID "+state.ID.ValueString()+": "+err.Error(),
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
func (r *cmEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cmEnvironmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state cmEnvironmentResourceModel
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

	err := r.client.UpdateEnvironment(plan.ProjectID.ValueString(), plan.ID.ValueString(), environment)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating CM environment",
			"Could not update CM environment, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated environment from API
	updatedEnvironment, err := r.client.GetEnvironment(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated CM environment",
			"Could not read updated CM environment ID "+plan.ID.ValueString()+": "+err.Error(),
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
func (r *cmEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cmEnvironmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the environment
	err := r.client.DeleteEnvironment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting CM environment",
			"Could not delete CM environment, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing CM environment into Terraform state
func (r *cmEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// Expected format: project_id,environment_id
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
