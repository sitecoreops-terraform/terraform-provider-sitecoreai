# SitecoreAI Provider

The SitecoreAI provider allows you to interact with SitecoreAI resources using Terraform or OpenTofu. This provider enables you to manage SitecoreAI projects, environments, and automation clients through infrastructure-as-code.

## Example Usage

```hcl
# Configure the SitecoreAI provider
provider "sitecore" {
  # Authentication can be provided via environment variables or directly in the provider configuration
  # client_id     = "your-client-id"
  # client_secret = "your-client-secret"
}

# You can also use environment variables for authentication:
# export SITECORE_CLIENT_ID="your-client-id"
# export SITECORE_CLIENT_SECRET="your-client-secret"
```

## Authentication

The SitecoreAI provider requires authentication to interact with the SitecoreAI Deploy API. You have two options for providing credentials:

1. **Environment Variables (Recommended)**:
   ```bash
   export SITECORE_CLIENT_ID="your-client-id"
   export SITECORE_CLIENT_SECRET="your-client-secret"
   ```

2. **Provider Configuration**:
   ```hcl
   provider "sitecore" {
     client_id     = "your-client-id"
     client_secret = "your-client-secret"
   }
   ```

To obtain credentials, follow the [SitecoreAI documentation](https://doc.sitecore.com/sai/en/developers/sitecoreai/manage-client-credentials-for-a-sitecoreai-organization-or-environment.html#create-an-automation-client-for-a-sitecoreai-organization) to create an Organization automation client.

## Argument Reference

The following arguments are supported in the provider configuration:

* `client_id` - (Optional) The client ID for Sitecore API authentication. If not provided, the `SITECORE_CLIENT_ID` environment variable will be used.
* `client_secret` - (Optional) The client secret for Sitecore API authentication. If not provided, the `SITECORE_CLIENT_SECRET` environment variable will be used.

!> **Warning**: Hardcoding credentials in the provider configuration is not recommended for production use. Use environment variables or a secure secrets management system instead.

## Proxy Configuration

If you need to use a proxy for API requests (e.g., for debugging or corporate networks), you can set the `HTTPS_PROXY` environment variable:

```bash
export HTTPS_PROXY="http://your-proxy-server:port"
```

## Getting Started

To get started with the SitecoreAI provider:

1. Install the provider (see installation instructions for your Terraform/OpenTofu setup)
2. Configure authentication using one of the methods above
3. Start managing SitecoreAI resources using the available resources and data sources

For more detailed examples, see the [examples directory](https://github.com/sitecoreops/terraform-provider-sitecoreai/tree/main/examples) in the provider repository.

## API Documentation

This provider interacts with the SitecoreAI Deploy API:

* [SitecoreAI Deploy API: Swagger UI](https://xmclouddeploy-api.sitecorecloud.io/)
* [SitecoreAI Deploy API: OpenAPI Specification](https://xmclouddeploy-api.sitecorecloud.io/swagger/v1/swagger.json)