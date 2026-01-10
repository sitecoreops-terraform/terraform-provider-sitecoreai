data "sitecore_project" "main" {
  name = var.project_name
}

data "sitecore_environment" "main" {
  project_id = data.sitecore_project.main.id
  name       = var.environment_name
}

resource "sitecore_deploy_client" "default" {
  name        = "Sample deploy client"
  description = "Created from Terraform example"
}

resource "sitecore_cm_client" "default" {
  project_id     = data.sitecore_environment.main.project_id
  environment_id = data.sitecore_environment.main.id

  name        = "Sample CM client"
  description = "Created from Terraform example"
}

resource "sitecore_edge_client" "default" {
  project_id     = data.sitecore_environment.main.project_id
  environment_id = data.sitecore_environment.main.id

  name        = "Sample Edge client"
  description = "Created from Terraform example"
}
