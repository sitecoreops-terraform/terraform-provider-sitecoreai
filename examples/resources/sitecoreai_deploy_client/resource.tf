# We need proper permissions to create clients, this usually involves CLI authentication
provider "sitecoreai" {
  use_cli = true
}

# Create a Deploy client (note: deploy clients are typically project-level, not environment-specific)
resource "sitecoreai_deploy_client" "deploy" {
  name        = "terraform-deploy-client"
  description = "Deploy client updated again via Terraform"
}

output "deploy_client_name" {
  value       = sitecoreai_deploy_client.deploy.name
  description = "The name of the Deploy client"
}

# Output client credentials (sensitive data - be careful with these!)
output "deploy_client_id" {
  value       = sitecoreai_deploy_client.deploy.client_id
  description = "The client ID for the Deploy client"
}

output "deploy_client_secret" {
  value       = sitecoreai_deploy_client.deploy.client_secret
  description = "The name of the Deploy client"
  sensitive   = true
}
