data "sitecore_project" "main" {
  name = var.project_name
}

resource "sitecore_cm_environment" "main" {
  project_id = data.sitecore_project.main.id
  name       = var.environment_name
  is_prod    = var.environment_is_prod
}
