data "sitecoreai_project" "default" {
  name = "XMC"
}

data "sitecoreai_environment" "default" {
  project_id = data.sitecoreai_project.default.id
  name       = "dev"
}

resource "sitecoreai_environment_variable" "env" {
  environment_id = data.sitecoreai_environment.default.id

  name  = "SXA_ENVIRONMENT_NAME"
  value = "dev"
}
