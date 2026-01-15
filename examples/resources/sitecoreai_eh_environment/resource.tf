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

# Output editing host environment details
output "eh_environment_id" {
  value       = sitecoreai_eh_environment.website.id
  description = "The ID of the EH environment"
}

output "eh_environment_name" {
  value       = sitecoreai_eh_environment.website.name
  description = "The name of the EH environment"
}
