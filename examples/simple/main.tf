terraform {
  required_providers {
    sitecore = {
      source = "github.com/jballe/sitecore-go-apiclient/terraform-provider-sitecore"
    }
  }
}

provider "sitecore" {
  client_id = var.sitecore_client_id
  client_secret = var.sitecore_client_secret
}

data "sitecore_project" "default" {
  name = "XMC"
}
