# We need proper permissions to create clients, this usually involves CLI authentication
provider "sitecoreai" {
  use_cli = true
}

# Use the project data source to get information about a project by name
data "sitecoreai_project" "default" {
  name = "XMC" # Replace with an existing project name
}

# Use the environment data source to get information about an environment
data "sitecoreai_environment" "default" {
  project_id = data.sitecoreai_project.default.id
  name       = "production" # Replace with existing environment name
}

# Create an Edge client
resource "sitecoreai_edge_client" "edge" {
  name           = "terraform-edge-client"
  description    = "Edge client created via Terraform"
  project_id     = data.sitecoreai_environment.default.project_id
  environment_id = data.sitecoreai_environment.default.id
}

# Output client credentials (sensitive data - be careful with these!)
output "edge_client_id" {
  value       = sitecoreai_edge_client.edge.client_id
  description = "The client ID for the Edge client"
  sensitive   = true
}

output "edge_client_name" {
  value       = sitecoreai_edge_client.edge.name
  description = "The name of the Edge client"
}