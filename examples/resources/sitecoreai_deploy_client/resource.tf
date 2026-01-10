# Use the project data source to get information about a project by name
data "sitecoreai_project" "default" {
  name = "XMC" # Replace with an existing project name
}

# Create a Deploy client (note: deploy clients are typically project-level, not environment-specific)
resource "sitecoreai_deploy_client" "deploy" {
  name        = "terraform-deploy-client"
  description = "Deploy client created via Terraform"
  project_id  = data.sitecoreai_project.default.id
}

# Output client credentials (sensitive data - be careful with these!)
output "deploy_client_id" {
  value       = sitecoreai_deploy_client.deploy.client_id
  description = "The client ID for the Deploy client"
}

output "deploy_client_name" {
  value       = sitecoreai_deploy_client.deploy.name
  description = "The name of the Deploy client"
}