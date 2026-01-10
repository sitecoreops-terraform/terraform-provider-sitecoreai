# First, get the project
data "sitecore_project" "main" {
  name = "XMC"
}

# Get an existing environment
data "sitecore_environment" "main" {
  project_id = data.sitecore_project.main.id
  name       = "production"
}

# Get the editing secret for the environment
data "sitecore_editing_secret" "main" {
  environment_id = data.sitecore_environment.main.id
}

# Output the secret (marked as sensitive)
output "editing_secret" {
  value     = data.sitecore_editing_secret.main.secret
  sensitive = true
}