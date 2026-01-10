# Sitecore AI Terraform Provider Examples

This directory contains comprehensive examples demonstrating how to use the Sitecore AI Terraform provider.

## Directory Structure

```
examples/
├── complete_example/          # Complete end-to-end example
├── data-sources/              # Data source examples
│   ├── sitecoreai_editing_secret/
│   ├── sitecoreai_environment/
│   └── sitecoreai_project/
├── resources/                 # Resource examples
│   ├── sitecoreai_cm_client/
│   ├── sitecoreai_cm_environment/
│   ├── sitecoreai_deploy_client/
│   ├── sitecoreai_edge_client/
│   ├── sitecoreai_editing_host_build_client/
│   ├── sitecoreai_environment/
│   └── sitecoreai_project/
└── README.md                  # This file
```

## Documentation generation

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** example file for the provider index page
* **data-sources/`full data source name`/example.tf** example file for the named data source page
* **resources/`full resource name`/example.tf** example file for the named resource page

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

## Usage

To use any of these examples:

1. **Navigate to the example directory**:
   ```bash
   cd examples/resources/sitecoreai_project
   ```

2. **Initialize Terraform**:
   ```bash
   terraform init
   ```

3. **Review the plan**:
   ```bash
   terraform plan
   ```

4. **Apply the configuration**:
   ```bash
   terraform apply
   ```

5. **Clean up** (when done):
   ```bash
   terraform destroy
   ```

## Best Practices

1. **Use variables** for customizable values
2. **Mark sensitive outputs** as `sensitive = true`
3. **Use data sources** to reference existing resources
4. **Organize resources** logically (projects → environments → clients)
5. **Use separate environments** for development, staging, and production
6. **Store state remotely** for team collaboration

#### Option A: Use variables

```hcl
terraform {
  required_providers {
    sitecore = {
      source = "sitecoreops/sitecoreai"
    }
  }
}

provider "sitecore" {
  client_id     = var.sitecore_client_id
  client_secret = var.sitecore_client_secret
}

variable "sitecore_client_id" {
  description = "The client ID for authentication"
  type        = string
  sensitive   = true
}

variable "sitecore_client_secret" {
  description = "The client secret for authentication"
  type        = string
  sensitive   = true
}
```

There are several ways to specify the variables, all standard terraform functionality:

* Specify the variables from the command line when running commands

    ```bash
    terraform apply -var="sitecore_client_id=your_client_id" -var="sitecore_client_secret=your_client_secret"
    ```

* Create a `.tfvars` file with your variables and use that when running

    ```hcl
    sitecore_client_id     = "y4AKvexbrg7WVSXnXfXyAgs63cu4Abt8"
    sitecore_client_secret = "a1mrZlL2bt94DplYUgrEqr2gu1iyyclRqRo3mnI6zeNDE6Gl95U4TnmyDzM7SoVq"
    ```

    ```bash
    terraform apply -var-file=".tfvars"
    ```

* Specify the variables using default terraform syntax for environment variables

    ```bash
    export TF_VAR_sitecore_client_id="your_client_id"
    export TF_VAR_sitecore_client_secret="your_client_secret"

    terraform apply
    ```

#### Option B: Provider Environment Variables

The provider implementation will use environment variables if those are specified:

```hcl
terraform {
  required_providers {
    sitecore = {
      source = "sitecoreops/sitecoreai"
    }
  }
}

provider "sitecore" {}
```

```bash
export SITECORE_CLIENT_ID="your_client_id"
export SITECORE_CLIENT_SECRET="your_client_secret"

terraform apply
```

## Examples

### Resources

#### [Project Resource](resources/sitecoreai_project/)
Demonstrates how to create and manage Sitecore projects.

#### [Environment Resource](resources/sitecoreai_environment/)
Shows how to create and manage Sitecore environments within projects.

#### [CM Environment Resource](resources/sitecoreai_cm_environment/)
Examples for managing Content Management environments.

#### [Client Resources](resources/)
Examples for creating various types of clients:
- **CM Client**: Content Management automation clients
- **Edge Client**: Edge delivery clients  
- **Deploy Client**: Deployment automation clients
- **Editing Host Build Client**: Editing host build clients

### Data Sources

#### [Project Data Source](data-sources/sitecoreai_project/)
How to retrieve information about existing projects.

#### [Environment Data Source](data-sources/sitecoreai_environment/)
Retrieving information about existing environments.

#### [Editing Secret Data Source](data-sources/sitecoreai_editing_secret/)
Getting editing secrets for environments.

### Complete Example

The [complete_example/](complete_example/) directory contains a comprehensive example that demonstrates:

- Creating a project
- Creating development and production environments
- Setting up all types of clients (CM, Edge, Deploy, Editing Host)
- Using data sources to reference existing resources
- Proper output configuration
- Variable usage for customization

### Existing Examples

The following examples were already present in the repository:

#### [1. Reading](1.reading/)

This is a very basic example, only reading from SitecoreAI but still showing how to:
- Configure provider credentials
- Look up existing project using data sources
- Use outputs to display project information

This example is safe to run even with your existing environment, as it will not make any changes, even when running `terraform apply`

#### [2. Create environment](2.environment/)

This is a simple example that creates a new environment for an existing project.

#### [3. Automation clients](3.automation_clients/)

This is an example to create automation clients to use for other integrations (eg. Azure DevOps pipelines etc.)
