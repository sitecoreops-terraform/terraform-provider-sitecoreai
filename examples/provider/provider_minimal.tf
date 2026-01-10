terraform {
  required_providers {
    sitecoreai = {
      source = "sitecoreops/sitecoreai"
    }
  }
  required_version = ">= 1.0.0"
}

# Configure the Sitecore AI provider
provider "sitecoreai" {
  # Specify client credentials in environment variables
  # export SITECORE_CLIENT_ID="your-client-id"
  # export SITECORE_CLIENT_SECRET="your-client-secret"
}
