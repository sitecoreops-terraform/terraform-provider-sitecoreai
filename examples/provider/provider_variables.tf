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
  client_id     = var.sitecore_client_id
  client_secret = var.sitecore_client_secret
}

# Variables for provider configuration
variable "sitecore_client_id" {
  description = "The client ID for SitecoreAI Deploy API authentication"
  type        = string
  sensitive   = true
}

variable "sitecore_client_secret" {
  description = "The client secret for SitecoreAI Deploy API authentication"
  type        = string
  sensitive   = true
}