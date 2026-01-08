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