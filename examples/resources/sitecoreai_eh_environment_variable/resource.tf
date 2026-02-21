data "sitecoreai_project" "default" {
  name = "XMC" # Replace with an existing project name
}

# Get the CM environment to associate with the EH environment
data "sitecoreai_environment" "default" {
  project_id = data.sitecoreai_project.default.id
  name       = "dev" # Replace with existing CM environment name
}

# Create editing host environment
resource "sitecoreai_eh_environment" "website" {
  name              = "website"
  project_id        = data.sitecoreai_project.default.id
  cm_environment_id = data.sitecoreai_environment.default.id
}

# Example: EH environment variable (target="EH" is implicit)
resource "sitecore_eh_environment_variable" "non_secret" {
  environment_id = sitecoreai_eh_environment.website.id
  name           = "EH_NON_SECRET_VAR"
  value          = "hello_world" # Non-sensitive value
}

# Example: Secret EH environment variable
resource "sitecore_eh_environment_variable" "secret" {
  environment_id = sitecoreai_eh_environment.website.id
  name           = "EH_SECRET_VAR"
  secret_value   = "s3cr3t_p@ssw0rd" # Sensitive value (masked in logs/state)
}
