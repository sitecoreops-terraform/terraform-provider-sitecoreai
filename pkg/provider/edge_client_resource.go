// Edge Client resource implementation
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
	_ resource.Resource                = &edgeClientResource{}
	_ resource.ResourceWithConfigure   = &edgeClientResource{}
	_ resource.ResourceWithImportState = &edgeClientResource{}
)

// NewEdgeClientResource is a helper function to simplify the provider implementation
func NewEdgeClientResource() resource.Resource {
	return &edgeClientResource{}
}

// edgeClientResource is the resource implementation
type edgeClientResource struct {
	client *apiclient.Client
}

// edgeClientResourceModel maps the resource schema data
type edgeClientResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ProjectID     types.String `tfsdk:"project_id"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	ClientID      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
}

// Metadata returns the resource type name
func (r *edgeClientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_edge_client"
}

// Schema defines the schema for the resource
func (r *edgeClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Sitecore Edge automation client",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the Edge client",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Edge client",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the Edge client",
				Optional:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID for the Edge client",
				Required:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "The environment ID for the Edge client",
				Required:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID for authentication",
				Computed:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The client secret for authentication",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *edgeClientResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*apiclient.Client)
}

// Create creates the resource and sets the initial Terraform state
func (r *edgeClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan edgeClientResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the Edge client
	clientResponse, err := r.client.CreateEdgeClient(
		plan.ProjectID.ValueString(),
		plan.EnvironmentID.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Edge client",
			"Could not create Edge client, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(clientResponse.ClientID) // Using client ID as the resource ID
	plan.ClientID = types.StringValue(clientResponse.ClientID)
	plan.ClientSecret = types.StringValue(clientResponse.ClientSecret)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data
func (r *edgeClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state edgeClientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For Edge clients, we can't currently retrieve the full client details from the API
	// So we'll just return the current state as-is

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success
func (r *edgeClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan edgeClientResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For Edge clients, we need to delete the old client and create a new one
	// since the API doesn't support updating clients

	// First, delete the old client
	err := r.client.DeleteClient(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Edge client",
			"Could not delete old Edge client: "+err.Error(),
		)
		return
	}

	// Then create a new client with the updated values
	clientResponse, err := r.client.CreateEdgeClient(
		plan.ProjectID.ValueString(),
		plan.EnvironmentID.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Edge client",
			"Could not create new Edge client: "+err.Error(),
		)
		return
	}

	// Update state with new values
	plan.ID = types.StringValue(clientResponse.ClientID)
	plan.ClientID = types.StringValue(clientResponse.ClientID)
	plan.ClientSecret = types.StringValue(clientResponse.ClientSecret)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success
func (r *edgeClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state edgeClientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the Edge client
	err := r.client.DeleteClient(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Edge client",
			"Could not delete Edge client, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing Edge client into Terraform state
func (r *edgeClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
