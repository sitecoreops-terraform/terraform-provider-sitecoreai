--- 
page_title: "Getting Started with SitecoreAI Provider"
--- 

# Getting Started with SitecoreAI Provider

This guide will walk you through setting up and using the SitecoreAI Terraform provider.

## Prerequisites

Before you begin, ensure you have:

1. Terraform or OpenTofu installed (version 1.0 or later)
2. SitecoreAI organization access with appropriate permissions
3. SitecoreAI deploy organization client credentials. 

## Installation

### Using Terraform Registry (Recommended)

```hcl
terraform {
  required_providers {
    sitecore = {
      source = "sitecoreops/sitecoreai"
      version = ">= 1.0.0"  # Use the latest version
    }
  }
}

provider "sitecore" {
  # Configuration will be covered in the next section
}
```


## Authentication Setup

### Step 1: Create SitecoreAI Automation Client

1. Go to your Sitecore portal
1. Open SitecoreAI Deploy
1. Navigate to Credentials [Link to organization credentials in SitecoreAI deploy](https://deploy.sitecorecloud.io/credentials/organization)
3. Create a new client by clicking "Create credentials" button
4. Note down the Client ID and Client Secret

### Step 2: Configure Provider Authentication

**Option A: Environment Variables (Recommended)**

```bash
# Set environment variables
export SITECORE_CLIENT_ID="your-client-id"
export SITECORE_CLIENT_SECRET="your-client-secret"

# Verify they're set
echo $SITECORE_CLIENT_ID
echo $SITECORE_CLIENT_SECRET
```

**Option B: Provider Configuration**

```hcl
provider "sitecore" {
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
}
```

## Basic Usage Examples

### Example 1: Reading Existing Resources

```hcl
# Get information about an existing project
data "sitecore_project" "example" {
  name = "My Sitecore Project"
}

# Get information about an environment
data "sitecore_environment" "example" {
  project_id = data.sitecore_project.example.id
  name       = "production"
}

# Get editing secret for an environment
data "sitecore_editing_secret" "example" {
  environment_id = data.sitecore_environment.example.id
}

output "project_id" {
  value = data.sitecore_project.example.id
}

output "editing_secret" {
  value     = data.sitecore_editing_secret.example.secret
  sensitive = true
}
```

### Example 2: Creating a CM Environment

```hcl
# Get the project first
data "sitecore_project" "main" {
  name = "My Project"
}

# Create a new CM environment
resource "sitecore_cm_environment" "staging" {
  project_id = data.sitecore_project.main.id
  name       = "staging"
  is_prod    = false
}

# Get the editing secret for the new environment
data "sitecore_editing_secret" "staging" {
  environment_id = sitecore_cm_environment.staging.id
}
```

### Example 3: Managing Automation Clients

```hcl
# Create automation clients for different purposes
resource "sitecore_deploy_client" "ci_cd" {
  name        = "CI/CD Pipeline Client"
  description = "Used by our CI/CD pipeline for deployments"
}

resource "sitecore_cm_client" "content_authoring" {
  project_id     = data.sitecore_project.main.id
  environment_id = data.sitecore_environment.main.id
  name           = "Content Authoring Client"
  description    = "Used by content authors for CMS access"
}

resource "sitecore_edge_client" "cdn_purge" {
  project_id     = data.sitecore_project.main.id
  environment_id = data.sitecore_environment.main.id
  name           = "CDN Purge Client"
  description    = "Used for CDN cache purging"
}
```

## Common Patterns

### Working with Multiple Environments

```hcl
locals {
  environments = ["development", "staging", "production"]
}

resource "sitecore_cm_environment" "envs" {
  for_each = toset(local.environments)
  
  project_id = data.sitecore_project.main.id
  name       = each.key
  is_prod    = each.key == "production"
}
```

### Using Outputs Safely

```hcl
# Always mark sensitive outputs as sensitive
output "editing_secrets" {
  value = {
    for env_name, env in sitecore_cm_environment.envs : 
      env_name => data.sitecore_editing_secret.envs[env_name].secret
  }
  sensitive = true
}
```

### Debugging

To enable debug logging:

```bash
# Set the TF_LOG environment variable
export TF_LOG=DEBUG

# Run Terraform commands
tf plan
```

## Next Steps

* Explore the [examples directory](https://github.com/sitecoreops/terraform-provider-sitecoreai/tree/main/examples) for more complex scenarios
* Review the documentation for specific resources and data sources
* Learn about advanced patterns like using modules and workspaces
