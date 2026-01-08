output "existing_project_id" {
  description = "The ID of the existing project found by data source"
  value       = data.sitecore_project.main.id
}

output "existing_project_name" {
  description = "The name of the existing project found by data source"
  value       = data.sitecore_project.main.name
}

output "existing_environment_name" {
  description = "The name of the existing environment found by data source"
  value       = data.sitecore_environment.main.name
}

output "existing_environment_id" {
  description = "The ID of the existing environment found by data source"
  value       = data.sitecore_environment.main.id
}

output "existing_environment_tenant_type" {
  description = "The tenant type of the existing environment"
  value       = data.sitecore_environment.main.tenant_type
}

output "existing_environment_live_context_id" {
  description = "The live context id of the existing environment"
  value       = data.sitecore_environment.main.live_context_id
}

output "existing_environment_preview_context_id" {
  description = "The preview context id of the existing environment"
  value       = data.sitecore_environment.main.preview_context_id
}
