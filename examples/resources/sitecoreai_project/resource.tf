# Create a new Sitecore project
resource "sitecore_project" "example" {
  name        = "XMC"
  description = "A project for managing our corporate website"
}

# Reference the project ID in outputs
output "project_id" {
  value       = sitecoreai_project.example.id
  description = "The ID of the created project"
}

output "project_name" {
  value       = sitecoreai_project.example.name
  description = "The name of the created project"
}
