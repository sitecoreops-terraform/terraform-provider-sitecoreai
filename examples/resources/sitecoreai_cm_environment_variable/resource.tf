data "sitecoreai_project" "default" {
  name = "XMC"
}

data "sitecoreai_environment" "default" {
  project_id = data.sitecoreai_project.default.id
  name       = "dev"
}

# Example: CM environment variable (target="CM" is implicit)
resource "sitecore_cm_environment_variable" "non_secret" {
  environment_id = sitecoreai_environment.default.id
  name           = "Sitecore_GraphQL_ExposePlayground"
  value          = "true" # Non-sensitive value
}

# Example: Secret CM environment variable
resource "sitecore_cm_environment_variable" "secret" {
  environment_id = sitecoreai_environment.default.id
  name           = "CM_SECRET_VAR"
  secret_value   = "s3cr3t_p@ssw0rd" # Sensitive value (masked in logs/state)
}
