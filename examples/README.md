# Sitecore AI Terraform Provider Examples

This directory contains comprehensive examples demonstrating how to use the Sitecore AI Terraform provider.

The examples are mostly to be included in documentation, but there are also a few full examples on their own.

## Documentation generation

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

| Path                                                                         | Description                                |
|------------------------------------------------------------------------------|--------------------------------------------|
| `examples/provider/provider<*>.tf`                                           | Provider example config(s)                 |
| `examples/data-sources/<data source name>/data-source<*>.tf`                 | Data source example config(s)              |
| `examples/resources/<resource name>/resource<*>.tf`                          | Resource example config(s)                 |
| `examples/resources/<resource name>/import.sh`                               | Resource example import command            |
| `examples/resources/<resource name>/import-by-string-id.tf`                  | Resource example import by id config       |
| `examples/resources/<resource name>/import-by-identity.tf`                   | Resource example import by identity config |

## Prerequisites

- Terraform 1.0+
- SitecoreAI API credentials (create organization credentials in SitecoreAI Deploy)

## Working with local provider implementation

Usually terraform providers are found through the registry, however you can specify a local override where a certain provider is found. This can be specified in a `.terraformrc` file in the user's home directory, or the environment variable `TF_CLI_CONFIG_FILE` can point to a `*.tfrc` file. We have created a `localdev.tfrc` and hereby you can run

```bash
go build
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



## Complete examples on their own

The following examples are intended to be used in their own and will not be included in documentation.

### [1. Reading](1.reading/)

This is a very basic example, only reading from SitecoreAI but still showing how to:
- Configure provider credentials
- Look up existing project using data sources
- Use outputs to display project information

This example is safe to run even with your existing environment, as it will not make any changes, even when running `terraform apply`

### [2. Create environment](2.environment/)

This is a simple example that creates a new environment for an existing project.

### [3. Automation clients](3.automation_clients/)

This is an example to create automation clients to use for other integrations (eg. Azure DevOps pipelines etc.)

Be aware that an organization client doesn't have access to create automation clients, so usually CLI authentication is required
