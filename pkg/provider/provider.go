//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name sitecore

// Package provider contains the Sitecore Terraform provider implementation
package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreai/pkg/apiclient"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &sitecoreProvider{}
)

// New is a helper function to simplify provider server and testing implementation
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sitecoreProvider{
			version: version,
		}
	}
}

// sitecoreProvider is the provider implementation
type sitecoreProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// sitecoreProviderModel maps provider schema data to a Go type
type sitecoreProviderModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	UseCLI       types.Bool   `tfsdk:"use_cli"`
}

// Metadata returns the provider type name
func (p *sitecoreProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sitecoreai"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data
func (p *sitecoreProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with SitecoreAI",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "The client ID for Sitecore API authentication",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The client secret for Sitecore API authentication",
				Optional:    true,
				Sensitive:   true,
			},
			"use_cli": schema.BoolAttribute{
				Description: "Use Sitecore CLI authentication (searches for .sitecore/user.json)",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a Sitecore API client for data sources and resources
func (p *sitecoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config sitecoreProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var client *apiclient.Client
	var err error

	// Handle environment variables
	// Check if CLI authentication is requested
	useCLI := os.Getenv("SITECOREAI_USE_CLI") == "1"
	if !config.UseCLI.IsNull() {
		useCLI = config.UseCLI.ValueBool()
	}

	// Use traditional client_id/client_secret authentication
	clientID := os.Getenv("SITECOREAI_CLIENT_ID")
	clientSecret := os.Getenv("SITECOREAI_CLIENT_SECRET")

	// Override with configuration values if provided
	if !config.ClientID.IsNull() && len(config.ClientID.ValueString()) > 0 {
		clientID = config.ClientID.ValueString()
	}
	if !config.ClientSecret.IsNull() && len(config.ClientSecret.ValueString()) > 0 {
		clientSecret = config.ClientSecret.ValueString()
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value
	if config.ClientID.IsUnknown() || config.ClientSecret.IsUnknown() || config.UseCLI.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Sitecore API Configuration",
			"Cannot use unknown values for Sitecore API configuration",
		)
		return
	}

	if useCLI {
		// Try CLI authentication
		client, err = apiclient.NewClientFromCLI("")
		if err != nil {
			resp.Diagnostics.AddError(
				"Sitecore CLI Authentication Failed",
				"Unable to authenticate using Sitecore CLI: "+err.Error(),
			)
			return
		}
	} else {
		// Create a new Sitecore API client
		client, err = apiclient.NewClient(clientID, clientSecret)
		if err != nil {
			resp.Diagnostics.AddError(
				"SitecoreAI API Authentication Failed",
				"Unable to authenticate using Client Id and Client Secret: "+err.Error(),
			)
			return
		}
	}

	// Authenticate the client
	err = client.Authenticate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Sitecore API Client Authentication Failed",
			"Unable to authenticate Sitecore API client: "+err.Error(),
		)
		return
	}

	// Make the Sitecore API client available during data source and resource
	// type Configure methods
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider
func (p *sitecoreProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewEnvironmentDataSource,
		NewEditingSecretDataSource,
	}
}

// Resources defines the resources implemented in the provider
func (p *sitecoreProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewEnvironmentResource,
		NewCMEnvironmentResource,
		NewCMClientResource,
		NewEdgeClientResource,
		NewDeployClientResource,
		NewEditingHostBuildClientResource,
		NewEnvironmentVariableResource,
	}
}
