output "existing_project_id" {
  description = "The ID of the existing project found by data source"
  value       = data.sitecoreai_project.main.id
}

output "existing_project_name" {
  description = "The name of the existing project found by data source"
  value       = data.sitecoreai_project.main.name
}

output "new_environment_name" {
  description = "The name of the created environment"
  value       = sitecoreai_cm_environment.main.name
}

output "new_environment_id" {
  description = "The ID of the created environment"
  value       = sitecoreai_cm_environment.main.id
}

output "new_environment_tenant_type" {
  description = "The tenant type of the environment"
  value       = sitecoreai_cm_environment.main.tenant_type
}

output "new_environment_live_context_id" {
  description = "The live context id of the environment"
  value       = sitecoreai_cm_environment.main.live_context_id
}

output "new_environment_preview_context_id" {
  description = "The preview context id of the environment"
  value       = sitecoreai_cm_environment.main.preview_context_id
}

output "existing_environment_editing_secret" {
  description = "The JSS Editing Secret to use for the environment"
  value       = data.sitecoreai_editing_secret.main.secret
}
