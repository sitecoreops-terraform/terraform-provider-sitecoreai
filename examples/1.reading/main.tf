data "sitecoreai_project" "main" {
  name = var.project_name
}

data "sitecoreai_environment" "main" {
  project_id = data.sitecoreai_project.main.id
  name       = var.environment_name
}

data "sitecoreai_editing_secret" "main" {
  environment_id = data.sitecoreai_environment.main.id
}
