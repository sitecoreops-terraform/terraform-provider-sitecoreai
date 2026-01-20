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
  # Specify client credentials in environment variables
  # export SITECOREAI_CLIENT_ID="your-client-id"
  # export SITECOREAI_CLIENT_SECRET="your-client-secret"
  # or alternative,
  # export SITECOREAI_USE_CLI=1
}
