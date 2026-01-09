# Sitecore Terraform Provider Examples

This directory contains examples demonstrating how to use the Sitecore Terraform Provider.

## Prerequisites

- Terraform 1.0+
- SitecoreAI API credentials (create organization credentials in SitecoreAI Deploy)

## Working with local provider implementation

Usually terraform providers are found through the registry, however you can specify a local override where a certain provider is found. This can be specified in a `.terraformrc` file in the user's home directory, or the environment variable `TF_CLI_CONFIG_FILE` can point to a `*.tfrc` file. We have created a `localdev.tfrc` and hereby you can run

```bash
go build -o out/terraform-provider-sitecoreai
cd examples
export TF_CLI_CONFIG_FILE=$(pwd)/localdev.tfrc
```

However, when there is only the local provider, there is no need to run `terraform init`

## Configure credentials

You will need credentials created in SitecoreAI deploy

The provider requires two configuration parameters:

- `client_id`: Your Sitecore AI API client ID
- `client_secret`: Your Sitecore AI API client secret

#### Option A: Environment Variables

```bash
export SITECORE_CLIENT_ID="your_client_id"
export SITECORE_CLIENT_SECRET="your_client_secret"
```

#### Option B: Variables File

Create `auto.tfvars` with your actual credentials and variables:

```hcl
sitecore_client_id     = "y4AKvexbrg7WVSXnXfXyAgs63cu4Abt8"
sitecore_client_secret = "a1mrZlL2bt94DplYUgrEqr2gu1iyyclRqRo3mnI6zeNDE6Gl95U4TnmyDzM7SoVq"
```

#### Option C: Command Line

```bash
terraform apply -var="sitecore_client_id=your_client_id" -var="sitecore_client_secret=your_client_secret"
```

## Available Examples

### [1. Reading/](1.reading/)

This is a very basic example, only reading from SitecoreAI but still showing how to
- Configure provider credentials
- Look up existing project using data sources
- Use outputs to display project information

This example is safe to run even with your existing environment, as it will not make any changes, even when running `terraform apply`

### [2. Create environment](2.environment/)

This is a simple example that creates a new environment for an exising project.
