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
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreai/pkg/apiclient"
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the Edge client",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID for the Edge client",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_id": schema.StringAttribute{
				Description: "The environment ID for the Edge client",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID for authentication",
				Computed:    true,
				Sensitive:   false,
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

	// We need to get all organization clients to find the newly created client and get its ID
	// as it is not returned from api and is needed for future calls.
	// Feature request XS-11108 to include id in api response
	clientsResponse, err := r.client.GetClientsForEnvironment()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving organization clients",
			"Could not retrieve organization clients to find newly created client: "+err.Error(),
		)
		return
	}

	// Find the newly created client by matching client IDs
	var clientID string
	var resourceID string
	for _, client := range clientsResponse.Items {
		if client.ClientID == clientResponse.ClientID {
			clientID = client.ClientID
			resourceID = client.ID
			break
		}
	}

	if resourceID == "" {
		resp.Diagnostics.AddError(
			"Error finding created client",
			"Could not find the newly created client in the organization clients list",
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(resourceID) // Using the resource ID from GetClientsForOrganization
	plan.ClientID = types.StringValue(clientID)
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

	// Get all organization clients to find the client by its resource ID
	clientsResponse, err := r.client.GetClientsForEnvironment()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving organization clients",
			"Could not retrieve organization clients: "+err.Error(),
		)
		return
	}

	// Find the client by matching resource IDs
	var foundClient *apiclient.ClientDto
	for _, client := range clientsResponse.Items {
		if client.ID == state.ID.ValueString() {
			foundClient = &client
			break
		}
	}

	if foundClient == nil {
		resp.Diagnostics.AddError(
			"Client not found",
			"Could not find Edge client with ID: "+state.ID.ValueString(),
		)
		return
	}

	// Update state with current values from API
	state.ID = types.StringValue(foundClient.ID)
	state.Name = types.StringValue(foundClient.Name)
	if foundClient.Description != "" {
		state.Description = types.StringValue(foundClient.Description)
	}
	state.ClientID = types.StringValue(foundClient.ClientID)

	// Note: We can't retrieve the client secret after creation, so we keep the existing value

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

	// Get current state to access the stored client ID
	var state edgeClientResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For Edge clients, we need to delete the old client and create a new one
	// since the API doesn't support updating clients

	// Delete the old client
	err := r.client.DeleteClient(state.ID.ValueString())
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

	// Get the new resource ID from GetClientsForOrganization
	newClientsResponse, err := r.client.GetClientsForEnvironment()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving organization clients",
			"Could not retrieve organization clients to find newly created client: "+err.Error(),
		)
		return
	}

	// Find the newly created client by matching client IDs
	var newResourceID string
	for _, client := range newClientsResponse.Items {
		if client.ClientID == clientResponse.ClientID {
			newResourceID = client.ID
			break
		}
	}

	if newResourceID == "" {
		resp.Diagnostics.AddError(
			"Error finding created client",
			"Could not find the newly created client in the organization clients list",
		)
		return
	}

	// Update state with new values
	plan.ID = types.StringValue(newResourceID)
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

	// Delete the Edge client using the stored resource ID
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
