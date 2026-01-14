// Deploy Client resource implementation
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
	_ resource.Resource                = &deployClientResource{}
	_ resource.ResourceWithConfigure   = &deployClientResource{}
	_ resource.ResourceWithImportState = &deployClientResource{}
)

// NewDeployClientResource is a helper function to simplify the provider implementation
func NewDeployClientResource() resource.Resource {
	return &deployClientResource{}
}

// deployClientResource is the resource implementation
type deployClientResource struct {
	client *apiclient.Client
}

// deployClientResourceModel maps the resource schema data
type deployClientResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

// Metadata returns the resource type name
func (r *deployClientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deploy_client"
}

// Schema defines the schema for the resource
func (r *deployClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Sitecore Deploy automation client",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the Deploy client",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Deploy client",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the Deploy client",
				Optional:    true,
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
func (r *deployClientResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*apiclient.Client)
}

// Create creates the resource and sets the initial Terraform state
func (r *deployClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan deployClientResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the Deploy client
	clientResponse, err := r.client.CreateDeployClient(
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Deploy client",
			"Could not create Deploy client, unexpected error: "+err.Error(),
		)
		return
	}

	// We need to get all clients to find the newly created client and get its ID
	// as it is not returned from api and is needed for future calls.
	// Feature request XS-11108 to include id in api response
	clientsResponse, err := r.client.GetClientsForOrganization()
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
func (r *deployClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state deployClientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all organization clients to find the client by its resource ID
	clientsResponse, err := r.client.GetClientsForOrganization()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving organization clients",
			"Could not retrieve organization clients: "+err.Error(),
		)
		return
	}

	// Find the client by matching resource IDs
	var foundClient *apiclient.OrganizationClientDto
	for _, client := range clientsResponse.Items {
		if client.ID == state.ID.ValueString() {
			foundClient = &client
			break
		}
	}

	if foundClient == nil {
		resp.Diagnostics.AddError(
			"Client not found",
			"Could not find Deploy client with ID: "+state.ID.ValueString(),
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
func (r *deployClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan deployClientResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to access the stored client ID
	var state deployClientResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For Deploy clients, we need to delete the old client and create a new one
	// since the API doesn't support updating clients

	// Delete the old client
	err := r.client.DeleteClient(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Deploy client",
			"Could not delete old Deploy client: "+err.Error(),
		)
		return
	}

	// Then create a new client with the updated values
	clientResponse, err := r.client.CreateDeployClient(
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Deploy client",
			"Could not create new Deploy client: "+err.Error(),
		)
		return
	}

	// Get the new resource ID from GetClientsForOrganization
	newClientsResponse, err := r.client.GetClientsForOrganization()
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
func (r *deployClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state deployClientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the Deploy client using the stored resource ID
	err := r.client.DeleteClient(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Deploy client",
			"Could not delete Deploy client, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing Deploy client into Terraform state
func (r *deployClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
