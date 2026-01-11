data "sitecoreai_project" "main" {
  name = var.project_name
}

resource "sitecoreai_cm_environment" "main" {
  project_id = data.sitecoreai_project.main.id
  name       = var.environment_name
  is_prod    = var.environment_is_prod
}

data "sitecoreai_editing_secret" "main" {
  environment_id = sitecoreai_cm_environment.main.id
}
