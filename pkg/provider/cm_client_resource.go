// CM Client resource implementation
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
	_ resource.Resource                = &cmClientResource{}
	_ resource.ResourceWithConfigure   = &cmClientResource{}
	_ resource.ResourceWithImportState = &cmClientResource{}
)

// NewCMClientResource is a helper function to simplify the provider implementation
func NewCMClientResource() resource.Resource {
	return &cmClientResource{}
}

// cmClientResource is the resource implementation
type cmClientResource struct {
	client *apiclient.Client
}

// cmClientResourceModel maps the resource schema data
type cmClientResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ProjectID     types.String `tfsdk:"project_id"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	ClientID      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
}

// Metadata returns the resource type name
func (r *cmClientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_client"
}

// Schema defines the schema for the resource
func (r *cmClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Automation Clients ¤ Manages a Sitecore CM automation client",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the CM client",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the CM client",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the CM client",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID for the CM client",
				Required:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "The environment ID for the CM client",
				Required:    true,
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
func (r *cmClientResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*apiclient.Client)
}

// Create creates the resource and sets the initial Terraform state
func (r *cmClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cmClientResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the CM client
	clientResponse, err := r.client.CreateCMClient(
		plan.ProjectID.ValueString(),
		plan.EnvironmentID.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating CM client",
			"Could not create CM client, unexpected error: "+err.Error(),
		)
		return
	}

	// We need to get all environment clients to find the newly created client and get its ID
	// as it is not returned from api and is needed for future calls.
	// Feature request XS-11108 to include id in api response
	clientsResponse, err := r.client.GetClientsForEnvironment()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving environment clients",
			"Could not retrieve environment clients to find newly created client: "+err.Error(),
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
			"Could not find the newly created client in the environment clients list",
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(resourceID) // Using the resource ID from GetClientsForEnvironment
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
func (r *cmClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cmClientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all environment clients to find the client by its resource ID
	clientsResponse, err := r.client.GetClientsForEnvironment()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving environment clients",
			"Could not retrieve environment clients: "+err.Error(),
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
			"Could not find CM client with ID: "+state.ID.ValueString(),
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
func (r *cmClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cmClientResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For CM clients, we need to delete the old client and create a new one
	// since the API doesn't support updating clients

	// Delete the old client
	err := r.client.DeleteClient(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating CM client",
			"Could not delete old CM client: "+err.Error(),
		)
		return
	}

	// Then create a new client with the updated values
	clientResponse, err := r.client.CreateCMClient(
		plan.ProjectID.ValueString(),
		plan.EnvironmentID.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating CM client",
			"Could not create new CM client: "+err.Error(),
		)
		return
	}

	// Get the new resource ID from GetClientsForEnvironment
	newClientsResponse, err := r.client.GetClientsForEnvironment()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving environment clients",
			"Could not retrieve environment clients to find newly created client: "+err.Error(),
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
			"Could not find the newly created client in the environment clients list",
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
func (r *cmClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cmClientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the CM client
	err := r.client.DeleteClient(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting CM client",
			"Could not delete CM client, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing CM client into Terraform state
func (r *cmClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
