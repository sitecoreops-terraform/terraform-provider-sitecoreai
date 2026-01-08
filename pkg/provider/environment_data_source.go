// Environment data source implementation
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops/terraform-provider-sitecoreai/pkg/apiclient"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &environmentDataSource{}
	_ datasource.DataSourceWithConfigure = &environmentDataSource{}
)

// NewEnvironmentDataSource is a helper function
func NewEnvironmentDataSource() datasource.DataSource {
	return &environmentDataSource{}
}

// environmentDataSource is the data source implementation
type environmentDataSource struct {
	client *apiclient.Client
}

// environmentDataSourceModel maps the data source schema data
type environmentDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	ProjectID   types.String `tfsdk:"project_id"`
	Description types.String `tfsdk:"description"`
	Variables   types.Map    `tfsdk:"variables"`
}

// Metadata returns the data source type name
func (d *environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// Schema defines the schema for the data source
func (d *environmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Sitecore environment by project ID and name",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the environment",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the environment to search for",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The ID of the project to search within",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the environment",
				Computed:    true,
			},
			"variables": schema.MapAttribute{
				Description: "Environment variables as key-value pairs",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *environmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*apiclient.Client)
}

// Read refreshes the Terraform state with the latest data
func (d *environmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get configuration
	var state environmentDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all environments for the specified project
	environments, err := d.client.GetProjectEnvironments(state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading environments",
			"Could not read environments for project "+state.ProjectID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Find the environment with the matching name
	var foundEnvironment *apiclient.Environment
	for i := range environments {
		if environments[i].Name == state.Name.ValueString() {
			foundEnvironment = &environments[i]
			break
		}
	}

	if foundEnvironment == nil {
		resp.Diagnostics.AddError(
			"Environment not found",
			fmt.Sprintf("Could not find environment with name '%s' in project '%s'",
				state.Name.ValueString(), state.ProjectID.ValueString()),
		)
		return
	}

	// Map the environment data to the schema
	state.ID = types.StringValue(foundEnvironment.ID)
	state.Name = types.StringValue(foundEnvironment.Name)
	state.ProjectID = types.StringValue(foundEnvironment.ProjectID)
	state.Description = types.StringValue(foundEnvironment.Description)

	// Convert variables map to Terraform types
	if len(foundEnvironment.Variables) > 0 {
		variablesAttr := make(map[string]attr.Value)
		for key, value := range foundEnvironment.Variables {
			variablesAttr[key] = types.StringValue(value)
		}
		state.Variables = types.MapValueMust(types.StringType, variablesAttr)
	} else {
		state.Variables = types.MapNull(types.StringType)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
