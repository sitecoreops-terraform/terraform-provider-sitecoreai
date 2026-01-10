# Use the project data source to get information about a project by name
data "sitecoreai_project" "default" {
  name = "XMC" # Replace with an existing project name
}

# Use the environment data source to get information about an environment
data "sitecoreai_environment" "default" {
  project_id = data.sitecoreai_project.default.id
  name       = "production" # Replace with existing environment name
}

# Create a CM automation client
resource "sitecoreai_cm_client" "automation" {
  name           = "terraform-cm-automation"
  description    = "CM automation client created via Terraform"
  project_id     = data.sitecoreai_environment.default.project_id
  environment_id = data.sitecoreai_environment.default.id
}

# Output client credentials (sensitive data - be careful with these!)
output "cm_client_id" {
  value       = sitecoreai_cm_client.automation.client_id
  description = "The client ID for the CM automation client"
  sensitive   = true
}

output "cm_client_name" {
  value       = sitecoreai_cm_client.automation.name
  description = "The name of the CM automation client"
}