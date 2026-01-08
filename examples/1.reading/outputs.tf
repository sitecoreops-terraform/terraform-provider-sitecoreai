output "existing_project_id" {
  description = "The ID of the existing project found by data source"
  value       = data.sitecore_project.main.id
}

output "existing_project_name" {
  description = "The name of the existing project found by data source"
  value       = data.sitecore_project.main.name
}
