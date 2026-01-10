# Use the project data source to get information about a project by name
data "sitecoreai_project" "default" {
  name = "XMC" # Replace with an existing project name
}

# Use the environment data source to get information about an environment
data "sitecoreai_environment" "default" {
  project_id = data.sitecoreai_project.default.id
  name       = "production" # Replace with existing environment name
}

# Create an Editing Host Build client
resource "sitecoreai_editing_host_build_client" "editing" {
  name           = "terraform-editing-host-client"
  description    = "Editing Host Build client created via Terraform"
  project_id     = data.sitecoreai_environment.default.project_id
  environment_id = data.sitecoreai_environment.default.id
}

# Output client credentials (sensitive data - be careful with these!)
output "editing_host_client_id" {
  value       = sitecoreai_editing_host_build_client.editing.client_id
  description = "The client ID for the Editing Host Build client"
}

output "editing_host_client_secret" {
  value       = sitecoreai_editing_host_build_client.editing.client_secret
  description = "The name of the Editing Host Build client"
  sensitive   = true
}