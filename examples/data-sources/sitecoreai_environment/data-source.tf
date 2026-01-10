data "sitecoreai_project" "existing" {
  name = "XMC" # Replace with an existing project name
}

# Use the environment data source to get information about an environment
data "sitecoreai_environment" "existing" {
  project_id = data.sitecoreai_project.existing.id
  name       = "production" # Replace with existing environment name
}

# Output the environment information
output "environment_id" {
  value       = data.sitecoreai_environment.existing.id
  description = "The ID of the existing environment"
}

output "environment_host" {
  value       = data.sitecoreai_environment.existing.host
  description = "The host URL of the existing environment"
}

output "environment_is_prod" {
  value       = data.sitecoreai_environment.existing.is_prod
  description = "Whether the environment is a production environment"
}