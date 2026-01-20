// Editing Secret data source implementation
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreai/pkg/apiclient"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &editingSecretDataSource{}
	_ datasource.DataSourceWithConfigure = &editingSecretDataSource{}
)

// NewEditingSecretDataSource is a helper function to simplify the provider implementation
func NewEditingSecretDataSource() datasource.DataSource {
	return &editingSecretDataSource{}
}

// editingSecretDataSource is the data source implementation
type editingSecretDataSource struct {
	client *apiclient.Client
}

// editingSecretDataSourceModel maps the data source schema data
type editingSecretDataSourceModel struct {
	EnvironmentID types.String `tfsdk:"environment_id"`
	Secret        types.String `tfsdk:"secret"`
}

// Metadata returns the data source type name
func (d *editingSecretDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_editing_secret"
}

// Schema defines the schema for the data source
func (d *editingSecretDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Environments ¤ Use this data source to get the editing secret for a Sitecore environment",
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Description: "The ID of the environment to get the editing secret for",
				Required:    true,
			},
			"secret": schema.StringAttribute{
				Description: "The editing secret for the environment",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *editingSecretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*apiclient.Client)
}

// Read refreshes the Terraform state with the latest data
func (d *editingSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get configuration
	var state editingSecretDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get editing secret from API
	secret, err := d.client.ObtainEditingSecret(state.EnvironmentID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading editing secret",
			"Could not read editing secret for environment "+state.EnvironmentID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Handle empty secret case
	if secret == "" {
		resp.Diagnostics.AddWarning(
			"Empty editing secret",
			"The editing secret for environment "+state.EnvironmentID.ValueString()+" is empty, this might be the case until there have been a deployment to the environment",
		)
		return
	}

	// Map the secret to the schema
	state.Secret = types.StringValue(secret)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
