terraform {
  required_providers {
    sitecoreai = {
      source = "sitecoreops-terraform/sitecoreai"
    }
  }
  required_version = ">= 0.1.0"
}

# Configure the Sitecore AI provider
provider "sitecoreai" {
  # Authenticate with CLI before running terraform
  # The .sitecore folder must be in terraform folder or a parent folder
  use_cli = true
}
