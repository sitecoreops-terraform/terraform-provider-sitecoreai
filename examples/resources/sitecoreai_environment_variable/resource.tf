data "sitecoreai_project" "default" {
  name = "XMC"
}

data "sitecoreai_environment" "default" {
  project_id = data.sitecoreai_project.default.id
  name       = "dev"
}

resource "sitecoreai_environment_variable" "shared" {
  environment_id = sitecoreai_environment.default.id
  name           = "SHARED_VAR"
  value          = "hello_world" # Non-sensitive value
  # No target specified (applies to all targets)
}

# Example: CM environment variable
resource "sitecoreai_environment_variable" "playground" {
  environment_id = ssitecoreai_environment.default.id
  target         = "CM" # CM target
  name           = "Sitecore_GraphQL_ExposePlayground"
  value          = "true" # Non-sensitive value
}

resource "sitecoreai_environment_variable" "cm_secret" {
  environment_id = sitecoreai_environment.default.id
  target         = "CM" # CM target
  name           = "CM_SECRET_VAR"
  secret_value   = "s3cr3t_p@ssw0rd" # Sensitive value (masked in logs/state)
}

# Example: EH environment variable
resource "sitecoreai_environment_variable" "eh_non_secret" {
  environment_id = sitecoreai_environment.default.id
  target         = "website" # Target editing host named website in xmcloud.deploy.json
  name           = "EH_NON_SECRET_VAR"
  value          = "hello_world" # Non-sensitive value
}
