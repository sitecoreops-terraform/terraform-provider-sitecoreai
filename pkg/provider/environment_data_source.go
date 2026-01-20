// Environment data source implementation
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
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	ProjectID               types.String `tfsdk:"project_id"`
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

// Metadata returns the data source type name
func (d *environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// Schema defines the schema for the data source
func (d *environmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Environments ¤ Use this data source to get information about a Sitecore environment by project ID and name",
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
				Description: "The tenant type",
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
	state.Host = types.StringValue(foundEnvironment.Host)
	state.PlatformTenantId = types.StringValue(foundEnvironment.PlatformTenantId)
	state.PlatformTenantName = types.StringValue(foundEnvironment.PlatformTenantName)
	state.TenantType = types.StringValue(foundEnvironment.TenantType)
	state.CreatedAt = types.StringValue(foundEnvironment.CreatedAt)
	state.CreatedBy = types.StringValue(foundEnvironment.CreatedBy)
	state.LastUpdatedBy = types.StringValue(foundEnvironment.LastUpdatedBy)
	state.LastUpdatedAt = types.StringValue(foundEnvironment.LastUpdatedAt)
	state.IsDeleted = types.BoolValue(foundEnvironment.IsDeleted)
	state.PreviewContextId = types.StringValue(foundEnvironment.PreviewContextId)
	state.LiveContextId = types.StringValue(foundEnvironment.LiveContextId)
	state.HighAvailabilityEnabled = types.BoolValue(foundEnvironment.HighAvailabilityEnabled)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
