# Use the project data source to get information about a project by name
data "sitecoreai_project" "default" {
  name = "XMC" # Replace with an existing project name
}

# Create a CM environment
resource "sitecoreai_cm_environment" "cm_production" {
  name       = "production"
  project_id = data.sitecoreai_project.default.id
  is_prod    = true
}

# Output CM environment details
output "cm_environment_id" {
  value       = sitecoreai_cm_environment.cm_production.id
  description = "The ID of the CM environment"
}

output "cm_environment_name" {
  value       = sitecoreai_cm_environment.cm_production.name
  description = "The name of the CM environment"
}
