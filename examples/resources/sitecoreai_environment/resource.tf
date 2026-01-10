# First, create a project (or use an existing one)
resource "sitecoreai_project" "example" {
  name        = "XMC"
  description = "Project for environment example"
}

# Create a new environment within the project
resource "sitecoreai_environment" "staging" {
  name       = "staging"
  project_id = sitecoreai_project.example.id
  is_prod    = false
}

# Output the environment details
output "environment_id" {
  value       = sitecoreai_environment.staging.id
  description = "The ID of the created environment"
}

output "environment_host" {
  value       = sitecoreai_environment.staging.host
  description = "The host URL of the environment"
}

output "environment_status" {
  value       = sitecoreai_environment.staging.is_deleted ? "deleted" : "active"
  description = "The status of the environment"
}