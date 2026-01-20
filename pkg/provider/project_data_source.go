// Project data source implementation
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreai/pkg/apiclient"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

// NewProjectDataSource is a helper function to simplify the provider implementation
func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

// projectDataSource is the data source implementation
type projectDataSource struct {
	client *apiclient.Client
}

// projectDataSourceModel maps the data source schema data
type projectDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the data source type name
func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the data source
func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Environments ¤ Use this data source to get information about a Sitecore project by name",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the project",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the project to search for",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the project",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*apiclient.Client)
}

// Read refreshes the Terraform state with the latest data
func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state projectDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all projects from API
	projects, err := d.client.GetProjects()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading projects",
			"Could not read projects: "+err.Error(),
		)
		return
	}

	// Find the project with the matching name
	var foundProject *apiclient.Project
	for i := range projects {
		if projects[i].Name == state.Name.ValueString() {
			foundProject = &projects[i]
			break
		}
	}

	if foundProject == nil {
		resp.Diagnostics.AddError(
			"Project not found",
			fmt.Sprintf("Could not find project with name: %s", state.Name.ValueString()),
		)
		return
	}

	// Map the project data to the schema
	state.ID = types.StringValue(foundProject.ID)
	state.Name = types.StringValue(foundProject.Name)
	state.Description = types.StringValue(foundProject.Description)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
