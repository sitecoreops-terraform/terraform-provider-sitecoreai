# Use the project data source to get information about a project by name
data "sitecoreai_project" "existing" {
  name = "My Existing Project" # Replace with an existing project name
}

# Output the project information
output "project_id" {
  value       = data.sitecoreai_project.existing.id
  description = "The ID of the existing project"
}

output "project_description" {
  value       = data.sitecoreai_project.existing.description
  description = "The description of the existing project"
}