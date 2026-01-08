data "sitecore_project" "main" {
  name = var.project_name
}

data "sitecore_environment" "main" {
  project_id = data.sitecore_project.main.id
  name       = var.environment_name
}
